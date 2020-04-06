package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-trellis/trellis/message"
	"github.com/go-trellis/trellis/registry"
)

var hClient = &http.Client{}

type HTTPClient struct{}

func (*HTTPClient) Call(req *message.Request) (*message.Response, error) {
	s, err := registry.GetService(req.GetMessage().GetServer())
	if err != nil {
		return nil, err
	}

	n, ok := s.GetNodes().NodeFor(req.ID())
	if !ok {
		return nil, fmt.Errorf("not found service nodes")
	}

	reqBody, err := req.Read()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(req.Method, n.Value, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	response, err := hClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status not ok, %d", response.StatusCode)
	}

	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	resp := &message.Response{}
	resp.SetBody(bs)
	return resp, nil
}

func (*HTTPClient) String() string {
	return "http"
}

func NewHTTPClient() Client {
	return (*HTTPClient)(nil)
}
