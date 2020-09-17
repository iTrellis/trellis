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
	"fmt"

	"github.com/go-trellis/node"

	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/message/proto"
	"github.com/go-trellis/trellis/registry"
)

// Caller 客户端请求对象
type Caller interface {
	CallService(nd *node.Node, msg *message.Message) (resp interface{}, err error)
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
		return nil, fmt.Errorf("not found service's node manager")
	}
	node, ok := nm.NodeFor(keys...)
	if !ok {
		return nil, fmt.Errorf("not found service's node")
	}

	protocol, err := node.Metadata.Int("protocol")
	if err != nil {
		return nil, err
	}
	c, ok := mapCallers[proto.Protocol(protocol)]
	if !ok {
		return nil, fmt.Errorf("unknown caller: %d", protocol)
	}

	return c.CallService(node, msg)
}
