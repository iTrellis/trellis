package grpc

var (
	// DefaultPoolMaxStreams maximum streams on a connectioin
	// (20)
	DefaultPoolMaxStreams = 20

	// DefaultPoolMaxIdle maximum idle conns of a pool
	// (50)
	DefaultPoolMaxIdle = 50

	// DefaultMaxRecvMsgSize maximum message that client can receive
	// (16 MB).
	DefaultMaxRecvMsgSize = 1024 * 1024 * 16

	// DefaultMaxSendMsgSize maximum message that client can send
	// (16 MB).
	DefaultMaxSendMsgSize = 1024 * 1024 * 16
)

type poolMaxStreams struct{}
type poolMaxIdle struct{}
