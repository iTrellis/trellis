// GNU GPL v3 License
// Copyright (c) 2017 go-trellis <hhh#rutcode.com>

package errors

import (
	"fmt"
	"strings"
)

const (
	defaultNamespace = "T:ERR"
)

// Errors error slice
type Errors []error

// NewErrors 生成错误数据对象
func NewErrors(errs ...error) Errors {
	var e Errors
	e = append(e, errs...)
	return e
}

func (p Errors) Error() string {
	return strings.Join(errorsString(p...), ";")
}

// Append 增补错误对象
func (p *Errors) Append(err ...error) {
	*p = append(*p, err...)
}

func errorsString(errs ...error) []string {
	var ss []string
	for _, e := range errs {
		switch ev := e.(type) {
		case ErrorCode:
			ss = append(ss, fmt.Sprintf("(%s#%d:%s) %s", ev.Namespace(), ev.Code(), ev.ID(), ev.Error()))
		case SimpleError:
			ss = append(ss, fmt.Sprintf("(%s:%s) %s", ev.Namespace(), ev.ID(), ev.Error()))
		default:
			ss = append(ss, ev.Error())
		}
	}
	return ss
}
