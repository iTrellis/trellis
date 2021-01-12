package component

import (
	"context"

	"github.com/go-trellis/trellis/service"
	"github.com/go-trellis/trellis/service/message"
	"github.com/go-trellis/trellis/service/router"

	"github.com/go-trellis/common/logger"
)

// NewComponentFunc 服务对象生成函数申明
type NewComponentFunc func(alias string, opts ...Option) (Component, error)

// Handler handle the message function
type Handler func(context.Context, message.Message) (interface{}, error)

// Middleware middlerwares for next handler
type Middleware func(Handler) Handler

// Component Component
type Component interface {
	Alias() string

	service.LifeCycle

	Route(topic string) Handler
}

type ComponentDescribe struct {
	Name         string
	RegisterFunc string
	Component    Component
}

// Option 处理参数函数
type Option func(*Options)

// Options 参数对象
type Options struct {
	Logger  logger.Logger
	Context context.Context
	Router  router.Router
}

// Context 注入配置
func Context(c context.Context) Option {
	return func(p *Options) {
		p.Context = c
	}
}

// Logger 日志记录
func Logger(l logger.Logger) Option {
	return func(p *Options) {
		p.Logger = l
	}
}

// Router 路由表，可用与服务间调用
func Router(r router.Router) Option {
	return func(p *Options) {
		p.Router = r
	}
}
