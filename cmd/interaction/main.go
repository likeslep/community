// Command interaction 启动 interaction-service：点赞/收藏/评论/浏览 + Outbox。
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

	interactionv1 "github.com/likeslep/community/api/gen/interaction/v1"
	"github.com/likeslep/community/internal/interaction/handler"
	imigrations "github.com/likeslep/community/internal/interaction/migrations"
	"github.com/likeslep/community/internal/interaction/repository"
	"github.com/likeslep/community/internal/interaction/service"
	"github.com/likeslep/community/pkg/config"
	"github.com/likeslep/community/pkg/grpcx"
	"github.com/likeslep/community/pkg/kafkax"
	"github.com/likeslep/community/pkg/logger"
	"github.com/likeslep/community/pkg/migrate"
	"github.com/likeslep/community/pkg/outbox"
	"github.com/likeslep/community/pkg/redisx"
	"github.com/likeslep/community/pkg/server"
)

func main() {
	var cfg struct {
		HTTPAddr     string        `env:"HTTP_ADDR" default:":8080"`
		GRPCAddr     string        `env:"GRPC_ADDR" default:":9090"`
		DBDSN        string        `env:"DB_DSN" required:"true"`
		RedisAddr    string        `env:"REDIS_ADDR" default:"localhost:6379"`
		KafkaBrokers string        `env:"KAFKA_BROKERS" default:"localhost:9092"`
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
	lg = lg.Named("interaction-service")

	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		lg.Fatal("连接 MySQL 失败", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		lg.Fatal("获取底层 *sql.DB 失败", zap.Error(err))
	}
	m := migrate.New(sqlDB)
	if err := m.Load(imigrations.FS, "."); err != nil {
		lg.Fatal("加载迁移脚本失败", zap.Error(err))
	}
	if err := m.Up(context.Background()); err != nil {
		lg.Fatal("执行迁移失败", zap.Error(err))
	}

	brokers := strings.Split(cfg.KafkaBrokers, ",")
	rdb := redisx.New(redisx.Config{Addr: cfg.RedisAddr})
	outboxStore := outbox.NewStore(db)
	producer := kafkax.NewProducer(kafkax.ProducerConfig{Brokers: brokers, Topic: "interaction.events"})
	repo := repository.NewGorm(db, outboxStore)
	svc := service.NewInteractionService(repo, service.Config{Producer: "interaction-service", Redis: rdb})
	h := handler.New(svc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	grpcSrv := grpcx.NewServer()
	interactionv1.RegisterInteractionServiceServer(grpcSrv.GRPC(), h)
	interactionv1.RegisterCommentServiceServer(grpcSrv.GRPC(), h)
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

	publisher := outbox.NewPublisher(outboxStore, producer, lg, outbox.PublisherConfig{
		Service: "interaction-service", Interval: time.Second,
	})
	go func() { _ = publisher.Run(ctx) }()

	httpSrv := server.New(server.Config{Addr: cfg.HTTPAddr, ShutdownWait: cfg.ShutdownWait}, lg)
	if err := httpSrv.Serve(ctx); err != nil {
		lg.Error("http server stopped with error", zap.Error(err))
	}

	grpcSrv.GRPC().GracefulStop()
	_ = producer.Close()
	_ = rdb.Close()
	lg.Info("interaction-service stopped")
}
