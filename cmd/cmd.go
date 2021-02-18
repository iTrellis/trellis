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
	"github.com/iTrellis/common/formats"
	"github.com/iTrellis/common/logger"
	"github.com/iTrellis/config"
	"github.com/iTrellis/trellis/configure"
	"github.com/iTrellis/trellis/doc"
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

	AddRegistryFunc(service.RegisterType, registry.NewRegistryFunc)

	GetRoutesManager() routes.Manager

	service.LifeCycle

	Run() error
}

type cmd struct {
	opts Options

	app *cli.App

	config configure.Configure

	routesManager routes.Manager

	newRegistryFuncs map[service.RegisterType]registry.NewRegistryFunc

	logger logger.Logger

	writers []logger.Writer

	hiddenVersions []string
}

func (p *cmd) AddRegistryFunc(t service.RegisterType, fn registry.NewRegistryFunc) {
	if _, ok := p.newRegistryFuncs[t]; ok {
		panic(fmt.Errorf("registry function isalready exits"))
	}

	p.newRegistryFuncs[t] = fn
}

func (p *cmd) Options() Options {
	return p.opts
}

func (p *cmd) Start() error {
	builder.Show()

	for _, regConfig := range p.config.Project.Registries {
		fn, ok := p.newRegistryFuncs[regConfig.Type]
		if !ok {
			return errors.New("unsupported registry type")
		}

		opts := []registry.Option{}

		// todo ctx
		ctx := context.Background()

		opts = append(opts,
			registry.Adds(regConfig.Address),
			registry.Timeout(regConfig.Timeout),
			registry.Context(ctx),
			registry.Logger(p.logger),
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
			p.routesManager.Router().WatchService(regConfig.Name, registry.WatchService(w.Service), registry.WatchContext(w.Options))
		}
	}

	for _, serviceConf := range p.config.Project.Services {

		node := serviceConf.ToNode()
		_, err := p.routesManager.CompManager().NewComponent(
			&serviceConf.Service, serviceConf.Alias,
			component.Logger(p.logger),
			component.Caller(p.routesManager),
			component.Config(node.Metadata.ToConfig()),
		)
		if err != nil {
			return err
		}

		if serviceConf.Registry == nil {
			continue
		}

		opts := []registry.RegisterOption{}
		// ctx := context.Background()
		// for k, v := range serviceConf.Registry.Options {
		// 	ctx = context.WithValue(ctx, k, v)
		// }

		opts = append(opts, registry.RegisterTTL(formats.ParseStringTime(serviceConf.Registry.TTL, 0)))
		// opts = append(opts, registry.RegisterContext(ctx))

		p.logger.Debug("regist service for registry", serviceConf)

		if err = p.routesManager.Router().RegisterService(
			serviceConf.Registry.Name,
			&registry.Service{
				Service: serviceConf.Service,
				// Nodes:   []*node.Node{serviceConf.ToNode()},
				Node: node,
			},
			opts...,
		); err != nil {
			return err
		}
	}

	return p.routesManager.Start()
}

func (p *cmd) Init(opts ...Option) error {
	for _, o := range opts {
		o(&p.opts)
	}

	reader, err := config.NewSuffixReader(config.ReaderOptionFilename(p.opts.ConfigFile))
	if err != nil {
		return err
	}

	err = reader.Read(&p.config)
	if err != nil {
		return err
	}

	return nil
}

func (p *cmd) Stop() error {
	if err := p.routesManager.Stop(); err != nil {
		return err
	}

	p.logger.ClearSubscribers()

	return nil
}

func (p *cmd) App() *cli.App {
	return p.app
}

func (p *cmd) Run() error {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	case <-ch:
	}

	fmt.Println("exit")
	return p.Stop()
}

func (p *cmd) document(ctx *cli.Context) (err error) {

	name := ctx.String("name")

	if len(name) == 0 {
		return
	}

	documenter, exist := doc.GetDocumenter(name)
	if !exist {
		err = fmt.Errorf("documenter of %s not exist", name)
		return
	}

	document := documenter.Document()

	docStr, err := document.JSON()

	if err != nil {
		return
	}

	fmt.Println(docStr)

	return
}

func (p *cmd) GetRoutesManager() routes.Manager {
	return p.routesManager
}

func New() Cmd {
	cmd := &cmd{
		newRegistryFuncs: DefaultNewRegistryFuncs,

		logger: logger.NewLogger(),

		app: cli.NewApp(),

		hiddenVersions: DefaultHiddenVersions,
	}

	logW, err := logger.ChanWriter(cmd.logger, logger.ChanWiterLevel(logger.DebugLevel))
	if err != nil {
		panic(err)
	}

	cmd.routesManager = routes.NewManager(
		routes.Logger(cmd.logger),
		routes.CompManager(DefaultCompManager),
		routes.WithRouter(routes.NewRoutes(cmd.logger)),
	)

	cmd.writers = append(cmd.writers, logW)

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
			Name:  "components",
			Usage: "list of local components",
			Action: func(ctx *cli.Context) error {
				cptsDes := cmd.routesManager.CompManager().ListComponents()
				for _, cpt := range cptsDes {
					fmt.Printf("%s: %s", cpt.Name, cpt.RegisterFunc)
				}
				return nil
			},
		},
		&cli.Command{
			Name:  "run",
			Usage: "run components",
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
			Name:   "document",
			Usage:  "print document",
			Action: cmd.document,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "name",
					Usage: "component name",
				},
			},
		},
	}

	for _, v := range cmd.hiddenVersions {
		if cmd.app.Version != v {
			continue
		}
		cmd.app.HideVersion = true
		break
	}

	return cmd
}
