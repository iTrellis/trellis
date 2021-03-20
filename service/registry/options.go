/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package registry

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/iTrellis/common/logger"
	"github.com/iTrellis/trellis/service"
)

// Option initial options' functions
type Option func(*Options)

// Options new registry Options
type Options struct {
	Endpoints []string
	Timeout   time.Duration
	Secure    bool
	TLSConfig *tls.Config

	ServerAddr string
	RetryTimes uint32

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	Logger logger.Logger
}

func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

func Endpoints(endpoints []string) Option {
	return func(o *Options) {
		o.Endpoints = endpoints
	}
}

func Timeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

func Secure(secure bool) Option {
	return func(o *Options) {
		o.Secure = secure
	}
}

func TLSConfig(tlsConfig *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = tlsConfig
	}
}

func ServerAddr(addr string) Option {
	return func(o *Options) {
		o.ServerAddr = addr
	}
}

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func RetryTimes(rTimes uint32) Option {
	return func(o *Options) {
		o.RetryTimes = rTimes
	}
}

func RegisterWeight(w uint32) RegisterOption {
	return func(o *RegisterOptions) {
		o.Weight = w
	}
}

func RegisterTTL(ttl time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.TTL = ttl
	}
}

func RegisterHeartbeat(hb time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.Heartbeat = hb
	}
}

// RegisterOption options' of registing service functions
type RegisterOption func(*RegisterOptions)

// RegisterOptions regist service Options
type RegisterOptions struct {
	TTL       time.Duration
	Heartbeat time.Duration
	Weight    uint32
}

func (p *RegisterOptions) Check() {
	if p.Heartbeat == 0 {
		p.Heartbeat = 10 * time.Second
	}
	if p.Weight == 0 {
		p.Weight = 1
	}
}

// DeregisterOption options' of deregistering service functions
type DeregisterOption func(*DeregisterOptions)

// DeregisterOptions deregister service Options
type DeregisterOptions struct {
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
	// // Other options for implementations of the interface
	// // can be stored in a context
	// Options config.Options

	Logger logger.Logger
}

func WatchLogger(l logger.Logger) WatchOption {
	return func(w *WatchOptions) {
		w.Logger = l
	}
}

func WatchService(service service.Service) WatchOption {
	return func(w *WatchOptions) {
		w.Service = service
	}
}

// func WatchContext(opts config.Options) WatchOption {
// 	return func(w *WatchOptions) {
// 		w.Options = opts
// 	}
// }

// GetOption options' of getting service functions
type GetOption func(*GetOptions)

// GetOptions get service Options
type GetOptions struct {
	Context context.Context
}
