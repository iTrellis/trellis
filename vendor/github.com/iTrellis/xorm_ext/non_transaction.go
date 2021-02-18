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

// NonTX do non transaction function by default database
func (p *TXorm) NonTX(fn interface{}, repos ...interface{}) error {
	return p.NonTXWithName(fn, DefaultDatabase, repos...)
}

// NonTXWithName do non transaction function with name of database
func (p *TXorm) NonTXWithName(fn interface{}, name string, repos ...interface{}) error {
	if err := p.checkRepos(fn, repos); err != nil {
		return err
	}

	_newRepos := []interface{}{}
	_newTXormRepos := []*TXorm{}

	for _, origin := range repos {

		_repo := getRepo(origin)
		if _repo == nil {
			return ErrStructCombineWithRepo
		}

		_newTXorm, _newRepoI, err := createNewTXorm(origin)
		if err != nil {
			return err
		}

		_newTXorm.engines = _repo.engines
		_newTXorm.defEngine = _repo.defEngine

		if err := _newTXorm.beginNonTransaction(name); err != nil {
			return err
		}

		_newRepos = append(_newRepos, _newRepoI)
		_newTXormRepos = append(_newTXormRepos, _newTXorm)
	}

	return _newTXormRepos[0].commitNonTransaction(fn, name, _newRepos...)
}

func (p *TXorm) beginNonTransaction(name string) error {
	if p.isTransaction {
		return ErrFailToConvetTXToNonTX
	}

	_engine, err := p.getEngine(name)
	if err != nil {
		return err
	}

	p.txSession = _engine.NewSession()

	return nil
}

func (p *TXorm) commitNonTransaction(txFunc interface{}, name string, repos ...interface{}) error {
	if p.isTransaction {
		return ErrNonTransactionCantCommit
	}

	_funcs := GetLogicFuncs(txFunc)

	var (
		_values []interface{}
		errcode error
	)

	if _funcs.BeforeLogic != nil {
		if _, errcode = CallFunc(_funcs.BeforeLogic, _funcs, repos); errcode != nil {
			return errcode
		}
	}

	if _funcs.Logic != nil {
		if _values, errcode = CallFunc(_funcs.Logic, _funcs, repos); errcode != nil {
			return errcode
		}
	}

	if _funcs.AfterLogic != nil {
		if _values, errcode = CallFunc(_funcs.AfterLogic, _funcs, repos); errcode != nil {
			return errcode
		}
	}

	if _funcs.AfterCommit != nil {
		if _, errcode = CallFunc(_funcs.AfterCommit, _funcs, _values); errcode != nil {
			return errcode
		}
	}

	return nil
}
