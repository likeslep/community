package service

import "github.com/likeslep/community/pkg/apperr"

// interaction-service 错误码（300xx，plan.md §14）。
const (
	errCodeInvalidTarget = apperr.CodeInteraction + 1 // 30001
)

// 预定义错误。
var (
	ErrInvalidTarget = apperr.New(errCodeInvalidTarget, "非法的互动目标类型", apperr.WithHTTP(400))
)
