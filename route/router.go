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

package route

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/message"
	"github.com/iTrellis/trellis/service/registry"
	"github.com/iTrellis/trellis/service/router"
	"github.com/iTrellis/trellis/service/transport"

	"github.com/iTrellis/node"
)

type routerCache struct {
	sync.RWMutex

	opts Options

	registries map[string]registry.Registry
	watchers   map[string][]registry.Watcher

	nodeManagers map[string]node.Manager

	local component.Manager
}

// NewRouter router constrator
func NewRouter(opts ...Option) router.Router {
	r := &routerCache{

		registries: make(map[string]registry.Registry),
		watchers:   make(map[string][]registry.Watcher),

		local: DefaultLocalRoute,
	}

	for _, o := range opts {
		o(&r.opts)
	}

	return r
}

func (p *routerCache) Call(ctx context.Context, msg message.Message) (interface{}, error) {
	s := msg.Service()

	cpt, err := p.local.GetComponent(s)
	if err == nil && cpt != nil {

		hf := cpt.Route(msg.Topic())
		if hf == nil {
			return nil, fmt.Errorf("handle function: [%s]-[%s] not found", s.FullPath(), msg.Topic())
		}

		return hf(ctx, msg)
	}

	p.RLock()
	nm, ok := p.nodeManagers[s.FullPath()]
	p.RUnlock()
	if !ok {
		return nil, fmt.Errorf("remote service: [%s] not found", s.FullPath())
	}

	nd, ok := nm.NodeFor()
	if !ok {
		return nil, fmt.Errorf("remote node: [%s] not found", s.FullPath())
	}

	return transport.Call(nd, ctx, msg)
}

func (p *routerCache) RegisterRegistry(string, registry.Registry) error {
	return nil
}

func (p *routerCache) DeregisterRegistry(string) error {
	return nil
}

func (p *routerCache) DeregisterService(name string, s *registry.Service) error {
	p.Lock()
	defer p.Unlock()

	reg, ok := p.registries[name]
	if !ok {
		return errors.New("not found registry")
	}

	return reg.Revoke(s)
}

func (p *routerCache) Deregister(name string) {
	p.Lock()
	defer p.Unlock()

	_, ok := p.registries[name]
	if !ok {
		return
	}
	delete(p.registries, name)

	watchers := p.watchers[name]

	for _, w := range watchers {
		w.Stop()
	}

	delete(p.watchers, name)
}

func (p *routerCache) RegisterService(name string, s *registry.Service, opts ...registry.RegisterOption) error {
	p.Lock()
	defer p.Unlock()

	reg, ok := p.registries[name]
	if !ok {
		return errors.New("not found registry")
	}

	return reg.Regist(s, opts...)
}

func (p *routerCache) Register(name string, r registry.Registry) error {
	p.Lock()
	defer p.Unlock()

	_, ok := p.registries[name]
	if ok {
		return errors.New("regitry isalready exists")
	}

	p.registries[name] = r

	return nil
}

func (p *routerCache) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	w := &Watcher{}
	for _, o := range opts {
		o(&w.opts)
	}

	p.watchers[w.opts.Service.FullPath()] = append(p.watchers[w.opts.Service.FullPath()], w)

	return w, nil
}

// Start running cache
func (p *routerCache) Start() error {
	p.Lock()
	defer p.Unlock()

	// TODO errors
	for _, cptDes := range p.local.ListComponents() {
		cptDes.Component.Start()
	}

	for _, ws := range p.watchers {
		for _, w := range ws {
			go func(wter registry.Watcher) {
				// todo watcher for node managers
				r, err := wter.Next()
				println(err, r)
			}(w)
		}
	}
	return nil
}

// Stop stop router and the watchers
func (p *routerCache) Stop() error {
	p.Lock()
	defer p.Unlock()

	for k, watchers := range p.watchers {
		for _, w := range watchers {
			w.Stop()
		}
		delete(p.watchers, k)
	}

	// TODO return errs
	for _, cptDes := range p.local.ListComponents() {
		cptDes.Component.Stop()
	}

	p.registries = make(map[string]registry.Registry)
	p.nodeManagers = make(map[string]node.Manager)

	return nil
}
