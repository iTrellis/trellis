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
	service.RegistNewServiceFunc("serviceB", "v1", NewServiceB)
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

func (p *ServiceB) Route(topic string) service.HandlerFunc {
	switch topic {
	case "test":
		return func(msg *message.Message) (interface{}, error) {
			req := &Ping{}
			if err := msg.ToObject(req); err != nil {
				return nil, err
			}
			return Pong{Name: fmt.Sprintf("serviceB test: %s", req.Name)}, nil
		}
	}
	return nil
}
