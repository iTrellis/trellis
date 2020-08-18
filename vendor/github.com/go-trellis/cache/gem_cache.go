// GNU GPL v3 License

// Copyright (c) 2016 github.com:go-trellis

package cache

import (
	"fmt"
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

func (p *gemCache) New(tab string, options ...Option) (err error) {
	p.Lock()
	defer p.Unlock()
	fmt.Println("in", tab)

	tabCache := p.getTable(tab)
	if tabCache != nil {
		fmt.Println(tab, tabCache)
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

func (p *gemCache) DeleteAllObjects(tab string) bool {

	tabCache := p.getTable(tab)
	if tabCache == nil {
		return true
	}
	return tabCache.DeleteObjects()
}

func (p *gemCache) DeleteObject(tab, key string) bool {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return true
	}
	return tabCache.DeleteObject(key)
}

func (p *gemCache) Insert(tab, key string, value interface{}) bool {

	tabCache := p.getTable(tab)
	if tabCache == nil {
		return false
	}

	return tabCache.Insert(key, value)
}

func (p *gemCache) InsertExpire(tab, key string, value interface{}, expire time.Duration) bool {

	tabCache := p.getTable(tab)
	if tabCache == nil {
		return false
	}

	return tabCache.InsertExpire(key, value, expire)
}

func (p *gemCache) getTable(tab string) TableCache {
	return p.tables[tab]
}

func (p *gemCache) Lookup(tab, key string) ([]interface{}, bool) {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return nil, false
	}
	return tabCache.Lookup(key)
}

func (p *gemCache) Member(tab, key string) bool {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return false
	}

	return tabCache.Member(key)
}

func (p *gemCache) Members(tab string) ([]string, bool) {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return nil, false
	}

	return tabCache.Members()
}

func (p *gemCache) SetExpire(tab, key string, expire time.Duration) bool {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return false
	}

	return tabCache.SetExpire(key, expire)
}

func (p *gemCache) LookupAll(tab string) (map[string][]interface{}, bool) {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return nil, false
	}

	return tabCache.LookupAll()
}

func (p *gemCache) LookupLimit(tab string, pos, limit uint) (map[string][]interface{}, bool) {
	tabCache := p.getTable(tab)
	if tabCache == nil {
		return nil, false
	}
	return tabCache.LookupLimit(pos, limit)
}
