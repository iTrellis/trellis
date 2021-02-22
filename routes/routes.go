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

package routes

import (
	"fmt"
	"sync"

	"github.com/iTrellis/common/errors"
	"github.com/iTrellis/node"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/registry"
	"github.com/iTrellis/trellis/service/router"
)

type routes struct {
	sync.RWMutex

	// logger logger.Logger

	registries   map[string]registry.Registry
	watchers     map[string][]registry.Watcher
	serviceNodes map[string]serviceNodes

	components        map[string]component.Component
	newComponentFuncs map[string]component.NewComponentFunc
	componentNames    []string
}

type serviceNodes map[string]*node.Node

// NewRoutes default routes
func NewRoutes() router.Router {
	return &routes{

		registries:   make(map[string]registry.Registry),
		watchers:     make(map[string][]registry.Watcher),
		serviceNodes: make(map[string]serviceNodes),

		components:        make(map[string]component.Component),
		newComponentFuncs: make(map[string]component.NewComponentFunc),
	}
}

func (p *routes) GetServiceNodes(opts ...router.ReadOption) ([]*node.Node, error) {
	options := router.ReadOptions{}

	for _, o := range opts {
		o(&options)
	}
	if options.Service.GetName() == "" {
		return nil, fmt.Errorf("serive is nil")
	}

	mapNodes := p.serviceNodes[options.Service.FullRegistryPath()]

	var nodes []*node.Node
	for _, v := range mapNodes {
		nodes = append(nodes, v)
	}

	return nodes, nil
}

func (p *routes) RegisterRegistry(name string, reg registry.Registry) error {
	p.RLock()
	_, ok := p.registries[name]
	p.RUnlock()
	if ok {
		return errors.New("registry isalready registered")
	}

	p.Lock()
	p.registries[name] = reg
	p.Unlock()
	return nil
}

func (p *routes) DeregisterRegistry(name string) error {

	p.RLock()
	_, ok := p.registries[name]
	if !ok {
		p.RUnlock()
		return errors.New("not found registry")
	}
	watchers, ok := p.watchers[name]
	p.RUnlock()

	p.Lock()
	defer p.Unlock()

	delete(p.watchers, name)
	delete(p.registries, name)
	if ok {

		for _, w := range watchers {
			w.Stop()
		}

	}
	return nil
}

func (p *routes) DeregisterService(name string, s *service.Service, opts ...registry.DeregisterOption) error {
	p.Lock()
	defer p.Unlock()

	reg, ok := p.registries[name]
	if !ok {
		return errors.New("not found registry")
	}

	delete(p.serviceNodes, s.FullRegistryPath())

	return reg.Deregister(s, opts...)
}

func (p *routes) RegisterService(name string, s *service.Service, opts ...registry.RegisterOption) error {
	p.Lock()
	defer p.Unlock()

	reg, ok := p.registries[name]
	if !ok {
		return errors.Newf("regsiter service, not found registry: %s", name)
	}

	return reg.Register(s, opts...)
}

func (p *routes) WatchService(name string, opts ...registry.WatchOption) error {

	p.RLock()
	reg, ok := p.registries[name]
	p.RUnlock()
	if !ok {
		return errors.New("not found registry")
	}
	w, err := reg.Watch(opts...)
	if err != nil {
		return err
	}

	p.Lock()
	p.watchers[name] = append(p.watchers[name], w)
	p.Unlock()

	for {
		result, err := w.Next()
		if err != nil {
			// p.logger.Warn("failed_get_next_node", err)
			// Todo
			return err
		}

		if result.Service == nil {
			continue
		}

		p.RLock()
		sNodes, ok := p.serviceNodes[result.Service.FullRegistryPath()]
		p.RUnlock()
		if !ok {
			sNodes = make(serviceNodes)
		}
		// p.logger.Debugf("watch nodes: %+v, %+v\n", *result, *result.Service.Node)

		switch result.Type {
		case service.EventType_create, service.EventType_update:
			sNodes[result.Service.Node.ID] = result.Service.Node
		case service.EventType_delete:
			delete(sNodes, result.Service.Node.ID)
		}

		p.Lock()
		p.serviceNodes[result.Service.FullRegistryPath()] = sNodes
		p.Unlock()
	}
}

// Start running router
func (p *routes) Start() error {
	p.Lock()
	defer p.Unlock()

	return nil
}

// Stop stop router and the watchers
func (p *routes) Stop() error {
	p.Lock()
	defer p.Unlock()

	for k, watchers := range p.watchers {
		for _, w := range watchers {
			w.Stop()
		}
		delete(p.watchers, k)
	}

	for _, reg := range p.registries {
		reg.Stop()
	}

	p.registries = make(map[string]registry.Registry)
	p.serviceNodes = make(map[string]serviceNodes)

	return nil
}
