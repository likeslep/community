package auth

import (
	"testing"
	"time"
)

func TestSignAndVerify(t *testing.T) {
	secret := []byte("test-secret")
	token, err := Sign(secret, "123", "alice", "author", time.Hour)
	if err != nil {
		t.Fatalf("Sign() err = %v", err)
	}
	claims, err := Verify(secret, token)
	if err != nil {
		t.Fatalf("Verify() err = %v", err)
	}
	if claims.Subject != "123" || claims.Username != "alice" || claims.Role != "author" {
		t.Fatalf("claims = %+v", claims)
	}
}

func TestVerifyInvalid(t *testing.T) {
	secret := []byte("test-secret")
	otherToken, _ := Sign([]byte("other-secret"), "1", "a", "author", time.Hour)

	tests := []struct {
		name  string
		token string
	}{
		{"空token", ""},
		{"乱码", "not.a.token"},
		{"错误密钥", otherToken},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Verify(secret, tt.token); err == nil {
				t.Fatal("期望校验失败，实际通过")
			}
		})
	}
}
