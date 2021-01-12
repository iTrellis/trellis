package components

import (
	"context"

	"github.com/go-trellis/trellis/service"
	"github.com/go-trellis/trellis/service/component"
	"github.com/go-trellis/trellis/service/message"
)

type ping struct {
	alias string

	opts component.Options
}

func NewPing(alias string, opts ...component.Option) (component.Component, error) {
	c := &ping{alias: alias}
	for _, o := range opts {
		o(&c.opts)
	}
	return c, nil
}

func (p *ping) Alias() string {
	return p.alias
}

func (p *ping) Route(topic string) component.Handler {
	switch topic {
	case "ping":
		return func(_ context.Context, _ message.Message) (interface{}, error) {
			return p.opts.Router.Call(nil, message.NewOptionMessage(
				message.Options{
					Service: &service.Service{
						Name: "component_pong", Version: "v1",
					},
					Topic: "ping",
				},
			))
		}
	}
	return nil
}

func (p *ping) Start() error {
	println("component ping started")
	return nil
}

func (p *ping) Stop() error {
	println("component ping stopped")
	return nil
}
