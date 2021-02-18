package main

import (
	"flag"

	"github.com/iTrellis/trellis/cmd"

	_ "github.com/iTrellis/trellis/examples/components"
	_ "github.com/iTrellis/trellis/server/grpc"
)

var config string

func init() {
	flag.StringVar(&config, "config", "config.yaml", "config path")
}

func main() {

	flag.Parse()

	c := cmd.New()
	if err := c.Init(cmd.ConfigFile(config)); err != nil {
		panic(err)
	}

	if err := c.Start(); err != nil {
		panic(err)
	}

	defer c.Stop()

	c.Run()
}
