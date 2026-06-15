# Simple Bank 重写计划

本文档用于按 Tech School `simplebank` 课程顺序，从零重写一个可运行、可测试、可部署的 Simple Bank 项目。

参考来源：

- 上游仓库：<https://github.com/techschool/simplebank>
- 当前本地参考实现：`D:\AAAJove\simplebank\simplebank`
- 后端课程顺序：上游 `README.md` 的 Backend course videos
- 前端课程顺序：上游 `README.md` 的 Frontend crash course videos

## 目标

1. 按课程顺序重写，不直接复制现有实现。
2. 每个阶段都能独立运行、测试、提交。
3. 最终得到一个包含数据库、REST API、gRPC API、异步任务、邮件验证、容器化、本地编排、CI、可选 AWS EKS 部署、可选 Vue 前端的完整项目。
4. 每个关键能力都保留自己的学习笔记、验证命令和提交记录，方便复盘。

## 重写原则

- 先实现最小可运行版本，再逐步扩展。
- 数据库 schema、SQL 查询、生成代码和业务代码分开演进。
- 生成代码必须可重复生成，包括 `sqlc`、`mockgen`、`protoc`、Swagger/statik。
- 所有涉及余额变更的逻辑必须以事务和测试为先。
- 密码、token、session、邮件验证码、生产 secrets 不硬编码。
- AWS 部署阶段默认视为可选，先完成本地 Docker Compose 与 CI，再进入云资源。
- 每完成一个 lecture 或一个小 milestone，提交一次有意义的 commit。

## 建议工作方式

### 仓库策略

推荐两种方式二选一：

1. 保留当前仓库作为参考，在新分支或新目录中重写。
2. 新建一个空仓库，将本仓库只作为只读参考。

如果在当前仓库内重写，建议：

- 新建分支：`rewrite-from-scratch`
- 删除或移动不需要的生成产物，例如根目录二进制 `simplebank`
- 先建立基础目录，再逐步添加代码
- 不在同一个 commit 中混合多个课程主题

### Commit 节奏

建议 commit 命名：

- `chore: bootstrap go module and tooling`
- `db: add initial schema and migrations`
- `db: generate sqlc queries for accounts`
- `api: add gin account endpoints`
- `auth: add paseto token maker`
- `grpc: add user service rpc`
- `worker: enqueue verify email task`
- `deploy: add docker compose`

### 每阶段固定检查

每个阶段结束前至少执行：

```bash
go test ./...
```

涉及前端时执行：

```bash
cd frontend
npm run build
npm run test:unit
```

涉及容器时执行：

```bash
docker compose up --build
```

## 总体里程碑

| 阶段 | 课程范围 | 主题 | 主要产物 |
| --- | --- | --- | --- |
| M0 | Lecture #0 | 环境与项目骨架 | Go module、Makefile、Docker/Postgres 工具链 |
| M1 | Lecture #1-#10 | 数据库与事务 | DBML、migration、sqlc、事务、DB 测试、CI |
| M2 | Lecture #11-#22 | REST API 与认证 | Gin API、config、mock、用户、密码、JWT/PASETO、middleware |
| M3 | Lecture #23-#36 | 容器化与 AWS 部署 | Dockerfile、Compose、ECR/RDS/Secrets/EKS/Ingress/HTTPS/GitHub Actions |
| M4 | Lecture #37-#53 | Session、gRPC、Gateway、日志 | session、protobuf、gRPC、gateway、Swagger、partial update、logger |
| M5 | Lecture #54-#64 | 异步任务与邮件验证 | Redis、Asynq、邮件发送、verify email、gRPC 单测 |
| M6 | Lecture #65-#77 | 稳定性、安全与补强 | sqlc v2、pgx、RBAC、graceful shutdown、CORS、JWT v5 |
| M7 | Frontend #1-#9 | Vue 前端 | 登录页、路由、auth state、API 联调 |

## M0：环境与项目骨架

对应 Lecture #0：Setup development environment on Windows: WSL2 + Go + VSCode + Docker + Make + Sqlc

### 目标

建立从零重写的最小工程：能启动 Go module，能用 Makefile 管理常用命令，能启动 Postgres。

### 任务

