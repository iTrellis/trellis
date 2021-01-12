package memory

import (
	"errors"

	"github.com/go-trellis/trellis/service/registry"
)

// Watcher watcher
type Watcher struct {
	id   string
	wo   registry.WatchOptions
	exit chan bool
	res  chan *registry.Result
}

// Next watch the regstry result
func (p *Watcher) Next() (*registry.Result, error) {
	for {
		select {
		case r := <-p.res:
			if p.wo.Service.Name != "" &&
				p.wo.Service.FullPath() != r.Service.Service.FullPath() {
				continue
			}
			return r, nil
		case <-p.exit:
			return nil, errors.New("watcher stopped")
		}
	}
}

// Stop stop watcher
func (p *Watcher) Stop() {
	select {
	case <-p.exit:
		return
	default:
		close(p.exit)
	}
}
