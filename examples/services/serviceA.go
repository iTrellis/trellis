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

package services

import (
	"fmt"

	"github.com/go-trellis/trellis/clients"
	"github.com/go-trellis/trellis/codec"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/message/proto"
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
			req := &Ping{}
			if err := msg.ToObject(req); err != nil {
				return nil, err
			}
			return Pong{Name: fmt.Sprintf("hello1: %s", req.Name)}, nil
		}
	case "test_remote":
		return func(msg *message.Message) (interface{}, error) {
			req := &Ping{}
			if err := msg.ToObject(req); err != nil {
				return nil, err
			}

			msgTo := msg.Copy()

			msgTo.Service = &proto.Service{Name: "remote_http", Version: "v1"}
			msgTo.Topic = "remote"

			if err := msgTo.SetBody(ReqRemote{Name: req.Name}); err != nil {
				return nil, err
			}

			body, err := clients.CallService(msgTo)
			if err != nil {
				return nil, err
			}

			r := &RespRemote{}

			if err := codec.Unmarshal(msgTo.GetHeader("Content-Type"), body.([]byte), r); err != nil {
				return nil, err
			}
			return Pong{Name: r.Msg}, nil
		}
	case "test_grpc":
		return func(msg *message.Message) (interface{}, error) {
			req := &Ping{}
			if err := msg.ToObject(req); err != nil {
				return nil, err
			}

			msgTo := msg.Copy()

			msgTo.Service = &proto.Service{Name: "remote_grpc", Version: "v1"}
			msgTo.Topic = "remote"

			if err := msgTo.SetBody(ReqRemote{Name: req.Name}); err != nil {
				return nil, err
			}

			body, err := clients.CallService(msgTo)
			if err != nil {
				return nil, err
			}

			r := &RespRemote{}

			if err := codec.Unmarshal(msgTo.GetHeader("Content-Type"), body.([]byte), r); err != nil {
				return nil, err
			}
			return Pong{Name: r.Msg}, nil
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
