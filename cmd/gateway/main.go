// Command gateway 启动 gateway-service：REST 路由 + JWT 认证 + gRPC 转发。
package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	adminv1 "github.com/likeslep/community/api/gen/admin/v1"
	contentv1 "github.com/likeslep/community/api/gen/content/v1"
	feedv1 "github.com/likeslep/community/api/gen/feed/v1"
	interactionv1 "github.com/likeslep/community/api/gen/interaction/v1"
	notificationv1 "github.com/likeslep/community/api/gen/notification/v1"
	searchv1 "github.com/likeslep/community/api/gen/search/v1"
	socialv1 "github.com/likeslep/community/api/gen/social/v1"
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
		HTTPAddr             string        `env:"HTTP_ADDR" default:":8080"`
		UserGRPCAddr         string        `env:"USER_GRPC_ADDR" default:"localhost:9090"`
		ContentGRPCAddr      string        `env:"CONTENT_GRPC_ADDR" default:"localhost:9091"`
		AdminGRPCAddr        string        `env:"ADMIN_GRPC_ADDR" default:"localhost:9092"`
		InteractionGRPCAddr  string        `env:"INTERACTION_GRPC_ADDR" default:"localhost:9093"`
		SocialGRPCAddr       string        `env:"SOCIAL_GRPC_ADDR" default:"localhost:9094"`
		SearchGRPCAddr       string        `env:"SEARCH_GRPC_ADDR" default:"localhost:9095"`
		FeedGRPCAddr         string        `env:"FEED_GRPC_ADDR" default:"localhost:9096"`
		NotificationGRPCAddr string        `env:"NOTIFICATION_GRPC_ADDR" default:"localhost:9097"`
		JWTSecret            string        `env:"JWT_SECRET" required:"true"`
		TraceEndpoint        string        `env:"TRACE_ENDPOINT" default:""`
		LogLevel             string        `env:"LOG_LEVEL" default:"info"`
		LogJSON              bool          `env:"LOG_JSON" default:"true"`
		ShutdownWait         time.Duration `env:"SHUTDOWN_WAIT" default:"10s"`
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

	userConn := mustDial(lg, "user-service", cfg.UserGRPCAddr)
	defer userConn.Close()
	contentConn := mustDial(lg, "content-service", cfg.ContentGRPCAddr)
	defer contentConn.Close()
	adminConn := mustDial(lg, "admin-service", cfg.AdminGRPCAddr)
	defer adminConn.Close()
	interactionConn := mustDial(lg, "interaction-service", cfg.InteractionGRPCAddr)
	defer interactionConn.Close()
	socialConn := mustDial(lg, "social-service", cfg.SocialGRPCAddr)
	defer socialConn.Close()
	searchConn := mustDial(lg, "search-service", cfg.SearchGRPCAddr)
	defer searchConn.Close()
	feedConn := mustDial(lg, "feed-service", cfg.FeedGRPCAddr)
	defer feedConn.Close()
	notificationConn := mustDial(lg, "notification-service", cfg.NotificationGRPCAddr)
	defer notificationConn.Close()

	h := handler.New(handler.Clients{
		Users:         userv1.NewUserServiceClient(userConn),
		Articles:      contentv1.NewArticleServiceClient(contentConn),
		Questions:     contentv1.NewQuestionServiceClient(contentConn),
		Admin:         adminv1.NewAdminServiceClient(adminConn),
		Interactions:  interactionv1.NewInteractionServiceClient(interactionConn),
		Comments:      interactionv1.NewCommentServiceClient(interactionConn),
		Social:        socialv1.NewSocialServiceClient(socialConn),
		Search:        searchv1.NewSearchServiceClient(searchConn),
		Feed:          feedv1.NewFeedServiceClient(feedConn),
		Notifications: notificationv1.NewNotificationServiceClient(notificationConn),
	})

	srv := server.New(server.Config{Addr: cfg.HTTPAddr, ShutdownWait: cfg.ShutdownWait}, lg)
	h.RegisterRoutes(srv.Engine(), gwmw.Auth([]byte(cfg.JWTSecret)))

	if err := srv.Run(); err != nil {
		lg.Fatal("server stopped with error", zap.Error(err))
	}
}

// mustDial 建立 gRPC 连接，失败时退出。
func mustDial(lg *zap.Logger, name, addr string) *grpc.ClientConn {
	conn, err := grpcx.Dial(addr)
	if err != nil {
		lg.Fatal("连接 "+name+" 失败", zap.String("addr", addr), zap.Error(err))
	}
	return conn
}
