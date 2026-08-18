// Package kafkax 提供 Kafka 事件基础设施：Envelope、Producer/Consumer、幂等与重试/DLQ。
package kafkax

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Envelope 是统一事件信封（plan.md §17）。
type Envelope struct {
	EventID          string          `json:"event_id"`
	EventType        string          `json:"event_type"`
	Version          int             `json:"version"`
	OccurredAt       time.Time       `json:"occurred_at"`
	Producer         string          `json:"producer"`
	AggregateType    string          `json:"aggregate_type"`
	AggregateID      string          `json:"aggregate_id"`
	AggregateVersion int             `json:"aggregate_version"`
	TraceID          string          `json:"trace_id"`
	Payload          json.RawMessage `json:"payload"`
}

// NewEnvelope 构造事件信封。eventVersion 为事件 schema 版本，aggregateVersion 为聚合版本。
func NewEnvelope(eventType, producer, aggregateType, aggregateID string, aggregateVersion int, payload any) (Envelope, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return Envelope{}, err
	}
	return Envelope{
		EventID:          uuid.NewString(),
		EventType:        eventType,
		Version:          1,
		OccurredAt:       time.Now(),
		Producer:         producer,
		AggregateType:    aggregateType,
		AggregateID:      aggregateID,
		AggregateVersion: aggregateVersion,
		Payload:          raw,
	}, nil
}

// PartitionKey 返回 Kafka 分区键，保证同一 aggregate 的事件有序（plan.md §17）。
func (e Envelope) PartitionKey() string { return e.AggregateType + ":" + e.AggregateID }
