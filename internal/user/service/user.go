// Package service 是 user-service 的业务逻辑层。
package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/likeslep/community/internal/user/model"
	"github.com/likeslep/community/pkg/apperr"
	"github.com/likeslep/community/pkg/auth"
	"github.com/likeslep/community/pkg/kafkax"
)

// Config 是 user-service 的业务配置。
type Config struct {
	Producer  string        // 服务名，用于事件 producer 字段
	JWTSecret []byte        // JWT 签名密钥
	TokenTTL  time.Duration // token 有效期
}

// UserService 是 user-service 业务逻辑层。
type UserService struct {
	repo Repository
	cfg  Config
}

// NewUserService 构造。
func NewUserService(repo Repository, cfg Config) *UserService {
	return &UserService{repo: repo, cfg: cfg}
}

// Register 注册用户：校验 → 哈希密码 → 事务写入用户 + outbox 事件。
func (s *UserService) Register(ctx context.Context, username, email, password string) (*model.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	if username == "" || len(password) < 6 {
		return nil, apperr.New(errCodeInvalidInput, "用户名不能为空且密码至少 6 位", apperr.WithHTTP(400))
	}

	if _, err := s.repo.FindByUsername(ctx, username); err == nil {
		return nil, ErrUsernameTaken
	} else if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}
	if email != "" {
		if _, err := s.repo.FindByEmail(ctx, email); err == nil {
			return nil, ErrEmailTaken
		} else if !errors.Is(err, ErrUserNotFound) {
			return nil, err
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeUser, "密码加密失败", apperr.WithKind(apperr.KindSystem))
	}

	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Role:         model.RoleAuthor,
		Status:       model.StatusActive,
	}

	err = s.repo.Create(ctx, user, func(u *model.User) (kafkax.Envelope, error) {
		return kafkax.NewEnvelope(kafkax.EventUserCreated, s.cfg.Producer, "user",
			strconv.FormatUint(u.ID, 10), 1, userCreatedPayload{UserID: u.ID, Username: u.Username})
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Login 校验凭证并签发 JWT。
func (s *UserService) Login(ctx context.Context, username, password string) (string, *model.User, error) {
	user, err := s.repo.FindByUsername(ctx, strings.TrimSpace(username))
	if err != nil {
		return "", nil, ErrInvalidPassword // 统一返回，不泄露用户是否存在
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", nil, ErrInvalidPassword
	}
	token, err := auth.Sign(s.cfg.JWTSecret, strconv.FormatUint(user.ID, 10), user.Username, user.Role, s.cfg.TokenTTL)
	if err != nil {
		return "", nil, apperr.Wrap(err, apperr.CodeUser, "签发 token 失败", apperr.WithKind(apperr.KindSystem))
	}
	return token, user, nil
}

// GetProfile 查询用户资料。
func (s *UserService) GetProfile(ctx context.Context, id uint64) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

// UpdateProfile 更新用户资料。
func (s *UserService) UpdateProfile(ctx context.Context, id uint64, email, bio, avatar string) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if email != "" {
		user.Email = strings.TrimSpace(email)
	}
	if bio != "" {
		user.Bio = bio
	}
	if avatar != "" {
		user.AvatarFileID = avatar
	}
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUser 查询用户（供内部/gateway 鉴权使用）。
func (s *UserService) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

// ListUsers 分页查询用户。
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]model.User, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.ListUsers(ctx, limit, offset)
}

// BanUser 封禁用户。
func (s *UserService) BanUser(ctx context.Context, id uint64) error {
	return s.repo.UpdateStatus(ctx, id, model.StatusBanned)
}

type userCreatedPayload struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
}
