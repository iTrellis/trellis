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
	"github.com/go-trellis/trellis/message/proto"
	"github.com/go-trellis/trellis/registry"

	"github.com/go-trellis/common/logger"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
)

// Register ETCD 注册对象
type Register struct {
	option *registry.RegistOption
	logger logger.Logger

	client *clientv3.Client

	services sync.Map
}

// New instance of server regitster
func New() registry.Registry {
	return &Register{}
}

func init() {
	registry.Regist(proto.RegisterType_ETCD, New)
}

// Init initial register
func (p *Register) Init(option *registry.RegistOption) (err error) {
	p.option = option
	p.logger = option.Logger.With("register", "etcd")

	// get endpoints for register dial address
	if p.client, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(p.option.Endpoint, ","),
		DialTimeout: p.option.TTL,
	}); err != nil {
		err = fmt.Errorf("grpclib: create clientv3 client failed: %v", err)
		p.logger.Error(err.Error())
		return
	}

	return nil
}

type worker struct {
	service *configure.RegistService
	client  *clientv3.Client

	fullServiceName string
	fullpath        string

	stopSignal chan bool

	// invoke self-register with ticker
	ticker *time.Ticker

	interval time.Duration
}

// Regist server regist into etcd
func (p *Register) Regist(s *configure.RegistService) (err error) {
	wer := &worker{
		service:    s,
		ticker:     time.NewTicker(p.option.Heartbeat),
		fullpath:   internal.WorkerETCDDomainPath(s.Name, s.Version, s.Domain),
		stopSignal: make(chan bool),
	}

	_, ok := p.services.Load(wer.fullpath)
	if ok {
		err = errors.New("service already registed")
		p.logger.Error("service", wer, "error", err)
		return err
	}

	p.logger.Debug("service", wer.fullpath)

	go func(w *worker) {
		var count uint32
		for {
			if err = p.regist(w); err != nil {
				if p.option.RetryTimes == 0 {
					continue
				}
				if p.option.RetryTimes <= count {
					panic(fmt.Errorf("%s regist into etcd failed times above: %d, %v", w.fullpath, count, err))
				}
				count++
				continue
			}
			count = 0
			select {
			case <-w.stopSignal:
				p.stop(w.fullpath)
				return
			case <-w.ticker.C:

			}
		}
	}(wer)

	p.services.Store(wer.fullpath, wer)

	return
}

func (p *Register) stop(path string) {

	ctxDel, cDel := context.WithTimeout(context.Background(), p.option.Heartbeat)
	defer cDel()
	resp, err := p.client.Delete(ctxDel, path)
	p.logger.Debug("stop", path, resp, err)
	p.client.Close()
}

func (p *Register) regist(wor *worker) (err error) {
	// minimum lease TTL is ttl-second
	ctxGrant, cGrant := context.WithTimeout(context.TODO(), p.option.TTL)
	defer cGrant()
	resp, ie := p.client.Grant(ctxGrant, int64(p.option.TTL/time.Second))
	if ie != nil {
		return fmt.Errorf("grpclib: set service %q clientv3 failed: %s", wor.fullpath, ie.Error())
	}

	ctxGet, cGet := context.WithTimeout(context.Background(), p.option.TTL)
	defer cGet()
	_, err = p.client.Get(ctxGet, wor.fullpath)
	// should get first, if not exist, set it
	if err != nil {
		if err != rpctypes.ErrKeyNotFound {
			return fmt.Errorf("grpclib: get service %q failed: %s", wor.fullpath, err.Error())
		}
	}

	// refresh set to true for not notifying the watcher
	ctxPut, cPut := context.WithTimeout(context.TODO(), p.option.TTL)
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

	fullpath := internal.WorkerETCDDomainPath(s.Name, s.Version, s.Domain)

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

	return nil
}

// Stop 结束进程
func (p *Register) Stop() {

	wg := sync.WaitGroup{}
	p.services.Range(func(key, value interface{}) bool {
		wg.Add(1)
		wer, ok := value.(*worker)
		if !ok {
			return false
		}
		go func(w *worker) {
			defer wg.Done()
			w.stopSignal <- true
		}(wer)
		return true
	})
	wg.Wait()
}

// Watcher 获取Watcher对象
func (p *Register) Watcher(option *configure.Watcher) (registry.Watcher, error) {
	return newWatcher(p, option)
}

func (p *Register) String() string {
	return "etcd"
}
