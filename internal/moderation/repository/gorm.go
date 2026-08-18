// Package repository 提供 moderation-service 的 GORM 持久化实现。
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/likeslep/community/internal/moderation/model"
	"github.com/likeslep/community/internal/moderation/service"
	"github.com/likeslep/community/pkg/kafkax"
	"github.com/likeslep/community/pkg/outbox"
)

// Gorm 是 service.Repository 的 GORM 实现。
type Gorm struct {
	db     *gorm.DB
	outbox *outbox.Store
}

// NewGorm 构造。
func NewGorm(db *gorm.DB, outboxStore *outbox.Store) *Gorm {
	return &Gorm{db: db, outbox: outboxStore}
}

func (r *Gorm) ListSensitiveWords(ctx context.Context) ([]model.SensitiveWord, error) {
	var words []model.SensitiveWord
	err := r.db.WithContext(ctx).Find(&words).Error
	return words, err
}

func (r *Gorm) CreateSensitiveWord(ctx context.Context, w *model.SensitiveWord) error {
	return r.db.WithContext(ctx).Create(w).Error
}

func (r *Gorm) CreateTask(ctx context.Context, t *model.ModerationTask) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *Gorm) ListTasks(ctx context.Context, limit int) ([]model.ModerationTask, error) {
	var tasks []model.ModerationTask
	err := r.db.WithContext(ctx).Order("id DESC").Limit(limit).Find(&tasks).Error
	return tasks, err
}

func (r *Gorm) FindTask(ctx context.Context, id uint64) (*model.ModerationTask, error) {
	var t model.ModerationTask
	if err := r.db.WithContext(ctx).First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrTaskNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *Gorm) ApproveTask(ctx context.Context, t *model.ModerationTask, env kafkax.Envelope) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(t).Error; err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

func (r *Gorm) RejectTask(ctx context.Context, t *model.ModerationTask, env kafkax.Envelope) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(t).Error; err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

func (r *Gorm) CreateReport(ctx context.Context, rep *model.Report) error {
	return r.db.WithContext(ctx).Create(rep).Error
}

func (r *Gorm) ListReports(ctx context.Context, limit int) ([]model.Report, error) {
	var reports []model.Report
	err := r.db.WithContext(ctx).Order("id DESC").Limit(limit).Find(&reports).Error
	return reports, err
}

func (r *Gorm) FindReport(ctx context.Context, id uint64) (*model.Report, error) {
	var rep model.Report
	if err := r.db.WithContext(ctx).First(&rep, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrReportNotFound
		}
		return nil, err
	}
	return &rep, nil
}

func (r *Gorm) UpdateReport(ctx context.Context, rep *model.Report) error {
	return r.db.WithContext(ctx).Save(rep).Error
}
