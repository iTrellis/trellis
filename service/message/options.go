package message

import "github.com/go-trellis/trellis/service"

// Option used by NewMessage
type Option func(*Options)

// Options parameters
type Options struct {
	Service *service.Service
	Topic   string
	Payload *BasePayload

	ContentType string
}

func ContentType(ct string) Option {
	return func(o *Options) {
		o.ContentType = ct
	}
}

func Topic(topic string) Option {
	return func(o *Options) {
		o.Topic = topic
	}
}

func MessagePayload(payload *BasePayload) Option {
	return func(o *Options) {
		o.Payload = payload
	}
}

func Service(s *service.Service) Option {
	return func(o *Options) {
		o.Service = s
	}
}
