package etcd

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/registry"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdWatcher struct {
	registryID string

	w       clientv3.WatchChan
	client  *clientv3.Client
	timeout time.Duration

	sync.Mutex

	cancel func()
	stop   chan bool
}

func newEtcdWatcher(c *clientv3.Client, regID string, timeout time.Duration, opts ...registry.WatchOption) (
	registry.Watcher, error) {

	var wo registry.WatchOptions
	for _, o := range opts {
		o(&wo)
	}

	if wo.Service.Name == "" {
		return nil, errors.New("service name not found")
	}

	ctx, cancel := context.WithCancel(context.Background())

	w := c.Watch(ctx, wo.Service.FullRegistry(), clientv3.WithPrefix(), clientv3.WithPrevKV())
	stop := make(chan bool, 1)

	return &etcdWatcher{
		cancel:  cancel,
		stop:    stop,
		w:       w,
		client:  c,
		timeout: timeout,
	}, nil
}

func (p *etcdWatcher) Next() (*registry.Result, error) {
	for resp := range p.w {
		if resp.Err() != nil {
			return nil, resp.Err()
		}

		if resp.Canceled {
			return nil, errors.New("could not get next")
		}

		for _, ev := range resp.Events {
			s := decode(ev.Kv.Value)
			var typ service.EventType

			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					typ = service.EventType_create
				} else if ev.IsModify() {
					typ = service.EventType_update
				}
			case clientv3.EventTypeDelete:
				typ = service.EventType_delete

				// get service from prevKv
				s = decode(ev.PrevKv.Value)
			}

			if s == nil {
				continue
			}
			return &registry.Result{
				ID:        p.registryID,
				Type:      typ,
				Timestamp: time.Now(),
				Service:   s,
			}, nil
		}
	}
	return nil, errors.New("could not get next")
}

func (p *etcdWatcher) Stop() {
	p.Lock()
	defer p.Unlock()

	select {
	case <-p.stop:
		return
	default:
		close(p.stop)
		p.cancel()
		p.client.Close()
	}
}
