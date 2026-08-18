package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/likeslep/community/internal/content/model"
	"github.com/likeslep/community/internal/content/service"
	"github.com/likeslep/community/pkg/kafkax"
)

// CreateQuestion 创建问题（绑定标签 + outbox）在同一事务内。
func (r *Gorm) CreateQuestion(ctx context.Context, q *model.Question, tagNames []string, build service.BuildQuestion) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(q).Error; err != nil {
			return err
		}
		if err := bindQuestionTags(ctx, tx, q.ID, tagNames); err != nil {
			return err
		}
		env, err := build(q)
		if err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

// FindQuestion 查询问题。
func (r *Gorm) FindQuestion(ctx context.Context, id uint64) (*model.Question, error) {
	var q model.Question
	if err := r.db.WithContext(ctx).First(&q, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrQuestionNotFound
		}
		return nil, err
	}
	return &q, nil
}

// UpdateQuestion 更新问题状态。
func (r *Gorm) UpdateQuestion(ctx context.Context, q *model.Question) error {
	return r.db.WithContext(ctx).Save(q).Error
}

// ListQuestions 分页查询问题。
func (r *Gorm) ListQuestions(ctx context.Context, limit, offset int) ([]model.Question, error) {
	var questions []model.Question
	err := r.db.WithContext(ctx).Order("id DESC").Limit(limit).Offset(offset).Find(&questions).Error
	return questions, err
}

// ListAnswers 查询某问题的回答列表（采纳的排最前）。
func (r *Gorm) ListAnswers(ctx context.Context, questionID uint64) ([]model.Answer, error) {
	var answers []model.Answer
	err := r.db.WithContext(ctx).
		Where("question_id = ?", questionID).
		Order("accepted DESC, id ASC").
		Find(&answers).Error
	return answers, err
}

// CreateAnswer 创建回答 + outbox。
func (r *Gorm) CreateAnswer(ctx context.Context, a *model.Answer, build service.BuildAnswer) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(a).Error; err != nil {
			return err
		}
		env, err := build(a)
		if err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

// FindAnswer 查询回答。
func (r *Gorm) FindAnswer(ctx context.Context, id uint64) (*model.Answer, error) {
	var a model.Answer
	if err := r.db.WithContext(ctx).First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrAnswerNotFound
		}
		return nil, err
	}
	return &a, nil
}

// AcceptAnswer 设置采纳回答（清除旧采纳 + 更新问题 + outbox）在同一事务内。
func (r *Gorm) AcceptAnswer(ctx context.Context, questionID, answerID uint64, env kafkax.Envelope) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Answer{}).
			Where("question_id = ? AND accepted = ?", questionID, true).
			Update("accepted", false).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Answer{}).
			Where("id = ?", answerID).
			Update("accepted", true).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Question{}).
			Where("id = ?", questionID).
			Update("accepted_answer_id", answerID).Error; err != nil {
			return err
		}
		return r.outbox.Insert(ctx, tx, env)
	})
}

func bindQuestionTags(ctx context.Context, tx *gorm.DB, questionID uint64, names []string) error {
	for _, raw := range names {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		var tag model.Tag
		if err := tx.WithContext(ctx).Where("name = ?", name).FirstOrCreate(&tag, model.Tag{Name: name}).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Create(&model.QuestionTag{QuestionID: questionID, TagID: tag.ID}).Error; err != nil {
			return err
		}
	}
	return nil
}
