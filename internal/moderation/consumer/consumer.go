// Package consumer 是 moderation-service 的 Kafka 事件消费者。
package consumer

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/likeslep/community/internal/moderation/service"
	"github.com/likeslep/community/pkg/kafkax"
)

// EventConsumer 消费 content 提交事件，生成审核任务（plan.md §22.3）。
type EventConsumer struct {
	svc    *service.ModerationService
	logger *zap.Logger
}

// New 构造。
func New(svc *service.ModerationService, logger *zap.Logger) *EventConsumer {
	return &EventConsumer{svc: svc, logger: logger}
}

// Handle 处理一条事件。
func (c *EventConsumer) Handle(ctx context.Context, env kafkax.Envelope) error {
	if env.EventType != kafkax.EventArticleSubmitted {
		return nil // 忽略其它事件类型
	}

	var payload struct {
		Title string `json:"title"`
	}
	if err := json.Unmarshal(env.Payload, &payload); err != nil {
		return err
	}
	return c.svc.CreateTask(ctx, env.AggregateType, env.AggregateID, payload.Title)
}
