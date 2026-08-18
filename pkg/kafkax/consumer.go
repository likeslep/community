package kafkax

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/segmentio/kafka-go"
)

// Handler 处理一条事件。
type Handler func(ctx context.Context, env Envelope) error

// Republisher 将事件重新投递到 retry / DLQ topic。
type Republisher interface {
	Republish(ctx context.Context, env Envelope, topic string, retryCount int) error
}

// ConsumerConfig 消费者配置。
type ConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

// Consumer 封装 kafka-go Reader，提供幂等与重试/DLQ 骨架（plan.md §20、§21）。
type Consumer struct {
	r     *kafka.Reader
	store ProcessedStore
	repub Republisher
	topic string
}

// NewConsumer 构造消费者。
func NewConsumer(cfg ConsumerConfig, store ProcessedStore, repub Republisher) *Consumer {
	return &Consumer{
		r: kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.Brokers,
			Topic:   cfg.Topic,
			GroupID: cfg.GroupID,
		}),
		store: store,
		repub: repub,
		topic: cfg.Topic,
	}
}

// Close 关闭消费者。
func (c *Consumer) Close() error { return c.r.Close() }

// Consume 阻塞消费直到 ctx 取消。
func (c *Consumer) Consume(ctx context.Context, h Handler) error {
	for {
		msg, err := c.r.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		if err := c.handle(ctx, msg, h); err != nil {
			return err
		}
	}
}

// handle 处理单条消息：解析 → 幂等检查 → 调用 handler → 按错误分类提交/重试/DLQ。
func (c *Consumer) handle(ctx context.Context, msg kafka.Message, h Handler) error {
	var env Envelope
	if err := json.Unmarshal(msg.Value, &env); err != nil {
		// 无法解析的消息直接提交，避免毒消息阻塞消费。
		return c.r.CommitMessages(ctx, msg)
	}

	if c.store != nil && c.store.Exists(ctx, env.EventID) {
		return c.r.CommitMessages(ctx, msg)
	}

	count := retryCount(msg)
	switch NextAction(count, Classify(h(ctx, env))) {
	case ActionRetry:
		return c.repub.Republish(ctx, env, RetryTopic(c.topic), count+1)
	case ActionDLQ:
		return c.repub.Republish(ctx, env, DLQTopic(c.topic), count+1)
	default: // ActionCommit
		if c.store != nil {
			if err := c.store.Mark(ctx, env.EventID); err != nil {
				return err
			}
		}
		return c.r.CommitMessages(ctx, msg)
	}
}

func retryCount(msg kafka.Message) int {
	for _, h := range msg.Headers {
		if h.Key == "x-retry-count" {
			n, _ := strconv.Atoi(string(h.Value))
			return n
		}
	}
	return 0
}
