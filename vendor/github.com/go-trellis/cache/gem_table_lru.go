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

package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// LRU implements a non-thread safe fixed size LRU cache
type LRU struct {
	name string

	locker    sync.RWMutex
	size      int
	evictList *list.List
	items     map[interface{}]*list.Element
	onEvict   EvictCallback

	valueMode ValueMode
}

// NewTableCache constructs a fixed size cache.
func NewTableCache(name string, opts ...OptionFunc) (TableCache, error) {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	lru, err := NewLRU(name, options)
	if err != nil {
		return nil, err
	}

	return lru, nil
}

// NewLRU constructs an LRU of the given options
func NewLRU(name string, opts Options) (*LRU, error) {
	if opts.Size < 0 {
		return nil, errors.New("must provide a positive size, 0 is unlimit")
	}
	c := &LRU{
		name:      name,
		size:      opts.Size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
		onEvict:   opts.Evict,
		valueMode: opts.ValueMode,
	}
	return c, nil
}

// DeleteObject deletes the provided key from the cache, returning if the
// key was contained.
func (p *LRU) DeleteObject(key interface{}) (present bool) {
	p.locker.Lock()
	defer p.locker.Unlock()
	if ent, ok := p.items[key]; ok {
		p.removeElement(ent)
		return true
	}
	return false
}

// DeleteObjects is used to completely clear the cache.
func (p *LRU) DeleteObjects() {
	p.locker.Lock()
	defer p.locker.Unlock()

	for k, v := range p.items {
		if p.onEvict != nil {
			p.onEvict(k, v.Value.(*DataValues).Values)
		}
		delete(p.items, k)
	}
	p.evictList.Init()
}

// InsertExpire insert a value to the cache. Returns true if insert kv successful.
func (p *LRU) InsertExpire(key, value interface{}, expire time.Duration) bool {
	p.locker.Lock()
	defer p.locker.Unlock()

	var dv *DataValues

	// Check for existing item
	entry, ok := p.items[key]
	if ok {
		if _, ok := p.isElementExpired(entry); ok {
			dv = &DataValues{Key: key, Exists: make(map[interface{}]bool)}
		} else {
			dv = entry.Value.(*DataValues)
		}
	} else {
		dv = &DataValues{Key: key, Exists: make(map[interface{}]bool)}
	}

	switch p.valueMode {
	case ValueModeBag:
		if !dv.Exists[value] {
			dv.Values = append(dv.Values, value)
			dv.Exists[value] = true
		}
	case ValueModeDuplicateBag:
		dv.Values = append(dv.Values, value)
	case ValueModeUnique:
		fallthrough
	default:
		dv.Values = []interface{}{value}
	}

	// set expired time
	if expire > NoExpire {
		t := time.Now().Add(expire)
		dv.Expire = &t
	}

	if ok {
		p.evictList.MoveToFront(entry)
	} else {
		entry = p.evictList.PushFront(dv)
	}

	p.items[key] = entry

	evict := p.size > 0 && p.evictList.Len() > p.size
	// Verify size not exceeded
	if evict {
		p.removeOldest()
	}
	return true
}

// Insert insert a value to the cache. Returns true if an eviction occurred.
func (p *LRU) Insert(key, value interface{}) (evicted bool) {
	return p.InsertExpire(key, value, 0)
}

// Lookup Look up values with key: Key.
func (p *LRU) Lookup(key interface{}) ([]interface{}, bool) {
	p.locker.RLock()
	entry, ok := p.items[key]
	if ok {
		values, ok := p.isElementExpired(entry)
		p.locker.RUnlock()
		if ok {
			go p.DeleteObject(entry)
			return nil, false
		}
		return values, true
	}
	p.locker.RUnlock()
	return nil, false
}

func (p *LRU) isElementExpired(e *list.Element) ([]interface{}, bool) {
	dv := e.Value.(*DataValues)
	if dv.Expire != nil && dv.Expire.UnixNano() < time.Now().UnixNano() {
		return nil, true
	}
	return dv.Values, false
}

// LookupAll Look up all key-value pairs.
func (p *LRU) LookupAll() (items map[interface{}][]interface{}, ok bool) {
	p.locker.RLock()
	for k, v := range p.items {
		values, ok := p.isElementExpired(v)
		if ok {
			continue
		}

		if items == nil {
			items = make(map[interface{}][]interface{})
			ok = true
		}

		items[k] = values
	}
	p.locker.RUnlock()
	return
}

// Member Returns true if one or more elements in the table has key: Key, otherwise false.
func (p *LRU) Member(key interface{}) bool {
	p.locker.RLock()
	entry, ok := p.items[key]
	p.locker.RUnlock()
	if !ok {
		return false
	}

	_, ok = p.isElementExpired(entry)
	return !ok
}

// Members Retruns all keys in the table Tab.
func (p *LRU) Members() (keys []interface{}, ok bool) {
	p.locker.RLock()
	for k, v := range p.items {
		if _, ok := p.isElementExpired(v); ok {
			continue
		}
		keys = append(keys, k)
		ok = true
	}
	p.locker.RUnlock()
	return
}

// SetExpire Set Key Expire time
func (p *LRU) SetExpire(key interface{}, expire time.Duration) bool {
	p.locker.Lock()
	defer p.locker.Unlock()
	entry, ok := p.items[key]
	if !ok {

		return false
	}
	ent := entry.Value.(*DataValues)

	expiredTime := time.Now().Add(expire)
	ent.Expire = &expiredTime

	entry = p.evictList.PushFront(ent)
	p.items[key] = entry

	return true
}

// RemoveOldest removes the oldest item from the cache.
func (p *LRU) RemoveOldest() (key, value interface{}, ok bool) {
	p.locker.Lock()
	defer p.locker.Unlock()
	return p.removeOldest()
}

// removeOldest removes the oldest item from the cache.
func (p *LRU) removeOldest() (key, value interface{}, ok bool) {
	ent := p.evictList.Back()
	if ent != nil {
		p.removeElement(ent)
		kv := ent.Value.(*DataValues)
		return kv.Key, kv.Values, true
	}
	return nil, nil, false
}

func (p *LRU) removeElement(e *list.Element) {
	p.evictList.Remove(e)
	kv := e.Value.(*DataValues)
	delete(p.items, kv.Key)
	if p.onEvict != nil {
		p.onEvict(kv.Key, kv.Values)
	}
}
