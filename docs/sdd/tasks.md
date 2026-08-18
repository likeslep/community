# Community Platform - Task Breakdown

Version: 1.0
Status: Draft

---

## 1. Document Purpose

本文档将 `plan.md` 的 Phase 0 ~ 15 拆分为可由 AI Agent **独立执行与验收**的工程任务。

文档优先级（Agent 必须遵循）：

```text
spec.md
   ↓
plan.md
   ↓
tasks.md
```

- `spec.md` —— 定义系统「做什么」（产品意图、功能边界、业务规则）
- `plan.md` —— 定义系统「怎么实现」（技术栈、架构、服务边界、Phase）
- `tasks.md`（本文档）—— 定义「拆成哪些任务、谁依赖谁、如何验收」

Agent 不得自行修改系统架构；不得新增未授权的 Service / Database / Broker / Framework；冲突时停止并报告。

---

## 2. Conventions

### 2.1 Task ID 方案

```text
TASK-<DOMAIN>-<NNN>
```

| Domain | 含义 | 对应 Phase |
|--------|------|-----------|
| FND | Foundation | Phase 0 |
| USER | User & Authentication | Phase 1 |
| ARTICLE | Article | Phase 2 |
| QA | Question & Answer | Phase 3 |
| MOD | Moderation | Phase 4 |
| INTERACTION | Interaction | Phase 5 |
| SOCIAL | Social | Phase 6 |
| NOTIF | Notification | Phase 7 |
| SEARCH | Search | Phase 8 |
| FEED | Feed | Phase 9 |
| FILE | File | Phase 10 |
| ADMIN | Admin | Phase 11 |
| OBS | Observability | Phase 12 |
| REL | Reliability | Phase 13 |
| PERF | Performance | Phase 14 |
| HARDEN | Production Hardening | Phase 15 |

### 2.2 Task 字段

每个 Task 必须包含（源自 plan.md §36）：

| 字段 | 说明 |
|------|------|
| Objective | 任务目标，一句话 |
| Context | 背景与前置上下文 |
| Dependencies | 显式依赖的 Task ID 列表 |
| Scope | 明确边界：做什么 / 不做什么 |
| Implementation Requirements | 实现必须满足的要求（引用 plan.md 相关章节） |
| Acceptance Criteria | 可验证的验收标准 |
| Tests | 需要覆盖的测试类型与场景 |
| Definition of Done | 完成定义 |

### 2.3 依赖与并行

- `Dependencies` 为空 = 无前置，可立即开始（需 Phase 0 完成）。
- 无依赖关系的 Task 可并行执行。
- Agent 不得假设尚未完成的 Task 已存在（plan.md §37）。
- Phase 表示能力依赖，不强制串行（plan.md §34）。跨 Phase 的依赖在具体 Task 的 `Dependencies` 中显式声明。

### 2.4 状态约定

所有 Task 初始状态为 `Not Started`。执行过程中标记：`Not Started` → `In Progress` → `Done`。

---

## 3. 任务总览

### 3.1 索引

| Task ID | 标题 | 所属服务 | 依赖 |
|---------|------|---------|------|
| TASK-FND-001 | Repository 结构（Monorepo） | 全部 | — |
| TASK-FND-002 | 服务模板（分层骨架） | 全部 | FND-001 |
| TASK-FND-003 | 配置系统 | 全部 | FND-002 |
| TASK-FND-004 | 错误系统（错误码空间） | 全部 | FND-002 |
| TASK-FND-005 | 结构化日志（Zap） | 全部 | FND-002 |
| TASK-FND-006 | Request ID / Trace ID 上下文传播 | 全部 | FND-002 |
| TASK-FND-007 | gRPC 基础设施 | 全部 | FND-002 |
| TASK-FND-008 | Kafka 基础设施（Envelope / Producer / Consumer） | 全部 | FND-002 |
| TASK-FND-009 | MySQL Migration 框架 | 全部 | FND-002 |
| TASK-FND-010 | Redis Client 封装 | 全部 | FND-002 |
| TASK-FND-011 | Health Check + Graceful Shutdown | 全部 | FND-002 |
| TASK-FND-012 | Docker Compose 编排 | 全部 | FND-001 |
| TASK-USER-001 | 用户注册 + `user.created` Outbox | user-service | FND-004/005/006/007/008/009/010/011 |
| TASK-USER-002 | 用户登录 + JWT 签发 | user-service | USER-001 |
| TASK-USER-003 | 用户资料（Profile） | user-service | USER-001 |
| TASK-USER-004 | 角色与权限（RBAC） | user-service | USER-001 |
| TASK-USER-005 | Gateway 认证 / 鉴权中间件 | gateway-service | USER-002/004 |
| TASK-ARTICLE-001 | Article 领域模型 + 状态机 + Markdown | content-service | USER-005 |
| TASK-ARTICLE-002 | Article 草稿 CRUD | content-service | ARTICLE-001 |
| TASK-ARTICLE-003 | Article 版本管理 | content-service | ARTICLE-001 |
| TASK-ARTICLE-004 | Tag 实体与内容绑定 | content-service | ARTICLE-001 |
| TASK-ARTICLE-005 | Article 提交审核（Outbox） | content-service | ARTICLE-001/004 |
| TASK-ARTICLE-006 | Article 发布 / 驳回状态流转 | content-service | ARTICLE-005, MOD-004 |
| TASK-QA-001 | Question 领域模型 + 状态机 | content-service | USER-005 |
| TASK-QA-002 | Question CRUD + 发布（Outbox） | content-service | QA-001, ARTICLE-004 |
| TASK-QA-003 | Answer 领域模型 + CRUD | content-service | QA-002 |
| TASK-QA-004 | Accept Answer（采纳回答） | content-service | QA-003 |
| TASK-QA-005 | Question-Tag 绑定 | content-service | QA-001, ARTICLE-004 |
| TASK-MOD-001 | 敏感词库 + 检测 | moderation-service | FND-004/009 |
| TASK-MOD-002 | Moderation Task + 消费提交事件 | moderation-service | MOD-001, ARTICLE-005 |
| TASK-MOD-003 | 人工审核接口 | moderation-service | MOD-002 |
| TASK-MOD-004 | 审核结果流转（Approve/Reject/Hide） | moderation-service | MOD-003 |
| TASK-MOD-005 | Report 举报与处理 | moderation-service | MOD-001 |
| TASK-INTERACTION-001 | Comment 领域模型 + CRUD | interaction-service | USER-005, ARTICLE-001 |
| TASK-INTERACTION-002 | Like（幂等 + 并发） | interaction-service | INTERACTION-001 |
| TASK-INTERACTION-003 | Collection（幂等） | interaction-service | INTERACTION-001 |
| TASK-INTERACTION-004 | View 浏览计数 | interaction-service | INTERACTION-001 |
| TASK-INTERACTION-005 | 互动计数器一致性 | interaction-service | INTERACTION-002/003/004 |
| TASK-SOCIAL-001 | User Follow / Unfollow | social-service | USER-005 |
| TASK-SOCIAL-002 | Tag Follow | social-service | SOCIAL-001, ARTICLE-004 |
| TASK-SOCIAL-003 | Social Graph 查询 | social-service | SOCIAL-001/002 |
| TASK-NOTIF-001 | Notification 模型 + 事件消费 | notification-service | FND-008, USER-005 |
| TASK-NOTIF-002 | 通知列表 + 未读数 | notification-service | NOTIF-001 |
| TASK-NOTIF-003 | 已读 / 全部已读 | notification-service | NOTIF-002 |
| TASK-SEARCH-001 | ES 索引与 Mapping | search-service | FND-009, ARTICLE-001, QA-001 |
| TASK-SEARCH-002 | 索引同步（消费事件） | search-service | SEARCH-001, ARTICLE-006 |
| TASK-SEARCH-003 | 搜索 API（过滤/高亮/排序/分页） | search-service | SEARCH-002 |
| TASK-SEARCH-004 | Index Rebuild（全量重建） | search-service | SEARCH-001 |
| TASK-FEED-001 | Feed 数据模型 + Redis Sorted Set | feed-service | FND-010, ARTICLE-001 |
| TASK-FEED-002 | Fan-out on Write（普通作者） | feed-service | FEED-001, ARTICLE-006 |
| TASK-FEED-003 | Fan-out on Read（高关注作者） | feed-service | FEED-001, ARTICLE-006 |
| TASK-FEED-004 | Rule Ranking（抽象接口） | feed-service | FEED-001 |
| TASK-FEED-005 | Feed API（Following/Tag/Latest/Hot） | feed-service | FEED-002/003/004, SOCIAL-001/002 |
| TASK-FILE-001 | Storage 抽象 + LocalStorage | file-service | FND-002 |
| TASK-FILE-002 | 文件上传（校验 + 元数据） | file-service | FILE-001 |
| TASK-FILE-003 | 文件下载（流式 + 权限） | file-service | FILE-002 |
| TASK-ADMIN-001 | 用户管理 | admin-service | USER-001, FND-007 |
| TASK-ADMIN-002 | 内容管理 | admin-service | ARTICLE-006, QA-002, FND-007 |
| TASK-ADMIN-003 | 评论管理 | admin-service | INTERACTION-001, FND-007 |
| TASK-ADMIN-004 | 举报处理 | admin-service | MOD-005, FND-007 |
| TASK-ADMIN-005 | 标签管理 | admin-service | ARTICLE-004, FND-007 |
| TASK-ADMIN-006 | 敏感词管理 | admin-service | MOD-001, FND-007 |
| TASK-ADMIN-007 | 数据统计 | admin-service | FND-007 |
| TASK-ADMIN-008 | 审计日志 | admin-service | FND-005/007 |
| TASK-OBS-001 | OpenTelemetry 接入 | 全部 | FND-006/007/008 |
| TASK-OBS-002 | Metrics（HTTP/gRPC/Kafka/Business） | 全部 | OBS-001 |
| TASK-OBS-003 | Grafana / Jaeger 部署与 Dashboard | 全部 | OBS-002 |
| TASK-OBS-004 | 完整 Trace Demo | 全部 | OBS-003, ARTICLE-006, SEARCH-002 |
| TASK-REL-001 | 依赖不可用容错 | 全部 | Phase 7/8/9 完成 |
| TASK-REL-002 | 消费者重试 / 恢复 | 全部 | REL-001 |
| TASK-REL-003 | 重复 / 乱序事件验证 | 全部 | REL-001 |
| TASK-REL-004 | Failure 测试套件 | 全部 | REL-002/003 |
| TASK-PERF-001 | HTTP / gRPC 压测 | 全部 | Phase 8/9 完成 |
| TASK-PERF-002 | Kafka 吞吐压测 | 全部 | PERF-001 |
| TASK-PERF-003 | Feed / Search 压测 | 全部 | PERF-001 |
| TASK-PERF-004 | 压测报告与指标记录 | 全部 | PERF-002/003 |
| TASK-HARDEN-001 | 限流（Rate Limit） | gateway-service | FND-012, USER-005 |
| TASK-HARDEN-002 | 安全加固（注入/XSS/文件/JWT/密码） | 全部 | Phase 13/14 完成 |
| TASK-HARDEN-003 | 弹性（Timeout/Retry/Circuit Breaker） | 全部 | HARDEN-002 |
| TASK-HARDEN-004 | 敏感数据保护与审计 | 全部 | HARDEN-002 |

