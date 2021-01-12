package router

import (
	"context"

	"github.com/go-trellis/trellis/service"
	"github.com/go-trellis/trellis/service/message"
	"github.com/go-trellis/trellis/service/registry"
)

// Router router
type Router interface {
	service.LifeCycle

	Register

	Caller
}

type Caller interface {
	Call(context.Context, message.Message) (interface{}, error)
}

// Register router register
type Register interface {
	RegisterRegistry(string, registry.Registry) error
	DeregisterRegistry(string) error

	RegisterService(string, *registry.Service, ...registry.RegisterOption) error
	DeregisterService(string, *registry.Service) error

	Watch(...registry.WatchOption) (registry.Watcher, error)
}
