# Community Platform - Development Plan

Version: 1.0
Status: Draft

---

## 1. Document Purpose

本文档定义 Community Platform 的技术实现规划。

`spec.md` 定义系统的产品意图、功能边界、业务规则与架构要求。

`plan.md` 定义系统如何被实现，包括：

- 技术栈
- 系统架构
- 微服务边界
- 数据架构
- 服务间通信
- 事件驱动架构
- 基础设施
- 测试策略
- 可观测性
- 开发阶段
- 交付策略

`tasks.md` 将基于本文档进一步拆分为可由 AI Agent 独立执行和验收的工程任务。

---

## 2. Engineering Goals

本项目不仅需要实现完整的社区业务，还需要体现现代 Go 后端与分布式系统工程能力。

主要工程目标：

1. 微服务架构
2. 服务自治
3. gRPC 内部通信
4. Kafka 事件驱动
5. Outbox Pattern
6. Idempotent Consumer
7. Redis
8. Elasticsearch
9. Feed System
10. Content Moderation
11. Distributed Tracing
12. Metrics
13. Structured Logging
14. Automated Testing
15. Integration Testing
16. E2E Testing
17. Load Testing
18. Failure Testing
19. Docker 化部署
20. AI Agent 协同开发

项目最终应能够作为高级 Go 后端 / 分布式系统方向的综合工程项目进行展示。

---

## 3. Product Scope

系统是一个面向程序员和技术学习者的知识社区。

核心内容类型：

- Article
- Question
- Answer

核心社区能力：

- Markdown 内容创作
- Tag
- Comment
- Like
- Collection
- View
- User Follow
- Tag Follow
- Notification
- Feed
- Search
- Moderation
- Report
- File Upload
- Admin
- Statistics

---

## 4. Architecture Overview

系统采用前后端分离的微服务架构。

整体结构：

```text
                    Client
                      |
                      | HTTP / JSON
                      v
                API Gateway
                      |
                      | gRPC
                      v
       +--------------+---------------+
       |              |               |
       v              v               v
     User          Content        Interaction
       |              |               |
       +--------------+---------------+
                      |
                   Kafka
                      |
        +-------------+-------------+
        |             |             |
        v             v             v
      Search         Feed     Notification
        |
        v
 Elasticsearch

Moderation
     ^
     |
   Kafka

File
     |
 Local Storage
``` 

---

## 5. Service Architecture

系统包含：
1. gateway-service
2. user-service
3. content-service
4. interaction-service
5. social-service
6. moderation-service
7. notification-service
8. feed-service
9. search-service
10. file-service
11. admin-service

其中：
- gateway-service 为 API Gateway
- 其余为业务服务

---

## 6. Service Responsibilities

**6.1 Gateway Service**
负责：
- HTTP Routing
- Authentication
- Request ID
- Trace ID
- CORS
- Rate Limiting
- REST API Response Formatting
- gRPC Client Routing
Gateway 不允许包含领域业务逻辑。

**6.2 User Service**
- 负责：
- 用户注册
- 用户登录
- 用户资料
- 用户状态
- 用户角色
- 密码管理
- 用户生命周期
拥有：`user_db`

**6.3 Content Service**
负责：
- Article
- Question
- Answer
- Tag
- Markdown 内容
- 内容状态
- 内容发布
- 内容编辑
- 内容版本
- 内容可见性
拥有：`content_db`，Article与Question必须保持独立领域模型。

**6.4 Interaction Service**
负责：
- Like
- Collection
- Comment
- View
- Interaction Counter
拥有：`interaction_db`


**6.5 Social Service**
负责：
- User Follow
- User Unfollow
- Tag Follow
- Social Graph
拥有：`social_db`

**6.6 Moderation Service**
负责：
- Content Moderation
- Report
- Sensitive Word
- Moderation Task
- Human Review
- Moderation Workflow
拥有：`moderation_db`

**6.7 Notification Service**
负责：
- Notification
- Notification History
- Unread Count
- Mark as Read
拥有：`notification_db`

