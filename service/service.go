package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message"

	"github.com/go-trellis/common/logger"
	"github.com/go-trellis/config"
)

var (
	serviceFuncs map[string]NewServiceFunc = make(map[string]NewServiceFunc)
	serverNames  []string
)

// OptionFunc 处理参数函数
type OptionFunc func(*Options)

// Options 参数对象
type Options struct {
	Config config.Options
	Logger logger.Logger
}

// Config 注入配置
func Config(c config.Options) OptionFunc {
	return func(p *Options) {
		p.Config = c
	}
}

// Logger 日志记录
func Logger(l logger.Logger) OptionFunc {
	return func(p *Options) {
		p.Logger = l
	}
}

// Service 服务对象
type Service interface {
	LifeCycle
	Handlers
}

// HandlerFunc 函数执行
type HandlerFunc func(*message.Message) (interface{}, error)

// Handlers 函数路由器
type Handlers interface {
	Route(topic string) HandlerFunc
}

// LifeCycle server的生命周期
type LifeCycle interface {
	Start() error
	Stop() error
}

// NewServiceFunc 服务对象生成函数申明
type NewServiceFunc func(opts ...OptionFunc) (Service, error)

// RegistNewServiceFunc 注册服务对象生成函数
func RegistNewServiceFunc(name, version string, fn NewServiceFunc) {

	if len(name) == 0 {
		panic("server name is empty")
	}

	if fn == nil {
		panic("server function is nil")
	}

	serviceKey := internal.WorkerTrellisPath(name, version)

	_, exist := serviceFuncs[serviceKey]

	if exist {
		panic(fmt.Sprintf("server is already registered: %s", serviceKey))
	}

	serviceFuncs[serviceKey] = fn
	serverNames = append(serverNames, serviceKey)
}

// New 生成函数对象
func New(name, version string, opts ...OptionFunc) (Service, error) {
	serviceKey := internal.WorkerTrellisPath(name, version)
	fn, exist := serviceFuncs[serviceKey]
	if !exist {
		return nil, fmt.Errorf("server '%s' not exist", serviceKey)
	}
	return fn(opts...)
}

// BlockStop 阻断式停止
func BlockStop() error {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	case <-ch:
	}
	return runner.Stop()
}

// Stop 停止服务
func Stop() error {
	defer time.Sleep(time.Second)
	return runner.Stop()
}

// Run 运行
func Run(cfg *configure.Project, l logger.Logger) (err error) {
	runner, err = NewRunner(cfg, l)
	if err != nil {
		return
	}

	err = runner.Run()
	if err != nil {
		return err
	}

	return nil
}
