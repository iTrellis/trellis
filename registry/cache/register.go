package cache

import (
	"fmt"

	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/registry"

	"github.com/go-trellis/cache"
)

const schema = "trellis://"

type worker struct {
	Cache cache.TableCache
}

func NewRegister() (register registry.Registry, err error) {
	w := &worker{}

	w.Cache, err = cache.NewTableCache("cache_register")
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (p *worker) Regist(service *configure.RegistService) error {
	servicePath := internal.WorkerPath(internal.SchemaTrellis, service.Name, service.Version)
	values, ok := p.Cache.Lookup(servicePath)
	if !ok {
		ss := configure.RegistServices{}
		ss = append(ss, service)
		if ok = p.Cache.Insert(servicePath, ss); !ok {
			return fmt.Errorf("regist service failed")
		}
		return nil
	}

	domainPath := internal.WorkerDomainPath(internal.SchemaTrellis, service.Name, service.Version, service.Domain)
	conf := values[0].(configure.RegistServices)

	for _, c := range conf {
		cPath := internal.WorkerDomainPath(internal.SchemaTrellis, c.Name, c.Version, c.Domain)

		if cPath == domainPath {
			return fmt.Errorf("service's domain exists: %s, %s, %s", c.Name, c.Version, c.Domain)
		}
	}

	return nil
}

func (p *worker) Revoke(service *configure.RegistService) error {
	servicePath := internal.WorkerPath(internal.SchemaTrellis, service.Name, service.Version)
	values, ok := p.Cache.Lookup(servicePath)
	if !ok {
		return nil
	}
	domainPath := internal.WorkerDomainPath(internal.SchemaTrellis, service.Name, service.Version, service.Domain)
	conf := values[0].(configure.RegistServices)
	for i, c := range conf {
		cPath := internal.WorkerDomainPath(internal.SchemaTrellis, c.Name, c.Version, c.Domain)
		if cPath == domainPath {
			conf = append(conf[:i], conf[i+1:]...)
			break
		}
	}

	p.Cache.Insert(servicePath, conf)
	return nil
}

func (p *worker) Watch(registry.WatchOption) (registry.Watcher, error) {
	return &watcher{
		registry: p,
	}, nil
}

type watcher struct {
	key      string
	registry *worker
}

// Stop 结束进程
func (p *worker) Stop() {
	p.Cache.DeleteObjects()
}

func (p *worker) String() string {
	return "cache"
}

func (p *watcher) Next() (*registry.Result, error) {
	values, ok := p.registry.Cache.Lookup(p.key)
	if !ok {
		return nil, fmt.Errorf("not found service info: %s", p.key)
	}

	services := values[0].(configure.RegistServices)
	return &registry.Result{Service: services}, nil
}

func (p *watcher) Stop() {
	if p != nil {
		p = nil
	}
}
