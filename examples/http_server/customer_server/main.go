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
	"log"

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
	c, err := cmd.New()
	if err != nil {
		log.Fatalln(err)
	}

	if err := c.Init(cmd.ConfigFile("config.yaml")); err != nil {
		log.Fatalln(err)
	}

	if err := c.Start(); err != nil {
		log.Fatalln(err)
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

func (p *compHandler) Route(msg message.Message) (interface{}, error) {
	switch msg.Topic() {
	case "ping":
		return p.Response, nil
	}
	return nil, nil
}

func (p *compHandler) Start() error {
	return nil
}

func (p *compHandler) Stop() error {
	return nil
}