- 确认本机工具：Go、Docker Desktop、Make、migrate、sqlc、mockgen、protoc、Node.js。
- 初始化 Go module，建议临时使用自己的模块路径，避免一直依赖 `github.com/techschool/simplebank`。
- 建立基础目录：
  - `db/migration`
  - `db/query`
  - `db/sqlc`
  - `api`
  - `util`
  - `token`
  - `proto`
  - `pb`
  - `gapi`
  - `worker`
  - `mail`
  - `doc`
- 建立 `.gitignore`，排除二进制、临时文件、`.env`、覆盖率文件、node_modules。
- 建立 `Makefile` 的初版命令：`postgres`、`createdb`、`dropdb`、`migrateup`、`migratedown`、`sqlc`、`test`。
- 建立 `app.env.example`，不要提交真实 secrets。

### 验收

```bash
go test ./...
make postgres
make createdb
```

### DoD

- 本地能启动 Postgres。
- 空项目测试能通过。
- README 或本计划记录已安装工具版本。

## M1：数据库、SQLC、事务与 CI

对应 Lecture #1-#10：Working with database [Postgres]

### M1.1 数据库设计

对应 Lecture #1：Design DB schema and generate SQL code with dbdiagram.io

任务：

- 设计 `accounts`、`entries`、`transfers` 三张核心表。
- 使用 DBML 维护 schema：`doc/db.dbml`。
- 生成 SQL schema：`doc/schema.sql`。
- 明确字段：
  - `accounts`: `id`、`owner`、`balance`、`currency`、`created_at`
  - `entries`: `id`、`account_id`、`amount`、`created_at`
  - `transfers`: `id`、`from_account_id`、`to_account_id`、`amount`、`created_at`
- 明确外键、索引、金额正负约定。

验收：

```bash
make db_schema
```

### M1.2 本地 Postgres 与 Migration

对应 Lecture #2-#3

任务：

- 用 Docker 启动 Postgres。
- 建立 `000001_init_schema.up.sql` 与 `.down.sql`。
- 使用 `golang-migrate` 管理版本。
- 确保 migration 可重复 up/down。

验收：

```bash
make migrateup
make migratedown
make migrateup
```

### M1.3 SQLC CRUD

对应 Lecture #4

任务：

- 建立 `sqlc.yaml`。
- 为 `accounts`、`entries`、`transfers` 编写 SQL 查询。
- 生成 Go 代码到 `db/sqlc`。
- 建立 `Store` 抽象，为后续事务和 mock 做准备。

验收：

```bash
make sqlc
go test ./db/...
```

### M1.4 DB 单元测试与随机数据

对应 Lecture #5

任务：

- 建立 `db/sqlc/main_test.go`，连接测试数据库。
- 编写 account、entry、transfer 的 CRUD 测试。
- 建立 `util/random.go` 生成 owner、money、currency。
- 测试使用真实 Postgres，不 mock 数据库。

验收：

```bash
go test -v ./db/sqlc
```

### M1.5 转账事务

对应 Lecture #6-#9

任务：

- 实现 `execTx`。
- 实现 `TransferTx`：
  - 创建 transfer 记录
  - 创建 from entry
  - 创建 to entry
  - 扣减 from account
  - 增加 to account
- 用并发测试覆盖多笔转账。
- 理解并避免 deadlock：固定账户更新顺序。
- 记录 PostgreSQL 隔离级别与锁行为笔记。

验收：

```bash
go test -v ./db/sqlc -run TestTransferTx
go test -v ./db/sqlc -run TestTransferTxDeadlock
```

### M1.6 GitHub Actions

对应 Lecture #10

任务：

- 建立 `.github/workflows/test.yml`。
- CI 启动 Postgres service。
- 安装 migrate。
- 执行 migration 与 `go test -v -cover -short ./...`。

验收：

- push 或 PR 后 CI 通过。
- CI 日志中能看到 migration 和测试执行。

## M2：REST API、配置、Mock 与认证

对应 Lecture #11-#22：Building RESTful HTTP JSON API [Gin]

### M2.1 Gin REST API

对应 Lecture #11

任务：

- 建立 `api/server.go`。
- 实现账户 API：
  - `POST /accounts`
  - `GET /accounts/:id`
  - `GET /accounts?page_id=&page_size=`
- 实现标准错误响应。
- 使用 validator 校验 URI、query、JSON body。

