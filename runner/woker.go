package runner

import (
	"github.com/go-trellis/trellis/errcode"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/server"

	"github.com/go-trellis/config"
	"github.com/go-trellis/errors"
)

type Worker struct {
	opts   WorkerOptions
	server server.Server
}

type WorkerOptionFunc func(*WorkerOptions)

type WorkerOptions struct {
	url string

	name    string
	version string

	conf config.Config

	serverOptionFuncs []server.OptionFunc
}

func WorkerServer(name string, opts ...server.OptionFunc) WorkerOptionFunc {
	return func(wOpts *WorkerOptions) {
		wOpts.name = name
		wOpts.serverOptionFuncs = opts
	}
}

func WorkerVersion(ver string) WorkerOptionFunc {
	return func(wOpts *WorkerOptions) {
		wOpts.version = ver
	}
}

func (p *Worker) Stop() error {
	p.Revoke()
	return p.server.Stop()
}

func (p *Worker) Revoke() {
	RevokeWorker(p.opts.name, p.opts.version)
}

func (p *Worker) Call(msg *message.Message) errors.ErrorCode {
	err := p.server.Route(msg)(msg)
	if err != nil {
		switch t := err.(type) {
		case errors.ErrorCode:
			return t
		default:
			return errcode.ErrFailedCallServer.New(errors.Params{"err": err.Error()})
		}
	}
	return nil
}
