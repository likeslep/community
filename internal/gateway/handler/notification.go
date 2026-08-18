package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	notificationv1 "github.com/likeslep/community/api/gen/notification/v1"
	gwmw "github.com/likeslep/community/internal/gateway/middleware"
	"github.com/likeslep/community/pkg/grpcx"
)

func (h *Handler) ListNotifications(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	resp, err := h.notifications.ListNotifications(ctx, &notificationv1.ListNotificationsRequest{Limit: int32(limit)})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

func (h *Handler) UnreadCount(c *gin.Context) {
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	resp, err := h.notifications.UnreadCount(ctx, &notificationv1.UnreadCountRequest{})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

func (h *Handler) MarkRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的通知 ID"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.notifications.MarkRead(ctx, &notificationv1.MarkReadRequest{Id: id}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func (h *Handler) MarkAllRead(c *gin.Context) {
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.notifications.MarkAllRead(ctx, &notificationv1.MarkAllReadRequest{}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}
