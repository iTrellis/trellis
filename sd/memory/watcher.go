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

package memory

import (
	"errors"

	"github.com/iTrellis/trellis/service/registry"
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
				p.wo.Service.FullRegistryPath() != r.Service.Service.FullRegistryPath() {
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
