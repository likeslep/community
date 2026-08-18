// Package middleware 提供 HTTP 中间件（plan.md §28：Request ID / Trace ID 传播）。
package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

// ctxKey 是未导出的上下文键类型，避免与其它包注入的键冲突。
type ctxKey string

const (
	ctxRequestID ctxKey = "request_id"
	ctxTraceID   ctxKey = "trace_id"
)

const (
	headerRequestID = "X-Request-ID"
	headerTraceID   = "X-Trace-ID"
)

// newID 生成 16 字节的十六进制随机 ID。
func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// RequestID 是 Gin 中间件：无 ID 时生成，有 ID 时透传，并注入 context 与响应头。
// 同一链路的多个服务共享相同的 trace_id，用于日志串联。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(headerRequestID)
		if rid == "" {
			rid = newID()
		}
		tid := c.GetHeader(headerTraceID)
		if tid == "" {
			tid = newID()
		}

		ctx := context.WithValue(c.Request.Context(), ctxRequestID, rid)
		ctx = context.WithValue(ctx, ctxTraceID, tid)
		c.Request = c.Request.WithContext(ctx)

		c.Header(headerRequestID, rid)
		c.Header(headerTraceID, tid)
		c.Next()
	}
}

// RequestIDFrom 从 context 读取 request id。
func RequestIDFrom(ctx context.Context) string {
	if v, ok := ctx.Value(ctxRequestID).(string); ok {
		return v
	}
	return ""
}

// TraceIDFrom 从 context 读取 trace id。
func TraceIDFrom(ctx context.Context) string {
	if v, ok := ctx.Value(ctxTraceID).(string); ok {
		return v
	}
	return ""
}
