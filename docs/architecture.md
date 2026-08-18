# 架构约定（服务模板）

本文档定义所有服务共同遵循的分层与目录约定（对应 `tasks.md` 的 `TASK-FND-002`）。

## 分层

每个业务服务内部遵循统一分层（plan.md §9）：

```text
Handler (transport)
    |
    v
Service (业务逻辑)
    |
    v
Repository Interface (接口，定义在消费方)
    |
    v
GORM Repository (实现)
    |
    v
MySQL
```

## 规则

1. **业务层禁止直接依赖 GORM**：`service` 包不得 import GORM，只依赖 Repository 接口。
2. **Repository 接口定义在消费方**（service 包内），实现放在 repository 包 —— 遵循「accept interfaces, return structs」。
3. **避免过度分层**：不引入 Controller/Application/Domain/Port/Adapter/DAO 等多层抽象（plan.md §9）。
4. **简单 CRUD 用 GORM**；复杂/性能敏感查询允许显式 SQL。
5. **依赖注入用构造函数**：不引入 DI 容器，依赖通过 `NewXxx(deps...)` 显式传递。

## 单服务目录约定

以 `user` 服务为例：

```text
cmd/user/main.go           入口：加载配置 → 构造依赖 → 启动服务
internal/user/
  handler/                 HTTP / gRPC 处理器（transport 层，不含业务）
  service/                 业务逻辑 + Repository 接口定义
  repository/              GORM 实现
  model/                   领域模型 / GORM 实体
```

## 共享库

跨服务复用的能力放在 `pkg/`：

| 包 | 职责 |
|----|------|
| `pkg/config` | 环境变量配置加载 |
| `pkg/apperr` | 统一错误类型与错误码 |
| `pkg/logger` | Zap 结构化日志 |
| `pkg/middleware` | HTTP/gRPC 中间件（Request/Trace ID 等） |
| `pkg/grpcx` | gRPC 服务端/客户端底座 |
| `pkg/kafkax` | Kafka 生产者/消费者底座 |
| `pkg/redisx` | Redis 客户端封装 |
| `pkg/version` | 构建版本信息 |

`internal/<service>/` 只存放该服务私有代码，服务之间不得直接引用对方的 internal 包。
