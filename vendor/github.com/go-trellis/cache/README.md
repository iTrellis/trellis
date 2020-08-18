# cache
light cache in go 

* [![GoDoc](http://godoc.org/github.com/go-trellis/cache?status.svg)](http://godoc.org/github.com/go-trellis/cache)
* [![Build Status](https://travis-ci.org/go-trellis/cache.png)](https://travis-ci.org/go-trellis/cache)

## Introduction

### Installation

```shell
go get "github.com/go-trellis/formats/inner-types"
go get "github.com/go-trellis/cache"
```

### Example

[See test file](table_cache_test.go)

### Features

#### Cache

cache is manager for k-vs tables base on TableCache

```go
// Cache Manager functions for executing k-v tables base on TableCache
type Cache interface {
	// Returns a list of all tables at the node.
	All() []string
	// Get TableCache
	GetTableCache(tab string) (TableCache, bool)
	// Creates a new table.
	New(tab string, options TableOptions) error
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
```

#### TableCache

table cache is manager for k-vs

```golang
// TableCache
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
```

#### Sample: NewTableCache with options

```go
// TableOptionSets
const (
	// OrederSet: true | false
	TableOptionOrderedSet = "ordered_set"
	// ValueMode
	TableOptionValueMode = "value_mode"
)

// ValueMode
const (
	// only one value
	ValueModeUnique ValueMode = iota
	// The table is a bag table, which can have many objects
	// but only one instance of each object, per key.
	ValueModeBag
	// The table is a duplicate_bag table, which can have many objects,
	// including multiple copies of the same object, per key.
	ValueModeDuplicateBag
)

	// ValueModeUnique table
	Cache.New("tab1", nil)
	// ValueModeDuplicateBag table
	Cache.New(tab2, cache.TableOptions{cache.TableOptionValueMode: cache.ValueModeDuplicateBag})
	// ValueModeBag table
	Cache.New(tab3, cache.TableOptions{cache.TableOptionValueMode: cache.ValueModeBag})
	// ValueModeUnique with keys order set (TableOptionOrderedSet) 
	Cache.New(tab4, cache.TableOptions{cache.TableOptionOrderedSet: true})
	// ValueModeBag with keys order set (TableOptionOrderedSet) 
	Cache.New(tab4, cache.TableOptions{
		cache.TableOptionValueMode: cache.ValueModeBag, 
		cache.TableOptionOrderedSet: true})
```
