// Package model 定义 file-service 的领域模型。
package model

import "time"

// File 是文件元数据（业务库只存元数据，不存二进制，plan.md §26）。
type File struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"index" json:"user_id"`
	Name      string    `gorm:"size:256;not null" json:"name"` // 原始文件名
	Path      string    `gorm:"size:256;not null" json:"-"`    // 存储 key（不对外暴露）
	Type      string    `gorm:"size:128" json:"type"`          // MIME 类型
	Size      int64     `gorm:"not null" json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