验收：

```bash
go test ./api/...
make server
```

### M2.2 配置管理

对应 Lecture #12

任务：

- 使用 Viper 从 `app.env` 和环境变量加载配置。
- 配置项至少包括：
  - `DB_SOURCE`
  - `HTTP_SERVER_ADDRESS`
  - `TOKEN_SYMMETRIC_KEY`
  - `ACCESS_TOKEN_DURATION`
  - `REFRESH_TOKEN_DURATION`
  - `REDIS_ADDRESS`
  - `EMAIL_*`
- 建立 `app.env.example`。

验收：

```bash
go test ./util/...
```

### M2.3 HTTP API Mock 测试

对应 Lecture #13-#14、#18

任务：

- 使用 gomock 生成 `db/mock/store.go`。
- 测试 account API：
  - 正常
  - 参数错误
  - 未找到
  - DB error
- 实现 transfer API：
  - `POST /transfers`
  - 自定义 currency validator
  - 校验转出账户与转入账户币种一致
- 使用 custom gomock matcher 校验 hashed password 等复杂参数。

验收：

```bash
make mock
go test -v ./api
```

### M2.4 用户、密码与错误处理

对应 Lecture #15-#17

任务：

- 新增 users migration。
- users 表字段：
  - `username`
  - `hashed_password`
  - `full_name`
  - `email`
  - `password_changed_at`
  - `created_at`
- 为 users 添加 SQL 查询和测试。
- 使用 bcrypt hash password。
- 正确处理唯一键冲突、外键冲突、未找到等 DB 错误。

验收：

```bash
make migrateup
make sqlc
go test ./...
```

### M2.5 JWT/PASETO 与登录

对应 Lecture #19-#22

任务：

- 实现 token maker interface。
- 实现 JWT maker。
- 实现 PASETO maker。
- 实现 `POST /users` 创建用户。
- 实现 `POST /users/login` 登录并返回 access token。
- 实现 auth middleware。
- 给账户和转账 API 加权限规则。

验收：

```bash
go test ./token/...
go test ./api/...
```

## M3：Docker、本地编排与 AWS 部署

对应 Lecture #23-#36：Deploying the application to production [Kubernetes + AWS]

这一阶段分为“必须完成的本地容器化”和“可选完成的 AWS 云部署”。如果重写目标主要是学习后端，AWS 可推迟，避免额外费用。

### M3.1 Docker 镜像

对应 Lecture #23

任务：

- 编写多阶段 `Dockerfile`。
- 输出尽量小的 runtime image。
- 将 `app.env`、`start.sh`、`wait-for.sh` 的职责拆清楚。
- 确认容器内能运行 migration 和 server。

验收：

```bash
docker build -t simplebank:rewrite .
docker run --rm simplebank:rewrite
```

### M3.2 Docker Network 与 Compose

对应 Lecture #24-#25、#68

任务：

- 用 docker network 连接 app 与 Postgres。
- 编写 `docker-compose.yaml`：
  - `postgres`
  - `redis`
  - `api`
- 配置端口映射和 Postgres volume。
- 用 wait-for 或 healthcheck 控制启动顺序。

验收：

```bash
docker compose up --build
```

### M3.3 AWS 基础设施

对应 Lecture #26-#36

任务：

- 创建 AWS 账号并设置预算告警。
- 创建 ECR repository。
- 创建 RDS Postgres。
- 使用 AWS Secrets Manager 管理生产配置。
- 创建 EKS cluster。
- 配置 kubectl/k9s。
- 编写 Kubernetes manifests：
  - deployment
  - service
  - ingress
  - issuer
  - aws-auth
- 配置 Route53 域名。
- 配置 NGINX ingress。
- 配置 cert-manager 与 Let's Encrypt。
- 配置 GitHub Actions 自动 build、push、deploy。

验收：

```bash
kubectl get pods
kubectl get svc
kubectl get ingress
```

### AWS 成本控制

- 只在真正学习部署时创建 EKS、RDS、NAT、Load Balancer。
- 课程结束当天清理不用的云资源。
- 所有云资源用 tag 标识：`project=simplebank-rewrite`。
- 保留“销毁清单”，防止忘记收费资源。

## M4：Session、gRPC、Gateway、Swagger 与日志

