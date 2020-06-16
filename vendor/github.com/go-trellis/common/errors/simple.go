// GNU GPL v3 License
// Copyright (c) 2020 go-trellis <hhh#rutcode.com>

package errors

import (
	"fmt"

	"github.com/google/uuid"
)

// SimpleError simple error functions
type SimpleError interface {
	ID() string
	Namespace() string
	Message() string
	Error() string
}

// Error error define
type Error struct {
	id        string
	namespace string
	message   string
}

// Newf 生成简单对象
func Newf(text string, params ...interface{}) SimpleError {
	return New(fmt.Sprintf(text, params...))
}

// New 生成简单对象
func New(text string) SimpleError {
	return new(defaultNamespace, uuid.New().String(), text)
}

func new(namespace, id, message string) *Error {
	return &Error{id: id, namespace: namespace, message: message}
}

func (p *Error) Error() string {
	return fmt.Sprintf("%s#%s:%s", p.namespace, p.id, p.message)
}

// ID 返回ID
func (p *Error) ID() string {
	return p.id
}

// Namespace 错误的域
func (p *Error) Namespace() string {
	return p.namespace
}

// Message 信息
func (p *Error) Message() string {
	return p.message
}
