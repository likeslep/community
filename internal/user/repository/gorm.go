// Package repository 提供 user-service 的 GORM 持久化实现。
package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/likeslep/community/internal/user/model"
	"github.com/likeslep/community/internal/user/service"
	"github.com/likeslep/community/pkg/kafkax"
	"github.com/likeslep/community/pkg/outbox"
)

// Gorm 是 service.Repository 的 GORM 实现。
type Gorm struct {
	db     *gorm.DB
	outbox *outbox.Store
}

// NewGorm 构造 GORM 实现。
func NewGorm(db *gorm.DB, outboxStore *outbox.Store) *Gorm {
	return &Gorm{db: db, outbox: outboxStore}
}

// Create 在事务内创建用户，并调用 buildEvent 构造事件写入 outbox。
func (r *Gorm) Create(ctx context.Context, user *model.User, buildEvent func(*model.User) (kafkax.Envelope, error)) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return mapDupError(err)
		}
		env, err := buildEvent(user)
		if err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

// FindByUsername 按用户名查询。
func (r *Gorm) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		return nil, mapNotFound(err)
	}
	return &u, nil
}

// FindByEmail 按邮箱查询。
func (r *Gorm) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return nil, mapNotFound(err)
	}
	return &u, nil
}

// FindByID 按 ID 查询。
func (r *Gorm) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, mapNotFound(err)
	}
	return &u, nil
}

// Update 更新用户。
func (r *Gorm) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// ListUsers 分页查询用户。
func (r *Gorm) ListUsers(ctx context.Context, limit, offset int) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).Order("id ASC").Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// UpdateStatus 更新用户状态。
func (r *Gorm) UpdateStatus(ctx context.Context, id uint64, status string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

func mapNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return service.ErrUserNotFound
	}
	return err
}

func mapDupError(err error) error {
	if field := duplicateField(err); field == "email" {
		return service.ErrEmailTaken
	} else if field == "username" {
		return service.ErrUsernameTaken
	}
	return err
}

// duplicateField 从 MySQL 1062 错误中解析冲突的索引字段（username/email）。
func duplicateField(err error) string {
	var me *mysql.MySQLError
	if errors.As(err, &me) && me.Number == 1062 {
		if strings.Contains(me.Message, "uk_email") {
			return "email"
		}
		return "username"
	}
	return ""
}
