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

package service

import (
	"net"

	"github.com/go-trellis/trellis/message/proto"

	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/router"
	"github.com/go-trellis/trellis/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	service.RegistNewServiceFunc("trellis-trans-grpc", NewService)
}

// GrpcService api service
type GrpcService struct {
	debug bool
	opts  service.Options

	Address string
}

// NewService new api service
func NewService(opts ...service.OptionFunc) (service.Service, error) {

	s := &GrpcService{}

	for _, o := range opts {
		o(&s.opts)
	}

	err := s.init()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (p *GrpcService) init() (err error) {
	return
}

// Start start service
func (p *GrpcService) Start() error {
	lis, err := net.Listen("tcp", p.Address)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	proto.RegisterGrpcServcieServer(s, p)
	reflection.Register(s)

	go func() error {
		if err := s.Serve(lis); err != nil {
			return err
		}
		return nil
	}()
	return nil
}

// Stop stop service
func (p *GrpcService) Stop() error {
	return nil
}

// Route 路由
func (p *GrpcService) Route(msg *message.Message) router.HandlerFunc {
	// async中处理callback
	return nil
}

// Call 路由
func (p *GrpcService) Call(msg *proto.Payload) error {
	// async中处理callback
	return nil
}
