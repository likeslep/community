// Package handler 是 social-service 的 gRPC 处理器。
package handler

import (
	"context"

	socialv1 "github.com/likeslep/community/api/gen/social/v1"
	"github.com/likeslep/community/internal/social/service"
	"github.com/likeslep/community/pkg/grpcx"
)

// Handler 实现 social gRPC 服务。
type Handler struct {
	socialv1.UnimplementedSocialServiceServer
	svc *service.SocialService
}

// New 构造。
func New(svc *service.SocialService) *Handler { return &Handler{svc: svc} }

func (h *Handler) FollowUser(ctx context.Context, req *socialv1.FollowUserRequest) (*socialv1.FollowUserResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.FollowUser(ctx, uid, req.GetFolloweeId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &socialv1.FollowUserResponse{}, nil
}

func (h *Handler) UnfollowUser(ctx context.Context, req *socialv1.UnfollowUserRequest) (*socialv1.UnfollowUserResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.UnfollowUser(ctx, uid, req.GetFolloweeId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &socialv1.UnfollowUserResponse{}, nil
}

func (h *Handler) FollowTag(ctx context.Context, req *socialv1.FollowTagRequest) (*socialv1.FollowTagResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.FollowTag(ctx, uid, req.GetTagId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &socialv1.FollowTagResponse{}, nil
}

func (h *Handler) UnfollowTag(ctx context.Context, req *socialv1.UnfollowTagRequest) (*socialv1.UnfollowTagResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.UnfollowTag(ctx, uid, req.GetTagId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &socialv1.UnfollowTagResponse{}, nil
}

func (h *Handler) ListFollowing(ctx context.Context, req *socialv1.ListFollowingRequest) (*socialv1.ListFollowingResponse, error) {
	ids, err := h.svc.ListFollowing(ctx, req.GetUserId())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &socialv1.ListFollowingResponse{UserIds: ids}, nil
}

func (h *Handler) ListFollowers(ctx context.Context, req *socialv1.ListFollowersRequest) (*socialv1.ListFollowersResponse, error) {
	ids, err := h.svc.ListFollowers(ctx, req.GetUserId())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &socialv1.ListFollowersResponse{UserIds: ids}, nil
}
