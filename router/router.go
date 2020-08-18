package router

import (
	"fmt"

	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/service"

	"github.com/go-trellis/common/errors"
	"github.com/go-trellis/common/logger"
)

// todo
// local worker => service
// remote worker => http | rpc | xxx to remote service

// Router 路由器
type Router interface {
	NewService(...OptionFunc) error
	StopService(name, version string) error
	GetService(name, version string) (service.Service, error)
	Run() error
	Stop() error
}

// OptionFunc 配置函数定义
type OptionFunc func(*Options)

// Options 配置
type Options struct {
	cfg *configure.Service

	logger logger.Logger
}

// OptionService 配置参数
func OptionService(c *configure.Service) OptionFunc {
	return func(w *Options) {
		w.cfg = c
	}
}

// OptionLogger 日志
func OptionLogger(l logger.Logger) OptionFunc {
	return func(w *Options) {
		w.logger = l
	}
}

type worker struct {
	opts Options

	// locker   sync.RWMutex
	services map[string]service.Service
}

// NewRouter gen router
func NewRouter() Router {
	return &worker{
		services: make(map[string]service.Service),
	}
}

func (p *worker) NewService(opts ...OptionFunc) (err error) {

	for _, o := range opts {
		o(&p.opts)
	}

	url := internal.WorkerPath(internal.SchemaTrellis, p.opts.cfg.GetName(), p.opts.cfg.GetVersion())
	if _, ok := p.services[url]; ok {
		err = fmt.Errorf("%s already exists", url)
		p.opts.logger.Error("new_service_failed", err.Error())
		return err
	}

	var s service.Service

	s, err = service.New(p.opts.cfg.GetName(), p.opts.cfg.GetVersion(),
		service.Config(p.opts.cfg.Options),
		service.Logger(p.opts.logger.With(url)),
	)
	if err != nil {
		p.opts.logger.Error("new_service_failed", err.Error())
		return err
	}

	p.services[url] = s

	return nil
}

// Run 停止工作者
func (p *worker) Run() error {
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
func (p *worker) Stop() error {
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
func (p *worker) StopService(name, version string) error {
	url := internal.WorkerPath(internal.SchemaTrellis, name, version)
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
func (p *worker) GetService(name, version string) (service.Service, error) {
	s, ok := p.services[internal.WorkerPath(internal.SchemaTrellis, name, version)]
	if !ok {
		return nil, fmt.Errorf("unknown service: %s, %s", name, version)
	}
	return s, nil
}
