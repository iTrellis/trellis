package transport

// import (
// 	"time"

// 	"github.com/go-trellis/trellis"
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
