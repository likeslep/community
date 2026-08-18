// Package consumer 是 content-service 的 Kafka 事件消费者。
package consumer

import (
	"context"
	"strconv"

	"go.uber.org/zap"

	"github.com/likeslep/community/internal/content/service"
	"github.com/likeslep/community/pkg/kafkax"
)

// EventConsumer 消费审核结果事件，驱动文章发布/驳回（plan.md §22.4）。
type EventConsumer struct {
	articles *service.ArticleService
	logger   *zap.Logger
}

// New 构造。
func New(articles *service.ArticleService, logger *zap.Logger) *EventConsumer {
	return &EventConsumer{articles: articles, logger: logger}
}

// Handle 处理一条事件。
func (c *EventConsumer) Handle(ctx context.Context, env kafkax.Envelope) error {
	if env.AggregateType != "article" {
		return nil // 仅处理文章审核结果
	}
	id, err := strconv.ParseUint(env.AggregateID, 10, 64)
	if err != nil {
		return err
	}
	switch env.EventType {
	case kafkax.EventModerationApproved:
		return c.articles.PublishArticle(ctx, id)
	case kafkax.EventModerationRejected:
		return c.articles.RejectArticle(ctx, id)
	default:
		return nil
	}
}
