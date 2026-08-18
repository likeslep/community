package apperr

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		code int
		msg  string
	}{
		{"user-not-found", CodeUser + 1, "用户不存在"},
		{"content-invalid", CodeContent + 2, "内容非法"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(tt.code, tt.msg)
			if e.Code != tt.code || e.Message != tt.msg {
				t.Fatalf("New() = %+v", e)
			}
			if e.Kind != KindBusiness {
				t.Fatalf("默认 Kind 应为 business，实际 %q", e.Kind)
			}
		})
	}
}

func TestWrapUnwrap(t *testing.T) {
	base := errors.New("db down")
	e := Wrap(base, CodeContent, "查询失败", WithKind(KindDependency), WithRetryable(true))
	if !errors.Is(e, base) {
		t.Fatal("Wrap 后 errors.Is 应命中底层错误")
	}
	if !e.Retryable || e.Kind != KindDependency {
		t.Fatalf("Wrap opts 未生效: %+v", e)
	}
}

func TestIsCode(t *testing.T) {
	e := New(CodeUser+42, "x")
	tests := []struct {
		name string
		err  error
		code int
		want bool
	}{
		{"命中", e, CodeUser + 42, true},
		{"未命中", e, CodeUser + 1, false},
		{"包装后命中", fmt.Errorf("wrap: %w", e), CodeUser + 42, true},
		{"普通错误", errors.New("plain"), CodeUser + 42, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCode(tt.err, tt.code); got != tt.want {
				t.Fatalf("IsCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKindOf(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want Kind
	}{
		{"业务错误", New(CodeUser, "x"), KindBusiness},
		{"系统错误", Wrap(errors.New("panic"), CodeUser, "x", WithKind(KindSystem)), KindSystem},
		{"普通错误", errors.New("plain"), KindSystem},
		{"nil", nil, KindSystem},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KindOf(tt.err); got != tt.want {
				t.Fatalf("KindOf() = %q, want %q", got, tt.want)
			}
		})
	}
}
