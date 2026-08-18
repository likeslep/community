// Command admin 启动 admin-service：后台管理 + 审计日志（通过 gRPC 编排）。
package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	adminv1 "github.com/likeslep/community/api/gen/admin/v1"
	contentv1 "github.com/likeslep/community/api/gen/content/v1"
	interactionv1 "github.com/likeslep/community/api/gen/interaction/v1"
	moderationv1 "github.com/likeslep/community/api/gen/moderation/v1"
	userv1 "github.com/likeslep/community/api/gen/user/v1"
	"github.com/likeslep/community/internal/admin/handler"
	amigrations "github.com/likeslep/community/internal/admin/migrations"
	"github.com/likeslep/community/internal/admin/repository"
	"github.com/likeslep/community/internal/admin/service"
	"github.com/likeslep/community/pkg/config"
	"github.com/likeslep/community/pkg/grpcx"
	"github.com/likeslep/community/pkg/logger"
	"github.com/likeslep/community/pkg/migrate"
	"github.com/likeslep/community/pkg/server"
)

func main() {
	var cfg struct {
		HTTPAddr            string        `env:"HTTP_ADDR" default:":8080"`
		GRPCAddr            string        `env:"GRPC_ADDR" default:":9090"`
		DBDSN               string        `env:"DB_DSN" required:"true"`
		UserGRPCAddr        string        `env:"USER_GRPC_ADDR" default:"localhost:9091"`
		ContentGRPCAddr     string        `env:"CONTENT_GRPC_ADDR" default:"localhost:9092"`
		InteractionGRPCAddr string        `env:"INTERACTION_GRPC_ADDR" default:"localhost:9093"`
		ModerationGRPCAddr  string        `env:"MODERATION_GRPC_ADDR" default:"localhost:9094"`
		LogLevel            string        `env:"LOG_LEVEL" default:"info"`
		LogJSON             bool          `env:"LOG_JSON" default:"true"`
		ShutdownWait        time.Duration `env:"SHUTDOWN_WAIT" default:"10s"`
	}
	if err := config.Load(&cfg); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	lg, err := logger.New(logger.Config{Level: cfg.LogLevel, JSON: cfg.LogJSON})
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	lg = lg.Named("admin-service")

	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		lg.Fatal("连接 MySQL 失败", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		lg.Fatal("获取底层 *sql.DB 失败", zap.Error(err))
	}
	m := migrate.New(sqlDB)
	if err := m.Load(amigrations.FS, "."); err != nil {
		lg.Fatal("加载迁移脚本失败", zap.Error(err))
	}
	if err := m.Up(context.Background()); err != nil {
		lg.Fatal("执行迁移失败", zap.Error(err))
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
	interactionConn, err := grpcx.Dial(cfg.InteractionGRPCAddr)
	if err != nil {
		lg.Fatal("连接 interaction-service 失败", zap.Error(err))
	}
	defer interactionConn.Close()
	modConn, err := grpcx.Dial(cfg.ModerationGRPCAddr)
	if err != nil {
		lg.Fatal("连接 moderation-service 失败", zap.Error(err))
	}
	defer modConn.Close()

	repo := repository.NewGorm(db)
	svc := service.NewAdminService(
		userv1.NewUserServiceClient(userConn),
		contentv1.NewArticleServiceClient(contentConn),
		interactionv1.NewCommentServiceClient(interactionConn),
		moderationv1.NewReportServiceClient(modConn),
		moderationv1.NewSensitiveWordServiceClient(modConn),
		repo,
	)
	h := handler.New(svc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	grpcSrv := grpcx.NewServer()
	adminv1.RegisterAdminServiceServer(grpcSrv.GRPC(), h)
	grpcSrv.SetServing()
	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		lg.Fatal("监听 gRPC 失败", zap.Error(err))
	}
	go func() {
		lg.Info("grpc server starting", zap.String("addr", cfg.GRPCAddr))
		if err := grpcSrv.GRPC().Serve(lis); err != nil {
			lg.Error("grpc server stopped", zap.Error(err))
		}
	}()

	httpSrv := server.New(server.Config{Addr: cfg.HTTPAddr, ShutdownWait: cfg.ShutdownWait}, lg)
	if err := httpSrv.Serve(ctx); err != nil {
		lg.Error("http server stopped with error", zap.Error(err))
	}

	grpcSrv.GRPC().GracefulStop()
	lg.Info("admin-service stopped")
}
