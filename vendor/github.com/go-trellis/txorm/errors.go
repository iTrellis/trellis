// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package txorm

import "errors"

// define connector errors
var (
	ErrNotFoundDefaultDatabase    = errors.New("not found default database")
	ErrAtLeastOneRepo             = errors.New("input one repo at least")
	ErrNotFoundTransationFunction = errors.New("not found transation function")
	ErrStructCombineWithRepo      = errors.New("your repository struct should combine repo")
	ErrFailToCreateRepo           = errors.New("fail to create an new repo")
	ErrFailToConvetTXToNonTX      = errors.New("could not convert TX to NON-TX")
	ErrTransactionIsAlreadyBegin  = errors.New("transaction is already begin")
	ErrNonTransactionCantCommit   = errors.New("non-transaction can't commit")
	ErrTransactionSessionIsNil    = errors.New("transaction session is nil")
	ErrNotFoundXormEngine         = errors.New("not found xorm engine")
)
