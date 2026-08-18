package model

import "time"

// 举报状态。
const (
	ReportPending  = "pending"
	ReportApproved = "approved"
	ReportRejected = "rejected"
)

// Report 是用户举报。
type Report struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	ReporterID uint64    `gorm:"index;not null" json:"reporter_id"`
	TargetType string    `gorm:"size:32;not null;index" json:"target_type"`
	TargetID   uint64    `gorm:"index;not null" json:"target_id"`
	Reason     string    `gorm:"size:256" json:"reason"`
	Status     string    `gorm:"size:32;not null;default:pending" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
