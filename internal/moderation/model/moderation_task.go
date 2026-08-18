package model

import "time"

// 审核任务状态。
const (
	TaskPending  = "pending"
	TaskApproved = "approved"
	TaskRejected = "rejected"
)

// ModerationTask 是审核任务（消费 content 提交事件生成）。
type ModerationTask struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	AggregateType string    `gorm:"size:64;not null;index" json:"aggregate_type"`
	AggregateID   string    `gorm:"size:64;not null;index" json:"aggregate_id"`
	Content       string    `gorm:"type:text" json:"content"`
	Status        string    `gorm:"size:32;not null;default:pending" json:"status"`
	Reason        string    `gorm:"size:256" json:"reason"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
