package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/likeslep/community/internal/notification/model"
	"github.com/likeslep/community/pkg/kafkax"
)

// fakeRepo 是 Repository 的内存实现，仅用于测试。
type fakeRepo struct {
	notifications []*model.Notification
}

func (f *fakeRepo) CreateNotification(_ context.Context, n *model.Notification) error {
	f.notifications = append(f.notifications, n)
	return nil
}
func (f *fakeRepo) ListNotifications(_ context.Context, _ uint64, _ int) ([]model.Notification, error) {
	return nil, nil
}
func (f *fakeRepo) UnreadCount(_ context.Context, _ uint64) (int64, error) { return 0, nil }
func (f *fakeRepo) MarkRead(_ context.Context, _, _ uint64) error          { return nil }
func (f *fakeRepo) MarkAllRead(_ context.Context, _ uint64) error          { return nil }

func TestHandleEventFollow(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewNotificationService(repo)

	env := kafkax.Envelope{
		EventType: kafkax.EventUserFollowed,
		Payload:   json.RawMessage(`{"follower_id":1,"followee_id":2}`),
	}
	if err := svc.HandleEvent(context.Background(), env); err != nil {
		t.Fatalf("HandleEvent() err = %v", err)
	}
	if len(repo.notifications) != 1 {
		t.Fatalf("应生成 1 条通知，实际 %d", len(repo.notifications))
	}
	if repo.notifications[0].UserID != 2 {
		t.Fatalf("通知接收者应为 2，实际 %d", repo.notifications[0].UserID)
	}
}

func TestHandleEventIgnoresUnknown(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewNotificationService(repo)
	env := kafkax.Envelope{EventType: "unknown.event", Payload: json.RawMessage(`{}`)}
	if err := svc.HandleEvent(context.Background(), env); err != nil {
		t.Fatalf("HandleEvent() err = %v", err)
	}
	if len(repo.notifications) != 0 {
		t.Fatal("未知事件不应生成通知")
	}
}