### 3.2 阶段级依赖图

```text
Phase 0 (FND)
   │
   ▼
Phase 1 (USER) ────────────────┬──────────────┬──────────────┐
   │                           │              │              │
   ▼                           ▼              │              │
Phase 2 (ARTICLE) ──┬─── Phase 3 (QA)         │              │
   │                │        │                │              │
   ├──────┬─────────┤        │                │              │
   ▼      │         ▼        │              Phase 6 (SOCIAL) │
Phase 4  │   Phase 5          │                │              │
(MOD)    │   (INTERACTION)    │                │              │
   │     │         │          │                │              │
   └─────┴─────────┴────┬─────┴────────────────┤              │
                        ▼                      │              │
                    Phase 7 (NOTIF) ◄──────────┘              │
                                                              │
Phase 2/3 ────────────────────────► Phase 8 (SEARCH)          │
Phase 2/3/5/6 ────────────────────► Phase 9 (FEED)            │
Phase 2 ──────────────────────────► Phase 10 (FILE)           │
Phase 1/4/5/6 ────────────────────► Phase 11 (ADMIN)          │
Phase 0 ──────────────────────────► Phase 12 (OBS)            │
Phase 7/8/9 ──────────────────────► Phase 13 (REL)            │
Phase 8/9 ────────────────────────► Phase 14 (PERF)           │
Phase 13/14 ──────────────────────► Phase 15 (HARDEN)         │
```

---

## 4. Phase 0 — Foundation

> 验收总目标（plan.md §33）：`docker compose up` 后所有基础设施 healthy；所有服务可启动；Gateway health 返回 200；gRPC health 返回 SERVING。

### TASK-FND-001 · Repository 结构（Monorepo）
- **Objective**：建立多服务 Go Monorepo 目录布局。
- **Context**：项目起点，无任何代码。需承载 11 个服务 + 共享库。
- **Dependencies**：—
- **Scope**：目录结构、go.mod / go.work 布局、共享库目录。不实现任何业务。
- **Implementation Requirements**：
  - 采用 `services/<name>/` 与 `internal/`、`pkg/`（或 `shared/`）分离的布局。
  - 共享库（config、errors、logging、middleware、grpc、kafka 等）放独立 module，可被各服务 import。
  - 遵循 Go 官方布局约定与 `golangci-lint` 检查。
- **Acceptance Criteria**：
  - 目录树可清晰区分 11 个服务与共享库。
  - `go build ./...` 在根目录通过。
  - `.golangci.yml` 就位且可运行。
- **Tests**：无（结构任务），以 `go build` + lint 通过为验证。
- **Definition of Done**：目录结构提交，根目录 `go build ./...` 与 lint 通过。

### TASK-FND-002 · 服务模板（分层骨架）
- **Objective**：提供一个所有服务复用的分层代码骨架。
- **Context**：plan.md §9 定义了 `Handler → Service → Repository Interface → GORM Repository → MySQL` 分层，业务层禁止直接依赖 GORM。
- **Dependencies**：`TASK-FND-001`
- **Scope**：分层骨架模板（handler/service/repository 接口与实现、启动入口、依赖注入方式）。不包含具体业务逻辑。
- **Implementation Requirements**：
  - 遵循 plan.md §9，Repository 用接口隔离业务与持久化。
  - 避免 plan.md §9 列出的过度多层抽象（Controller/Application/Domain/Port/Adapter/DAO）。
  - 提供 `cmd/<service>/main.go` 统一启动范式。
- **Acceptance Criteria**：任意新建服务可复制骨架并快速启动空服务。
- **Tests**：一个最小 smoke test（骨架可实例化）。
- **Definition of Done**：骨架模板可运行，含最小测试。

### TASK-FND-003 · 配置系统
- **Objective**：统一配置加载与环境变量管理。
- **Context**：所有服务需要一致的配置（DB/Redis/Kafka/gRPC/日志级别）。
- **Dependencies**：`TASK-FND-002`
- **Scope**：config 加载库、环境变量约定、各环境 profile。不实现业务配置项。
- **Implementation Requirements**：
  - 支持环境变量覆盖默认值。
  - 敏感信息（密码/密钥）通过环境变量注入，禁止硬编码。
  - 配置缺失时启动失败并给出清晰错误。
- **Acceptance Criteria**：服务可通过 env 覆盖默认配置启动。
- **Tests**：表格驱动测试覆盖默认值 / 覆盖 / 缺失报错。
- **Definition of Done**：配置库就位并有单元测试。

### TASK-FND-004 · 错误系统（错误码空间）
- **Objective**：建立统一错误类型与各服务独立错误码空间。
- **Context**：plan.md §14 定义 `User 100xx / Content 200xx / Interaction 300xx / Social 400xx / Moderation 500xx / File 600xx`，错误码必须稳定。
- **Dependencies**：`TASK-FND-002`
- **Scope**：统一 error 类型、错误码常量、错误码↔HTTP/gRPC 状态映射。不实现业务错误文案。
- **Implementation Requirements**：
  - 禁止用无语义错误字符串作为程序判断依据。
  - 区分 Expected Business Error / Unexpected System Error / Dependency Failure / Timeout / Retryable / Non-Retryable（plan.md §43）。
  - 错误需可携带 code、message、字段级信息。
- **Acceptance Criteria**：业务错误与系统错误可被统一包装与映射。
- **Tests**：错误码映射、error 包装与 `errors.Is/As` 表格驱动测试。
- **Definition of Done**：错误系统就位，含测试。

### TASK-FND-005 · 结构化日志（Zap）
- **Objective**：统一 Zap 结构化日志封装。
- **Context**：plan.md §28 要求日志统一 Zap，禁止生产代码 `fmt.Println()`。
- **Dependencies**：`TASK-FND-002`
- **Scope**：logger 构造、字段约定、level 管理。不含具体业务日志点。
- **Implementation Requirements**：
  - 每个重要请求至少含 `service / request_id / trace_id / user_id`。
  - 业务事件日志含 `event_type / aggregate_id`。
  - 支持 JSON 输出，适配日志采集。
- **Acceptance Criteria**：logger 输出结构化 JSON，含约定字段。
- **Tests**：字段注入与输出格式测试。
- **Definition of Done**：日志封装就位，含测试。

### TASK-FND-006 · Request ID / Trace ID 上下文传播
- **Objective**：实现 Request ID / Trace ID 的生成、注入与跨服务传播。
- **Context**：plan.md §28/§29 要求每个请求可被 trace 串联；Gateway 负责生成 Request/Trace ID。
- **Dependencies**：`TASK-FND-002`
- **Scope**：middleware（HTTP）与 interceptor（gRPC）中 ID 的生成、context 注入、响应头回传。不含 OpenTelemetry 完整链路（属 Phase 12）。
- **Implementation Requirements**：
  - 无 ID 时生成；有 ID 时透传（trace 串联）。
  - ID 可被 logger 自动读取。
  - gRPC metadata 中传播 trace 上下文（为 Phase 12 铺路）。
- **Acceptance Criteria**：同链路多个服务的日志携带相同 trace_id。
- **Tests**：ID 生成、透传、缺失生成的单元测试。
- **Definition of Done**：ID 上下文传播就位，含测试。

