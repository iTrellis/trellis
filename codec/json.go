package codec

import (
	"encoding/json"

	"github.com/go-trellis/trellis/message/proto"
)

type jsonCodec struct{}

func (*jsonCodec) Unmarshal(bytes []byte) (*proto.Payload, error) {
	resp := &proto.Payload{}
	err := resp.UnmarshalJSON(bytes)
	return resp, err
}

func (*jsonCodec) UnmarshalObject(bytes []byte, obj interface{}) error {
	return json.Unmarshal(bytes, obj)
}

func (*jsonCodec) Marshal(body interface{}) ([]byte, error) {
	switch t := body.(type) {
	case *proto.Payload:
		return t.MarshalJSON()
	}
	return json.Marshal(body)
}

func (*jsonCodec) String() string {
	return JSON
}

func newJSONCodec() (Codec, error) {
	return (*jsonCodec)(nil), nil
}
