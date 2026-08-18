// Package repository 提供 notification-service 的 GORM 持久化实现。
package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/likeslep/community/internal/notification/model"
)

// Gorm 是 service.Repository 的 GORM 实现。
type Gorm struct {
	db *gorm.DB
}

// NewGorm 构造。
func NewGorm(db *gorm.DB) *Gorm { return &Gorm{db: db} }

func (r *Gorm) CreateNotification(ctx context.Context, n *model.Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *Gorm) ListNotifications(ctx context.Context, userID uint64, limit int) ([]model.Notification, error) {
	var ns []model.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("id DESC").
		Limit(limit).
		Find(&ns).Error
	return ns, err
}

func (r *Gorm) UnreadCount(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *Gorm) MarkRead(ctx context.Context, userID, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *Gorm) MarkAllRead(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}
