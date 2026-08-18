// Package repository 提供 file-service 的 GORM 持久化实现。
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/likeslep/community/internal/file/model"
	"github.com/likeslep/community/internal/file/service"
)

// Gorm 是 service.Repository 的 GORM 实现。
type Gorm struct {
	db *gorm.DB
}

// NewGorm 构造。
func NewGorm(db *gorm.DB) *Gorm { return &Gorm{db: db} }

func (r *Gorm) Create(ctx context.Context, f *model.File) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *Gorm) Find(ctx context.Context, id uint64) (*model.File, error) {
	var f model.File
	if err := r.db.WithContext(ctx).First(&f, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrFileNotFound
		}
		return nil, err
	}
	return &f, nil
}
