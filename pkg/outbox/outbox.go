// Package outbox 提供事务性 Outbox：业务数据与事件在同一事务内写入（plan.md §19）。
package outbox

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/likeslep/community/pkg/kafkax"
)

// Event 是 outbox_events 表的一行。
type Event struct {
	ID               string     `gorm:"column:id;primaryKey"`
	EventType        string     `gorm:"column:event_type"`
	AggregateType    string     `gorm:"column:aggregate_type"`
	AggregateID      string     `gorm:"column:aggregate_id"`
	AggregateVersion int        `gorm:"column:aggregate_version"`
	Payload          string     `gorm:"column:payload"`
	TraceID          string     `gorm:"column:trace_id"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	PublishedAt      *time.Time `gorm:"column:published_at"`
}

// TableName 指定表名。
func (Event) TableName() string { return "outbox_events" }

// Store 操作 outbox 表。
type Store struct {
	db *gorm.DB
}

// NewStore 构造 Store。
func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

// Insert 在给定事务内写入一条事件。
func (s *Store) Insert(ctx context.Context, tx *gorm.DB, env kafkax.Envelope) error {
	e := Event{
		ID:               env.EventID,
		EventType:        env.EventType,
		AggregateType:    env.AggregateType,
		AggregateID:      env.AggregateID,
		AggregateVersion: env.AggregateVersion,
		Payload:          string(env.Payload),
		TraceID:          env.TraceID,
		CreatedAt:        env.OccurredAt,
	}
	return tx.WithContext(ctx).Create(&e).Error
}

// Pending 查询未发布的最近 limit 条事件（按创建时间升序）。
func (s *Store) Pending(ctx context.Context, limit int) ([]Event, error) {
	var events []Event
	err := s.db.WithContext(ctx).
		Where("published_at IS NULL").
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

// MarkPublished 标记事件已发布。
func (s *Store) MarkPublished(ctx context.Context, id string) error {
	now := time.Now()
	return s.db.WithContext(ctx).
		Model(&Event{}).
		Where("id = ?", id).
		Update("published_at", &now).Error
}
