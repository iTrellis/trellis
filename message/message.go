package message

import (
	"fmt"

	"github.com/go-trellis/trellis/codec"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message/proto"
	"github.com/google/uuid"
)

// Message Message
type Message struct {
	proto.Payload `json:",inline" yaml:",inline"`

	codecer codec.Codec
}

// SetBody set request body
func (p *Message) SetBody(body []byte) error {
	p.ReqBody = body
	return nil
}

// ToObject codec with request body to object
func (p *Message) ToObject(obj interface{}) error {
	if err := p.getCodecer(); err != nil {
		return err
	}
	return p.codecer.Unmarshal(p.GetReqBody(), obj)
}

func (p *Message) getCodecer() error {
	if p.codecer != nil {
		return nil
	}

	header := p.GetHeader()
	if header == nil {
		return fmt.Errorf("header is nil")
	}
	c, err := codec.GetCodec(header["Content-Type"])
	if err != nil {
		return err
	}
	p.codecer = c

	return nil
}

func NewMessage() *Message {
	return &Message{
		Payload: proto.Payload{
			TraceId: uuid.New().String(),
			TraceIp: func() string {
				ip, err := internal.ExternalIP()
				if err != nil {
					return ""
				}
				return ip.String()
			}(),
			Id: uuid.New().String(),
		},
	}
}

func (p *Message) Copy() *Message {
	if p == nil {
		return nil
	}
	return &Message{
		Payload: proto.Payload{
			TraceId: p.TraceId,
			TraceIp: p.TraceIp,
			Id:      uuid.New().String(),
		},
	}
}
