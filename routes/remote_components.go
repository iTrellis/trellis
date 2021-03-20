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
	"encoding/json"
	"errors"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/iTrellis/node"
	"github.com/iTrellis/trellis/server"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/client/grpc"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/message"
	"github.com/iTrellis/trellis/service/registry"
)

type RemoteComponent interface {
	Init(opts ...component.Option)

	component.Component
}

type remoteComponents struct {
	sync.RWMutex

	reg registry.Registry

	s *service.Service

	options  component.Options
	wOpts    []registry.WatchOption
	woptions registry.WatchOptions

	nodeManager node.Manager
}

func NewRemoteComponent(nodeType node.Type, r registry.Registry, wOpts ...registry.WatchOption) (
	RemoteComponent, error) {
	c := &remoteComponents{
		reg:   r,
		wOpts: wOpts,
	}
	for _, o := range wOpts {
		o(&c.woptions)
	}
	var err error
	c.nodeManager, err = node.New(nodeType, c.woptions.Service.TrellisPath())
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (p *remoteComponents) Init(opts ...component.Option) {
	for _, o := range opts {
		o(&p.options)
	}
	return
}

func (p *remoteComponents) Start() error {
	w, err := p.reg.Watch(p.wOpts...)
	if err != nil {
		return err
	}
	go func() {
		for {
			result, err := w.Next()
			if err != nil {
				p.options.Logger.Warn("msg", "failed_get_next_node", "error", err)
				return
			}

			if result.Service == nil {
				continue
			}

			p.options.Logger.Debugf("watch nodes: %+v, %+v\n", *result, *result.Service.Node)

			switch result.Type {
			case service.EventType_create, service.EventType_update:
				p.nodeManager.Add(result.Service.Node)
			case service.EventType_delete:
				p.nodeManager.RemoveByID(result.Service.Node.ID)
			}
		}
	}()
	return nil
}

func (p *remoteComponents) Stop() error {
	return p.reg.Deregister(p.s)
}

func (p *remoteComponents) Route(msg message.Message) (interface{}, error) {
	nd, ok := p.nodeManager.NodeFor(msg.Topic(), msg.GetPayload().Get(service.HeaderXClientIP))
	if !ok || nd == nil {
		err := errors.New("not found remote server to call")
		return nil, err
	}

	var rep interface{}

	protocol := nd.Metadata["protocol"]
	if protocol == nil {
		protocol = service.Protocol_HTTP
	}

	switch protocol {
	case service.Protocol_HTTP:
		client := resty.New()

		remoteMsg := msg.ToRemoteMessage()

		req := client.NewRequest().SetBody(remoteMsg)

		resp, err := req.Post(nd.Value)
		if err != nil {
			return nil, err
		}

		r := &server.Response{}

		err = json.Unmarshal(resp.Body(), r)

		if err != nil {
			return nil, err
		}

		return r.Result, nil
	case service.Protocol_GRPC:
		fallthrough
	default:
		c := grpc.NewClient()

		// todo options
		req := c.NewRequest(msg.Service(), nd.Value, msg.GetPayload().GetBody())
		ctx := context.Background()

		err := c.Call(ctx, req, rep)
		if err != nil {
			return nil, err
		}
	}

	return rep, nil
}
