package service

import (
	"context"

	"github.com/likeslep/community/internal/user/model"
	"github.com/likeslep/community/pkg/kafkax"
)

// Repository 是 service 层定义的持久化接口（plan.md §9），业务层不依赖 GORM。
type Repository interface {
	// Create 创建用户，并在同一事务内调用 buildEvent 构造事件写入 outbox。
	Create(ctx context.Context, user *model.User, buildEvent func(*model.User) (kafkax.Envelope, error)) error
	// FindByUsername 按用户名查询，未找到返回 ErrUserNotFound。
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	// FindByEmail 按邮箱查询，未找到返回 ErrUserNotFound。
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	// FindByID 按 ID 查询，未找到返回 ErrUserNotFound。
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	// Update 更新用户。
	Update(ctx context.Context, user *model.User) error
	// ListUsers 分页查询用户。
	ListUsers(ctx context.Context, limit, offset int) ([]model.User, error)
	// UpdateStatus 更新用户状态（封禁/解禁）。
	UpdateStatus(ctx context.Context, id uint64, status string) error
}
