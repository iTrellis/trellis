package clients

import (
	"fmt"
	"net/http"

	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/message/proto"

	"github.com/go-resty/resty/v2"
	"github.com/go-trellis/node"
)

// TODO
// dial options

func init() {
	RegistCaller(proto.Protocol_HTTP, NewHTTPCaller())
}

type HTTPCaller struct{}

func NewHTTPCaller() Caller {
	return &HTTPCaller{}
}

func (p *HTTPCaller) CallService(node *node.Node, msg *message.Message) (interface{}, error) {

	client := resty.New()

	resp, err := client.R().
		SetBody(msg.Payload).
		SetHeader("Content-Type", msg.GetHeader("Content-Type")).
		Post(node.Value)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("call server failed")
	}

	return resp.Body(), nil
}
