package kafkax

import (
	"errors"
	"testing"

	"github.com/likeslep/community/pkg/apperr"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want RetryClass
	}{
		{"成功", nil, ClassSuccess},
		{"业务错误", apperr.New(10001, "参数非法", apperr.WithRetryable(false)), ClassPermanent},
		{"系统错误", errors.New("db down"), ClassRetryable},
		{"依赖可重试", apperr.Wrap(errors.New("timeout"), 20001, "x", apperr.WithRetryable(true)), ClassRetryable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Classify(tt.err); got != tt.want {
				t.Fatalf("Classify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextAction(t *testing.T) {
	tests := []struct {
		name  string
		count int
		class RetryClass
		want  Action
	}{
		{"成功提交", 0, ClassSuccess, ActionCommit},
		{"业务提交", 0, ClassPermanent, ActionCommit},
		{"首次重试", 0, ClassRetryable, ActionRetry},
		{"第二次重试", 1, ClassRetryable, ActionRetry},
		{"第三次重试", 2, ClassRetryable, ActionRetry},
		{"超限进 DLQ", 3, ClassRetryable, ActionDLQ},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextAction(tt.count, tt.class); got != tt.want {
				t.Fatalf("NextAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTopics(t *testing.T) {
	if RetryTopic("user.events") != "user.events.retry" {
		t.Fatalf("RetryTopic = %q", RetryTopic("user.events"))
	}
	if DLQTopic("user.events") != "user.events.dlq" {
		t.Fatalf("DLQTopic = %q", DLQTopic("user.events"))
	}
}
