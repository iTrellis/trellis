package errcode

import "github.com/go-trellis/common/errors"

const namespace = "trellis"

// 错误集合
var (
	ErrBadRequest      = errors.TN(namespace, 10, "bad request: {{.err}}")
	ErrAPINotFound     = errors.TN(namespace, 11, "api not found")
	ErrGetService      = errors.TN(namespace, 12, "{{.err}}")
	ErrGetServiceTopic = errors.TN(namespace, 13, "failed get service's topic")
	ErrCallService     = errors.TN(namespace, 14, "failed call service: {{.err}}")
)
