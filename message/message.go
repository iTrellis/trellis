package message

import (
	"github.com/go-trellis/trellis/codec"
	"github.com/go-trellis/trellis/message/proto"
)

type Message struct {
	proto.Payload `json:",inline"`

	Codec codec.Codec `json:"codec"`
}

func (p *Message) SetBody(obj interface{}) (err error) {

	p.Payload.Body, err = p.Codec.Marshal(obj)
	if err != nil {
		return err
	}

	return nil
}

func (p *Message) GetService() *proto.BaseService {
	return &proto.BaseService{Name: p.ServiceName, Version: p.ServiceVersion}
}

func (p *Message) ToObject(obj interface{}) error {
	codec, err := codec.GetCodec(p.Payload.GetContentType())
	if err != nil {
		return err
	}

	return codec.UnmarshalObject(p.Payload.GetBody(), obj)
}
