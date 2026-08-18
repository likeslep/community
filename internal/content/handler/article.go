// Package handler 是 content-service 的 gRPC 处理器（transport 层）。
package handler

import (
	"context"

	contentv1 "github.com/likeslep/community/api/gen/content/v1"
	"github.com/likeslep/community/internal/content/service"
	"github.com/likeslep/community/pkg/grpcx"
)

// ArticleHandler 实现 Article gRPC 服务。
type ArticleHandler struct {
	contentv1.UnimplementedArticleServiceServer
	svc *service.ArticleService
}

// NewArticleHandler 构造。
func NewArticleHandler(svc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

// CreateArticle 创建草稿。作者身份来自 gRPC metadata。
func (h *ArticleHandler) CreateArticle(ctx context.Context, req *contentv1.CreateArticleRequest) (*contentv1.CreateArticleResponse, error) {
	authorID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	a, err := h.svc.CreateArticle(ctx, authorID, req.GetTitle(), req.GetContent(), req.GetTags())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &contentv1.CreateArticleResponse{Id: a.ID, Status: string(a.Status)}, nil
}

// UpdateArticle 更新草稿（仅作者）。
func (h *ArticleHandler) UpdateArticle(ctx context.Context, req *contentv1.UpdateArticleRequest) (*contentv1.UpdateArticleResponse, error) {
	authorID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if _, err := h.svc.UpdateArticle(ctx, authorID, req.GetId(), req.GetTitle(), req.GetContent()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &contentv1.UpdateArticleResponse{}, nil
}

// GetArticle 查询文章。
func (h *ArticleHandler) GetArticle(ctx context.Context, req *contentv1.GetArticleRequest) (*contentv1.GetArticleResponse, error) {
	a, err := h.svc.GetArticle(ctx, req.GetId())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &contentv1.GetArticleResponse{
		Id: a.ID, AuthorId: a.AuthorID, Title: a.Title, Content: a.Content,
		Status: string(a.Status), Version: int32(a.Version),
	}, nil
}

// DeleteArticle 删除文章（仅作者）。
func (h *ArticleHandler) DeleteArticle(ctx context.Context, req *contentv1.DeleteArticleRequest) (*contentv1.DeleteArticleResponse, error) {
	authorID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.DeleteArticle(ctx, authorID, req.GetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &contentv1.DeleteArticleResponse{}, nil
}

// SubmitArticle 提交审核（仅作者）。
func (h *ArticleHandler) SubmitArticle(ctx context.Context, req *contentv1.SubmitArticleRequest) (*contentv1.SubmitArticleResponse, error) {
	authorID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.SubmitArticle(ctx, authorID, req.GetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &contentv1.SubmitArticleResponse{}, nil
}

// HideArticle 隐藏文章（管理员）。
func (h *ArticleHandler) HideArticle(ctx context.Context, req *contentv1.HideArticleRequest) (*contentv1.HideArticleResponse, error) {
	if err := h.svc.HideArticle(ctx, req.GetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &contentv1.HideArticleResponse{}, nil
}

// ListTags 查询所有标签。
func (h *ArticleHandler) ListTags(ctx context.Context, _ *contentv1.ListTagsRequest) (*contentv1.ListTagsResponse, error) {
	tags, err := h.svc.ListTags(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &contentv1.ListTagsResponse{}
	for _, t := range tags {
		resp.Tags = append(resp.Tags, &contentv1.Tag{Id: t.ID, Name: t.Name})
	}
	return resp, nil
}
