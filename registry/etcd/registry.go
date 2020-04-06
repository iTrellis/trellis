package etcd

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-trellis/trellis/message/proto"
	"github.com/go-trellis/trellis/registry"

	"github.com/go-trellis/etcdnaming"
)

func init() {
	registry.RegistRegistry(proto.RegistryType_ETCD, NewRegistry)
}

type etcdRegistry struct {
	Locker sync.RWMutex

	Registers map[string]etcdnaming.ServerRegister
}

// NewRegistry 生成对象
func NewRegistry() (registry.Registry, error) {
	r := &etcdRegistry{}
	return r, nil
}

func (p *etcdRegistry) GetType() proto.RegistryType {
	return proto.RegistryType_ETCD
}

func (p *etcdRegistry) RegistService(s *proto.Service) error {

	sbts, err := s.MarshalJSON()
	if err != nil {
		return err
	}
	cfg := etcdnaming.ServerRegisterConfig{
		Name:             s.BaseService.GetName(),
		Version:          s.BaseService.GetVersion(),
		TTL:              8,
		Interval:         time.Second * 12,
		RegistRetryTimes: 10,
	}

	cfg.Service = string(sbts)
	register := etcdnaming.NewDefaultServerRegister(cfg)

	p.Registers[genServiceNameByService(s)] = register

	return nil
}

// RevokeService 注销服务
func (p *etcdRegistry) RevokeService(s *proto.Service) error {

	register, ok := p.Registers[genServiceNameByService(s)]
	if !ok {
		return fmt.Errorf("not found register")
	}

	return register.Revoke()
}
func genServiceNameByService(s *proto.Service) string {
	return fmt.Sprintf("%s-%s", s.GetBaseService().GetName(), s.GetBaseService().GetVersion())
}

func genServiceName(name, version string) string {
	return fmt.Sprintf("%s-%s", name, version)
}
