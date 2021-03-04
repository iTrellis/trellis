package gin_middlewares

import (
	"net/http"
	"runtime"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/iTrellis/config"
)

func LoadPprof(engine *gin.Engine, conf config.Config) {

	if conf == nil || engine == nil {
		return
	}

	if !conf.GetBoolean("enabled", false) {
		return
	}
	prefix := conf.GetString("prefix", "/")
	authorization := conf.GetString("authorization")
	if authorization != "" {
		authorGroup := engine.Group(prefix,
			func(c *gin.Context) {
				if c.Request.Header.Get("Authorization") != authorization {
					c.AbortWithStatus(http.StatusForbidden)
					return
				}
				c.Next()
			})
		pprof.RouteRegister(authorGroup)
	} else {
		pprof.Register(engine, prefix)
	}

	runtime.SetBlockProfileRate(int(conf.GetInt("block-profile-rate", 0)))
}