**6.8 Feed Service**
负责：
- Personalized Feed
- Following Feed
- Tag Feed
- Feed Ranking
- Fan-out on Write
- Fan-out on Read
拥有：`feed_db`，Redis用于Feed热数据。

**6.9 Search Service**
负责：
- Elasticsearch Index
- Search
- Search Filter
- Search Ranking
- Highlight
- Index Rebuild
Search数据属于Derived Data，MySQL是Source of Truth。

**6.10 File Service**
负责：
- File Upload
- File Download
- File Metadata
- File Validation
- Storage Abstraction
初期：`Local Storage`，未来可扩展：oss s3  MinIO

**6.11 Admin Service**
负责：
- User Management
- Content Management
- Comment Management
- Report Management
- Tag Management
- Sensitive Word Management
- Statistics
- Audit Log
Admin Service 不允许直接访问其他 Service 的数据库。


---

## 7. Data Ownership
每个业务服务拥有独立的数据库。
```
user-service
    -> user_db

content-service
    -> content_db

interaction-service
    -> interaction_db

social-service
    -> social_db

moderation-service
    -> moderation_db

notification-service
    -> notification_db

feed-service
    -> feed_db

file-service
    -> file_db

admin-service
    -> admin_db
```
服务之间禁止直接访问其他服务数据库。跨服务数据访问必须使用：gRPC，Kafka Event，Service-owned Read Model

---

## 8. Technology Stack 
**Backend**
- GO
- Gin
- gRPC
- Protocol Buffers
- GORM

**Database**
- MySQL8

**Cache**
- Redis

**Message Queue**
- Kafka
- KRaft

**Search**
- Elasticsearch

**Observability**
- OpenTelemetry
- Prometheus
- Grafana
- Jaeger
- Zap

**Deployment**
- Docker
- Docker Compose

**Frontend**
- Vue3
- Vite
- JavaScript

---

## 9. Database Access
- GORM 是默认 ORM。

- 简单 CRUD 使用 GORM。

- 复杂查询或性能敏感查询允许使用显式 SQL。

- 业务层禁止直接依赖 GORM。

推荐结构：
```
Handler
   |
   v
Service
   |
   v
Repository Interface
   |
   v
GORM Repository
   |
   v
MySQL
```
Repository Interface 用于隔离业务逻辑和持久化实现。

不采用过度复杂的 Clean Architecture。

避免出现无必要的多层抽象，例如：
```
Controller
Application
Domain
Port
Adapter
DAO
ORM
```

---

## 10. Internal Communication
外部通信：
```
REST + JSON
```
内部服务通信：
```
gRPC + Protobuf
```
内部 gRPC 不直接暴露给公共客户端。


---

## 11. Authentication 
使用JWT.

请求：
```http
Authorization: Bearer <token>
```

JWT Payload:
```
{
  "sub": "123",
  "username": "username",
  "role": "author",
  "iat": 123456,
  "exp": 123456
}
```
Gateway 负责基础认证。

涉及资源权限的服务必须执行最终授权检查。

---

## 12. Authorization
采用：`RBAC + Resource Ownership`

角色：
- author
- moderator
- admin
权限与角色分离。

例如：
```
article:create
article:update:own
article:update:any
article:delete:own
article:delete:any
article:moderate
user:ban
tag:manage
```
资源操作必须同时考虑：

- Role
- Permission
- Ownership
- Resource State

## 13. API Response
成功：
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```
失败：
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```
HTTP Status Code 必须正确使用。

Business Code 不能替代 HTTP Status。

## 14. Error Code
不同 Service 使用独立错误码空间。
```
User        100xx
Content     200xx
Interaction 300xx
Social      400xx
Moderation  500xx
File        600xx
```
example:
```
User        100xx
Content     200xx
Interaction 300xx
Social      400xx
Moderation  500xx
File        600xx
```
错误码必须稳定。

禁止随意创建无语义的错误字符串作为程序判断依据。

---

