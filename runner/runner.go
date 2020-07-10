package runner

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-trellis/common/errors"
	"github.com/go-trellis/trellis/message/proto"
	"github.com/go-trellis/trellis/service"

	"github.com/go-trellis/config"
)

// Runner 启动对象
type Runner struct {
	locker sync.RWMutex

	conf config.Config

	workers map[string]*Worker

	orderWorkers []string
}

// OptionFunc 处理参数函数
type OptionFunc func(*Options)

// Options 参数对象
type Options struct {
	Config config.Config
}

// Config 注入配置
func Config(c config.Config) OptionFunc {
	return func(p *Options) {
		p.Config = c
	}
}

var runner *Runner

// Run 运行
func Run(opts ...OptionFunc) (err error) {
	err = NewRunner(opts...)
	if err != nil {
		return
	}

	if err = runner.Run(); err != nil {
		return err
	}

	return nil
}

// NewRunner 生成启动对象
func NewRunner(opts ...OptionFunc) error {

	rOpts := &Options{}

	for _, o := range opts {
		o(rOpts)
	}

	t := &Runner{
		conf:    rOpts.Config,
		workers: make(map[string]*Worker),
	}

	if err := t.registServices(); err != nil {
		return err
	}

	runner = t

	return nil
}

// GetWorker 获取worker
func GetWorker(base *proto.BaseService) (*Worker, error) {
	runner.locker.RLock()
	defer runner.locker.RUnlock()
	worker, ok := runner.workers[runner.genWorkerPath(base.GetName(), base.GetVersion())]
	if !ok {
		return nil, fmt.Errorf("not found service: %s, %s", base.GetName(), base.GetVersion())
	}
	return worker, nil
}

func (p *Runner) registServices() error {
	servicesConf := p.conf.GetValuesConfig("project.services")
	if servicesConf == nil {
		return errors.New("services nil")
	}

	for _, name := range servicesConf.GetKeys() {
		serviceConf := servicesConf.GetValuesConfig(fmt.Sprintf("%s.options", name))
		err := p.newWorker(
			WorkerService(name, service.Config(serviceConf)),
			WorkerVersion(servicesConf.GetString(name+".version")))

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Runner) checkRegistry() {

}

// RevokeWorker 注销worker
func RevokeWorker(name, version string) {
	runner.locker.Lock()
	defer runner.locker.Unlock()
	runner.removeWorker(runner.genWorkerPath(name, version))
}

func (p *Runner) removeWorker(elt string) {
	delete(p.workers, elt)
	j := 0
	for _, val := range p.orderWorkers {
		if val == elt {
			p.orderWorkers[j] = val
			j++
		}
	}
	p.orderWorkers = p.orderWorkers[:j]
}

// Run 启动进程
func (p *Runner) Run() error {
	p.locker.Lock()
	defer p.locker.Unlock()
	for _, worker := range p.workers {
		if err := worker.service.Start(); err != nil {
			return err
		}
	}

	return nil
}

// BlockStop 阻断式停止
func BlockStop() error {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	case <-ch:
	}

	errs := errors.Errors(runner.Stop())
	if len(errs) != 0 {
		return errors.New(errs.Error())
	}
	return nil
}

// Stop 停止服务
func Stop() error {
	errs := errors.Errors(runner.Stop())
	if len(errs) != 0 {
		return errors.New(errs.Error())
	}

	return nil
}

// Stop 停止所有的Worker
func (p *Runner) Stop() []error {
	p.locker.Lock()
	defer p.locker.Unlock()
	var errs []error
	for _, worker := range p.workers {
		runner.removeWorker(runner.genWorkerPath(worker.opts.name, worker.opts.version))
		err := worker.service.Stop()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (p *Runner) genWorkerPath(name, version string) string {
	return fmt.Sprintf("trellis://workers/%s/%s", name, version)
}
