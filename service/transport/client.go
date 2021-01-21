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

package transport

import (
	"context"

	"github.com/iTrellis/node"
	"github.com/iTrellis/trellis/service/message"
)

func Call(nd *node.Node, ctx context.Context, msg message.Message) (interface{}, error) {
	return nil, nil
}

// import (
// 	"context"

// 	"github.com/iTrellis/trellis/service/codec"
// 	"github.com/iTrellis/trellis/service/message"
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