## 15. Event-Driven Architecture
Kafka 用于异步领域事件。

Event 表示：
>某个事件已经发生

Command表示：
>请求某个Service执行操作

Command与Event必须明且区分

---

## 16. Kafka Topics
Topic按Domain划分，例如：
```
user.events
content.events
interaction.events
social.events
moderation.events
```
不为每个 Event 创建独立 Topic，除非存在明确的运维或吞吐需求。

Event Type 放在消息 Envelope 中。

---

## 17. Event Envelope
统一Event Envelope：
```json
{
  "event_id": "uuid",
  "event_type": "article.published",
  "version": 1,
  "occurred_at": "2026-08-18T10:00:00Z",
  "producer": "content-service",
  "aggregate_type": "article",
  "aggregate_id": "10086",
  "aggregate_version": 3,
  "trace_id": "trace-id",
  "payload": {}
}
```
kafka Partition key默认使用
```
aggregate_type + aggregate_id
```
保证同一 Aggregate 的事件尽可能保持顺序。

---

## 18. Event Types
初期包括：
```
user.created

article.created
article.updated
article.submitted
article.published
article.rejected

question.created
question.published

answer.created
answer.accepted

comment.created
like.created
collection.created

user.followed
tag.followed

moderation.approved
moderation.rejected
```
后续新增 Event 必须遵循 Event Naming Convention。

---

## 19. Outbox Pattern
任何：
```
Database State Change
+
Event Publication
```
必须使用 Outbox Pattern。

事务：
```
BEGIN

Update Business Data

Insert Outbox Event

COMMIT
```

然后：
```
Outbox Publisher
      |
      v
Kafka
```

禁止依赖：
```
DB Transaction
+
Direct Kafka Publish
```
作为单一原子事务。

---

## 20. Idempotency
Kafka Consumer 必须支持幂等处理。

系统必须能够处理：

- Duplicate Delivery
- Retry
- Consumer Restart
- Consumer Rebalance

实现方式可根据业务选择：

- Processed Event Table
- Unique Event ID
- Idempotency Key
- Database Unique Constraint

---

## 21. Retry and DLQ
Kafka Consumer失败后：
```
Consume
  |
  v
Process
  |
  +---- Success ---> Commit
  |
  +---- Failure
           |
           v
         Retry
           |
           v
         Retry
           |
           v
         Retry
           |
           v
          DLQ
```
初始：
```
3 retries
```
DLQ 必须保留：

- Event ID
- Event Type
- Aggregate ID
- Original Payload
- Error
- Retry Count
- Trace ID

## 22. Core Business Flows
**22.1 Article Creation**
```
Client
  |
  v
Gateway
  |
  v
Content Service
  |
  +--> Validate
  |
  +--> Create Article
  |
  +--> Create Outbox
  |
  v
Commit
```
新 Article 默认：`DRAFT`

**22.2 Article Submission**
```
Client
  |
  v
Content Service
  |
  +--> Permission Check
  +--> State Check
  +--> Sensitive Word Check
  |
  +--> PENDING_REVIEW
  |
  +--> Outbox
  |
  v
Kafka
```

**22.3 Article Moderation**
```
article.submitted
       |
       v
Moderation Service
       |
       v
Moderation Task
       |
       v
Admin Review
       |
   +---+---+
   |       |
Approve   Reject
   |       |
   v       v
Event     Event
```
审核服务不能直接修改 Content Service 数据库。

**22.4 Article Publishing**
```
moderation.approved
        |
        v
Content Service
        |
        v
Article -> PUBLISHED
        |
        v
article.published
```

**22.5 Search Indexing**
```
article.published
        |
        v
Search Service
        |
        v
Elasticsearch
```
Search Service 不将 Elasticsearch 作为 Source of Truth。

**22.6 Notification**
```
like.created
comment.created
user.followed
answer.accepted
moderation.rejected
        |
        v
Notification Service
        |
        v
Notification
```

