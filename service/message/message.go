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
	"context"

	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/codec"
)

// header keys
const (
	XAPIToken   = "x-api-token"
	XAPITraceID = "x-api-trace-id"
	XClientIP   = "x-client-ip"
)

// Message is the interface for publishing asynchronously
type Message interface {
	Service() *service.Service
	Codec() codec.Codec
	Topic() string
	SetTopic(string)
	SetBody(v interface{}) error
	GetPayload() *Payload
	ToObject(v interface{}) error
}

// Caller caller for calling component or server
type Caller interface {
	CallComponent(context.Context, Message) (interface{}, error)
	CallServer(ctx context.Context, msg Message) (interface{}, error)
}

func (p *Payload) Set(key, value string) {
	header := p.GetHeader()
	if header == nil {
		header = make(map[string]string)
	}

	header[key] = value

	p.Header = header
}

func (p *Payload) Get(key string) string {
	header := p.GetHeader()
	if header == nil {
		return ""
	}

	return header[key]
}
