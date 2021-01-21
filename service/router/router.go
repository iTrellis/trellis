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
	"context"

	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/message"
	"github.com/iTrellis/trellis/service/registry"
)

// Router router
type Router interface {
	service.LifeCycle

	Register

	Caller
}

type Caller interface {
	Call(context.Context, message.Message) (interface{}, error)
}

// Register router register
type Register interface {
	RegisterRegistry(string, registry.Registry) error
	DeregisterRegistry(string) error

	RegisterService(string, *registry.Service, ...registry.RegisterOption) error
	DeregisterService(string, *registry.Service) error

	Watch(...registry.WatchOption) (registry.Watcher, error)
}
