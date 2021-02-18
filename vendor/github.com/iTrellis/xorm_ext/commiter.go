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

// Committer gorm committer
type committer struct {
	Name string
}

// NewCommitter get trellis gorm committer
func NewCommitter() Committer {
	return &committer{Name: "iTrellis::txorm::committer"}
}

// NonTX do non transaction function by default database
func (p *committer) NonTX(fn interface{}, repos ...interface{}) error {
	return p.NonTXWithName(fn, DefaultDatabase, repos...)
}

// NonTXWithName do non transaction function with name of database
func (p *committer) NonTXWithName(fn interface{}, name string, repos ...interface{}) error {
	if err := p.checkRepos(fn, repos); err != nil {
		return err
	}

	_newRepos := []interface{}{}
	_newTXormRepos := []*TXorm{}

	for _, origin := range repos {
		repo := getRepo(origin)
		if repo == nil {
			return ErrStructCombineWithRepo
		}

		_newTxorm, _newRepoI, err := createNewTXorm(origin)
		if err != nil {
			return err
		}

		_newRepos = append(_newRepos, _newRepoI)

		_newTxorm.engines = repo.engines
		_newTxorm.defEngine = repo.defEngine

		if err := _newTxorm.beginNonTransaction(name); err != nil {
			return err
		}

		_newTXormRepos = append(_newTXormRepos, _newTxorm)
	}

	return _newTXormRepos[0].commitNonTransaction(fn, name, _newRepos...)
}

// TX do transaction function by default database
func (p *committer) TX(fn interface{}, repos ...interface{}) error {
	return p.TXWithName(fn, DefaultDatabase, repos...)
}

// TXWithName do transaction function with name of database
func (p *committer) TXWithName(fn interface{}, name string, repos ...interface{}) error {
	if err := p.checkRepos(fn, repos); err != nil {
		return err
	}

	_newRepos := []interface{}{}
	_newTXormRepos := []*TXorm{}

	for _, origin := range repos {

		repo := getRepo(origin)
		if repo == nil {
			return ErrStructCombineWithRepo
		}

		_newTXorm, _newRepoI, err := createNewTXorm(origin)
		if err != nil {
			return err
		}

		_newTXorm.engines = repo.engines
		_newTXorm.defEngine = repo.defEngine
		_newRepos = append(_newRepos, _newRepoI)
		_newTXormRepos = append(_newTXormRepos, _newTXorm)
	}

	if err := _newTXormRepos[0].beginTransaction(name); err != nil {
		return err
	}

	for i := range _newTXormRepos {
		_newTXormRepos[i].txSession = _newTXormRepos[0].txSession
		_newTXormRepos[i].isTransaction = _newTXormRepos[0].isTransaction
	}

	return _newTXormRepos[0].commitTransaction(fn, _newRepos...)
}

func (p *committer) checkRepos(txFunc interface{}, originRepos ...interface{}) error {
	if reposLen := len(originRepos); reposLen < 1 {
		return ErrAtLeastOneRepo
	}

	if txFunc == nil {
		return ErrNotFoundTransationFunction
	}
	return nil
}
