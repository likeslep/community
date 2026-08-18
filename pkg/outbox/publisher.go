package outbox

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/likeslep/community/pkg/kafkax"
)

// Publisher 轮询 outbox 并发布到 Kafka。
type Publisher struct {
	store    *Store
	producer *kafkax.Producer
	logger   *zap.Logger
	svcName  string
	interval time.Duration
	batch    int
}

// PublisherConfig 发布器配置。
type PublisherConfig struct {
	Service  string        // 服务名，用于事件 producer 字段
	Interval time.Duration // 轮询间隔
	Batch    int           // 每批最多发布条数
}

// NewPublisher 构造发布器。
func NewPublisher(store *Store, producer *kafkax.Producer, logger *zap.Logger, cfg PublisherConfig) *Publisher {
	if cfg.Interval <= 0 {
		cfg.Interval = time.Second
	}
	if cfg.Batch <= 0 {
		cfg.Batch = 100
	}
	return &Publisher{
		store:    store,
		producer: producer,
		logger:   logger,
		svcName:  cfg.Service,
		interval: cfg.Interval,
		batch:    cfg.Batch,
	}
}

// Run 阻塞轮询直到 ctx 取消。
func (p *Publisher) Run(ctx context.Context) error {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := p.publishBatch(ctx); err != nil {
				p.logger.Error("outbox 发布失败", zap.Error(err))
			}
		}
	}
}

func (p *Publisher) publishBatch(ctx context.Context) error {
	events, err := p.store.Pending(ctx, p.batch)
	if err != nil {
		return err
	}
	for _, e := range events {
		env := p.toEnvelope(e)
		if err := p.producer.Publish(ctx, env); err != nil {
			return err
		}
		if err := p.store.MarkPublished(ctx, e.ID); err != nil {
			return err
		}
	}
	return nil
}

func (p *Publisher) toEnvelope(e Event) kafkax.Envelope {
	return kafkax.Envelope{
		EventID:          e.ID,
		EventType:        e.EventType,
		Version:          1,
		OccurredAt:       e.CreatedAt,
		Producer:         p.svcName,
		AggregateType:    e.AggregateType,
		AggregateID:      e.AggregateID,
		AggregateVersion: e.AggregateVersion,
		TraceID:          e.TraceID,
		Payload:          json.RawMessage(e.Payload),
	}
}
