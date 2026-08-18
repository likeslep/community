package service

import (
	"context"

	"github.com/likeslep/community/internal/content/model"
	"github.com/likeslep/community/pkg/kafkax"
)

// BuildEvent 在文章 ID 已知后构造事件信封。
type BuildEvent func(a *model.Article) (kafkax.Envelope, error)

// Repository 是 content-service 的持久化接口（plan.md §9），业务层不依赖 GORM。
type Repository interface {
	// CreateArticle 创建文章（绑定标签 + outbox）在同一事务内。
	CreateArticle(ctx context.Context, a *model.Article, tagNames []string, build BuildEvent) error
	// UpdateArticle 更新文章（写版本历史 + outbox）在同一事务内。
	UpdateArticle(ctx context.Context, a *model.Article, build BuildEvent) error
	// FindArticle 查询文章。
	FindArticle(ctx context.Context, id uint64) (*model.Article, error)
	// DeleteArticle 软删文章（状态→deleted）。
	DeleteArticle(ctx context.Context, id uint64) error
	// SubmitArticle 更新状态（→pending_review）+ outbox 在同一事务内。
	SubmitArticle(ctx context.Context, a *model.Article, build BuildEvent) error
	// SaveWithEvent 保存文章状态 + outbox 在同一事务内（用于发布/驳回）。
	SaveWithEvent(ctx context.Context, a *model.Article, build BuildEvent) error
	// TagsByArticle 查询文章标签。
	TagsByArticle(ctx context.Context, articleID uint64) ([]model.Tag, error)
	// ListTags 查询所有标签。
	ListTags(ctx context.Context) ([]model.Tag, error)
	// UpdateArticleStatus 更新文章状态（管理员隐藏等）。
	UpdateArticleStatus(ctx context.Context, id uint64, status model.Status) error
}
