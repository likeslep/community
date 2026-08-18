package kafkax

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	"github.com/segmentio/kafka-go"
)

// ProducerRepublisher 将事件投递到任意 topic（用于 retry / DLQ）。
// 并发安全：用 mutex 保护 writers map，避免并发写 map 竞态。
type ProducerRepublisher struct {
	brokers []string
	mu      sync.Mutex
	writers map[string]*kafka.Writer
}

// NewProducerRepublisher 构造投递器。
func NewProducerRepublisher(brokers []string) *ProducerRepublisher {
	return &ProducerRepublisher{brokers: brokers, writers: make(map[string]*kafka.Writer)}
}

// Republish 将事件投递到目标 topic，并写入重试次数 header。
func (p *ProducerRepublisher) Republish(ctx context.Context, env Envelope, topic string, retryCount int) error {
	w := p.writer(topic)
	b, err := json.Marshal(env)
	if err != nil {
		return err
	}
	return w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(env.PartitionKey()),
		Value: b,
		Headers: []kafka.Header{
			{Key: "x-retry-count", Value: []byte(strconv.Itoa(retryCount))},
		},
	})
}

func (p *ProducerRepublisher) writer(topic string) *kafka.Writer {
	p.mu.Lock()
	defer p.mu.Unlock()
	if w, ok := p.writers[topic]; ok {
		return w
	}
	w := &kafka.Writer{
		Addr:     kafka.TCP(p.brokers...),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}
	p.writers[topic] = w
	return w
}
