// Package model 定义 notification-service 的领域模型。
package model

import "time"

// Notification 是站内通知。
type Notification struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	UserID     uint64    `gorm:"index;not null" json:"user_id"` // 接收者
	ActorID    uint64    `json:"actor_id"`                      // 触发者
	Type       string    `gorm:"size:64;not null" json:"type"`  // like/comment/follow/answer_accepted/moderation
	TargetType string    `gorm:"size:32" json:"target_type"`
	TargetID   uint64    `json:"target_id"`
	Content    string    `gorm:"size:512" json:"content"`
	Read       bool      `gorm:"not null;default:false" json:"read"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
