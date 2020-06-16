package runner

import (
	"errors"
	"fmt"

	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/service"

	"github.com/go-trellis/config"
)

// Worker 工作对象
type Worker struct {
	opts    WorkerOptions
	service service.Service
}

// WorkerOptionFunc 工作者对象函数
type WorkerOptionFunc func(*WorkerOptions)

// WorkerOptions 工作者参数
type WorkerOptions struct {
	url string

	name    string
	version string

	conf config.Config

	serviceOptionFuncs []service.OptionFunc
}

// WorkerService 工作者服务名称
func WorkerService(name string, opts ...service.OptionFunc) WorkerOptionFunc {
	return func(wOpts *WorkerOptions) {
		wOpts.name = name
		wOpts.serviceOptionFuncs = opts
	}
}

// WorkerVersion 工作者版本
func WorkerVersion(ver string) WorkerOptionFunc {
	return func(wOpts *WorkerOptions) {
		wOpts.version = ver
	}
}

func (p *Runner) newWorker(opts ...WorkerOptionFunc) error {

	workerOpts := WorkerOptions{}

	for _, o := range opts {
		o(&workerOpts)
	}

	if len(workerOpts.name) == 0 {
		return fmt.Errorf("empty servers")
	}

	if len(workerOpts.url) == 0 {
		workerOpts.url = p.genWorkerPath(workerOpts.name, workerOpts.version)
	}

	_, exist := p.workers[workerOpts.url]
	if exist {
		return fmt.Errorf("worker is already registerd, url: %s", workerOpts.url)
	}

	service, err := service.New(
		workerOpts.name,
		workerOpts.serviceOptionFuncs...,
	)

	if err != nil {
		return err
	}

	worker := &Worker{
		opts:    workerOpts,
		service: service,
	}

	p.locker.Lock()
	p.workers[workerOpts.url] = worker
	p.orderWorkers = append(p.orderWorkers, workerOpts.url)
	p.locker.Unlock()

	return nil
}

// Stop 停止工作者
func (p *Worker) Stop() error {
	RevokeWorker(p.opts.name, p.opts.version)
	return p.service.Stop()
}

// Revoke 注销工作者
func (p *Worker) Revoke() {
	RevokeWorker(p.opts.name, p.opts.version)
}

// Call 访问工作者
func (p *Worker) Call(msg *message.Message) error {

	hFunc := p.service.Route(msg)
	if hFunc == nil {
		return errors.New("not found handler")
	}
	err := hFunc(msg)
	if err != nil {
		return err
	}
	return nil
}
