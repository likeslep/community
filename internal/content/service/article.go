package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/likeslep/community/internal/content/model"
	"github.com/likeslep/community/pkg/apperr"
	"github.com/likeslep/community/pkg/kafkax"
)

// Config 是 content-service 的业务配置。
type Config struct {
	Producer string // 服务名，用于事件 producer 字段
}

// ArticleService 是文章业务逻辑层。
type ArticleService struct {
	repo Repository
	cfg  Config
}

// NewArticleService 构造。
func NewArticleService(repo Repository, cfg Config) *ArticleService {
	return &ArticleService{repo: repo, cfg: cfg}
}

// CreateArticle 创建草稿（默认 DRAFT）。
func (s *ArticleService) CreateArticle(ctx context.Context, authorID uint64, title, content string, tags []string) (*model.Article, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, apperr.New(errCodeInvalidInput, "标题不能为空", apperr.WithHTTP(400))
	}
	a := &model.Article{
		AuthorID: authorID,
		Title:    title,
		Content:  content,
		Status:   model.StatusDraft,
		Version:  1,
	}
	err := s.repo.CreateArticle(ctx, a, tags, func(a *model.Article) (kafkax.Envelope, error) {
		return kafkax.NewEnvelope(kafkax.EventArticleCreated, s.cfg.Producer, "article",
			strconv.FormatUint(a.ID, 10), a.Version, articlePayload{ArticleID: a.ID, AuthorID: a.AuthorID, Title: a.Title})
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

// UpdateArticle 更新草稿并写入版本历史。
func (s *ArticleService) UpdateArticle(ctx context.Context, authorID, id uint64, title, content string) (*model.Article, error) {
	a, err := s.repo.FindArticle(ctx, id)
	if err != nil {
		return nil, err
	}
	if a.AuthorID != authorID {
		return nil, ErrForbidden
	}
	if a.Status != model.StatusDraft && a.Status != model.StatusRejected {
		return nil, ErrIllegalState
	}
	if strings.TrimSpace(title) == "" {
		return nil, apperr.New(errCodeInvalidInput, "标题不能为空", apperr.WithHTTP(400))
	}
	a.Title = strings.TrimSpace(title)
	a.Content = content
	a.Version++

	err = s.repo.UpdateArticle(ctx, a, func(a *model.Article) (kafkax.Envelope, error) {
		return kafkax.NewEnvelope(kafkax.EventArticleUpdated, s.cfg.Producer, "article",
			strconv.FormatUint(a.ID, 10), a.Version, articlePayload{ArticleID: a.ID, AuthorID: a.AuthorID, Title: a.Title})
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

// GetArticle 查询文章。
func (s *ArticleService) GetArticle(ctx context.Context, id uint64) (*model.Article, error) {
	return s.repo.FindArticle(ctx, id)
}

// DeleteArticle 软删文章（仅作者）。
func (s *ArticleService) DeleteArticle(ctx context.Context, authorID, id uint64) error {
	a, err := s.repo.FindArticle(ctx, id)
	if err != nil {
		return err
	}
	if a.AuthorID != authorID {
		return ErrForbidden
	}
	if err := a.Transition(model.StatusDeleted); err != nil {
		return ErrIllegalState
	}
	return s.repo.DeleteArticle(ctx, id)
}

// SubmitArticle 提交审核（DRAFT/REJECTED → PENDING_REVIEW）。
func (s *ArticleService) SubmitArticle(ctx context.Context, authorID, id uint64) error {
	a, err := s.repo.FindArticle(ctx, id)
	if err != nil {
		return err
	}
	if a.AuthorID != authorID {
		return ErrForbidden
	}
	if err := a.Transition(model.StatusPendingReview); err != nil {
		return ErrIllegalState
	}
	return s.repo.SubmitArticle(ctx, a, func(a *model.Article) (kafkax.Envelope, error) {
		return kafkax.NewEnvelope(kafkax.EventArticleSubmitted, s.cfg.Producer, "article",
			strconv.FormatUint(a.ID, 10), a.Version, articlePayload{ArticleID: a.ID, AuthorID: a.AuthorID, Title: a.Title})
	})
}

// Tags 查询文章标签。
func (s *ArticleService) Tags(ctx context.Context, id uint64) ([]model.Tag, error) {
	return s.repo.TagsByArticle(ctx, id)
}

// HideArticle 隐藏文章（管理员）。
func (s *ArticleService) HideArticle(ctx context.Context, id uint64) error {
	a, err := s.repo.FindArticle(ctx, id)
	if err != nil {
		return err
	}
	if err := a.Transition(model.StatusHidden); err != nil {
		return ErrIllegalState
	}
	return s.repo.UpdateArticleStatus(ctx, id, model.StatusHidden)
}

// ListTags 查询所有标签。
func (s *ArticleService) ListTags(ctx context.Context) ([]model.Tag, error) {
	return s.repo.ListTags(ctx)
}

// PublishArticle 审核通过后发布文章（PENDING_REVIEW → PUBLISHED）。
func (s *ArticleService) PublishArticle(ctx context.Context, id uint64) error {
	a, err := s.repo.FindArticle(ctx, id)
	if err != nil {
		return err
	}
	if err := a.Transition(model.StatusPublished); err != nil {
		return ErrIllegalState
	}
	return s.repo.SaveWithEvent(ctx, a, func(a *model.Article) (kafkax.Envelope, error) {
		return kafkax.NewEnvelope(kafkax.EventArticlePublished, s.cfg.Producer, "article",
			strconv.FormatUint(a.ID, 10), a.Version, articlePayload{ArticleID: a.ID, AuthorID: a.AuthorID, Title: a.Title})
	})
}

// RejectArticle 审核驳回文章（PENDING_REVIEW → REJECTED）。
func (s *ArticleService) RejectArticle(ctx context.Context, id uint64) error {
	a, err := s.repo.FindArticle(ctx, id)
	if err != nil {
		return err
	}
	if err := a.Transition(model.StatusRejected); err != nil {
		return ErrIllegalState
	}
	return s.repo.SaveWithEvent(ctx, a, func(a *model.Article) (kafkax.Envelope, error) {
		return kafkax.NewEnvelope(kafkax.EventArticleRejected, s.cfg.Producer, "article",
			strconv.FormatUint(a.ID, 10), a.Version, articlePayload{ArticleID: a.ID, AuthorID: a.AuthorID, Title: a.Title})
	})
}

type articlePayload struct {
	ArticleID uint64 `json:"article_id"`
	AuthorID  uint64 `json:"author_id"`
	Title     string `json:"title"`
}
