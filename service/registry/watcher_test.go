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

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/iTrellis/common/testutils"
// 	"github.com/iTrellis/config"
// 	"github.com/iTrellis/node"
// 	"github.com/iTrellis/trellis/service"
// )

// func TestServiceToNodes(t *testing.T) {
// 	service := &Service{
// 		Service: service.Service{
// 			Domain:  "trellis",
// 			Name:    "test",
// 			Version: "v1",
// 		},
// 		Nodes: []*node.Node{
// 			&node.Node{
// 				ID:       "1",
// 				Weight:   1,
// 				Metadata: config.Options{"nKey1": "nValue1"},
// 				Value:    "node1"},
// 			&node.Node{
// 				ID:       "2",
// 				Weight:   2,
// 				Metadata: config.Options{"nKey2": "nValue2", "testk": "ntestv"},
// 				Value:    "node2"},
// 			&node.Node{
// 				ID:     "3",
// 				Weight: 3,
// 				Value:  "node3"},
// 		},
// 	}

// 	for _, node := range service.Nodes {
// 		switch node.ID {
// 		case "1":
// 			testutils.Equals(t, "nValue1", node.Metadata["nKey1"])
// 			testutils.Equals(t, nil, node.Metadata["testk"])
// 			testutils.Equals(t, "node1", node.Value)
// 			testutils.Equals(t, uint32(1), node.Weight)
// 		case "2":
// 			testutils.Equals(t, nil, node.Metadata["nKey1"])
// 			testutils.Equals(t, "nValue2", node.Metadata["nKey2"])
// 			testutils.Equals(t, "ntestv", node.Metadata["testk"])
// 			testutils.Equals(t, "node2", node.Value)
// 			testutils.Equals(t, uint32(2), node.Weight)
// 		case "3":
// 			testutils.Equals(t, nil, node.Metadata["nKey1"])
// 			testutils.Equals(t, nil, node.Metadata["testk"])
// 			testutils.Equals(t, "node3", node.Value)
// 			testutils.Equals(t, uint32(3), node.Weight)
// 		}
// 	}
// }

// // ToNodeManager service to node manager for nodes
// func ToNodeManager(s *Service, t node.Type) (node.Manager, error) {
// 	nm, err := node.New(t, s.FullPath())
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, n := range s.Nodes {
// 		item := *n
// 		nm.Add(&item)
// 	}
// 	return nm, nil
// }

// func TestServiceToNodeManager(t *testing.T) {
// 	service := &Service{
// 		Service: service.Service{
// 			Domain:  "trellis",
// 			Name:    "test",
// 			Version: "v1",
// 		},
// 		Nodes: []*node.Node{
// 			&node.Node{
// 				ID:       "1",
// 				Weight:   10,
// 				Metadata: config.Options{"nKey1": "nValue1"},
// 				Value:    "node1"},
// 			&node.Node{
// 				ID:       "2",
// 				Weight:   20,
// 				Metadata: config.Options{"nKey2": "nValue2", "testk": "ntestv"},
// 				Value:    "node2"},
// 			&node.Node{
// 				ID:     "3",
// 				Weight: 30,
// 				Value:  "node3"},
// 		},
// 	}

// 	nm, err := ToNodeManager(service, node.NodeTypeDirect)
// 	testutils.Ok(t, err)
// 	n := 0
// 	for n < 3 {
// 		node, ok := nm.NodeFor()
// 		testutils.Equals(t, true, ok)
// 		testutils.Equals(t, "3", node.ID)
// 		testutils.Equals(t, "node3", node.Value)
// 		n++
// 	}

// 	nm, err = ToNodeManager(service, node.NodeTypeRoundRobin)
// 	testutils.Ok(t, err)
// 	n = 0
// 	for n < 100 {
// 		node, ok := nm.NodeFor()
// 		testutils.Equals(t, true, ok)
// 		switch n % 3 {
// 		case 0:
// 			testutils.Equals(t, "1", node.ID)
// 			testutils.Equals(t, "node1", node.Value)
// 		case 1:
// 			testutils.Equals(t, "2", node.ID)
// 			testutils.Equals(t, "node2", node.Value)
// 		case 2:
// 			testutils.Equals(t, "3", node.ID)
// 			testutils.Equals(t, "node3", node.Value)
// 		}
// 		n++
// 	}

// 	nm, err = ToNodeManager(service, node.NodeTypeRandom)
// 	testutils.Ok(t, err)
// 	n = 0
// 	counts := make(map[string]int, 3)
// 	for n < 6000 {
// 		node, ok := nm.NodeFor()
// 		testutils.Equals(t, true, ok)
// 		counts[node.ID]++
// 		n++
// 	}
// 	for k, v := range counts {
// 		switch k {
// 		case "1":
// 			if 900 < v || v > 1100 {
// 				testutils.NotOk(t, fmt.Errorf("node1 should be 1000"))
// 			}
// 		case "2":
// 			if 1900 < v || v > 2100 {
// 				testutils.NotOk(t, fmt.Errorf("node2 should be 2000"))
// 			}
// 		case "3":
// 			if 2900 < v || v > 3100 {
// 				testutils.NotOk(t, fmt.Errorf("node3 should be 3000"))
// 			}
// 		}
// 	}

// }
