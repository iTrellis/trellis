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

	"github.com/iTrellis/trellis/service"
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
