package service

import (
	"strings"

	"github.com/go-trellis/trellis/clients"
	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message/proto"
	"github.com/go-trellis/trellis/registry"

	// 注册机
	_ "github.com/go-trellis/trellis/registry/cache"
	_ "github.com/go-trellis/trellis/registry/etcd"

	"github.com/go-trellis/common/errors"
	"github.com/go-trellis/common/formats"
	"github.com/go-trellis/common/logger"
	"github.com/go-trellis/node"
)

// Runner 启动对象
type Runner struct {
	conf   *configure.Project
	logger logger.Logger
	router Router
}

var runner *Runner

// NewRunner 生成启动对象
func NewRunner(cfg *configure.Project, l logger.Logger) (*Runner, error) {

	t := &Runner{
		conf:   cfg,
		router: NewRouter(),

		logger: l,
	}

	if err := t.initRegistries(); err != nil {
		t.logger.Error("init_registries_failed", err)
		return nil, err
	}

	if err := t.newServices(); err != nil {
		t.logger.Error("new_services_failed", err)
		return nil, err
	}
	t.logger.Info("new services ok")

	if err := t.registServices(); err != nil {
		t.logger.Error("regist_services_failed", err)
		return nil, err
	}

	return t, nil
}

func (p *Runner) newServices() error {

	for name, service := range p.conf.Services {
		service.Name = name

		p.logger.Debug("new service", service.String())

		err := p.router.NewService(RouterOptionService(service), RouterOptionLogger(p.logger))
		if err != nil {
			return err
		}

		path := internal.WorkerTrellisPath(service.GetName(), service.GetVersion())

		nm := node.NewDirect(path)
		nm.Add(&node.Node{
			ID:       path,
			Weight:   1,
			Value:    service.String(),
			Metadata: map[string]interface{}{"protocol": proto.Protocol_LOCAL},
		})

		registry.SetNodeManager(path, nm)
	}

	clients.RegistCaller(proto.Protocol_LOCAL, p.router)
	return nil
}

func (p *Runner) registServices() error {
	p.logger.Info("regist service start")
	for name, service := range p.conf.Services {
		service.Name = name

		if service.Registry == nil {
			continue
		}
		p.logger.Info("regist service start", name, service.Registry)

		regConf := &configure.RegistService{
			Name:     service.Name,
			Version:  service.GetVersion(),
			Domain:   service.Registry.Domain,
			Protocol: service.Registry.Protocol,
			Weight:   service.Registry.Weight,
		}
		if err := registry.RegistService(service.Registry.Name, regConf); err != nil {
			p.logger.Error("regist service failed", regConf, err)
			return err
		}
	}

	return nil
}

// Run 启动进程
func (p *Runner) Run() error {
	if err := p.router.Run(); err != nil {
		return err
	}
	return nil
}

// Stop 停止所有的Worker
func (p *Runner) Stop() error {
	var errs errors.Errors
	if err := runner.router.Stop(); err != nil {
		errs = append(errs, err)
	}

	registry.Stop()
	return errs
}

// runRegistries 启动注册器
func (p *Runner) initRegistries() (err error) {

	for name, reg := range p.conf.Registries {
		retryTimes, _ := reg.Options.Int("retry_times")
		rOpts := &registry.RegistOption{
			RegisterType: proto.RegisterType(proto.RegisterType_value[strings.ToUpper(reg.Type)]),
			Endpoint:     reg.Options.Get("endpoint"),
			TTL:          formats.ParseStringTime(reg.Options.Get("ttl")),
			Heartbeat:    formats.ParseStringTime(reg.Options.Get("heartbeat")),
			RetryTimes:   uint32(retryTimes),
			Logger:       p.logger,
		}
		p.logger.Debug("new registry", rOpts)

		if err := registry.NewRegistry(name, rOpts); err != nil {
			p.logger.Error("failed new registry", rOpts.Endpoint, err)
			return err
		}

		for _, wConfig := range reg.Watchers {

			if err = registry.NewRegistryWatcher(name, wConfig); err != nil {
				p.logger.Error("new watcher failed", *wConfig, err)
				return err
			}
		}

		p.logger.Info("initial registry ok", name, reg)

	}
	return nil
}
