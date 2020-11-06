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
const (
	DefaultPoolMaxStreams = 20
	DefaultPoolMaxIdle    = 50
)

type GRPCCaller struct {
	opts Options

	pool *pool
}

// TODO remove
func init() {
	RegistCaller(proto.Protocol_GRPC, NewGRPCCaller())
}

func NewGRPCCaller(opts ...OptionFunc) Caller {
	c := &GRPCCaller{}
	for _, o := range opts {
		o(&c.opts)
	}

	conf := c.opts.Options.ToConfig()

	maxIdle := conf.GetInt("pool_max_idle", DefaultPoolMaxIdle)
	maxStreams := conf.GetInt("pool_max_streams", DefaultPoolMaxStreams)

	c.pool = newPool(c.opts.PoolSize, c.opts.PoolTTL, maxIdle, maxStreams)

	return c
}

func (p *GRPCCaller) CallService(ctx context.Context, node *node.Node, msg *message.Message) (interface{}, error) {

	conn, err := grpc.Dial(node.Value, grpc.WithInsecure())
	// poolConn, err := p.pool.getConn(node.Value)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	// var grr error
	// defer p.pool.release(node.Value, poolConn, grr)

	// client := proto.NewRPCServiceClient(poolConn.ClientConn)

	client := proto.NewRPCServiceClient(conn)

	cbResp, err := client.Call(ctx, msg.Payload)

	if err != nil {
		return nil, err
	}
	return cbResp.GetBody(), nil
}
