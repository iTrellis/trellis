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

package grpc

import (
	"context"
	"net"

	"github.com/go-trellis/trellis/clients"
	"github.com/go-trellis/trellis/codec"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/message/proto"
	"github.com/go-trellis/trellis/service"

	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	service.RegistNewServiceFunc("trellis-grpcserver", "v1", NewService)
}

// Service api service
type Service struct {
	opts service.Options

	Address string
}

// NewService new api service
func NewService(opts ...service.OptionFunc) (service.Service, error) {

	s := &Service{}

	for _, o := range opts {
		o(&s.opts)
	}

	err := s.init()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (p *Service) init() (err error) {
	p.Address = p.opts.Config.GetString("addr")
	return
}

// Start start service
// TODO server options
func (p *Service) Start() error {
	lis, err := net.Listen("tcp", p.Address)
	if err != nil {
		return err
	}

	s := ggrpc.NewServer()
	proto.RegisterRPCServiceServer(s, p)
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
func (p *Service) Stop() error {
	return nil
}

// Route 路由
func (p *Service) Route(string) service.HandlerFunc {
	return nil
}

// Call 路由
func (p *Service) Call(ctx context.Context, payload *proto.Payload) (*proto.Response, error) {
	// async中处理callback
	msg := &message.Message{}
	msg.Payload = payload

	resp, err := clients.CallService(msg)
	if err != nil {
		return nil, err
	}
	cb := &proto.Response{}
	cb.Body, err = codec.Marshal(msg.GetHeader("Content-Type"), resp)
	if err != nil {
		return nil, err
	}

	return cb, nil
}
