package components

import (
	"context"

	"github.com/go-trellis/trellis/route"
	"github.com/go-trellis/trellis/service"
	"github.com/go-trellis/trellis/service/component"
	"github.com/go-trellis/trellis/service/message"
)

func init() {
	route.RegisterComponentFunc(&service.Service{Name: "component_pong", Version: "v1"}, NewPong)
}

type pong struct {
	alias string
}

func NewPong(alias string, opts ...component.Option) (component.Component, error) {
	return &pong{alias: alias}, nil
}

func (p *pong) Alias() string {
	return p.alias
}

func (p *pong) Route(topic string) component.Handler {
	switch topic {
	case "ping":
		return func(_ context.Context, _ message.Message) (interface{}, error) {
			return "pong", nil
		}
	}
	return nil
}

func (p *pong) Start() error {
	println("component pong started")
	return nil
}

func (p *pong) Stop() error {
	println("component pong stopped")
	return nil
}
