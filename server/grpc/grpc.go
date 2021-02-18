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
	cmd.DefaultCompManager.RegisterComponentFunc(&service.Service{Name: "trellis-grpcserver", Version: "v1"}, NewService)
}

// Service api service
type Service struct {
	alias string

	opts component.Options

	grpcServer *grpc.Server

	Address string
}

// NewService new api service
func NewService(alias string, opts ...component.Option) (component.Component, error) {

	s := Service{alias: alias}

	for _, o := range opts {
		o(&s.opts)
	}

	err := s.init()
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (p *Service) Alias() string {
	return p.alias
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
func (p *Service) Route(string) component.Handler {
	return nil
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
