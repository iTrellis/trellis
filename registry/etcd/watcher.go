package etcd

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/registry"

	"go.etcd.io/etcd/clientv3"
)

type watcher struct {
	option registry.WatchOption

	client *clientv3.Client

	fullpath string

	ctx    context.Context
	cancel context.CancelFunc

	watchChan clientv3.WatchChan
}

func newWatcher(r *Register, opts registry.WatchOption) (registry.Watcher, error) {
	w := &watcher{
		client: r.client,
		option: opts,
	}

	w.ctx, w.cancel = context.WithTimeout(context.Background(), r.option.TTL)

	w.fullpath = internal.WorkerPath(
		internal.SchemaETCDNaming, w.option.Service.GetName(), w.option.Service.GetVersion())

	w.watchChan = w.client.Watch(w.ctx, w.fullpath, clientv3.WithPrefix(), clientv3.WithPrevKV())

	return w, w.ctx.Err()
}

func (p *watcher) Stop() {
	p.cancel()
}

func (p *watcher) Next() (*registry.Result, error) {

	for {
		select {
		case wresp := <-p.watchChan:
			if wresp.Err() != nil {
				return nil, wresp.Err()
			}

			if wresp.Canceled {
				return nil, errors.New("watcher was canceled")
			}

			for _, ev := range wresp.Events {

				var action string
				var services configure.RegistServices

				switch ev.Type {
				case clientv3.EventTypePut:
					services, _ = p.decode(ev.Kv.Value)
					if ev.IsCreate() {
						action = registry.ActionCreate
					} else if ev.IsModify() {
						action = registry.ActionUpdate
					}
				case clientv3.EventTypeDelete:
					action = registry.ActionDelete
					services, _ = p.decode(ev.PrevKv.Value)
				}

				if services == nil {
					continue
				}

				return &registry.Result{
					Action:  action,
					Service: services,
				}, nil
			}
		}
	}
}

func (p *watcher) decode(bs []byte) (configure.RegistServices, error) {
	ss := configure.RegistServices{}
	if err := json.Unmarshal(bs, &ss); err != nil {
		return nil, err
	}
	return ss, nil
}
