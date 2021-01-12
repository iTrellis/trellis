package configure

import (
	"github.com/go-trellis/node"
	"github.com/go-trellis/trellis/service"
)

type Registry struct {
	Name string               `json:"name" yaml:"name"`
	Type service.RegisterType `json:"type" yaml:"type"`

	Address []string `json:"address" yaml:"address"`
	Secure  bool     `json:"secure" yaml:"secure"`

	Watcher `json:",inline" yaml:",inline"`
}

type Watcher struct {
	Services []service.Service `json:"services" yaml:"services"`
	NodeType node.Type         `json:"type" yaml:"type"`
}
