// Package handler 是 moderation-service 的 gRPC 处理器。
package handler

import (
	"context"

	moderationv1 "github.com/likeslep/community/api/gen/moderation/v1"
	"github.com/likeslep/community/internal/moderation/service"
	"github.com/likeslep/community/pkg/grpcx"
)

// Handler 实现 moderation 的 gRPC 服务。
type Handler struct {
	moderationv1.UnimplementedSensitiveWordServiceServer
	moderationv1.UnimplementedModerationServiceServer
	moderationv1.UnimplementedReportServiceServer
	svc *service.ModerationService
}

// New 构造。
func New(svc *service.ModerationService) *Handler { return &Handler{svc: svc} }

// CheckText 敏感词检测。
func (h *Handler) CheckText(ctx context.Context, req *moderationv1.CheckTextRequest) (*moderationv1.CheckTextResponse, error) {
	action, matched, err := h.svc.CheckText(ctx, req.GetText())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &moderationv1.CheckTextResponse{Action: action, Matched: matched}, nil
}

// CreateSensitiveWord 添加敏感词。
func (h *Handler) CreateSensitiveWord(ctx context.Context, req *moderationv1.CreateSensitiveWordRequest) (*moderationv1.CreateSensitiveWordResponse, error) {
	if err := h.svc.CreateSensitiveWord(ctx, req.GetWord(), req.GetLevel()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &moderationv1.CreateSensitiveWordResponse{}, nil
}

// ListSensitiveWords 查询敏感词列表。
func (h *Handler) ListSensitiveWords(ctx context.Context, _ *moderationv1.ListSensitiveWordsRequest) (*moderationv1.ListSensitiveWordsResponse, error) {
	words, err := h.svc.ListSensitiveWords(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &moderationv1.ListSensitiveWordsResponse{}
	for _, w := range words {
		resp.Words = append(resp.Words, &moderationv1.SensitiveWord{Id: w.ID, Word: w.Word, Level: w.Level})
	}
	return resp, nil
}

// ListTasks 列出审核任务。
func (h *Handler) ListTasks(ctx context.Context, req *moderationv1.ListTasksRequest) (*moderationv1.ListTasksResponse, error) {
	tasks, err := h.svc.ListTasks(ctx, int(req.GetLimit()))
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &moderationv1.ListTasksResponse{}
	for _, t := range tasks {
		resp.Tasks = append(resp.Tasks, &moderationv1.ModerationTask{
			Id: t.ID, AggregateType: t.AggregateType, AggregateId: t.AggregateID,
			Content: t.Content, Status: t.Status,
		})
	}
	return resp, nil
}

// Approve 通过审核。
func (h *Handler) Approve(ctx context.Context, req *moderationv1.ApproveRequest) (*moderationv1.ApproveResponse, error) {
	if err := h.svc.Approve(ctx, req.GetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &moderationv1.ApproveResponse{}, nil
}

// Reject 驳回审核。
func (h *Handler) Reject(ctx context.Context, req *moderationv1.RejectRequest) (*moderationv1.RejectResponse, error) {
	if err := h.svc.Reject(ctx, req.GetId(), req.GetReason()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &moderationv1.RejectResponse{}, nil
}

// CreateReport 创建举报。
func (h *Handler) CreateReport(ctx context.Context, req *moderationv1.CreateReportRequest) (*moderationv1.CreateReportResponse, error) {
	reporterID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.CreateReport(ctx, reporterID, req.GetTargetType(), req.GetTargetId(), req.GetReason()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &moderationv1.CreateReportResponse{}, nil
}

// ListReports 列出举报。
func (h *Handler) ListReports(ctx context.Context, req *moderationv1.ListReportsRequest) (*moderationv1.ListReportsResponse, error) {
	reports, err := h.svc.ListReports(ctx, int(req.GetLimit()))
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &moderationv1.ListReportsResponse{}
	for _, r := range reports {
		resp.Reports = append(resp.Reports, &moderationv1.Report{
			Id: r.ID, ReporterId: r.ReporterID, TargetType: r.TargetType,
			TargetId: r.TargetID, Reason: r.Reason, Status: r.Status,
		})
	}
	return resp, nil
}

// HandleReport 处理举报。
func (h *Handler) HandleReport(ctx context.Context, req *moderationv1.HandleReportRequest) (*moderationv1.HandleReportResponse, error) {
	if err := h.svc.HandleReport(ctx, req.GetId(), req.GetAction()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &moderationv1.HandleReportResponse{}, nil
}
