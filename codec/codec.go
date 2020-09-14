package codec

import (
	"fmt"
)

var (
	codecs map[string]NewCodecFunc = make(map[string]NewCodecFunc)

	defaultCodecs map[string]Codec = make(map[string]Codec)
)

// codec的类型申明
const (
	JSON = "application/json"
)

func init() {
	RegisterCodec(JSON, newJSONCodec)
}

// NewCodecFunc 生成编码器方法
type NewCodecFunc func() (Codec, error)

// Codec 编码器
type Codec interface {
	Unmarshal([]byte, interface{}) error
	// Marshal(*proto.Payload) ([]byte, error)
	Marshal(interface{}) ([]byte, error)
	String() string
}

// GetCodec 获取编码器
func GetCodec(name string) (c Codec, err error) {
	if len(name) == 0 {
		name = JSON
	}
	fn, exist := codecs[name]

	if !exist {
		err = fmt.Errorf("codec not found: '%s'", name)
		return
	}

	c, err = fn()

	return
}

// RegisterCodec 注册编码器
func RegisterCodec(name string, fn NewCodecFunc) {
	if len(name) == 0 {
		panic("codec name is empty")
	}

	if fn == nil {
		panic("codec fn is nil")
	}

	c, err := fn()
	if err != nil {
		panic(err)
	}

	defaultCodecs[name] = c

	codecs[name] = fn
}

func Unmarshal(name string, body []byte, obj interface{}) error {
	c, err := GetCodec(name)
	if err != nil {
		return err
	}
	return c.Unmarshal(body, obj)
}

func Marshal(name string, obj interface{}) ([]byte, error) {
	c, err := GetCodec(name)
	if err != nil {
		return nil, err
	}
	return c.Marshal(obj)
}
