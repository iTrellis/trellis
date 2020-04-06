package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-trellis/trellis/router"
	"github.com/go-trellis/trellis/runner"

	"github.com/go-trellis/trellis/codec"
	"github.com/go-trellis/trellis/errcode"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/server"
	"github.com/go-trellis/trellis/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-trellis/errors"
)

func init() {
	server.RegistNewServerFunc("trellis-api", NewAPIServer)
}

type APIServer struct {
	debug bool
	opts  server.Options

	forwardHeaders []string

	srv *http.Server
}

func NewAPIServer(opts ...server.OptionFunc) (server.Server, error) {

	s := &APIServer{}

	for _, o := range opts {
		o(&s.opts)
	}

	err := s.init()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (p *APIServer) init() (err error) {

	p.debug = p.opts.Config.GetBoolean("debug", false)
	if !p.debug {
		gin.SetMode("release")
	}

	httpConf := p.opts.Config.GetConfig("http")

	urlPath := httpConf.GetString("path", "/")

	router := gin.New()

	router.Use(gin.Recovery())

	address := httpConf.GetString("address", ":8080")

	forwardHeaders := httpConf.GetStringList("forward.headers")

	p.forwardHeaders = forwardHeaders

	utils.LoadCors(router, httpConf.GetConfig("cors"))
	utils.LoadPprof(router, httpConf.GetConfig("pprof"))

	// router.ServeHTTP()
	router.POST(urlPath, p.serve)

	p.srv = &http.Server{
		Addr:    address,
		Handler: router,
	}
	return
}

func (p *APIServer) Start() error {
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

func (p *APIServer) Stop() error {

	dur := p.opts.Config.GetTimeDuration("http.shutdown-timeout", time.Second*30)

	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	if err := p.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("postapi shutdown failure, err: %s", err)
	}
	fmt.Println("api stop")
	return nil
}

func (p *APIServer) serve(ctx *gin.Context) {
	msg := &message.Message{}

	ct := ctx.Request.Header.Get("content-type")
	r := &message.Response{}
	switch ct {
	case codec.JSON:
		err := ctx.BindJSON(msg)
		if err != nil {
			r.SetError(err)
			ctx.JSON(200, r)
			return
		}
	default:
		r.SetError(errcode.ErrUnsupportedContentType.New(errors.Params{"err": ct}))
		ctx.JSON(200, r)
		return
	}

	worker, err := runner.GetWorker(msg.GetServer())
	if err != nil {
		r.SetError(err)
		ctx.JSON(200, r)
		return
	}

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
func (p *APIServer) Route(msg *message.Message) router.HandlerFunc {
	// async中处理callback
	return nil
}
