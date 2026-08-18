// Command notification 启动 notification-service：消费事件生成站内通知。
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

	notificationv1 "github.com/likeslep/community/api/gen/notification/v1"
	"github.com/likeslep/community/internal/notification/handler"
	nmigrations "github.com/likeslep/community/internal/notification/migrations"
	"github.com/likeslep/community/internal/notification/repository"
	"github.com/likeslep/community/internal/notification/service"
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
	lg = lg.Named("notification-service")

	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		lg.Fatal("连接 MySQL 失败", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		lg.Fatal("获取底层 *sql.DB 失败", zap.Error(err))
	}
	m := migrate.New(sqlDB)
	if err := m.Load(nmigrations.FS, "."); err != nil {
		lg.Fatal("加载迁移脚本失败", zap.Error(err))
	}
	if err := m.Up(context.Background()); err != nil {
		lg.Fatal("执行迁移失败", zap.Error(err))
	}

	brokers := strings.Split(cfg.KafkaBrokers, ",")
	repo := repository.NewGorm(db)
	svc := service.NewNotificationService(repo)
	h := handler.New(svc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	grpcSrv := grpcx.NewServer()
	notificationv1.RegisterNotificationServiceServer(grpcSrv.GRPC(), h)
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

	// 消费多个领域事件 topic 生成通知。
	var consumers []*kafkax.Consumer
	for _, topic := range []string{"social.events", "content.events", "interaction.events"} {
		c := kafkax.NewConsumer(
			kafkax.ConsumerConfig{Brokers: brokers, Topic: topic, GroupID: "notification-consumer"},
			outbox.NewProcessedStore(db),
			kafkax.NewProducerRepublisher(brokers),
		)
		consumers = append(consumers, c)
		go func(topic string, c *kafkax.Consumer) {
			lg.Info("kafka consumer starting", zap.String("topic", topic))
			if err := c.Consume(ctx, svc.HandleEvent); err != nil {
				lg.Error("kafka consumer stopped", zap.String("topic", topic), zap.Error(err))
			}
		}(topic, c)
	}

	httpSrv := server.New(server.Config{Addr: cfg.HTTPAddr, ShutdownWait: cfg.ShutdownWait}, lg)
	if err := httpSrv.Serve(ctx); err != nil {
		lg.Error("http server stopped with error", zap.Error(err))
	}

	grpcSrv.GRPC().GracefulStop()
	for _, c := range consumers {
		_ = c.Close()
	}
	lg.Info("notification-service stopped")
}
