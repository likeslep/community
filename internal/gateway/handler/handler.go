// Package handler 是 gateway 的 REST 处理器，通过 gRPC 转发到后端服务。
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	contentv1 "github.com/likeslep/community/api/gen/content/v1"
	userv1 "github.com/likeslep/community/api/gen/user/v1"
	gwmw "github.com/likeslep/community/internal/gateway/middleware"
	"github.com/likeslep/community/pkg/grpcx"
)

// Handler 是 gateway 的 REST 处理器。
type Handler struct {
	users     userv1.UserServiceClient
	articles  contentv1.ArticleServiceClient
	questions contentv1.QuestionServiceClient
}

// New 构造。
func New(users userv1.UserServiceClient, articles contentv1.ArticleServiceClient, questions contentv1.QuestionServiceClient) *Handler {
	return &Handler{users: users, articles: articles, questions: questions}
}

// RegisterRoutes 注册路由。auth 为受保护路由的认证中间件。
func (h *Handler) RegisterRoutes(engine *gin.Engine, auth gin.HandlerFunc) {
	api := engine.Group("/api/v1")
	api.POST("/auth/register", h.Register)
	api.POST("/auth/login", h.Login)
	api.GET("/users/:id", h.GetProfile)
	api.PUT("/users/me", auth, h.UpdateProfile)

	api.POST("/articles", auth, h.CreateArticle)
	api.PUT("/articles/:id", auth, h.UpdateArticle)
	api.GET("/articles/:id", h.GetArticle)
	api.DELETE("/articles/:id", auth, h.DeleteArticle)
	api.POST("/articles/:id/submit", auth, h.SubmitArticle)

	api.POST("/questions", auth, h.CreateQuestion)
	api.GET("/questions/:id", h.GetQuestion)
	api.POST("/questions/:id/close", auth, h.CloseQuestion)
	api.POST("/questions/:id/answers", auth, h.CreateAnswer)
	api.POST("/questions/:id/accept", auth, h.AcceptAnswer)
}

// Register 处理注册。
func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	resp, err := h.users.Register(c.Request.Context(), &userv1.RegisterRequest{
		Username: req.Username, Email: req.Email, Password: req.Password,
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// Login 处理登录。
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	resp, err := h.users.Login(c.Request.Context(), &userv1.LoginRequest{
		Username: req.Username, Password: req.Password,
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// GetProfile 查询用户资料。
func (h *Handler) GetProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的用户 ID"})
		return
	}
	resp, err := h.users.GetProfile(c.Request.Context(), &userv1.GetProfileRequest{UserId: id})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}

// UpdateProfile 更新当前用户资料（受保护）。
func (h *Handler) UpdateProfile(c *gin.Context) {
	var req struct {
		Email        string `json:"email"`
		Bio          string `json:"bio"`
		AvatarFileID string `json:"avatar_file_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	if _, err := h.users.UpdateProfile(ctx, &userv1.UpdateProfileRequest{
		Email: req.Email, Bio: req.Bio, AvatarFileId: req.AvatarFileID,
	}); err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func writeErr(c *gin.Context, err error) {
	status := grpcx.HTTPStatus(err)
	c.JSON(status, gin.H{"code": status, "message": grpcx.Message(err)})
}
