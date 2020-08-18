package codec

import (
	"encoding/json"
)

type jsonCodec struct{}

func (*jsonCodec) Unmarshal(bytes []byte, obj interface{}) error {
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
