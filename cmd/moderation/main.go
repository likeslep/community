// Command moderation 启动 moderation-service：敏感词检测 + 审核任务 + 举报 + Outbox。
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

	moderationv1 "github.com/likeslep/community/api/gen/moderation/v1"
	"github.com/likeslep/community/internal/moderation/consumer"
	"github.com/likeslep/community/internal/moderation/handler"
	mmigrations "github.com/likeslep/community/internal/moderation/migrations"
	"github.com/likeslep/community/internal/moderation/repository"
	"github.com/likeslep/community/internal/moderation/service"
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
	lg = lg.Named("moderation-service")

	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		lg.Fatal("连接 MySQL 失败", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		lg.Fatal("获取底层 *sql.DB 失败", zap.Error(err))
	}
	m := migrate.New(sqlDB)
	if err := m.Load(mmigrations.FS, "."); err != nil {
		lg.Fatal("加载迁移脚本失败", zap.Error(err))
	}
	if err := m.Up(context.Background()); err != nil {
		lg.Fatal("执行迁移失败", zap.Error(err))
	}

	brokers := strings.Split(cfg.KafkaBrokers, ",")
	outboxStore := outbox.NewStore(db)
	producer := kafkax.NewProducer(kafkax.ProducerConfig{Brokers: brokers, Topic: "moderation.events"})
	repo := repository.NewGorm(db, outboxStore)
	svc := service.NewModerationService(repo, service.Config{Producer: "moderation-service"})
	h := handler.New(svc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// gRPC 服务。
	grpcSrv := grpcx.NewServer()
	moderationv1.RegisterSensitiveWordServiceServer(grpcSrv.GRPC(), h)
	moderationv1.RegisterModerationServiceServer(grpcSrv.GRPC(), h)
	moderationv1.RegisterReportServiceServer(grpcSrv.GRPC(), h)
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

	// outbox 发布器。
	publisher := outbox.NewPublisher(outboxStore, producer, lg, outbox.PublisherConfig{
		Service: "moderation-service", Interval: time.Second,
	})
	go func() { _ = publisher.Run(ctx) }()

	// Kafka 消费者：消费 content 提交事件生成审核任务。
	eventConsumer := consumer.New(svc, lg)
	kafkaConsumer := kafkax.NewConsumer(
		kafkax.ConsumerConfig{Brokers: brokers, Topic: "content.events", GroupID: "moderation-consumer"},
		outbox.NewProcessedStore(db),
		kafkax.NewProducerRepublisher(brokers),
	)
	go func() {
		lg.Info("kafka consumer starting", zap.String("topic", "content.events"))
		if err := kafkaConsumer.Consume(ctx, eventConsumer.Handle); err != nil {
			lg.Error("kafka consumer stopped", zap.Error(err))
		}
	}()

	// HTTP 健康服务。
	httpSrv := server.New(server.Config{Addr: cfg.HTTPAddr, ShutdownWait: cfg.ShutdownWait}, lg)
	if err := httpSrv.Serve(ctx); err != nil {
		lg.Error("http server stopped with error", zap.Error(err))
	}

	grpcSrv.GRPC().GracefulStop()
	_ = kafkaConsumer.Close()
	_ = producer.Close()
	lg.Info("moderation-service stopped")
}
