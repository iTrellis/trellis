package local

import (
	"fmt"
	"sync"

	"github.com/go-trellis/trellis/message/proto"
	"github.com/go-trellis/trellis/registry"
)

func init() {
	registry.RegistRegistry(proto.RegistryType_LOCAL, NewRegistry)
}

type localRegistry struct {
	Locker   sync.RWMutex
	Services map[string]*proto.Service
}

// NewRegistry 生成注册
func NewRegistry() (registry.Registry, error) {
	return &localRegistry{
		Services: make(map[string]*proto.Service),
	}, nil
}

func (p *localRegistry) GetType() proto.RegistryType {
	return proto.RegistryType_LOCAL
}

func (p *localRegistry) RegistService(s *proto.Service) error {
	name := fmt.Sprintf("%s:%s", s.BaseService.GetName(), s.BaseService.GetVersion())

	p.Locker.Lock()
	defer p.Locker.Unlock()
	_, ok := p.Services[name]
	if ok {
		return fmt.Errorf("service is already exists: %s, %s", s.BaseService.GetName(), s.BaseService.GetVersion())
	}
	p.Services[name] = s
	return nil
}

// // GetService 获取名称对应的邮箱
// func (p *localRegistry) GetService(base *proto.BaseService) (registry.Service, bool) {

// 	p.Locker.RLock()
// 	defer p.Locker.RUnlock()
// 	s, ok := p.Services[fmt.Sprintf("%s:%s", base.Name, base.Version)]
// 	return s, ok
// }

// // AllServices 获取所有的邮箱
// func (p *localRegistry) AllServices() []registry.Service {

// 	p.Locker.RLock()
// 	defer p.Locker.RUnlock()
// 	var ss []registry.Service
// 	for k := range p.Services {
// 		ss = append(ss, p.Services[k])
// 	}
// 	return ss
// }

// RevokeService 注销服务
func (p *localRegistry) RevokeService(s *proto.Service) error {

	name := fmt.Sprintf("%s:%s", s.BaseService.GetName(), s.BaseService.GetVersion())
	p.Locker.Lock()
	defer p.Locker.Unlock()
	_, ok := p.Services[name]
	if !ok {
		return fmt.Errorf("service not exists: %s, %s", s.BaseService.GetName(), s.BaseService.GetVersion())
	}
	delete(p.Services, name)
	return nil
}

// // Watch 不存在监控，无实现
// func (p *localRegistry) Watch(...registry.WatchOptionFunc) (registry.Watcher, error) {
// 	return nil, nil
// }