### TASK-FND-007 · gRPC 基础设施
- **Objective**：建立 gRPC 服务端 / 客户端通用底座与 proto 组织。
- **Context**：plan.md §10 内部通信用 gRPC；§40 要求 proto 稳定（不随意复用 field number、保留 reserved）。
- **Dependencies**：`TASK-FND-002`
- **Scope**：proto 目录组织、server/client 构造、interceptor 骨架、health 服务、gRPC 状态码↔业务错误映射。不实现业务 RPC。
- **Implementation Requirements**：
  - proto 定义遵循 plan.md §40 稳定性约束。
  - 提供 gRPC health 端点（返回 SERVING）。
  - 服务间调用通过 gRPC client，不直接暴露给公共客户端。
- **Acceptance Criteria**：任意两个空服务可通过 gRPC 互调并返回 health。
- **Tests**：health 端点 + client/server 联通测试。
- **Definition of Done**：gRPC 底座就位，含测试。

### TASK-FND-008 · Kafka 基础设施（Envelope / Producer / Consumer）
- **Objective**：建立 Kafka 生产 / 消费通用底座与统一 Envelope。
- **Context**：plan.md §15~§21 定义事件驱动规范：Command/Event 区分、Topic 按 Domain 划分、统一 Envelope、Outbox、幂等消费、重试/DLQ。
- **Dependencies**：`TASK-FND-002`
- **Scope**：Envelope 结构、Producer、Consumer 框架（含幂等抽象与重试/DLQ 骨架）、Topic 常量。不含具体业务事件。
- **Implementation Requirements**：
  - Envelope 字段遵循 plan.md §17。
  - Partition key 默认 `aggregate_type + aggregate_id`（plan.md §17）。
  - 提供幂等消费钩子（processed event / idempotency key 抽象，plan.md §20）。
  - 提供重试（初始 3 次）与 DLQ 骨架，DLQ 保留 plan.md §21 要求的字段。
- **Acceptance Criteria**：可生产 / 消费一条 Envelope 消息并观察幂等与 DLQ 机制生效。
- **Tests**：Envelope 序列化、幂等去重、DLQ 投递的单元/集成测试。
- **Definition of Done**：Kafka 底座就位，含测试。

### TASK-FND-009 · MySQL Migration 框架
- **Objective**：建立数据库 migration 流程与工具。
- **Context**：plan.md §41 要求所有结构变化走 Migration，禁止手动改库，含索引/约束。
- **Dependencies**：`TASK-FND-002`
- **Scope**：migration 框架接入、migration 文件组织、可重复执行/版本化。不含业务表结构。
- **Implementation Requirements**：
  - 可追踪、可版本化、与代码版本对应。
  - 支持必要索引与约束声明。
  - 数据库 Schema 属于各 Service 自己的边界（plan.md §41）。
- **Acceptance Criteria**：一个示例 migration 可 up/down 执行。
- **Tests**：migration 幂等性测试。
- **Definition of Done**：migration 框架就位，含示例。

### TASK-FND-010 · Redis Client 封装
- **Objective**：封装统一 Redis Client。
- **Context**：plan.md §27 Redis 用于 Cache/Feed/Counter/Rate Limit/Hot Data/Unread Count/锁，但非 Source of Truth。
- **Dependencies**：`TASK-FND-002`
- **Scope**：client 构造、连接池、基础操作封装（String/Hash/Set/ZSet）、健康检查。不含业务使用。
- **Implementation Requirements**：
  - 统一连接池与超时配置。
  - 支持 cluster/单机可切换的构造。
  - 提供 Ping 健康检查。
- **Acceptance Criteria**：服务可连 Redis 执行基础读写。
- **Tests**：封装操作单元测试（可用 miniredis/testcontainers）。
- **Definition of Done**：Redis 封装就位，含测试。

### TASK-FND-011 · Health Check + Graceful Shutdown
- **Objective**：为所有服务提供健康检查与优雅停机。
- **Context**：plan.md §39 要求 Graceful Shutdown / Timeout / Context Cancellation。
- **Dependencies**：`TASK-FND-002`
- **Scope**：HTTP/gRPC health 端点、信号处理、优雅停机流程（停止接新请求、等待在途完成、关闭依赖）。
- **Implementation Requirements**：
  - 捕获 SIGINT/SIGTERM 触发优雅停机。
  - 停机有超时上限，超时强杀。
  - 正确关闭 gRPC/HTTP/MySQL/Redis/Kafka 连接。
- **Acceptance Criteria**：发 SIGTERM 后服务在途请求完成再退出。
- **Tests**：停机流程集成测试（mock 依赖）。
- **Definition of Done**：健康检查与优雅停机就位，含测试。

### TASK-FND-012 · Docker Compose 编排
- **Objective**：提供一键启动全部基础设施与服务。
- **Context**：Phase 0 验收以 `docker compose up` 为基准。
- **Dependencies**：`TASK-FND-001`
- **Scope**：MySQL8 / Redis / Kafka(KRaft) / Elasticsearch / 各服务 Dockerfile / compose 编排。不含生产 K8s。
- **Implementation Requirements**：
  - 基础设施容器含健康检查。
  - 各服务可 Docker 化部署。
  - 环境变量通过 env 文件注入，密码不硬编码。
- **Acceptance Criteria**：`docker compose up` 后基础设施 healthy，服务可启动。
- **Tests**：compose 启动冒烟（CI 可选）。
- **Definition of Done**：本地一键环境可运行，验收目标达成。

---

## 5. Phase 1 — User and Authentication

> 验收总目标：Register→MySQL→`user.created`→Kafka；Login→JWT；JWT→Gateway→受保护 API。

### TASK-USER-001 · 用户注册 + `user.created` Outbox
- **Objective**：实现用户注册，并用 Outbox 发布 `user.created`。
- **Context**：plan.md §19 Outbox Pattern；§22 新文章默认 DRAFT 同类事务；§42 密码禁止明文。
- **Dependencies**：`TASK-FND-004 / FND-005 / FND-006 / FND-007 / FND-008 / FND-009 / FND-010 / FND-011`
- **Scope**：注册接口、密码哈希、用户表 migration、Outbox 表、Outbox 发布器。不含登录/JWT。
- **Implementation Requirements**：
  - 密码使用安全哈希（如 bcrypt），禁止明文。
  - 注册事务内：插入用户 + 插入 Outbox 事件，再提交（plan.md §19）。
  - Outbox 发布器异步投递到 Kafka，禁止 DB 事务 + 直发 Kafka。
  - 用户名唯一约束；错误码用 `100xx` 空间。
- **Acceptance Criteria**：注册成功后 MySQL 出现用户记录且 Kafka 出现 `user.created`。
- **Tests**：表格驱动单测（参数校验、重复用户名）；集成测试（MySQL + Kafka，验证 Outbox 一致性）。
- **Definition of Done**：注册闭环完成，单测 + 集成测试通过，日志/错误码就位。

### TASK-USER-002 · 用户登录 + JWT 签发
- **Objective**：实现登录与 JWT 签发。
- **Context**：plan.md §11 JWT Payload 含 `sub/username/role/iat/exp`。
- **Dependencies**：`TASK-USER-001`
- **Scope**：登录接口、密码校验、JWT 生成与过期。不含鉴权中间件（属 USER-005）。
- **Implementation Requirements**：
  - JWT 签名密钥从配置注入，禁止硬编码。
  - 错误统一返回，登录失败不泄露用户是否存在（可选加固）。
- **Acceptance Criteria**：合法凭证返回 JWT，非法凭证返回 401/业务错误。
- **Tests**：登录成功 / 密码错误 / 用户不存在 / token 过期 表格驱动测试。
- **Definition of Done**：登录闭环完成，含测试。

### TASK-USER-003 · 用户资料（Profile）
- **Objective**：实现用户资料查看与更新。
- **Context**：核心用户角色（plan.md §3）需要 profile；头像依赖 file-service（Phase 10），本阶段仅存 `file_id`。
- **Dependencies**：`TASK-USER-001`
- **Scope**：资料查询、更新（昵称/简介/头像 file_id）。不含密码修改。
- **Implementation Requirements**：
  - 仅本人可更新（Resource Ownership，plan.md §12）。
  - 头像只存 `file_id`，不存物理路径（plan.md §26）。
- **Acceptance Criteria**：可查询他人公开资料，本人可更新自己资料。
- **Tests**：查询 / 更新 / 越权更新 表格驱动测试。
- **Definition of Done**：资料接口完成，含权限测试。

### TASK-USER-004 · 角色与权限（RBAC）
- **Objective**：建立角色与权限模型。
- **Context**：plan.md §12 RBAC + Resource Ownership，角色 author/moderator/admin，权限与角色分离。
- **Dependencies**：`TASK-USER-001`
- **Scope**：角色/权限数据模型、权限点常量（如 `article:create`、`user:ban`）、权限查询接口。不实现各资源的 ownership 检查（在各业务 Task 内）。
- **Implementation Requirements**：
  - 权限与角色分离，权限点命名遵循 plan.md §12 示例。
  - 提供判定接口供 Gateway/业务服务调用。
- **Acceptance Criteria**：可为用户分配角色并查询其权限集合。
- **Tests**：角色→权限映射、权限点解析 表格驱动测试。
- **Definition of Done**：RBAC 模型完成，含测试。

