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
		return nil, fmt.Errorf("not found service")
	}
	node, ok := nm.NodeFor(keys...)
	if !ok {
		return nil, fmt.Errorf("not found service")
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
