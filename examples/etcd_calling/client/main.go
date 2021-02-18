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
	c.GetRoutesManager().CompManager().RegisterComponentFunc(
		&service.Service{Name: "component_ping", Version: "v1"},
		components.NewPing)

	if err := c.Start(); err != nil {
		panic(err)
	}

	defer c.Stop()

	time.Sleep(time.Second)

	cpt, err := c.GetRoutesManager().CompManager().GetComponent(&service.Service{Name: "component_ping", Version: "v1"})
	if err != nil {
		panic(err)
	}

	hf := cpt.Route("etcd_ping")
	if hf == nil {
		panic("not found handler function")
	}
	for i := 0; i < 10000; i++ {
		time.Sleep(time.Second)
		resp, err := hf(message.NewMessage())
		if err != nil {
			panic(err)
		}
		fmt.Println("get response:", resp)
	}
	c.Run()
}
