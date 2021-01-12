package route

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"

	"github.com/go-trellis/trellis/service"
	"github.com/go-trellis/trellis/service/component"
	"github.com/go-trellis/trellis/service/message"
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
