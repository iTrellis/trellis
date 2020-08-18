// GNU GPL v3 License
// Copyright (c) 2016 github.com:go-trellis

package cache

import "time"

// Cache Manager functions for executing k-v tables base on TableCache
type Cache interface {
	// Returns a list of all tables at the node.
	All() []string
	// Get TableCache
	GetTableCache(tab string) (TableCache, bool)
	// Creates a new table.
	New(tab string, options ...Option) error
	// Inserts the object or all of the objects in list.
	Insert(tab, key string, value interface{}) bool
	// Inserts the object or all of the objects with expired time in list.
	InsertExpire(tab, key string, value interface{}, expire time.Duration) bool
	// Deletes the entire table Tab.
	Delete(tab string) bool
	// Deletes all objects with key, Key from table Tab.
	DeleteObject(tab, key string) bool
	// Delete all objects in the table Tab. Remain table in cache.
	DeleteAllObjects(tab string) bool
	// Look up values with key, Key from table Tab.
	Lookup(tab, key string) ([]interface{}, bool)
	// Look up all values in the Tab.
	LookupAll(tab string) (map[string][]interface{}, bool)
	// Look up pos to limit, table order set must be true
	// if limit equals 0, to table's end
	LookupLimit(tab string, pos, limit uint) (map[string][]interface{}, bool)
	// Returns true if one or more elements in the table has key Key, otherwise false.
	Member(tab, key string) bool
	// Retruns all keys in the table Tab.
	Members(tab string) ([]string, bool)
	// Set key Key expire time in the table Tab.
	SetExpire(tab, key string, expire time.Duration) bool
}
