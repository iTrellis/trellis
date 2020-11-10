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

package router

import (
	"github.com/go-trellis/trellis/clients"
	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/message/proto"

	"github.com/go-trellis/common/logger"
)

// Router 路由器
type Router interface {
	NewService(...OptionFunc) error
	StopService(*proto.Service) error
	Run() error
	Stop() error

	clients.Caller
}

// OptionFunc 配置函数定义
type OptionFunc func(*Options)

// Options 配置
type Options struct {
	Service *configure.Service

	Logger logger.Logger
}

// OptionService 配置参数
func OptionService(s *configure.Service) OptionFunc {
	return func(w *Options) {
		w.Service = s
	}
}

// OptionLogger 日志
func OptionLogger(l logger.Logger) OptionFunc {
	return func(w *Options) {
		w.Logger = l
	}
}
