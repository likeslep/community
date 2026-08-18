// Package rbac 提供角色与权限模型（plan.md §12：RBAC + Resource Ownership）。
package rbac

// 角色（与 user-service model 保持一致）。
const (
	RoleAuthor    = "author"
	RoleModerator = "moderator"
	RoleAdmin     = "admin"
)

// 权限点（plan.md §12）。
const (
	PermArticleCreate    = "article:create"
	PermArticleUpdateOwn = "article:update:own"
	PermArticleUpdateAny = "article:update:any"
	PermArticleDeleteOwn = "article:delete:own"
	PermArticleDeleteAny = "article:delete:any"
	PermArticleModerate  = "article:moderate"
	PermUserBan          = "user:ban"
	PermTagManage        = "tag:manage"
)

// rolePermissions 是角色 → 权限集合的映射。管理员拥有全部权限。
var rolePermissions = map[string]map[string]struct{}{
	RoleAuthor: set(
		PermArticleCreate,
		PermArticleUpdateOwn,
		PermArticleDeleteOwn,
	),
	RoleModerator: set(
		PermArticleCreate,
		PermArticleUpdateOwn,
		PermArticleDeleteOwn,
		PermArticleUpdateAny,
		PermArticleDeleteAny,
		PermArticleModerate,
	),
	RoleAdmin: set(
		PermArticleCreate,
		PermArticleUpdateOwn,
		PermArticleDeleteOwn,
		PermArticleUpdateAny,
		PermArticleDeleteAny,
		PermArticleModerate,
		PermUserBan,
		PermTagManage,
	),
}

// HasPermission 判断角色是否拥有权限。
func HasPermission(role, perm string) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}
	_, ok = perms[perm]
	return ok
}

func set(perms ...string) map[string]struct{} {
	m := make(map[string]struct{}, len(perms))
	for _, p := range perms {
		m[p] = struct{}{}
	}
	return m
}
