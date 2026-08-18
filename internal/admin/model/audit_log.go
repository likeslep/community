// Package model 定义 admin-service 的领域模型。
package model

import "time"

// AuditLog 是管理员操作审计日志（plan.md §6.11 Audit Log）。
type AuditLog struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	AdminID    uint64    `json:"admin_id"`
	Action     string    `gorm:"size:128;not null" json:"action"`
	TargetType string    `gorm:"size:64" json:"target_type"`
	TargetID   uint64    `json:"target_id"`
	Detail     string    `gorm:"size:512" json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}
