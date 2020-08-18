// GNU GPL v3 License
// Copyright (c) 2016 github.com:go-trellis

package cache

import "errors"

// errors
var (
	ErrTableExists           = errors.New("table already exists")
	ErrNewTableCache         = errors.New("failed new table")
	ErrOrderSetMustBeBool    = errors.New("order set must be bool")
	ErrUnknownTableOption    = errors.New("unknown table option")
	ErrUnknownTableValueMode = errors.New("unknown table value mode")
)
