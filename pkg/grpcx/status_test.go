package grpcx

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/likeslep/community/pkg/apperr"
)

func TestToStatus(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode codes.Code
	}{
		{"nil", nil, codes.OK},
		{"业务404", apperr.New(20001, "not found", apperr.WithHTTP(404)), codes.NotFound},
		{"业务401", apperr.New(10001, "unauthorized", apperr.WithHTTP(401)), codes.Unauthenticated},
		{"系统错误", errors.New("boom"), codes.Internal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToStatus(tt.err)
			if tt.err == nil {
				if got != nil {
					t.Fatalf("ToStatus(nil) = %v, want nil", got)
				}
				return
			}
			st, ok := status.FromError(got)
			if !ok {
				t.Fatal("ToStatus 未产生 gRPC status")
			}
			if st.Code() != tt.wantCode {
				t.Fatalf("code = %v, want %v", st.Code(), tt.wantCode)
			}
		})
	}
}

func TestFromStatus(t *testing.T) {
	e := FromStatus(status.Error(codes.NotFound, "x"))
	if !apperr.IsCode(e, 404) {
		t.Fatalf("FromStatus 应转换为 code=404 的 apperr，got %v", e)
	}
}
