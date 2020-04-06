package message

import (
	"github.com/go-trellis/trellis/codec"
	"github.com/go-trellis/trellis/message/proto"
)

type Message struct {
	*proto.Payload

	codec codec.Codec
}

func (p *Message) SetBody(obj interface{}) (err error) {

	p.Payload.Body, err = p.codec.Marshal(obj)
	if err != nil {
		return err
	}

	return nil
}

func (p *Message) ToObject(obj interface{}) error {
	codec, err := codec.GetCodec(p.Payload.GetContentType())
	if err != nil {
		return err
	}

	return codec.UnmarshalObject(p.Payload.GetBody(), obj)
}

// // Request is a synchronous request interface
// type Request interface {
// 	ID() string
// 	// Service name requested
// 	Service() string
// 	// Service version requested
// 	Version() string
// 	// The action requested
// 	Method() string
// 	// The handler's name in service
// 	Topic() string
// 	// Endpoint name requested
// 	Endpoint() string
// 	// The encoded message stream
// 	Codec() codec.Codec
// 	// Read the undecoded request body
// 	Read() ([]byte, error)
// 	// Set Request header
// 	SetHeader(string, string)
// 	Message() *Message
// }

// // Response 服务端返回的内容
// type Response interface {
// 	Code() uint64
// 	Codec() codec.Codec
// 	// Get the header value by key
// 	GetHeader(string) string
// 	GetHeaders() map[string]string

// 	ToObject(interface{}) error
// }
