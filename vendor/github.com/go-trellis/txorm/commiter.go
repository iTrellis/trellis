// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package txorm

// Committer gorm committer
type committer struct {
	Name string
}

// NewCommitter get trellis gorm committer
func NewCommitter() Committer {
	return &committer{Name: "go-trellis::txorm::committer"}
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
