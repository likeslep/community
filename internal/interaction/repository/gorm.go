// Package repository 提供 interaction-service 的 GORM 持久化实现。
package repository

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/likeslep/community/internal/interaction/model"
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

// CreateComment 创建评论 + 计数 + outbox 在同一事务内。
func (r *Gorm) CreateComment(ctx context.Context, c *model.Comment, build func(*model.Comment) (kafkax.Envelope, error)) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		if err := incrCounter(tx, c.TargetType, c.TargetID, "comment_count"); err != nil {
			return err
		}
		env, err := build(c)
		if err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

// ListComments 查询评论。
func (r *Gorm) ListComments(ctx context.Context, targetType string, targetID uint64) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.WithContext(ctx).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("id ASC").
		Find(&comments).Error
	return comments, err
}

// Like 点赞：唯一约束兜底幂等，计数器同一事务内维护。
func (r *Gorm) Like(ctx context.Context, userID uint64, targetType string, targetID uint64, build func(*model.Like) (kafkax.Envelope, error)) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		l := &model.Like{UserID: userID, TargetType: targetType, TargetID: targetID}
		if err := tx.Create(l).Error; err != nil {
			if isDuplicate(err) {
				return nil // 幂等：重复点赞
			}
			return err
		}
		if err := incrCounter(tx, targetType, targetID, "like_count"); err != nil {
			return err
		}
		env, err := build(l)
		if err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

// Unlike 取消点赞：删除后计数减一。
func (r *Gorm) Unlike(ctx context.Context, userID uint64, targetType string, targetID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Delete(&model.Like{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil // 未点赞，幂等
		}
		return decrCounter(tx, targetType, targetID, "like_count")
	})
}

// Collect 收藏：幂等。
func (r *Gorm) Collect(ctx context.Context, userID uint64, targetType string, targetID uint64, build func(*model.Collection) (kafkax.Envelope, error)) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		c := &model.Collection{UserID: userID, TargetType: targetType, TargetID: targetID}
		if err := tx.Create(c).Error; err != nil {
			if isDuplicate(err) {
				return nil
			}
			return err
		}
		if err := incrCounter(tx, targetType, targetID, "collect_count"); err != nil {
			return err
		}
		env, err := build(c)
		if err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

// Uncollect 取消收藏。
func (r *Gorm) Uncollect(ctx context.Context, userID uint64, targetType string, targetID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Delete(&model.Collection{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}
		return decrCounter(tx, targetType, targetID, "collect_count")
	})
}

// GetCounter 查询计数。
func (r *Gorm) GetCounter(ctx context.Context, targetType string, targetID uint64) (*model.Counter, error) {
	var c model.Counter
	if err := r.db.WithContext(ctx).Where("target_type = ? AND target_id = ?", targetType, targetID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.Counter{TargetType: targetType, TargetID: targetID}, nil
		}
		return nil, err
	}
	return &c, nil
}

// DeleteComment 删除评论（计数减一）。
func (r *Gorm) DeleteComment(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var c model.Comment
		if err := tx.First(&c, id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.Comment{}, id).Error; err != nil {
			return err
		}
		return decrCounter(tx, c.TargetType, c.TargetID, "comment_count")
	})
}

// incrCounter 原子递增计数器（UPSERT）。
func incrCounter(tx *gorm.DB, targetType string, targetID uint64, column string) error {
	sql := "INSERT INTO counters (target_type, target_id, " + column + ", updated_at) VALUES (?, ?, 1, NOW()) " +
		"ON DUPLICATE KEY UPDATE " + column + " = " + column + " + 1, updated_at = NOW()"
	return tx.Exec(sql, targetType, targetID).Error
}

// decrCounter 原子递减计数器。
func decrCounter(tx *gorm.DB, targetType string, targetID uint64, column string) error {
	sql := "UPDATE counters SET " + column + " = GREATEST(" + column + " - 1, 0), updated_at = NOW() " +
		"WHERE target_type = ? AND target_id = ?"
	return tx.Exec(sql, targetType, targetID).Error
}

func isDuplicate(err error) bool {
	var me *mysql.MySQLError
	return errors.As(err, &me) && me.Number == 1062
}
