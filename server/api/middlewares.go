package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iTrellis/common/logger"
	"github.com/iTrellis/trellis/service"
)

func StatFunc(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqTime := time.Now().UnixNano() / int64(time.Microsecond)

		c.Next()

		respTime := time.Now().UnixNano() / int64(time.Microsecond)

		reqID := c.Request.Header.Get(service.HeaderXRequestID)

		api := c.Request.Header.Get("X-API")
		logger.Info("msg", "api_cost(us)", "request_id", reqID, "url", c.Request.URL.String(),
			"method", c.Request.Method, "X-API", api,
			"request_time", reqTime, "response_time", respTime, "cost_time", respTime-reqTime)
	}
}
