// Package middleware 提供 gateway 的认证中间件。
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/likeslep/community/pkg/auth"
)

// gin context 键，供下游 handler 读取认证身份。
const (
	CtxUserID = "user_id"
	CtxRole   = "role"
)

// Auth 校验 JWT 并将 user_id / role 注入 gin context（plan.md §11）。
func Auth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c.GetHeader("Authorization"))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未提供 token"})
			return
		}
		claims, err := auth.Verify(secret, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "token 无效"})
			return
		}
		c.Set(CtxUserID, claims.Subject)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

func extractToken(header string) string {
	if token, ok := strings.CutPrefix(header, "Bearer "); ok {
		return token
	}
	return ""
}

// RequireAdmin 校验当前用户角色为 admin（plan.md §12）。
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString(CtxRole) != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "message": "需要管理员权限"})
			return
		}
		c.Next()
	}
}
