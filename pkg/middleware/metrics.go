package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/likeslep/community/pkg/metrics"
)

// Metrics 是 Gin 中间件：记录 HTTP 请求计数、耗时与错误（plan.md §30）。
func Metrics(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}
		status := strconv.Itoa(c.Writer.Status())
		metrics.HTTPRequestsTotal.WithLabelValues(service, c.Request.Method, path, status).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(service, c.Request.Method, path).Observe(time.Since(start).Seconds())
		if c.Writer.Status() >= 400 {
			metrics.HTTPRequestErrors.WithLabelValues(service, c.Request.Method, path, status).Inc()
		}
	}
}
