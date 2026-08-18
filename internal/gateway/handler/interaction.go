package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	interactionv1 "github.com/likeslep/community/api/gen/interaction/v1"
	gwmw "github.com/likeslep/community/internal/gateway/middleware"
	"github.com/likeslep/community/pkg/grpcx"
)

func (h *Handler) Like(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint64 `json:"target_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.interactions.Like(ctx, &interactionv1.LikeRequest{TargetType: req.TargetType, TargetId: req.TargetID}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func (h *Handler) Unlike(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint64 `json:"target_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.interactions.Unlike(ctx, &interactionv1.UnlikeRequest{TargetType: req.TargetType, TargetId: req.TargetID}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func (h *Handler) Collect(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint64 `json:"target_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.interactions.Collect(ctx, &interactionv1.CollectRequest{TargetType: req.TargetType, TargetId: req.TargetID}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func (h *Handler) Uncollect(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint64 `json:"target_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.interactions.Uncollect(ctx, &interactionv1.UncollectRequest{TargetType: req.TargetType, TargetId: req.TargetID}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func (h *Handler) View(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint64 `json:"target_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	resp, err := h.interactions.View(c.Request.Context(), &interactionv1.ViewRequest{TargetType: req.TargetType, TargetId: req.TargetID})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

func (h *Handler) CreateComment(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint64 `json:"target_id" binding:"required"`
		Content    string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	resp, err := h.comments.CreateComment(ctx, &interactionv1.CreateCommentRequest{
		TargetType: req.TargetType, TargetId: req.TargetID, Content: req.Content,
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

func (h *Handler) ListComments(c *gin.Context) {
	targetType := c.Query("target_type")
	targetID, _ := strconv.ParseUint(c.Query("target_id"), 10, 64)
	resp, err := h.comments.ListComments(c.Request.Context(), &interactionv1.ListCommentsRequest{
		TargetType: targetType, TargetId: targetID,
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}
