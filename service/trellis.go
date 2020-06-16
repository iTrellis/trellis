package service

// import (
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"sync"
// 	"syscall"

// 	"github.com/go-trellis/common/errors"
// 	"github.com/go-trellis/config"
// )

// // Runner 启动对象
// type Runner struct {
// 	locker sync.RWMutex

// 	conf config.Config

// 	workers map[string]*Worker

// 	orderWorkers []string
// }

// func (p *Runner) registServices() error {
// 	servicesConf := p.conf.GetMap("services")
// 	if servicesConf == nil {
// 		return nil
// 	}

// 	for _, name := range servicesConf.Keys() {
// 		serviceConf := p.conf.GetValuesConfig("services." + name)
// 		if serviceConf == nil {
// 			continue
// 		}

// 		err := p.newWorker(
// 			WorkerService(serviceConf.GetString("name", name), Config(serviceConf)),
// 			WorkerVersion(serviceConf.GetString("version", "v1")))

// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (p *Runner) genWorkerPath(name, version string) string {
// 	return fmt.Sprintf("trellis://workers/%s/%s", name, version)
// }

// func (p *Runner) newWorker(opts ...WorkerOptionFunc) error {

// 	workerOpts := WorkerOptions{}

// 	for _, o := range opts {
// 		o(&workerOpts)
// 	}

// 	if len(workerOpts.name) == 0 {
// 		return fmt.Errorf("empty servers")
// 	}

// 	if len(workerOpts.url) == 0 {
// 		workerOpts.url = p.genWorkerPath(workerOpts.name, workerOpts.version)
// 	}

// 	_, exist := p.workers[workerOpts.url]
// 	if exist {
// 		return fmt.Errorf("worker is already registerd, url: %s", workerOpts.url)
// 	}

// 	service, err := New(
// 		workerOpts.name,
// 		workerOpts.serviceOptionFuncs...,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	worker := &Worker{
// 		opts:    workerOpts,
// 		service: service,
// 	}

// 	p.locker.Lock()
// 	p.workers[workerOpts.url] = worker
// 	p.orderWorkers = append(p.orderWorkers, workerOpts.url)
// 	p.locker.Unlock()

// 	return nil
// }

// func (p *Runner) removeWorker(elt string) {
// 	delete(p.workers, elt)
// 	j := 0
// 	for _, val := range p.orderWorkers {
// 		if val == elt {
// 			p.orderWorkers[j] = val
// 			j++
// 		}
// 	}
// 	p.orderWorkers = p.orderWorkers[:j]
// }

// // Run 启动进程
// func (p *Runner) Run() error {
// 	for _, worker := range p.workers {
// 		if err := worker.service.Start(); err != nil {
// 			return err
// 		}
// 	}

// 	ch := make(chan os.Signal, 1)
// 	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

// 	select {
// 	case <-ch:
// 	}

// 	errs := errors.Errors(p.Stop())
// 	if len(errs) != 0 {
// 		return errors.New(errs.Error())
// 	}

// 	return nil
// }

// // Stop 停止所有的Worker
// func (p *Runner) Stop() []error {
// 	var errs []error
// 	for _, worker := range p.workers {
// 		e := worker.Stop()
// 		if e != nil {
// 			errs = append(errs, e)
// 		}
// 	}

// 	return errs
// }
