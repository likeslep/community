package service

import (
	"context"

	"github.com/likeslep/community/internal/social/model"
	"github.com/likeslep/community/pkg/kafkax"
)

// Repository 是 social-service 的持久化接口（plan.md §9）。
type Repository interface {
	FollowUser(ctx context.Context, followerID, followeeID uint64, build func(*model.Follow) (kafkax.Envelope, error)) error
	UnfollowUser(ctx context.Context, followerID, followeeID uint64) error
	FollowTag(ctx context.Context, userID, tagID uint64, build func(*model.TagFollow) (kafkax.Envelope, error)) error
	UnfollowTag(ctx context.Context, userID, tagID uint64) error
	ListFollowing(ctx context.Context, userID uint64) ([]uint64, error)
	ListFollowers(ctx context.Context, userID uint64) ([]uint64, error)
}
