// Package service 是 interaction-service 的业务逻辑层。
package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/likeslep/community/internal/interaction/model"
	"github.com/likeslep/community/pkg/kafkax"
	"github.com/likeslep/community/pkg/redisx"
)

// 合法的互动目标类型。
var validTargets = map[string]bool{
	"question": true,
	"answer":   true,
	"article":  true,
	"comment":  true,
}

// Config 是 interaction-service 的业务配置。
type Config struct {
	Producer string
	Redis    *redisx.Client
}

// InteractionService 是互动业务逻辑层。
type InteractionService struct {
	repo Repository
	cfg  Config
}

// NewInteractionService 构造。
func NewInteractionService(repo Repository, cfg Config) *InteractionService {
	return &InteractionService{repo: repo, cfg: cfg}
}

// CreateComment 创建评论。
func (s *InteractionService) CreateComment(ctx context.Context, userID uint64, targetType string, targetID uint64, content string) (*model.Comment, error) {
	if !validTargets[targetType] {
		return nil, ErrInvalidTarget
	}
	if strings.TrimSpace(content) == "" {
		return nil, ErrInvalidTarget
	}
	c := &model.Comment{UserID: userID, TargetType: targetType, TargetID: targetID, Content: content}
	err := s.repo.CreateComment(ctx, c, func(c *model.Comment) (kafkax.Envelope, error) {
		return kafkax.NewEnvelope(kafkax.EventCommentCreated, s.cfg.Producer, "comment",
			strconv.FormatUint(c.ID, 10), 1, commentPayload{CommentID: c.ID, UserID: c.UserID, TargetType: c.TargetType, TargetID: c.TargetID})
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ListComments 查询评论。
func (s *InteractionService) ListComments(ctx context.Context, targetType string, targetID uint64) ([]model.Comment, error) {
	return s.repo.ListComments(ctx, targetType, targetID)
}

// DeleteComment 删除评论（管理员）。
func (s *InteractionService) DeleteComment(ctx context.Context, id uint64) error {
	return s.repo.DeleteComment(ctx, id)
}

// Like 点赞（幂等）。
func (s *InteractionService) Like(ctx context.Context, userID uint64, targetType string, targetID uint64) error {
	if !validTargets[targetType] {
		return ErrInvalidTarget
	}
	return s.repo.Like(ctx, userID, targetType, targetID, func(l *model.Like) (kafkax.Envelope, error) {
		return kafkax.NewEnvelope(kafkax.EventLikeCreated, s.cfg.Producer, "like",
			strconv.FormatUint(l.ID, 10), 1, interactionPayload{UserID: l.UserID, TargetType: l.TargetType, TargetID: l.TargetID})
	})
}

// Unlike 取消点赞（幂等）。
func (s *InteractionService) Unlike(ctx context.Context, userID uint64, targetType string, targetID uint64) error {
	if !validTargets[targetType] {
		return ErrInvalidTarget
	}
	return s.repo.Unlike(ctx, userID, targetType, targetID)
}

// Collect 收藏（幂等）。
func (s *InteractionService) Collect(ctx context.Context, userID uint64, targetType string, targetID uint64) error {
	if !validTargets[targetType] {
		return ErrInvalidTarget
	}
	return s.repo.Collect(ctx, userID, targetType, targetID, func(c *model.Collection) (kafkax.Envelope, error) {
		return kafkax.NewEnvelope(kafkax.EventCollectionCreated, s.cfg.Producer, "collection",
			strconv.FormatUint(c.ID, 10), 1, interactionPayload{UserID: c.UserID, TargetType: c.TargetType, TargetID: c.TargetID})
	})
}

// Uncollect 取消收藏（幂等）。
func (s *InteractionService) Uncollect(ctx context.Context, userID uint64, targetType string, targetID uint64) error {
	if !validTargets[targetType] {
		return ErrInvalidTarget
	}
	return s.repo.Uncollect(ctx, userID, targetType, targetID)
}

// View 浏览计数（Redis 热计数，后续异步落库）。
func (s *InteractionService) View(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	if !validTargets[targetType] {
		return 0, ErrInvalidTarget
	}
	key := fmt.Sprintf("view:%s:%d", targetType, targetID)
	return s.cfg.Redis.Redis().Incr(ctx, key).Result()
}

// GetCounter 查询计数。
func (s *InteractionService) GetCounter(ctx context.Context, targetType string, targetID uint64) (*model.Counter, error) {
	return s.repo.GetCounter(ctx, targetType, targetID)
}

type commentPayload struct {
	CommentID  uint64 `json:"comment_id"`
	UserID     uint64 `json:"user_id"`
	TargetType string `json:"target_type"`
	TargetID   uint64 `json:"target_id"`
}

type interactionPayload struct {
	UserID     uint64 `json:"user_id"`
	TargetType string `json:"target_type"`
	TargetID   uint64 `json:"target_id"`
}