### TASK-USER-005 · Gateway 认证 / 鉴权中间件
- **Objective**：Gateway 实现 JWT 校验与基础鉴权。
- **Context**：plan.md §11 Gateway 负责基础认证；涉及资源权限的服务必须执行最终授权。
- **Dependencies**：`TASK-USER-002 / USER-004`
- **Scope**：JWT 校验中间件、身份注入下游（user_id/role）、受保护路由。不含具体资源 ownership 授权（业务服务执行）。
- **Implementation Requirements**：
  - 解析 `Authorization: Bearer <token>`。
  - 校验失败返回 401，禁止业务 code 替代 HTTP 状态（plan.md §13）。
  - 将 user_id/role 注入 context 并透传到下游 gRPC。
- **Acceptance Criteria**：无 token / 非法 token 被拒，合法 token 可访问受保护 API。
- **Tests**：认证中间件 表格驱动测试（无/错/过期/合法）。
- **Definition of Done**：认证闭环完成，含测试。

---

## 6. Phase 2 — Article

> 状态机：`DRAFT / PENDING_REVIEW / PUBLISHED / REJECTED / HIDDEN / DELETED`（plan.md §33 Phase 2）。

### TASK-ARTICLE-001 · Article 领域模型 + 状态机 + Markdown
- **Objective**：建立 Article 独立领域模型与状态机。
- **Context**：plan.md §6.3 Article 与 Question 必须独立建模；§22 新 Article 默认 DRAFT；spec §4.2 Markdown 是核心表达方式。
- **Dependencies**：`TASK-USER-005`
- **Scope**：Article 实体、状态枚举与合法流转、Markdown 字段存储。不含 CRUD 接口（属 ARTICLE-002）。
- **Implementation Requirements**：
  - Article 与 Question 独立表、独立领域模型（不共用一张内容表）。
  - 状态机集中定义，非法流转报错。
  - Markdown 原文存正文，渲染/净化后置（Phase 15）。
- **Acceptance Criteria**：状态机所有合法/非法流转被测试覆盖。
- **Tests**：状态机流转表格驱动测试（全状态两两）。
- **Definition of Done**：Article 模型 + 状态机完成，含测试。

### TASK-ARTICLE-002 · Article 草稿 CRUD
- **Objective**：实现 Article 创建/编辑/查询/删除。
- **Context**：plan.md §12 ownership 授权；§22 创建后 DRAFT。
- **Dependencies**：`TASK-ARTICLE-001`
- **Scope**：创建草稿、编辑、按 id 查询、删除（软删）。不含提交审核（ARTICLE-005）。
- **Implementation Requirements**：
  - 仅作者本人可编辑/删除自己文章（`article:update:own` / `article:delete:own`）。
  - 删除为状态 `DELETED`（软删），不物理删除。
- **Acceptance Criteria**：作者可增删改查；非作者操作被拒。
- **Tests**：CRUD + ownership 越权 表格驱动测试。
- **Definition of Done**：Article CRUD 完成，含权限测试。

### TASK-ARTICLE-003 · Article 版本管理
- **Objective**：为 Article 编辑保留版本。
- **Context**：plan.md §6.3 内容版本；§33 Phase 2 列 Version。
- **Dependencies**：`TASK-ARTICLE-001`
- **Scope**：版本表、每次编辑生成版本、版本查询/回滚（可选）。
- **Implementation Requirements**：
  - 编辑时在同一事务记录新版本。
  - 版本可追溯，不破坏主表数据。
- **Acceptance Criteria**：多次编辑后可查询历史版本。
- **Tests**：版本生成与查询 表格驱动测试。
- **Definition of Done**：版本管理完成，含测试。

### TASK-ARTICLE-004 · Tag 实体与内容绑定
- **Objective**：建立 Tag 实体与内容-标签多对多关系（content-service 侧）。
- **Context**：spec §5 Tag 是重要基础实体；plan.md §6.3 Tag 属 content-service；Tag Follow 属 social（Phase 6），Tag 管理属 admin（Phase 11）。
- **Dependencies**：`TASK-ARTICLE-001`
- **Scope**：Tag 实体、Article-Tag 绑定/解绑、Tag 使用计数基础。不含关注（SOCIAL-002）、后台管理（ADMIN-005）。
- **Implementation Requirements**：
  - 内容-标签为多对多关系，走关联表。
  - Tag 参与后续搜索/Feed，命名规范化。
- **Acceptance Criteria**：文章可绑定多个标签，可查询标签下内容。
- **Tests**：绑定/解绑/多对多查询 表格驱动测试。
- **Definition of Done**：Tag 实体与绑定完成，含测试。

### TASK-ARTICLE-005 · Article 提交审核（Outbox）
- **Objective**：实现 Article 提交审核流程。
- **Context**：plan.md §22.2 提交含敏感词检查，状态→PENDING_REVIEW，Outbox→Kafka。
- **Dependencies**：`TASK-ARTICLE-001 / ARTICLE-004`
- **Scope**：提交接口、状态流转（DRAFT→PENDING_REVIEW）、敏感词基础检查、`article.submitted` Outbox。
- **Implementation Requirements**：
  - 提交前权限检查 + 状态检查。
  - 同步敏感词检查（调用 moderation gRPC，若 Phase 4 未就绪可先做本地占位）。
  - Outbox 发布 `article.submitted`。
- **Acceptance Criteria**：提交后状态为 PENDING_REVIEW，Kafka 出现 `article.submitted`。
- **Tests**：状态流转 + Outbox 集成测试。
- **Definition of Done**：提交闭环完成，含测试。

### TASK-ARTICLE-006 · Article 发布 / 驳回状态流转
- **Objective**：消费审核结果事件，驱动 Article 发布。
- **Context**：plan.md §22.3/22.4 审核服务不能改 Content DB；Content 消费 `moderation.approved/rejected` 流转状态并发布 `article.published`。
- **Dependencies**：`TASK-ARTICLE-005, TASK-MOD-004`
- **Scope**：消费 `moderation.approved/rejected/hidden`、状态流转（PENDING_REVIEW→PUBLISHED/REJECTED/HIDDEN）、`article.published` Outbox。幂等消费。
- **Implementation Requirements**：
  - 幂等消费（同 event_id 不重复处理，plan.md §20）。
  - 只有 PENDING_REVIEW 状态才接受审核结果，防乱序。
- **Acceptance Criteria**：审核通过后 Article 状态为 PUBLISHED，Kafka 出现 `article.published`。
- **Tests**：流转 + 幂等（重复事件）+ 非法状态 表格驱动/集成测试。
- **Definition of Done**：发布闭环完成，含幂等测试。

---

## 7. Phase 3 — Q&A

### TASK-QA-001 · Question 领域模型 + 状态机
- **Objective**：建立 Question 独立领域模型与状态机。
- **Context**：plan.md §6.3 Question 与 Article 必须独立建模（不共用表）。
- **Dependencies**：`TASK-USER-005`
- **Scope**：Question 实体、状态（如 OPEN/CLOSED/DELETED）、Markdown 字段。不含接口。
- **Implementation Requirements**：独立表与独立模型；状态流转集中定义。
- **Acceptance Criteria**：Question 状态机被测试覆盖。
- **Tests**：状态机流转 表格驱动测试。
- **Definition of Done**：Question 模型完成，含测试。

### TASK-QA-002 · Question CRUD + 发布（Outbox）
- **Objective**：实现提问/编辑/删除与发布。
- **Context**：plan.md §18 `question.created` / `question.published`。
- **Dependencies**：`TASK-QA-001, TASK-ARTICLE-004`
- **Scope**：提问、编辑、删除（软删）、提交审核与发布事件。ownership 权限。
- **Implementation Requirements**：
  - 仅作者可编辑/删除自己问题。
  - 发布走 Outbox（`question.created/published`）。
  - 提交审核复用 content 审核链路（同 Article）。
- **Acceptance Criteria**：问题可发布，Kafka 出现对应事件。
- **Tests**：CRUD + ownership + Outbox 测试。
- **Definition of Done**：Question 闭环完成，含测试。

### TASK-QA-003 · Answer 领域模型 + CRUD
- **Objective**：实现回答的创建/编辑/删除。
- **Context**：plan.md §18 `answer.created`。
- **Dependencies**：`TASK-QA-002`
- **Scope**：Answer 实体、创建/编辑/删除、Markdown、`answer.created` Outbox。
- **Implementation Requirements**：仅作者可改删；回答归属某 Question。
- **Acceptance Criteria**：问题下可发布回答，Kafka 出现 `answer.created`。
- **Tests**：CRUD + ownership + Outbox 测试。
- **Definition of Done**：Answer 闭环完成，含测试。

### TASK-QA-004 · Accept Answer（采纳回答）
- **Objective**：实现提问者采纳回答。
- **Context**：plan.md §18 `answer.accepted`；§22.6 该事件触发通知。
- **Dependencies**：`TASK-QA-003`
- **Scope**：采纳接口、状态流转、`answer.accepted` Outbox。
- **Implementation Requirements**：仅提问者可采纳；可变更采纳。
- **Acceptance Criteria**：采纳后 Kafka 出现 `answer.accepted`。
- **Tests**：采纳/变更/越权 表格驱动测试。
- **Definition of Done**：采纳闭环完成，含测试。

