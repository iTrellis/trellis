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