**22.7 Feed**
```
article.published
        |
        v
Feed Service
        |
        +--> Normal Author
        |       |
        |       v
        |   Fan-out on Write
        |
        +--> High Follower Author
                |
                v
            Fan-out on Read
```

## 23. Feed Strategy
第一版采用规则推荐。

不引入机器学习推荐模型。

Ranking 可使用：
```
score =
    freshness_score
    + engagement_score
    + author_score
    + tag_match_score
```
Feed 数据使用 Redis Sorted Set 等结构实现。

普通作者采用：
```
Fan-out on Write
```
高关注作者采用：
```
Fan-out on Read
```
阈值由后续压测确定。

Ranking Engine 必须通过抽象接口实现，以便未来替换为更复杂推荐系统。

---

## 24. Search Strategy

Source of Truth：
```
MySQL
```
Search：
```
Elasticsearch
```
必须支持：
```
Article Search
Question Search
Answer Search
Tag Filter
Author Filter
Highlight
Pagination
```
必须支持 Index Rebuild：
```
MySQL
   |
   v
Rebuild Process
   |
   v
Elasticsearch
```

---

## 25. Moderation Strategy
提交内容时进行同步基础规则检查。

例如：
```
Sensitive Word
```
高风险内容可直接阻止。

其他内容进入：
```
PENDING_REVIEW
```
Moderation Service 负责：

- Moderation Task
- Report
- Sensitive Word
- Human Review

Moderation Service 不直接修改 Content DB。

---

## 26. File Strategy
File Service 对外提供统一接口。

例如：
```
POST /files
GET /files/{id}
```
业务系统只保存：
```
file_id
```
而不是依赖物理文件路径。

Storage：
```
type Storage interface {
    Put(...)
    Get(...)
    Delete(...)
}
```
第一版：
```
LocalStorage
```
未来：
```
OSS
S3
MinIO
```
客户端 API 不应因 Storage 实现变化而改变。

---

## 27. Redis Strategy
Redis 用于：

- Cache
- Feed
- Counter
- Rate Limit
- Hot Data
- Unread Count
- Optional Distributed Lock

Redis 不作为核心 Source of Truth。

---

## 28. Observability
所有 Service 必须支持：

- Structured Logging
- Metrics
- Distributed Tracing

日志统一使用 Zap。

禁止生产代码使用：
```
fmt.Println()
```
每个重要请求至少包含：
```
service
request_id
trace_id
user_id
```
业务事件日志还应包含：
```
event_type
aggregate_id
```

---

## 29. Distributed Tracing
使用：
```
OpenTelemetry
```
Trace 应覆盖：
```
Gateway
  |
  v
Service
  |
  v
gRPC
  |
  v
MySQL / Redis
  |
  v
Kafka
  |
  v
Consumer
```
Kafka Event 应尽可能传播 Trace Context

## 30. Metrics
HTTP：
```
http_request_total
http_request_duration_seconds
http_request_errors_total
```
gRPC：
```
grpc_request_total
grpc_request_duration_seconds
grpc_request_errors_total
```
Kafka：
```
kafka_consumer_lag
kafka_consume_errors_total
kafka_retry_total
kafka_dlq_total
```
Business：
```
user_register_total
article_publish_total
question_publish_total
answer_create_total
comment_create_total
like_total
collection_total
report_total
```

## 31. Testing Strategy
测试分为：

- Unit Test
- Integration Test
- Contract Test
- E2E Test
- Load Test
- Failure Test

**31.1 Unit Test**

重点：

- Domain Rule
- State Transition
- Permission
- Ownership
- Ranking
- Moderation Rule

**31.2 Integration Test**

验证：

- MySQL
- Redis
- Kafka
- Elasticsearch

优先使用 Testcontainers。

**31.3 Contract Test**

验证：

- gRPC Contract
- Kafka Event Schema
- Event Version Compatibility


**31.4 E2E Test**
核心 E2E：
```
Register
   ↓
Login
   ↓
Create Article
   ↓
Submit Review
   ↓
Approve
   ↓
Publish
   ↓
Search
   ↓
Like
   ↓
Comment
   ↓
Notification
```

