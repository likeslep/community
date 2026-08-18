# Failure / Reliability 测试（Phase 13）

本文档描述可靠性验证方法（plan.md §33 Phase 13）。

## 已通过单元测试覆盖的能力

| 能力 | 测试位置 |
|------|---------|
| 消费者错误分类（业务/系统/重试） | `pkg/kafkax/retry_test.go` |
| 重试 / DLQ 决策（3 次上限） | `pkg/kafkax/retry_test.go` |
| 幂等消费（内存 + DB 实现） | `pkg/kafkax/idempotency_test.go`、`pkg/outbox/processed.go` |
| 点赞并发幂等（唯一约束） | `internal/interaction/repository/gorm.go` |
| 优雅停机 | `pkg/server`（真实冒烟验证 exit 0） |

## 故障注入场景（需集成环境验证）

在 `docker compose up` 后，用以下方式注入故障并观察恢复：

1. **Kafka 不可用**：`docker compose stop kafka` → 各服务 outbox 事件积压，Kafka 恢复后自动投递。
2. **Redis 不可用**：`docker compose stop redis` → 限流器降级放行（`ratelimit.go` 中 `err != nil` 时 `c.Next()`）。
3. **MySQL 超时**：依赖 `context` 超时，消费者按 `Retryable` 分类重试。
4. **重复事件**：通过 `processed_events` 表（DB 幂等）或唯一约束去重。
5. **乱序事件**：文章状态机守卫（只有 `PENDING_REVIEW` 接受审核结果）。

## 定义（Definition of Done）

- 故障注入后服务不 panic、不雪崩。
- 恢复后数据最终一致。
- 重复/乱序事件不产生错误状态。
