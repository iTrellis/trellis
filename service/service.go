/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package service

import (
	"fmt"

	"github.com/go-trellis/trellis/message/proto"

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
	Config config.Config
	Logger logger.Logger
}

// Config 注入配置
func Config(c config.Config) OptionFunc {
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

	s := proto.Service{Name: name, Version: version}

	_, exist := serviceFuncs[s.String()]

	if exist {
		panic(fmt.Sprintf("server is already registered: %s", s.String()))
	}

	serviceFuncs[s.String()] = fn
	serverNames = append(serverNames, s.String())
}

// New 生成函数对象
func New(service *proto.Service, opts ...OptionFunc) (Service, error) {
	fn, exist := serviceFuncs[service.String()]
	if !exist {
		return nil, fmt.Errorf("server '%s' not exist", service.String())
	}
	return fn(opts...)
}
