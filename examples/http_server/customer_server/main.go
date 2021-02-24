package main

import (
	"github.com/iTrellis/trellis/cmd"
	"github.com/iTrellis/trellis/server/api"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/message"

	"github.com/gin-gonic/gin"
)

// custom handler
// curl -X 'POST' 'http://localhost:8080/ch'

// component handler
// curl -X 'POST' -H 'X-Api: custom.ping' 'http://localhost:8080/v1'

func init() {
	cmd.DefaultCompManager.RegisterComponentFunc(
		&service.Service{Domain: "custom", Name: "component_handler", Version: "v1"}, NewCompHandler)
	api.RegistCustomHandlers("custom_handler", "ch", "post", customHandler)
}

func main() {
	c := cmd.New()
	if err := c.Init(cmd.ConfigFile("config.yaml")); err != nil {
		panic(err)
	}

	if err := c.Start(); err != nil {
		panic(err)
	}

	c.BlockRun()
}

var defHandler *compHandler

func customHandler(c *gin.Context) {
	defHandler.options.Logger.Info("msg", "custom_handler")
	c.JSON(200, map[string]string{"message": defHandler.Response})
}

func NewCompHandler(opts ...component.Option) (component.Component, error) {
	defHandler = &compHandler{
		Response: "pong",
	}

	for _, o := range opts {
		o(&defHandler.options)
	}
	return defHandler, nil
}

type compHandler struct {
	Response string
	options  component.Options
}

func (p *compHandler) Route(topic string) component.Handler {
	switch topic {
	case "ping":
		return func(_ message.Message) (interface{}, error) {
			return p.Response, nil
		}
	}
	return nil
}

func (p *compHandler) Start() error {
	return nil
}

func (p *compHandler) Stop() error {
	return nil
}
