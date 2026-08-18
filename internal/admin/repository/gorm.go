// Package repository 提供 admin-service 的 GORM 持久化实现。
package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/likeslep/community/internal/admin/model"
)

// Gorm 是 service.Repository 的 GORM 实现。
type Gorm struct {
	db *gorm.DB
}

// NewGorm 构造。
func NewGorm(db *gorm.DB) *Gorm { return &Gorm{db: db} }

func (r *Gorm) CreateAuditLog(ctx context.Context, l *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(l).Error
}

func (r *Gorm) ListAuditLogs(ctx context.Context, limit int) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	err := r.db.WithContext(ctx).Order("id DESC").Limit(limit).Find(&logs).Error
	return logs, err
}
