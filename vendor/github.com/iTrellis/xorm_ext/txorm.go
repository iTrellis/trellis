/*
Copyright © 2019 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package xorm_ext

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
