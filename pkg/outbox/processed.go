package outbox

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProcessedEvent 是已处理事件记录，用于幂等消费（plan.md §20）。
type ProcessedEvent struct {
	EventID     string    `gorm:"column:event_id;primaryKey"`
	ProcessedAt time.Time `gorm:"column:processed_at"`
}

// TableName 指定表名。
func (ProcessedEvent) TableName() string { return "processed_events" }

// ProcessedStore 是基于数据库的幂等存储，实现 kafkax.ProcessedStore 接口。
type ProcessedStore struct {
	db *gorm.DB
}

// NewProcessedStore 构造。
func NewProcessedStore(db *gorm.DB) *ProcessedStore { return &ProcessedStore{db: db} }

// Exists 判断事件是否已处理。
func (s *ProcessedStore) Exists(ctx context.Context, eventID string) bool {
	var count int64
	_ = s.db.WithContext(ctx).Model(&ProcessedEvent{}).Where("event_id = ?", eventID).Count(&count)
	return count > 0
}

// Mark 记录事件已处理（重复插入被忽略）。
func (s *ProcessedStore) Mark(ctx context.Context, eventID string) error {
	return s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&ProcessedEvent{EventID: eventID, ProcessedAt: time.Now()}).Error
}
