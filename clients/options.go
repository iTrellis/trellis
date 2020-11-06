package clients

import (
	"time"

	"github.com/go-trellis/config"
)

// OptionFunc used by the Client
type OptionFunc func(*Options)

// Options
type Options struct {
	// Used to select codec
	ContentType string
	// Proxy address to send requests via
	Proxy string

	// // Plugged interfaces
	// Broker    broker.Broker
	// Codecs    map[string]codec.NewCodec
	// Router    router.Router
	// Selector  selector.Selector
	// Transport transport.Transport

	// // Lookup used for looking up routes
	// Lookup LookupFunc

	// Connection Pool
	PoolSize int
	PoolTTL  time.Duration

	// // Middleware for client
	// Wrappers []Wrapper

	// Default Call Options
	CallOptions CallOptions

	// Other options for implementations of the interface
	// can be stored in a context
	Options config.Options
}

type CallOptions struct {
	// Address of remote hosts
	Address []string
	// // Backoff func
	// Backoff BackoffFunc
	// Transport Dial Timeout
	DialTimeout time.Duration
	// Number of Call attempts
	Retries int
	// // Check if retriable func
	// Retry RetryFunc
	// Request/Response timeout
	RequestTimeout time.Duration
	// // Router to use for this call
	// Router router.Router
	// // Selector to use for the call
	// Selector selector.Selector
	// // SelectOptions to use when selecting a route
	// SelectOptions []selector.SelectOption
	// // Stream timeout for the stream
	// StreamTimeout time.Duration
	// Use the auth token as the authorization header
	AuthToken bool
	// Network to lookup the route within
	Network string

	// // Middleware for low level call func
	// CallWrappers []CallWrapper

	// Other options for implementations of the interface
	// can be stored in a context
	Options config.Options
}
