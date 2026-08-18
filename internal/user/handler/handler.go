// Package handler 是 user-service 的 gRPC 处理器（transport 层）。
package handler

import (
	"context"
	"strconv"

	userv1 "github.com/likeslep/community/api/gen/user/v1"
	"github.com/likeslep/community/internal/user/service"
	"github.com/likeslep/community/pkg/apperr"
	"github.com/likeslep/community/pkg/grpcx"
)

// Handler 实现 user gRPC 服务。
type Handler struct {
	userv1.UnimplementedUserServiceServer
	svc *service.UserService
}

// New 构造 Handler。
func New(svc *service.UserService) *Handler { return &Handler{svc: svc} }

// Register 注册用户。
func (h *Handler) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	u, err := h.svc.Register(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &userv1.RegisterResponse{Id: u.ID, Username: u.Username}, nil
}

// Login 登录并签发 JWT。
func (h *Handler) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	token, u, err := h.svc.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &userv1.LoginResponse{Token: token, UserId: u.ID, Username: u.Username, Role: u.Role}, nil
}

// GetProfile 查询用户资料。
func (h *Handler) GetProfile(ctx context.Context, req *userv1.GetProfileRequest) (*userv1.GetProfileResponse, error) {
	u, err := h.svc.GetProfile(ctx, req.GetUserId())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &userv1.GetProfileResponse{
		Id: u.ID, Username: u.Username, Email: u.Email, Role: u.Role,
		Status: u.Status, Bio: u.Bio, AvatarFileId: u.AvatarFileID,
	}, nil
}

// UpdateProfile 更新用户资料。从 gRPC metadata 取认证用户 ID 做所有权校验（plan.md §12）。
func (h *Handler) UpdateProfile(ctx context.Context, req *userv1.UpdateProfileRequest) (*userv1.UpdateProfileResponse, error) {
	userID, err := authenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if _, err := h.svc.UpdateProfile(ctx, userID, req.GetEmail(), req.GetBio(), req.GetAvatarFileId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &userv1.UpdateProfileResponse{}, nil
}

// authenticatedUserID 从 gRPC metadata 解析认证用户 ID，未认证返回 401。
func authenticatedUserID(ctx context.Context) (uint64, error) {
	uid := grpcx.UserIDFrom(ctx)
	if uid == "" {
		return 0, apperr.New(apperr.CodeUser+6, "未认证", apperr.WithHTTP(401))
	}
	id, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		return 0, apperr.New(apperr.CodeUser+6, "未认证", apperr.WithHTTP(401))
	}
	return id, nil
}

// GetUser 查询用户（供内部/gateway 鉴权使用）。
func (h *Handler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	u, err := h.svc.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &userv1.GetUserResponse{Id: u.ID, Username: u.Username, Role: u.Role, Status: u.Status}, nil
}

// ListUsers 分页查询用户。
func (h *Handler) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	users, err := h.svc.ListUsers(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &userv1.ListUsersResponse{}
	for _, u := range users {
		resp.Users = append(resp.Users, &userv1.UserBrief{
			Id: u.ID, Username: u.Username, Role: u.Role, Status: u.Status,
		})
	}
	return resp, nil
}

// BanUser 封禁用户。
func (h *Handler) BanUser(ctx context.Context, req *userv1.BanUserRequest) (*userv1.BanUserResponse, error) {
	if err := h.svc.BanUser(ctx, req.GetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &userv1.BanUserResponse{}, nil
}
