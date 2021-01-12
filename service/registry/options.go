package registry

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/go-trellis/trellis/service"
)

// Option initial options' functions
type Option func(*Options)

// Options new registry Options
type Options struct {
	Addrs     []string
	Timeout   time.Duration
	Secure    bool
	TLSConfig *tls.Config
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// RegisterOption options' of registing service functions
type RegisterOption func(*RegisterOptions)

func RegisterTTL(ttl time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.TTL = ttl
	}
}

func RegisterContext(ctx context.Context) RegisterOption {
	return func(o *RegisterOptions) {
		o.Context = ctx
	}
}

// RegisterOptions regist service Options
type RegisterOptions struct {
	TTL time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// RevokeOption options' of revoking service functions
type RevokeOption func(*RevokeOptions)

// RevokeOptions revoke service Options
type RevokeOptions struct {
	TTL time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// WatchOption options' of watching service functions
type WatchOption func(*WatchOptions)

// WatchOptions watch service Options
type WatchOptions struct {
	Service service.Service
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

func WatchService(service service.Service) WatchOption {
	return func(w *WatchOptions) {
		w.Service = service
	}
}

// GetOption options' of getting service functions
type GetOption func(*GetOptions)

// GetOptions get service Options
type GetOptions struct {
	Context context.Context
}
