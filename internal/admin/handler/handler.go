// Package handler 是 admin-service 的 gRPC 处理器。
package handler

import (
	"context"

	adminv1 "github.com/likeslep/community/api/gen/admin/v1"
	"github.com/likeslep/community/internal/admin/service"
	"github.com/likeslep/community/pkg/grpcx"
)

// Handler 实现 admin gRPC 服务。
type Handler struct {
	adminv1.UnimplementedAdminServiceServer
	svc *service.AdminService
}

// New 构造。
func New(svc *service.AdminService) *Handler { return &Handler{svc: svc} }

// ListReports 查询举报。
func (h *Handler) ListReports(ctx context.Context, req *adminv1.ListReportsRequest) (*adminv1.ListReportsResponse, error) {
	reports, err := h.svc.ListReports(ctx, int(req.GetLimit()))
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &adminv1.ListReportsResponse{}
	for _, r := range reports {
		resp.Reports = append(resp.Reports, &adminv1.Report{
			Id: r.Id, ReporterId: r.ReporterId, TargetType: r.TargetType,
			TargetId: r.TargetId, Reason: r.Reason, Status: r.Status,
		})
	}
	return resp, nil
}

// HandleReport 处理举报。
func (h *Handler) HandleReport(ctx context.Context, req *adminv1.HandleReportRequest) (*adminv1.HandleReportResponse, error) {
	adminID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.HandleReport(ctx, adminID, req.GetId(), req.GetAction()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &adminv1.HandleReportResponse{}, nil
}

// ListAuditLogs 查询审计日志。
func (h *Handler) ListAuditLogs(ctx context.Context, req *adminv1.ListAuditLogsRequest) (*adminv1.ListAuditLogsResponse, error) {
	logs, err := h.svc.ListAuditLogs(ctx, int(req.GetLimit()))
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &adminv1.ListAuditLogsResponse{}
	for _, l := range logs {
		resp.Logs = append(resp.Logs, &adminv1.AuditLog{
			Id: l.ID, AdminId: l.AdminID, Action: l.Action,
			TargetType: l.TargetType, TargetId: l.TargetID, Detail: l.Detail,
		})
	}
	return resp, nil
}

// ListUsers 查询用户。
func (h *Handler) ListUsers(ctx context.Context, req *adminv1.ListUsersRequest) (*adminv1.ListUsersResponse, error) {
	users, err := h.svc.ListUsers(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &adminv1.ListUsersResponse{}
	for _, u := range users {
		resp.Users = append(resp.Users, &adminv1.UserBrief{
			Id: u.Id, Username: u.Username, Role: u.Role, Status: u.Status,
		})
	}
	return resp, nil
}

// BanUser 封禁用户。
func (h *Handler) BanUser(ctx context.Context, req *adminv1.BanUserRequest) (*adminv1.BanUserResponse, error) {
	adminID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.BanUser(ctx, adminID, req.GetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &adminv1.BanUserResponse{}, nil
}

// HideArticle 隐藏文章。
func (h *Handler) HideArticle(ctx context.Context, req *adminv1.HideArticleRequest) (*adminv1.HideArticleResponse, error) {
	adminID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.HideArticle(ctx, adminID, req.GetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &adminv1.HideArticleResponse{}, nil
}

// DeleteComment 删除评论。
func (h *Handler) DeleteComment(ctx context.Context, req *adminv1.DeleteCommentRequest) (*adminv1.DeleteCommentResponse, error) {
	adminID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.DeleteComment(ctx, adminID, req.GetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &adminv1.DeleteCommentResponse{}, nil
}

// ListTags 查询标签。
func (h *Handler) ListTags(ctx context.Context, _ *adminv1.ListTagsRequest) (*adminv1.ListTagsResponse, error) {
	tags, err := h.svc.ListTags(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &adminv1.ListTagsResponse{}
	for _, t := range tags {
		resp.Tags = append(resp.Tags, &adminv1.Tag{Id: t.Id, Name: t.Name})
	}
	return resp, nil
}

// ListSensitiveWords 查询敏感词。
func (h *Handler) ListSensitiveWords(ctx context.Context, _ *adminv1.ListSensitiveWordsRequest) (*adminv1.ListSensitiveWordsResponse, error) {
	words, err := h.svc.ListSensitiveWords(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &adminv1.ListSensitiveWordsResponse{}
	for _, w := range words {
		resp.Words = append(resp.Words, &adminv1.SensitiveWord{Id: w.Id, Word: w.Word, Level: w.Level})
	}
	return resp, nil
}

// CreateSensitiveWord 添加敏感词。
func (h *Handler) CreateSensitiveWord(ctx context.Context, req *adminv1.CreateSensitiveWordRequest) (*adminv1.CreateSensitiveWordResponse, error) {
	adminID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.CreateSensitiveWord(ctx, adminID, req.GetWord(), req.GetLevel()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &adminv1.CreateSensitiveWordResponse{}, nil
}

// Statistics 数据统计。
func (h *Handler) Statistics(ctx context.Context, _ *adminv1.StatisticsRequest) (*adminv1.StatisticsResponse, error) {
	stats, err := h.svc.Statistics(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &adminv1.StatisticsResponse{
		UserCount: stats.UserCount, ReportCount: stats.ReportCount, TagCount: stats.TagCount,
		SensitiveWordCount: stats.SensitiveWordCount, AuditLogCount: stats.AuditLogCount,
	}, nil
}
