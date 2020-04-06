// GNU GPL v3 License
// Copyright (c) 2016 github.com:go-trellis

package registry

import (
	"time"

	"github.com/go-trellis/trellis/message/proto"
)

// import (
// 	"context"

// 	"github.com/go-trellis/trellis/message/proto"
// )

type ResolverOptionFunc func(*ResolverOptions)

type ResolverOptions struct {
	// Specify a service to Resolver
	// If blank, the Resolver is for all services
	proto.BaseService
	// Other options for implementations of the interface
	// can be stored in a context
	Timeout time.Duration

	Target string
}

// ResolverService Resolver a service
func ResolverService(name string) ResolverOptionFunc {
	return func(o *ResolverOptions) {
		o.BaseService.Name = name
	}
}

// ResolverVersion Resolver a service
func ResolverVersion(version string) ResolverOptionFunc {
	return func(o *ResolverOptions) {
		o.Version = version
	}
}

// ResolverTimeout timeout
func ResolverTimeout(timeout time.Duration) ResolverOptionFunc {
	return func(o *ResolverOptions) {
		o.Timeout = timeout
	}
}

// ResolverTarget registry target
func ResolverTarget(target string) ResolverOptionFunc {
	return func(o *ResolverOptions) {
		o.Target = target
	}
}

// Resolver Resolver the services
type Resolver interface {
	Next() (*Result, error)
	Stop()
}

// Result service info
type Result struct {
	Action  string
	Service *proto.Service
}