---

## 32. Development Strategy
开发不按照“一个 Service 一个阶段”的方式进行。

采用：
```
Foundation
    ↓
Vertical Business Slice
    ↓
Cross-Service Capability
    ↓
Infrastructure Capability
    ↓
Reliability
    ↓
Performance
    ↓
Production Hardening
```
开发阶段以“可运行、可验收的能力闭环”为基本单位，而不是单纯按照服务数量推进。

---

## 33. Development Phases
**Phase 0 - Foundation**

目标：

建立所有服务共同使用的工程基础设施。

包括：

- Repository Structure
- Service Template
- Configuration
- Environment
- Error System
- Logging
- Request ID
- Trace ID
- gRPC Foundation
- Kafka Foundation
- MySQL Migration
- Redis Client
- Health Check
- Graceful Shutdown
- Docker Compose

验收：
```
docker compose up


All infrastructure services are healthy.


All application services can start.


Gateway health endpoint returns 200.


gRPC health endpoint returns SERVING.
```

**Phase 1 - User and Authentication**

实现：

- Register
- Login
- JWT
- Profile
- Role
- Authorization
- user.created

验收：
```
Register
  ↓
MySQL
  ↓
user.created
  ↓
Kafka


Login
  ↓
JWT


JWT
  ↓
Gateway
  ↓
Protected API
```

**Phase 2 - Article**

实现：

- Article Entity
- Markdown
- Draft
- Edit
- Tag
- Version
- Permission
- Submission
- Outbox

状态：
```
DRAFT
PENDING_REVIEW
PUBLISHED
REJECTED
HIDDEN
DELETED
```

**Phase 3 - Q&A**

实现：

- Question
- Answer
- Accept Answer
- Question State
- Answer State
- Tag
- Markdown

Question 与 Article 必须独立建模。


**Phase 4 - Moderation**

实现：

- Moderation Task
- Sensitive Word
- Report
- Human Review
- Approve
- Reject
- Hide

建立完整：
```
Content
  ↓
Moderation
  ↓
Content
```
事件闭环。


**Phase 5 - Interaction**

实现：

- Like
- Collection
- Comment
- View
- Counter
- Idempotency

重点测试：

- Concurrent Like
- Duplicate Like
- Duplicate Event
- Counter Consistency

**Phase 6 - Social**

实现：

- User Follow
- User Unfollow
- Tag Follow

建立 Social Graph。

**Phase 7 - Notification**

消费：

- like.created
- comment.created
- user.followed
- answer.accepted
- moderation.rejected

提供：

- Notification List
- Unread Count
- Mark Read
- Mark All Read

**Phase 8 - Search**

实现：

- Elasticsearch Index
- Search API
- Filtering
- Highlight
- Ranking
- Index Rebuild

**Phase 9 - Feed**

实现：

- Following Feed
- Tag Feed
- Rule Ranking
- Redis Sorted Set
- Fan-out on Write
- Fan-out on Read

**Phase 10 - File**

实现：

- Upload
- Download
- Metadata
- Validation
- Local Storage
- Storage Interface


**Phase 11 - Admin**
实现：

- User Management
- Content Management
- Comment Management
- Report Management
- Tag Management
- Sensitive Word Management
- Statistics
- Audit Log

**Phase 12 - Observability**

完善：

- OpenTelemetry
- Prometheus
- Grafana
- Jaeger
- Kafka Metrics
- Business Metrics
- Dashboards

建立完整 Trace Demo。

重点演示：
```
HTTP Request
   ↓
Gateway
   ↓
Content Service
   ↓
MySQL
   ↓
Outbox
   ↓
Kafka
   ↓
Moderation
   ↓
Content
   ↓
Search
   ↓
Elasticsearch
```
整个过程应尽可能由同一个 Trace ID 串联。

**Phase 13 - Reliability**

测试：

