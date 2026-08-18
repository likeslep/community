package service

import "github.com/likeslep/community/pkg/apperr"

// content-service 错误码（200xx，plan.md §14）。
const (
	errCodeArticleNotFound  = apperr.CodeContent + 1 // 20001
	errCodeForbidden        = apperr.CodeContent + 2 // 20002
	errCodeIllegalState     = apperr.CodeContent + 3 // 20003
	errCodeInvalidInput     = apperr.CodeContent + 4 // 20004
	errCodeAnswerNotFound   = apperr.CodeContent + 5 // 20005
	errCodeAnswerMismatch   = apperr.CodeContent + 6 // 20006
	errCodeQuestionNotFound = apperr.CodeContent + 7 // 20007
)

// 预定义错误。
var (
	ErrArticleNotFound  = apperr.New(errCodeArticleNotFound, "文章不存在", apperr.WithHTTP(404))
	ErrForbidden        = apperr.New(errCodeForbidden, "无权限操作该资源", apperr.WithHTTP(403))
	ErrIllegalState     = apperr.New(errCodeIllegalState, "当前状态不允许该操作", apperr.WithHTTP(409))
	ErrAnswerNotFound   = apperr.New(errCodeAnswerNotFound, "回答不存在", apperr.WithHTTP(404))
	ErrAnswerMismatch   = apperr.New(errCodeAnswerMismatch, "回答不属于该问题", apperr.WithHTTP(400))
	ErrQuestionNotFound = apperr.New(errCodeQuestionNotFound, "问题不存在", apperr.WithHTTP(404))
)
