package service

import (
	"context"

	"github.com/likeslep/community/internal/admin/model"
)

// Repository 是 admin-service 的持久化接口（plan.md §9）。
type Repository interface {
	CreateAuditLog(ctx context.Context, l *model.AuditLog) error
	ListAuditLogs(ctx context.Context, limit int) ([]model.AuditLog, error)
}
