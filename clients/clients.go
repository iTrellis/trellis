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
	"fmt"
	"time"

	"github.com/go-trellis/node"

	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/message/proto"
	"github.com/go-trellis/trellis/registry"
)

const (
	DefaultTimeout = 10 * time.Second
)

// Caller 客户端请求对象
type Caller interface {
	CallService(ctx context.Context, nd *node.Node, msg *message.Message) (resp interface{}, err error)
}

var mapCallers = map[proto.Protocol]Caller{}

// RegistCaller 注册caller
func RegistCaller(protocol proto.Protocol, caller Caller) {
	mapCallers[protocol] = caller
}

// CallService 请求服务
func CallService(msg *message.Message, keys ...string) (resp interface{}, err error) {

	path := internal.WorkerTrellisPath(msg.GetService().GetName(), msg.GetService().GetVersion())
	nm, ok := registry.GetNodeManager(path)
	if !ok {
		return nil, fmt.Errorf("not found service's node manager: %+v", msg.GetService())
	}
	node, ok := nm.NodeFor(keys...)
	if !ok {
		return nil, fmt.Errorf("not found service's node")
	}

	mdConfig := node.Metadata.ToConfig()

	protocol := mdConfig.GetInt("protocol")

	c, ok := mapCallers[proto.Protocol(protocol)]
	if !ok {
		return nil, fmt.Errorf("unknown caller: %d", protocol)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	return c.CallService(ctx, node, msg)
}