### TASK-QA-005 · Question-Tag 绑定
- **Objective**：复用 Tag 实体实现 Question 标签绑定。
- **Context**：Tag 实体由 ARTICLE-004 提供。
- **Dependencies**：`TASK-QA-001, TASK-ARTICLE-004`
- **Scope**：Question-Tag 多对多绑定/解绑/查询。
- **Implementation Requirements**：复用 Tag 实体，不重复建模。
- **Acceptance Criteria**：问题可绑定多标签并查询。
- **Tests**：绑定/解绑 表格驱动测试。
- **Definition of Done**：Question-Tag 完成，含测试。

---

## 8. Phase 4 — Moderation

> 关键约束：Moderation Service 不直接修改 Content DB（plan.md §25/§6.6）。

### TASK-MOD-001 · 敏感词库 + 检测
- **Objective**：建立敏感词库与文本检测接口。
- **Context**：plan.md §25 提交时同步基础规则检查；spec §11 敏感词流程（通过/拦截/人工审核）。
- **Dependencies**：`TASK-FND-004 / FND-009`
- **Scope**：敏感词 CRUD、检测接口（返回通过/拦截/待审）、敏感词匹配算法。不含后台管理界面（ADMIN-006）。
- **Implementation Requirements**：
  - 检测接口供 content-service 同步调用（gRPC）。
  - 高风险词可拦截，其余进入人工审核。
- **Acceptance Criteria**：给定文本可返回命中结果与处理建议。
- **Tests**：敏感词匹配 表格驱动测试（命中/未命中/大小写/边界）。
- **Definition of Done**：敏感词检测完成，含测试。

### TASK-MOD-002 · Moderation Task + 消费提交事件
- **Objective**：消费 content 提交事件，生成审核任务。
- **Context**：plan.md §22.3 `article.submitted` → Moderation Task；§20 幂等消费。
- **Dependencies**：`TASK-MOD-001, TASK-ARTICLE-005`
- **Scope**：消费 `article.submitted`（及 question 提交）、Moderation Task 实体、任务状态机（PENDING/APPROVED/REJECTED）。幂等。
- **Implementation Requirements**：
  - 幂等消费（event_id 去重）。
  - Task 记录 aggregate_type/id、提交内容快照、trace_id。
- **Acceptance Criteria**：提交内容后生成对应 Moderation Task。
- **Tests**：消费生成 + 幂等 + 快照 测试。
- **Definition of Done**：任务生成闭环完成，含测试。

### TASK-MOD-003 · 人工审核接口
- **Objective**：提供审核列表/详情/操作接口。
- **Context**：spec §11 举报处理与内容审核流程。
- **Dependencies**：`TASK-MOD-002`
- **Scope**：待审列表、详情、通过/驳回操作入口。RBAC（moderator/admin）。
- **Implementation Requirements**：仅 moderator/admin 可审核（plan.md §12）。
- **Acceptance Criteria**：审核人可查看并操作任务。
- **Tests**：权限 + 列表/详情 表格驱动测试。
- **Definition of Done**：人工审核接口完成，含测试。

### TASK-MOD-004 · 审核结果流转（Approve/Reject/Hide）
- **Objective**：审核结果通过事件驱动回 content-service。
- **Context**：plan.md §22.3 审核不能改 Content DB，发 `moderation.approved/rejected` 事件。
- **Dependencies**：`TASK-MOD-003`
- **Scope**：通过/驳回/隐藏操作、Outbox 发布 `moderation.approved/rejected`、`moderation.*` 事件闭环。
- **Implementation Requirements**：
  - 审核结果走 Outbox → Kafka，不改 Content DB。
  - 事件含 aggregate_type/id，供 content-service 幂等消费。
- **Acceptance Criteria**：审核通过后 Kafka 出现 `moderation.approved`。
- **Tests**：结果流转 + Outbox 集成测试。
- **Definition of Done**：审核结果闭环完成，含测试。

### TASK-MOD-005 · Report 举报与处理
- **Objective**：实现举报提交与处理流程。
- **Context**：spec §11 举报处理流程（举报→待处理→审核→通过/驳回→记录结果）。
- **Dependencies**：`TASK-MOD-001`
- **Scope**：举报创建、待处理列表、处理（通过/驳回）、结果记录。举报对象覆盖 question/answer/article/comment。
- **Implementation Requirements**：举报去重/频率控制基础；处理结果留痕。
- **Acceptance Criteria**：举报可提交并被处理，结果可查。
- **Tests**：举报创建/去重/处理 表格驱动测试。
- **Definition of Done**：举报闭环完成，含测试。

---

## 9. Phase 5 — Interaction

> 重点：幂等、并发、计数一致性（spec §9）。

### TASK-INTERACTION-001 · Comment 领域模型 + CRUD
- **Objective**：实现评论的创建/编辑/删除。
- **Context**：plan.md §18 `comment.created`；评论覆盖 question/answer/article。
- **Dependencies**：`TASK-USER-005, TASK-ARTICLE-001`
- **Scope**：Comment 实体（多态目标）、CRUD、`comment.created` Outbox。
- **Implementation Requirements**：多态目标用 `target_type + target_id` 建模；仅作者可改删。
- **Acceptance Criteria**：任意内容可评论，Kafka 出现 `comment.created`。
- **Tests**：CRUD + ownership + 多态目标 表格驱动测试。
- **Definition of Done**：Comment 闭环完成，含测试。

### TASK-INTERACTION-002 · Like（幂等 + 并发）
- **Objective**：实现点赞，保证幂等与并发安全。
- **Context**：spec §9 点赞需幂等/并发/计数；plan.md §5 Phase 5 重点测试并发 Like、重复 Like。
- **Dependencies**：`TASK-INTERACTION-001`
- **Scope**：Like 创建/取消、唯一约束（user+target 唯一）、`like.created` Outbox。
- **Implementation Requirements**：
  - 幂等：重复点赞不产生重复记录/计数（数据库唯一约束 + 幂等键）。
  - 并发安全：明确竞态风险并说明措施（DB 唯一约束兜底）。
  - 计数通过计数器任务（INTERACTION-005）维护，Like 只负责记录与事件。
- **Acceptance Criteria**：并发点赞 N 次，最终计数正确、记录唯一。
- **Tests**：并发 Like（goroutine 竞态测试，`-race`）、重复 Like、幂等 表格驱动 + 集成测试。
- **Definition of Done**：Like 完成，竞态测试通过。

### TASK-INTERACTION-003 · Collection（幂等）
- **Objective**：实现收藏，保证幂等。
- **Context**：spec §9；plan.md §18 `collection.created`。
- **Dependencies**：`TASK-INTERACTION-001`
- **Scope**：收藏创建/取消、唯一约束、`collection.created` Outbox。
- **Implementation Requirements**：幂等（唯一约束）；并发安全措施说明。
- **Acceptance Criteria**：重复收藏不重复计数。
- **Tests**：重复收藏、取消、幂等 表格驱动测试。
- **Definition of Done**：Collection 完成，含测试。

### TASK-INTERACTION-004 · View 浏览计数
- **Objective**：实现浏览计数（异步聚合）。
- **Context**：plan.md §6.4 View 属 interaction；高频操作需异步。
- **Dependencies**：`TASK-INTERACTION-001`
- **Scope**：浏览上报接口、异步计数（Redis 累加 → 定期落库）。
- **Implementation Requirements**：
  - 浏览计数高并发，用 Redis 累加、异步批量落库，避免逐次写库。
  - 明确最终一致性语义。
- **Acceptance Criteria**：多次浏览最终计数正确。
- **Tests**：异步聚合 + 计数一致性 测试。
- **Definition of Done**：View 完成，含测试。

### TASK-INTERACTION-005 · 互动计数器一致性
- **Objective**：统一维护 Like/Collect/Comment 计数，保证一致性。
- **Context**：plan.md §6.4 Interaction Counter；§27 Redis 计数器但非 Source of Truth。
- **Dependencies**：`TASK-INTERACTION-002 / 003 / 004`
- **Scope**：计数器表（MySQL 为准）、Redis 缓存、异步对账。
- **Implementation Requirements**：
  - MySQL 为计数 Source of Truth，Redis 为缓存/热数据。
  - 提供缓存回源与对账机制，保证最终一致。
- **Acceptance Criteria**：计数与真实记录一致（含缓存失效后）。
- **Tests**：计数一致性（缓存失效/对账）集成测试。
- **Definition of Done**：计数器一致性完成，含测试。

---

## 10. Phase 6 — Social

### TASK-SOCIAL-001 · User Follow / Unfollow
- **Objective**：实现用户关注/取关。
- **Context**：spec §6；plan.md §18 `user.followed`；§6.5 Social Graph。
- **Dependencies**：`TASK-USER-005`
- **Scope**：关注/取关接口、唯一约束（follower+followee）、`user.followed` Outbox、幂等。
- **Implementation Requirements**：
  - 幂等（唯一约束）；不能关注自己（校验）。
  - 关注计数维护。
- **Acceptance Criteria**：关注后 Kafka 出现 `user.followed`，重复关注幂等。
- **Tests**：关注/取关/自关注/幂等 表格驱动测试。
- **Definition of Done**：User Follow 完成，含测试。

### TASK-SOCIAL-002 · Tag Follow
- **Objective**：实现标签关注。
- **Context**：spec §5/§6；plan.md §18 `tag.followed`。
- **Dependencies**：`TASK-SOCIAL-001, TASK-ARTICLE-004`
- **Scope**：Tag 关注/取关、`tag.followed` Outbox。
- **Implementation Requirements**：幂等；tag_id 来自 content-service 的 Tag。
- **Acceptance Criteria**：关注标签后 Kafka 出现 `tag.followed`。
- **Tests**：关注/取关/幂等 表格驱动测试。
- **Definition of Done**：Tag Follow 完成，含测试。

