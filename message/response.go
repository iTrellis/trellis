package message

import (
	"github.com/go-trellis/errors"
	"github.com/go-trellis/trellis/errcode"
)

// Response 返回体
type Response struct {
	errors.Error
	Result interface{} `json:"result"`
}

func (p *Response) SetError(err error) {
	if err == nil {
		return
	}
	var e errors.ErrorCode
	switch t := err.(type) {
	case errors.ErrorCode:
		e = t
	default:
		e = errcode.ErrTrellisResponse.New(errors.Params{"err": err.Error()})
	}
	p.ID = e.ID()
	p.Namespace = e.Namespace()
	p.Code = e.Code()
	p.Message = e.Error()
}

func (p *Response) SetBody(body interface{}) {
	p.Result = body
}
