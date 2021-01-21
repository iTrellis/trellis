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
	"github.com/iTrellis/trellis/service"
)

type message struct {
	opts Options
}

func NewMessage(fs ...Option) Message {
	m := &message{}
	for _, o := range fs {
		o(&m.opts)
	}

	return m
}

func NewOptionMessage(opts Options) Message {
	return &message{opts: opts}
}

func (p *message) ContentType() string {
	return p.opts.ContentType
}

func (p *message) Payload() *BasePayload {
	return p.opts.Payload
}

func (p *message) Service() *service.Service {
	return p.opts.Service
}

func (p *message) Topic() string {
	return p.opts.Topic
}