对应 Lecture #37-#53：Advanced Backend Topics [Sessions + gRPC]

### M4.1 Refresh Session

对应 Lecture #37

任务：

- 新增 `sessions` 表。
- 登录时创建 session。
- 返回 access token 与 refresh token。
- 记录 user agent、client IP、blocked、expires_at。
- 实现 refresh token API。

验收：

```bash
go test ./db/sqlc -run Session
go test ./api -run Login
```

### M4.2 DB 文档

对应 Lecture #38

任务：

- 使用 DBML 生成 dbdocs。
- 保持 `doc/db.dbml`、migration、`doc/schema.sql` 一致。

验收：

```bash
make db_docs
make db_schema
```

### M4.3 gRPC 基础

对应 Lecture #39-#44

任务：

- 建立 proto 文件：
  - `user.proto`
  - `rpc_create_user.proto`
  - `rpc_login_user.proto`
  - `service_simple_bank.proto`
- 使用 protoc 生成 Go、gRPC、grpc-gateway 代码。
- 建立 `gapi/server.go`。
- 实现 create user 和 login user gRPC API。
- 启动 gRPC server。
- 用 Evans 或 grpcurl 调用。
- 提取 metadata 中的 user agent、client IP。

验收：

```bash
make proto
make server
evans --host localhost --port 9090 -r repl
```

### M4.4 gRPC Gateway 与 Swagger

对应 Lecture #43、#45-#46

任务：

- 通过 grpc-gateway 同时提供 HTTP JSON API。
- 生成 Swagger JSON。
- 将 Swagger UI 静态文件嵌入 Go binary。
- 确保 `/swagger/` 可访问。

验收：

```bash
curl http://localhost:8080/swagger/
```

### M4.5 参数验证、Migration in Go、Partial Update

对应 Lecture #47-#50

任务：

- 建立 `val` 包，封装字段校验。
- gRPC 错误返回兼顾 human friendly 和 machine friendly。
- 在 Go 启动流程中执行 migration。
- 使用 sqlc nullable params 实现用户部分更新。
- 实现 `UpdateUser` gRPC API。

验收：

```bash
go test ./gapi/...
go test ./val/...
```

### M4.6 gRPC 授权与结构化日志

对应 Lecture #51-#53

任务：

- gRPC API 加 authorization。
- 使用 zerolog 编写 gRPC unary interceptor。
- 编写 HTTP logger middleware。
- 日志字段包含 protocol、method、status、duration、user agent、client IP。

验收：

```bash
go test ./gapi/...
curl http://localhost:8080/v1/...
```

## M5：异步任务、Redis、邮件与验证

对应 Lecture #54-#64：Asynchronous processing with background workers [Asynq + Redis]

### M5.1 Worker 基础

对应 Lecture #54-#58

任务：

- 使用 Redis 作为消息队列。
- 使用 Asynq 实现 task distributor 与 task processor。
- 定义 `TaskSendVerifyEmail`。
- Web server 创建用户后 enqueue task。
- 支持任务延迟、重试、错误日志。
- 在 DB transaction 内安全派发任务。

验收：

```bash
make redis
go test ./worker/...
```

### M5.2 邮件发送

对应 Lecture #59-#60

任务：

- 建立 `mail` 包。
- 支持 Gmail SMTP 或本地 Mailhog/假 SMTP。
- 邮件测试默认可跳过或用测试配置运行。
- 配置 VSCode/test flag 或 Go test short mode。

验收：

```bash
go test -short ./mail/...
```

### M5.3 Email Verification

对应 Lecture #61-#64

任务：

- 新增 `verify_emails` 表。
- 创建用户时生成 verify email 记录。
- 异步发送验证邮件。
- 实现 verify email API。
- 验证成功后激活用户或标记 `is_email_verified`。
- 为涉及 DB、Redis、task distributor 的 gRPC API 编写 mock 测试。
- 测试需要鉴权的 gRPC API。

验收：

```bash
go test ./db/sqlc -run VerifyEmail
go test ./gapi/...
```

## M6：稳定性、安全与现代化补强

对应 Lecture #65-#77：Improve the stability and security of the server

### M6.1 sqlc v2 与 pgx

对应 Lecture #65-#67

任务：

