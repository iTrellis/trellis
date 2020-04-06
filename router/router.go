package router

import (
	"github.com/go-trellis/trellis/message"
)

// HandlerFunc 函数执行
type HandlerFunc func(*message.Message) error

// Router 函数路由器
type Router interface {
	Route(*message.Message) HandlerFunc
}
