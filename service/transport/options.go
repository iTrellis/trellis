package transport

// import (
// 	"context"
// 	"crypto/tls"
// 	"time"

// 	"github.com/go-trellis/trellis/service/codec"
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