- 将 `sqlc.yaml` 升级到 v2 格式。
- 从 `lib/pq` 切换到 `pgx/v5`。
- 更新 DB 错误处理逻辑。
- 确认所有 DB 测试继续通过。

验收：

```bash
make sqlc
go test ./db/sqlc ./api ./gapi
```

### M6.2 工具安装与依赖治理

对应 Lecture #69、#75、#77

任务：

- 固定工具安装方式，例如 `go install tool@version`。
- 记录 Go 版本对 for-loop 变量行为的影响。
- 升级 JWT package 到 v5。
- 定期执行依赖审计和 `go test ./...`。

验收：

```bash
go mod tidy
go test ./...
```

### M6.3 RBAC

对应 Lecture #70

任务：

- users 表增加 `role` 字段。
- 定义角色：`depositor`、`banker`。
- 管理员或 banker 才能访问指定接口。
- 测试普通用户不能越权。

验收：

```bash
go test ./api ./gapi -run Role
```

### M6.4 EKS 生产补强

对应 Lecture #71-#73

任务：

- 配置 EKS 到 RDS/Redis 的 security group。
- 部署 gRPC + HTTP server 到 EKS。
- 验证 HTTP ingress 与 gRPC ingress。
- 整理 AWS 资源清理清单。

验收：

```bash
kubectl get ingress
grpcurl <grpc-domain>:443 list
curl https://<http-domain>/swagger/
```

### M6.5 Graceful Shutdown 与 CORS

对应 Lecture #74、#76

任务：

- 对 gRPC server、HTTP gateway、worker 增加 graceful shutdown。
- 捕获 `SIGINT`、`SIGTERM`。
- 给 Vue 前端配置 CORS。
- 确认本地和 Docker Compose 都能正常退出。

验收：

```bash
go test ./...
docker compose up --build
```

手动验证：

- 启动服务后按 Ctrl+C。
- 日志应显示 gRPC、HTTP、worker 都已停止。

## M7：Vue 前端 Crash Course

对应 Frontend Lecture #1-#9。

这个阶段可在后端 M2 登录 API 完成后开始，也可以等 M6 后统一做。

### 目标

实现一个最小可用前端，完成登录、保存登录状态、调用后端 API、展示用户信息、退出登录。

### 任务顺序

1. 建立 Vue 3 + Vite 项目。
2. 配置 Vue Router。
3. 使用 PrimeVue、PrimeFlex 构建登录表单。
4. 用 `v-model` 和 computed 管理表单状态。
5. 用 axios 调用后端 login API。
6. 联调 CORS。
7. 将 auth state 存到 store/localStorage。
8. 用 props 拆分用户信息组件。
9. 用 emit 实现 logout。

### 验收

```bash
cd frontend
npm install
npm run build
npm run test:unit
```

手动验证：

- 登录失败显示错误。
- 登录成功跳转首页。
- 刷新页面后仍能识别登录状态。
- 点击 logout 清除状态并返回登录页。

## 课程映射清单

### Backend

