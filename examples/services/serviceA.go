package services

import (
	"fmt"

	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/router"
	"github.com/go-trellis/trellis/service"
)

func init() {
	service.RegistNewServiceFunc("serviceA", NewServiceA)
}

func NewServiceA(optFuncs ...service.OptionFunc) (service.Service, error) {
	s := &ServiceA{}

	for _, o := range optFuncs {
		o(&s.opts)
	}

	fmt.Println(s.opts.Config.GetString("key1"))
	fmt.Println(s.opts.Config.GetString("key2"))
	return s, nil
}

type ServiceA struct {
	opts service.Options
}

type ResponseA struct {
	Name string `json:"name"`
}

func (p *ServiceA) Route(msg *message.Message) router.HandlerFunc {
	fmt.Println(msg)
	resp := ResponseA{}
	switch msg.GetTopic() {
	case "test1":
		return func(msg *message.Message) error {
			resp.Name = "test1"
			msg.SetBody(resp)
			return nil
		}
	case "test2":
		return func(msg *message.Message) error {
			resp.Name = "test2"
			msg.SetBody(resp)
			return nil
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
