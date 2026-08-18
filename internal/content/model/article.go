// Package model 定义 content-service 的领域模型（Article 与 Question 独立建模，plan.md §6.3）。
package model

import (
	"fmt"
	"time"
)

// Status 是文章状态（plan.md §33 Phase 2）。
type Status string

const (
	StatusDraft         Status = "draft"
	StatusPendingReview Status = "pending_review"
	StatusPublished     Status = "published"
	StatusRejected      Status = "rejected"
	StatusHidden        Status = "hidden"
	StatusDeleted       Status = "deleted"
)

// Article 是文章实体，Markdown 是核心内容表达方式（spec.md §4.2）。
type Article struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	AuthorID  uint64    `gorm:"index;not null" json:"author_id"`
	Title     string    `gorm:"size:256;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"` // Markdown 原文
	Status    Status    `gorm:"size:32;not null;default:draft" json:"status"`
	Version   int       `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// transitions 是文章状态机的合法流转表。
var transitions = map[Status]map[Status]bool{
	StatusDraft: {
		StatusPendingReview: true, // 提交审核
		StatusDeleted:       true, // 作者删除草稿
	},
	StatusPendingReview: {
		StatusPublished: true, // 审核通过
		StatusRejected:  true, // 审核驳回
		StatusDeleted:   true, // 作者撤回
	},
	StatusPublished: {
		StatusHidden:  true, // 管理员隐藏
		StatusDeleted: true, // 删除
	},
	StatusRejected: {
		StatusDraft:         true, // 打回后编辑
		StatusPendingReview: true, // 重新提交
		StatusDeleted:       true,
	},
	StatusHidden: {
		StatusPublished: true, // 恢复
		StatusDeleted:   true,
	},
	StatusDeleted: {}, // 终态
}

// CanTransition 判断 from → to 是否合法。
func CanTransition(from, to Status) bool {
	return transitions[from][to]
}

// Transition 执行状态流转，非法流转返回错误。
func (a *Article) Transition(to Status) error {
	if !CanTransition(a.Status, to) {
		return fmt.Errorf("非法状态流转: %s → %s", a.Status, to)
	}
	a.Status = to
	return nil
}
