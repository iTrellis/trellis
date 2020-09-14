package services

import (
	"fmt"

	"github.com/go-trellis/trellis/message"

	"github.com/go-trellis/trellis/service"
)

func init() {
	service.RegistNewServiceFunc("remote_grpc", "v1", NewRemoteGRPC)
}

func NewRemoteGRPC(opts ...service.OptionFunc) (service.Service, error) {
	return &RemoteGRPC{}, nil
}

type RemoteGRPC struct{}

func (p *RemoteGRPC) Start() error {
	fmt.Println("RemoteGRPC Start")
	return nil
}

func (p *RemoteGRPC) Stop() error {
	fmt.Println("RemoteGRPC Stop")
	return nil
}

func (p *RemoteGRPC) Route(topic string) service.HandlerFunc {
	switch topic {
	case "remote":
		return func(msg *message.Message) (interface{}, error) {
			req := &ReqRemote{}
			if err := msg.ToObject(req); err != nil {
				return nil, err
			}
			fmt.Println(string(msg.GetReqBody()))
			return &RespRemote{Msg: fmt.Sprintf("grpc: %s", req.Name)}, nil
		}
	}
	return nil
}
