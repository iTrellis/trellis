// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package txorm

// TX do transaction function by default database
func (p *TXorm) TX(fn interface{}, repos ...interface{}) error {
	return p.TXWithName(fn, DefaultDatabase, repos...)
}

// TXWithName do transaction function with name of database
func (p *TXorm) TXWithName(fn interface{}, name string, repos ...interface{}) error {
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

		_newTxorm, newRepoI, err := createNewTXorm(origin)
		if err != nil {
			return err
		}

		_newTxorm.engines = repo.engines
		_newTxorm.defEngine = repo.defEngine
		_newRepos = append(_newRepos, newRepoI)
		_newTXormRepos = append(_newTXormRepos, _newTxorm)
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

func (p *TXorm) beginTransaction(name string) error {
	if !p.isTransaction {
		p.isTransaction = true
		_engine, err := p.getEngine(name)
		if err != nil {
			return err
		}
		p.txSession = _engine.NewSession()
		if p.txSession == nil {
			return ErrTransactionSessionIsNil
		}
		return nil
	}
	return ErrTransactionIsAlreadyBegin
}

func (p *TXorm) commitTransaction(txFunc interface{}, repos ...interface{}) error {
	if !p.isTransaction {
		return ErrNonTransactionCantCommit
	}

	if p.txSession == nil {
		return ErrTransactionSessionIsNil
	}
	defer p.txSession.Close()

	if txFunc == nil {
		return ErrNotFoundTransationFunction
	}

	isNeedRollBack := true

	if err := p.txSession.Begin(); err != nil {
		return err
	}

	defer func() {
		if isNeedRollBack {
			_ = p.txSession.Rollback()
		}
	}()

	_funcs := GetLogicFuncs(txFunc)

	var (
		_values []interface{}
		ecode   error
	)

	if _funcs.BeforeLogic != nil {
		if _, ecode = CallFunc(_funcs.BeforeLogic, _funcs, repos); ecode != nil {
			return ecode
		}
	}

	if _funcs.Logic != nil {
		if _values, ecode = CallFunc(_funcs.Logic, _funcs, repos); ecode != nil {
			return ecode
		}
	}

	if _funcs.AfterLogic != nil {
		if _values, ecode = CallFunc(_funcs.AfterLogic, _funcs, repos); ecode != nil {
			return ecode
		}
	}

	isNeedRollBack = false
	if err := p.txSession.Commit(); err != nil {
		return err
	}

	if _funcs.AfterCommit != nil {
		if _, ecode = CallFunc(_funcs.AfterCommit, _funcs, _values); ecode != nil {
			return ecode
		}
	}

	return nil
}
