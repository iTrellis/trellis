package registry

import (
	"fmt"
	"time"

	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/message/proto"

	"github.com/go-trellis/common/logger"
)

// NewRegistryFunc 注册机生成函数
type NewRegistryFunc = func() Registry

var mapRegistries = make(map[proto.RegisterType]NewRegistryFunc)

// Regist 注册注册机
func Regist(name proto.RegisterType, fn NewRegistryFunc) {
	_, ok := mapRegistries[name]
	if ok {
		panic(fmt.Errorf("registry'name (%s) is already exist", name))
	}
	mapRegistries[name] = fn
}

// GetNewRegistryFunc 获取注册机生成函数
func GetNewRegistryFunc(name proto.RegisterType) (NewRegistryFunc, error) {
	r, ok := mapRegistries[name]
	if !ok {
		return nil, fmt.Errorf("registry'name (%s) isnot exist", name)
	}
	return r, nil
}

// Registry the registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Init(option *RegistOption, log logger.Logger) error
	// 注册的不只是服务本身，还需要第三方客户端的配置
	Regist(*configure.RegistService) error
	Revoke(*configure.RegistService) error
	Watcher(*configure.Watcher) (Watcher, error)

	Stop()
	String() string
}

// RegistOption the configure of registry
type RegistOption struct {
	// registry url
	Endpoint string
	// expired time
	TTL time.Duration
	// Rotation time to regist serv into endpoint
	Heartbeat time.Duration
	// allow failed to regist server and retry times; -1 alaways retry
	RetryTimes uint32
}
