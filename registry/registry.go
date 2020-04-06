// GNU GPL v3 License
// Copyright (c) 2016 github.com:go-trellis

package registry

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-trellis/trellis/message/proto"
)

// trellis://registry/service/ServerName/ServerVersion

var registries map[proto.RegistryType]Registry = make(map[proto.RegistryType]Registry)
var servicesType map[*proto.BaseService]proto.RegistryType = make(map[*proto.BaseService]proto.RegistryType)

// var services map[*proto.BaseService]*proto.Service = make(map[*proto.BaseService]*proto.Service)

// Registry regist or revoke server interface
type Registry interface {
	GetType() proto.RegistryType
	// RegistService 注册服务
	RegistService(*proto.Service) error
	// // GetService 获取名称对应的邮箱
	// GetService(*proto.BaseService) (Service, bool)
	// // AllServices 获取所有的邮箱
	// AllServices() []Service
	// RevokeService 注销服务
	RevokeService(*proto.Service) error
	// // Watch 监测注册的服务
	// Watch(...WatchOptionFunc) (Watcher, error)
}

// NewRegistryFunc 生成注册函数
type NewRegistryFunc func() (Registry, error)

// RegistRegistry 注册
func RegistRegistry(typ proto.RegistryType, fn NewRegistryFunc) {

	if fn == nil {
		panic("registry fn is nil")
	}

	_, exist := registries[typ]

	if exist {
		panic(fmt.Sprintf("registry already registered: %s", typ))
	}

	r, err := fn()
	if err != nil {
		panic(err)
	}

	registries[typ] = r
}

// RegistService 注册服务
func RegistService(s *proto.Service) error {
	r, err := GetRegistry(s.Type)
	if err != nil {
		return err
	}
	return r.RegistService(s)
}

// RegistServiceByPath 注册服务
func RegistServiceByPath(path string, endpoint string, metadata map[string]string) error {
	endpoint = strings.TrimSpace(endpoint)
	if len(endpoint) == 0 {
		return errors.New("service's endpoint must not be empty")
	}

	s, err := ParseServicePath(path)
	if err != nil {
		return err
	}

	s.Endpoint = endpoint
	s.Metadata = metadata

	r, err := GetRegistry(s.GetType())
	if err != nil {
		return err
	}

	// s, err := NewService(ServiceBaseService(base), ServiceMetadata(metadata), ServiceNodes(n))
	// if err != nil {
	// 	return err
	// }

	return r.RegistService(s)
}

// GetRegistry 通过名字获取注册对象
func GetRegistry(typ proto.RegistryType) (Registry, error) {
	r, ok := registries[typ]
	if ok {
		return r, nil
	}
	return nil, fmt.Errorf("not found registry: %s", typ)
}

// GetService 通过注册机类型和服务名字获取服务对象
func GetService(base *proto.BaseService) (*proto.Service, error) {
	if base == nil || base.Name == "" {
		return nil, errors.New("unknown service's name")
	}
	// r, err := GetRegistry(base.Type)
	// if err != nil {
	// 	return nil, err
	// }
	// s, ok := r.GetService(base)
	// if !ok {
	// 	return nil, fmt.Errorf("not found service: %d - %s:%s", base.Type, base.Name, base.Version)
	// }

	// t, ok := servicesType[base]
	// if !ok {
	// 	return nil, fmt.Errorf("not found service: %s:%s", base.Name, base.Version)
	// }

	return nil, nil
}

// // GetServiceByPath 通过路径获取服务信息
// // trellis://registry/service/proto.RegistryType/ServiceName/ServiceVersion
// func GetServiceByPath(path string) (Service, error) {
// 	base, typ, err := ParseServicePath(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return GetService(base)
// }

// ParseServicePath 解析trellis的基础服务路径
// trellis://registry/service/proto.RegistryType/ServiceName/ServiceVersion
func ParseServicePath(path string) (*proto.Service, error) {
	if !strings.HasPrefix(path, "trellis://registry/service/") {
		return nil, fmt.Errorf("prefix is not trellis://registry/service/")
	}

	ss := strings.Split(strings.TrimLeft(path, "trellis://registry/service/"), "/")
	lenSS := len(ss)
	base := &proto.BaseService{}
	if lenSS < 2 {
		return nil, fmt.Errorf("path is incorrect")
	} else if lenSS > 3 {
		base.Version = ss[2]
	}

	typ, ok := proto.RegistryType_value[ss[0]]
	if !ok {
		return nil, fmt.Errorf("not found registry: %s", ss[0])
	}

	base.Name = ss[1]

	s := &proto.Service{
		BaseService: base,
		Type:        proto.RegistryType(typ),
	}

	return s, nil
}
