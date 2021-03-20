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

	"github.com/iTrellis/trellis/cmd"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/message"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	cmd.DefaultCompManager.RegisterComponentFunc(&service.Service{Name: "trellis-server-grpc", Version: "v1"}, NewService)
}

// Service api service
type Service struct {
	opts component.Options

	grpcServer *grpc.Server

	Address string
}

// NewService new api service
func NewService(opts ...component.Option) (component.Component, error) {

	s := Service{}

	for _, o := range opts {
		o(&s.opts)
	}

	err := s.init()
	if err != nil {
		return nil, err
	}

	return &s, nil
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

	s := grpc.NewServer()
	RegisterClientServer(s, p)
	reflection.Register(s)

	go func() error {
		if err := s.Serve(lis); err != nil {
			return err
		}
		return nil
	}()

	p.grpcServer = s
	return nil
}

// Stop stop service
func (p *Service) Stop() error {
	p.grpcServer.Stop()
	return nil
}

// Route 路由
func (p *Service) Route(_ message.Message) (interface{}, error) {
	return nil, nil
}

// Call 路由
func (p *Service) Call(ctx context.Context, req *message.Request) (*message.Response, error) {
	// callResp, err := p.opts.Router.Call(ctx, req)
	// if err != nil {
	// 	return nil, err
	// }

	// msg := &message.Response{
	// 	Header: make(map[string]string),
	// }
	// t := reflect.TypeOf(callResp)
	// switch t.Kind() {
	// case reflect.String:
	// 	msg.Body = []byte(callResp.(string))
	// 	msg.Header["Content-Type"] = "text/plain"
	// case reflect.Ptr:
	// 	msg.Body, _ = json.Marshal(callResp)
	// 	msg.Header["Content-Type"] = "application/json"
	// default:
	// 	return nil, errors.Newf("unsupported type: %s", t.String())
	// }

	// return msg, nil
	return nil, nil
}

// Publish 路由
func (p *Service) Publish(context.Context, *message.Payload) (*message.Payload, error) {

	return nil, nil
}

// Stream 路由
func (p *Service) Stream(Client_StreamServer) error {

	return nil
}
