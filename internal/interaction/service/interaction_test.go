package service

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"

	"github.com/likeslep/community/internal/interaction/model"
	"github.com/likeslep/community/pkg/kafkax"
	"github.com/likeslep/community/pkg/redisx"
)

// fakeRepo 是 Repository 的内存实现，仅用于测试。
type fakeRepo struct {
	likeCount int
}

func (f *fakeRepo) CreateComment(_ context.Context, _ *model.Comment, _ func(*model.Comment) (kafkax.Envelope, error)) error {
	return nil
}
func (f *fakeRepo) ListComments(_ context.Context, _ string, _ uint64) ([]model.Comment, error) {
	return nil, nil
}
func (f *fakeRepo) Like(_ context.Context, _ uint64, _ string, _ uint64, _ func(*model.Like) (kafkax.Envelope, error)) error {
	f.likeCount++
	return nil
}
func (f *fakeRepo) Unlike(_ context.Context, _ uint64, _ string, _ uint64) error { return nil }
func (f *fakeRepo) Collect(_ context.Context, _ uint64, _ string, _ uint64, _ func(*model.Collection) (kafkax.Envelope, error)) error {
	return nil
}
func (f *fakeRepo) Uncollect(_ context.Context, _ uint64, _ string, _ uint64) error { return nil }
func (f *fakeRepo) GetCounter(_ context.Context, _ string, _ uint64) (*model.Counter, error) {
	return &model.Counter{}, nil
}
func (f *fakeRepo) DeleteComment(_ context.Context, _ uint64) error { return nil }

func newTestService(t *testing.T) *InteractionService {
	mr := miniredis.RunT(t)
	return NewInteractionService(&fakeRepo{}, Config{
		Producer: "interaction-service",
		Redis:    redisx.New(redisx.Config{Addr: mr.Addr()}),
	})
}

func TestLikeInvalidTarget(t *testing.T) {
	svc := newTestService(t)
	if err := svc.Like(context.Background(), 1, "invalid_type", 1); err == nil {
		t.Fatal("非法目标类型应报错")
	}
}

func TestLikeValid(t *testing.T) {
	svc := newTestService(t)
	if err := svc.Like(context.Background(), 1, "article", 1); err != nil {
		t.Fatalf("Like() err = %v", err)
	}
}

func TestView(t *testing.T) {
	svc := newTestService(t)
	ctx := context.Background()
	c1, err := svc.View(ctx, "article", 1)
	if err != nil {
		t.Fatalf("View() err = %v", err)
	}
	c2, err := svc.View(ctx, "article", 1)
	if err != nil {
		t.Fatalf("View() err = %v", err)
	}
	if c2 != c1+1 {
		t.Fatalf("View 计数应递增：c1=%d c2=%d", c1, c2)
	}
}
