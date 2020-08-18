package client

// import (
// 	"fmt"

// 	"github.com/go-trellis/trellis/message"
// )

// // Type 客户端类型
// type Type string

// // 各种客户端类型
// const (
// 	HTTP  Type = "http"
// 	INNER Type = "inner"
// 	RPC   Type = "rpc"
// )

// // type Options struct {
// // 	cfg config.Config
// // }

// // Client is the interface used to make requests to services.
// // It supports Request/Response via Transport and Publishing via the Broker.
// // It also supports bidirectional streaming of requests.
// type Client interface {
// 	Call(req *message.Request) (*message.Response, error)
// 	// Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error
// 	String() string
// }

// // Call 根据类型发送请求
// func Call(ct Type, req *message.Request) (*message.Response, error) {
// 	switch ct {
// 	case HTTP:
// 		return NewHTTPClient().Call(req)
// 	case RPC:
// 		return nil, nil
// 	default:
// 		return nil, fmt.Errorf("unkown supperted client type: %s", ct)
// 	}
// }
