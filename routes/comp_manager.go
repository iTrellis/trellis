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
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
)

type compManager struct {
	sync.RWMutex

	// logger logger.Logger

	components map[string]component.Component

	newComponentFuncs map[string]component.NewComponentFunc

	componentNames    []string
	startedComponents map[string]bool
}

// NewCompManager new default component manager
func NewCompManager() component.Manager {
	return &compManager{

		components:        make(map[string]component.Component),
		newComponentFuncs: make(map[string]component.NewComponentFunc),
		startedComponents: make(map[string]bool),
	}
}

// RegisterComponent regist component function
func (p *compManager) RegisterComponentFunc(s *service.Service, fn component.NewComponentFunc) {

	if s.GetName() == "" {
		panic("component name is empty")
	}

	if fn == nil {
		panic("component fn is nil")
	}

	p.RLock()
	_, exist := p.newComponentFuncs[s.TrellisPath()]
	p.RUnlock()
	if exist {
		panic(fmt.Sprintf("component already registered: %s", s.TrellisPath()))
	}

	p.Lock()
	p.newComponentFuncs[s.TrellisPath()] = fn
	p.componentNames = append(p.componentNames, s.TrellisPath())
	p.Unlock()
}

// ListComponents get components
func (p *compManager) ListComponents() []component.Describe {

	var descs []component.Describe

	for _, name := range p.componentNames {
		p.RLock()
		cpt := p.components[name]
		started := p.startedComponents[name]
		p.RUnlock()

		desc := component.Describe{
			Name:    name,
			Started: started,
		}

		if cpt != nil {
			desc.RegisterFunc = runtime.FuncForPC(reflect.ValueOf(cpt).Pointer()).Name()
			desc.Component = cpt
		}

		descs = append(descs, desc)
	}

	return descs
}

// NewComponent new component
func (p *compManager) NewComponent(s *service.Service, opts ...component.Option) (
	component.Component, error) {
	p.RLock()
	fn, ok := p.newComponentFuncs[s.TrellisPath()]
	p.RUnlock()
	if !ok {
		return nil, fmt.Errorf("component driver '%s' not exist", s.TrellisPath())
	}

	cpt, err := fn(opts...)
	if err != nil {
		return nil, err
	}

	p.Lock()
	p.components[s.TrellisPath()] = cpt
	p.startedComponents[s.TrellisPath()] = true
	p.Unlock()

	return cpt, nil
}

// GetComponent get component
func (p *compManager) GetComponent(s *service.Service) (cpt component.Component, err error) {
	p.RLock()
	cpt, ok := p.components[s.TrellisPath()]
	p.RUnlock()
	if !ok {
		return nil, errors.New("component is not exists")
	}
	return cpt, nil
}

// type compResp struct {
// 	r   interface{}
// 	err error
// }

// // GetComponent get component
// func (p *compManager) Call(msg message.Message, opts ...component.CallOption) (interface{}, error) {

// 	cpt, err := p.GetComponent(msg.Service())
// 	if err != nil {
// 		return nil, err
// 	}

// 	options := component.CallOptions{}
// 	for _, o := range opts {
// 		o(&options)
// 	}

// 	if options.Timeout == 0 {
// 		options.Timeout = 10 * time.Second
// 	}

// 	h := cpt.Route(msg.Topic())
// 	if h == nil {
// 		return nil, errors.New("not found handler")
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
// 	defer cancel()

// 	ch := make(chan compResp)

// 	go func() {
// 		respH, err := h(msg)
// 		ch <- compResp{
// 			r:   respH,
// 			err: err,
// 		}
// 	}()

// 	select {
// 	case res := <-ch:
// 		// return res
// 		return res.r, res.err
// 	case <-ctx.Done():
// 		return nil, errors.New("component exceed timeout")
// 	}
// }
