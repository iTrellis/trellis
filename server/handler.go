package server

import "github.com/gin-gonic/gin"

type Handler struct {
	Name    string
	URLPath string
	Method  string
	Func    gin.HandlerFunc
}
