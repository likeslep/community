package service

import (
	"context"

	"github.com/likeslep/community/internal/content/model"
	"github.com/likeslep/community/pkg/kafkax"
)

// BuildQuestion 在问题 ID 已知后构造事件信封。
type BuildQuestion func(q *model.Question) (kafkax.Envelope, error)

// BuildAnswer 在回答 ID 已知后构造事件信封。
type BuildAnswer func(a *model.Answer) (kafkax.Envelope, error)

// QARepository 是问答的持久化接口（plan.md §9）。
type QARepository interface {
	// CreateQuestion 创建问题（绑定标签 + outbox）在同一事务内。
	CreateQuestion(ctx context.Context, q *model.Question, tagNames []string, build BuildQuestion) error
	// FindQuestion 查询问题。
	FindQuestion(ctx context.Context, id uint64) (*model.Question, error)
	// UpdateQuestion 更新问题状态。
	UpdateQuestion(ctx context.Context, q *model.Question) error
	// CreateAnswer 创建回答 + outbox 在同一事务内。
	CreateAnswer(ctx context.Context, a *model.Answer, build BuildAnswer) error
	// FindAnswer 查询回答。
	FindAnswer(ctx context.Context, id uint64) (*model.Answer, error)
	// AcceptAnswer 设置采纳回答（清除旧采纳 + 更新问题 + outbox）在同一事务内。
	AcceptAnswer(ctx context.Context, questionID, answerID uint64, env kafkax.Envelope) error
	// ListQuestions 分页查询问题。
	ListQuestions(ctx context.Context, limit, offset int) ([]model.Question, error)
}
