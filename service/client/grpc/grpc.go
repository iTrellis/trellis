package grpc

import (
	"context"
	"sync/atomic"

	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/client"
)

type grpcClient struct {
	opts client.Options
	pool *pool
	once atomic.Value
}

func NewClient(opts ...client.Option) client.Client {
	return newClient(opts...)
}

func newClient(opts ...client.Option) client.Client {
	options := client.NewOptions()
	// default content type for grpc
	options.ContentType = "application/grpc+proto"

	for _, o := range opts {
		o(&options)
	}

	rc := &grpcClient{
		opts: options,
	}
	rc.once.Store(false)

	rc.pool = newPool(options.PoolSize, options.PoolTTL, rc.poolMaxIdle(), rc.poolMaxStreams())

	c := client.Client(rc)

	// wrap in reverse
	for i := len(options.Wrappers); i > 0; i-- {
		c = options.Wrappers[i-1](c)
	}

	return c
}

func (p *grpcClient) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return nil
}

func (p *grpcClient) NewMessage(msg interface{}, opts ...client.MessageOption) client.Message {
	return nil
}

func (p *grpcClient) NewRequest(service *service.Service, endpoint string, req interface{}, reqOpts ...client.RequestOption) client.Request {
	return nil
}

func (p *grpcClient) Publish(ctx context.Context, msg client.Message, opts ...client.PublishOption) error {
	return nil
}

func (p *grpcClient) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	return nil, nil
}

func (p *grpcClient) String() string {
	return "grpc"
}

func (p *grpcClient) poolMaxIdle() int {
	if p.opts.Context == nil {
		return DefaultPoolMaxIdle
	}
	v := p.opts.Context.Value(poolMaxIdle{})
	if v == nil {
		return DefaultPoolMaxIdle
	}
	return v.(int)
}

func (p *grpcClient) poolMaxStreams() int {
	if p.opts.Context == nil {
		return DefaultPoolMaxStreams
	}
	v := p.opts.Context.Value(poolMaxStreams{})
	if v == nil {
		return DefaultPoolMaxStreams
	}
	return v.(int)
}
