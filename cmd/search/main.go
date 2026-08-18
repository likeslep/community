// Command search 启动 search-service：Elasticsearch 索引 + 搜索 + 事件同步。
package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"

	searchv1 "github.com/likeslep/community/api/gen/search/v1"
	"github.com/likeslep/community/internal/search/consumer"
	"github.com/likeslep/community/internal/search/handler"
	"github.com/likeslep/community/internal/search/service"
	"github.com/likeslep/community/pkg/config"
	"github.com/likeslep/community/pkg/grpcx"
	"github.com/likeslep/community/pkg/kafkax"
	"github.com/likeslep/community/pkg/logger"
	"github.com/likeslep/community/pkg/server"
)

func main() {
	var cfg struct {
		HTTPAddr     string        `env:"HTTP_ADDR" default:":8080"`
		GRPCAddr     string        `env:"GRPC_ADDR" default:":9090"`
		ESAddr       string        `env:"ELASTICSEARCH_ADDR" default:"http://localhost:9200"`
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
	lg = lg.Named("search-service")

	es, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{cfg.ESAddr}})
	if err != nil {
		lg.Fatal("连接 Elasticsearch 失败", zap.Error(err))
	}
	svc := service.NewSearchService(es)
	if err := svc.EnsureIndex(context.Background()); err != nil {
		lg.Warn("创建索引失败（可能已存在）", zap.Error(err))
	}
	h := handler.New(svc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	grpcSrv := grpcx.NewServer()
	searchv1.RegisterSearchServiceServer(grpcSrv.GRPC(), h)
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

	// 消费 content 事件同步索引。search 无独立 DB，幂等用内存实现（生产可换 Redis/DB）。
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	eventConsumer := consumer.New(svc, lg)
	kafkaConsumer := kafkax.NewConsumer(
		kafkax.ConsumerConfig{Brokers: brokers, Topic: "content.events", GroupID: "search-consumer"},
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
	lg.Info("search-service stopped")
}
