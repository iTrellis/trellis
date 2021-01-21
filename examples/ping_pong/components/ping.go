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

package components

import (
	"context"

	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/message"
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
