// Command file 启动 file-service：文件上传/下载（LocalStorage）。
package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/likeslep/community/internal/file/handler"
	fmigrations "github.com/likeslep/community/internal/file/migrations"
	"github.com/likeslep/community/internal/file/repository"
	"github.com/likeslep/community/internal/file/service"
	"github.com/likeslep/community/internal/file/storage"
	"github.com/likeslep/community/pkg/config"
	"github.com/likeslep/community/pkg/logger"
	"github.com/likeslep/community/pkg/migrate"
	"github.com/likeslep/community/pkg/server"
)

func main() {
	var cfg struct {
		HTTPAddr     string        `env:"HTTP_ADDR" default:":8080"`
		DBDSN        string        `env:"DB_DSN" required:"true"`
		StorageDir   string        `env:"STORAGE_DIR" default:"./data/files"`
		MaxSize      int64         `env:"MAX_SIZE" default:"10485760"`
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
	lg = lg.Named("file-service")

	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		lg.Fatal("连接 MySQL 失败", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		lg.Fatal("获取底层 *sql.DB 失败", zap.Error(err))
	}
	m := migrate.New(sqlDB)
	if err := m.Load(fmigrations.FS, "."); err != nil {
		lg.Fatal("加载迁移脚本失败", zap.Error(err))
	}
	if err := m.Up(context.Background()); err != nil {
		lg.Fatal("执行迁移失败", zap.Error(err))
	}

	repo := repository.NewGorm(db)
	st := storage.NewLocalStorage(cfg.StorageDir)
	svc := service.NewFileService(repo, service.Config{MaxSize: cfg.MaxSize, Storage: st})
	h := handler.New(svc)

	srv := server.New(server.Config{Addr: cfg.HTTPAddr, ShutdownWait: cfg.ShutdownWait}, lg)
	h.RegisterRoutes(srv.Engine())
	if err := srv.Run(); err != nil {
		lg.Fatal("server stopped with error", zap.Error(err))
	}
}
