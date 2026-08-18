// Package model 定义 interaction-service 的领域模型。
package model

import "time"

// Comment 是评论（多态目标：question/answer/article）。
type Comment struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	UserID     uint64    `gorm:"index;not null" json:"user_id"`
	TargetType string    `gorm:"size:32;not null;index" json:"target_type"`
	TargetID   uint64    `gorm:"index;not null" json:"target_id"`
	Content    string    `gorm:"type:text" json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Like 是点赞（用户+目标唯一，幂等）。
type Like struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	UserID     uint64    `gorm:"uniqueIndex:uk_like_user_target;not null" json:"user_id"`
	TargetType string    `gorm:"uniqueIndex:uk_like_user_target;size:32;not null" json:"target_type"`
	TargetID   uint64    `gorm:"uniqueIndex:uk_like_user_target;not null" json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// Collection 是收藏（用户+目标唯一，幂等）。
type Collection struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	UserID     uint64    `gorm:"uniqueIndex:uk_collect_user_target;not null" json:"user_id"`
	TargetType string    `gorm:"uniqueIndex:uk_collect_user_target;size:32;not null" json:"target_type"`
	TargetID   uint64    `gorm:"uniqueIndex:uk_collect_user_target;not null" json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// Counter 是互动计数器（MySQL 为 Source of Truth）。
type Counter struct {
	TargetType   string    `gorm:"primaryKey;size:32" json:"target_type"`
	TargetID     uint64    `gorm:"primaryKey" json:"target_id"`
	LikeCount    int64     `gorm:"not null;default:0" json:"like_count"`
	CommentCount int64     `gorm:"not null;default:0" json:"comment_count"`
	CollectCount int64     `gorm:"not null;default:0" json:"collect_count"`
	ViewCount    int64     `gorm:"not null;default:0" json:"view_count"`
	UpdatedAt    time.Time `json:"updated_at"`
}
