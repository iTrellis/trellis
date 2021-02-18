package etcd

import (
	"time"

	"github.com/iTrellis/trellis/service/registry"
)

type worker struct {
	service *registry.Service

	options registry.RegisterOptions

	// client *clientv3.Client

	fullServiceName string
	fullpath        string

	stopSignal chan bool

	// invoke self-register with ticker
	ticker *time.Ticker

	interval time.Duration
}
