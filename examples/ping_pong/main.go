package main

import (
	"fmt"

	"github.com/go-trellis/trellis/examples/ping_pong/components"

	"github.com/go-trellis/trellis/cmd"
	"github.com/go-trellis/trellis/service"
)

func main() {
	c := cmd.New()
	if err := c.Init(cmd.ConfigFile("config.yaml")); err != nil {
		panic(err)
	}

	// Explicit to register component function
	c.RegisterComponentFunc(&service.Service{Name: "component_ping", Version: "v1"}, components.NewPing)

	// // implicit in pong.go
	// c.RegisterComponentFunc(&service.Service{Name: "component_pong", Version: "v1"}, components.NewPong)

	if err := c.Start(); err != nil {
		panic(err)
	}

	defer c.Stop()

	cpt, err := c.GetComponent(&service.Service{Name: "component_ping", Version: "v1"})
	if err != nil {
		panic(err)
	}

	hf := cpt.Route("ping")
	resp, err := hf(nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("get response:", resp)
}
