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

package component

import (
	"github.com/iTrellis/common/logger"
	"github.com/iTrellis/config"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/message"
)

// NewComponentFunc 服务对象生成函数申明
type NewComponentFunc func(alias string, opts ...Option) (Component, error)

// Handler handle the message function
type Handler func(message.Message) (interface{}, error)

// Middleware middlerwares for next handler
type Middleware func(Handler) Handler

// Component Component
type Component interface {
	Alias() string

	service.LifeCycle

	Route(topic string) Handler
}

// Describe description of component
type Describe struct {
	Name         string
	RegisterFunc string
	Component    Component
}

// Option 处理参数函数
type Option func(*Options)

// Options 参数对象
type Options struct {
	Logger logger.Logger
	Config config.Config
	Caller message.Caller
}

// Config 注入配置
func Config(c config.Config) Option {
	return func(p *Options) {
		p.Config = c
	}
}

// Logger 日志记录
func Logger(l logger.Logger) Option {
	return func(p *Options) {
		p.Logger = l
	}
}

// Caller remote service
func Caller(c message.Caller) Option {
	return func(p *Options) {
		p.Caller = c
	}
}
