package service

import "github.com/likeslep/community/pkg/apperr"

// social-service 错误码（400xx，plan.md §14）。
const (
	errCodeSelfFollow = apperr.CodeSocial + 1 // 40001
)

// 预定义错误。
var (
	ErrSelfFollow = apperr.New(errCodeSelfFollow, "不能关注自己", apperr.WithHTTP(400))
)
