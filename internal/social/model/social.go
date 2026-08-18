// Package model 定义 social-service 的领域模型。
package model

import "time"

// Follow 是用户关注关系（follower 关注 followee）。
type Follow struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	FollowerID uint64    `gorm:"uniqueIndex:uk_follow;not null" json:"follower_id"`
	FolloweeID uint64    `gorm:"uniqueIndex:uk_follow;not null" json:"followee_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// TagFollow 是用户关注标签。
type TagFollow struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"uniqueIndex:uk_tag_follow;not null" json:"user_id"`
	TagID     uint64    `gorm:"uniqueIndex:uk_tag_follow;not null" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}
