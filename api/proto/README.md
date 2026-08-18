# Protobuf 约定

内部服务通信使用 gRPC + Protobuf（plan.md §10、§40）。

## 目录

- `api/proto/<service>/`：各服务的 proto 定义。
- 每个服务一个 proto package，如 `package community.user.v1`。

## 稳定性规则（plan.md §40）

- 不随意复用 field number。
- 删除字段时保留 `reserved`。
- 不随意修改已有字段语义。
- 新增字段使用递增的 field number。

## 代码生成

使用 `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc` 生成到 `api/gen/`。

> 具体业务 proto 在 Phase 1+ 随各服务定义；Phase 0 仅使用 gRPC 标准 health 服务（见 `pkg/grpcx`）。
