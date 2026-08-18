// Package repository 提供 content-service 的 GORM 持久化实现。
package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/likeslep/community/internal/content/model"
	"github.com/likeslep/community/internal/content/service"
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

// CreateArticle 创建文章（绑定标签 + outbox）在同一事务内。
func (r *Gorm) CreateArticle(ctx context.Context, a *model.Article, tagNames []string, build service.BuildEvent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(a).Error; err != nil {
			return err
		}
		if err := bindTags(ctx, tx, a.ID, tagNames); err != nil {
			return err
		}
		env, err := build(a)
		if err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

// UpdateArticle 更新文章（写版本历史 + outbox）在同一事务内。
func (r *Gorm) UpdateArticle(ctx context.Context, a *model.Article, build service.BuildEvent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(a).Error; err != nil {
			return err
		}
		v := &model.ArticleVersion{ArticleID: a.ID, Title: a.Title, Content: a.Content, Version: a.Version}
		if err := tx.Create(v).Error; err != nil {
			return err
		}
		env, err := build(a)
		if err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

// FindArticle 查询文章。
func (r *Gorm) FindArticle(ctx context.Context, id uint64) (*model.Article, error) {
	var a model.Article
	if err := r.db.WithContext(ctx).First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrArticleNotFound
		}
		return nil, err
	}
	return &a, nil
}

// DeleteArticle 软删文章。
func (r *Gorm) DeleteArticle(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.Article{}).
		Where("id = ?", id).
		Update("status", model.StatusDeleted).Error
}

// SubmitArticle 更新状态并写 outbox。
func (r *Gorm) SubmitArticle(ctx context.Context, a *model.Article, build service.BuildEvent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(a).Error; err != nil {
			return err
		}
		env, err := build(a)
		if err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

// SaveWithEvent 保存文章状态并写 outbox（用于发布/驳回）。
func (r *Gorm) SaveWithEvent(ctx context.Context, a *model.Article, build service.BuildEvent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(a).Error; err != nil {
			return err
		}
		env, err := build(a)
		if err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

// TagsByArticle 查询文章标签。
func (r *Gorm) TagsByArticle(ctx context.Context, articleID uint64) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.WithContext(ctx).
		Joins("JOIN article_tags ON article_tags.tag_id = tags.id").
		Where("article_tags.article_id = ?", articleID).
		Find(&tags).Error
	return tags, err
}

// ListTags 查询所有标签。
func (r *Gorm) ListTags(ctx context.Context) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.WithContext(ctx).Find(&tags).Error
	return tags, err
}

// UpdateArticleStatus 更新文章状态。
func (r *Gorm) UpdateArticleStatus(ctx context.Context, id uint64, status model.Status) error {
	return r.db.WithContext(ctx).Model(&model.Article{}).Where("id = ?", id).Update("status", status).Error
}

// bindTags 查找或创建标签并绑定到文章。
func bindTags(ctx context.Context, tx *gorm.DB, articleID uint64, names []string) error {
	for _, raw := range names {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		var tag model.Tag
		if err := tx.WithContext(ctx).Where("name = ?", name).FirstOrCreate(&tag, model.Tag{Name: name}).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Create(&model.ArticleTag{ArticleID: articleID, TagID: tag.ID}).Error; err != nil {
			return err
		}
	}
	return nil
}
