package runner

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-trellis/config"
	"github.com/go-trellis/errors"

	"github.com/go-trellis/trellis/server"
)

type Runner struct {
	locker sync.RWMutex

	conf config.Config

	workers map[string]*Worker

	orderWorkers []string
}

type OptionFunc func(*Options)

type Options struct {
	conf config.Config
}

func Config(conf config.Config) OptionFunc {
	return func(p *Options) {
		p.conf = conf
	}
}

func (p *Runner) registServers() error {
	serversConf := p.conf.GetMap("servers")
	if serversConf == nil {
		return nil
	}

	for _, name := range serversConf.Keys() {
		serverConf := p.conf.GetValuesConfig("servers." + name)
		if serverConf == nil {
			continue
		}

		err := p.newWorker(
			WorkerServer(serverConf.GetString("name", name), server.Config(serverConf)),
			WorkerVersion(serverConf.GetString("version", "v1")))

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Runner) genWorkerPath(name, version string) string {
	return fmt.Sprintf("trellis://workers/%s/%s", name, version)
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

	server, err := server.New(
		workerOpts.name,
		workerOpts.serverOptionFuncs...,
	)

	if err != nil {
		return err
	}

	// warnNoDocsComp(name, actOpts.componentDriver, comp)

	worker := &Worker{
		opts:   workerOpts,
		server: server,
	}

	p.locker.Lock()
	p.workers[workerOpts.url] = worker
	p.orderWorkers = append(p.orderWorkers, workerOpts.url)
	p.locker.Unlock()

	return nil
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
	for _, worker := range p.workers {
		if err := worker.server.Start(); err != nil {
			return err
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	case <-ch:
	}

	errs := errors.ErrsString(p.Stop())
	if len(errs) != 0 {
		return fmt.Errorf(errs)
	}

	return nil
}

// Stop 停止所有的Worker
func (p *Runner) Stop() []error {
	var errs []error
	for _, worker := range p.workers {
		e := worker.Stop()
		if e != nil {
			errs = append(errs, e)
		}
	}

	return errs
}
