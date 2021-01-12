package registry

import "github.com/go-trellis/common/logger"

// NewRegistryFunc new registry function
type NewRegistryFunc func(logger logger.Logger, opts ...Option) (Registry, error)

// Registry The registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Init(...Option) error
	Options() Options

	Regist(*Service, ...RegisterOption) error
	Revoke(*Service, ...RevokeOption) error

	Watch(...WatchOption) (Watcher, error)

	ID() string
	String() string
}
