package gin_middlewares

import (
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

	pprof.Register(engine)
	runtime.SetBlockProfileRate(int(conf.GetInt("block-profile-rate", 0)))
}