### TASK-SOCIAL-003 · Social Graph 查询
- **Objective**：提供关注/粉丝列表与计数查询。
- **Context**：plan.md §6.5 Social Graph。
- **Dependencies**：`TASK-SOCIAL-001 / 002`
- **Scope**：following 列表、followers 列表、关注计数。
- **Implementation Requirements**：分页查询；计数走缓存/计数器。
- **Acceptance Criteria**：可查询某用户关注列表与粉丝列表。
- **Tests**：列表/分页/计数 表格驱动测试。
- **Definition of Done**：Social Graph 查询完成，含测试。

---

## 11. Phase 7 — Notification

### TASK-NOTIF-001 · Notification 模型 + 事件消费
- **Objective**：消费互动/社交/审核事件，生成站内通知。
- **Context**：plan.md §22.6 消费 `like.created/comment.created/user.followed/answer.accepted/moderation.rejected`。
- **Dependencies**：`TASK-FND-008, TASK-USER-005`
- **Scope**：Notification 实体、事件消费者（幂等）、通知生成规则。
- **Implementation Requirements**：
  - 幂等消费（event_id 去重）。
  - 按事件类型生成对应通知文案与关联对象。
  - 忽略「自己给自己的操作」产生的通知。
- **Acceptance Criteria**：消费上述事件后生成对应通知。
- **Tests**：各事件类型生成 + 幂等 表格驱动/集成测试。
- **Definition of Done**：通知生成闭环完成，含测试。

### TASK-NOTIF-002 · 通知列表 + 未读数
- **Objective**：提供通知列表与未读数查询。
- **Context**：plan.md §6.7 Unread Count；§27 Redis 存未读数。
- **Dependencies**：`TASK-NOTIF-001`
- **Scope**：通知列表（分页）、未读数查询（Redis 缓存）。
- **Implementation Requirements**：未读数用 Redis 缓存，MySQL 为准。
- **Acceptance Criteria**：可查通知列表与未读数。
- **Tests**：列表/分页/未读数 表格驱动测试。
- **Definition of Done**：通知查询完成，含测试。

### TASK-NOTIF-003 · 已读 / 全部已读
- **Objective**：实现单条/全部标记已读。
- **Context**：plan.md §6.7 Mark as Read / Mark All Read。
- **Dependencies**：`TASK-NOTIF-002`
- **Scope**：单条已读、全部已读接口，未读数清零。
- **Implementation Requirements**：仅本人可标记；已读状态持久化。
- **Acceptance Criteria**：标记后未读数正确更新。
- **Tests**：单条/全部已读 表格驱动测试。
- **Definition of Done**：已读闭环完成，含测试。

---

## 12. Phase 8 — Search

> Search 数据属 Derived Data，MySQL 是 Source of Truth（plan.md §6.9）。

### TASK-SEARCH-001 · ES 索引与 Mapping
- **Objective**：定义 ES 索引与 Mapping。
- **Context**：plan.md §24 支持 Article/Question/Answer/Tag/Author 过滤、高亮、分页。
- **Dependencies**：`TASK-FND-009, TASK-ARTICLE-001, TASK-QA-001`
- **Scope**：article/question/answer/user/tag 索引 mapping、analyzer（分词）、重建脚本骨架。
- **Implementation Requirements**：分词、相关性字段设计；mapping 版本化。
- **Acceptance Criteria**：索引可创建，mapping 符合检索需求。
- **Tests**：mapping 校验测试。
- **Definition of Done**：索引定义完成，含测试。

### TASK-SEARCH-002 · 索引同步（消费事件）
- **Objective**：消费 content 事件同步索引。
- **Context**：plan.md §22.5 `article.published` → ES；§6.9 ES 非 Source of Truth。
- **Dependencies**：`TASK-SEARCH-001, TASK-ARTICLE-006`
- **Scope**：消费 `article.published`（及 question/answer）写入 ES，删除/隐藏同步。幂等。
- **Implementation Requirements**：幂等；ES 写入失败可重试，不影响主流程。
- **Acceptance Criteria**：发布文章后可被搜索到。
- **Tests**：同步 + 幂等 集成测试。
- **Definition of Done**：索引同步完成，含测试。

### TASK-SEARCH-003 · 搜索 API（过滤/高亮/排序/分页）
- **Objective**：实现搜索接口。
- **Context**：spec §8 关键词搜索/全文检索/相关性排序/高亮；plan.md §24 过滤分页。
- **Dependencies**：`TASK-SEARCH-002`
- **Scope**：关键词搜索、按类型/Tag/作者过滤、高亮、相关性排序、分页。
- **Implementation Requirements**：搜索排序与高亮实现；请求参数校验。
- **Acceptance Criteria**：关键词可召回相关文章/问题/回答并高亮。
- **Tests**：过滤/排序/分页/高亮 集成测试。
- **Definition of Done**：搜索 API 完成，含测试。

### TASK-SEARCH-004 · Index Rebuild（全量重建）
- **Objective**：实现 MySQL→ES 全量重建。
- **Context**：plan.md §24 必须支持 Index Rebuild。
- **Dependencies**：`TASK-SEARCH-001`
- **Scope**：全量重建流程（从 MySQL 批量读取写入 ES）、重建期间一致性处理。
- **Implementation Requirements**：可增量或全量触发；重建不影响在线查询（可切换索引）。
- **Acceptance Criteria**：重建后索引与 MySQL 数据一致。
- **Tests**：重建一致性测试。
- **Definition of Done**：Index Rebuild 完成，含测试。

---

## 13. Phase 9 — Feed

> Feed 用 Redis Sorted Set；普通作者 Fan-out on Write，高关注作者 Fan-out on Read（阈值压测确定）；Ranking 走抽象接口（plan.md §23）。

### TASK-FEED-001 · Feed 数据模型 + Redis Sorted Set
- **Objective**：建立 Feed 收件箱（inbox）数据模型。
- **Context**：plan.md §23 Feed 用 Redis Sorted Set。
- **Dependencies**：`TASK-FND-010, TASK-ARTICLE-001`
- **Scope**：用户 inbox 的 Sorted Set 结构、score 设计、读写封装。
- **Implementation Requirements**：定义 score 字段（时间 + 热度）；封装 ZAdd/ZRevRange。
- **Acceptance Criteria**：可向 inbox 写入并按 score 读取。
- **Tests**：Sorted Set 读写 单元测试（miniredis/testcontainers）。
- **Definition of Done**：Feed 数据模型完成，含测试。

### TASK-FEED-002 · Fan-out on Write（普通作者）
- **Objective**：消费发布事件，写扩散到关注者 inbox。
- **Context**：plan.md §22.7/§23 普通作者写扩散。
- **Dependencies**：`TASK-FEED-001, TASK-ARTICLE-006`
- **Scope**：消费 `article.published`（及 question/answer），查询作者粉丝，写入粉丝 inbox。幂等。
- **Implementation Requirements**：
  - 幂等消费。
  - 粉丝数低于阈值走写扩散；写入失败可重试。
- **Acceptance Criteria**：关注者可在 Feed 看到所关注作者的新内容。
- **Tests**：写扩散 + 幂等 集成测试。
- **Definition of Done**：Fan-out on Write 完成，含测试。

### TASK-FEED-003 · Fan-out on Read（高关注作者）
- **Objective**：高关注作者内容读时合并。
- **Context**：plan.md §22.7/§23 读扩散。
- **Dependencies**：`TASK-FEED-001, TASK-ARTICLE-006`
- **Scope**：识别高关注作者、读时拉取其最近内容合并进 Feed。
- **Implementation Requirements**：阈值可配；读扩散与写扩散结果合并排序。
- **Acceptance Criteria**：高关注作者内容可通过读扩散被消费。
- **Tests**：读扩散合并 集成测试。
- **Definition of Done**：Fan-out on Read 完成，含测试。

### TASK-FEED-004 · Rule Ranking（抽象接口）
- **Objective**：实现规则排序引擎（抽象接口）。
- **Context**：plan.md §23 score = freshness + engagement + author + tag_match；Ranking Engine 抽象接口可替换。
- **Dependencies**：`TASK-FEED-001`
- **Scope**：Ranking 接口定义、规则实现、可插拔。
- **Implementation Requirements**：
  - 定义 `Ranker` 接口，规则实现可替换为未来推荐系统。
  - 实现 plan.md §23 的 score 公式。
- **Acceptance Criteria**：不同内容按规则得分排序。
- **Tests**：ranking 公式 表格驱动测试（新鲜度/热度/作者/标签）。
- **Definition of Done**：Ranking 完成，含测试。

### TASK-FEED-005 · Feed API（Following/Tag/Latest/Hot）
- **Objective**：提供四类信息流接口。
- **Context**：spec §7 推荐/关注/最新/热门；plan.md §6.8。
- **Dependencies**：`TASK-FEED-002 / 003 / 004, TASK-SOCIAL-001 / 002`
- **Scope**：Following Feed、Tag Feed、Latest、Hot 四个接口。
- **Implementation Requirements**：结合关注用户 + 关注标签 + 热门 + 新鲜度；分页。
- **Acceptance Criteria**：可分别获取四类 Feed。
- **Tests**：四类 Feed 集成测试。
- **Definition of Done**：Feed API 完成，含测试。

