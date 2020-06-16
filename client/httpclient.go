package client

// import (
// 	"bytes"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"

// 	"github.com/go-trellis/trellis/message"
// )

// var hClient = &http.Client{}

// type HTTPClient struct{}

// func (*HTTPClient) Call(req *message.Request) (*message.Response, error) {

// 	reqBody, err := req.Read()
// 	if err != nil {
// 		return nil, err
// 	}

// 	request, err := http.NewRequest(req.Method, req.get, bytes.NewReader(reqBody))
// 	if err != nil {
// 		return nil, err
// 	}

// 	response, err := hClient.Do(request)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer response.Body.Close()

// 	if response.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("status not ok, %d", response.StatusCode)
// 	}

// 	bs, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp := &message.Response{}
// 	resp.SetBody(bs)
// 	return resp, nil
// }

// func (*HTTPClient) String() string {
// 	return "http"
// }

// func NewHTTPClient() Client {
// 	return (*HTTPClient)(nil)
// }
