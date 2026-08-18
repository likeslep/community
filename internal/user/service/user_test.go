package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/likeslep/community/internal/user/model"
	"github.com/likeslep/community/pkg/kafkax"
)

// fakeRepo 是 Repository 的内存实现，仅用于测试。
// 并发安全：用 mutex 保护 map，避免并发读写竞态。
type fakeRepo struct {
	mu     sync.Mutex
	users  map[string]*model.User
	emails map[string]*model.User
	nextID uint64
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{users: map[string]*model.User{}, emails: map[string]*model.User{}, nextID: 1}
}

func (f *fakeRepo) Create(_ context.Context, user *model.User, buildEvent func(*model.User) (kafkax.Envelope, error)) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	user.ID = f.nextID
	f.nextID++
	f.users[user.Username] = user
	f.emails[user.Email] = user
	_, err := buildEvent(user)
	return err
}

func (f *fakeRepo) FindByUsername(_ context.Context, username string) (*model.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.users[username]; ok {
		return u, nil
	}
	return nil, ErrUserNotFound
}

func (f *fakeRepo) FindByEmail(_ context.Context, email string) (*model.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.emails[email]; ok {
		return u, nil
	}
	return nil, ErrUserNotFound
}

func (f *fakeRepo) FindByID(_ context.Context, id uint64) (*model.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, u := range f.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}

func (f *fakeRepo) Update(_ context.Context, _ *model.User) error { return nil }

func (f *fakeRepo) ListUsers(_ context.Context, _ int, _ int) ([]model.User, error) {
	return nil, nil
}

func (f *fakeRepo) UpdateStatus(_ context.Context, _ uint64, _ string) error { return nil }

func newTestService() *UserService {
	return NewUserService(newFakeRepo(), Config{
		Producer:  "user-service",
		JWTSecret: []byte("test-secret"),
		TokenTTL:  time.Hour,
	})
}

func TestRegister(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	t.Run("成功注册", func(t *testing.T) {
		u, err := svc.Register(ctx, "alice", "alice@example.com", "password123")
		if err != nil {
			t.Fatalf("Register() err = %v", err)
		}
		if u.ID == 0 || u.Username != "alice" {
			t.Fatalf("user = %+v", u)
		}
		if u.PasswordHash == "password123" {
			t.Fatal("密码不应明文存储")
		}
	})

	t.Run("重复用户名", func(t *testing.T) {
		if _, err := svc.Register(ctx, "alice", "other@example.com", "password123"); !errors.Is(err, ErrUsernameTaken) {
			t.Fatalf("期望 ErrUsernameTaken，got %v", err)
		}
	})

	t.Run("重复邮箱", func(t *testing.T) {
		if _, err := svc.Register(ctx, "bob", "alice@example.com", "password123"); !errors.Is(err, ErrEmailTaken) {
			t.Fatalf("期望 ErrEmailTaken，got %v", err)
		}
	})

	t.Run("密码过短", func(t *testing.T) {
		if _, err := svc.Register(ctx, "carol", "", "123"); err == nil {
			t.Fatal("期望密码过短报错，实际通过")
		}
	})
}

func TestLogin(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	if _, err := svc.Register(ctx, "alice", "alice@example.com", "password123"); err != nil {
		t.Fatalf("Register() err = %v", err)
	}

	t.Run("成功登录", func(t *testing.T) {
		token, u, err := svc.Login(ctx, "alice", "password123")
		if err != nil {
			t.Fatalf("Login() err = %v", err)
		}
		if token == "" || u.Username != "alice" {
			t.Fatalf("token=%q user=%+v", token, u)
		}
	})

	t.Run("密码错误", func(t *testing.T) {
		if _, _, err := svc.Login(ctx, "alice", "wrongpass"); !errors.Is(err, ErrInvalidPassword) {
			t.Fatalf("期望 ErrInvalidPassword，got %v", err)
		}
	})

	t.Run("用户不存在", func(t *testing.T) {
		if _, _, err := svc.Login(ctx, "nobody", "whatever"); !errors.Is(err, ErrInvalidPassword) {
			t.Fatalf("期望 ErrInvalidPassword（不泄露存在性），got %v", err)
		}
	})
}
