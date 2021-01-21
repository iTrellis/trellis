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

package transport

// import (
// 	"time"

// 	"github.com/iTrellis/trellis"
// )

// const (
// 	proxyAuthHeader = "Proxy-Authorization"
// )

// var (
// 	// DefaultTransport Transport = NewHTTPTransport()

// 	DefaultDialTimeout = time.Second * 5
// )

// // Transport is an interface which is used for communication between
// // services. It uses connection based socket send/recv semantics and
// // has various implementations; http, grpc, quic.
// type Transport interface {
// 	Init(...Option) error
// 	Options() Options
// 	Dial(addr string, opts ...DialOption) (Client, error)
// 	Listen(addr string, opts ...ListenOption) (Listener, error)
// 	String() string
// }

// type Socket interface {
// 	Recv(*trellis.BasePayload) error
// 	Send(*trellis.BasePayload) error
// 	Close() error
// 	Local() string
// 	Remote() string
// }

// type Client interface {
// 	Socket
// }

// type Listener interface {
// 	Addr() string
// 	Close() error
// 	Accept(func(Socket)) error
// }

// type Option func(*Options)

// type DialOption func(*DialOptions)

// type ListenOption func(*ListenOptions)
