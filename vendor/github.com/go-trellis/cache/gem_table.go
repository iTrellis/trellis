// GNU GPL v3 License
// Copyright (c) 2016 github.com:go-trellis

package cache

import (
	"sort"
	"sync"
	"time"

	"github.com/go-trellis/common/formats"
)

type table struct {
	name string
	sync.RWMutex

	expiredMutex sync.Mutex
	expiredTimer time.Duration
	expiredKeys  map[string]*time.Time

	items map[string]*DataValues

	orderSetKeys formats.Strings

	backgroundExpiredFunc        func()
	backgroundExpiredFuncRunning bool

	orderSets bool
	valueMode ValueMode
}

// Option 参数处理函数
type Option func(*table)

// OptionValueMode 设置值的数据方式
func OptionValueMode(mode ValueMode) Option {
	return func(t *table) {
		t.valueMode = mode
	}
}

// OptionOrderedSet 设置是否排序
func OptionOrderedSet() Option {
	return func(t *table) {
		t.orderSets = true
	}
}

// NewTableCache return table cache with input options
func NewTableCache(name string, options ...Option) (TableCache, error) {

	p := &table{
		name:         name,
		expiredKeys:  make(map[string]*time.Time),
		items:        make(map[string]*DataValues),
		expiredTimer: DefaultTimer,
	}

	if err := p.init(options...); err != nil {
		return nil, err
	}

	p.backgroundExpiredFunc = p.backgroundDeleteExpiredKey
	go p.backgroundExpiredFunc()

	return p, nil
}

func (p *table) init(options ...Option) error {
	for _, o := range options {
		o(p)
	}
	switch p.valueMode {
	case ValueModeUnique, ValueModeDuplicateBag, ValueModeBag:
		return nil
	default:
		return ErrUnknownTableValueMode
	}
}

func (p *table) Insert(key string, value interface{}) bool {
	return p.InsertExpire(key, value, NoExpire)
}

func (p *table) InsertExpire(key string, value interface{}, expire time.Duration) bool {

	if key == "" || value == nil || expire < 0 {
		return false
	}
	// force to delete old expired data values
	p.mutexExpiredDelete(key)

	// set key value with expire or noexpire
	return p.setKeyValue(key, value, expire)
}

func (p *table) Lookup(key string) ([]interface{}, bool) {
	p.RLock()
	defer p.RUnlock()
	if p.isKeyExpired(key) {
		return nil, false
	}
	item, ok := p.getKeyValue(key)
	if !ok {
		return nil, false
	}
	return item.Values, true
}

func (p *table) LookupAll() (items map[string][]interface{}, ok bool) {
	p.RLock()
	for k, v := range p.items {
		if p.isKeyExpired(k) {
			continue
		}
		if items == nil {
			items = make(map[string][]interface{}, 1)
		}
		items[k] = v.Values
		ok = true
	}
	p.RUnlock()
	return
}

func (p *table) LookupLimit(pos, limit uint) (items map[string][]interface{}, ok bool) {
	p.RLock()
	if !p.orderSets {
		p.RUnlock()
		return nil, false
	}
	lenOrderSetKeys, count := len(p.orderSetKeys), limit
	if limit == 0 {
		count = uint(lenOrderSetKeys) - pos
	}
	for i := int(pos); i < lenOrderSetKeys; i++ {
		k := p.orderSetKeys[i]
		if p.isKeyExpired(k) {
			continue
		}
		if items == nil {
			items = make(map[string][]interface{}, 1)
		}

		items[k] = p.items[k].Values
		ok = true

		if count--; count == 0 {
			break
		}
	}
	p.RUnlock()
	return
}

func (p *table) isKeyExpired(key string) bool {
	if expiredTime := p.expiredKeys[key]; expiredTime != nil {
		return time.Now().After(*expiredTime)
	}
	return false
}

func (p *table) Member(key string) bool {
	p.RLock()
	defer p.RUnlock()
	if p.isKeyExpired(key) {
		return false
	}
	_, ok := p.getKeyValue(key)
	return ok
}

func (p *table) Members() ([]string, bool) {
	p.RLock()
	s := formats.Strings{}
	ok := false
	for k := range p.items {
		if p.isKeyExpired(k) {
			continue
		}
		s = append(s, k)
	}
	if len(s) != 0 {
		ok = true
	}

	if p.orderSets {
		sort.Sort(s)
	}

	p.RUnlock()
	return s, ok
}

func (p *table) DeleteObject(key string) bool {
	p.mutexDelete(key)
	return true
}

func (p *table) DeleteObjects() bool {
	p.Lock()
	p.items = make(map[string]*DataValues)
	p.expiredKeys = make(map[string]*time.Time)
	p.Unlock()
	return true
}

func (p *table) SetExpire(key string, expire time.Duration) bool {
	p.RLock()
	p.setExpire(key, expire)
	p.RUnlock()
	return true
}

func (p *table) setExpire(key string, expire time.Duration) {
	if expire == NoExpire {
		return
	}
	t := time.Now().Add(expire)
	p.expiredKeys[key] = &t
}

func (p *table) setKeyValue(key string, value interface{}, expire time.Duration) bool {

	p.Lock()
	defer p.Unlock()

	p.setExpire(key, expire)

	item, ok := p.getKeyValue(key)
	if !ok {
		item = &DataValues{Key: key}
	}
	switch p.valueMode {
	case ValueModeUnique:
		{
			prepend := make([]interface{}, 1)
			prepend[0] = value

			item.Length = 1
			item.Values = prepend
		}
	case ValueModeBag:
		{
			if item.Exists == nil {
				item.Exists = make(map[interface{}]bool, 1)
			}

			if !item.Exists[value] {
				item.Exists[value] = true
				item.Values = append(item.Values, value)
				item.Length++
			}
		}
	case ValueModeDuplicateBag:
		{
			item.Values = append(item.Values, value)
			item.Length++
		}
	default:
		return false
	}
	p.items[key] = item

	p.sortOrderSetKeys()

	return true
}

func (p *table) getKeyValue(key string) (item *DataValues, ok bool) {
	item, ok = p.items[key]
	return
}

func (p *table) mutexDelete(key string) {
	p.Lock()
	p.delete(key)
	p.Unlock()
}

func (p *table) delete(key string) {
	delete(p.items, key)
	delete(p.expiredKeys, key)
	p.sortOrderSetKeys()
}

func (p *table) mutexExpiredDelete(key string) {
	p.Lock()
	p.expiredDelete(key)
	p.Unlock()
}

func (p *table) expiredDelete(key string) {
	if p.isKeyExpired(key) {
		p.delete(key)
	}
}

func (p *table) sortOrderSetKeys() {
	if !p.orderSets {
		return
	}
	p.orderSetKeys = nil
	for k := range p.items {
		p.orderSetKeys = append(p.orderSetKeys, k)
	}

	sort.Sort(p.orderSetKeys)
}

func (p *table) SetBackgroundExpiredTime(t time.Duration) {
	p.expiredTimer = t
}

func (p *table) backgroundDeleteExpiredKey() {
	for {
		func(timer time.Duration) {
			if p.backgroundExpiredFuncRunning {
				return
			}
			p.backgroundExpiredFuncRunning = true
			defer func() {
				p.backgroundExpiredFuncRunning = false
			}()

			p.expiredMutex.Lock()
			for key := range p.expiredKeys {
				p.mutexExpiredDelete(key)
			}
			p.expiredMutex.Unlock()
			time.Sleep(timer)
		}(p.expiredTimer)
	}
}
