package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		reqID, traceID string
		wantGen        bool
	}{
		{"透传已有", "req-123", "trace-456", false},
		{"生成缺失", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(RequestID())
			r.GET("/", func(c *gin.Context) {
				c.String(http.StatusOK, RequestIDFrom(c.Request.Context())+"|"+TraceIDFrom(c.Request.Context()))
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.reqID != "" {
				req.Header.Set(headerRequestID, tt.reqID)
			}
			if tt.traceID != "" {
				req.Header.Set(headerTraceID, tt.traceID)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if tt.reqID != "" && w.Header().Get(headerRequestID) != tt.reqID {
				t.Fatalf("request id 未透传: got %q want %q", w.Header().Get(headerRequestID), tt.reqID)
			}
			if tt.traceID != "" && w.Header().Get(headerTraceID) != tt.traceID {
				t.Fatalf("trace id 未透传: got %q want %q", w.Header().Get(headerTraceID), tt.traceID)
			}
			if tt.wantGen && (w.Header().Get(headerRequestID) == "" || w.Header().Get(headerTraceID) == "") {
				t.Fatal("期望生成 request/trace id，实际为空")
			}
		})
	}
}
