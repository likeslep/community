// Package redisx 提供 Redis 客户端封装（plan.md §27）。
package redisx

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Config Redis 连接配置。
type Config struct {
	Addr     string `env:"REDIS_ADDR" default:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD" default:""`
	DB       int    `env:"REDIS_DB" default:"0"`
}

// Client 封装 go-redis Client，统一连接池与超时。
type Client struct {
	rdb *redis.Client
}

// New 构造 Redis 客户端。
func New(cfg Config) *Client {
	return &Client{rdb: redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})}
}

// Redis 返回底层 *redis.Client，供业务使用具体数据结构（String/Hash/Set/ZSet）。
func (c *Client) Redis() *redis.Client { return c.rdb }

// Ping 健康检查。
func (c *Client) Ping(ctx context.Context) error { return c.rdb.Ping(ctx).Err() }

// Close 关闭连接池。
func (c *Client) Close() error { return c.rdb.Close() }
