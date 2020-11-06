/*
Copyright © 2016 Henry Huang <hhh@rutcode.com>

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

package cache

import (
	"sort"
	"sync"
	"time"

	"github.com/go-trellis/common/formats"
)

type gemCache struct {
	sync.RWMutex

	tables map[string]TableCache
}

// New return cache manager
func New() Cache {
	return &gemCache{
		tables: make(map[string]TableCache),
	}
}

func (p *gemCache) All() []string {
	p.RLock()
	defer p.RUnlock()

	var tmpTables formats.Strings
	for k := range p.tables {
		tmpTables = append(tmpTables, k)
	}

	sort.Sort(tmpTables)
	return tmpTables
}

func (p *gemCache) GetTableCache(tab string) (TableCache, bool) {
	t, ok := p.tables[tab]
	return t, ok
}

func (p *gemCache) New(tab string, options ...OptionFunc) (err error) {
	p.Lock()
	defer p.Unlock()

	tabCache := p.getTable(tab)
	if tabCache != nil {
		return ErrTableExists
	}

	if tabCache, err = NewTableCache(tab, options...); err != nil {
		return
	}

	p.tables[tab] = tabCache

	return nil
}

func (p *gemCache) Delete(tab string) bool {
	p.Lock()
	defer p.Unlock()

	tabCache := p.getTable(tab)
	if tabCache != nil {
		tabCache.DeleteObjects()
		delete(p.tables, tab)
	}
	return true
}

func (p *gemCache) DeleteObjects(tab string) {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return
	}
	tabCache.DeleteObjects()
}

func (p *gemCache) DeleteObject(tab string, key interface{}) bool {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return true
	}
	return tabCache.DeleteObject(key)
}

func (p *gemCache) Insert(tab string, key, value interface{}) bool {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return false
	}
	return tabCache.Insert(key, value)
}

func (p *gemCache) InsertExpire(tab string, key, value interface{}, expire time.Duration) bool {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return false
	}
	return tabCache.InsertExpire(key, value, expire)
}

func (p *gemCache) getTable(tab string) TableCache {
	return p.tables[tab]
}

func (p *gemCache) Lookup(tab string, key interface{}) ([]interface{}, bool) {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return nil, false
	}
	return tabCache.Lookup(key)
}

func (p *gemCache) Member(tab string, key interface{}) bool {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return false
	}
	return tabCache.Member(key)
}

func (p *gemCache) Members(tab string) ([]interface{}, bool) {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return nil, false
	}
	return tabCache.Members()
}

func (p *gemCache) SetExpire(tab string, key interface{}, expire time.Duration) bool {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return false
	}
	return tabCache.SetExpire(key, expire)
}

func (p *gemCache) LookupAll(tab string) (map[interface{}][]interface{}, bool) {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return nil, false
	}
	return tabCache.LookupAll()
}
