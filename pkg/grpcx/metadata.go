package grpcx

import (
	"context"
	"strconv"

	"google.golang.org/grpc/metadata"

	"github.com/likeslep/community/pkg/apperr"
)

// gRPC metadata 键，用于在网关与业务服务之间传播认证身份（plan.md §11 Gateway 负责基础认证）。
const (
	MDUserID = "x-user-id"
	MDRole   = "x-role"
)

// WithUserID 将用户 ID 注入出站 metadata。
func WithUserID(ctx context.Context, userID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, MDUserID, userID)
}

// WithRole 将角色注入出站 metadata。
func WithRole(ctx context.Context, role string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, MDRole, role)
}

// UserIDFrom 从入站 metadata 读取用户 ID。
func UserIDFrom(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if vals := md.Get(MDUserID); len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// RoleFrom 从入站 metadata 读取角色。
func RoleFrom(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if vals := md.Get(MDRole); len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// AuthenticatedUserID 从 metadata 解析认证用户 ID，未认证返回 401 错误。
func AuthenticatedUserID(ctx context.Context) (uint64, error) {
	uid := UserIDFrom(ctx)
	if uid == "" {
		return 0, apperr.New(apperr.CodeUser+6, "未认证", apperr.WithHTTP(401))
	}
	id, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		return 0, apperr.New(apperr.CodeUser+6, "未认证", apperr.WithHTTP(401))
	}
	return id, nil
}