- Kafka Unavailable
- Redis Unavailable
- Elasticsearch Unavailable
- MySQL Timeout
- gRPC Timeout
- Consumer Crash
- Consumer Restart
- Duplicate Event
- Out-of-order Event
- Network Failure

验证：

- Retry
- Idempotency
- DLQ
- Recovery
- Data Consistency

**Phase 14 - Performance**

进行：

- HTTP Load Test
- gRPC Load Test
- Kafka Throughput Test
- Feed Load Test
- Search Load Test
- Like Concurrency Test

记录：

- RPS
- P50
- P95
- P99
- Error Rate
- Kafka Lag
- DB Connections
- Redis Hit Rate

**Phase 15 - Production Hardening**

完善：

- Rate Limit
- Security
- SQL Injection Protection
- XSS Protection
- Markdown Sanitization
- File Security
- JWT Security
- Password Security
- Timeout
- Retry
- Circuit Breaker
- Graceful Shutdown
- Sensitive Data Protection


---

## 34. Phase Dependencies

Phase 表示能力依赖，不强制要求所有任务严格串行。

基础依赖：
```
Phase 0
   |
   v
Phase 1
```
内容能力：
```
Phase 1
   |
   +--------> Phase 2 Article
   |
   +--------> Phase 3 Q&A
```
审核：
```
Phase 2 ----+
            |
Phase 3 ----+----> Phase 4 Moderation
```
互动：
```
Phase 2 ----+
            |
Phase 3 ----+----> Phase 5 Interaction
```
社交：
```
Phase 1
   |
   v
Phase 6 Social
```
通知：
```
Phase 4 ----+
            |
Phase 5 ----+----> Phase 7 Notification
            |
Phase 6 ----+
```
搜索：
```
Phase 2 ----+
            |
Phase 3 ----+----> Phase 8 Search

Feed：

Phase 2 ----+
Phase 3 ----+
Phase 5 ----+----> Phase 9 Feed
Phase 6 ----+

File：

Phase 2
   |
   v
Phase 10 File

Admin：

Phase 4 ----+
Phase 5 ----+
Phase 6 ----+----> Phase 11 Admin
Phase 1 ----+
```
可观测性：
```
Phase 0
   |
   v
Phase 12 Observability
```
可靠性：
```
Phase 7
Phase 8
Phase 9
   |
   v
Phase 13 Reliability
```
性能：
```
Phase 8
Phase 9
   |
   v
Phase 14 Performance
```
最终：
```
Phase 13
Phase 14
   |
   v
Phase 15 Production Hardening
```
注意：

Phase 2 Article 与 Phase 3 Q&A 在满足 Phase 0 和 Phase 1 依赖后，可以由不同 Agent 并行开发。


---

## 35. Agent Development Rules

AI Agent 必须遵循以下优先级：
```
spec.md
   ↓
plan.md
   ↓
tasks.md
```
Agent 不得自行修改系统架构。

未经明确批准，Agent 不得新增：

- Service
- Database
- Message Broker
- Framework
- Communication Pattern
- External Infrastructure

如果实现过程中发现规范与现实实现冲突：

- 停止扩大修改范围
- 明确报告冲突
- 说明产生冲突的原因
- 提出候选解决方案
- 等待负责人决策

不得静默修改架构。

---

## 36. Task Execution Model

每个 Task 必须包含：
```
Task ID
Objective
Context
Dependencies
Scope
Implementation Requirements
Acceptance Criteria
Tests
Definition of Done
```
Agent 必须只修改当前 Task 所需范围。

禁止：
```
无关重构
未授权架构调整
修改其他服务行为
删除已有测试
为通过测试而降低业务约束
```
如果发现需要修改当前 Task 之外的架构或规范：

必须停止并报告。

---

## 37. Task Dependency Model

Task 必须声明显式依赖。

例如：
```
TASK-USER-001
      |
      v
TASK-USER-002
      |
      v
TASK-CONTENT-001
```
不存在依赖关系的任务可以并行执行。

Agent 不得假设尚未完成的 Task 已经存在。

