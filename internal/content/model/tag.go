package model

import "time"

// Tag 是标签实体（spec.md §5，系统重要基础实体）。
type Tag struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:64;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ArticleTag 是文章-标签多对多关联。
type ArticleTag struct {
	ArticleID uint64 `gorm:"primaryKey"`
	TagID     uint64 `gorm:"primaryKey"`
}
