package kafkax

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer 封装 kafka-go Writer，按分区键哈希保证同 aggregate 有序。
type Producer struct {
	w *kafka.Writer
}

// ProducerConfig 生产者配置。
type ProducerConfig struct {
	Brokers []string
	Topic   string
}

// NewProducer 构造生产者。
func NewProducer(cfg ProducerConfig) *Producer {
	return &Producer{w: &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}}
}

// Publish 序列化并发送事件。
func (p *Producer) Publish(ctx context.Context, env Envelope) error {
	b, err := json.Marshal(env)
	if err != nil {
		return err
	}
	return p.w.WriteMessages(ctx, kafka.Message{Key: []byte(env.PartitionKey()), Value: b})
}

// Close 关闭生产者。
func (p *Producer) Close() error { return p.w.Close() }
