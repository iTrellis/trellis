# cache
light cache in go 

* [![GoDoc](http://godoc.org/github.com/go-trellis/cache?status.svg)](http://godoc.org/github.com/go-trellis/cache)

## Introduction

### Installation

```shell
go get "github.com/go-trellis/cache"
```

### Features

* Simple lru
* It can set Unique | Bag | DuplicateBag values per key

### TODO

* main node: to manage cache
* consistent hash to several nodes to install keys

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
	New(tab string, options ...OptionFunc) error
	// Inserts the object or all of the objects in list.
	Insert(tab string, key, value interface{}) bool
	// Inserts the object or all of the objects with expired time in list.
	InsertExpire(tab string, key, value interface{}, expire time.Duration) bool
	// Deletes the entire table Tab.
	Delete(tab string) bool
	// Deletes all objects with key, Key from table Tab.
	DeleteObject(tab string, key interface{}) bool
	// Delete all objects in the table Tab. Remain table in cache.
	DeleteObjects(tab string)
	// Look up values with key, Key from table Tab.
	Lookup(tab string, key interface{}) ([]interface{}, bool)
	// Look up all values in the Tab.
	LookupAll(tab string) (map[interface{}][]interface{}, bool)
	// Returns true if one or more elements in the table has key Key, otherwise false.
	Member(tab string, key interface{}) bool
	// Retruns all keys in the table Tab.
	Members(tab string) ([]interface{}, bool)
	// Set key Key expire time in the table Tab.
	SetExpire(tab string, key interface{}, expire time.Duration) bool
}
```

#### TableCache

table cache is manager for k-vs

```golang
// TableCache
type TableCache interface {
	// Inserts the object or all of the objects in list.
	Insert(key, values interface{}) bool
	// Inserts the object or all of the objects with expired time in list.
	InsertExpire(key, value interface{}, expire time.Duration) bool
	// Deletes all objects with key: Key.
	DeleteObject(key interface{}) bool
	// Delete all objects in the table Tab. Remain table in cache.
	DeleteObjects()
	// Returns true if one or more elements in the table has key: Key, otherwise false.
	Member(key interface{}) bool
	// Retruns all keys in the table Tab.
	Members() ([]interface{}, bool)
	// Look up values with key: Key.
	Lookup(key interface{}) ([]interface{}, bool)
	// Look up all values in the Tab.
	LookupAll() (map[interface{}][]interface{}, bool)
	// Set Key Expire time
	SetExpire(key interface{}, expire time.Duration) bool
}
```

#### Sample: NewTableCache with options

[Examples](examples/main.go)
