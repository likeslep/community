package grpcx

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Dial 建立到目标服务的 gRPC 连接（内部通信用，默认明文，后续可换 mTLS）。
// 连接为惰性建立；各 RPC 调用时传入自己的 context 控制超时。
func Dial(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	dialOpts = append(dialOpts, opts...)
	return grpc.NewClient(addr, dialOpts...)
}
