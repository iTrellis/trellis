package runner

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message/proto"
	"github.com/go-trellis/trellis/registry"
	"github.com/go-trellis/trellis/registry/cache"
	"github.com/go-trellis/trellis/registry/etcd"
	"github.com/go-trellis/trellis/router"
	"github.com/go-trellis/trellis/service"

	"github.com/go-trellis/common/formats"
	"github.com/go-trellis/common/logger"
	"github.com/go-trellis/node"
)

// Runner 启动对象
type Runner struct {
	locker sync.RWMutex

	conf *configure.Project

	originLogger, logger logger.Logger

	router router.Router

	nodeManagers map[string]node.Manager
}

var runner *Runner

// Run 运行
func Run(cfg *configure.Project, l logger.Logger) (err error) {
	err = NewRunner(cfg, l)
	if err != nil {
		return
	}

	if err = runner.Run(); err != nil {
		return err
	}

	return nil
}

// NewRunner 生成启动对象
func NewRunner(cfg *configure.Project, l logger.Logger) error {

	t := &Runner{
		conf:   cfg,
		router: router.NewRouter(),

		originLogger: l,
		logger:       l.WithPrefix("runner"),

		nodeManagers: make(map[string]node.Manager),
	}

	if err := t.registServices(); err != nil {
		return err
	}

	if err := t.initRegistries(); err != nil {
		return err
	}

	runner = t

	return nil
}

// GetService 获取service
func GetService(name, version string, keys ...string) (service.Service, error) {
	path := internal.WorkerPath(internal.SchemaTrellis, name, version)
	runner.locker.RLock()
	nm := runner.nodeManagers[path]
	runner.locker.RUnlock()
	node, ok := nm.NodeFor()
	if !ok {
		return nil, fmt.Errorf("not found service")
	}

	protocol := node.Metadata.Get("protocol")

	switch proto.ServiceType(proto.ServiceType_value[protocol]) {
	case proto.ServiceType_LOCAL:
		return runner.router.GetService(name, version)
	case proto.ServiceType_HTTP:
		// TODO
	case proto.ServiceType_GRPC:
		// TODO
	}
	return nil, fmt.Errorf("not found service")
}

func (p *Runner) registServices() error {
	for name, service := range p.conf.Services {
		service.Name = name
		err := p.router.NewService(
			router.OptionService(service),
			router.OptionLogger(p.originLogger),
		)
		if err != nil {
			return err
		}

		path := internal.WorkerPath(internal.SchemaTrellis, service.GetName(), service.GetVersion())

		nm := node.NewDirect(path)
		nm.Add(&node.Node{
			ID:       path,
			Weight:   1,
			Value:    service.String(),
			Metadata: map[string]interface{}{"protocol": "cache"},
		})

		p.logger.Debug("new service", service.String())

		p.setService(path, nm)
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

// BlockStop 阻断式停止
func BlockStop() error {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	case <-ch:
	}

	return runner.router.Stop()
}

// Stop 停止服务
func Stop() error {
	return runner.router.Stop()
}

// Stop 停止所有的Worker
func (p *Runner) Stop() error {
	return runner.router.Stop()
}

// runRegistries 启动注册器
func (p *Runner) initRegistries() (err error) {

	fnc := func(wr registry.Watcher, conf *configure.Watcher) {
		result, err := wr.Next()
		if err != nil {
			return
		}

		nm := result.Service.ToNodeManager(conf.LoadBalancing)

		path := internal.WorkerPath(internal.SchemaTrellis, conf.GetName(), conf.GetVersion())

		p.setService(path, nm)
	}

	for _, reg := range p.conf.Registries {
		var r registry.Registry
		switch reg.Type {
		case "etcd":
			retryTimes, _ := reg.Options.Int("retry_times")
			rOpts := registry.RegistOption{
				Endpoint:   reg.Options.Get("endpoint"),
				TTL:        formats.ParseStringTime(reg.Options.Get("ttl")),
				Heartbeat:  formats.ParseStringTime(reg.Options.Get("heartbeat")),
				RetryTimes: uint32(retryTimes),
			}
			p.logger.Debug("new etcd registry", rOpts.Endpoint)
			r, err = etcd.NewRegister(rOpts)
			if err != nil {
				p.logger.Error("failed new etcd registry", rOpts.Endpoint, err)
				return err
			}
		case configure.RegistryTypeCache:
			p.logger.Debug("new cache registry")
			r, err = cache.NewRegister()
			if err != nil {
				p.logger.Error("failed new etcd registry")
				return err
			}
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
			p.logger.Debug("new watcher", *wConfig)
			watcher, err := r.Watch(registry.WatchOption{Service: wConfig.Service})
			if err != nil {
				p.logger.Error("new watcher failed", *wConfig, err)
				return err
			}
			switch reg.Type {
			case configure.RegistryTypeCache:
				fnc(watcher, wConfig)
			default:
				go func(wr registry.Watcher, conf *configure.Watcher) {
					for {
						fnc(wr, conf)
					}
				}(watcher, wConfig)
			}
		}
	}
	return nil
}

func (p *Runner) setService(key string, nm node.Manager) {

	p.locker.Lock()
	p.nodeManagers[key] = nm
	p.locker.Unlock()
}
