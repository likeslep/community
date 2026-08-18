// Package server 提供服务启动骨架：Gin 路由 + 健康检查 + 优雅停机（plan.md §6.1、§39）。
package server

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"

	"github.com/likeslep/community/pkg/metrics"
	"github.com/likeslep/community/pkg/middleware"
	"github.com/likeslep/community/pkg/version"
)

// Config HTTP 服务配置。
type Config struct {
	Addr         string        // 监听地址，如 ":8080"
	ShutdownWait time.Duration // 优雅停机等待在途请求的最大时长
}

// Server 是各服务的启动骨架，封装 gin.Engine 与生命周期管理。
type Server struct {
	cfg    Config
	logger *zap.Logger
	engine *gin.Engine
}

// New 构造服务骨架，注册 RequestID 中间件、Recovery 与健康检查端点。
func New(cfg Config, logger *zap.Logger) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(otelgin.Middleware(serviceName(logger)))
	engine.Use(middleware.RequestID())
	engine.Use(middleware.Metrics(serviceName(logger)))
	engine.Use(gin.Recovery())

	s := &Server{cfg: cfg, logger: logger, engine: engine}
	engine.GET("/healthz", s.healthz)
	engine.GET("/metrics", gin.WrapH(metrics.Handler()))
	return s
}

// serviceName 从 logger 名称提取服务名，未命名时回退为 unknown。
func serviceName(logger *zap.Logger) string {
	if name := logger.Name(); name != "" {
		return name
	}
	return "unknown"
}

// Engine 返回 gin.Engine，供各服务注册自己的业务路由。
func (s *Server) Engine() *gin.Engine { return s.engine }

func (s *Server) healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": version.Get(),
	})
}

// Serve 启动 HTTP 服务，直到 ctx 取消后执行优雅停机。
func (s *Server) Serve(ctx context.Context) error {
	srv := &http.Server{Addr: s.cfg.Addr, Handler: s.engine}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("http server starting", zap.String("addr", s.cfg.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownWait)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	s.logger.Info("http server stopped")
	return nil
}

// Run 以 SIGINT/SIGTERM 驱动 Serve，阻塞直到收到退出信号。
func (s *Server) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	return s.Serve(ctx)
}
