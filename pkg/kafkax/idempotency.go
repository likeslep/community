package kafkax

import (
	"context"
	"sync"
)

// ProcessedStore 记录已处理事件 ID，用于幂等消费（plan.md §20）。
// 生产环境应由数据库唯一约束实现；此处提供接口 + 内存实现。
type ProcessedStore interface {
	Exists(ctx context.Context, eventID string) bool
	Mark(ctx context.Context, eventID string) error
}

// MemProcessedStore 是内存实现，仅用于测试。
// 并发安全：用 mutex 保护底层 map，避免并发读写竞态。
type MemProcessedStore struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

// NewMemProcessedStore 构造内存幂等存储。
func NewMemProcessedStore() *MemProcessedStore {
	return &MemProcessedStore{seen: make(map[string]struct{})}
}

func (m *MemProcessedStore) Exists(_ context.Context, eventID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.seen[eventID]
	return ok
}

func (m *MemProcessedStore) Mark(_ context.Context, eventID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.seen[eventID] = struct{}{}
	return nil
}
