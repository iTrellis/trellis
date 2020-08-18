package registry

import (
	"fmt"

	"github.com/go-trellis/trellis/configure"
)

// Registry the registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	// 注册的不只是服务本身，还需要第三方客户端的配置
	Regist(*configure.RegistService) error
	Revoke(*configure.RegistService) error
	Watch(WatchOption) (Watcher, error)

	Stop()
	String() string
}

// FullName fullname
func (p *WatchOption) FullName() string {
	return fmt.Sprintf("%s/%s", p.Service.GetName(), p.Service.GetVersion())
}
