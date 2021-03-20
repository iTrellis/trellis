/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package gin_middlewares

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iTrellis/config"
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

func LoadGZip(conf config.Config) gin.HandlerFunc {

	if conf == nil || !conf.GetBoolean("enabled", true) {
		return nil
	}

	compressLevel := conf.GetString("level", "default")

	level := gzip.DefaultCompression

	switch compressLevel {
	case "best-compression":
		level = gzip.BestCompression
	case "best-speed":
		level = gzip.BestSpeed
	}

	return gzip.Gzip(level)
}
