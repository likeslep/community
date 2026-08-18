// Package handler 是 interaction-service 的 gRPC 处理器。
package handler

import (
	"context"

	interactionv1 "github.com/likeslep/community/api/gen/interaction/v1"
	"github.com/likeslep/community/internal/interaction/service"
	"github.com/likeslep/community/pkg/grpcx"
)

// Handler 实现互动 gRPC 服务。
type Handler struct {
	interactionv1.UnimplementedInteractionServiceServer
	interactionv1.UnimplementedCommentServiceServer
	svc *service.InteractionService
}

// New 构造。
func New(svc *service.InteractionService) *Handler { return &Handler{svc: svc} }

func (h *Handler) Like(ctx context.Context, req *interactionv1.LikeRequest) (*interactionv1.LikeResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.Like(ctx, uid, req.GetTargetType(), req.GetTargetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &interactionv1.LikeResponse{}, nil
}

func (h *Handler) Unlike(ctx context.Context, req *interactionv1.UnlikeRequest) (*interactionv1.UnlikeResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.Unlike(ctx, uid, req.GetTargetType(), req.GetTargetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &interactionv1.UnlikeResponse{}, nil
}

func (h *Handler) Collect(ctx context.Context, req *interactionv1.CollectRequest) (*interactionv1.CollectResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.Collect(ctx, uid, req.GetTargetType(), req.GetTargetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &interactionv1.CollectResponse{}, nil
}

func (h *Handler) Uncollect(ctx context.Context, req *interactionv1.UncollectRequest) (*interactionv1.UncollectResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.Uncollect(ctx, uid, req.GetTargetType(), req.GetTargetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &interactionv1.UncollectResponse{}, nil
}

func (h *Handler) View(ctx context.Context, req *interactionv1.ViewRequest) (*interactionv1.ViewResponse, error) {
	count, err := h.svc.View(ctx, req.GetTargetType(), req.GetTargetId())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &interactionv1.ViewResponse{ViewCount: count}, nil
}

func (h *Handler) CreateComment(ctx context.Context, req *interactionv1.CreateCommentRequest) (*interactionv1.CreateCommentResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	c, err := h.svc.CreateComment(ctx, uid, req.GetTargetType(), req.GetTargetId(), req.GetContent())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &interactionv1.CreateCommentResponse{Id: c.ID}, nil
}

func (h *Handler) ListComments(ctx context.Context, req *interactionv1.ListCommentsRequest) (*interactionv1.ListCommentsResponse, error) {
	comments, err := h.svc.ListComments(ctx, req.GetTargetType(), req.GetTargetId())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &interactionv1.ListCommentsResponse{}
	for _, c := range comments {
		resp.Comments = append(resp.Comments, &interactionv1.Comment{
			Id: c.ID, UserId: c.UserID, TargetType: c.TargetType, TargetId: c.TargetID, Content: c.Content,
		})
	}
	return resp, nil
}

// DeleteComment 删除评论（管理员）。
func (h *Handler) DeleteComment(ctx context.Context, req *interactionv1.DeleteCommentRequest) (*interactionv1.DeleteCommentResponse, error) {
	if err := h.svc.DeleteComment(ctx, req.GetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &interactionv1.DeleteCommentResponse{}, nil
}
