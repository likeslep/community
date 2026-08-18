package service

import "github.com/likeslep/community/pkg/apperr"

// moderation-service 错误码（500xx，plan.md §14）。
const (
	errCodeTaskNotFound   = apperr.CodeModeration + 1 // 50001
	errCodeReportNotFound = apperr.CodeModeration + 2 // 50002
	errCodeIllegalState   = apperr.CodeModeration + 3 // 50003
)

// 预定义错误。
var (
	ErrTaskNotFound   = apperr.New(errCodeTaskNotFound, "审核任务不存在", apperr.WithHTTP(404))
	ErrReportNotFound = apperr.New(errCodeReportNotFound, "举报不存在", apperr.WithHTTP(404))
	ErrIllegalState   = apperr.New(errCodeIllegalState, "当前状态不允许该操作", apperr.WithHTTP(409))
)