---

## 14. Phase 10 — File

> Storage 抽象接口，第一版 LocalStorage，未来 OSS/S3/MinIO（plan.md §26）。

### TASK-FILE-001 · Storage 抽象 + LocalStorage
- **Objective**：定义 Storage 接口并实现本地存储。
- **Context**：plan.md §26 `Storage{Put/Get/Delete}`；spec §12 业务库只存元数据。
- **Dependencies**：`TASK-FND-002`
- **Scope**：Storage 接口、LocalStorage 实现、file 元数据表（file_id/name/path/type/size）。
- **Implementation Requirements**：接口抽象，实现可替换；元数据与物理文件解耦。
- **Acceptance Criteria**：可通过接口 Put/Get/Delete 文件。
- **Tests**：Storage 接口实现 表格驱动测试。
- **Definition of Done**：Storage 抽象完成，含测试。

### TASK-FILE-002 · 文件上传（校验 + 元数据）
- **Objective**：实现文件上传接口。
- **Context**：plan.md §26 `POST /files`；§42 文件类型/大小限制、路径穿越防护。
- **Dependencies**：`TASK-FILE-001`
- **Scope**：上传接口、类型/大小校验、元数据入库、返回 file_id。
- **Implementation Requirements**：
  - 校验文件类型与大小上限。
  - 防路径穿越（文件名清洗、不信任客户端路径）。
  - 业务只保存 file_id。
- **Acceptance Criteria**：上传成功返回 file_id 并可通过 id 访问。
- **Tests**：类型/大小/路径穿越 表格驱动测试。
- **Definition of Done**：上传完成，含安全测试。

### TASK-FILE-003 · 文件下载（流式 + 权限）
- **Objective**：实现文件下载接口。
- **Context**：plan.md §26 `GET /files/{id}`。
- **Dependencies**：`TASK-FILE-002`
- **Scope**：下载接口、流式返回、基础权限控制。
- **Implementation Requirements**：流式读取避免内存膨胀；鉴权（公开/私有文件区分）。
- **Acceptance Criteria**：可通过 file_id 下载文件。
- **Tests**：下载/权限/不存在 表格驱动测试。
- **Definition of Done**：下载完成，含测试。

---

## 15. Phase 11 — Admin

> 关键约束：Admin Service 不直接访问其他 Service 数据库，跨服务通过 gRPC 编排（plan.md §6.11）。

### TASK-ADMIN-001 · 用户管理
- **Objective**：提供后台用户管理。
- **Context**：spec §11 用户管理；通过 user-service gRPC。
- **Dependencies**：`TASK-USER-001, TASK-FND-007`
- **Scope**：用户查询、状态（封禁/解禁）、活跃情况。RBAC（admin）。
- **Implementation Requirements**：仅 admin；调用 user-service gRPC，不直连 user_db。
- **Acceptance Criteria**：可查询并变更用户状态。
- **Tests**：查询/状态变更/权限 表格驱动测试。
- **Definition of Done**：用户管理完成，含测试。

### TASK-ADMIN-002 · 内容管理
- **Objective**：提供后台内容管理。
- **Context**：spec §11 内容查看/审核/隐藏/删除/恢复；通过 content-service gRPC。
- **Dependencies**：`TASK-ARTICLE-006, TASK-QA-002, TASK-FND-007`
- **Scope**：内容查看、隐藏、删除、恢复。RBAC（moderator/admin）。
- **Implementation Requirements**：调用 content-service gRPC，不直连 content_db。
- **Acceptance Criteria**：可对内容执行隐藏/删除/恢复。
- **Tests**：各操作 + 权限 表格驱动测试。
- **Definition of Done**：内容管理完成，含测试。

### TASK-ADMIN-003 · 评论管理
- **Objective**：提供后台评论管理。
- **Context**：spec §11 评论查看/审核/删除；通过 interaction-service gRPC。
- **Dependencies**：`TASK-INTERACTION-001, TASK-FND-007`
- **Scope**：评论查看、审核、删除。RBAC。
- **Implementation Requirements**：调用 interaction-service gRPC。
- **Acceptance Criteria**：可审核/删除评论。
- **Tests**：操作 + 权限 表格驱动测试。
- **Definition of Done**：评论管理完成，含测试。

### TASK-ADMIN-004 · 举报处理
- **Objective**：提供后台举报处理。
- **Context**：spec §11 举报处理流程；通过 moderation-service gRPC。
- **Dependencies**：`TASK-MOD-005, TASK-FND-007`
- **Scope**：举报列表、处理（通过/驳回）、结果记录。
- **Implementation Requirements**：调用 moderation-service gRPC。
- **Acceptance Criteria**：可处理举报并记录结果。
- **Tests**：处理流程 + 权限 表格驱动测试。
- **Definition of Done**：举报处理完成，含测试。

### TASK-ADMIN-005 · 标签管理
- **Objective**：提供后台标签管理。
- **Context**：spec §11 创建/编辑/合并/禁用标签；通过 content-service gRPC。
- **Dependencies**：`TASK-ARTICLE-004, TASK-FND-007`
- **Scope**：标签创建、编辑、合并、禁用、使用情况。RBAC。
- **Implementation Requirements**：调用 content-service gRPC。
- **Acceptance Criteria**：可对标签执行创建/编辑/合并/禁用。
- **Tests**：各操作 + 权限 表格驱动测试。
- **Definition of Done**：标签管理完成，含测试。

### TASK-ADMIN-006 · 敏感词管理
- **Objective**：提供后台敏感词管理。
- **Context**：spec §11 敏感词管理；通过 moderation-service gRPC。
- **Dependencies**：`TASK-MOD-001, TASK-FND-007`
- **Scope**：敏感词 CRUD。RBAC。
- **Implementation Requirements**：调用 moderation-service gRPC。
- **Acceptance Criteria**：可管理敏感词库。
- **Tests**：CRUD + 权限 表格驱动测试。
- **Definition of Done**：敏感词管理完成，含测试。

### TASK-ADMIN-007 · 数据统计
- **Objective**：提供社区运营数据统计。
- **Context**：spec §11 数据统计（用户/内容/互动/增长趋势）。
- **Dependencies**：`TASK-FND-007`
- **Scope**：用户数、活跃、新增、各内容量、举报量、趋势。
- **Implementation Requirements**：通过各服务 gRPC 聚合或 service-owned read model；不直连他库。
- **Acceptance Criteria**：可查询各项统计指标。
- **Tests**：统计聚合 测试。
- **Definition of Done**：数据统计完成，含测试。

### TASK-ADMIN-008 · 审计日志
- **Objective**：实现管理员操作审计日志。
- **Context**：plan.md §6.11 Audit Log；§42 Audit Logging。
- **Dependencies**：`TASK-FND-005 / FND-007`
- **Scope**：记录 admin 关键操作（谁、何时、做了什么）。
- **Implementation Requirements**：审计日志结构化，含操作者/时间/动作/对象。
- **Acceptance Criteria**：admin 操作可被审计追溯。
- **Tests**：审计记录 表格驱动测试。
- **Definition of Done**：审计日志完成，含测试。

---

## 16. Phase 12 — Observability

### TASK-OBS-001 · OpenTelemetry 接入
- **Objective**：接入 OpenTelemetry 分布式追踪。
- **Context**：plan.md §29 Trace 覆盖 Gateway→Service→gRPC→MySQL/Redis→Kafka→Consumer；Kafka 传播 Trace Context。
- **Dependencies**：`TASK-FND-006 / FND-007 / FND-008`
- **Scope**：OTel SDK 接入、gRPC/HTTP/DB/Redis/Kafka 埋点、Trace Context 跨 Kafka 传播。
- **Implementation Requirements**：Kafka 事件 Envelope 携带 trace 上下文；消费者恢复 trace。
- **Acceptance Criteria**：一次请求全链路在同一 trace 中可见。
- **Tests**：trace 传播集成测试。
- **Definition of Done**：OTel 接入完成，含测试。

### TASK-OBS-002 · Metrics（HTTP/gRPC/Kafka/Business）
- **Objective**：实现指标采集。
- **Context**：plan.md §30 定义 HTTP/gRPC/Kafka/Business 指标名。
- **Dependencies**：`TASK-OBS-001`
- **Scope**：Prometheus 指标暴露、HTTP/gRPC/Kafka 指标、业务指标（注册/发布/点赞等）。
- **Implementation Requirements**：指标名遵循 plan.md §30；暴露 `/metrics`。
- **Acceptance Criteria**：指标可被 Prometheus 抓取。
- **Tests**：指标存在性测试。
- **Definition of Done**：Metrics 完成，含测试。

### TASK-OBS-003 · Grafana / Jaeger 部署与 Dashboard
- **Objective**：部署可视化与看板。
- **Context**：plan.md §28/§33 Phase 12 Dashboards。
- **Dependencies**：`TASK-OBS-002`
- **Scope**：Grafana/Jaeger 容器、数据源配置、基础 Dashboard（HTTP/QPS/延迟/Kafka lag）。
- **Implementation Requirements**：纳入 docker compose；Dashboard 可导入。
- **Acceptance Criteria**：Grafana 可展示服务指标，Jaeger 可查询 trace。
- **Tests**：无（部署验证）。
- **Definition of Done**：可视化部署完成。

