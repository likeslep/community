// Package model 定义 user-service 的领域模型（GORM 实体）。
package model

import "time"

// User 是用户实体。
type User struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Email        string    `gorm:"uniqueIndex;size:128;not null;default:''" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:32;not null;default:author" json:"role"`
	Status       string    `gorm:"size:32;not null;default:active" json:"status"`
	Bio          string    `gorm:"size:512" json:"bio"`
	AvatarFileID string    `gorm:"size:64" json:"avatar_file_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// 用户角色（plan.md §12）。
const (
	RoleAuthor    = "author"
	RoleModerator = "moderator"
	RoleAdmin     = "admin"
)

// 用户状态。
const (
	StatusActive = "active"
	StatusBanned = "banned"
)