| Lecture | 主题 | 重写检查点 |
| --- | --- | --- |
| #0 | Windows/WSL2/Go/Docker/Make/sqlc 环境 | 工具链和项目骨架可用 |
| #1 | DB schema 与 DBML | `doc/db.dbml`、schema 草图 |
| #2 | Docker + Postgres | 本地 Postgres 可运行 |
| #3 | DB migration | migration up/down 可重复 |
| #4 | sqlc CRUD | 查询和生成代码可用 |
| #5 | DB CRUD 单测 | 数据库测试通过 |
| #6 | 事务封装 | `execTx` 与 `TransferTx` |
| #7 | transaction lock/deadlock | 并发测试暴露锁行为 |
| #8 | 避免 deadlock | 固定账户更新顺序 |
| #9 | 隔离级别 | 记录 PostgreSQL 事务笔记 |
| #10 | GitHub Actions | CI 跑 migration 和测试 |
| #11 | Gin REST API | account REST endpoints |
| #12 | Viper config | `util/config.go` |
| #13 | Mock DB 测试 API | gomock HTTP tests |
| #14 | Transfer API validator | 转账 API 和币种校验 |
| #15 | Users table | users migration/query |
| #16 | DB errors | 错误映射 |
| #17 | Bcrypt | 密码 hash/verify |
| #18 | Custom gomock matcher | 创建用户测试更严格 |
| #19 | PASETO vs JWT | token 策略笔记 |
| #20 | JWT/PASETO maker | token 包测试通过 |
| #21 | Login API | 登录返回 token |
| #22 | Auth middleware | API 授权规则 |
| #23 | Minimal Docker image | 多阶段 Dockerfile |
| #24 | Docker network | app 与 DB 容器互通 |
| #25 | Docker Compose | compose 一键启动 |
| #26 | AWS account | 预算与 IAM 准备 |
| #27 | ECR GitHub Actions | 自动构建推送镜像 |
| #28 | RDS | 生产 Postgres |
| #29 | Secrets Manager | 生产 secrets |
| #30 | EKS architecture | 集群创建 |
| #31 | kubectl/k9s | 集群连接 |
| #32 | Deploy to EKS | deployment/service |
| #33 | Route53 | 域名解析 |
| #34 | Ingress | HTTP/gRPC 路由 |
| #35 | TLS cert | cert-manager/Let's Encrypt |
| #36 | Auto deploy | release 分支自动部署 |
| #37 | Refresh session | sessions 表与刷新 token |
| #38 | DB docs | dbdocs/schema dump |
| #39 | gRPC intro | gRPC 基础调用 |
| #40 | Protobuf | proto 与生成代码 |
| #41 | gRPC server | gRPC server 可运行 |
| #42 | gRPC user APIs | create/login RPC |
| #43 | gRPC gateway | HTTP JSON gateway |
| #44 | gRPC metadata | user agent/client IP |
| #45 | Swagger docs | 生成 OpenAPI |
| #46 | Embed static files | Swagger UI 嵌入 binary |
| #47 | gRPC validation | 结构化验证错误 |
| #48 | Migrations in Go | 启动时执行 migration |
| #49 | sqlc nullable update | partial update query |
| #50 | gRPC update API | update user RPC |
| #51 | gRPC authorization | RPC 权限控制 |
| #52 | gRPC structured logs | unary interceptor |
| #53 | HTTP logger middleware | HTTP 结构化日志 |
| #54 | Asynq + Redis | worker/distributor |
| #55 | Integrate worker | server 注入 task distributor |
| #56 | Task in DB transaction | transaction 内发任务 |
| #57 | Worker logs/errors | worker 错误处理 |
| #58 | Delay tasks | 延迟与重试策略 |
| #59 | Gmail SMTP | 邮件发送 |
| #60 | Skip tests/config flag | 外部依赖测试隔离 |
| #61 | Email verification design | verify_emails 表 |
| #62 | Email verification API | verify email RPC/API |
| #63 | Mock DB & Redis tests | 多依赖 gRPC 单测 |
| #64 | Authenticated gRPC tests | 鉴权 gRPC 测试 |
| #65 | sqlc v2 | 升级配置 |
| #66 | pgx | 切换 DB driver |
| #67 | PGX DB errors | pgx 错误处理 |
| #68 | Compose port/volume | compose 持久化 |
| #69 | Binary packages | 工具版本固定 |
| #70 | RBAC | 角色权限 |
| #71 | EKS security group | EKS/RDS/Redis 网络权限 |
| #72 | Deploy gRPC + HTTP | 生产双协议部署 |
| #73 | AWS cost | 资源清理 |
| #74 | Graceful shutdown | gRPC/HTTP/worker 优雅退出 |
| #75 | Go for-loop trap | Go 版本差异记录 |
| #76 | CORS + VueJS | 前后端联调 |
| #77 | JWT v5 | JWT 依赖升级 |

### Frontend

| Lecture | 主题 | 重写检查点 |
| --- | --- | --- |
| #1 | Vue reactive app | Vite/Vue 项目启动 |
| #2 | Vue router/component | 登录页与首页路由 |
| #3 | Login form | PrimeVue 登录表单 |
| #4 | v-model/computed | 表单状态与校验 |
| #5 | HTTP request | axios 调用后端 |
| #6 | CORS | 前后端跨域打通 |
| #7 | Auth state | 保存登录状态 |
| #8 | Props | 用户信息组件 |
| #9 | Logout emit | 退出登录 |

## 验证矩阵

