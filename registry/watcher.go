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
	"strings"

	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message/proto"

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
		Value:  p.Service.Domain,
		Metadata: map[string]interface{}{
			"protocol": proto.Protocol(proto.Protocol_value[strings.ToUpper(p.Service.Protocol)]),
		}}
}
