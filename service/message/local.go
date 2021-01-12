package message

import (
	"github.com/go-trellis/trellis/service"
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
