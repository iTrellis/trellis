package codec

import (
	"encoding/json"

	"github.com/go-trellis/trellis/message/proto"
)

type jsonCodec struct{}

func (*jsonCodec) Unmarshal(bytes []byte) (*proto.Payload, error) {
	resp := &proto.Payload{}
	err := json.Unmarshal(bytes, resp)
	return resp, err
}

func (*jsonCodec) UnmarshalObject(bytes []byte, obj interface{}) error {
	return json.Unmarshal(bytes, obj)
}

func (*jsonCodec) Marshal(body interface{}) ([]byte, error) {
	return json.Marshal(body)
}

func (*jsonCodec) String() string {
	return JSON
}

func newJSONCodec() (Codec, error) {
	return (*jsonCodec)(nil), nil
}
