package message

import "github.com/go-trellis/trellis/service"

// Message is the interface for publishing asynchronously
type Message interface {
	Service() *service.Service
	Topic() string
	Payload() *BasePayload
	ContentType() string
}

// Payload payload between services
type Payload struct {
	ID       string
	Target   string
	Method   string
	Endpoint string
	Error    string

	BasePayload
}

// BasePayload payload between services
type BasePayload struct {
	// The values read from the socket
	Header map[string]string
	Body   []byte
}
