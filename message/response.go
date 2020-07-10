package message

// Response 返回体
type Response struct {
	err    error
	Result interface{} `json:"result"`
}

// SetError 设置错误信息
func (p *Response) SetError(err error) {
	if err == nil {
		return
	}
	p.err = err
}

// SetBody 设置返回体
func (p *Response) SetBody(body interface{}) {
	p.Result = body
}

// GetError 获取错误信息
func (p *Response) GetError() error {
	return p.err
}
