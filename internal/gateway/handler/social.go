package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	socialv1 "github.com/likeslep/community/api/gen/social/v1"
	gwmw "github.com/likeslep/community/internal/gateway/middleware"
	"github.com/likeslep/community/pkg/grpcx"
)

func (h *Handler) FollowUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的用户 ID"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.social.FollowUser(ctx, &socialv1.FollowUserRequest{FolloweeId: id}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func (h *Handler) UnfollowUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的用户 ID"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.social.UnfollowUser(ctx, &socialv1.UnfollowUserRequest{FolloweeId: id}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func (h *Handler) FollowTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的标签 ID"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.social.FollowTag(ctx, &socialv1.FollowTagRequest{TagId: id}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func (h *Handler) UnfollowTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的标签 ID"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.social.UnfollowTag(ctx, &socialv1.UnfollowTagRequest{TagId: id}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}
