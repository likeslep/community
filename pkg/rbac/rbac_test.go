package rbac

import "testing"

func TestHasPermission(t *testing.T) {
	tests := []struct {
		role, perm string
		want       bool
	}{
		{RoleAuthor, PermArticleCreate, true},
		{RoleAuthor, PermArticleUpdateOwn, true},
		{RoleAuthor, PermArticleUpdateAny, false},
		{RoleAuthor, PermUserBan, false},
		{RoleModerator, PermArticleUpdateAny, true},
		{RoleModerator, PermArticleModerate, true},
		{RoleModerator, PermUserBan, false},
		{RoleAdmin, PermUserBan, true},
		{RoleAdmin, PermTagManage, true},
		{"unknown-role", PermArticleCreate, false},
	}
	for _, tt := range tests {
		t.Run(tt.role+"/"+tt.perm, func(t *testing.T) {
			if got := HasPermission(tt.role, tt.perm); got != tt.want {
				t.Fatalf("HasPermission(%q,%q) = %v, want %v", tt.role, tt.perm, got, tt.want)
			}
		})
	}
}
