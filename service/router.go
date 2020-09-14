package service

import (
	"fmt"

	"github.com/go-trellis/node"
	"github.com/go-trellis/trellis/clients"
	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/errcode"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message"

	"github.com/go-trellis/common/errors"
	"github.com/go-trellis/common/logger"
)

// todo
// local worker => service
// remote worker => http | rpc | xxx to remote service

// Router 路由器
type Router interface {
	NewService(...RouterOptionFunc) error
	StopService(name, version string) error
	Run() error
	Stop() error

	clients.Caller
}

// RouterOptionFunc 配置函数定义
type RouterOptionFunc func(*RouterOptions)

// RouterOptions 配置
type RouterOptions struct {
	cfg *configure.Service

	logger logger.Logger
}

// RouterOptionService 配置参数
func RouterOptionService(c *configure.Service) RouterOptionFunc {
	return func(w *RouterOptions) {
		w.cfg = c
	}
}

// RouterOptionLogger 日志
func RouterOptionLogger(l logger.Logger) RouterOptionFunc {
	return func(w *RouterOptions) {
		w.logger = l
	}
}

type router struct {
	opts RouterOptions

	// locker   sync.RWMutex
	services map[string]Service
}

// NewRouter gen router
func NewRouter() Router {
	return &router{
		services: make(map[string]Service),
	}
}

func (p *router) NewService(opts ...RouterOptionFunc) (err error) {

	for _, o := range opts {
		o(&p.opts)
	}

	url := internal.WorkerTrellisPath(p.opts.cfg.GetName(), p.opts.cfg.GetVersion())
	if _, ok := p.services[url]; ok {
		err = fmt.Errorf("%s already exists", url)
		p.opts.logger.Error("new_service_failed", err.Error())
		return err
	}

	s, err := New(p.opts.cfg.GetName(), p.opts.cfg.GetVersion(),
		Config(p.opts.cfg.Options),
		Logger(p.opts.logger.With(url)),
	)
	if err != nil {
		p.opts.logger.Error("new_service_failed", err.Error())
		return err
	}

	p.services[url] = s

	return nil
}

// Run 停止工作者
func (p *router) Run() error {
	var errs errors.Errors
	for _, s := range p.services {
		err := s.Start()
		if err != nil {
			errs.Append(err)
		}
	}

	if len(errs) != 0 {
		p.opts.logger.Error("run_service_failed", errs.Error())
		return errs
	}
	return nil
}

// Stop 停止工作者
func (p *router) Stop() error {
	var errs errors.Errors
	for _, s := range p.services {
		err := s.Stop()
		if err != nil {
			errs.Append(err)
		}
	}
	p.services = nil

	if len(errs) != 0 {
		p.opts.logger.Error("stop_service_failed", errs.Error())
		return errs
	}

	return nil
}

// StopService stop service
func (p *router) StopService(name, version string) error {
	url := internal.WorkerTrellisPath(name, version)
	s, ok := p.services[url]
	if !ok {
		err := fmt.Errorf("unknown service: %s, %s", name, version)
		p.opts.logger.Error("stop_service_failed", err.Error())
		return err
	}
	if err := s.Stop(); err != nil {
		p.opts.logger.Error("stop_service_failed", err.Error())
		return err
	}

	delete(p.services, url)
	return nil
}

// GetService get service
func (p *router) GetService(name, version string) (Service, error) {
	s, ok := p.services[internal.WorkerTrellisPath(name, version)]
	if !ok {
		return nil, fmt.Errorf("unknown service: %s, %s", name, version)
	}
	return s, nil
}

// CallService call service
func (p *router) CallService(_ *node.Node, msg *message.Message) (interface{}, error) {
	s, err := p.GetService(msg.GetService().GetName(), msg.GetService().GetVersion())
	if err != nil {
		return nil, err
	}

	fn := s.Route(msg.GetTopic())
	if fn == nil {
		return nil, errcode.ErrGetServiceTopic.New()
	}
	return fn(msg)
}
