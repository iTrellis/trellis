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
	"fmt"
	"strings"

	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/codec"
	"github.com/iTrellis/trellis/service/codec/json"
)

var (
	DefaultCodecs = map[string]codec.NewCodec{
		// "application/grpc":         grpc.NewCodec,
		// "application/grpc+json":    grpc.NewCodec,
		// "application/grpc+proto":   grpc.NewCodec,
		// "application/protobuf":     proto.NewCodec,
		// "application/json-rpc":     jsonrpc.NewCodec,
		// "application/proto-rpc":    protorpc.NewCodec,
		// "application/octet-stream": raw.NewCodec,
		service.MIMEApplicationJSON: json.NewCodec,
	}
)

type local struct {
	service *service.Service

	payload *Payload

	codec codec.Codec
}

func NewMessage(opts ...Option) Message {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	m := &local{
		service: options.Service,
		payload: options.Payload,
	}

	return m
}

func (p *local) Codec() codec.Codec {
	codec, _ := p.getCodec()
	return codec
}

func (p *local) contentType() string {
	header := p.payload.GetHeader()
	if header == nil {
		return service.MIMEApplicationJSON
	}
	v, ok := header["Content-Type"]
	if !ok {
		return service.MIMEApplicationJSON
	}
	cts := strings.Split(v, ";")
	switch len(cts) {
	case 0:
		return service.MIMEApplicationJSON
	default:
		return cts[0]
	}
}

func (p *local) getCodec() (codec.Codec, error) {
	if p.codec != nil {
		return p.codec, nil
	}

	ct := p.contentType()
	fn := DefaultCodecs[ct]
	if fn == nil {
		return nil, fmt.Errorf("unknown content-type: %s", ct)
	}

	p.codec = fn()
	return p.codec, nil
}

func (p *local) GetPayload() *Payload {
	return p.payload
}

func (p *local) Service() *service.Service {
	return p.service
}

func (p *local) Topic() string {
	return p.service.GetTopic()
}

func (p *local) SetBody(v interface{}) error {
	codec, err := p.getCodec()
	if err != nil {
		return err
	}
	if p.payload == nil {
		p.payload = &Payload{}
	}
	p.payload.Body, err = codec.Marshal(v)
	return err
}

func (p *local) ToObject(v interface{}) error {
	codec, err := p.getCodec()
	if err != nil {
		return err
	}
	return codec.Unmarshal(p.payload.GetBody(), v)
}

func (p *local) SetTopic(topic string) {
	if p.service == nil {
		p.service = &service.Service{}
	}
	p.service.Topic = topic
}

func (p *local) ToRemoteMessage() *RemoteMessage {
	return &RemoteMessage{
		Service: p.service,
		Payload: p.payload,
	}
}
