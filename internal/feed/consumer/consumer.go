// Package consumer 是 feed-service 的 Kafka 事件消费者。
package consumer

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/likeslep/community/internal/feed/service"
	"github.com/likeslep/community/pkg/kafkax"
)

// FollowersResolver 查询某用户的粉丝列表（生产环境通过 social-service gRPC）。
type FollowersResolver func(ctx context.Context, userID uint64) ([]uint64, error)

// EventConsumer 消费内容发布事件做 Fan-out on Write（plan.md §22.7）。
type EventConsumer struct {
	feed      *service.FeedService
	followers FollowersResolver
	logger    *zap.Logger
}

// New 构造。
func New(feed *service.FeedService, followers FollowersResolver, logger *zap.Logger) *EventConsumer {
	return &EventConsumer{feed: feed, followers: followers, logger: logger}
}

// Handle 处理一条事件。
func (c *EventConsumer) Handle(ctx context.Context, env kafkax.Envelope) error {
	if env.EventType != kafkax.EventArticlePublished {
		return nil
	}
	var p struct {
		ArticleID uint64 `json:"article_id"`
		AuthorID  uint64 `json:"author_id"`
		Title     string `json:"title"`
	}
	if err := json.Unmarshal(env.Payload, &p); err != nil {
		return err
	}
	if c.followers == nil {
		return nil
	}
	followers, err := c.followers(ctx, p.AuthorID)
	if err != nil {
		return err
	}
	item := service.FeedItem{ID: env.AggregateID, Type: "article", Title: p.Title, AuthorID: p.AuthorID, PublishedAt: env.OccurredAt}
	for _, uid := range followers {
		if err := c.feed.AddToFeed(ctx, uid, item); err != nil {
			return err
		}
	}
	return nil
}
