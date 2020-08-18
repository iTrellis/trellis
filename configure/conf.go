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
	"encoding/json"
	"time"

	"github.com/go-trellis/common/logger"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message/proto"

	"github.com/go-trellis/config"
	"github.com/go-trellis/node"

	// cobra in main
	_ "github.com/spf13/cobra"
)

type Config struct {
	Project *Project `yaml:"project"`
}

// Project project config info
type Project struct {
	Logger LoggerConfig `yaml:"logger"`

	Services map[string]*Service `yaml:"services"`

	Registries map[string]*Registry `yam:"registries"`
}
type LoggerConfig struct {
	Level      logger.Level `yaml:"level"`
	ChanBuffer int          `yaml:"chan_buffer"`
	Separator  string       `yaml:"separator"`
}

// Service service info
type Service struct {
	proto.Service `yaml:",inline"`

	// Protocol string         `yaml:"protocol"`
	Options config.Options `yaml:"options"`
}

// Registry run configure
type Registry struct {
	Type    RegistryType   `yaml:"type"`
	Options config.Options `yaml:"options"`

	Services []*RegistService `yaml:"services"`
	Watchers []*Watcher       `yaml:"watchers"`
}

// RegistryType 注册机类型
type RegistryType string

// 注册类型实例
const (
	RegistryTypeCache = "cache"
	RegistryTypeETCD  = "etcd"
)

// RegistService service which should regist into registry
type RegistService struct {
	Name    string `yaml:"name" jpath:"name"`
	Version string `yaml:"version" jpath:"version"`
	// Field    string            `yaml:"field"`
	Protocol string            `yaml:"protocal" jpath:"protocol"`
	Domain   string            `yaml:"domain" jpath:"domain"`
	Weight   uint32            `yaml:"option" jpath:"weight"`
	Metadata map[string]string `yaml:"metadata" jpath:"metadata"`
}

func (p *RegistService) String() string {
	bs, _ := json.Marshal(p)
	return string(bs)
}

func ToRegistService(str string) (*RegistService, error) {
	s := &RegistService{}
	err := json.Unmarshal([]byte(str), s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

type RegistServices []*RegistService

func (p *RegistServices) UnmarshalJSON(bs []byte) error {
	if p == nil {
		p = &RegistServices{}
	}
	return json.Unmarshal(bs, p)
}

func (p RegistServices) ToNodeManager(nt node.Type) node.Manager {
	if len(p) == 0 {
		return nil
	}

	var nm node.Manager

	for _, s := range p {
		if nm == nil {
			nm = node.New(nt, internal.WorkerPath(internal.SchemaTrellis, s.Name, s.Version))
		}
		nm.Add(&node.Node{
			ID:     internal.WorkerDomainPath(internal.SchemaTrellis, s.Name, s.Version, s.Domain),
			Weight: s.Weight,
			Value:  s.String(),
			Metadata: map[string]interface{}{
				"protocol": s.Protocol,
				"domain":   s.Domain},
		})
	}
	return nm
}

type Watcher struct {
	proto.Service `yaml:",inline"`

	TTL time.Duration `yaml:"ttl"`

	LoadBalancing node.Type `yaml:"load_balancing" jpath:"load_balancing"`
}
