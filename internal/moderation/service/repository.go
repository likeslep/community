package service

import (
	"context"

	"github.com/likeslep/community/internal/moderation/model"
	"github.com/likeslep/community/pkg/kafkax"
)

// Repository 是 moderation-service 的持久化接口（plan.md §9）。
type Repository interface {
	// 敏感词
	ListSensitiveWords(ctx context.Context) ([]model.SensitiveWord, error)
	CreateSensitiveWord(ctx context.Context, w *model.SensitiveWord) error

	// 审核任务
	CreateTask(ctx context.Context, t *model.ModerationTask) error
	ListTasks(ctx context.Context, limit int) ([]model.ModerationTask, error)
	FindTask(ctx context.Context, id uint64) (*model.ModerationTask, error)
	ApproveTask(ctx context.Context, t *model.ModerationTask, env kafkax.Envelope) error
	RejectTask(ctx context.Context, t *model.ModerationTask, env kafkax.Envelope) error

	// 举报
	CreateReport(ctx context.Context, r *model.Report) error
	ListReports(ctx context.Context, limit int) ([]model.Report, error)
	FindReport(ctx context.Context, id uint64) (*model.Report, error)
	UpdateReport(ctx context.Context, r *model.Report) error
}
