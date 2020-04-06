package server

import (
	"fmt"

	"github.com/go-trellis/config"
	"github.com/go-trellis/trellis/router"
)

var (
	serverFuncs map[string]NewServerFunc = make(map[string]NewServerFunc)
	serverNames []string
)

// OptionFunc 处理参数函数
type OptionFunc func(*Options)

// Options 参数对象
type Options struct {
	Config config.Config
}

// Config 注入配置
func Config(c config.Config) OptionFunc {
	return func(p *Options) {
		p.Config = c
	}
}

// Server 服务对象
type Server interface {
	LifeCycle
	router.Router
}

// LifeCycle server的生命周期
type LifeCycle interface {
	Start() error
	Stop() error
}

// NewServerFunc 服务对象生成函数申明
type NewServerFunc func(opts ...OptionFunc) (Server, error)

// RegistNewServerFunc 注册服务对象生成函数
func RegistNewServerFunc(name string, fn NewServerFunc) {

	if len(name) == 0 {
		panic("server name is empty")
	}

	if fn == nil {
		panic("server function is nil")
	}

	_, exist := serverFuncs[name]

	if exist {
		panic(fmt.Sprintf("server is already registered: %s", name))
	}

	serverFuncs[name] = fn
	serverNames = append(serverNames, name)
}

// New 生成函数对象
func New(name string, opts ...OptionFunc) (Server, error) {
	fn, exist := serverFuncs[name]
	if !exist {
		return nil, fmt.Errorf("server '%s' not exist", name)
	}
	return fn(opts...)
}
