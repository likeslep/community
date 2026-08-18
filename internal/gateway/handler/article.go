package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	contentv1 "github.com/likeslep/community/api/gen/content/v1"
	gwmw "github.com/likeslep/community/internal/gateway/middleware"
	"github.com/likeslep/community/pkg/grpcx"
)

// ListArticles 分页查询文章（默认 published）。
func (h *Handler) ListArticles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.DefaultQuery("status", "published")
	resp, err := h.articles.ListArticles(c.Request.Context(), &contentv1.ListArticlesRequest{
		Limit: int32(limit), Offset: int32(offset), Status: status,
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// CreateArticle 创建文章（受保护）。
func (h *Handler) CreateArticle(c *gin.Context) {
	var req struct {
		Title   string   `json:"title" binding:"required"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	resp, err := h.articles.CreateArticle(ctx, &contentv1.CreateArticleRequest{
		Title: req.Title, Content: req.Content, Tags: req.Tags,
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// UpdateArticle 更新文章（受保护）。
func (h *Handler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的文章 ID"})
		return
	}
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.articles.UpdateArticle(ctx, &contentv1.UpdateArticleRequest{
		Id: id, Title: req.Title, Content: req.Content,
	}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

// GetArticle 查询文章。
func (h *Handler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的文章 ID"})
		return
	}
	resp, err := h.articles.GetArticle(c.Request.Context(), &contentv1.GetArticleRequest{Id: id})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// DeleteArticle 删除文章（受保护）。
func (h *Handler) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的文章 ID"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.articles.DeleteArticle(ctx, &contentv1.DeleteArticleRequest{Id: id}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

// SubmitArticle 提交审核（受保护）。
func (h *Handler) SubmitArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的文章 ID"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.articles.SubmitArticle(ctx, &contentv1.SubmitArticleRequest{Id: id}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}
