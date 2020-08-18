package registry

import (
	"errors"
	"time"

	"github.com/go-trellis/trellis/message/proto"
)

// RegistOption the configure of registry
type RegistOption struct {
	// registry url
	Endpoint string
	// expired time
	TTL time.Duration
	// Rotation time to regist serv into endpoint
	Heartbeat time.Duration
	// allow failed to regist server and retry times; -1 alaways retry
	RetryTimes uint32
}

type WatchOption struct {
	proto.Service `yaml:",inline"`
}

type GetOptionFunc func(*GetOption)

type GetOption struct {
	proto.Service
}

func (p *GetOption) check() error {
	if p.GetName() == "" {
		return errors.New("unknown service name")
	}
	return nil
}

func GetServiceVersion(version string) GetOptionFunc {
	return func(o *GetOption) {
		o.Version = version
	}
}
