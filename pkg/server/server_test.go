package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/likeslep/community/pkg/logger"
)

func TestHealthz(t *testing.T) {
	s := New(Config{Addr: ":0", ShutdownWait: 0}, logger.NewNop())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	s.Engine().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("healthz status = %d, want 200", w.Code)
	}
	if w.Header().Get("X-Request-ID") == "" {
		t.Fatal("healthz 应通过 RequestID 中间件注入 X-Request-ID")
	}
}
