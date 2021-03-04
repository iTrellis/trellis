package api

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

		logger.Info("msg", "request_info", "request_id", reqID, "url", c.Request.URL.String(),
			"method", c.Request.Method, "X-API", api, "request_time", reqTime)

		c.Next()

		respTime := time.Now().UnixNano() / int64(time.Microsecond)

		logger.Info("msg", "api_cost(us)", "request_id", reqID, "url", c.Request.URL.String(),
			"method", c.Request.Method, "X-API", api,
			"request_time", reqTime, "response_time", respTime, "cost_time", respTime-reqTime)
	}
}
