// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package txorm

import (
	"reflect"

	"xorm.io/xorm"
)

// TXorm trellis xorm
type TXorm struct {
	isTransaction bool
	txSession     *xorm.Session

	engines   map[string]*xorm.Engine
	defEngine *xorm.Engine
}

// New get trellis xorm committer
func New() Committer {
	return &TXorm{}
}

// SetEngines set xorm engines
func (p *TXorm) SetEngines(engines map[string]*xorm.Engine) {
	if defEngine, exist := engines[DefaultDatabase]; exist {
		p.engines = engines
		p.defEngine = defEngine
	} else {
		panic(ErrNotFoundDefaultDatabase)
	}
}

// Session get session
func (p *TXorm) Session() *xorm.Session {
	return p.txSession
}

// GetEngine get engine by name
func (p *TXorm) getEngine(name string) (*xorm.Engine, error) {
	if engine, _exist := p.engines[name]; _exist {
		return engine, nil
	}
	return nil, ErrNotFoundXormEngine
}

func (p *TXorm) checkRepos(txFunc interface{}, originRepos ...interface{}) error {
	if reposLen := len(originRepos); reposLen < 1 {
		return ErrAtLeastOneRepo
	}

	if txFunc == nil {
		return ErrNotFoundTransationFunction
	}
	return nil
}

func getRepo(v interface{}) *TXorm {
	_deepRepo := DeepFields(v, reflect.TypeOf(new(TXorm)), []reflect.Value{})
	if deepRepo, ok := _deepRepo.(*TXorm); ok {
		return deepRepo
	}
	return nil
}

func createNewTXorm(origin interface{}) (*TXorm, interface{}, error) {
	if repo, err := Derive(origin); err != nil {
		return nil, nil, err
	} else if repo != nil {
		return getRepo(repo), repo, nil
	}

	newRepoV := reflect.New(reflect.ValueOf(
		reflect.Indirect(reflect.ValueOf(origin)).Interface()).Type())
	if !newRepoV.IsValid() {
		return nil, nil, ErrFailToCreateRepo
	}

	newRepoI := newRepoV.Interface()

	if err := Inherit(newRepoI, origin); err != nil {
		return nil, nil, err
	}

	newTxorm := getRepo(newRepoI)

	if newTxorm == nil {
		return nil, nil, ErrFailToConvetTXToNonTX
	}
	return newTxorm, newRepoI, nil
}
