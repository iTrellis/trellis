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
	"github.com/iTrellis/node"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/registry"
)

// Router router
type Router interface {
	service.LifeCycle

	// // Caller handle serving messages
	// type Caller interface {
	// Remote for handler TODO
	GetServiceNodes(...ReadOption) ([]*node.Node, error) // []*node.Node || component
	// }

	// // Registry router registry
	// type Registry interface {
	RegisterRegistry(string, registry.Registry) error
	DeregisterRegistry(string) error

	RegisterService(string, *service.Service, ...registry.RegisterOption) error
	DeregisterService(string, *service.Service, ...registry.DeregisterOption) error
	WatchService(string, ...registry.WatchOption) error
}
