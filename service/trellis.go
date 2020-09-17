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

package service

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-trellis/trellis/clients"
	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/errcode"
	"github.com/go-trellis/trellis/internal"
	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/message/proto"
	"github.com/go-trellis/trellis/registry"

	"github.com/go-trellis/common/errors"
	"github.com/go-trellis/common/formats"
	"github.com/go-trellis/common/logger"
	"github.com/go-trellis/node"

	// 注册机
	_ "github.com/go-trellis/trellis/registry/cache"
	_ "github.com/go-trellis/trellis/registry/etcd"
)

// Trellis 启动对象
type Trellis struct {
	conf   *configure.Project
	logger logger.Logger

	services map[string]Service

	opts RouterOptions
}

// var trellis Router

// NewTrellis 生成启动对象
func NewTrellis(cfg *configure.Project, l logger.Logger) (Router, error) {

	t := &Trellis{
		conf:     cfg,
		services: make(map[string]Service),

		logger: l,
	}

	if err := t.initRegistries(); err != nil {
		t.logger.Error("init_registries_failed", err)
		return nil, err
	}

	if err := t.newServices(); err != nil {
		t.logger.Error("new_services_failed", err)
		return nil, err
	}
	t.logger.Info("new services ok")

	if err := t.registServices(); err != nil {
		t.logger.Error("regist_services_failed", err)
		return nil, err
	}

	return t, nil
}

// BlockStop 阻断式停止
func BlockStop() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	case <-ch:
	}
	return
}

// Run 运行
func Run(cfg *configure.Project, l logger.Logger) (Router, error) {
	trellis, err := NewTrellis(cfg, l)
	if err != nil {
		return nil, err
	}

	err = trellis.Run()
	if err != nil {
		return nil, err
	}

	return trellis, nil
}

func (p *Trellis) newServices() error {

	for name, service := range p.conf.Services {
		service.Name = name

		p.logger.Debug("new service", service.String())

		err := p.NewService(RouterOptionService(service), RouterOptionLogger(p.logger))
		if err != nil {
			return err
		}

		path := internal.WorkerTrellisPath(service.GetName(), service.GetVersion())

		nm := node.NewDirect(path)
		nm.Add(&node.Node{
			ID:       path,
			Weight:   1,
			Value:    service.String(),
			Metadata: map[string]interface{}{"protocol": proto.Protocol_LOCAL},
		})

		registry.SetNodeManager(path, nm)
	}

	clients.RegistCaller(proto.Protocol_LOCAL, p)
	return nil
}

// NewService new service
func (p *Trellis) NewService(opts ...RouterOptionFunc) (err error) {

	for _, o := range opts {
		o(&p.opts)
	}

	url := internal.WorkerTrellisPath(p.opts.Service.GetName(), p.opts.Service.GetVersion())
	if _, ok := p.services[url]; ok {
		err = fmt.Errorf("%s already exists", url)
		p.opts.Logger.Error("new_service_failed", err.Error())
		return err
	}

	s, err := New(p.opts.Service.GetName(), p.opts.Service.GetVersion(),
		Config(p.opts.Service.Options),
		Logger(p.opts.Logger.With(url)),
	)
	if err != nil {
		p.opts.Logger.Error("new_service_failed", err.Error())
		return err
	}

	p.services[url] = s

	return nil
}

func (p *Trellis) registServices() error {
	p.logger.Info("regist service start")
	for name, service := range p.conf.Services {
		service.Name = name

		if service.Registry == nil {
			continue
		}
		p.logger.Info("regist service start", name, service.Registry)

		regConf := &configure.RegistService{
			Name:     service.Name,
			Version:  service.GetVersion(),
			Domain:   service.Registry.Domain,
			Protocol: service.Registry.Protocol,
			Weight:   service.Registry.Weight,
		}
		if err := registry.RegistService(service.Registry.Name, regConf); err != nil {
			p.logger.Error("regist service failed", regConf, err)
			return err
		}
	}

	return nil
}

// Run 启动进程
func (p *Trellis) Run() error {
	var errs errors.Errors
	for _, s := range p.services {
		err := s.Start()
		if err != nil {
			errs.Append(err)
		}
	}
	if len(errs) != 0 {
		p.opts.Logger.Error("run_service_failed", errs.Error())
		return errs
	}
	return nil
}

// Stop 停止工作者
func (p *Trellis) Stop() error {
	var errs errors.Errors
	for _, s := range p.services {
		err := s.Stop()
		if err != nil {
			errs.Append(err)
		}
	}
	p.services = nil

	if len(errs) != 0 {
		p.opts.Logger.Error("stop_service_failed", errs.Error())
		return errs
	}

	return nil
}

// initRegistries 启动注册器
func (p *Trellis) initRegistries() (err error) {

	for name, reg := range p.conf.Registries {
		retryTimes, _ := reg.Options.Int("retry_times")
		rOpts := &registry.RegistOption{
			RegisterType: proto.RegisterType(proto.RegisterType_value[strings.ToUpper(reg.Type)]),
			Endpoint:     reg.Options.Get("endpoint"),
			TTL:          formats.ParseStringTime(reg.Options.Get("ttl")),
			Heartbeat:    formats.ParseStringTime(reg.Options.Get("heartbeat")),
			RetryTimes:   uint32(retryTimes),
			Logger:       p.logger,
		}
		p.logger.Debug("new registry", rOpts)

		if err := registry.NewRegistry(name, rOpts); err != nil {
			p.logger.Error("failed new registry", rOpts.Endpoint, err)
			return err
		}

		for _, wConfig := range reg.Watchers {

			if err = registry.NewRegistryWatcher(name, wConfig); err != nil {
				p.logger.Error("new watcher failed", *wConfig, err)
				return err
			}
		}

		p.logger.Info("initial registry ok", name, reg)

	}
	return nil
}

// StopService stop service
func (p *Trellis) StopService(name, version string) error {
	url := internal.WorkerTrellisPath(name, version)
	s, ok := p.services[url]
	if !ok {
		err := fmt.Errorf("unknown service: %s, %s", name, version)
		p.opts.Logger.Error("stop_service_failed", err.Error())
		return err
	}
	if err := s.Stop(); err != nil {
		p.opts.Logger.Error("stop_service_failed", err.Error())
		return err
	}

	delete(p.services, url)
	return nil
}

// GetService get service
func (p *Trellis) GetService(name, version string) (Service, error) {
	s, ok := p.services[internal.WorkerTrellisPath(name, version)]
	if !ok {
		return nil, fmt.Errorf("unknown service: %s, %s", name, version)
	}
	return s, nil
}

// CallService call service
func (p *Trellis) CallService(_ *node.Node, msg *message.Message) (interface{}, error) {
	s, err := p.GetService(msg.GetService().GetName(), msg.GetService().GetVersion())
	if err != nil {
		return nil, err
	}

	fn := s.Route(msg.GetTopic())
	if fn == nil {
		return nil, errcode.ErrGetServiceTopic.New()
	}
	return fn(msg)
}
