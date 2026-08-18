package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/likeslep/community/internal/content/model"
	"github.com/likeslep/community/pkg/apperr"
	"github.com/likeslep/community/pkg/kafkax"
)

// QAService 是问答业务逻辑层。
type QAService struct {
	repo QARepository
	cfg  Config
}

// NewQAService 构造。
func NewQAService(repo QARepository, cfg Config) *QAService {
	return &QAService{repo: repo, cfg: cfg}
}

// CreateQuestion 提问。
func (s *QAService) CreateQuestion(ctx context.Context, authorID uint64, title, content string, tags []string) (*model.Question, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, apperr.New(errCodeInvalidInput, "标题不能为空", apperr.WithHTTP(400))
	}
	q := &model.Question{AuthorID: authorID, Title: title, Content: content, Status: model.QuestionOpen}
	err := s.repo.CreateQuestion(ctx, q, tags, func(q *model.Question) (kafkax.Envelope, error) {
		return kafkax.NewEnvelope(kafkax.EventQuestionCreated, s.cfg.Producer, "question",
			strconv.FormatUint(q.ID, 10), 1, questionPayload{QuestionID: q.ID, AuthorID: q.AuthorID, Title: q.Title})
	})
	if err != nil {
		return nil, err
	}
	return q, nil
}

// GetQuestion 查询问题。
func (s *QAService) GetQuestion(ctx context.Context, id uint64) (*model.Question, error) {
	return s.repo.FindQuestion(ctx, id)
}

// CloseQuestion 关闭问题（仅提问者）。
func (s *QAService) CloseQuestion(ctx context.Context, authorID, id uint64) error {
	q, err := s.repo.FindQuestion(ctx, id)
	if err != nil {
		return err
	}
	if q.AuthorID != authorID {
		return ErrForbidden
	}
	if err := q.Transition(model.QuestionClosed); err != nil {
		return ErrIllegalState
	}
	return s.repo.UpdateQuestion(ctx, q)
}

// CreateAnswer 回答问题。
func (s *QAService) CreateAnswer(ctx context.Context, authorID, questionID uint64, content string) (*model.Answer, error) {
	q, err := s.repo.FindQuestion(ctx, questionID)
	if err != nil {
		return nil, err
	}
	if q.Status != model.QuestionOpen {
		return nil, ErrIllegalState
	}
	if strings.TrimSpace(content) == "" {
		return nil, apperr.New(errCodeInvalidInput, "回答内容不能为空", apperr.WithHTTP(400))
	}
	a := &model.Answer{QuestionID: questionID, AuthorID: authorID, Content: content}
	err = s.repo.CreateAnswer(ctx, a, func(a *model.Answer) (kafkax.Envelope, error) {
		return kafkax.NewEnvelope(kafkax.EventAnswerCreated, s.cfg.Producer, "answer",
			strconv.FormatUint(a.ID, 10), 1, answerPayload{AnswerID: a.ID, QuestionID: a.QuestionID, AuthorID: a.AuthorID})
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

// AcceptAnswer 采纳回答（仅提问者）。
func (s *QAService) AcceptAnswer(ctx context.Context, authorID, questionID, answerID uint64) error {
	q, err := s.repo.FindQuestion(ctx, questionID)
	if err != nil {
		return err
	}
	if q.AuthorID != authorID {
		return ErrForbidden
	}
	a, err := s.repo.FindAnswer(ctx, answerID)
	if err != nil {
		return err
	}
	if a.QuestionID != questionID {
		return ErrAnswerMismatch
	}
	env, err := kafkax.NewEnvelope(kafkax.EventAnswerAccepted, s.cfg.Producer, "answer",
		strconv.FormatUint(answerID, 10), 1, answerPayload{AnswerID: a.ID, QuestionID: a.QuestionID, AuthorID: a.AuthorID})
	if err != nil {
		return err
	}
	return s.repo.AcceptAnswer(ctx, questionID, answerID, env)
}

type questionPayload struct {
	QuestionID uint64 `json:"question_id"`
	AuthorID   uint64 `json:"author_id"`
	Title      string `json:"title"`
}

type answerPayload struct {
	AnswerID   uint64 `json:"answer_id"`
	QuestionID uint64 `json:"question_id"`
	AuthorID   uint64 `json:"author_id"`
}
