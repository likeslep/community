package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	adminv1 "github.com/likeslep/community/api/gen/admin/v1"
	gwmw "github.com/likeslep/community/internal/gateway/middleware"
	"github.com/likeslep/community/pkg/grpcx"
)

// adminCtx 返回带管理员身份的 context。
func adminCtx(c *gin.Context) context.Context {
	return grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
}

// AdminListUsers 查询用户。
func (h *Handler) AdminListUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	resp, err := h.admin.ListUsers(c.Request.Context(), &adminv1.ListUsersRequest{
		Limit: int32(limit), Offset: int32(offset),
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// AdminBanUser 封禁用户。
func (h *Handler) AdminBanUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的用户 ID"})
		return
	}
	if _, err := h.admin.BanUser(adminCtx(c), &adminv1.BanUserRequest{Id: id}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

// AdminHideArticle 隐藏文章。
func (h *Handler) AdminHideArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的文章 ID"})
		return
	}
	if _, err := h.admin.HideArticle(adminCtx(c), &adminv1.HideArticleRequest{Id: id}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

// AdminDeleteComment 删除评论。
func (h *Handler) AdminDeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的评论 ID"})
		return
	}
	if _, err := h.admin.DeleteComment(adminCtx(c), &adminv1.DeleteCommentRequest{Id: id}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

// AdminListTags 查询标签。
func (h *Handler) AdminListTags(c *gin.Context) {
	resp, err := h.admin.ListTags(c.Request.Context(), &adminv1.ListTagsRequest{})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// AdminListSensitiveWords 查询敏感词。
func (h *Handler) AdminListSensitiveWords(c *gin.Context) {
	resp, err := h.admin.ListSensitiveWords(c.Request.Context(), &adminv1.ListSensitiveWordsRequest{})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// AdminCreateSensitiveWord 添加敏感词。
func (h *Handler) AdminCreateSensitiveWord(c *gin.Context) {
	var req struct {
		Word  string `json:"word" binding:"required"`
		Level string `json:"level"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	if _, err := h.admin.CreateSensitiveWord(adminCtx(c), &adminv1.CreateSensitiveWordRequest{
		Word: req.Word, Level: req.Level,
	}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

// AdminListReports 查询举报。
func (h *Handler) AdminListReports(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	resp, err := h.admin.ListReports(c.Request.Context(), &adminv1.ListReportsRequest{Limit: int32(limit)})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// AdminHandleReport 处理举报。
func (h *Handler) AdminHandleReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的举报 ID"})
		return
	}
	var req struct {
		Action string `json:"action" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	if _, err := h.admin.HandleReport(adminCtx(c), &adminv1.HandleReportRequest{Id: id, Action: req.Action}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

// AdminListAuditLogs 查询审计日志。
func (h *Handler) AdminListAuditLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	resp, err := h.admin.ListAuditLogs(c.Request.Context(), &adminv1.ListAuditLogsRequest{Limit: int32(limit)})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// AdminStatistics 数据统计。
func (h *Handler) AdminStatistics(c *gin.Context) {
	resp, err := h.admin.Statistics(c.Request.Context(), &adminv1.StatisticsRequest{})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}
