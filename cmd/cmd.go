package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/doc"
	"github.com/go-trellis/trellis/route"
	"github.com/go-trellis/trellis/service"
	"github.com/go-trellis/trellis/service/component"
	"github.com/go-trellis/trellis/service/registry"
	"github.com/go-trellis/trellis/service/router"
	"github.com/go-trellis/trellis/version"

	"github.com/go-trellis/common/builder"
	"github.com/go-trellis/common/formats"
	"github.com/go-trellis/common/logger"
	"github.com/go-trellis/config"
	"github.com/go-trellis/node"
	"github.com/urfave/cli/v2"
)

type Cmd interface {
	Init(opts ...Option) error

	App() *cli.App

	AddRegistryFunc(service.RegisterType, registry.NewRegistryFunc)

	component.Manager

	service.LifeCycle

	Run() error
}

type cmd struct {
	opts Options

	app *cli.App

	config configure.Configure

	router router.Router

	component.Manager

	newRegistryFuncs map[service.RegisterType]registry.NewRegistryFunc

	logger logger.Logger

	hiddenVersions []string
}

func (p *cmd) AddRegistryFunc(t service.RegisterType, fn registry.NewRegistryFunc) {
	if _, ok := p.newRegistryFuncs[t]; ok {
		panic(fmt.Errorf("registry function isalready exits"))
	}

	p.newRegistryFuncs[t] = fn
}

func (p *cmd) Start() error {
	builder.Show()

	for _, regConfig := range p.config.Project.Registries {
		fn, ok := p.newRegistryFuncs[regConfig.Type]
		if !ok {
			return errors.New("unsupported registry type")
		}

		reg, err := fn(p.logger)
		if err != nil {
			return err
		}

		err = p.router.RegisterRegistry(regConfig.Name, reg)
		if err != nil {
			return err
		}

		// todo watcher
		// for _, w := range regConfig.Watcher.Services {
		// 	// p.route
		// }
	}

	for _, serviceConf := range p.config.Project.Services {

		_, err := route.DefaultLocalRoute.NewComponent(
			&serviceConf.Service, serviceConf.Alias,
			component.Logger(p.logger),
			component.Router(p.router),
		)
		if err != nil {
			return err
		}

		if serviceConf.Registry == nil {
			continue
		}

		ctx := context.Background()

		for k, v := range serviceConf.Registry.Options {
			ctx = context.WithValue(ctx, k, v)
		}

		regService := &registry.Service{
			Service: serviceConf.Service,
			Nodes:   []*node.Node{serviceConf.ToNode()},
		}

		err = p.router.RegisterService(
			serviceConf.Registry.Name,
			regService,
			registry.RegisterTTL(formats.ParseStringTime(serviceConf.Registry.TTL, 0)),
			registry.RegisterContext(ctx),
		)

		if err != nil {
			return err
		}
	}

	return p.router.Start()
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
	return p.router.Stop()
}

func (p *cmd) App() *cli.App {
	return p.app
}

func (p *cmd) Run() error {
	return nil
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

func New() Cmd {
	cmd := &cmd{
		newRegistryFuncs: DefaultNewRegistryFuncs,

		Manager: route.DefaultLocalRoute,

		logger: logger.NewLogger(),

		app: cli.NewApp(),

		hiddenVersions: DefaultHiddenVersions,
	}

	cmd.router = route.NewRouter(route.Logger(cmd.logger), route.LocalRouter(route.DefaultLocalRoute))

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

				cptsDes := cmd.Manager.ListComponents()
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
