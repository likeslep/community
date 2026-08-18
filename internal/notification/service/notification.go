// Package service 是 notification-service 的业务逻辑层。
package service

import (
	"context"
	"encoding/json"

	"github.com/likeslep/community/internal/notification/model"
	"github.com/likeslep/community/pkg/kafkax"
)

// NotificationService 是通知业务逻辑层。
type NotificationService struct {
	repo Repository
}

// NewNotificationService 构造。
func NewNotificationService(repo Repository) *NotificationService {
	return &NotificationService{repo: repo}
}

// HandleEvent 消费领域事件并生成站内通知（plan.md §22.6）。
// 注：like/comment 等需要查询目标归属者的事件，后续通过 gRPC 解析 recipient。
func (s *NotificationService) HandleEvent(ctx context.Context, env kafkax.Envelope) error {
	switch env.EventType {
	case kafkax.EventUserFollowed:
		var p struct {
			FollowerID uint64 `json:"follower_id"`
			FolloweeID uint64 `json:"followee_id"`
		}
		if err := json.Unmarshal(env.Payload, &p); err != nil {
			return err
		}
		return s.repo.CreateNotification(ctx, &model.Notification{
			UserID: p.FolloweeID, ActorID: p.FollowerID, Type: "follow", Content: "有人关注了你",
		})
	case kafkax.EventAnswerAccepted:
		var p struct {
			AnswerID   uint64 `json:"answer_id"`
			QuestionID uint64 `json:"question_id"`
			AuthorID   uint64 `json:"author_id"`
		}
		if err := json.Unmarshal(env.Payload, &p); err != nil {
			return err
		}
		return s.repo.CreateNotification(ctx, &model.Notification{
			UserID: p.AuthorID, Type: "answer_accepted", TargetType: "question", TargetID: p.QuestionID,
			Content: "你的回答被采纳了",
		})
	default:
		return nil // 忽略其它事件
	}
}

// ListNotifications 查询通知列表。
func (s *NotificationService) ListNotifications(ctx context.Context, userID uint64, limit int) ([]model.Notification, error) {
	return s.repo.ListNotifications(ctx, userID, limit)
}

// UnreadCount 查询未读数。
func (s *NotificationService) UnreadCount(ctx context.Context, userID uint64) (int64, error) {
	return s.repo.UnreadCount(ctx, userID)
}

// MarkRead 标记单条已读。
func (s *NotificationService) MarkRead(ctx context.Context, userID, id uint64) error {
	return s.repo.MarkRead(ctx, userID, id)
}

// MarkAllRead 标记全部已读。
func (s *NotificationService) MarkAllRead(ctx context.Context, userID uint64) error {
	return s.repo.MarkAllRead(ctx, userID)
}
