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

import (
	"time"

	"github.com/iTrellis/trellis/service"

	"github.com/iTrellis/node"
)

// Watcher is an interface that returns updates
// about services within the registry.
type Watcher interface {
	// Next is a blocking call
	Next() (*Result, error)
	Stop()
}

// Service service
type Service struct {
	service.Service `json:",inline" yaml:",inline"`

	Nodes []*node.Node `json:"nodes" yaml:"nodes"`
}

// Result is registry result
type Result struct {
	// Id is registry id
	ID string
	// Type defines type of event
	Type service.EventType
	// Timestamp is event timestamp
	Timestamp time.Time
	// Service is registry service
	Service *Service
}
