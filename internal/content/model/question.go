package model

import (
	"fmt"
	"time"
)

// QuestionStatus 是问题状态。
type QuestionStatus string

const (
	QuestionOpen    QuestionStatus = "open"
	QuestionClosed  QuestionStatus = "closed"
	QuestionDeleted QuestionStatus = "deleted"
)

// Question 是问题实体（与 Article 独立建模，plan.md §6.3）。
type Question struct {
	ID               uint64         `gorm:"primaryKey" json:"id"`
	AuthorID         uint64         `gorm:"index;not null" json:"author_id"`
	Title            string         `gorm:"size:256;not null" json:"title"`
	Content          string         `gorm:"type:text" json:"content"`
	Status           QuestionStatus `gorm:"size:32;not null;default:open" json:"status"`
	AcceptedAnswerID *uint64        `gorm:"default:null" json:"accepted_answer_id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// questionTransitions 是问题状态机的合法流转表。
var questionTransitions = map[QuestionStatus]map[QuestionStatus]bool{
	QuestionOpen: {
		QuestionClosed:  true,
		QuestionDeleted: true,
	},
	QuestionClosed: {
		QuestionOpen:    true, // 重新打开
		QuestionDeleted: true,
	},
	QuestionDeleted: {}, // 终态
}

// CanTransition 判断问题状态 from → to 是否合法。
func CanQuestionTransition(from, to QuestionStatus) bool {
	return questionTransitions[from][to]
}

// Transition 执行状态流转，非法流转返回错误。
func (q *Question) Transition(to QuestionStatus) error {
	if !CanQuestionTransition(q.Status, to) {
		return fmt.Errorf("非法状态流转: %s → %s", q.Status, to)
	}
	q.Status = to
	return nil
}
