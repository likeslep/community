package service

import "github.com/likeslep/community/pkg/apperr"

// file-service 错误码（600xx，plan.md §14）。
const (
	errCodeFileNotFound = apperr.CodeFile + 1 // 60001
	errCodeTooLarge     = apperr.CodeFile + 2 // 60002
	errCodeInvalidType  = apperr.CodeFile + 3 // 60003
)

// 预定义错误。
var (
	ErrFileNotFound = apperr.New(errCodeFileNotFound, "文件不存在", apperr.WithHTTP(404))
	ErrTooLarge     = apperr.New(errCodeTooLarge, "文件超出大小限制", apperr.WithHTTP(400))
	ErrInvalidType  = apperr.New(errCodeInvalidType, "不支持的文件类型", apperr.WithHTTP(400))
)
