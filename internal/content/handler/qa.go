package handler

import (
	"context"

	contentv1 "github.com/likeslep/community/api/gen/content/v1"
	"github.com/likeslep/community/internal/content/service"
	"github.com/likeslep/community/pkg/grpcx"
)

// QuestionHandler 实现 Question gRPC 服务。
type QuestionHandler struct {
	contentv1.UnimplementedQuestionServiceServer
	svc *service.QAService
}

// NewQuestionHandler 构造。
func NewQuestionHandler(svc *service.QAService) *QuestionHandler {
	return &QuestionHandler{svc: svc}
}

// CreateQuestion 提问。
func (h *QuestionHandler) CreateQuestion(ctx context.Context, req *contentv1.CreateQuestionRequest) (*contentv1.CreateQuestionResponse, error) {
	authorID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	q, err := h.svc.CreateQuestion(ctx, authorID, req.GetTitle(), req.GetContent(), req.GetTags())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &contentv1.CreateQuestionResponse{Id: q.ID, Status: string(q.Status)}, nil
}

// GetQuestion 查询问题。
func (h *QuestionHandler) GetQuestion(ctx context.Context, req *contentv1.GetQuestionRequest) (*contentv1.GetQuestionResponse, error) {
	q, err := h.svc.GetQuestion(ctx, req.GetId())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	var acceptedID uint64
	if q.AcceptedAnswerID != nil {
		acceptedID = *q.AcceptedAnswerID
	}
	return &contentv1.GetQuestionResponse{
		Id: q.ID, AuthorId: q.AuthorID, Title: q.Title, Content: q.Content,
		Status: string(q.Status), AcceptedAnswerId: acceptedID,
	}, nil
}

// CloseQuestion 关闭问题（仅提问者）。
func (h *QuestionHandler) CloseQuestion(ctx context.Context, req *contentv1.CloseQuestionRequest) (*contentv1.CloseQuestionResponse, error) {
	authorID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.CloseQuestion(ctx, authorID, req.GetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &contentv1.CloseQuestionResponse{}, nil
}

// CreateAnswer 回答问题。
func (h *QuestionHandler) CreateAnswer(ctx context.Context, req *contentv1.CreateAnswerRequest) (*contentv1.CreateAnswerResponse, error) {
	authorID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	a, err := h.svc.CreateAnswer(ctx, authorID, req.GetQuestionId(), req.GetContent())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &contentv1.CreateAnswerResponse{Id: a.ID}, nil
}

// AcceptAnswer 采纳回答（仅提问者）。
func (h *QuestionHandler) AcceptAnswer(ctx context.Context, req *contentv1.AcceptAnswerRequest) (*contentv1.AcceptAnswerResponse, error) {
	authorID, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.AcceptAnswer(ctx, authorID, req.GetQuestionId(), req.GetAnswerId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &contentv1.AcceptAnswerResponse{}, nil
}

// ListQuestions 分页查询问题。
func (h *QuestionHandler) ListQuestions(ctx context.Context, req *contentv1.ListQuestionsRequest) (*contentv1.ListQuestionsResponse, error) {
	questions, err := h.svc.ListQuestions(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &contentv1.ListQuestionsResponse{}
	for _, q := range questions {
		resp.Questions = append(resp.Questions, &contentv1.QuestionBrief{
			Id: q.ID, Title: q.Title, AuthorId: q.AuthorID, Status: string(q.Status),
		})
	}
	return resp, nil
}
