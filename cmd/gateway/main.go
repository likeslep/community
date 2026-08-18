// Command gateway 启动 gateway-service：REST 路由 + JWT 认证 + gRPC 转发。
package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/zap"

	contentv1 "github.com/likeslep/community/api/gen/content/v1"
	userv1 "github.com/likeslep/community/api/gen/user/v1"
	"github.com/likeslep/community/internal/gateway/handler"
	gwmw "github.com/likeslep/community/internal/gateway/middleware"
	"github.com/likeslep/community/pkg/config"
	"github.com/likeslep/community/pkg/grpcx"
	"github.com/likeslep/community/pkg/logger"
	"github.com/likeslep/community/pkg/server"
	"github.com/likeslep/community/pkg/tracing"
)

func main() {
	var cfg struct {
		HTTPAddr        string        `env:"HTTP_ADDR" default:":8080"`
		UserGRPCAddr    string        `env:"USER_GRPC_ADDR" default:"localhost:9090"`
		ContentGRPCAddr string        `env:"CONTENT_GRPC_ADDR" default:"localhost:9091"`
		JWTSecret       string        `env:"JWT_SECRET" required:"true"`
		TraceEndpoint   string        `env:"TRACE_ENDPOINT" default:""`
		LogLevel        string        `env:"LOG_LEVEL" default:"info"`
		LogJSON         bool          `env:"LOG_JSON" default:"true"`
		ShutdownWait    time.Duration `env:"SHUTDOWN_WAIT" default:"10s"`
	}
	if err := config.Load(&cfg); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	lg, err := logger.New(logger.Config{Level: cfg.LogLevel, JSON: cfg.LogJSON})
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	lg = lg.Named("gateway-service")

	if cfg.TraceEndpoint != "" {
		shutdown, err := tracing.Setup(context.Background(), "gateway-service", cfg.TraceEndpoint)
		if err != nil {
			lg.Warn("初始化 tracing 失败", zap.Error(err))
		} else {
			defer func() { _ = shutdown(context.Background()) }()
		}
	}

	userConn, err := grpcx.Dial(cfg.UserGRPCAddr)
	if err != nil {
		lg.Fatal("连接 user-service 失败", zap.Error(err))
	}
	defer userConn.Close()

	contentConn, err := grpcx.Dial(cfg.ContentGRPCAddr)
	if err != nil {
		lg.Fatal("连接 content-service 失败", zap.Error(err))
	}
	defer contentConn.Close()

	h := handler.New(
		userv1.NewUserServiceClient(userConn),
		contentv1.NewArticleServiceClient(contentConn),
		contentv1.NewQuestionServiceClient(contentConn),
	)
	srv := server.New(server.Config{Addr: cfg.HTTPAddr, ShutdownWait: cfg.ShutdownWait}, lg)
	h.RegisterRoutes(srv.Engine(), gwmw.Auth([]byte(cfg.JWTSecret)))

	if err := srv.Run(); err != nil {
		lg.Fatal("server stopped with error", zap.Error(err))
	}
}
