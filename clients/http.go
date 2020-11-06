/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package clients

import (
	"context"
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

func (p *HTTPCaller) CallService(ctx context.Context, node *node.Node, msg *message.Message) (interface{}, error) {

	client := resty.New()

	resp, err := client.R().
		SetBody(msg.Payload).
		SetContext(ctx).
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
