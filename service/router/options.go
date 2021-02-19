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

package router

import (
	"time"

	"github.com/google/uuid"
	"github.com/iTrellis/trellis/service"
)

// Option options' of registing service functions
type Option func(*Options)

func ID(id string) Option {
	return func(o *Options) {
		o.ID = id
	}
}
func TTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.TTL = ttl
	}
}

func Heartbeat(hb time.Duration) Option {
	return func(o *Options) {
		o.Heartbeat = hb
	}
}

func RetryTimes(rTimes uint32) Option {
	return func(o *Options) {
		o.RetryTimes = rTimes
	}
}

// Options regist service Options
type Options struct {
	ID string

	TTL time.Duration

	Heartbeat time.Duration

	RetryTimes uint32
}

// DefaultOptions returns router default options
func DefaultOptions() Options {
	return Options{
		ID:        uuid.New().String(),
		Heartbeat: 10 * time.Second,
	}
}

type ReadOption func(o *ReadOptions)

type ReadOptions struct {
	Service *service.Service
	Keys    []string
}

// ReadService sets the service to read from the table
func ReadService(s *service.Service) ReadOption {
	return func(o *ReadOptions) {
		o.Service = s
	}
}

// ReadKeys if need keys for node
func ReadKeys(keys ...string) ReadOption {
	return func(o *ReadOptions) {
		o.Keys = keys
	}
}
