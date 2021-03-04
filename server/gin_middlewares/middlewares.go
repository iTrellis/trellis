package gin_middlewares

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	tnet "github.com/iTrellis/trellis/internal/net"
	"github.com/iTrellis/trellis/service"
)

type Handler struct {
	Name    string
	URLPath string
	Method  string
	Func    gin.HandlerFunc
}

var UseFuncs = make(map[string]gin.HandlerFunc)
var IndexGinFuncs []string

// RegistUseFuncs 注册
func RegistUseFuncs(name string, fn gin.HandlerFunc) error {
	_, ok := UseFuncs[name]
	if ok {
		return fmt.Errorf("use funcs (%s) is already exist", name)
	}
	UseFuncs[name] = fn
	IndexGinFuncs = append(IndexGinFuncs, name)
	return nil
}

func NewRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Header.Set("X-Request-ID", uuid.NewString())
		c.Next()
	}
}

type CSRFConfig struct {
	QueryAllowHosts   func() []string
	QueryAllowPattern func() []string
	Validator         func(c *gin.Context) bool
}

func NewCSRF(c CSRFConfig) gin.HandlerFunc {

	var validations []func(*url.URL) bool

	if c.QueryAllowHosts != nil {
		for _, r := range c.QueryAllowHosts() {
			validations = append(validations, tnet.MatchHostSuffix(r))
		}
	}

	if c.QueryAllowPattern != nil {
		for _, p := range c.QueryAllowPattern() {
			validations = append(validations, tnet.MatchPattern(regexp.MustCompile(p)))
		}
	}

	return func(ctx *gin.Context) {

		referer := ctx.Request.Header.Get(service.HeaderReferer)
		if referer == "" {
			ctx.AbortWithStatus(403)
			return
		}

		illegal := true
		if uri, err := url.Parse(referer); err == nil && uri.Host != "" {
			for _, validate := range validations {
				if validate(uri) {
					illegal = false
					break
				}
			}
		}
		if illegal {
			ctx.AbortWithStatus(403)
			return
		}

		// 添加隐藏csrf-token的认证
		if c.Validator != nil {
			if !c.Validator(ctx) {
				ctx.AbortWithStatus(403)
				return
			}
		}
	}
}
