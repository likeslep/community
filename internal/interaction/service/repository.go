package service

import (
	"context"

	"github.com/likeslep/community/internal/interaction/model"
	"github.com/likeslep/community/pkg/kafkax"
)

// Repository 是 interaction-service 的持久化接口（plan.md §9）。
type Repository interface {
	CreateComment(ctx context.Context, c *model.Comment, build func(*model.Comment) (kafkax.Envelope, error)) error
	ListComments(ctx context.Context, targetType string, targetID uint64) ([]model.Comment, error)

	// Like / Unlike：幂等（唯一约束），计数器在同一事务内维护。
	Like(ctx context.Context, userID uint64, targetType string, targetID uint64, build func(*model.Like) (kafkax.Envelope, error)) error
	Unlike(ctx context.Context, userID uint64, targetType string, targetID uint64) error

	// Collect / Uncollect：幂等。
	Collect(ctx context.Context, userID uint64, targetType string, targetID uint64, build func(*model.Collection) (kafkax.Envelope, error)) error
	Uncollect(ctx context.Context, userID uint64, targetType string, targetID uint64) error

	// GetCounter 查询计数（MySQL 为准）。
	GetCounter(ctx context.Context, targetType string, targetID uint64) (*model.Counter, error)
	// DeleteComment 删除评论（计数减一）。
	DeleteComment(ctx context.Context, id uint64) error
}
