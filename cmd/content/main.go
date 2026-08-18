// Command content 启动 content-service：文章创建/编辑/提交审核 + Outbox。
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

	contentv1 "github.com/likeslep/community/api/gen/content/v1"
	"github.com/likeslep/community/internal/content/consumer"
	"github.com/likeslep/community/internal/content/handler"
	cmigrations "github.com/likeslep/community/internal/content/migrations"
	"github.com/likeslep/community/internal/content/repository"
	"github.com/likeslep/community/internal/content/service"
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
	lg = lg.Named("content-service")

	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		lg.Fatal("连接 MySQL 失败", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		lg.Fatal("获取底层 *sql.DB 失败", zap.Error(err))
	}
	m := migrate.New(sqlDB)
	if err := m.Load(cmigrations.FS, "."); err != nil {
		lg.Fatal("加载迁移脚本失败", zap.Error(err))
	}
	if err := m.Up(context.Background()); err != nil {
		lg.Fatal("执行迁移失败", zap.Error(err))
	}

	brokers := strings.Split(cfg.KafkaBrokers, ",")
	outboxStore := outbox.NewStore(db)
	kafkaProducer := kafkax.NewProducer(kafkax.ProducerConfig{
		Brokers: brokers,
		Topic:   "content.events",
	})
	repo := repository.NewGorm(db, outboxStore)
	articleSvc := service.NewArticleService(repo, service.Config{Producer: "content-service"})
	articleH := handler.NewArticleHandler(articleSvc)
	qaSvc := service.NewQAService(repo, service.Config{Producer: "content-service"})
	qaH := handler.NewQuestionHandler(qaSvc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	grpcSrv := grpcx.NewServer()
	contentv1.RegisterArticleServiceServer(grpcSrv.GRPC(), articleH)
	contentv1.RegisterQuestionServiceServer(grpcSrv.GRPC(), qaH)
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

	publisher := outbox.NewPublisher(outboxStore, kafkaProducer, lg, outbox.PublisherConfig{
		Service:  "content-service",
		Interval: time.Second,
	})
	go func() { _ = publisher.Run(ctx) }()

	// Kafka 消费者：消费审核结果驱动发布/驳回。
	eventConsumer := consumer.New(articleSvc, lg)
	kafkaConsumer := kafkax.NewConsumer(
		kafkax.ConsumerConfig{Brokers: brokers, Topic: "moderation.events", GroupID: "content-consumer"},
		outbox.NewProcessedStore(db),
		kafkax.NewProducerRepublisher(brokers),
	)
	go func() {
		lg.Info("kafka consumer starting", zap.String("topic", "moderation.events"))
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
	_ = kafkaProducer.Close()
	lg.Info("content-service stopped")
}
