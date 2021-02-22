/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/iTrellis/common/builder"
	"github.com/iTrellis/config"
	"github.com/iTrellis/trellis/configure"
	"github.com/iTrellis/trellis/routes"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/registry"
	"github.com/iTrellis/trellis/version"

	"github.com/urfave/cli/v2"
)

// Cmd command
type Cmd interface {
	// Init initialises options
	// Note: Use Run to parse command line
	Init(opts ...Option) error
	// Options set within this command
	Options() Options

	App() *cli.App

	service.LifeCycle

	BlockRun() error
}

type cmd struct {
	options Options

	app *cli.App

	config configure.Configure

	routesManager routes.Manager
}

func (p *cmd) Options() Options {
	return p.options
}

func (p *cmd) Start() error {
	builder.Show()

	for _, regConfig := range p.config.Project.Registries {
		fn, ok := DefaultNewRegistryFuncs[regConfig.Type]
		if !ok {
			return errors.New("unsupported registry type")
		}

		opts := []registry.Option{}

		opts = append(opts,
			registry.Endpoints(regConfig.Endpoints),
			registry.Timeout(regConfig.Timeout),
			registry.Context(context.Background()),
			// registry.Logger(p.logger),
		)

		reg, err := fn(opts...)
		if err != nil {
			return err
		}

		err = p.routesManager.Router().RegisterRegistry(regConfig.Name, reg)
		if err != nil {
			return err
		}

		for _, w := range regConfig.Watchers {
			go p.routesManager.Router().WatchService(
				regConfig.Name,
				registry.WatchService(w.Service),
				registry.WatchContext(w.Options))
		}
	}

	for _, serviceConf := range p.config.Project.Services {

		if _, err := p.routesManager.CompManager().NewComponent(
			&serviceConf.Service,
			// component.Logger(p.logger),
			component.Caller(p.routesManager),
			component.Config(serviceConf.Options.ToConfig()),
		); err != nil {
			return err
		}

		if serviceConf.Registry == nil {
			continue
		}

		opts := []registry.RegisterOption{}

		opts = append(opts,
			registry.RegisterWeight(serviceConf.Registry.Weight),
			registry.RegisterTTL(serviceConf.Registry.TTL),
			registry.RegisterHeartbeat(serviceConf.Registry.Heartbeat),
		)

		// p.logger.Debug("regist service for registry", serviceConf)

		if err := p.routesManager.Router().RegisterService(
			serviceConf.Registry.Name,
			&serviceConf.Service,
			opts...,
		); err != nil {
			return err
		}
	}

	return p.routesManager.Start()
}

func (p *cmd) Init(opts ...Option) error {
	for _, o := range opts {
		o(&p.options)
	}

	if p.options.configFile != "" {
		reader, err := config.NewSuffixReader(config.ReaderOptionFilename(p.options.configFile))
		if err != nil {
			return err
		}

		err = reader.Read(&p.config)
		if err != nil {
			return err
		}
	} else if p.options.config != nil {
		p.config = *p.options.config
	}

	return nil
}

func (p *cmd) Stop() error {
	// defer p.logger.ClearSubscribers()
	if err := p.routesManager.Stop(); err != nil {
		// p.logger.Error("trellis_routes_stop_failed", err)
		return err
	}

	return nil
}

func (p *cmd) Run() error {
	if err := p.Start(); err != nil {
		return err
	}

	return p.Stop()
}

func (p *cmd) BlockRun() error {
	if err := p.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	case <-ch:
	}

	return p.Stop()
}

func (p *cmd) App() *cli.App {
	return p.app
}

// New new command interface
func New(opts ...Option) Cmd {
	cmd := &cmd{}

	cmd.Init(opts...)

	// cmd.logger = logger.NewLogger()
	cmd.app = cli.NewApp()

	cmd.routesManager = routes.NewManager(
		routes.CompManager(DefaultCompManager),
		routes.WithRouter(routes.NewRoutes()),
	)

	cmd.app.Commands = cli.Commands{
		&cli.Command{
			Name:  "version",
			Usage: "print project version",
			Action: func(ctx *cli.Context) error {
				println(version.Version())
				return nil
			},
		},
		&cli.Command{
			Name:  "build_info",
			Usage: "print project build info",
			Action: func(ctx *cli.Context) error {
				println(version.BuildInfo())
				return nil
			},
		},
		&cli.Command{
			Name:  "list",
			Usage: "list of local informations",
			Subcommands: append([]*cli.Command{},
				&cli.Command{
					Name:  "components",
					Usage: "list of local components",
					Action: func(ctx *cli.Context) error {
						for _, cpt := range cmd.routesManager.CompManager().ListComponents() {
							fmt.Printf("components: %s - started: %t", cpt.Name, cpt.Started)
						}
						return nil
					},
				},
			),
		},
		&cli.Command{
			Name:  "run",
			Usage: "start & stop components",
			Action: func(ctx *cli.Context) error {
				return cmd.Run()
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "config",
					Usage: "config file",
				},
			},
		},
		&cli.Command{
			Name:  "brun",
			Usage: "start & block stop components",
			Action: func(ctx *cli.Context) error {
				return cmd.BlockRun()
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "config",
					Usage: "config file",
				},
			},
		},
	}

	for _, v := range DefaultHiddenVersions {
		if cmd.app.Version != v {
			continue
		}
		cmd.app.HideVersion = true
		break
	}

	return cmd
}
