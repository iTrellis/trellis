package codec

// Codec
type Codec interface {
	Reader
	Writer
	Close() error
	String() string
}

// Reader
type Reader interface {
	ReadHeader(*Payload) error
	ReadBody(interface{}) error
}

// Writer
type Writer interface {
	Write(*Payload, interface{}) error
}

// Marshaler is a simple encoding interface used for the broker/transport
// where headers are not supported by the underlying implementation.
type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}

// Payload payload between services
type Payload struct {
	// The values read from the socket
	Header map[string]string
	Body   []byte
}
