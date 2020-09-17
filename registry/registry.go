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
	"time"

	"github.com/go-trellis/common/logger"
	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/message/proto"
)

// NewRegistryFunc 注册机生成函数
type NewRegistryFunc = func() Registry

// Registry the registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Init(*RegistOption) error
	// 注册的不只是服务本身，还需要第三方客户端的配置
	Regist(*configure.RegistService) error
	Revoke(*configure.RegistService) error
	Watcher(*configure.Watcher) (Watcher, error)

	Stop()
	String() string
}

// RegistOption the configure of registry
type RegistOption struct {
	RegisterType proto.RegisterType

	// registry url
	Endpoint string
	// expired time
	TTL time.Duration
	// Rotation time to regist serv into endpoint
	Heartbeat time.Duration
	// allow failed to regist server and retry times; -1 alaways retry
	RetryTimes uint32

	Logger logger.Logger
}
