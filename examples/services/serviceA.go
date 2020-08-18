package services

import (
	"fmt"

	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/service"
)

func init() {
	service.RegistNewServiceFunc("serviceA", "v1", NewServiceA)
}

func NewServiceA(opts ...service.OptionFunc) (service.Service, error) {
	s := &ServiceA{}

	for _, o := range opts {
		o(&s.opts)
	}

	s.opts.Logger.Info("serviceA_init", "key1", s.opts.Config.Get("key1"))
	s.opts.Logger.Info("serviceA_init", "key2", s.opts.Config.Get("key2"))

	return s, nil
}

type ServiceA struct {
	opts service.Options
}

type Ping struct {
	Name string `json:"name"`
}

type Pong struct {
	Name string `json:"name"`
}

func (p *ServiceA) Route(topic string) service.HandlerFunc {
	switch topic {
	case "test1":
		return func(msg *message.Message) (interface{}, error) {
			fmt.Println(string(msg.ReqBody))
			req := &Ping{}
			if err := msg.ToObject(req); err != nil {
				return nil, err
			}
			return Pong{Name: fmt.Sprintf("hello1: %s", req.Name)}, nil
		}
	case "test2":
		return func(msg *message.Message) (interface{}, error) {
			req := &Ping{}
			if err := msg.ToObject(req); err != nil {
				return nil, err
			}
			return Pong{Name: fmt.Sprintf("hello2: %s", req.Name)}, nil
		}
	}
	return nil
}

func (p *ServiceA) Start() error {
	fmt.Println("serviceA Start")
	return nil
}

func (p *ServiceA) Stop() error {
	fmt.Println("serviceA Stop")
	return nil
}