### TASK-OBS-004 · 完整 Trace Demo
- **Objective**：建立端到端 Trace Demo。
- **Context**：plan.md §33 Phase 12 演示 HTTP→Gateway→Content→MySQL→Outbox→Kafka→Moderation→Content→Search→ES 同 trace。
- **Dependencies**：`TASK-OBS-003, TASK-ARTICLE-006, TASK-SEARCH-002`
- **Scope**：构造一次完整业务流，验证同 trace_id 串联全链路。
- **Implementation Requirements**：演示脚本/文档；确认 Kafka 环节 trace 不断裂。
- **Acceptance Criteria**：端到端全链路 trace 可完整串联。
- **Tests**：端到端验证脚本。
- **Definition of Done**：Trace Demo 完成并可复现。

---

## 17. Phase 13 — Reliability

### TASK-REL-001 · 依赖不可用容错
- **Objective**：验证并实现依赖不可用时的容错。
- **Context**：plan.md §33 Phase 13 场景：Kafka/Redis/ES 不可用、MySQL 超时、gRPC 超时、网络故障。
- **Dependencies**：`TASK-NOTIF-003, TASK-SEARCH-003, TASK-FEED-005`
- **Scope**：各服务对依赖故障的降级、超时、错误处理。
- **Implementation Requirements**：区分 Retryable/Non-Retryable（plan.md §43）；依赖不可用不 panic、不无限重试。
- **Acceptance Criteria**：注入依赖故障后服务可降级/报错且不崩溃。
- **Tests**：故障注入测试。
- **Definition of Done**：依赖容错完成，含测试。

### TASK-REL-002 · 消费者重试 / 恢复
- **Objective**：验证消费者重试、崩溃恢复、rebalance。
- **Context**：plan.md §21 重试 3 次→DLQ；§20 幂等；§33 Consumer Crash/Restart。
- **Dependencies**：`TASK-REL-001`
- **Scope**：消费者失败重试、崩溃重启恢复、rebalance 后不重/漏处理。
- **Implementation Requirements**：重试有上限 + Backoff；DLQ 保留关键字段。
- **Acceptance Criteria**：重试后成功提交；超限进 DLQ；重启不重复处理。
- **Tests**：重试/DLQ/恢复 集成测试。
- **Definition of Done**：消费者可靠性完成，含测试。

### TASK-REL-003 · 重复 / 乱序事件验证
- **Objective**：验证重复与乱序事件下的幂等与一致性。
- **Context**：plan.md §20 幂等；§33 Duplicate/Out-of-order Event。
- **Dependencies**：`TASK-REL-001`
- **Scope**：重复事件、乱序事件的幂等与最终一致验证。
- **Implementation Requirements**：依赖 event_id 幂等 + 状态机守卫（如只有 PENDING_REVIEW 接受审核结果）。
- **Acceptance Criteria**：重复/乱序事件不产生错误状态。
- **Tests**：重复/乱序事件测试。
- **Definition of Done**：幂等/一致性验证完成，含测试。

### TASK-REL-004 · Failure 测试套件
- **Objective**：建立故障测试套件。
- **Context**：plan.md §31.2/§33 Failure Test。
- **Dependencies**：`TASK-REL-002 / 003`
- **Scope**：自动化故障测试（依赖宕机、超时、重复事件等场景集）。
- **Implementation Requirements**：可重复执行；覆盖 plan.md §33 场景清单。
- **Acceptance Criteria**：故障测试套件可一键运行并报告通过。
- **Tests**：全套故障测试通过。
- **Definition of Done**：Failure 测试套件完成。

---

## 18. Phase 14 — Performance

### TASK-PERF-001 · HTTP / gRPC 压测
- **Objective**：HTTP 与 gRPC 负载测试。
- **Context**：plan.md §33 Phase 14 HTTP/gRPC Load Test。
- **Dependencies**：`TASK-SEARCH-003, TASK-FEED-005`
- **Scope**：HTTP/gRPC 压测脚本，记录 RPS/P50/P95/P99/Error Rate。
- **Implementation Requirements**：压测工具（k6/ghz/vegeta）脚本化。
- **Acceptance Criteria**：产出压测报告与关键指标。
- **Tests**：压测脚本可运行。
- **Definition of Done**：HTTP/gRPC 压测完成并有报告。

### TASK-PERF-002 · Kafka 吞吐压测
- **Objective**：Kafka 吞吐与 lag 测试。
- **Context**：plan.md §33 Phase 14 Kafka Throughput Test。
- **Dependencies**：`TASK-PERF-001`
- **Scope**：Kafka 吞吐压测，记录吞吐与 consumer lag。
- **Implementation Requirements**：压测脚本；观察 lag 指标。
- **Acceptance Criteria**：产出 Kafka 吞吐与 lag 数据。
- **Tests**：压测脚本可运行。
- **Definition of Done**：Kafka 压测完成并有报告。

### TASK-PERF-003 · Feed / Search 压测
- **Objective**：Feed 与 Search 负载测试。
- **Context**：plan.md §33 Phase 14 Feed/Search Load Test。
- **Dependencies**：`TASK-PERF-001`
- **Scope**：Feed 读路径、Search 查询压测，记录延迟与 Redis/ES 表现。
- **Implementation Requirements**：压测脚本；记录 Redis Hit Rate、ES 延迟。
- **Acceptance Criteria**：产出 Feed/Search 性能数据。
- **Tests**：压测脚本可运行。
- **Definition of Done**：Feed/Search 压测完成并有报告。

### TASK-PERF-004 · 压测报告与指标记录
- **Objective**：汇总压测报告。
- **Context**：plan.md §33 Phase 14 记录 RPS/P50/P95/P99/Error Rate/Kafka Lag/DB Connections/Redis Hit Rate。
- **Dependencies**：`TASK-PERF-002 / 003`
- **Scope**：汇总所有压测数据成报告。
- **Implementation Requirements**：记录全部要求指标。
- **Acceptance Criteria**：完整压测报告。
- **Tests**：无（报告）。
- **Definition of Done**：压测报告完成。

---

## 19. Phase 15 — Production Hardening

### TASK-HARDEN-001 · 限流（Rate Limit）
- **Objective**：Gateway 实现限流。
- **Context**：plan.md §42 Rate Limiting；§6.1 Gateway 负责。
- **Dependencies**：`TASK-FND-012, TASK-USER-005`
- **Scope**：基于 Redis 的限流（IP/用户/接口维度）。
- **Implementation Requirements**：限流策略可配；超限返回 429。
- **Acceptance Criteria**：超限请求被 429 拒绝。
- **Tests**：限流阈值 表格驱动测试。
- **Definition of Done**：限流完成，含测试。

### TASK-HARDEN-002 · 安全加固（注入/XSS/文件/JWT/密码）
- **Objective**：落实安全基线。
- **Context**：plan.md §42 全部项：SQL 注入、XSS、Markdown Sanitization、文件类型/大小、JWT、密码。
- **Dependencies**：`TASK-REL-004, TASK-PERF-004`
- **Scope**：SQL 注入防护（参数化）、XSS/Markdown 净化、文件安全、JWT/密码安全复查。
- **Implementation Requirements**：
  - Markdown 渲染前净化（防 XSS）。
  - 文件类型/大小/路径穿越校验。
  - JWT 校验强度、密码哈希强度复查。
- **Acceptance Criteria**：安全基线项全部落实。
- **Tests**：注入/净化/文件安全 安全测试。
- **Definition of Done**：安全加固完成，含测试。

### TASK-HARDEN-003 · 弹性（Timeout/Retry/Circuit Breaker）
- **Objective**：实现超时/重试/熔断。
- **Context**：plan.md §33 Phase 15 Timeout/Retry/Circuit Breaker。
- **Dependencies**：`TASK-HARDEN-002`
- **Scope**：gRPC/HTTP 超时、重试上限、熔断器。
- **Implementation Requirements**：重试有上限 + Backoff；熔断阈值可配。
- **Acceptance Criteria**：下游故障时服务不雪崩。
- **Tests**：熔断/超时 集成测试。
- **Definition of Done**：弹性机制完成，含测试。

### TASK-HARDEN-004 · 敏感数据保护与审计
- **Objective**：敏感数据脱敏与审计完善。
- **Context**：plan.md §42 Sensitive Data Protection / Audit Logging。
- **Dependencies**：`TASK-HARDEN-002`
- **Scope**：日志脱敏（密码/token）、敏感字段保护、审计覆盖。
- **Implementation Requirements**：日志不落敏感信息；审计留痕完整。
- **Acceptance Criteria**：日志无明文敏感数据。
- **Tests**：脱敏 表格驱动测试。
- **Definition of Done**：敏感数据保护完成，含测试。

---

## 20. 执行建议

1. **从 Phase 0 开始**：`TASK-FND-001` → `FND-002` → 其余 FND 任务可部分并行。
2. **并行窗口**：Phase 0 完成后，`ARTICLE-001` 与 `QA-001` 可由不同 Agent 并行（plan.md §34）。
3. **验收口径**：每个 Task 以 §2.2 的 `Definition of Done` 为准，不以外观「能编译」为完成（plan.md §38）。
4. **冲突处理**：若实现中发现与 plan.md 冲突，停止扩大范围、报告冲突、提出候选方案、等待决策（plan.md §35）。
