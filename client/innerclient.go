package client

import (
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/runner"
)

// // type InnerClient struct{}

func InnerCall(req *message.Request) error {

	worker, err := runner.GetWorker(req.Service())

	if err != nil {
		return err
	}

	return worker.Call(req.Message)
}