// Package apperr 提供统一错误类型与各服务独立错误码空间（plan.md §14、§43）。
package apperr

import (
	"errors"
	"fmt"
)

// Kind 描述错误的类别，用于区分预期业务错误与系统故障（plan.md §43）。
type Kind string

const (
	KindBusiness   Kind = "business"   // 预期内的业务错误
	KindSystem     Kind = "system"     // 未预期的系统错误
	KindDependency Kind = "dependency" // 外部依赖故障
	KindTimeout    Kind = "timeout"    // 超时
)

// 各服务错误码前缀（plan.md §14），服务在此之上累加自己的错误码。
const (
	CodeUser        = 10000
	CodeContent     = 20000
	CodeInteraction = 30000
	CodeSocial      = 40000
	CodeModeration  = 50000
	CodeFile        = 60000
)

// Option 用于在构造 Error 时设置可选属性。
type Option func(*Error)

// WithHTTP 设置错误对应的 HTTP 状态码。
func WithHTTP(status int) Option {
	return func(e *Error) { e.HTTP = status }
}

// WithKind 设置错误类别。
func WithKind(k Kind) Option {
	return func(e *Error) { e.Kind = k }
}

// WithRetryable 标记错误是否可重试。
func WithRetryable(r bool) Option {
	return func(e *Error) { e.Retryable = r }
}

// Error 是本项目的统一错误类型。
type Error struct {
	Code      int
	Message   string
	HTTP      int
	Kind      Kind
	Retryable bool
	cause     error
}

func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 支持 errors.Is / errors.As 沿错误链回溯。
func (e *Error) Unwrap() error { return e.cause }

// New 创建一个业务错误，默认类别为 business、HTTP 500。
func New(code int, msg string, opts ...Option) *Error {
	e := &Error{Code: code, Message: msg, HTTP: 500, Kind: KindBusiness}
	for _, o := range opts {
		o(e)
	}
	return e
}

// Wrap 将底层错误包装为业务错误。
func Wrap(cause error, code int, msg string, opts ...Option) *Error {
	e := New(code, msg, opts...)
	e.cause = cause
	return e
}

// As 返回错误链中的 *Error，未找到时返回 nil。
func As(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return nil
}

// IsCode 判断错误链中是否存在指定错误码。
func IsCode(err error, code int) bool {
	e := As(err)
	return e != nil && e.Code == code
}

// KindOf 返回错误类别，非业务错误视为系统错误。
func KindOf(err error) Kind {
	if e := As(err); e != nil {
		return e.Kind
	}
	return KindSystem
}