| 能力 | 命令 | 通过标准 |
| --- | --- | --- |
| Go 编译与单测 | `go test ./...` | 全部 package 通过 |
| DB migration | `make migrateup && make migratedown && make migrateup` | 可重复执行 |
| SQLC | `make sqlc` | 无生成错误且测试通过 |
| Mock | `make mock` | mock 文件可生成，API 测试通过 |
| Protobuf | `make proto` | `pb` 和 Swagger 产物可生成 |
| REST API | curl/Postman | 账户、用户、登录、转账接口符合预期 |
| gRPC API | Evans/grpcurl | create/login/update/verify RPC 可调用 |
| Worker | Redis + Asynq logs | 用户创建后任务入队并被处理 |
| Docker Compose | `docker compose up --build` | postgres、redis、api 都健康 |
| Frontend | `npm run build` | 前端构建成功 |
| CI | GitHub Actions | push/PR 自动测试通过 |

## 风险与规避

| 风险 | 影响 | 规避 |
| --- | --- | --- |
| 直接复制参考项目 | 学不到设计和排错 | 每个阶段先自己实现，再对照参考 |
| migration 与 DBML 不一致 | 生成文档失真 | 每次 schema 改动同时更新 DBML、migration、query |
| 生成代码不可重复 | 后续重构困难 | 所有生成命令进入 Makefile |
| 并发转账死锁 | 资金业务错误 | 并发测试，固定账户更新顺序 |
| 密码/token/secrets 泄露 | 安全事故 | 只提交 example env，真实 secrets 放环境变量或 Secrets Manager |
| 外部邮件测试不稳定 | CI flaky | short mode 跳过真实 SMTP，使用 mock/fake sender |
| AWS 持续扣费 | 费用风险 | 预算告警，部署阶段结束立即清理 |
| Windows 与 shell 命令差异 | 命令执行失败 | 优先 Makefile + Docker，必要时记录 PowerShell 等价命令 |

## 学习笔记要求

建议在 `notes/` 下按阶段写短笔记：

- `notes/01-database.md`
- `notes/02-rest-api.md`
- `notes/03-docker-deploy.md`
- `notes/04-grpc.md`
- `notes/05-worker-email.md`
- `notes/06-security-stability.md`
- `notes/07-frontend.md`

每篇笔记包含：

- 本阶段核心概念
- 遇到的问题
- 最终解决方式
- 关键命令
- 与参考实现不同的设计选择

## 最终交付物

- 后端服务：
  - REST API
  - gRPC API
  - gRPC Gateway
  - Swagger UI
  - Postgres persistence
  - Redis async worker
  - email verification
  - token/session auth
  - RBAC
  - graceful shutdown
- 本地运行：
  - `make server`
  - `docker compose up --build`
- 自动化：
  - migration
  - sqlc
  - mockgen
  - protoc
  - Go tests
  - GitHub Actions
- 可选生产部署：
  - ECR
  - RDS
  - Secrets Manager
  - EKS
  - Ingress
  - TLS
  - Route53
- 可选前端：
  - Vue login flow
  - auth state
  - API integration

## 推荐执行顺序

1. 完成 M0-M1，确保数据库和事务扎实。
2. 完成 M2，得到可用 REST 后端。
3. 先做 M3.1-M3.2，本地容器化完成后暂缓 AWS。
4. 完成 M4，切入 gRPC/gateway/logging。
5. 完成 M5，补齐异步邮件验证。
6. 完成 M6，本地稳定性、安全和依赖现代化。
7. 如果需要部署，再执行 M3.3 和 M6.4。
8. 如果需要完整体验，再执行 M7 前端。

## 第一周可执行清单

第一周只追求打稳基础，不碰 REST/gRPC/AWS。

- Day 1：整理环境、初始化 module、Makefile、Postgres。
- Day 2：设计 DBML 和初始 migration。
- Day 3：接入 sqlc，完成 accounts CRUD。
- Day 4：完成 entries/transfers CRUD 与测试。
- Day 5：实现 `TransferTx`。
- Day 6：写并发转账测试，修 deadlock。
- Day 7：补 CI，整理数据库阶段笔记。

完成第一周后，项目应达到：

```bash
make migrateup
make sqlc
go test ./db/...
```

全部通过。
