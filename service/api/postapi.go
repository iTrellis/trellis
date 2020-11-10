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
	"net/http"
	"sync"
	"time"

	"github.com/go-trellis/trellis/clients"
	"github.com/go-trellis/trellis/errcode"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/message/proto"
	"github.com/go-trellis/trellis/service"

	"github.com/gin-gonic/gin"
	"github.com/go-trellis/common/errors"
	"github.com/go-trellis/common/formats"
	"github.com/go-trellis/txorm"
	"xorm.io/xorm"
)

func init() {
	service.RegistNewServiceFunc("trellis-postapi", "v1", NewPostAPI)
}

// PostAPI api service
type PostAPI struct {
	mode string
	opts service.Options

	forwardHeaders []string

	apis map[string]*API

	srv *http.Server

	ticker    *time.Ticker
	syncer    sync.RWMutex
	apiEngine *xorm.Engine
}

// Response response
type Response struct {
	TraceID   string      `json:"trace_id"`
	TraceIP   string      `json:"trace_ip"`
	Host      string      `json:"host"`
	Code      uint64      `json:"code"`
	Namespace string      `json:"namespace,omitempty"`
	Msg       string      `json:"msg,omitempty"`
	Result    interface{} `json:"result"`
}

// NewPostAPI new api service
func NewPostAPI(opts ...service.OptionFunc) (service.Service, error) {

	s := &PostAPI{
		apis: make(map[string]*API),
	}

	for _, o := range opts {
		o(&s.opts)
	}

	err := s.init()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (p *PostAPI) init() (err error) {

	p.mode = p.opts.Config.GetString("mode")

	gin.SetMode(p.mode)

	apisConf := p.opts.Config.GetValuesConfig("apis")

	typ := apisConf.GetString("type")
	switch typ {
	case "file":
		apis := apisConf.GetValuesConfig(typ)
		for _, apiKey := range apis.GetKeys() {
			apiConf := apisConf.GetValuesConfig("file." + apiKey)
			if apiConf == nil {
				return fmt.Errorf("init api failed: %s", apiKey)
			}

			api := &API{
				Name:           apiConf.GetString("api"),
				Topic:          apiConf.GetString("topic"),
				ServiceName:    apiConf.GetString("service_name"),
				ServiceVersion: apiConf.GetString("service_version"),
			}

			p.apis[api.Name] = api
		}
	case "mysql":

		databaseConf := apisConf.GetValuesConfig(typ)

		engines, err := txorm.NewEnginesFromConfig(databaseConf)
		if err != nil {
			return err
		}
		p.apiEngine = engines[txorm.DefaultDatabase]

		ticker := formats.ParseStringTime(apisConf.GetString("ticker", "30s"))

		p.ticker = time.NewTicker(ticker)

		go p.syncAPIs()

	default:
		return fmt.Errorf("unknown apis' config type")
	}

	httpConf := p.opts.Config.GetValuesConfig("http")

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

	engine.POST(urlPath, p.serve)

	p.srv = &http.Server{
		Addr:    httpConf.GetString("address", ":8080"),
		Handler: engine,
	}

	return
}

// Start start service
func (p *PostAPI) Start() error {

	ch := make(chan error)
	go func() {

		var err error

		sslConf := p.opts.Config.GetValuesConfig("http.ssl")

		if sslConf != nil && sslConf.GetBoolean("enabled", false) {
			err = p.srv.ListenAndServeTLS(
				sslConf.GetString("cert-file"),
				sslConf.GetString("cert-key"),
			)
		} else {
			err = p.srv.ListenAndServe()
		}

		if err != http.ErrServerClosed {
			p.opts.Logger.Error("http_server_closed", err.Error())
		}

		ch <- err
	}()

	return <-ch
}

// Stop stop service
func (p *PostAPI) Stop() error {

	dur := p.opts.Config.GetTimeDuration("http.shutdown-timeout", time.Second*30)

	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	if err := p.srv.Shutdown(ctx); err != nil {
		return errors.Newf("api shutdown failure, err: %s", err)
	}
	return nil
}

func (p *PostAPI) serve(gCtx *gin.Context) {

	apiName := gCtx.Request.Header.Get("X-API")
	clientIP := internal.GetClientIP(gCtx)

	msg := message.NewMessage()

	r := &Response{
		TraceID: msg.GetTraceId(),
		TraceIP: msg.GetTraceIp(),
	}
	p.opts.Logger.Info("request", "trace_id", r.TraceID, "api_name", apiName, "client_ip", clientIP)
	api, ok := p.getAPI(apiName)
	if !ok {
		apiErr := errcode.ErrAPINotFound.New()
		r.Code = apiErr.Code()
		r.Msg = apiErr.Error()
		r.Namespace = apiErr.Namespace()
		gCtx.JSON(http.StatusBadRequest, r)
		p.opts.Logger.Error("api_not_found", "trace_id", r.TraceID, "api_name", apiName, "client_ip", clientIP)
		return
	}

	body, err := gCtx.GetRawData()
	if err != nil {
		getErr := errcode.ErrBadRequest.New(errors.Params{"err": err.Error()})
		r.Code = getErr.Code()
		r.Msg = getErr.Error()
		r.Namespace = getErr.Namespace()
		gCtx.JSON(http.StatusBadRequest, r)
		p.opts.Logger.Error("get_raw_data", "trace_id", r.TraceID,
			"api_name", apiName, "client_ip", clientIP, "err", err)
		return
	}
	msg.SetBody(body)

	msg.Service = &proto.Service{Name: api.ServiceName, Version: api.ServiceVersion}
	msg.Topic = api.Topic

	msg.SetHeader("Client-IP", clientIP)
	for _, h := range p.forwardHeaders {
		msg.SetHeader(h, gCtx.GetHeader(h))
	}

	resp, err := clients.CallService(msg, fmt.Sprintf("%s-%s", msg.GetService().String(), clientIP))
	if err == nil {
		r.Result = resp
		gCtx.JSON(200, r)
		return
	}

	// errors
	switch et := err.(type) {
	case errors.ErrorCode:
		r.Code = et.Code()
		r.Msg = et.Error()
		r.Namespace = et.Namespace()
	case errors.SimpleError:
		cErr := errcode.ErrCallService.New(errors.Params{"err": et.Error()})
		r.Code = cErr.Code()
		r.Msg = et.Error()
		r.Namespace = et.Namespace()
	default:
		cErr := errcode.ErrCallService.New(errors.Params{"err": et.Error()})
		r.Code = cErr.Code()
		r.Msg = cErr.Error()
		r.Namespace = cErr.Namespace()
	}

	p.opts.Logger.Error("call_server_failed", "trace_id", r.TraceID,
		"api_name", apiName, "client_ip", clientIP, "err", r)
	gCtx.JSON(200, r)
}

func (p *PostAPI) getAPI(name string) (*API, bool) {
	p.syncer.RLock()
	api, ok := p.apis[name]
	p.syncer.RUnlock()
	return api, ok
}

// Route 路由
func (p *PostAPI) Route(string) service.HandlerFunc {
	return nil
}
