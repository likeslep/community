package model

import "time"

// ArticleVersion 是文章编辑的历史版本快照。
type ArticleVersion struct {
	ID        uint64 `gorm:"primaryKey"`
	ArticleID uint64 `gorm:"index"`
	Title     string `gorm:"size:256;not null"`
	Content   string `gorm:"type:text"`
	Version   int    `gorm:"not null"`
	CreatedAt time.Time
}
