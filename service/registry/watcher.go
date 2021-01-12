package registry

import (
	"time"

	"github.com/go-trellis/trellis/service"

	"github.com/go-trellis/node"
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
