package route

import "github.com/go-trellis/trellis/service/registry"

// todo

type Watcher struct {
	opts registry.WatchOptions
}

func (p *Watcher) Next() (*registry.Result, error) {
	return nil, nil
}

func (p *Watcher) Stop() {

}
