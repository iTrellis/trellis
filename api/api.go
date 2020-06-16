package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-trellis/trellis/codec"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/router"
	"github.com/go-trellis/trellis/runner"
	"github.com/go-trellis/trellis/service"
	"github.com/go-trellis/trellis/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-trellis/common/errors"
)

func init() {
	service.RegistNewServiceFunc("trellis-api", NewAPIService)
}

type APIService struct {
	debug bool
	opts  service.Options

	forwardHeaders []string

	srv *http.Server
}

func NewAPIService(opts ...service.OptionFunc) (service.Service, error) {

	s := &APIService{}

	for _, o := range opts {
		o(&s.opts)
	}

	err := s.init()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (p *APIService) init() (err error) {

	p.debug = p.opts.Config.GetBoolean("debug", false)
	if !p.debug {
		gin.SetMode("release")
	}

	httpConf := p.opts.Config.GetConfig("http")

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

	fmt.Println(p.srv.Addr)
	return
}

func (p *APIService) Start() error {
	go func() {

		var err error

		sslConf := p.opts.Config.GetConfig("http.ssl")

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

func (p *APIService) Stop() error {

	dur := p.opts.Config.GetTimeDuration("http.shutdown-timeout", time.Second*30)

	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	if err := p.srv.Shutdown(ctx); err != nil {
		return errors.Newf("api shutdown failure, err: %s", err)
	}
	return nil
}

func (p *APIService) serve(ctx *gin.Context) {
	msg := &message.Message{}
	r := &message.Response{}
	var err error
	msg.Codec, err = codec.GetCodec(ctx.Request.Header.Get("content-type"))
	if err != nil {
		r.SetError(err)
		ctx.JSON(200, r)
		return
	}
	switch msg.Codec.String() {
	case codec.JSON:
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(ctx.Request.Body)
		if err != nil {
			r.SetError(err)
			ctx.JSON(200, r)
			return
		}
		err = msg.UnmarshalJSON(body.Bytes())
		if err != nil {
			r.SetError(err)
			ctx.JSON(200, r)
			return
		}
	}

	worker, err := runner.GetWorker(msg.GetService())
	if err != nil {
		r.SetError(err)
		ctx.JSON(200, r)
		return
	}

	fmt.Println(worker)

	err = worker.Call(msg)
	if err != nil {
		r.SetError(err)
		ctx.JSON(200, r)
		return
	}

	r.SetBody(msg)

	ctx.JSON(200, msg)
}

// Route 路由
func (p *APIService) Route(msg *message.Message) router.HandlerFunc {
	// async中处理callback
	return nil
}
