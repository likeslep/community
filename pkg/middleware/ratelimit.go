package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/likeslep/community/pkg/redisx"
)

// RateLimiter 是基于 Redis 的滑动窗口限流器（plan.md §42）。
type RateLimiter struct {
	rdb    *redisx.Client
	limit  int
	window time.Duration
}

// NewRateLimiter 构造限流器。
func NewRateLimiter(rdb *redisx.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{rdb: rdb, limit: limit, window: window}
}

// Middleware 按客户端 IP 限流，超限返回 429。
func (r *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "ratelimit:" + c.ClientIP()
		count, err := r.rdb.Redis().Incr(c.Request.Context(), key).Result()
		if err != nil {
			c.Next() // Redis 故障时放行（降级）
			return
		}
		if count == 1 {
			_ = r.rdb.Redis().Expire(c.Request.Context(), key, r.window).Err()
		}
		if count > int64(r.limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"code": 429, "message": "请求过于频繁"})
			return
		}
		c.Next()
	}
}
