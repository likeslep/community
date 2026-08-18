package kafkax

import (
	"testing"
)

func TestNewEnvelope(t *testing.T) {
	env, err := NewEnvelope(EventArticleCreated, "content-service", "article", "10086", 3, map[string]string{"title": "hello"})
	if err != nil {
		t.Fatalf("NewEnvelope() err = %v", err)
	}
	if env.EventID == "" {
		t.Fatal("EventID 不应为空")
	}
	if env.EventType != EventArticleCreated {
		t.Fatalf("EventType = %q", env.EventType)
	}
	if env.Version != 1 {
		t.Fatalf("事件 schema 版本应为 1，got %d", env.Version)
	}
	if env.AggregateVersion != 3 {
		t.Fatalf("AggregateVersion = %d, want 3", env.AggregateVersion)
	}
	if env.OccurredAt.IsZero() {
		t.Fatal("OccurredAt 不应为零值")
	}
	if string(env.Payload) != `{"title":"hello"}` {
		t.Fatalf("Payload = %s", env.Payload)
	}
}

func TestPartitionKey(t *testing.T) {
	env, _ := NewEnvelope(EventArticleCreated, "content-service", "article", "10086", 1, nil)
	if got := env.PartitionKey(); got != "article:10086" {
		t.Fatalf("PartitionKey() = %q, want %q", got, "article:10086")
	}
}
