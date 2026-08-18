package service

import (
	"context"

	"github.com/likeslep/community/internal/notification/model"
)

// Repository 是 notification-service 的持久化接口（plan.md §9）。
type Repository interface {
	CreateNotification(ctx context.Context, n *model.Notification) error
	ListNotifications(ctx context.Context, userID uint64, limit int) ([]model.Notification, error)
	UnreadCount(ctx context.Context, userID uint64) (int64, error)
	MarkRead(ctx context.Context, userID, id uint64) error
	MarkAllRead(ctx context.Context, userID uint64) error
}
