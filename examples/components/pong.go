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
	"fmt"

	"github.com/iTrellis/trellis/cmd"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/message"
)

func init() {
	cmd.RegisterComponentFunc(&service.Service{Name: "component_pong", Version: "v1"}, NewPong)
}

type pong struct{}

// NewPong pong constrator
func NewPong(opts ...component.Option) (component.Component, error) {
	return &pong{}, nil
}

func (p *pong) Route(msg message.Message) (interface{}, error) {
	switch msg.Topic() {
	case "ping":
		return "pong", nil
	}
	return nil, fmt.Errorf("unknown topic")
}

func (p *pong) Start() error {
	println("component pong started")
	return nil
}

func (p *pong) Stop() error {
	println("component pong stopped")
	return nil
}
