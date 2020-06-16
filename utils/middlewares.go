package utils

import (
	"runtime"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/go-trellis/config"
)

func LoadPprof(engine *gin.Engine, conf config.Config) {

	if conf == nil {
		return
	}

	if !conf.GetBoolean("enabled", false) {
		return
	}

	pprof.Register(engine)
	runtime.SetBlockProfileRate(int(conf.GetInt("block-profile-rate", 0)))
}

func LoadCors(engine *gin.Engine, conf config.Config) {

	var corsConf cors.Config
	if conf == nil {
		corsConf = cors.DefaultConfig()
		corsConf.AllowMethods = []string{"POST"}
		corsConf.AllowOrigins = []string{"*"}
		corsConf.AllowOriginFunc = func(origin string) bool {
			return true
		}
	} else {
		corsConf = cors.Config{
			AllowOrigins:     conf.GetStringList("allow-origins"),
			AllowMethods:     conf.GetStringList("allow-methods"),
			AllowHeaders:     conf.GetStringList("allow-headers"),
			ExposeHeaders:    conf.GetStringList("expose-headers"),
			AllowCredentials: conf.GetBoolean("allow-credentials", false),
			MaxAge:           conf.GetTimeDuration("max-age", time.Hour*12),
		}

		corsConf.AllowOriginFunc = wildcardMatchFunc(corsConf.AllowOrigins)
	}

	corsConf.AllowHeaders = append(corsConf.AllowHeaders, "X-Api", "X-Api-Timeout")

	engine.Use(cors.New(corsConf))
}

type wildcard struct {
	prefix string
	suffix string
}

func wildcardMatchFunc(allowedOrigins []string) func(string) bool {

	allowedWOrigins := []wildcard{}
	allowedOriginsAll := false

	for _, origin := range allowedOrigins {
		origin = strings.ToLower(origin)
		if origin == "*" {
			allowedOriginsAll = true
			allowedWOrigins = nil
			break
		} else if i := strings.IndexByte(origin, '*'); i >= 0 {
			w := wildcard{origin[0:i], origin[i+1:]}
			allowedWOrigins = append(allowedWOrigins, w)
		}
	}

	return func(origin string) bool {
		if allowedOriginsAll {
			return true
		}

		for _, w := range allowedWOrigins {
			if w.match(origin) {
				return true
			}
		}

		return false
	}
}

func (w wildcard) match(s string) bool {
	return len(s) >= len(w.prefix+w.suffix) && strings.HasPrefix(s, w.prefix) && strings.HasSuffix(s, w.suffix)
}
