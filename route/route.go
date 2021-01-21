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
	"reflect"
	"runtime"

	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/message"
)

var (
	// DefaultLocalRoute default local route
	DefaultLocalRoute = NewLocalRoute()
)

// RegisterComponentFunc regist component funciton into default local route
func RegisterComponentFunc(service *service.Service, fn component.NewComponentFunc) {
	DefaultLocalRoute.RegisterComponentFunc(service, fn)
}

// ListComponents list components in local route
func ListComponents() []component.ComponentDescribe {
	return DefaultLocalRoute.ListComponents()
}

// NewComponent new component
func NewComponent(service *service.Service, alias string, opts ...component.Option) (component.Component, error) {
	return DefaultLocalRoute.NewComponent(service, alias, opts...)
}

type localComponentManager struct {
	components map[string]component.Component

	newComponentFuncs map[string]component.NewComponentFunc

	componentNames []string
}

// NewLocalRoute new default local route manager
func NewLocalRoute() component.Manager {
	return &localComponentManager{
		components:        make(map[string]component.Component),
		newComponentFuncs: make(map[string]component.NewComponentFunc),
	}
}

func (p *localComponentManager) Call(ctx context.Context, msg message.Message) (interface{}, error) {
	cpt, err := p.GetComponent(msg.Service())
	if err != nil {
		return nil, err
	}

	return cpt.Route(msg.Topic())(ctx, msg)
}

// RegisterComponent regist component function
func (p *localComponentManager) RegisterComponentFunc(service *service.Service, fn component.NewComponentFunc) {

	if service == nil || len(service.Name) == 0 {
		panic("component name is empty")
	}

	if fn == nil {
		panic("component fn is nil")
	}

	_, exist := p.newComponentFuncs[service.FullPath()]

	if exist {
		panic(fmt.Sprintf("component already registered: %s", service.FullPath()))
	}

	p.newComponentFuncs[service.FullPath()] = fn
	p.componentNames = append(p.componentNames, service.FullPath())
}

// ListComponents get components
func (p *localComponentManager) ListComponents() []component.ComponentDescribe {

	var desc []component.ComponentDescribe

	for _, name := range p.componentNames {
		cpt := p.components[name]

		desc = append(desc, component.ComponentDescribe{
			Name:         name,
			RegisterFunc: runtime.FuncForPC(reflect.ValueOf(cpt).Pointer()).Name(),
			Component:    cpt,
		})
	}

	return desc
}

// NewComponent new component
func (p *localComponentManager) NewComponent(service *service.Service, alias string, opts ...component.Option) (
	component.Component, error) {
	fn, ok := p.newComponentFuncs[service.FullPath()]

	if !ok {
		return nil, fmt.Errorf("component driver '%s' not exist", service.FullPath())
	}

	cpt, err := fn(alias, opts...)
	if err != nil {
		return nil, err
	}

	p.components[service.FullPath()] = cpt

	return cpt, nil
}

// GetComponent get component
func (p *localComponentManager) GetComponent(s *service.Service) (cpt component.Component, err error) {
	cpt, ok := p.components[s.FullPath()]
	if !ok {
		return nil, errors.New("component is not exists")
	}
	return cpt, nil
}
