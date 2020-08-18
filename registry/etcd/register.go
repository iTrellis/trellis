package etcd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/registry"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
)

type Register struct {
	option registry.RegistOption

	client *clientv3.Client

	services sync.Map
}

// NewRegister instance of server regitster
func NewRegister(option registry.RegistOption) (register registry.Registry, err error) {
	r := &Register{
		option: option,
	}
	// get endpoints for register dial address
	if r.client, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(option.Endpoint, ","),
		DialTimeout: option.TTL,
	}); err != nil {
		return nil, fmt.Errorf("grpclib: create clientv3 client failed: %v", err)
	}
	return r, nil
}

type worker struct {
	service *configure.RegistService
	client  *clientv3.Client

	fullServiceName string
	fullpath        string

	// client *clientv3.Client

	stopSignal chan bool

	// invoke self-register with ticker
	ticker *time.Ticker

	interval time.Duration

	retryTimes uint32
}

// Regist server regist into etcd
func (p *Register) Regist(service *configure.RegistService) (err error) {
	wer := &worker{
		service: service,
		ticker:  time.NewTicker(p.option.Heartbeat),
	}
	wer.fullpath =
		internal.WorkerDomainPath(internal.SchemaETCDNaming, wer.service.Name, wer.service.Version, wer.service.Domain)

	_, ok := p.services.LoadOrStore(wer.fullpath, wer)
	if ok {
		return errors.New("service already registed")
	}

	go func(w *worker) {
		for {
			if err = p.regist(w); err != nil {
				if p.option.RetryTimes < 0 {
					continue
				}
				fmt.Println(err)

				w.retryTimes++
				if p.option.RetryTimes < w.retryTimes {
					panic(fmt.Errorf("%s regist into etcd failed times above: %d, %v",
						w.fullpath, w.retryTimes, err))
				}

				continue
			}
			w.retryTimes = 0
			select {
			case <-w.stopSignal:
				w.client.Close()
				return
			case <-w.ticker.C:
			}
		}
	}(wer)

	return
}

func (p *Register) regist(wor *worker) (err error) {
	// minimum lease TTL is ttl-second
	ctxGrant, cGrant := context.WithTimeout(context.TODO(), p.option.Heartbeat)
	defer cGrant()
	resp, ie := p.client.Grant(ctxGrant, int64(p.option.TTL/time.Second))
	if ie != nil {
		return fmt.Errorf("grpclib: set service %q clientv3 failed: %s", wor.fullpath, ie.Error())
	}

	ctxGet, cGet := context.WithTimeout(context.Background(), p.option.Heartbeat)
	defer cGet()
	_, err = p.client.Get(ctxGet, wor.fullpath)
	// should get first, if not exist, set it
	if err != nil {
		if err != rpctypes.ErrKeyNotFound {
			return fmt.Errorf("grpclib: get service %q failed: %s", wor.fullpath, err.Error())
		}
	}

	// refresh set to true for not notifying the watcher
	ctxPut, cPut := context.WithTimeout(context.TODO(), p.option.Heartbeat)
	defer cPut()
	if _, err = p.client.Put(ctxPut, wor.fullpath, wor.service.String(), clientv3.WithLease(resp.ID)); err != nil {
		return fmt.Errorf("grpclib: refresh service %q failed: %s", wor.fullpath, err.Error())
	}
	return
}

// Revoke 注销服务
func (p *Register) Revoke(s *configure.RegistService) error {

	if len(s.Domain) == 0 {
		return errors.New("domain is empty")
	}

	fullpath :=
		internal.WorkerDomainPath(internal.SchemaETCDNaming, s.Name, s.Version, s.Domain)

	lWorker, ok := p.services.Load(fullpath)
	if !ok {
		return nil
	}
	p.services.Delete(fullpath)

	wor, ok := lWorker.(*worker)
	if !ok {
		return nil
	}

	wor.stopSignal <- true

	ctx, cancel := context.WithTimeout(context.Background(), p.option.Heartbeat)
	defer cancel()

	if _, err := p.client.Delete(ctx, fullpath); err != nil {
		return err
	}

	return nil
}

// Stop 结束进程
func (p *Register) Stop() {

	wg := sync.WaitGroup{}
	p.services.Range(func(key, value interface{}) bool {
		wg.Add(1)
		wer := value.(*worker)
		go func(w *worker) {
			defer wg.Done()
			w.stopSignal <- true
		}(wer)
		return true
	})
	wg.Wait()
}

// Watch 获取watch对象
func (p *Register) Watch(option registry.WatchOption) (registry.Watcher, error) {
	return newWatcher(p, option)
}

func (p *Register) String() string {
	return "etcd"
}
