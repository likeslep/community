package service

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/likeslep/community/internal/content/model"
)

// fakeRepo 是 Repository 的内存实现，仅用于测试。
// 并发安全：用 mutex 保护 map，避免并发读写竞态。
type fakeRepo struct {
	mu       sync.Mutex
	articles map[uint64]*model.Article
	nextID   uint64
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{articles: map[uint64]*model.Article{}, nextID: 1}
}

func (f *fakeRepo) CreateArticle(_ context.Context, a *model.Article, _ []string, build BuildEvent) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	a.ID = f.nextID
	f.nextID++
	f.articles[a.ID] = a
	_, err := build(a)
	return err
}

func (f *fakeRepo) UpdateArticle(_ context.Context, a *model.Article, build BuildEvent) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.articles[a.ID] = a
	_, err := build(a)
	return err
}

func (f *fakeRepo) FindArticle(_ context.Context, id uint64) (*model.Article, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if a, ok := f.articles[id]; ok {
		return a, nil
	}
	return nil, ErrArticleNotFound
}

func (f *fakeRepo) DeleteArticle(_ context.Context, id uint64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if a, ok := f.articles[id]; ok {
		a.Status = model.StatusDeleted
	}
	return nil
}

func (f *fakeRepo) SubmitArticle(_ context.Context, a *model.Article, build BuildEvent) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.articles[a.ID] = a
	_, err := build(a)
	return err
}

func (f *fakeRepo) SaveWithEvent(_ context.Context, a *model.Article, build BuildEvent) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.articles[a.ID] = a
	_, err := build(a)
	return err
}

func (f *fakeRepo) TagsByArticle(_ context.Context, _ uint64) ([]model.Tag, error) {
	return nil, nil
}

func (f *fakeRepo) ListTags(_ context.Context) ([]model.Tag, error) { return nil, nil }

func (f *fakeRepo) UpdateArticleStatus(_ context.Context, _ uint64, _ model.Status) error { return nil }

func newTestArticleService() *ArticleService {
	return NewArticleService(newFakeRepo(), Config{Producer: "content-service"})
}

func TestCreateArticle(t *testing.T) {
	svc := newTestArticleService()
	ctx := context.Background()

	a, err := svc.CreateArticle(ctx, 1, "我的文章", "# 标题", []string{"go", "kafka"})
	if err != nil {
		t.Fatalf("CreateArticle() err = %v", err)
	}
	if a.ID == 0 || a.Status != model.StatusDraft || a.Version != 1 {
		t.Fatalf("article = %+v", a)
	}

	if _, err := svc.CreateArticle(ctx, 1, "", "x", nil); err == nil {
		t.Fatal("空标题应报错")
	}
}

func TestUpdateArticleOwnership(t *testing.T) {
	svc := newTestArticleService()
	ctx := context.Background()
	a, _ := svc.CreateArticle(ctx, 1, "标题", "内容", nil)

	// 非作者更新应被拒
	if _, err := svc.UpdateArticle(ctx, 2, a.ID, "新标题", "新内容"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("期望 ErrForbidden，got %v", err)
	}
	// 作者更新成功且版本递增
	updated, err := svc.UpdateArticle(ctx, 1, a.ID, "新标题", "新内容")
	if err != nil {
		t.Fatalf("UpdateArticle() err = %v", err)
	}
	if updated.Version != 2 {
		t.Fatalf("版本应为 2，got %d", updated.Version)
	}
}

func TestSubmitArticleFlow(t *testing.T) {
	svc := newTestArticleService()
	ctx := context.Background()
	a, _ := svc.CreateArticle(ctx, 1, "标题", "内容", nil)

	if err := svc.SubmitArticle(ctx, 1, a.ID); err != nil {
		t.Fatalf("SubmitArticle() err = %v", err)
	}
	got, _ := svc.GetArticle(ctx, a.ID)
	if got.Status != model.StatusPendingReview {
		t.Fatalf("状态应为 pending_review，got %s", got.Status)
	}
	// 已提交不能再提交（pending_review → pending_review 非法）
	if err := svc.SubmitArticle(ctx, 1, a.ID); !errors.Is(err, ErrIllegalState) {
		t.Fatalf("期望 ErrIllegalState，got %v", err)
	}
}

func TestPublishRejectFlow(t *testing.T) {
	svc := newTestArticleService()
	ctx := context.Background()

	// 发布：draft → submit → publish
	a, _ := svc.CreateArticle(ctx, 1, "标题", "内容", nil)
	_ = svc.SubmitArticle(ctx, 1, a.ID)
	if err := svc.PublishArticle(ctx, a.ID); err != nil {
		t.Fatalf("PublishArticle() err = %v", err)
	}
	pub, _ := svc.GetArticle(ctx, a.ID)
	if pub.Status != model.StatusPublished {
		t.Fatalf("状态应为 published，got %s", pub.Status)
	}

	// 驳回：draft → submit → reject
	b, _ := svc.CreateArticle(ctx, 1, "标题2", "内容2", nil)
	_ = svc.SubmitArticle(ctx, 1, b.ID)
	if err := svc.RejectArticle(ctx, b.ID); err != nil {
		t.Fatalf("RejectArticle() err = %v", err)
	}
	rej, _ := svc.GetArticle(ctx, b.ID)
	if rej.Status != model.StatusRejected {
		t.Fatalf("状态应为 rejected，got %s", rej.Status)
	}
}
