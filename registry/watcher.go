package registry

import (
	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/internal"

	"github.com/go-trellis/node"
)

// Watcher the watcher provides an interface for
// watching the services' configures.
type Watcher interface {
	// Watch
	// Next is a blocking call
	Next(ch chan *Result)
	Stop()

	Fullpath() string
}

// Actions
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// Result is returned by a call to Next on
// the watcher. Actions can be create, update, delete
type Result struct {
	NodeType node.Type
	Err      error
	Action   string
	Service  *configure.RegistService
}

// ToNode 封装 node
func (p *Result) ToNode() *node.Node {
	if p == nil || p.Service == nil {
		return nil
	}

	return &node.Node{
		ID:     internal.WorkerTrellisDomainPath(p.Service.Name, p.Service.Version, p.Service.Domain),
		Weight: p.Service.Weight,
		Value:  p.Service.String(),
		Metadata: map[string]interface{}{
			"protocol": p.Service.Protocol,
			"domain":   p.Service.Domain},
	}
}
