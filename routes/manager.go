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
	"fmt"
	"reflect"

	"github.com/iTrellis/common/logger"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/message"
)

// Manager routes manager
type Manager interface {
	Init(...Option)

	service.LifeCycle

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
	manager component.Manager
	logger  logger.Logger
}

func (p *manager) Init(opts ...Option) {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	if options.manager != nil {
		p.manager = options.manager
	}

	if p.manager == nil {
		p.manager = NewCompManager()
	}

	p.logger = options.logger
}

func (p *manager) CallComponent(msg message.Message) (interface{}, error) {

	cpt, err := p.manager.GetComponent(msg.Service())
	if err != nil {
		return nil, err
	} else if cpt == nil {
		return nil, fmt.Errorf("unknown component")
	}
	p.logger.Debug("call_component",
		"component", msg.Service().TrellisPath(), "topic", msg.Topic(), "component_type", reflect.TypeOf(cpt))

	return cpt.Route(msg)
}

func (p *manager) Start() (err error) {

	for _, cpt := range p.manager.ListComponents() {
		if !cpt.Started {
			continue
		}
		p.logger.Info("start_component", "component", cpt.Name)

		if cpt.Component == nil {
			err = fmt.Errorf("component not found: %s", cpt.Name)
			p.logger.Error("failed_start_component", "component", cpt.Name, "err", err.Error())
			return
		}

		if err = cpt.Component.Start(); err != nil {
			p.logger.Error("failed_start_component", "component", cpt.Name, "err", err.Error())
			return
		}
		p.logger.Error("start_component", "component", cpt.Name, "result", "ok")
	}

	return nil
}

func (p *manager) Stop() error {

	for _, cpt := range p.manager.ListComponents() {
		if !cpt.Started {
			continue
		}
		if err := cpt.Component.Stop(); err != nil {
			p.logger.Error("stop_component", "component", cpt.Name, "err", err.Error())
			return err
		}
	}
	// return p.router.Stop()
	return nil
}

func (p *manager) CompManager() component.Manager {
	return p.manager
}
