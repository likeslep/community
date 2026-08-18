# Community Platform

面向程序员与技术学习者的多用户知识学习社区（问答 + 文章 + 社交 + 推荐 + 搜索 + 治理）。

## 文档

| 文档 | 作用 |
|------|------|
| [`docs/sdd/spec.md`](docs/sdd/spec.md) | 产品意图、功能边界、业务规则 |
| [`docs/sdd/plan.md`](docs/sdd/plan.md) | 技术实现规划（架构、服务边界、Phase） |
| [`docs/sdd/tasks.md`](docs/sdd/tasks.md) | 可执行与验收的工程任务拆分 |

优先级：`spec.md` → `plan.md` → `tasks.md`。

## 技术栈

Go · Gin · gRPC · Protobuf · GORM · MySQL8 · Redis · Kafka(KRaft) · Elasticsearch · OpenTelemetry · Prometheus · Grafana · Jaeger · Zap · Docker。

## 目录结构

```
cmd/                  各服务 main 入口
internal/<service>/   各服务私有代码
pkg/                  跨服务共享库
api/proto/            Protobuf 定义
migrations/           各服务数据库迁移
deploy/               Docker 与编排
docs/sdd/             规范文档
```

## 常用命令

```text
make build   编译
make test    测试
make vet     go vet
make lint    golangci-lint（需先安装 golangci-lint）
make fmt     格式化（gofmt + goimports）
```

## 本地启动

```text
docker compose up
```

> 完整环境编排见 Phase 0 的 `deploy/compose`。
