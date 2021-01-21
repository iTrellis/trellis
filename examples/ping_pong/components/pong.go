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

	"github.com/iTrellis/trellis/route"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/message"
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
