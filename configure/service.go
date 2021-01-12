package configure

import (
	"github.com/go-trellis/common/encryption/hash"
	"github.com/go-trellis/config"
	"github.com/go-trellis/node"
	"github.com/go-trellis/trellis/service"
)

type Service struct {
	Alias string `json:"alias" yaml:"alias"`

	service.Service `json:",inline" yaml:",inline"`

	Address string `json:"address" yaml:"address"`
	Weight  uint32 `json:"weight" yaml:"weight"`

	TransportType service.TransportType `json:"transport_type" yaml:"transport_type"`

	Registry *ServiceRegistry `json:"registry" yaml:"registry"`
}

type ServiceRegistry struct {
	Name string `json:"name" yaml:"name"`
	TTL  string `json:"ttl" yaml:"ttl"`

	Options config.Options `json:"options" yaml:"options"`
}

func (s *Service) ToNode() *node.Node {
	n := &node.Node{
		ID:     hash.NewCRCIEEE().Sum(s.Service.FullPath(s.Address)),
		Value:  s.Address,
		Weight: s.Weight,
	}

	if n.Metadata == nil {
		n.Metadata = config.Options{}
	}

	n.Metadata["transport_type"] = s.TransportType

	return n
}
