package service

import (
	"context"
	"strconv"

	"github.com/likeslep/community/internal/moderation/model"
	"github.com/likeslep/community/pkg/kafkax"
)

// Config 是 moderation-service 的业务配置。
type Config struct {
	Producer string
}

// ModerationService 是审核业务逻辑层。
type ModerationService struct {
	repo Repository
	cfg  Config
}

// NewModerationService 构造。
func NewModerationService(repo Repository, cfg Config) *ModerationService {
	return &ModerationService{repo: repo, cfg: cfg}
}

// CheckText 敏感词检测，返回处理建议与命中词。
func (s *ModerationService) CheckText(ctx context.Context, text string) (string, []string, error) {
	words, err := s.repo.ListSensitiveWords(ctx)
	if err != nil {
		return "", nil, err
	}
	hits := NewDetector(words).Check(text)
	matched := make([]string, 0, len(hits))
	for _, h := range hits {
		matched = append(matched, h.Word)
	}
	return Classify(hits), matched, nil
}

// CreateSensitiveWord 添加敏感词。
func (s *ModerationService) CreateSensitiveWord(ctx context.Context, word, level string) error {
	if level == "" {
		level = model.LevelReview
	}
	return s.repo.CreateSensitiveWord(ctx, &model.SensitiveWord{Word: word, Level: level})
}

// ListSensitiveWords 查询敏感词列表。
func (s *ModerationService) ListSensitiveWords(ctx context.Context) ([]model.SensitiveWord, error) {
	return s.repo.ListSensitiveWords(ctx)
}

// CreateTask 从提交事件创建审核任务（幂等由消费者保证）。
func (s *ModerationService) CreateTask(ctx context.Context, aggregateType, aggregateID, content string) error {
	t := &model.ModerationTask{
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		Content:       content,
		Status:        model.TaskPending,
	}
	return s.repo.CreateTask(ctx, t)
}

// ListTasks 列出审核任务。
func (s *ModerationService) ListTasks(ctx context.Context, limit int) ([]model.ModerationTask, error) {
	return s.repo.ListTasks(ctx, limit)
}

// Approve 通过审核，发布 moderation.approved 事件。
func (s *ModerationService) Approve(ctx context.Context, id uint64) error {
	t, err := s.repo.FindTask(ctx, id)
	if err != nil {
		return err
	}
	if t.Status != model.TaskPending {
		return ErrIllegalState
	}
	t.Status = model.TaskApproved
	env, err := s.moderationEvent(t, kafkax.EventModerationApproved)
	if err != nil {
		return err
	}
	return s.repo.ApproveTask(ctx, t, env)
}

// Reject 驳回审核，发布 moderation.rejected 事件。
func (s *ModerationService) Reject(ctx context.Context, id uint64, reason string) error {
	t, err := s.repo.FindTask(ctx, id)
	if err != nil {
		return err
	}
	if t.Status != model.TaskPending {
		return ErrIllegalState
	}
	t.Status = model.TaskRejected
	t.Reason = reason
	env, err := s.moderationEvent(t, kafkax.EventModerationRejected)
	if err != nil {
		return err
	}
	return s.repo.RejectTask(ctx, t, env)
}

// CreateReport 创建举报。
func (s *ModerationService) CreateReport(ctx context.Context, reporterID uint64, targetType string, targetID uint64, reason string) error {
	r := &model.Report{
		ReporterID: reporterID,
		TargetType: targetType,
		TargetID:   targetID,
		Reason:     reason,
		Status:     model.ReportPending,
	}
	return s.repo.CreateReport(ctx, r)
}

// ListReports 列出举报。
func (s *ModerationService) ListReports(ctx context.Context, limit int) ([]model.Report, error) {
	return s.repo.ListReports(ctx, limit)
}

// HandleReport 处理举报（通过/驳回）。
func (s *ModerationService) HandleReport(ctx context.Context, id uint64, action string) error {
	r, err := s.repo.FindReport(ctx, id)
	if err != nil {
		return err
	}
	switch action {
	case model.ReportApproved, model.ReportRejected:
		r.Status = action
	default:
		return ErrIllegalState
	}
	return s.repo.UpdateReport(ctx, r)
}

func (s *ModerationService) moderationEvent(t *model.ModerationTask, eventType string) (kafkax.Envelope, error) {
	return kafkax.NewEnvelope(eventType, s.cfg.Producer, t.AggregateType, t.AggregateID, 1,
		moderationPayload{TaskID: t.ID, AggregateType: t.AggregateType, AggregateID: t.AggregateID, Reason: t.Reason})
}

type moderationPayload struct {
	TaskID        uint64 `json:"task_id"`
	AggregateType string `json:"aggregate_type"`
	AggregateID   string `json:"aggregate_id"`
	Reason        string `json:"reason"`
}

// ParseAggregateID 辅助解析事件 aggregate_id 为数字（供内容服务使用）。
func ParseAggregateID(s string) uint64 {
	id, _ := strconv.ParseUint(s, 10, 64)
	return id
}
