// Command feed 启动 feed-service：个性化信息流（Redis Sorted Set + 规则排序）。
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

	feedv1 "github.com/likeslep/community/api/gen/feed/v1"
	"github.com/likeslep/community/internal/feed/consumer"
	"github.com/likeslep/community/internal/feed/handler"
	"github.com/likeslep/community/internal/feed/service"
	"github.com/likeslep/community/pkg/config"
	"github.com/likeslep/community/pkg/grpcx"
	"github.com/likeslep/community/pkg/kafkax"
	"github.com/likeslep/community/pkg/logger"
	"github.com/likeslep/community/pkg/redisx"
	"github.com/likeslep/community/pkg/server"
)

func main() {
	var cfg struct {
		HTTPAddr     string        `env:"HTTP_ADDR" default:":8080"`
		GRPCAddr     string        `env:"GRPC_ADDR" default:":9090"`
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
	lg = lg.Named("feed-service")

	rdb := redisx.New(redisx.Config{Addr: cfg.RedisAddr})
	svc := service.NewFeedService(rdb, service.RuleRanker{})
	h := handler.New(svc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	grpcSrv := grpcx.NewServer()
	feedv1.RegisterFeedServiceServer(grpcSrv.GRPC(), h)
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

	// 消费 content 事件做 Fan-out on Write。followers 解析后续通过 social-service gRPC。
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	eventConsumer := consumer.New(svc, nil, lg)
	kafkaConsumer := kafkax.NewConsumer(
		kafkax.ConsumerConfig{Brokers: brokers, Topic: "content.events", GroupID: "feed-consumer"},
		kafkax.NewMemProcessedStore(),
		kafkax.NewProducerRepublisher(brokers),
	)
	go func() {
		lg.Info("kafka consumer starting", zap.String("topic", "content.events"))
		if err := kafkaConsumer.Consume(ctx, eventConsumer.Handle); err != nil {
			lg.Error("kafka consumer stopped", zap.Error(err))
		}
	}()

	httpSrv := server.New(server.Config{Addr: cfg.HTTPAddr, ShutdownWait: cfg.ShutdownWait}, lg)
	if err := httpSrv.Serve(ctx); err != nil {
		lg.Error("http server stopped with error", zap.Error(err))
	}

	grpcSrv.GRPC().GracefulStop()
	_ = kafkaConsumer.Close()
	_ = rdb.Close()
	lg.Info("feed-service stopped")
}
