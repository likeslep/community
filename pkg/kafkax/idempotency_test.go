package kafkax

import (
	"context"
	"testing"
)

func TestMemProcessedStore(t *testing.T) {
	ctx := context.Background()
	s := NewMemProcessedStore()
	if s.Exists(ctx, "evt-1") {
		t.Fatal("初始不应存在 evt-1")
	}
	if err := s.Mark(ctx, "evt-1"); err != nil {
		t.Fatalf("Mark() err = %v", err)
	}
	if !s.Exists(ctx, "evt-1") {
		t.Fatal("Mark 后应存在 evt-1")
	}
}
