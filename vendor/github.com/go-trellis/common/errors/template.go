// GNU GPL v3 License
// Copyright (c) 2017 go-trellis <hhh#rutcode.com>

package errors

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/go-trellis/common/encryption/hash"
)

const (
	errorcodeParseTemplate   = 1
	errorcodeExecuteTemplate = 2
)

// ErrorCodeTmpl error code template
type ErrorCodeTmpl struct {
	namespace string
	code      uint64
	message   string
}

// Params template params
type Params map[string]interface{}

var tmplDefined = make(map[string]bool)

// TN returns a new error code template
func TN(namespace string, code uint64, message string) *ErrorCodeTmpl {
	eKey := fmt.Sprintf("%s:%d", namespace, code)
	if tmplDefined[eKey] {
		panic(fmt.Errorf("error code (%s) is already exists", eKey))
	}
	tmplDefined[eKey] = true
	tmpl := &ErrorCodeTmpl{
		namespace: namespace,
		code:      code,
		message:   message,
	}
	return tmpl
}

// New ErrorCodeTmpl new error code by template
func (p *ErrorCodeTmpl) New(v ...Params) ErrorCode {
	params := Params{}
	if len(v) != 0 {
		for _, v := range v {
			for k, param := range v {
				params[k] = param
			}
		}
	}

	eCode := &errorCode{
		code:       p.code,
		stackTrace: callersDeepth(5),
		context:    make(map[string]interface{}),
	}

	errID := hash.NewCRCIEEE().Sum(fmt.Sprintf("%s.%d.%s.%d",
		p.namespace, p.code, p.message, time.Now().UnixNano()))

	t, e := template.New(genErrorCodeKey(p.namespace, p.code)).Parse(p.message)
	if e != nil {
		eCode.code = errorcodeParseTemplate

		eCode.err = new(defaultNamespace, errID, fmt.Sprintf(
			"parser template error, namespace: %s, code: %d, error: %s",
			p.namespace, p.code, e.Error()))
		return eCode
	}

	var buf bytes.Buffer
	if e := t.Execute(&buf, params); e != nil {
		eCode.code = errorcodeExecuteTemplate

		eCode.err = new(defaultNamespace, errID, fmt.Sprintf(
			"execute template error, namespace: %s code: %d, error: %s",
			p.message, p.code, e.Error()))
		return eCode
	}
	eCode.err = new(defaultNamespace, errID, buf.String())

	return eCode
}
