// Package service 是 social-service 的业务逻辑层。
package service

import (
	"context"
	"strconv"

	"github.com/likeslep/community/internal/social/model"
	"github.com/likeslep/community/pkg/kafkax"
)

// Config 是 social-service 的业务配置。
type Config struct {
	Producer string
}

// SocialService 是社交业务逻辑层。
type SocialService struct {
	repo Repository
	cfg  Config
}

// NewSocialService 构造。
func NewSocialService(repo Repository, cfg Config) *SocialService {
	return &SocialService{repo: repo, cfg: cfg}
}

// FollowUser 关注用户（幂等）。
func (s *SocialService) FollowUser(ctx context.Context, followerID, followeeID uint64) error {
	if followerID == followeeID {
		return ErrSelfFollow
	}
	return s.repo.FollowUser(ctx, followerID, followeeID, func(f *model.Follow) (kafkax.Envelope, error) {
		return kafkax.NewEnvelope(kafkax.EventUserFollowed, s.cfg.Producer, "follow",
			strconv.FormatUint(f.ID, 10), 1, followPayload{FollowerID: f.FollowerID, FolloweeID: f.FolloweeID})
	})
}

// UnfollowUser 取关用户（幂等）。
func (s *SocialService) UnfollowUser(ctx context.Context, followerID, followeeID uint64) error {
	return s.repo.UnfollowUser(ctx, followerID, followeeID)
}

// FollowTag 关注标签（幂等）。
func (s *SocialService) FollowTag(ctx context.Context, userID, tagID uint64) error {
	return s.repo.FollowTag(ctx, userID, tagID, func(t *model.TagFollow) (kafkax.Envelope, error) {
		return kafkax.NewEnvelope(kafkax.EventTagFollowed, s.cfg.Producer, "tag_follow",
			strconv.FormatUint(t.ID, 10), 1, tagFollowPayload{UserID: t.UserID, TagID: t.TagID})
	})
}

// UnfollowTag 取消关注标签（幂等）。
func (s *SocialService) UnfollowTag(ctx context.Context, userID, tagID uint64) error {
	return s.repo.UnfollowTag(ctx, userID, tagID)
}

// ListFollowing 查询关注列表。
func (s *SocialService) ListFollowing(ctx context.Context, userID uint64) ([]uint64, error) {
	return s.repo.ListFollowing(ctx, userID)
}

// ListFollowers 查询粉丝列表。
func (s *SocialService) ListFollowers(ctx context.Context, userID uint64) ([]uint64, error) {
	return s.repo.ListFollowers(ctx, userID)
}

type followPayload struct {
	FollowerID uint64 `json:"follower_id"`
	FolloweeID uint64 `json:"followee_id"`
}

type tagFollowPayload struct {
	UserID uint64 `json:"user_id"`
	TagID  uint64 `json:"tag_id"`
}
