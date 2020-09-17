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

package message

import (
	"github.com/go-trellis/trellis/codec"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message/proto"

	"github.com/google/uuid"
)

// Message Message
type Message struct {
	*proto.Payload `json:",inline" yaml:",inline"`

	codecer codec.Codec
}

// SetBody set request body
func (p *Message) SetBody(body interface{}) error {

	switch bs := body.(type) {
	case []byte:
		p.Payload.ReqBody = bs
		return nil
	case string:
		p.Payload.ReqBody = []byte(bs)
		return nil
	}
	err := p.getCodecer()
	if err != nil {
		return err
	}
	p.Payload.ReqBody, err = p.codecer.Marshal(body)
	return err
}

// ToObject codec with request body to object
func (p *Message) ToObject(obj interface{}) error {
	if err := p.getCodecer(); err != nil {
		return err
	}
	return p.codecer.Unmarshal(p.Payload.GetReqBody(), obj)
}

func (p *Message) getCodecer() error {
	if p.codecer != nil {
		return nil
	}

	header := p.GetHeader("Content-Type")
	if header == "" {
		// default json
		header = codec.JSON
	}
	c, err := codec.GetCodec(header)
	if err != nil {
		return err
	}
	p.codecer = c

	return nil
}

// NewMessage new message
func NewMessage() *Message {
	return &Message{
		Payload: &proto.Payload{
			TraceId: uuid.New().String(),
			TraceIp: func() string {
				ip, err := internal.ExternalIP()
				if err != nil {
					return ""
				}
				return ip.String()
			}(),
			Id:     uuid.New().String(),
			Header: make(map[string]string),
		},
	}
}

// Copy copy message by base message
func (p *Message) Copy() *Message {
	if p == nil {
		return nil
	}
	return &Message{
		Payload: &proto.Payload{
			TraceId: p.Payload.TraceId,
			TraceIp: p.Payload.TraceIp,
			Id:      uuid.New().String(),
			Header:  p.Payload.Header,
		},
	}
}

// GetHeader get header value with key
func (p *Message) GetHeader(key string) string {
	return p.Payload.Header[key]
}

// SetHeader set header value with key
func (p *Message) SetHeader(key, value string) {
	p.Payload.Header[key] = value
}
