// Package handler 是 gateway 的 REST 处理器，通过 gRPC 转发到后端服务。
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	adminv1 "github.com/likeslep/community/api/gen/admin/v1"
	contentv1 "github.com/likeslep/community/api/gen/content/v1"
	feedv1 "github.com/likeslep/community/api/gen/feed/v1"
	interactionv1 "github.com/likeslep/community/api/gen/interaction/v1"
	notificationv1 "github.com/likeslep/community/api/gen/notification/v1"
	searchv1 "github.com/likeslep/community/api/gen/search/v1"
	socialv1 "github.com/likeslep/community/api/gen/social/v1"
	userv1 "github.com/likeslep/community/api/gen/user/v1"
	gwmw "github.com/likeslep/community/internal/gateway/middleware"
	"github.com/likeslep/community/pkg/grpcx"
)

// Clients 汇总所有后端服务的 gRPC 客户端。
type Clients struct {
	Users         userv1.UserServiceClient
	Articles      contentv1.ArticleServiceClient
	Questions     contentv1.QuestionServiceClient
	Admin         adminv1.AdminServiceClient
	Interactions  interactionv1.InteractionServiceClient
	Comments      interactionv1.CommentServiceClient
	Social        socialv1.SocialServiceClient
	Search        searchv1.SearchServiceClient
	Feed          feedv1.FeedServiceClient
	Notifications notificationv1.NotificationServiceClient
}

// Handler 是 gateway 的 REST 处理器。
type Handler struct {
	users         userv1.UserServiceClient
	articles      contentv1.ArticleServiceClient
	questions     contentv1.QuestionServiceClient
	admin         adminv1.AdminServiceClient
	interactions  interactionv1.InteractionServiceClient
	comments      interactionv1.CommentServiceClient
	social        socialv1.SocialServiceClient
	search        searchv1.SearchServiceClient
	feed          feedv1.FeedServiceClient
	notifications notificationv1.NotificationServiceClient
}

// New 构造。
func New(c Clients) *Handler {
	return &Handler{
		users: c.Users, articles: c.Articles, questions: c.Questions, admin: c.Admin,
		interactions: c.Interactions, comments: c.Comments, social: c.Social,
		search: c.Search, feed: c.Feed, notifications: c.Notifications,
	}
}

// RegisterRoutes 注册全部路由。
func (h *Handler) RegisterRoutes(engine *gin.Engine, auth gin.HandlerFunc) {
	api := engine.Group("/api/v1")

	// 认证
	api.POST("/auth/register", h.Register)
	api.POST("/auth/login", h.Login)

	// 用户
	api.GET("/users/:id", h.GetProfile)
	api.PUT("/users/me", auth, h.UpdateProfile)

	// 文章
	api.GET("/articles", h.ListArticles)
	api.POST("/articles", auth, h.CreateArticle)
	api.GET("/articles/:id", h.GetArticle)
	api.PUT("/articles/:id", auth, h.UpdateArticle)
	api.DELETE("/articles/:id", auth, h.DeleteArticle)
	api.POST("/articles/:id/submit", auth, h.SubmitArticle)

	// 问答
	api.GET("/questions", h.ListQuestions)
	api.POST("/questions", auth, h.CreateQuestion)
	api.GET("/questions/:id", h.GetQuestion)
	api.POST("/questions/:id/close", auth, h.CloseQuestion)
	api.POST("/questions/:id/answers", auth, h.CreateAnswer)
	api.POST("/questions/:id/accept", auth, h.AcceptAnswer)

	// 互动
	api.POST("/interactions/like", auth, h.Like)
	api.POST("/interactions/unlike", auth, h.Unlike)
	api.POST("/interactions/collect", auth, h.Collect)
	api.POST("/interactions/uncollect", auth, h.Uncollect)
	api.POST("/interactions/view", h.View)
	api.POST("/comments", auth, h.CreateComment)
	api.GET("/comments", h.ListComments)

	// 社交
	api.POST("/users/:id/follow", auth, h.FollowUser)
	api.DELETE("/users/:id/follow", auth, h.UnfollowUser)
	api.POST("/tags/:id/follow", auth, h.FollowTag)
	api.DELETE("/tags/:id/follow", auth, h.UnfollowTag)

	// 搜索 / 信息流 / 通知
	api.GET("/search", h.Search)
	api.GET("/feed", auth, h.GetFeed)
	api.GET("/notifications", auth, h.ListNotifications)
	api.GET("/notifications/unread-count", auth, h.UnreadCount)
	api.POST("/notifications/:id/read", auth, h.MarkRead)
	api.POST("/notifications/read-all", auth, h.MarkAllRead)

	// 后台管理
	admin := api.Group("/admin", auth, gwmw.RequireAdmin())
	admin.GET("/users", h.AdminListUsers)
	admin.POST("/users/:id/ban", h.AdminBanUser)
	admin.POST("/articles/:id/hide", h.AdminHideArticle)
	admin.DELETE("/comments/:id", h.AdminDeleteComment)
	admin.GET("/tags", h.AdminListTags)
	admin.GET("/sensitive-words", h.AdminListSensitiveWords)
	admin.POST("/sensitive-words", h.AdminCreateSensitiveWord)
	admin.GET("/reports", h.AdminListReports)
	admin.POST("/reports/:id/handle", h.AdminHandleReport)
	admin.GET("/audit-logs", h.AdminListAuditLogs)
	admin.GET("/statistics", h.AdminStatistics)
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
