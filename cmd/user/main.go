// Command user 启动 user-service：注册/登录/资料 + user.created Outbox。
package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	userv1 "github.com/likeslep/community/api/gen/user/v1"
	"github.com/likeslep/community/internal/user/handler"
	usermigrations "github.com/likeslep/community/internal/user/migrations"
	"github.com/likeslep/community/internal/user/repository"
	"github.com/likeslep/community/internal/user/service"
	"github.com/likeslep/community/pkg/config"
	"github.com/likeslep/community/pkg/grpcx"
	"github.com/likeslep/community/pkg/kafkax"
	"github.com/likeslep/community/pkg/logger"
	"github.com/likeslep/community/pkg/migrate"
	"github.com/likeslep/community/pkg/outbox"
	"github.com/likeslep/community/pkg/server"
)

func main() {
	var cfg struct {
		HTTPAddr     string        `env:"HTTP_ADDR" default:":8080"`
		GRPCAddr     string        `env:"GRPC_ADDR" default:":9090"`
		DBDSN        string        `env:"DB_DSN" required:"true"`
		KafkaBrokers string        `env:"KAFKA_BROKERS" default:"localhost:9092"`
		JWTSecret    string        `env:"JWT_SECRET" required:"true"`
		TokenTTL     time.Duration `env:"TOKEN_TTL" default:"24h"`
		LogLevel     string        `env:"LOG_LEVEL" default:"info"`
		LogJSON      bool          `env:"LOG_JSON" default:"true"`
		ShutdownWait time.Duration `env:"SHUTDOWN_WAIT" default:"10s"`
	}
	if err := config.Load(&cfg); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	lg, err := logger.New(logger.Config{Level: cfg.LogLevel, JSON: cfg.LogJSON})
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	lg = lg.Named("user-service")

	// 连接 MySQL 并执行迁移。
	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		lg.Fatal("连接 MySQL 失败", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		lg.Fatal("获取底层 *sql.DB 失败", zap.Error(err))
	}
	m := migrate.New(sqlDB)
	if err := m.Load(usermigrations.FS, "."); err != nil {
		lg.Fatal("加载迁移脚本失败", zap.Error(err))
	}
	if err := m.Up(context.Background()); err != nil {
		lg.Fatal("执行迁移失败", zap.Error(err))
	}

	// 依赖装配。
	outboxStore := outbox.NewStore(db)
	kafkaProducer := kafkax.NewProducer(kafkax.ProducerConfig{
		Brokers: strings.Split(cfg.KafkaBrokers, ","),
		Topic:   "user.events",
	})
	repo := repository.NewGorm(db, outboxStore)
	svc := service.NewUserService(repo, service.Config{
		Producer:  "user-service",
		JWTSecret: []byte(cfg.JWTSecret),
		TokenTTL:  cfg.TokenTTL,
	})
	h := handler.New(svc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// gRPC 服务。
	grpcSrv := grpcx.NewServer()
	userv1.RegisterUserServiceServer(grpcSrv.GRPC(), h)
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

	// outbox 发布器（后台轮询 → Kafka）。
	publisher := outbox.NewPublisher(outboxStore, kafkaProducer, lg, outbox.PublisherConfig{
		Service:  "user-service",
		Interval: time.Second,
	})
	go func() { _ = publisher.Run(ctx) }()

	// HTTP 健康服务（阻塞直到收到信号）。
	httpSrv := server.New(server.Config{Addr: cfg.HTTPAddr, ShutdownWait: cfg.ShutdownWait}, lg)
	if err := httpSrv.Serve(ctx); err != nil {
		lg.Error("http server stopped with error", zap.Error(err))
	}

	// 优雅停机。
	grpcSrv.GRPC().GracefulStop()
	_ = kafkaProducer.Close()
	lg.Info("user-service stopped")
}
