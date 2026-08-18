package model

import "time"

// Answer 是问题回答实体。
type Answer struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	QuestionID uint64    `gorm:"index;not null" json:"question_id"`
	AuthorID   uint64    `gorm:"index;not null" json:"author_id"`
	Content    string    `gorm:"type:text" json:"content"`
	Accepted   bool      `gorm:"not null;default:false" json:"accepted"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// QuestionTag 是问题-标签多对多关联。
type QuestionTag struct {
	QuestionID uint64 `gorm:"primaryKey"`
	TagID      uint64 `gorm:"primaryKey"`
}
