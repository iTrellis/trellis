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
	"strings"
	"sync"
	"time"

	"github.com/iTrellis/trellis/cmd"
	"github.com/iTrellis/trellis/internal/addr"
	"github.com/iTrellis/trellis/server"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/message"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iTrellis/common/errors"
	"github.com/iTrellis/common/formats"
	"github.com/iTrellis/xorm_ext"
	"xorm.io/xorm"
)

func init() {
	cmd.DefaultCompManager.RegisterComponentFunc(
		&service.Service{Name: "trellis-postapi", Version: "v1"},
		NewHTTPServer,
	)
}

var handlers = make(map[string]*server.Handler)

func RegistCustomHandlers(name, path, method string, fn gin.HandlerFunc) {

	if fn == nil {
		panic("handler function should not be nil")
	}
	_, ok := handlers[name]
	if ok {
		panic(fmt.Errorf("handler isalread exists"))
	}

	handlers[name] = &server.Handler{Name: name, URLPath: path, Method: strings.ToUpper(method), Func: fn}
}

type httpServer struct {
	ginMode string

	mode string // LOCAL, REMOTE

	forwardHeaders []string

	apis map[string]*API

	options component.Options

	srv *http.Server

	ticker    *time.Ticker
	syncer    sync.RWMutex
	apiEngine *xorm.Engine
}

// Response response
type Response struct {
	TraceID   string      `json:"trace_id"`
	TraceIP   string      `json:"trace_ip"`
	Code      uint64      `json:"code"`
	Namespace string      `json:"namespace,omitempty"`
	Msg       string      `json:"msg,omitempty"`
	Result    interface{} `json:"result"`
}

// NewHTTPServer new api service
func NewHTTPServer(opts ...component.Option) (component.Component, error) {

	s := &httpServer{
		apis: make(map[string]*API),
	}

	for _, o := range opts {
		o(&s.options)
	}

	err := s.init()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (p *httpServer) init() error {

	p.mode = p.options.Config.GetString("mode")
	p.ginMode = p.options.Config.GetString("gin_mode")

	gin.SetMode(p.ginMode)

	apisConf := p.options.Config.GetValuesConfig("apis")

	typ := apisConf.GetString("type", "file")
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
				ServiceDomain:  apiConf.GetString("service_domain"),
				ServiceName:    apiConf.GetString("service_name"),
				ServiceVersion: apiConf.GetString("service_version"),
			}

			p.apis[api.Name] = api
		}
	case "mysql":

		databaseConf := apisConf.GetValuesConfig(typ)

		engines, err := xorm_ext.NewEnginesFromConfig(databaseConf)
		if err != nil {
			return err
		}
		p.apiEngine = engines[xorm_ext.DefaultDatabase]

		ticker := formats.ParseStringTime(apisConf.GetString("ticker", "30s"))

		p.ticker = time.NewTicker(ticker)

		go p.syncAPIs()
	default:
		return fmt.Errorf("unknown apis' config type")
	}

	httpConf := p.options.Config.GetValuesConfig("http")

	engine := gin.New()

	engine.Use(gin.Recovery())

	loadCors(engine, httpConf.GetValuesConfig("cors"))
	loadPprof(engine, httpConf.GetValuesConfig("pprof"))

	for _, fn := range useFuncs {
		engine.Use(fn)
	}

	for _, v := range handlers {
		// p.options.Logger.Info("start_costomer_handler", v.Name, v.URLPath, v.Method)
		engine.Handle(v.Method, v.URLPath, v.Func)
	}

	urlPath := httpConf.GetString("postapi")
	if len(urlPath) != 0 {
		engine.POST(urlPath, p.serve)
	}

	p.forwardHeaders = httpConf.GetStringList("forward.headers")

	p.srv = &http.Server{
		Addr:    httpConf.GetString("address", ":8080"),
		Handler: engine,
	}

	return nil
}

func (p *httpServer) Route(topic string) component.Handler {
	return nil
}

func (p *httpServer) Start() error {

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
			// p.options.Logger.Error("http_server_closed", err.Error())
		}

		ch <- err
	}()

	return <-ch
}

func (p *httpServer) Stop() error {

	dur := p.options.Config.GetTimeDuration("http.shutdown-timeout", time.Second*30)

	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	if err := p.srv.Shutdown(ctx); err != nil {
		return errors.Newf("api shutdown failure, err: %s", err)
	}
	return nil
}

func (p *httpServer) serve(gCtx *gin.Context) {

	apiName := gCtx.Request.Header.Get("X-API")
	clientIP := addr.GetClientIP(gCtx.Request)

	traceID := gCtx.GetHeader(message.XAPITraceID)
	if traceID == "" {
		traceID = uuid.NewString()
	}

	r := &Response{
		TraceID: traceID,
		TraceIP: addr.ExternalIPs()[0],
	}

	// p.options.Logger.Info("request", "trace_id", traceID, "api_name", apiName, "client_ip", clientIP)
	api, ok := p.getAPI(apiName)
	if !ok {
		r.Code = 11
		r.Msg = "api not found"
		r.Namespace = "trellis"
		gCtx.JSON(http.StatusBadRequest, r)
		// p.options.Logger.Error("api_not_found", "trace_id", traceID, "api_name", apiName, "client_ip", clientIP)
		return
	}

	body, err := gCtx.GetRawData()
	if err != nil {
		r.Code = 10
		r.Msg = fmt.Sprintf("bad request: %s", err.Error())
		r.Namespace = "trellis"
		gCtx.JSON(http.StatusBadRequest, r)
		// p.options.Logger.Error("get_raw_data", "trace_id", r.TraceID, "api_name", apiName, "client_ip", clientIP, "err", err)
		return
	}

	payload := &message.Payload{
		Header: make(map[string]string),
		Body:   body,
	}

	payload.Set(message.XClientIP, clientIP)
	payload.Set(message.XAPITraceID, traceID)
	for _, h := range p.forwardHeaders {
		payload.Set(h, gCtx.GetHeader(h))
	}

	msg := message.NewMessage(message.Service(
		service.Service{
			Domain:  api.ServiceDomain,
			Name:    api.ServiceName,
			Version: api.ServiceVersion,
			Topic:   api.Topic}),
		message.MessagePayload(payload),
	)

	var resp interface{}
	switch p.mode {
	case "local", "":
		// resp, err = p.options.Manager.Call(msg,
		// 	component.Timeout(30*time.Second),
		// 	component.Keys(msg.Service().FullPath(), clientIP))

		resp, err = p.options.Caller.CallComponent(context.Background(), msg)
	case "remote":
		resp, err = p.options.Caller.CallServer(context.Background(), msg)
	}

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
		r.Code = 14
		r.Msg = et.Error()
		r.Namespace = et.Namespace()
	default:
		r.Code = 15
		r.Msg = et.Error()
		r.Namespace = "trellis"
	}

	// p.options.Logger.Error("call_server_failed", "trace_id", r.TraceID, "api_name", apiName, "client_ip", clientIP, "err", r)
	gCtx.JSON(200, r)
}

func (p *httpServer) getAPI(name string) (*API, bool) {
	p.syncer.RLock()
	api, ok := p.apis[name]
	p.syncer.RUnlock()
	return api, ok
}
