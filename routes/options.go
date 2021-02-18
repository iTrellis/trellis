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

package routes

import (
	"github.com/iTrellis/common/logger"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/router"
)

type Option func(*Options)

type Options struct {
	logger  logger.Logger
	router  router.Router
	manager component.Manager
}

func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.logger = l
	}
}

func WithRouter(r router.Router) Option {
	return func(o *Options) {
		o.router = r
	}
}

func CompManager(m component.Manager) Option {
	return func(o *Options) {
		o.manager = m
	}
}
