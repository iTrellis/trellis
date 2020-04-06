package client

import (
	"github.com/go-trellis/errors"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/runner"
)

// type InnerClient struct{}

func InnerCall(req *message.Request) errors.ErrorCode {

	worker, err := runner.GetWorker(req.Server())

	if err != nil {
		return err
	}

	return worker.Call(req.Message)
}

// func (*InnerClient) String() string {
// 	return "inner"
// }

// func NewInnerClient() Client {
// 	return (*InnerClient)(nil)
// }
