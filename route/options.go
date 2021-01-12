package route

import (
	"github.com/go-trellis/common/logger"
	"github.com/go-trellis/trellis/service/component"
)

type Option func(*Options)

type Options struct {
	logger logger.Logger
	local  component.Manager
}

func Logger(logger logger.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

func LocalRouter(local component.Manager) Option {
	return func(o *Options) {
		o.local = local
	}
}
