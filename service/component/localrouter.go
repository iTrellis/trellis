package component

import (
	"github.com/go-trellis/trellis/service"
	"github.com/go-trellis/trellis/service/router"
)

// Manager local router
type Manager interface {
	router.Caller

	RegisterComponentFunc(service *service.Service, fn NewComponentFunc)
	ListComponents() []ComponentDescribe
	NewComponent(service *service.Service, alias string, opts ...Option) (Component, error)
	GetComponent(*service.Service) (Component, error)
}
