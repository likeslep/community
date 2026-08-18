// Package repository 提供 social-service 的 GORM 持久化实现。
package repository

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/likeslep/community/internal/social/model"
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

// FollowUser 关注（唯一约束兜底幂等 + outbox）。
func (r *Gorm) FollowUser(ctx context.Context, followerID, followeeID uint64, build func(*model.Follow) (kafkax.Envelope, error)) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		f := &model.Follow{FollowerID: followerID, FolloweeID: followeeID}
		if err := tx.Create(f).Error; err != nil {
			if isDuplicate(err) {
				return nil
			}
			return err
		}
		env, err := build(f)
		if err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

// UnfollowUser 取关（幂等）。
func (r *Gorm) UnfollowUser(ctx context.Context, followerID, followeeID uint64) error {
	return r.db.WithContext(ctx).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&model.Follow{}).Error
}

// FollowTag 关注标签。
func (r *Gorm) FollowTag(ctx context.Context, userID, tagID uint64, build func(*model.TagFollow) (kafkax.Envelope, error)) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		t := &model.TagFollow{UserID: userID, TagID: tagID}
		if err := tx.Create(t).Error; err != nil {
			if isDuplicate(err) {
				return nil
			}
			return err
		}
		env, err := build(t)
		if err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

// UnfollowTag 取消关注标签。
func (r *Gorm) UnfollowTag(ctx context.Context, userID, tagID uint64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND tag_id = ?", userID, tagID).
		Delete(&model.TagFollow{}).Error
}

// ListFollowing 查询关注列表。
func (r *Gorm) ListFollowing(ctx context.Context, userID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.db.WithContext(ctx).
		Model(&model.Follow{}).
		Where("follower_id = ?", userID).
		Pluck("followee_id", &ids).Error
	return ids, err
}

// ListFollowers 查询粉丝列表。
func (r *Gorm) ListFollowers(ctx context.Context, userID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.db.WithContext(ctx).
		Model(&model.Follow{}).
		Where("followee_id = ?", userID).
		Pluck("follower_id", &ids).Error
	return ids, err
}

func isDuplicate(err error) bool {
	var me *mysql.MySQLError
	return errors.As(err, &me) && me.Number == 1062
}
