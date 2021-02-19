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
	"flag"
	"fmt"
	"time"

	"github.com/iTrellis/trellis/cmd"
	"github.com/iTrellis/trellis/examples/components"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/message"
)

var config string

func init() {
	flag.StringVar(&config, "config", "config.yaml", "config path")
}
func main() {

	c := cmd.New()
	if err := c.Init(cmd.ConfigFile(config)); err != nil {
		panic(err)
	}

	// Explicit to register component function
	cmd.DefaultCompManager.RegisterComponentFunc(&service.Service{Name: "component_ping", Version: "v1"},
		components.NewPing)

	if err := c.Start(); err != nil {
		panic(err)
	}

	defer c.Stop()

	time.Sleep(time.Second)

	cpt, err := cmd.DefaultCompManager.GetComponent(&service.Service{Name: "component_ping", Version: "v1"})
	if err != nil {
		panic(err)
	}

	hf := cpt.Route("etcd_ping")
	if hf == nil {
		panic("not found handler function")
	}
	for i := 0; i < 60; i++ {
		time.Sleep(time.Second)
		resp, err := hf(message.NewMessage())
		if err != nil {
			panic(err)
		}
		fmt.Println("get response:", resp)
	}
}