---

## 38. Definition of Done

一个功能只有满足以下条件才算完成：
```
Implementation completed
Unit tests completed
Integration tests completed where required
Database migration completed where required
API / gRPC contract updated
Event schema updated where required
Logging implemented
Metrics implemented where required
Tracing implemented where required
Error handling implemented
Documentation updated
Acceptance Criteria all passed
No unrelated behavior changed
```
代码能够编译不等于 Task 完成。

测试通过也不自动等于 Task 完成。

---

## 39. Code Quality Rules

代码必须：
```
遵循 Go 官方风格
使用 gofmt
使用 go vet
使用 golangci-lint
避免明显重复代码
避免无必要抽象
避免全局可变状态
正确处理 context
正确处理 error
正确释放资源
正确处理 goroutine 生命周期
```
服务必须支持：
```
Graceful Shutdown
Timeout
Context Cancellation
```

---

## 40. API and Contract Stability

已经发布的 API Contract 不得随意破坏。

修改 API 时必须考虑：
```
Backward Compatibility
Versioning
Client Compatibility
```
gRPC Protobuf：
```
不随意复用 Field Number
删除字段时保留 reserved
不随意修改已有字段语义
```
Kafka Event：
```
Event Version 必须明确
不允许静默修改已有字段语义
消费者应尽可能兼容未来字段
```

---

## 41. Database Migration Rules

所有数据库结构变化必须通过 Migration。

禁止：
```
手动修改生产数据库结构
```
Migration 必须：
```
可重复执行或具有明确版本
可追踪
与代码版本对应
包含必要索引
包含必要约束
```
数据库 Schema 属于 Service 自己的边界。


---

## 42. Security Baseline

必须考虑：
```
Password Hashing
JWT Validation
Authorization
Resource Ownership
SQL Injection
XSS
Markdown Sanitization
File Type Validation
File Size Limit
Path Traversal
Rate Limiting
Sensitive Data Protection
Audit Logging
```
密码等敏感信息禁止明文存储。


---

## 43. Failure Handling Principles

系统必须区分：
```
Expected Business Error
Unexpected System Error
Dependency Failure
Timeout
Retryable Error
Non-Retryable Error
```
不得对所有错误进行无限 Retry。

Retry 必须：
```
有上限
有 Backoff
有明确 Retryable 条件
```

---

## 44. Architecture Constraints

v1.0 阶段冻结以下基础技术：

- Go
- Gin
- gRPC
- Protobuf
- GORM
- MySQL
- Redis
- Kafka
- KRaft
- Elasticsearch
- Docker
- Docker Compose
- OpenTelemetry
- Prometheus
- Grafana
- Jaeger
- Zap
- JWT
- Vue 3
- Vite
- JavaScript

除非出现明确技术阻塞，否则 Agent 不得擅自更换。

---

## 45. Future Extensions

以下能力不属于第一版本核心范围，但架构应保留扩展能力：

- OSS / S3
- Kubernetes
- CI/CD
- AI Content Moderation
- Advanced Recommendation
- Machine Learning Ranking
- Real-time WebSocket Notification
- Multi-region Deployment
- Distributed Cache Optimization
- Advanced Analytics

---

## 46. Final Delivery Criteria

项目最终交付必须满足：

Functional

核心业务完整：

- User
- Article
- Question
- Answer
- Comment
- Like
- Collection
- Follow
- Tag
- Notification
- Feed
- Search
- Moderation
- Report
- File
- Admin

Architectural

必须体现：

- Microservices
- gRPC
- Kafka
- Outbox
- Idempotency
- Redis
- Elasticsearch
- Eventual Consistency
  
Engineering

必须具备：

- Unit Tests
- Integration Tests
- Contract Tests
- E2E Tests
- Load Tests
- Failure Tests

Observability

必须具备：

- Structured Logging
- Metrics
- Distributed Tracing
- Kafka Monitoring

Deployment

必须支持：
```
docker compose up
```
完成本地完整环境启动。