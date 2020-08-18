// GNU GPL v3 License
// Copyright (c) 2016 github.com:go-trellis

package cache

// ValueMode define value mode
type ValueMode int

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

// DataValues define k-vs struct
type DataValues struct {
	Key    string
	Values []interface{}
	Exists map[interface{}]bool
	Length int
}
