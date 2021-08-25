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

package static

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iTrellis/trellis/cmd"
	"github.com/iTrellis/trellis/internal/gin_middlewares"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/message"
)

func init() {
	cmd.DefaultCompManager.RegisterComponentFunc(
		&service.Service{Name: "trellis-server-static", Version: "v1"},
		NewStaticServer,
	)
}

type Handler struct {
	options component.Options

	ginMode string

	srv *http.Server
}

func NewStaticServer(opts ...component.Option) (component.Component, error) {
	h := &Handler{}
	for _, o := range opts {
		o(&h.options)
	}

	err := h.init()
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (p *Handler) init() error {
	p.ginMode = p.options.Config.GetString("gin_mode")

	gin.SetMode(p.ginMode)
	engine := gin.New()

	engine.Use(gin.Recovery())

	httpConf := p.options.Config.GetValuesConfig("http")

	ginHanlders := []gin.HandlerFunc{
		gin_middlewares.LoadCors(httpConf.GetValuesConfig("cors")),
	}

	for _, name := range gin_middlewares.IndexGinFuncs {
		ginHanlders = append(ginHanlders, gin_middlewares.UseFuncs[name])
	}
	engine.Use(ginHanlders...)

	pathConfig := httpConf.GetValuesConfig("path")

	for _, key := range pathConfig.GetKeys() {
		staticPath := pathConfig.GetString(key+".static", "/")
		staticRedirect := pathConfig.GetString(key + ".redirect")
		staticRoot := pathConfig.GetString(key+".root", "./root")

		if staticRedirect != "" {
			engine.GET(staticPath, func(c *gin.Context) {
				c.Redirect(http.StatusFound, staticRedirect)
			})
			engine.Static(staticRedirect, staticRoot)
		} else {
			engine.Static(staticPath, staticRoot)
		}
	}

	p.srv = &http.Server{
		Addr:    httpConf.GetString("address", ":8080"),
		Handler: engine,
	}
	return nil
}

func (*Handler) Route(message.Message) (interface{}, error) {
	return nil, nil
}

func (p *Handler) Start() error {

	ch := make(chan error)
	go func() {

		var err error

		sslConf := p.options.Config.GetValuesConfig("http.ssl")

		if sslConf != nil && sslConf.GetBoolean("enabled", false) {
			err = p.srv.ListenAndServeTLS(
				sslConf.GetString("cert-file"),
				sslConf.GetString("cert-key"),
			)
		} else {
			err = p.srv.ListenAndServe()
		}

		if err != http.ErrServerClosed {
			p.options.Logger.Error("http_server_closed", "err", err.Error())
		}

		ch <- err
	}()

	return <-ch
}

func (*Handler) Stop() error {
	return nil
}
