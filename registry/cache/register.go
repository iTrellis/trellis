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

package cache

import (
	"fmt"

	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message/proto"
	"github.com/go-trellis/trellis/registry"

	"github.com/go-trellis/cache"
	"github.com/go-trellis/common/logger"
)

const schema = "trellis://"

type worker struct {
	Cache cache.TableCache

	logger logger.Logger
}

// New 生成新对象
func New() registry.Registry {
	return &worker{}
}

func init() {
	registry.Regist(proto.RegisterType_Cache, New)
}

// Init initial register
func (p *worker) Init(opts *registry.RegistOption) (err error) {

	p.logger = opts.Logger

	p.Cache, err = cache.NewTableCache("cache_register")
	if err != nil {
		return err
	}
	return nil
}

func (p *worker) Regist(service *configure.RegistService) error {
	fullpath := internal.WorkerTrellisPath(service.Name, service.Version)
	values, ok := p.Cache.Lookup(fullpath)
	if !ok {
		ss := configure.RegistServices{}
		ss = append(ss, service)
		if ok = p.Cache.Insert(fullpath, ss); !ok {
			return fmt.Errorf("regist service failed")
		}
		return nil
	}

	domainPath := internal.WorkerTrellisDomainPath(service.Name, service.Version, service.Domain)
	conf := values[0].(configure.RegistServices)

	for _, c := range conf {
		cPath := internal.WorkerTrellisDomainPath(c.Name, c.Version, c.Domain)

		if cPath == domainPath {
			return fmt.Errorf("service's domain exists: %s, %s, %s", c.Name, c.Version, c.Domain)
		}
	}
	conf = append(conf, service)
	p.Cache.Insert(fullpath, conf)
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

func (p *worker) Watcher(conf *configure.Watcher) (registry.Watcher, error) {
	w := &watcher{
		registry: p,
		conf:     conf,
		fullpath: conf.Fullpath(),
	}

	return w, nil
}

type watcher struct {
	registry *worker

	conf *configure.Watcher

	fullpath string
}

// Stop 结束进程
func (p *worker) Stop() {
	p.Cache.DeleteObjects()
}

func (p *worker) String() string {
	return "cache"
}

func (p *watcher) Next(ch chan *registry.Result) {
	resp := &registry.Result{
		NodeType: p.conf.LoadBalancing,
		Action:   registry.ActionCreate,
	}
	values, ok := p.registry.Cache.Lookup(p.fullpath)
	if !ok || len(values) == 0 {
		resp.Err = fmt.Errorf("not found service info: %s", p.fullpath)
		ch <- resp
	} else {
		for _, service := range values[0].(configure.RegistServices) {
			resp.Service = service
			ch <- resp
		}
	}

}

func (p *watcher) Stop() {
	if p != nil {
		p = nil
	}
}

func (p *watcher) Fullpath() string {
	return p.fullpath
}
