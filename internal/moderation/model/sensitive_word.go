// Package model 定义 moderation-service 的领域模型。
package model

import "time"

// 敏感词级别：block（直接拦截）或 review（进入人工审核）。
const (
	LevelBlock  = "block"
	LevelReview = "review"
)

// SensitiveWord 是敏感词实体。
type SensitiveWord struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Word      string    `gorm:"uniqueIndex;size:128;not null" json:"word"`
	Level     string    `gorm:"size:16;not null;default:review" json:"level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
