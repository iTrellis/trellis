package runner

import (
	"github.com/go-trellis/errors"
	"github.com/go-trellis/trellis/errcode"
	"github.com/go-trellis/trellis/message/proto"
)

var runner *Runner

func Run(opts ...OptionFunc) (err error) {
	runner, err = New(opts...)
	if err != nil {
		return
	}

	return runner.Run()
}

func New(opts ...OptionFunc) (*Runner, error) {

	tOpts := &Options{}

	for _, o := range opts {
		o(tOpts)
	}

	t := &Runner{
		conf:    tOpts.conf,
		workers: make(map[string]*Worker),
	}

	if err := t.registServers(); err != nil {
		return nil, err
	}

	return t, nil
}

func GetWorker(base *proto.BaseService) (*Worker, errors.ErrorCode) {
	runner.locker.RLock()
	defer runner.locker.RUnlock()
	worker, ok := runner.workers[runner.genWorkerPath(base.GetName(), base.GetVersion())]
	if !ok {
		return nil, errcode.ErrNotFoundServerWorker.New(
			errors.Params{"name": base.GetName(), "version": base.GetVersion()})
	}
	return worker, nil
}

func RevokeWorker(name, version string) {
	runner.locker.Lock()
	defer runner.locker.Unlock()
	runner.removeWorker(runner.genWorkerPath(name, version))
}
