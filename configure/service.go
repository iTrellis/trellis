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

package configure

import (
	"time"

	"github.com/iTrellis/config"
	"github.com/iTrellis/trellis/service"
)

// Service service
type Service struct {
	service.Service `json:",inline" yaml:",inline"`

	Options config.Options `json:"options" yaml:"options"`

	Registry *ServiceRegistry `json:"registry" yaml:"registry"`
}

// ServiceRegistry service's registry infor
type ServiceRegistry struct {
	// registry name
	Name string `json:"name" yaml:"name"`
	// node weight
	Weight uint32 `json:"weight" yaml:"weight"`
	// protocol between two servers
	Protocol service.Protocol `json:"protocol" yaml:"protocol"`

	TTL       time.Duration `json:"ttl" yaml:"ttl"`
	Heartbeat time.Duration `json:"heartbeat" yaml:"heartbeat"`
}

// func (p *Service) ToNode(*Registry) *node.Node {
// 	n := &node.Node{
// 		Metadata: p.Options,
// 	}

// 	if n.Metadata == nil {
// 		n.Metadata = config.Options{}
// 	}

// 	if p.Registry == nil {
// 		n.ID = p.ID("127.0.0.1")
// 		n.Weight = 1
// 		n.Metadata["protocol"] = p.Registry.Protocol
// 	} else {
// 		n.ID = p.ID(p.Registry.Address)
// 		n.Value = p.Registry.Address
// 		n.Weight = p.Registry.Weight
// 	}

// 	return n
// }
