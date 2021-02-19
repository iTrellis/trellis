/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

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

	if wo.Service.GetName() == "" {
		return nil, errors.New("service name not found")
	}

	ctx, cancel := context.WithCancel(context.Background())

	w := c.Watch(ctx, wo.Service.FullRegistryPath(), clientv3.WithPrefix(), clientv3.WithPrevKV())
	stop := make(chan bool, 1)

	return &etcdWatcher{
		registryID: regID,
		cancel:     cancel,
		stop:       stop,
		w:          w,
		client:     c,
		timeout:    timeout,
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
