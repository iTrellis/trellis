package services

import (
	"fmt"

	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/router"
	"github.com/go-trellis/trellis/service"
)

func init() {
	service.RegistNewServiceFunc("serviceB", NewServiceB)
}

func NewServiceB(opts ...service.OptionFunc) (service.Service, error) {
	return &ServiceB{}, nil
}

type ServiceB struct{}

func (p *ServiceB) Start() error {
	fmt.Println("serviceB Start")
	return nil
}

func (p *ServiceB) Stop() error {
	fmt.Println("serviceB Stop")
	return nil
}

func (p *ServiceB) Route(*message.Message) router.HandlerFunc {
	return nil
}
