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
// 	"context"
// 	"crypto/tls"
// 	"time"

// 	"github.com/iTrellis/trellis/service/codec"
// )

// type Options struct {
// 	// Addrs is the list of intermediary addresses to connect to
// 	Addrs []string
// 	// Codec is the codec interface to use where headers are not supported
// 	// by the transport and the entire payload must be encoded
// 	Codec codec.Marshaler
// 	// Secure tells the transport to secure the connection.
// 	// In the case TLSConfig is not specified best effort self-signed
// 	// certs should be used
// 	Secure bool
// 	// TLSConfig to secure the connection. The assumption is that this
// 	// is mTLS keypair
// 	TLSConfig *tls.Config
// 	// Timeout sets the timeout for Send/Recv
// 	Timeout time.Duration
// 	// Other options for implementations of the interface
// 	// can be stored in a context
// 	Context context.Context
// }

// type DialOptions struct {
// 	// Tells the transport this is a streaming connection with
// 	// multiple calls to send/recv and that send may not even be called
// 	Stream bool
// 	// Timeout for dialing
// 	Timeout time.Duration

// 	// TODO: add tls options when dialling
// 	// Currently set in global options

// 	// Other options for implementations of the interface
// 	// can be stored in a context
// 	Context context.Context
// }

// type ListenOptions struct {
// 	// TODO: add tls options when listening
// 	// Currently set in global options

// 	// Other options for implementations of the interface
// 	// can be stored in a context
// 	Context context.Context
// }
