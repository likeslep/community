package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	contentv1 "github.com/likeslep/community/api/gen/content/v1"
	gwmw "github.com/likeslep/community/internal/gateway/middleware"
	"github.com/likeslep/community/pkg/grpcx"
)

// ListQuestions 分页查询问题。
func (h *Handler) ListQuestions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	resp, err := h.questions.ListQuestions(c.Request.Context(), &contentv1.ListQuestionsRequest{
		Limit: int32(limit), Offset: int32(offset),
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// CreateQuestion 提问（受保护）。
func (h *Handler) CreateQuestion(c *gin.Context) {
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
	resp, err := h.questions.CreateQuestion(ctx, &contentv1.CreateQuestionRequest{
		Title: req.Title, Content: req.Content, Tags: req.Tags,
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// GetQuestion 查询问题。
func (h *Handler) GetQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的问题 ID"})
		return
	}
	resp, err := h.questions.GetQuestion(c.Request.Context(), &contentv1.GetQuestionRequest{Id: id})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// CloseQuestion 关闭问题（受保护）。
func (h *Handler) CloseQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的问题 ID"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.questions.CloseQuestion(ctx, &contentv1.CloseQuestionRequest{Id: id}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

// CreateAnswer 回答问题（受保护）。
func (h *Handler) CreateAnswer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的问题 ID"})
		return
	}
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	resp, err := h.questions.CreateAnswer(ctx, &contentv1.CreateAnswerRequest{
		QuestionId: id, Content: req.Content,
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// AcceptAnswer 采纳回答（受保护）。
func (h *Handler) AcceptAnswer(c *gin.Context) {
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的问题 ID"})
		return
	}
	var req struct {
		AnswerID uint64 `json:"answer_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.questions.AcceptAnswer(ctx, &contentv1.AcceptAnswerRequest{
		QuestionId: questionID, AnswerId: req.AnswerID,
	}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}
