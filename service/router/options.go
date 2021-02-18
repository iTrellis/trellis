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
