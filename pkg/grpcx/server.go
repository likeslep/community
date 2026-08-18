// Package grpcx 提供 gRPC 服务端/客户端底座与业务错误映射（plan.md §10、§40）。
package grpcx

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Server 封装 gRPC 服务端，自动注册标准 health 服务。
type Server struct {
	server *grpc.Server
	health *health.Server
}

// NewServer 构造 gRPC 服务端（含 OpenTelemetry 追踪）。
func NewServer(opts ...grpc.ServerOption) *Server {
	baseOpts := []grpc.ServerOption{grpc.StatsHandler(otelgrpc.NewServerHandler())}
	baseOpts = append(baseOpts, opts...)
	s := grpc.NewServer(baseOpts...)
	hs := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, hs)
	return &Server{server: s, health: hs}
}

// GRPC 返回底层 *grpc.Server，供注册各业务 gRPC 服务。
func (s *Server) GRPC() *grpc.Server { return s.server }

// SetServing 将整体健康状态标记为 SERVING。
func (s *Server) SetServing() {
	s.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
}
