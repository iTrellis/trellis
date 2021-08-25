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

	"github.com/iTrellis/trellis/configure"
	"github.com/iTrellis/trellis/routes"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/registry"
	"github.com/iTrellis/trellis/version"

	"github.com/iTrellis/common/builder"
	"github.com/iTrellis/common/logger"
	"github.com/iTrellis/config"
	"github.com/iTrellis/node"
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

	config config.Config

	routesManager routes.Manager

	registries map[string]registry.Registry

	logger logger.Logger
}

func (p *cmd) Options() Options {
	return p.options
}

func (p *cmd) Start() error {
	if p.config == nil {
		return nil
	}

	registriesConfig := p.config.GetValuesConfig("project.registries")

	for _, rKey := range registriesConfig.GetKeys() {

		regConfig := &configure.Registry{}

		if err := registriesConfig.ToObject(rKey, regConfig); err != nil {
			return err
		}

		fn, ok := DefaultNewRegistryFuncs[regConfig.Type]
		if !ok {
			return errors.New("unsupported registry type")
		}

		opts := []registry.Option{}

		opts = append(opts,
			registry.Endpoints(regConfig.Endpoints),
			registry.ServerAddr(regConfig.ServerAddr),
			registry.Timeout(regConfig.Timeout),
			registry.RetryTimes(regConfig.RetryTimes),
			registry.Context(context.Background()),
			registry.Logger(p.logger.With("registry", regConfig.Name)),
		)

		// option secure

		p.logger.Debug("new_registry", "name", regConfig.Name, "address", regConfig.ServerAddr)

		reg, err := fn(opts...)
		if err != nil {
			return err
		}

		p.registries[rKey] = reg

		for _, w := range regConfig.Watchers {
			p.logger.Debug("new_registry_watcher", "name", regConfig.Name, "address", regConfig.ServerAddr,
				"watch_service", w.Service.FullRegistryPath())
			rCpt, err := routes.NewRemoteComponent(node.NodeTypeRandom, reg,
				registry.WatchService(w.Service),
				registry.WatchLogger(
					p.logger.With("registry", regConfig.Name, "watcher", w.Service.FullRegistryPath())),
			)
			if err != nil {
				return err
			}

			rCpt.Init(component.Caller(p.routesManager),
				component.Config(w.Options.ToConfig()),
				component.Logger(p.logger.With("remote_component", w.Service.TrellisPath())))

			if err = p.routesManager.CompManager().RegisterComponent(&w.Service, rCpt); err != nil {
				return err
			}
		}
	}

	servicesConfig := p.config.GetValuesConfig("project.services")

	for _, sKey := range servicesConfig.GetKeys() {

		serviceConf := &configure.Service{}

		if err := servicesConfig.ToObject(sKey, serviceConf); err != nil {
			return err
		}

		p.logger.Debug("new_component", "component", serviceConf.Service.TrellisPath())
		if _, err := p.routesManager.CompManager().NewComponent(
			&serviceConf.Service,
			component.Caller(p.routesManager),
			component.Config(serviceConf.Options.ToConfig()),
			component.Logger(p.logger.With("component", serviceConf.Service.TrellisPath())),
		); err != nil {
			p.logger.Error("new_component", "component", serviceConf.Service.TrellisPath(), "err", err.Error())
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

		p.logger.Debug("regist service for registry", "config", serviceConf)

		reg, ok := p.registries[serviceConf.Registry.Name]
		if !ok {
			return fmt.Errorf("not found registry: %s", serviceConf.Registry.Name)
		}

		err := reg.Register(&serviceConf.Service, opts...)

		if err != nil {
			return err
		}
	}

	return p.routesManager.Start()
}

func (p *cmd) Init(opts ...Option) (err error) {

	for _, o := range opts {
		o(&p.options)
	}

	if p.options.configFile != "" {
		p.config, err = config.NewConfig(p.options.configFile)
		if err != nil {
			return
		}
	} else if p.options.config != nil {
		p.config, err = config.NewConfigOptions(config.OptionStruct(config.ReaderTypeYAML, p.options.config))
		if err != nil {
			return err
		}
	} else {
		return nil
	}

	err = p.initLogger()
	if err != nil {
		return err
	}

	p.logger.Debug("start_initial", "cammand_initial", "new manager")

	p.routesManager = routes.NewManager(
		routes.CompManager(DefaultCompManager),
		routes.Logger(p.logger.With("component", "routes_manager")),
	)
	return
}

func (p *cmd) initLogger() error {
	var loggerConfig logger.LogConfig
	err := p.config.ToObject("project.logger", &loggerConfig)
	if err != nil {
		return err
	}

	var loggerOptions []logger.Option

	if loggerConfig.Caller {
		loggerOptions = append(loggerOptions, logger.Caller())
		loggerOptions = append(loggerOptions, logger.CallerSkip(1))
	}

	if loggerConfig.EncoderConfig != nil {
		loggerOptions = append(loggerOptions, logger.EncoderConfig(loggerConfig.EncoderConfig))
	}

	if loggerConfig.StackTrace {
		loggerOptions = append(loggerOptions, logger.StackTrace())
	}
	if loggerConfig.Encoding == "" {
		loggerConfig.Encoding = "json"
	}
	loggerOptions = append(loggerOptions, logger.Encoding(loggerConfig.Encoding))
	loggerOptions = append(loggerOptions, logger.LogLevel(loggerConfig.Level))
	loggerOptions = append(loggerOptions, logger.LogFileOptions(&loggerConfig.FileOptions))

	zLogger, err := logger.NewLogger(loggerOptions...)
	if err != nil {
		return err
	}

	p.logger = zLogger
	return nil
}

func (p *cmd) Stop() error {
	if p.routesManager == nil {
		return nil
	}
	if err := p.routesManager.Stop(); err != nil {
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

	<-ch

	return p.Stop()
}

func (p *cmd) App() *cli.App {
	return p.app
}

// New new command interface
func New(opts ...Option) (Cmd, error) {
	builder.Show()

	cmd := &cmd{
		app:        cli.NewApp(),
		registries: make(map[string]registry.Registry),
	}

	if err := cmd.Init(opts...); err != nil {
		return nil, err
	}

	cmd.app.Commands = cli.Commands{
		&cli.Command{
			Name:  "version",
			Usage: "print project version",
			Action: func(ctx *cli.Context) error {
				fmt.Println(version.Version())
				return nil
			},
		},
		&cli.Command{
			Name:  "build_info",
			Usage: "print project build info",
			Action: func(ctx *cli.Context) error {
				fmt.Println(version.BuildInfo())
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
							fmt.Printf("components: %s - started: %t\n", cpt.Name, cpt.Started)
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
				err := cmd.Init(ConfigFile(ctx.String("config")))
				if err != nil {
					return err
				}
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
				err := cmd.Init(ConfigFile(ctx.String("config")))
				if err != nil {
					return err
				}
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

	return cmd, nil
}
