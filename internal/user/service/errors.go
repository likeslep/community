package service

import "github.com/likeslep/community/pkg/apperr"

// user-service 错误码（100xx，plan.md §14）。
const (
	errCodeUsernameTaken   = apperr.CodeUser + 1 // 10001
	errCodeEmailTaken      = apperr.CodeUser + 2 // 10002
	errCodeUserNotFound    = apperr.CodeUser + 3 // 10003
	errCodeInvalidPassword = apperr.CodeUser + 4 // 10004
	errCodeInvalidInput    = apperr.CodeUser + 5 // 10005
)

// 预定义错误。
var (
	ErrUsernameTaken   = apperr.New(errCodeUsernameTaken, "用户名已被占用", apperr.WithHTTP(409))
	ErrEmailTaken      = apperr.New(errCodeEmailTaken, "邮箱已被注册", apperr.WithHTTP(409))
	ErrUserNotFound    = apperr.New(errCodeUserNotFound, "用户不存在", apperr.WithHTTP(404))
	ErrInvalidPassword = apperr.New(errCodeInvalidPassword, "用户名或密码错误", apperr.WithHTTP(401))
)
