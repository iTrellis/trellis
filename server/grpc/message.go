package grpc

import "github.com/iTrellis/trellis/service"

type Message struct {
	service service.Service
}

func (p *Message) Service() *service.Service {
	return &p.service
}

// Service() *service.Service
// Topic() string
// Payload() *BasePayload
// ContentType() string
