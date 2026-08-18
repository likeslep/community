package grpcx

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/likeslep/community/pkg/apperr"
)

// ToStatus 将业务/系统错误映射为 gRPC status 错误，用于服务端拦截器统一返回。
func ToStatus(err error) error {
	if err == nil {
		return nil
	}
	if e := apperr.As(err); e != nil {
		return status.Error(httpToCode(e.HTTP), e.Message)
	}
	return status.Error(codes.Internal, err.Error())
}

// FromStatus 将 gRPC status 错误转换回 apperr，用于客户端拦截器统一处理。
// 注意：gRPC status 不携带业务错误码空间（如 200xx），此处用 HTTP 语义近似还原。
func FromStatus(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	return apperr.New(codeToHTTP(st.Code()), st.Message())
}

func httpToCode(s int) codes.Code {
	switch s {
	case 400:
		return codes.InvalidArgument
	case 401:
		return codes.Unauthenticated
	case 403:
		return codes.PermissionDenied
	case 404:
		return codes.NotFound
	case 409:
		return codes.AlreadyExists
	case 429:
		return codes.ResourceExhausted
	case 503:
		return codes.Unavailable
	case 504:
		return codes.DeadlineExceeded
	default:
		if s >= 500 {
			return codes.Internal
		}
		return codes.Unknown
	}
}

func codeToHTTP(c codes.Code) int {
	switch c {
	case codes.InvalidArgument:
		return 400
	case codes.Unauthenticated:
		return 401
	case codes.PermissionDenied:
		return 403
	case codes.NotFound:
		return 404
	case codes.AlreadyExists:
		return 409
	case codes.ResourceExhausted:
		return 429
	case codes.Unavailable:
		return 503
	case codes.DeadlineExceeded:
		return 504
	default:
		return 500
	}
}

// HTTPStatus 从错误（apperr 或 gRPC status）中提取 HTTP 状态码。
func HTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if e := apperr.As(err); e != nil && e.HTTP != 0 {
		return e.HTTP
	}
	if st, ok := status.FromError(err); ok {
		return codeToHTTP(st.Code())
	}
	return http.StatusInternalServerError
}

// Message 从错误中提取面向用户的消息。
func Message(err error) string {
	if err == nil {
		return "success"
	}
	if e := apperr.As(err); e != nil {
		return e.Message
	}
	if st, ok := status.FromError(err); ok {
		return st.Message()
	}
	return err.Error()
}
