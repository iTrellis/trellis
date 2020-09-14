package services

import (
	"fmt"

	"github.com/go-trellis/trellis/message"

	"github.com/go-trellis/trellis/service"
)

func init() {
	service.RegistNewServiceFunc("remote_http", "v1", NewRemoteHTTP)
}

func NewRemoteHTTP(opts ...service.OptionFunc) (service.Service, error) {
	return &RemoteHTTP{}, nil
}

type RemoteHTTP struct{}

func (p *RemoteHTTP) Start() error {
	fmt.Println("RemoteHTTP Start")
	return nil
}

func (p *RemoteHTTP) Stop() error {
	fmt.Println("RemoteHTTP Stop")
	return nil
}

type ReqRemote struct {
	Name string `json:"name"`
}

type RespRemote struct {
	Msg string `json:"message"`
}

func (p *RemoteHTTP) Route(topic string) service.HandlerFunc {
	switch topic {
	case "remote":
		return func(msg *message.Message) (interface{}, error) {
			req := &ReqRemote{}
			if err := msg.ToObject(req); err != nil {
				return nil, err
			}
			fmt.Println(string(msg.GetReqBody()))
			return &RespRemote{Msg: fmt.Sprintf("remote: %s", req.Name)}, nil
		}
	}
	return nil
}
