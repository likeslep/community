package service

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/likeslep/community/internal/content/model"
	"github.com/likeslep/community/pkg/kafkax"
)

// fakeQARepo 是 QARepository 的内存实现，仅用于测试。
// 并发安全：用 mutex 保护 map，避免并发读写竞态。
type fakeQARepo struct {
	mu        sync.Mutex
	questions map[uint64]*model.Question
	answers   map[uint64]*model.Answer
	nextQID   uint64
	nextAID   uint64
}

func newFakeQARepo() *fakeQARepo {
	return &fakeQARepo{
		questions: map[uint64]*model.Question{},
		answers:   map[uint64]*model.Answer{},
		nextQID:   1, nextAID: 1,
	}
}

func (f *fakeQARepo) CreateQuestion(_ context.Context, q *model.Question, _ []string, build BuildQuestion) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	q.ID = f.nextQID
	f.nextQID++
	f.questions[q.ID] = q
	_, err := build(q)
	return err
}

func (f *fakeQARepo) FindQuestion(_ context.Context, id uint64) (*model.Question, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if q, ok := f.questions[id]; ok {
		return q, nil
	}
	return nil, ErrQuestionNotFound
}

func (f *fakeQARepo) UpdateQuestion(_ context.Context, q *model.Question) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.questions[q.ID] = q
	return nil
}

func (f *fakeQARepo) CreateAnswer(_ context.Context, a *model.Answer, build BuildAnswer) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	a.ID = f.nextAID
	f.nextAID++
	f.answers[a.ID] = a
	_, err := build(a)
	return err
}

func (f *fakeQARepo) FindAnswer(_ context.Context, id uint64) (*model.Answer, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if a, ok := f.answers[id]; ok {
		return a, nil
	}
	return nil, ErrAnswerNotFound
}

func (f *fakeQARepo) AcceptAnswer(_ context.Context, questionID, answerID uint64, _ kafkax.Envelope) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, a := range f.answers {
		if a.QuestionID == questionID {
			a.Accepted = false
		}
	}
	f.answers[answerID].Accepted = true
	f.questions[questionID].AcceptedAnswerID = &answerID
	return nil
}

func newTestQAService() *QAService {
	return NewQAService(newFakeQARepo(), Config{Producer: "content-service"})
}

func TestCreateQuestion(t *testing.T) {
	svc := newTestQAService()
	ctx := context.Background()

	q, err := svc.CreateQuestion(ctx, 1, "如何学习 Go？", "内容", []string{"go"})
	if err != nil {
		t.Fatalf("CreateQuestion() err = %v", err)
	}
	if q.ID == 0 || q.Status != model.QuestionOpen {
		t.Fatalf("question = %+v", q)
	}
	if _, err := svc.CreateQuestion(ctx, 1, "", "x", nil); err == nil {
		t.Fatal("空标题应报错")
	}
}

func TestCloseQuestionOwnership(t *testing.T) {
	svc := newTestQAService()
	ctx := context.Background()
	q, _ := svc.CreateQuestion(ctx, 1, "标题", "内容", nil)

	if err := svc.CloseQuestion(ctx, 2, q.ID); !errors.Is(err, ErrForbidden) {
		t.Fatalf("非提问者关闭应返回 ErrForbidden，got %v", err)
	}
	if err := svc.CloseQuestion(ctx, 1, q.ID); err != nil {
		t.Fatalf("CloseQuestion() err = %v", err)
	}
}

func TestAcceptAnswer(t *testing.T) {
	svc := newTestQAService()
	ctx := context.Background()
	q, _ := svc.CreateQuestion(ctx, 1, "标题", "内容", nil)
	a, _ := svc.CreateAnswer(ctx, 2, q.ID, "回答内容")

	if err := svc.AcceptAnswer(ctx, 2, q.ID, a.ID); !errors.Is(err, ErrForbidden) {
		t.Fatalf("非提问者采纳应返回 ErrForbidden，got %v", err)
	}
	if err := svc.AcceptAnswer(ctx, 1, q.ID, a.ID); err != nil {
		t.Fatalf("AcceptAnswer() err = %v", err)
	}
}
