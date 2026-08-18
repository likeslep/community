// Package handler 是 search-service 的 gRPC 处理器。
package handler

import (
	"context"

	searchv1 "github.com/likeslep/community/api/gen/search/v1"
	"github.com/likeslep/community/internal/search/service"
	"github.com/likeslep/community/pkg/grpcx"
)

// Handler 实现 search gRPC 服务。
type Handler struct {
	searchv1.UnimplementedSearchServiceServer
	svc *service.SearchService
}

// New 构造。
func New(svc *service.SearchService) *Handler { return &Handler{svc: svc} }

// Search 搜索。
func (h *Handler) Search(ctx context.Context, req *searchv1.SearchRequest) (*searchv1.SearchResponse, error) {
	results, total, err := h.svc.Search(ctx, req.GetKeyword(), req.GetType(), req.GetTags(), int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &searchv1.SearchResponse{Total: total}
	for _, r := range results {
		resp.Results = append(resp.Results, &searchv1.SearchResult{
			Id: r.ID, Type: r.Type, Title: r.Title, Snippet: r.Snippet, AuthorId: r.AuthorID,
		})
	}
	return resp, nil
}

// Rebuild 全量重建索引（骨架）。
func (h *Handler) Rebuild(ctx context.Context, _ *searchv1.RebuildRequest) (*searchv1.RebuildResponse, error) {
	if err := h.svc.EnsureIndex(ctx); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &searchv1.RebuildResponse{}, nil
}
