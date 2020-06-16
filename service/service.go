package service

import (
	"fmt"

	"github.com/go-trellis/config"
	"github.com/go-trellis/trellis/router"
)

var (
	serviceFuncs map[string]NewServiceFunc = make(map[string]NewServiceFunc)
	serverNames  []string
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

// Service 服务对象
type Service interface {
	LifeCycle
	router.Router
}

// LifeCycle server的生命周期
type LifeCycle interface {
	Start() error
	Stop() error
}

// NewServiceFunc 服务对象生成函数申明
type NewServiceFunc func(opts ...OptionFunc) (Service, error)

// RegistNewServiceFunc 注册服务对象生成函数
func RegistNewServiceFunc(name string, fn NewServiceFunc) {

	if len(name) == 0 {
		panic("server name is empty")
	}

	if fn == nil {
		panic("server function is nil")
	}

	_, exist := serviceFuncs[name]

	if exist {
		panic(fmt.Sprintf("server is already registered: %s", name))
	}

	serviceFuncs[name] = fn
	serverNames = append(serverNames, name)
}

// New 生成函数对象
func New(name string, opts ...OptionFunc) (Service, error) {
	fn, exist := serviceFuncs[name]
	if !exist {
		return nil, fmt.Errorf("server '%s' not exist", name)
	}
	return fn(opts...)
}
