package registry

import "github.com/go-trellis/trellis/message/proto"

// Registry the registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Regist(*proto.BaseService) error
	Revoke(*proto.BaseService) error
	Name() string
	Type() string
	// Deregister(*Service, ...DeregisterOption) error
	// GetService(string, ...GetOption) ([]*Service, error)
	// ListServices(...ListOption) ([]*Service, error)
	// Watch(...WatchOption) (Watcher, error)
	// String() string
}
