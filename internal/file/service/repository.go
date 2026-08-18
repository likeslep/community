package service

import (
	"context"

	"github.com/likeslep/community/internal/file/model"
)

// Repository 是 file-service 的持久化接口（plan.md §9）。
type Repository interface {
	Create(ctx context.Context, f *model.File) error
	Find(ctx context.Context, id uint64) (*model.File, error)
}
