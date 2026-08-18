// Package consumer 是 search-service 的 Kafka 事件消费者。
package consumer

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/likeslep/community/internal/search/service"
	"github.com/likeslep/community/pkg/kafkax"
)

// EventConsumer 消费 content 事件同步搜索索引（plan.md §22.5）。
type EventConsumer struct {
	search *service.SearchService
	logger *zap.Logger
}

// New 构造。
func New(search *service.SearchService, logger *zap.Logger) *EventConsumer {
	return &EventConsumer{search: search, logger: logger}
}

// Handle 处理一条事件。
func (c *EventConsumer) Handle(ctx context.Context, env kafkax.Envelope) error {
	var payload struct {
		ArticleID  uint64 `json:"article_id"`
		QuestionID uint64 `json:"question_id"`
		Title      string `json:"title"`
	}

	switch env.EventType {
	case kafkax.EventArticlePublished:
		_ = json.Unmarshal(env.Payload, &payload)
		return c.search.Index(ctx, service.Document{
			ID: env.AggregateID, Type: "article", Title: payload.Title,
		})
	case kafkax.EventQuestionPublished:
		_ = json.Unmarshal(env.Payload, &payload)
		return c.search.Index(ctx, service.Document{
			ID: env.AggregateID, Type: "question", Title: payload.Title,
		})
	default:
		// 删除/隐藏同步可在此扩展（调用 search.Delete）。
		return nil
	}
}
