// Package handler 是 notification-service 的 gRPC 处理器。
package handler

import (
	"context"

	notificationv1 "github.com/likeslep/community/api/gen/notification/v1"
	"github.com/likeslep/community/internal/notification/service"
	"github.com/likeslep/community/pkg/grpcx"
)

// Handler 实现 notification gRPC 服务。
type Handler struct {
	notificationv1.UnimplementedNotificationServiceServer
	svc *service.NotificationService
}

// New 构造。
func New(svc *service.NotificationService) *Handler { return &Handler{svc: svc} }

func (h *Handler) ListNotifications(ctx context.Context, req *notificationv1.ListNotificationsRequest) (*notificationv1.ListNotificationsResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	ns, err := h.svc.ListNotifications(ctx, uid, int(req.GetLimit()))
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	resp := &notificationv1.ListNotificationsResponse{}
	for _, n := range ns {
		resp.Notifications = append(resp.Notifications, &notificationv1.Notification{
			Id: n.ID, ActorId: n.ActorID, Type: n.Type, TargetType: n.TargetType,
			TargetId: n.TargetID, Content: n.Content, Read: n.Read,
		})
	}
	return resp, nil
}

func (h *Handler) UnreadCount(ctx context.Context, _ *notificationv1.UnreadCountRequest) (*notificationv1.UnreadCountResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	count, err := h.svc.UnreadCount(ctx, uid)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &notificationv1.UnreadCountResponse{Count: count}, nil
}

func (h *Handler) MarkRead(ctx context.Context, req *notificationv1.MarkReadRequest) (*notificationv1.MarkReadResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.MarkRead(ctx, uid, req.GetId()); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &notificationv1.MarkReadResponse{}, nil
}

func (h *Handler) MarkAllRead(ctx context.Context, _ *notificationv1.MarkAllReadRequest) (*notificationv1.MarkAllReadResponse, error) {
	uid, err := grpcx.AuthenticatedUserID(ctx)
	if err != nil {
		return nil, grpcx.ToStatus(err)
	}
	if err := h.svc.MarkAllRead(ctx, uid); err != nil {
		return nil, grpcx.ToStatus(err)
	}
	return &notificationv1.MarkAllReadResponse{}, nil
}
