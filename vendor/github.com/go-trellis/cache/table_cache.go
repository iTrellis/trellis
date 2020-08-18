// GNU GPL v3 License

// Copyright (c) 2016 github.com:go-trellis

package cache

import (
	"time"
)

// Timers
const (
	NoExpire time.Duration = 0

	DefaultTimer = time.Second * 30
)

// TableOptionSets
const (
	// OrederSet: true | false
	TableOptionOrderedSet = "ordered_set"
	// ValueMode
	TableOptionValueMode = "value_mode"
)

// TableCache table manager for k-vs functions
type TableCache interface {
	// Inserts the object or all of the objects in list.
	Insert(key string, values interface{}) bool
	// Inserts the object or all of the objects with expired time in list.
	InsertExpire(key string, value interface{}, expire time.Duration) bool
	// Deletes all objects with key: Key.
	DeleteObject(key string) bool
	// Delete all objects in the table Tab. Remain table in cache.
	DeleteObjects() bool
	// Returns true if one or more elements in the table has key: Key, otherwise false.
	Member(key string) bool
	// Retruns all keys in the table Tab.
	Members() ([]string, bool)
	// Look up values with key: Key.
	Lookup(key string) ([]interface{}, bool)
	// Look up all values in the Tab.
	LookupAll() (map[string][]interface{}, bool)
	// Look up pos to limit, table order set must be true
	// if limit equals 0, to table's end
	LookupLimit(pos, limit uint) (map[string][]interface{}, bool)
	// Set Key Expire time
	SetExpire(key string, expire time.Duration) bool
	// Set background expired time, default: 30s
	SetBackgroundExpiredTime(time.Duration)
}
