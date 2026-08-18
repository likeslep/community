// Package metrics 提供 Prometheus 指标定义（plan.md §30）。
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HTTP 指标（plan.md §30）。
var (
	HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "HTTP 请求总数",
	}, []string{"service", "method", "path", "status"})

	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP 请求耗时（秒）",
		Buckets: prometheus.DefBuckets,
	}, []string{"service", "method", "path"})

	HTTPRequestErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_errors_total",
		Help: "HTTP 请求错误总数",
	}, []string{"service", "method", "path", "status"})
)

// 业务指标（plan.md §30）。
var (
	UserRegisterTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_register_total", Help: "用户注册总数",
	})
	ArticlePublishTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "article_publish_total", Help: "文章发布总数",
	})
	QuestionPublishTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "question_publish_total", Help: "问题发布总数",
	})
	AnswerCreateTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "answer_create_total", Help: "回答创建总数",
	})
	CommentCreateTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "comment_create_total", Help: "评论创建总数",
	})
	LikeTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "like_total", Help: "点赞总数",
	})
	CollectionTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "collection_total", Help: "收藏总数",
	})
	ReportTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "report_total", Help: "举报总数",
	})
)

// Handler 返回 /metrics 的 HTTP handler。
func Handler() http.Handler {
	return promhttp.Handler()
}
