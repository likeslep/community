// Package logger 提供基于 zap 的统一结构化日志构造（plan.md §28）。
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 日志配置。
type Config struct {
	Level string // debug / info / warn / error
	JSON  bool   // 是否输出 JSON 格式
}

// New 构造生产级 Logger。Level 非法时返回错误。
func New(cfg Config) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var enc zapcore.Encoder
	if cfg.JSON {
		enc = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		enc = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(enc, zapcore.Lock(os.Stdout), level)
	return zap.New(core, zap.AddCaller()), nil
}

// NewNop 返回一个丢弃所有输出的 Logger，用于测试。
func NewNop() *zap.Logger { return zap.NewNop() }
