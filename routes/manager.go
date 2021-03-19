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

package routes

import (
	"context"
	"fmt"

	"github.com/iTrellis/common/logger"
	"github.com/iTrellis/node"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/client/grpc"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/message"
	"github.com/iTrellis/trellis/service/router"
)

type Manager interface {
	Init(...Option)

	service.LifeCycle

	Router() router.Router

	CompManager() component.Manager

	message.Caller
}

// NewManager routes manager
func NewManager(opts ...Option) Manager {

	r := &manager{}

	r.Init(opts...)

	return r
}

type manager struct {
	router  router.Router
	manager component.Manager
	logger  logger.Logger
}

func (p *manager) Init(opts ...Option) {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	if options.router != nil {
		p.router = options.router
	}

	if options.manager != nil {
		p.manager = options.manager
	}

	if p.router == nil {
		p.router = NewRoutes(options.logger)
	}

	if p.manager == nil {
		p.manager = NewCompManager()
	}

	p.logger = options.logger
}

func (p *manager) CallComponent(ctx context.Context, msg message.Message) (interface{}, error) {

	cpt, err := p.manager.GetComponent(msg.Service())
	if err != nil {
		return nil, err
	} else if cpt == nil {
		return nil, fmt.Errorf("unknown component")
	}

	return cpt.Route(msg)
}

func (p *manager) CallServer(ctx context.Context, msg message.Message) (interface{}, error) {

	nodes, err := p.router.GetServiceNodes(router.ReadService(msg.Service()))
	if err != nil {
		return nil, err
	}

	nm, err := node.NewWithNodes(node.NodeTypeConsistent, msg.Service().TrellisPath(), nodes)
	if err != nil {
		return nil, err
	}

	var keys []string

	node, ok := nm.NodeFor(keys...)
	if !ok {
		return nil, fmt.Errorf("not found service node")
	}

	var rep interface{}
	switch node.Metadata["protocol"] {
	case service.Protocol_GRPC:
		fallthrough
	default:
		c := grpc.NewClient()

		// todo options
		req := c.NewRequest(msg.Service(), node.Value, msg.GetPayload().GetBody())
		ctx := context.Background()

		err := c.Call(ctx, req, rep)
		if err != nil {
			return nil, err
		}
	}

	return rep, nil
}

func (p *manager) Start() error {

	for _, cpt := range p.manager.ListComponents() {
		if cpt.Started {
			continue
		}
		p.logger.Info("start_component", cpt.Name)

		if cpt.Component == nil {
			err := fmt.Errorf("component not found: %s", cpt.Name)
			p.logger.Error("start_component", cpt.Name, "err", "not found")
			return err
		}
		if err := cpt.Component.Start(); err != nil {
			p.logger.Error("start_component", cpt.Name, "err", err.Error())
			return err
		}
	}

	return p.router.Start()
}

func (p *manager) Stop() error {

	for _, cpt := range p.manager.ListComponents() {
		if cpt.Started {
			continue
		}
		if err := cpt.Component.Stop(); err != nil {
			p.logger.Error("stop_component", cpt.Name, "err", err.Error())
			return err
		}
	}
	return p.router.Stop()
}

func (p *manager) Router() router.Router {
	return p.router
}

func (p *manager) CompManager() component.Manager {
	return p.manager
}
