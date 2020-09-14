package clients

import (
	"context"

	"github.com/go-trellis/node"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/message/proto"
	"google.golang.org/grpc"
)

// TODO
// 1. pool
// 2. dial options

type GRPCCaller struct {
}

func init() {
	RegistCaller(proto.Protocol_GRPC, NewGRPCCaller())
}

func NewGRPCCaller() Caller {
	return &GRPCCaller{}
}

func (p *GRPCCaller) CallService(node *node.Node, msg *message.Message) (interface{}, error) {

	conn, err := grpc.Dial(node.Value, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := proto.NewRPCServiceClient(conn)

	cb, err := client.Call(context.Background(), msg.Payload)
	if err != nil {
		return nil, err
	}
	return cb.GetBody(), nil
}
