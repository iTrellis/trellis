package message

import (
	"github.com/go-trellis/trellis/codec"
	"github.com/go-trellis/trellis/message/proto"

	"github.com/google/uuid"
)

// RequestOptionFunc 请求结构体的参数赋值方法
type RequestOptionFunc func(*Request)

type Request struct {
	Endpoint string
	Method   string

	body interface{}

	*Message
}

func (p *Request) initMessage() {
	p.Message.Payload = &proto.Payload{
		ID:     uuid.New().String(),
		Header: make(map[string]string),
		Server: &proto.BaseService{},
	}

}

// NewRequest 新生成请求体
func NewRequest(opts ...RequestOptionFunc) (req *Request, err error) {
	r := &Request{Message: &Message{}}
	r.initMessage()

	for _, o := range opts {
		o(r)
	}

	r.codec, err = codec.GetCodec(r.Payload.ContentType)
	if err != nil {
		return
	}

	if r.body != nil {
		r.Payload.Body, err = r.codec.Marshal(r.body)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

// RequestServer 请求的服务
func RequestServer(server string) RequestOptionFunc {
	return func(r *Request) {
		r.Payload.Server.Name = server
	}
}

// RequestVersion 请求的版本号
func RequestVersion(version string) RequestOptionFunc {
	return func(r *Request) {
		r.Payload.Server.Version = version
	}
}

// RequestEndpoint 请求的端点
func RequestEndpoint(endpoint string) RequestOptionFunc {
	return func(r *Request) {
		r.Endpoint = endpoint
	}
}

// RequestMethod 请求的方法
func RequestMethod(method string) RequestOptionFunc {
	return func(r *Request) {
		r.Method = method
	}
}

// RequestTopic 请求的方法
func RequestTopic(topic string) RequestOptionFunc {
	return func(r *Request) {
		r.Payload.Topic = topic
	}
}

// RequestContentType 请求头类型
func RequestContentType(contentType string) RequestOptionFunc {
	return func(r *Request) {
		r.Payload.ContentType = contentType
	}
}

// RequestPayload 请求体内容
func RequestPayload(body interface{}) RequestOptionFunc {
	return func(r *Request) {
		r.body = body
	}
}

// RequestHeader 请求头信息
func RequestHeader(header map[string]string) RequestOptionFunc {
	return func(r *Request) {
		r.Payload.Header = header
	}
}

func (p *Request) ID() string {
	return p.Payload.GetID()
}

func (p *Request) Server() *proto.BaseService {
	return p.Payload.GetServer()
}

func (p *Request) Read() ([]byte, error) {
	return p.codec.Marshal(p.Payload)
}

func (p *Request) Topic() string {
	return p.Payload.Topic
}

func (p *Request) SetHeader(key, value string) {
	p.Payload.Header[key] = value
}

func (p *Request) GetMessage() *Message {
	return p.Message
}
