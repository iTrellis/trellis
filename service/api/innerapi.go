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
	"fmt"
	"net/http"
	"time"

	"github.com/go-trellis/trellis/codec"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/service"

	"github.com/gin-gonic/gin"
	"github.com/go-trellis/common/errors"
	"github.com/go-trellis/config"
)

func init() {
	service.RegistNewServiceFunc("trellis-innerapi", "v1", NewService)
}

var useFuncs = make(map[string]gin.HandlerFunc)

// RegistUseFuncs 注册
func RegistUseFuncs(name string, fn gin.HandlerFunc) error {
	_, ok := useFuncs[name]
	if ok {
		return fmt.Errorf("use funcs (%s) is already exist", name)
	}
	useFuncs[name] = fn
	return nil
}

// Service api service
type Service struct {
	mode string
	opts service.Options

	cfg config.Config

	forwardHeaders []string

	srv *http.Server
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

	p.mode = p.cfg.GetString("mode")

	gin.SetMode(p.mode)

	httpConf := p.cfg.GetValuesConfig("http")

	urlPath := httpConf.GetString("path", "/")

	staticPaths := httpConf.GetMap("static_path")

	engine := gin.New()

	engine.Use(gin.Recovery())

	for _, fn := range useFuncs {
		engine.Use(fn)
	}

	for path, static := range staticPaths {
		s, ok := static.(string)
		if !ok {
			return fmt.Errorf("static path is invalid: %s", path)
		}
		engine.Static(path, s)
	}

	p.forwardHeaders = httpConf.GetStringList("forward.headers")

	internal.LoadCors(engine, httpConf.GetValuesConfig("cors"))
	internal.LoadPprof(engine, httpConf.GetValuesConfig("pprof"))

	// router.ServeHTTP()
	engine.POST(urlPath, p.serve)

	p.srv = &http.Server{
		Addr:    httpConf.GetString("address", ":8080"),
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
			p.opts.Logger.Error("http_server_closed", err)
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

	body := &bytes.Buffer{}
	_, err := body.ReadFrom(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, msg)
		return
	}

	msgcodeC, err := codec.GetCodec(ctx.Request.Header.Get("content-type"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, msg)
		return
	}

	err = msgcodeC.Unmarshal(body.Bytes(), msg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, msg)
		return
	}

	for _, h := range p.forwardHeaders {
		msg.SetHeader(h, ctx.GetHeader(h))
	}

	p.opts.Logger.Info("request", "message", msg)

	resp, err := service.CallServer(msg,
		fmt.Sprintf("%s-%s", msg.GetService().String(), msg.GetHeader("Client-IP")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, msg)
		return
	}

	bs, err := msgcodeC.Marshal(resp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, msg)
		return
	}

	msg.SetBody(bs)

	ctx.JSON(200, msg)
}

// Route 路由
func (p *Service) Route(string) service.HandlerFunc {
	// async中处理callback
	return nil
}
