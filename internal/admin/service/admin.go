// Package service 是 admin-service 的业务逻辑层（通过 gRPC 编排，不直连他库，plan.md §6.11）。
package service

import (
	"context"

	contentv1 "github.com/likeslep/community/api/gen/content/v1"
	interactionv1 "github.com/likeslep/community/api/gen/interaction/v1"
	moderationv1 "github.com/likeslep/community/api/gen/moderation/v1"
	userv1 "github.com/likeslep/community/api/gen/user/v1"
	"github.com/likeslep/community/internal/admin/model"
)

// AdminService 是后台管理业务逻辑层。
type AdminService struct {
	users          userv1.UserServiceClient
	articles       contentv1.ArticleServiceClient
	comments       interactionv1.CommentServiceClient
	reports        moderationv1.ReportServiceClient
	sensitiveWords moderationv1.SensitiveWordServiceClient
	repo           Repository
}

// NewAdminService 构造。
func NewAdminService(
	users userv1.UserServiceClient,
	articles contentv1.ArticleServiceClient,
	comments interactionv1.CommentServiceClient,
	reports moderationv1.ReportServiceClient,
	sensitiveWords moderationv1.SensitiveWordServiceClient,
	repo Repository,
) *AdminService {
	return &AdminService{
		users: users, articles: articles, comments: comments,
		reports: reports, sensitiveWords: sensitiveWords, repo: repo,
	}
}

// ListUsers 查询用户（委托 user-service）。
func (s *AdminService) ListUsers(ctx context.Context, limit, offset int) ([]*userv1.UserBrief, error) {
	resp, err := s.users.ListUsers(ctx, &userv1.ListUsersRequest{Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	return resp.GetUsers(), nil
}

// BanUser 封禁用户（委托 user-service + 审计）。
func (s *AdminService) BanUser(ctx context.Context, adminID, id uint64) error {
	if _, err := s.users.BanUser(ctx, &userv1.BanUserRequest{Id: id}); err != nil {
		return err
	}
	return s.audit(ctx, adminID, "ban_user", "user", id, "")
}

// HideArticle 隐藏文章（委托 content-service + 审计）。
func (s *AdminService) HideArticle(ctx context.Context, adminID, id uint64) error {
	if _, err := s.articles.HideArticle(ctx, &contentv1.HideArticleRequest{Id: id}); err != nil {
		return err
	}
	return s.audit(ctx, adminID, "hide_article", "article", id, "")
}

// DeleteComment 删除评论（委托 interaction-service + 审计）。
func (s *AdminService) DeleteComment(ctx context.Context, adminID, id uint64) error {
	if _, err := s.comments.DeleteComment(ctx, &interactionv1.DeleteCommentRequest{Id: id}); err != nil {
		return err
	}
	return s.audit(ctx, adminID, "delete_comment", "comment", id, "")
}

// ListTags 查询标签（委托 content-service）。
func (s *AdminService) ListTags(ctx context.Context) ([]*contentv1.Tag, error) {
	resp, err := s.articles.ListTags(ctx, &contentv1.ListTagsRequest{})
	if err != nil {
		return nil, err
	}
	return resp.GetTags(), nil
}

// ListSensitiveWords 查询敏感词（委托 moderation-service）。
func (s *AdminService) ListSensitiveWords(ctx context.Context) ([]*moderationv1.SensitiveWord, error) {
	resp, err := s.sensitiveWords.ListSensitiveWords(ctx, &moderationv1.ListSensitiveWordsRequest{})
	if err != nil {
		return nil, err
	}
	return resp.GetWords(), nil
}

// CreateSensitiveWord 添加敏感词（委托 moderation-service + 审计）。
func (s *AdminService) CreateSensitiveWord(ctx context.Context, adminID uint64, word, level string) error {
	if _, err := s.sensitiveWords.CreateSensitiveWord(ctx, &moderationv1.CreateSensitiveWordRequest{Word: word, Level: level}); err != nil {
		return err
	}
	return s.audit(ctx, adminID, "create_sensitive_word", "sensitive_word", 0, word)
}

// ListReports 查询举报（委托 moderation-service）。
func (s *AdminService) ListReports(ctx context.Context, limit int) ([]*moderationv1.Report, error) {
	resp, err := s.reports.ListReports(ctx, &moderationv1.ListReportsRequest{Limit: int32(limit)})
	if err != nil {
		return nil, err
	}
	return resp.GetReports(), nil
}

// HandleReport 处理举报（委托 moderation-service + 审计）。
func (s *AdminService) HandleReport(ctx context.Context, adminID, id uint64, action string) error {
	if _, err := s.reports.HandleReport(ctx, &moderationv1.HandleReportRequest{Id: id, Action: action}); err != nil {
		return err
	}
	return s.audit(ctx, adminID, "handle_report", "report", id, action)
}

// ListAuditLogs 查询审计日志。
func (s *AdminService) ListAuditLogs(ctx context.Context, limit int) ([]model.AuditLog, error) {
	return s.repo.ListAuditLogs(ctx, limit)
}

// Statistics 是数据统计结果。
type Statistics struct {
	UserCount          int64
	ReportCount        int64
	TagCount           int64
	SensitiveWordCount int64
	AuditLogCount      int64
}

// Statistics 聚合各服务数据统计。
func (s *AdminService) Statistics(ctx context.Context) (*Statistics, error) {
	const big = 1000
	stats := &Statistics{}

	users, err := s.ListUsers(ctx, big, 0)
	if err != nil {
		return nil, err
	}
	stats.UserCount = int64(len(users))

	reports, err := s.ListReports(ctx, big)
	if err != nil {
		return nil, err
	}
	stats.ReportCount = int64(len(reports))

	tags, err := s.ListTags(ctx)
	if err != nil {
		return nil, err
	}
	stats.TagCount = int64(len(tags))

	words, err := s.ListSensitiveWords(ctx)
	if err != nil {
		return nil, err
	}
	stats.SensitiveWordCount = int64(len(words))

	logs, err := s.ListAuditLogs(ctx, big)
	if err != nil {
		return nil, err
	}
	stats.AuditLogCount = int64(len(logs))

	return stats, nil
}

func (s *AdminService) audit(ctx context.Context, adminID uint64, action, targetType string, targetID uint64, detail string) error {
	return s.repo.CreateAuditLog(ctx, &model.AuditLog{
		AdminID: adminID, Action: action, TargetType: targetType, TargetID: targetID, Detail: detail,
	})
}
