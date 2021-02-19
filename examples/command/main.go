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

package main

import (
	"fmt"
	"os"

	"github.com/iTrellis/trellis/cmd"
	"github.com/iTrellis/trellis/configure"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
)

var s = service.Service{Name: "command_example", Version: "v1"}

func init() {
	cmd.DefaultCompManager.RegisterComponentFunc(&s, newSimpleComp)
}

type command struct{}

func newSimpleComp(...component.Option) (component.Component, error) {
	return &command{}, nil
}

func (p *command) Start() error {
	fmt.Println("do something")
	return nil
}

func (p *command) Stop() error {
	fmt.Println("stop something")
	return nil
}

func (p *command) Route(topic string) component.Handler {
	return nil
}

func main() {

	cs := &configure.Service{Service: s}
	c := cmd.New(cmd.WithConfig(&configure.Configure{Project: configure.Project{
		Services: []*configure.Service{cs},
	}}))

	// if err := c.Start(); err != nil {
	// 	panic(err)
	// }

	app := c.App()

	app.Run(os.Args)
}
