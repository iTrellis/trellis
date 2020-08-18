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
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/message/proto"
	"github.com/go-trellis/trellis/runner"
	"github.com/go-trellis/trellis/service"
	"github.com/go-trellis/trellis/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-trellis/common/errors"
	"github.com/go-trellis/config"
	"github.com/google/uuid"
)

func init() {
	service.RegistNewServiceFunc("trellis-trans-postapi", "v1", NewPostAPI)
}

// PostAPI api service
type PostAPI struct {
	mode string
	opts service.Options

	cfg config.Config

	forwardHeaders []string

	apis map[string]*api

	srv *http.Server
}

type api struct {
	proto.Service
	Topic string
}

// PostAPIResponse response
type PostAPIResponse struct {
	TraceID string      `json:"trace_id"`
	TraceIP string      `json:"trace_ip"`
	Code    uint64      `json:"code"`
	Msg     string      `json:"msg,omitempty"`
	Result  interface{} `json:"result"`
}

// NewPostAPI new api service
func NewPostAPI(opts ...service.OptionFunc) (service.Service, error) {

	s := &PostAPI{
		apis: make(map[string]*api),
	}

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

func (p *PostAPI) init() (err error) {

	p.mode = p.cfg.GetString("mode")

	gin.SetMode(p.mode)

	apisConf := p.cfg.GetValuesConfig("apis")

	typ := apisConf.GetString("type")
	switch typ {
	case "file":
		apis := apisConf.GetValuesConfig(typ)
		for _, apiKey := range apis.GetKeys() {
			apiConf := apisConf.GetMap("file." + apiKey)
			if apiConf == nil {
				return fmt.Errorf("init api failed: %s", apiKey)
			}

			api := &api{}

			api.Name = apiConf.Get("service_name")
			api.Version = apiConf.Get("service_version")
			api.Topic = apiConf.Get("topic")

			p.apis[apiConf.Get("api")] = api
		}
	case "mysql":
	default:
		return fmt.Errorf("unknown apis' config type")
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
func (p *PostAPI) Start() error {
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
func (p *PostAPI) Stop() error {

	dur := p.cfg.GetTimeDuration("http.shutdown-timeout", time.Second*30)

	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	if err := p.srv.Shutdown(ctx); err != nil {
		return errors.Newf("api shutdown failure, err: %s", err)
	}
	return nil
}

func (p *PostAPI) serve(ctx *gin.Context) {

	apiName := ctx.Request.Header.Get("X-API")

	r := &PostAPIResponse{
		TraceID: uuid.New().String(),
		TraceIP: utils.IPs()[0],
	}
	api, ok := p.apis[apiName]
	if !ok {
		r.Code = 1
		r.Msg = "api not found"
		ctx.JSON(http.StatusBadRequest, r)
		return
	}

	service := &proto.Service{Name: api.GetName(), Version: api.GetVersion()}
	rService, err := runner.GetService(api.GetName(), api.GetVersion())
	if err != nil {
		r.Code = 2
		r.Msg = err.Error()
		ctx.JSON(http.StatusInternalServerError, r)
		return
	}

	body, err := ctx.GetRawData()
	if err != nil {
		r.Code = 3
		r.Msg = err.Error()
		ctx.JSON(http.StatusBadRequest, r)
		return
	}

	msg := &message.Message{
		Payload: proto.Payload{
			TraceId: r.TraceID,
			Id:      uuid.New().String(),
			Service: service,
			ReqBody: body,
			Topic:   api.Topic,
			Header: map[string]string{
				"Content-Type": ctx.GetHeader("Content-Type"),
				"X-API":        ctx.GetHeader("X-API"),
				"Client-IP":    p.getClientIP(ctx),
				"Host": func() string {
					ip, err := internal.ExternalIP()
					if err != nil {
						return ""
					}
					return ip.String()
				}(),
			},
		},
	}

	fn := rService.Route(api.Topic)
	if fn == nil {
		r.Code = 4
		r.Msg = "topic not found"
		ctx.JSON(200, r)
		return
	}
	resp, err := fn(msg)
	if err != nil {
		r.Code = 5
		r.Msg = err.Error()
		ctx.JSON(200, r)
		return
	}

	r.Result = resp

	ctx.JSON(200, r)
}

// Route 路由
func (p *PostAPI) Route(string) service.HandlerFunc {
	// async中处理callback
	return nil
}

// getClientIP 获取客户端IP
func (*PostAPI) getClientIP(ctx *gin.Context) string {

	// Cdn-Src-Ip
	if ip := ctx.GetHeader("Cdn-Src-Ip"); ip != "" {
		return ip
	}

	// X-Forwarded-For
	if ips := ctx.GetHeader("X-Forwarded-For"); ips != "" {
		addr := strings.Split(ips, ",")
		if len(addr) > 0 && addr[0] != "" {
			rip, _, err := net.SplitHostPort(addr[0])
			if err != nil {
				rip = addr[0]
			}
			return rip
		}
	}

	// Client_Ip
	if ip := ctx.GetHeader("Client-Ip"); ip != "" {
		return ip
	}

	// RemoteAddr
	if ip, _, err := net.SplitHostPort(ctx.Request.RemoteAddr); err == nil {
		return ip
	}

	return ""
}
