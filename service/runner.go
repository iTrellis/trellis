package service

import (
	"fmt"
	"sync"

	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/internal"
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
	locker sync.RWMutex

	conf *configure.Project

	logger logger.Logger

	router Router

	nodeManagers map[string]node.Manager

	registries map[string]registry.Registry
}

var runner *Runner

// NewRunner 生成启动对象
func NewRunner(cfg *configure.Project, l logger.Logger) (*Runner, error) {

	t := &Runner{
		conf:   cfg,
		router: NewRouter(),

		logger: l,

		nodeManagers: make(map[string]node.Manager),
		registries:   make(map[string]registry.Registry),
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
			Metadata: map[string]interface{}{"protocol": "cache"},
		})

		p.setNodeManager(path, nm)
	}
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

		r, ok := p.registries[service.Registry.Name]
		if !ok {
			return fmt.Errorf("not found service(%s)'s registry (%s)", service.Name, service.Registry.Name)
		}
		regConf := &configure.RegistService{
			Name:     service.Name,
			Version:  service.GetVersion(),
			Domain:   service.Registry.Domain,
			Protocol: service.Registry.Protocol,
			Weight:   service.Registry.Weight,
			Metadata: service.Registry.Metadata,
		}
		if err := r.Regist(regConf); err != nil {
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

	for _, r := range p.registries {
		r.Stop()
	}

	return errs
}

// runRegistries 启动注册器
func (p *Runner) initRegistries() (err error) {

	fnc := func(w registry.Watcher) error {
		ch := make(chan *registry.Result)
		go w.Next(ch)

		for {
			result := <-ch

			if result.Err != nil {
				p.logger.Error("get registry result failed", result, result.Err)
				continue
			}

			path := internal.WorkerTrellisPath(result.Service.Name, result.Service.Version)

			nm, ok := p.getNodeManager(path)

			if !ok {
				nm = node.New(result.NodeType, path)
			}

			nd := result.ToNode()
			p.logger.Info("get registry node", result, nd)
			switch result.Action {
			case registry.ActionCreate, registry.ActionUpdate:
				nm.Add(nd)
			case registry.ActionDelete:
				nm.RemoveByID(nd.ID)
			default:
			}

			p.setNodeManager(path, nm)

			nm.PrintNodes()
		}
	}

	for name, reg := range p.conf.Registries {

		retryTimes, _ := reg.Options.Int("retry_times")
		rOpts := &registry.RegistOption{
			Endpoint:   reg.Options.Get("endpoint"),
			TTL:        formats.ParseStringTime(reg.Options.Get("ttl")),
			Heartbeat:  formats.ParseStringTime(reg.Options.Get("heartbeat")),
			RetryTimes: uint32(retryTimes),
		}
		p.logger.Debug("new registry", rOpts)

		fn, err := registry.GetNewRegistryFunc(reg.Type)
		if err != nil {
			p.logger.Error("failed get registry func", reg.Type, err)
			return err
		}

		r := fn()

		if err := r.Init(rOpts, p.logger); err != nil {
			p.logger.Error("failed new registry", rOpts.Endpoint, err)
			return err
		}

		// 注册服务
		for _, serv := range reg.Services {
			p.logger.Debug("regist service", serv.String())
			if err = r.Regist(serv); err != nil {
				p.logger.Error("regist service failed", serv.String(), err)
				return err
			}
		}

		for _, wConfig := range reg.Watchers {
			p.logger.Debug("new watcher", wConfig)
			watcher, err := r.Watcher(wConfig)
			if err != nil {
				p.logger.Error("new watcher failed", *wConfig, err)
				return err
			}

			go fnc(watcher)
		}

		p.logger.Info("initical registry ok", name, reg)

		p.registries[name] = r
	}
	return nil
}

func (p *Runner) setNodeManager(key string, nm node.Manager) {
	p.locker.Lock()
	p.nodeManagers[key] = nm
	p.locker.Unlock()
}

func (p *Runner) delNodeManager(key string) {
	p.locker.Lock()
	delete(p.nodeManagers, key)
	p.locker.Unlock()
}

func (p *Runner) getNodeManager(key string) (node.Manager, bool) {
	p.locker.RLock()
	nm, ok := p.nodeManagers[key]
	p.locker.RUnlock()
	return nm, ok
}
