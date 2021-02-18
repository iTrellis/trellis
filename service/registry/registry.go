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

package registry

// NewRegistryFunc new registry function
type NewRegistryFunc func(...Option) (Registry, error)

// Registry The registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Init(...Option) error
	Options() Options

	ProcessService

	Stop() error

	ID() string
	String() string
}

// ProcessService Process register service
type ProcessService interface {
	Register(*Service, ...RegisterOption) error
	Deregister(*Service, ...DeregisterOption) error
	Watch(...WatchOption) (Watcher, error)
}
