// Package handler 是 feed-service 的 gRPC 处理器。
package handler

import (
	"context"

	feedv1 "github.com/likeslep/community/api/gen/feed/v1"
	"github.com/likeslep/community/internal/feed/service"
	"github.com/likeslep/community/pkg/grpcx"
)

// Handler 实现 feed gRPC 服务。
type Handler struct {
	feedv1.UnimplementedFeedServiceServer
	svc *service.FeedService
}

// New 构造。
func New(svc *service.FeedService) *Handler { return &Handler{svc: svc} }

// GetFeed 查询信息流。
func (h *Handler) GetFeed(ctx context.Context, req *feedv1.GetFeedRequest) (*feedv1.GetFeedResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	items, err := h.svc.GetFeed(ctx, uid, int(req.GetLimit()))
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &feedv1.GetFeedResponse{}
	for _, it := range items {
		resp.Items = append(resp.Items, &feedv1.FeedItem{
			Id: it.ID, Type: it.Type, Title: it.Title, AuthorId: it.AuthorID,
		})
	}
	return resp, nil
}
