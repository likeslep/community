package kafkax

import "github.com/likeslep/community/pkg/apperr"

// RetryClass 描述错误的可重试性（plan.md §43）。
type RetryClass int

const (
	ClassSuccess   RetryClass = iota // 无错误
	ClassPermanent                   // 预期业务错误：不重试、不进 DLQ，直接提交
	ClassRetryable                   // 系统/依赖/超时错误：可重试
)

// Classify 根据错误返回其重试类别。
func Classify(err error) RetryClass {
	if err == nil {
		return ClassSuccess
	}
	if e := apperr.As(err); e != nil && !e.Retryable {
		return ClassPermanent
	}
	return ClassRetryable
}

// DefaultMaxRetries 初始重试上限（plan.md §21）。
const DefaultMaxRetries = 3

// Action 描述消息处理后的下一步动作。
type Action int

const (
	ActionCommit Action = iota // 提交（成功或业务拒绝）
	ActionRetry                // 投递到 retry topic
	ActionDLQ                  // 投递到 DLQ
)

// NextAction 根据重试次数与错误类别决定下一步动作。
func NextAction(retryCount int, class RetryClass) Action {
	switch class {
	case ClassSuccess, ClassPermanent:
		return ActionCommit
	case ClassRetryable:
		if retryCount >= DefaultMaxRetries {
			return ActionDLQ
		}
		return ActionRetry
	}
	return ActionCommit
}

// RetryTopic 返回对应 topic 的重试 topic 名。
func RetryTopic(topic string) string { return topic + ".retry" }

// DLQTopic 返回对应 topic 的死信 topic 名。
func DLQTopic(topic string) string { return topic + ".dlq" }
