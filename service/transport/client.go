package transport

import (
	"context"

	"github.com/go-trellis/node"
	"github.com/go-trellis/trellis/service/message"
)

func Call(nd *node.Node, ctx context.Context, msg message.Message) (interface{}, error) {
	return nil, nil
}

// import (
// 	"context"

// 	"github.com/go-trellis/trellis/service/codec"
// 	"github.com/go-trellis/trellis/service/message"
// )

// type Client interface {
// 	NewMessage(topic string, msg interface{}, opts ...MessageOption) message.Message
// 	NewRequest(service, endpoint string, req interface{}, reqOpts ...RequestOption) Request
// 	Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error

// 	// Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error)
// 	// Publish(ctx context.Context, msg Message, opts ...PublishOption) error
// 	String() string
// }

// type CallOption func(*CallOptions)

// type CallOptions struct {
// 	Keys []string
// }

// // Request is the interface for a synchronous request used by Call or Stream
// type Request interface {
// 	// The service to call
// 	Service() string
// 	// The action to take
// 	Method() string
// 	// The endpoint to invoke
// 	Endpoint() string
// 	// The content type
// 	ContentType() string
// 	// The unencoded request body
// 	Body() interface{}
// 	// Write to the encoded request writer. This is nil before a call is made
// 	Codec() codec.Writer
// 	// indicates whether the request will be a streaming one rather than unary
// 	Stream() bool
// }

// // RequestOption used by NewRequest
// type RequestOption func(*RequestOptions)

// type RequestOptions struct {
// 	ContentType string
// 	Stream      bool

// 	// Other options for implementations of the interface
// 	// can be stored in a context
// 	Context context.Context
// }

// // Response is the response received from a service
// type Response interface {
// 	// Read the response
// 	Codec() codec.Reader
// 	// read the header
// 	Header() map[string]string
// 	// Read the undecoded response
// 	Read() ([]byte, error)
// }
