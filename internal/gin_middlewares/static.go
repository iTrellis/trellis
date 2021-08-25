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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iTrellis/common/logger"
	"github.com/iTrellis/trellis/service"
)

// StatFunc log request & response
// TODO added prometheus metrics
func StatFunc(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.Request.Header.Get(service.HeaderXRequestID)
		api := c.Request.Header.Get("X-API")
		reqTime := time.Now().UnixNano() / int64(time.Microsecond)

		logger.Info("msg", "audit", "position", "request", "request_id", reqID, "url", c.Request.URL.String(),
			"method", c.Request.Method, "X-API", api, "request_time", reqTime)

		c.Next()

		respTime := time.Now().UnixNano() / int64(time.Microsecond)

		logger.Info("msg", "audit", "position", "response", "request_id", reqID, "url", c.Request.URL.String(),
			"method", c.Request.Method, "X-API", api,
			"request_time", reqTime, "response_time", respTime, "cost_time", respTime-reqTime)
	}
}
