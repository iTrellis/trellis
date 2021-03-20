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
