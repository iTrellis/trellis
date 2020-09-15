// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package txorm

// Committer Defination
type Committer interface {
	TX(fn interface{}, repos ...interface{}) error
	TXWithName(fn interface{}, name string, repos ...interface{}) error
	NonTX(fn interface{}, repos ...interface{}) error
	NonTXWithName(fn interface{}, name string, repos ...interface{}) error
}

// TXFunc Transation function
type TXFunc func(repos []interface{}) error

// Inheritor inherit function
type Inheritor interface {
	Inherit(repo interface{}) error
}

// Inherit inherit a new repository from origin repository
func Inherit(new, origin interface{}) error {
	if inheritor, ok := new.(Inheritor); ok {
		return inheritor.Inherit(origin)
	}
	return nil
}

// Deriver derive function
type Deriver interface {
	Derive() (repo interface{}, err error)
}

// Derive derive from developer function
func Derive(origin interface{}) (interface{}, error) {
	if deriver, ok := origin.(Deriver); ok {
		return deriver.Derive()
	}
	return nil, nil
}
