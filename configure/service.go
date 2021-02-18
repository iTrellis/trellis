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

	"github.com/google/uuid"
	"github.com/iTrellis/config"
	"github.com/iTrellis/node"
	"github.com/iTrellis/trellis/service"
)

type Service struct {
	Alias string `json:"alias" yaml:"alias"`

	service.Service `json:",inline" yaml:",inline"`

	Protocol service.Protocol `json:"protocol" yaml:"protocol"`

	Options config.Options `json:"options" yaml:"options"`

	Registry *ServiceRegistry `json:"registry" yaml:"registry"`
}

type ServiceRegistry struct {
	Name string `json:"name" yaml:"name"`

	TTL        string        `json:"ttl" yaml:"ttl"`
	Heartbeat  time.Duration `json:"heartbeat" yaml:"heartbeat"`
	RetryTimes uint32        `json:"retry_times" yaml:"retry_times"`

	Address string `json:"address" yaml:"address"`
	Weight  uint32 `json:"weight" yaml:"weight"`

	Options config.Options `json:"options" yaml:"options"`
}

func (p *Service) ToNode() *node.Node {
	n := &node.Node{}

	if p.Registry == nil {
		n.ID = p.ID(uuid.New().String())
		n.Weight = 1
	} else {
		n.ID = p.ID(p.Registry.Address)
		n.Value = p.Registry.Address
		n.Weight = p.Registry.Weight
	}
	n.Metadata = p.Options

	if n.Metadata == nil {
		n.Metadata = config.Options{}
	}

	n.Metadata["protocol"] = p.Protocol

	return n
}
