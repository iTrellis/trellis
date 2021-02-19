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

package memory

import (
	"sync"
	"time"

	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/registry"

	"github.com/google/uuid"
	"github.com/iTrellis/node"
)

var (
	sendEventTime = 10 * time.Millisecond
)

type memory struct {
	id string

	sync.RWMutex

	options registry.Options

	// map[serviceFullName]map[version]*registry.Service
	services map[string]map[string]*registry.Service
	watchers map[string]*Watcher
}

// NewRegistry 生成新对象
func NewRegistry(opts ...registry.Option) (registry.Registry, error) {
	options := registry.Options{}
	for _, o := range opts {
		o(&options)
	}

	r := &memory{
		id: uuid.New().String(),

		options: options,

		// domain/service version
		services: make(map[string]map[string]*registry.Service),
		watchers: make(map[string]*Watcher),
	}

	return r, nil
}

// Init initial register
func (p *memory) Init(fs ...registry.Option) (err error) {

	for _, f := range fs {
		f(&p.options)
	}

	return nil
}

func (p *memory) Options() registry.Options {
	return p.options
}

func (p *memory) Register(s *service.Service, ofs ...registry.RegisterOption) error {
	p.Lock()
	defer p.Unlock()
	serviceName := s.FullRegistryPath()
	nodes, ok := p.services[serviceName]
	if !ok || nodes == nil {
		nodes = make(map[string]*registry.Service)
	}

	regService := &registry.Service{
		Service: *s,

		Node: &node.Node{
			ID:     s.ID(uuid.NewString()),
			Weight: 1,
			Value:  p.options.ServerAddr,
		},
	}
	p.options.Logger.Debugf("Registry (memory) added new service: %+v", *regService)

	nodes[s.GetVersion()] = regService

	go p.sendEvent(&registry.Result{
		ID:        p.id,
		Timestamp: time.Now(),
		Type:      service.EventType_update,
		Service:   regService})

	return nil
}

func (p *memory) Deregister(s *service.Service, ofs ...registry.DeregisterOption) error {
	p.Lock()
	defer p.Unlock()
	serviceName := s.FullRegistryPath()
	nodes, ok := p.services[serviceName]
	if !ok {
		return nil
	}

	if _, ok := nodes[s.GetVersion()]; ok {
		p.options.Logger.Debugf("Registry (memory) removed service' version: %+v", s)
		delete(p.services[serviceName], s.GetVersion())
	}

	regServices, ok := p.services[serviceName]
	if ok {
		p.options.Logger.Debugf("Registry (memory) removed service: %+v", s)
		delete(p.services, serviceName)
	}

	for _, v := range regServices {
		item := *v
		go p.sendEvent(&registry.Result{
			ID:        p.id,
			Timestamp: time.Now(),
			Type:      service.EventType_delete,
			Service:   &item})
	}

	return nil
}

func (p *memory) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	var wo registry.WatchOptions
	for _, o := range opts {
		o(&wo)
	}

	w := &Watcher{
		id:   uuid.New().String(),
		exit: make(chan bool),
		res:  make(chan *registry.Result),
		wo:   wo,
	}

	p.Lock()
	p.watchers[w.id] = w
	p.Unlock()

	return w, nil
}

func (p *memory) Stop() error {
	return nil
}

func (p *memory) ID() string {
	return p.id
}

func (p *memory) String() string {
	return service.RegisterType_memory.String()
}

func (p *memory) sendEvent(r *registry.Result) {
	p.RLock()
	watchers := make([]*Watcher, 0, len(p.watchers))
	for _, w := range p.watchers {
		watchers = append(watchers, w)
	}
	p.RUnlock()

	for _, w := range watchers {
		select {
		case <-w.exit:
			p.Lock()
			delete(p.watchers, w.id)
			p.Unlock()
		default:
			select {
			case w.res <- r:
			case <-time.After(sendEventTime):
			}
		}
	}
}
