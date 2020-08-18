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

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-trellis/config"

	"github.com/go-trellis/trellis/codec"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/runner"
	"github.com/go-trellis/trellis/service"
	"github.com/go-trellis/trellis/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-trellis/common/errors"
)

func init() {
	service.RegistNewServiceFunc("trellis-trans-api", "v1", NewService)
}

// Service api service
type Service struct {
	debug bool
	opts  service.Options

	cfg config.Config

	forwardHeaders []string

	srv *http.Server
}

type Response struct {
	Code   uint64      `json:"code"`
	Msg    string      `json:"msg,omitempty"`
	Result interface{} `json:"result"`
}

// NewService new api service
func NewService(opts ...service.OptionFunc) (service.Service, error) {

	s := &Service{}

	for _, o := range opts {
		o(&s.opts)
	}

	s.cfg = config.DefaultGetter.GenMapConfig(s.opts.Config)

	err := s.init()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (p *Service) init() (err error) {

	p.debug = p.cfg.GetBoolean("debug")
	if !p.debug {
		gin.SetMode("release")
	}

	httpConf := p.cfg.GetConfig("http")

	urlPath := httpConf.GetString("path", "/")

	engine := gin.New()

	engine.Use(gin.Recovery())

	address := httpConf.GetString("address", ":8080")

	forwardHeaders := httpConf.GetStringList("forward.headers")

	p.forwardHeaders = forwardHeaders

	utils.LoadCors(engine, httpConf.GetConfig("cors"))
	utils.LoadPprof(engine, httpConf.GetConfig("pprof"))

	// router.ServeHTTP()
	engine.POST(urlPath, p.serve)

	p.srv = &http.Server{
		Addr:    address,
		Handler: engine,
	}

	return
}

// Start start service
func (p *Service) Start() error {
	go func() {

		var err error

		sslConf := p.cfg.GetConfig("http.ssl")

		if sslConf != nil && sslConf.GetBoolean("enabled", false) {
			err = p.srv.ListenAndServeTLS(
				sslConf.GetString("cert-file"),
				sslConf.GetString("cert-key"),
			)
		} else {
			err = p.srv.ListenAndServe()
		}

		if err != http.ErrServerClosed {
			// print log
		}
	}()
	return nil
}

// Stop stop service
func (p *Service) Stop() error {

	dur := p.cfg.GetTimeDuration("http.shutdown-timeout", time.Second*30)

	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	if err := p.srv.Shutdown(ctx); err != nil {
		return errors.Newf("api shutdown failure, err: %s", err)
	}
	return nil
}

func (p *Service) serve(ctx *gin.Context) {
	msg := &message.Message{}
	// msg := &proto.Payload{}
	r := &Response{}
	msgcodeC, err := codec.GetCodec(ctx.Request.Header.Get("content-type"))
	if err != nil {
		r.Msg = err.Error()
		ctx.JSON(http.StatusBadRequest, r)
		return
	}
	switch msgcodeC.String() {
	case codec.JSON:
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(ctx.Request.Body)
		if err != nil {
			r.Msg = err.Error()
			ctx.JSON(http.StatusBadRequest, r)
			return
		}
		err = json.Unmarshal(body.Bytes(), msg)
		if err != nil {
			r.Msg = err.Error()
			ctx.JSON(http.StatusBadRequest, r)
			return
		}
	default:
		r.Msg = fmt.Sprintf("unsupported codec, %s", msgcodeC.String())
		ctx.JSON(http.StatusBadRequest, r)
		return
	}

	rService, err := runner.GetService(msg.GetService().GetName(), msg.GetService().GetVersion())
	if err != nil {
		r.Msg = err.Error()
		ctx.JSON(http.StatusInternalServerError, r)
		return
	}

	hf := rService.Route(msg.GetTopic())
	if hf == nil {
		r.Msg = "topic not found"
		ctx.JSON(200, r)
		return
	}
	resp, err := hf(msg)
	if err != nil {
		r.Msg = err.Error()
		ctx.JSON(200, r)
		return
	}

	r.Result = resp

	ctx.JSON(200, r)
}

// Route 路由
func (p *Service) Route(string) service.HandlerFunc {
	// async中处理callback
	return nil
}
