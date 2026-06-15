# GeekAI 系统完整架构图

> 基于 geekai v4.2.7 源码梳理（后端 Go + 前端 Vue3，单体应用 + Docker 部署）。
> 所有图均为 Mermaid，可在 GitHub / VS Code / Typora 中直接渲染。

---

## 1. 总体系统架构图

```mermaid
flowchart TB
    subgraph Clients["🖥️ 客户端层"]
        PC["Web 前端<br/>Vue3 + Element Plus + Pinia"]
        H5["移动端 H5<br/>Vant 组件库"]
        ADM["管理后台<br/>同一 SPA 的 /admin 路由"]
        APP["桌面 / App 客户端<br/>Win / Mac / Linux / Android"]
    end

    subgraph Gateway["🌐 接入层"]
        NGINX["Nginx (geekai-web :8080)<br/>托管 Vue 静态资源 + /api 反向代理"]
    end

    subgraph Backend["⚙️ GeekAI API — Go 单体服务 :5678"]
        direction TB
        subgraph Core["core 核心层"]
            FX["uber-go/fx 依赖注入容器<br/>(main.go 组装全部组件)"]
            GIN["Gin Engine (AppServer)"]
            MW["中间件链<br/>参数清洗 / 全局 Recover<br/>JWT+Redis 双重鉴权 / 静态资源 / 限流"]
            CFG["双层配置<br/>AppConfig ← config.toml<br/>SystemConfig ← DB configs 表"]
        end

        subgraph Handlers["handler — HTTP 处理层"]
            CHATH["ChatHandler<br/>SSE 流式对话"]
            USERH["User / Sms / Captcha<br/>Invite / Redeem"]
            DRAWH["绘画: MJ / SD<br/>DALL-E / 即梦"]
            MEDIAH["Suno 音乐 / Video 视频<br/>MarkMap(WS) / Realtime 语音(WS)"]
            PAYH["Payment / Order / Product"]
            FUNCH["Function 插件 / Menu<br/>Config / Prompt / PowerLog"]
            ADMINH["admin/* 管理后台<br/>用户/模型/Key/订单/审核/仪表盘"]
        end

        subgraph Services["service — 业务服务层"]
            USERS["UserService<br/>算力计费: 锁+事务+流水"]
            TASKS["异步任务服务 ×6<br/>mj / sd / dalle / suno / video / jimeng<br/>(RedisQueue 生产者-消费者)"]
            PAYS["支付服务<br/>Alipay / WxPay / GeekPay易支付"]
            MODS["内容审核<br/>百度 / 腾讯 / Gitee AI"]
            OSSS["UploaderManager<br/>Local / MinIO / 七牛 / 阿里OSS"]
            SMSS["SmsManager<br/>阿里云 / 短信宝"]
            MISCS["License / Captcha / Smtp<br/>WxLogin / Snowflake / Migration"]
        end
    end

    subgraph Data["💾 数据层"]
        MYSQL[("MySQL 8<br/>GORM · 表前缀 geekai_<br/>用户/消息/订单/任务Job")]
        REDIS[("Redis<br/>登录会话 users/{id}<br/>6 条任务队列 · 验证码")]
        LEVELDB[("LevelDB<br/>本地 KV 存储")]
    end

    subgraph External["☁️ 外部服务"]
        LLM["大模型 API (OpenAI 兼容协议)<br/>OpenAI / Claude / 通义 / Kimi<br/>DeepSeek / Gitee AI<br/>← API-Key 池轮换"]
        MJAPI["MidJourney 代理服务"]
        SDAPI["Stable Diffusion WebUI"]
        GENAPI["Suno / Luma / 可灵 / 即梦"]
        PAYAPI["支付宝 / 微信支付 / 易支付网关"]
        SMSAPI["阿里云短信 / 短信宝"]
        MODAPI["百度/腾讯/Gitee 审核 API"]
        OSSAPI["阿里OSS / 七牛 / MinIO"]
        TIKA["Apache Tika :9998<br/>附件文档文本提取"]
        PLUGIN["函数插件 API<br/>微博热搜/头条/早报 (sapi.geekai.me)"]
        WXAPI["微信扫码登录服务"]
    end

    PC --> NGINX
    H5 --> NGINX
    ADM --> NGINX
    APP --> NGINX
    NGINX -->|"/api/* 反代<br/>SSE / WebSocket / REST"| GIN

    GIN --> MW --> Handlers
    FX -.->|启动时注入依赖+注册路由| Handlers
    FX -.->|启动消费者 goroutine| TASKS
    CFG -.-> GIN

    CHATH --> USERS
    CHATH --> MODS
    CHATH -->|流式请求| LLM
    CHATH -->|工具调用回调| PLUGIN
    CHATH -->|附件解析| TIKA
    DRAWH --> TASKS
    MEDIAH --> TASKS
    PAYH --> PAYS
    USERH --> SMSS
    USERH --> MISCS
    ADMINH --> USERS

    TASKS -->|提交/查询任务| MJAPI
    TASKS --> SDAPI
    TASKS --> GENAPI
    TASKS -->|生成结果转存| OSSS
    PAYS <-->|下单 / 异步回调| PAYAPI
    SMSS --> SMSAPI
    MODS --> MODAPI
    OSSS --> OSSAPI
    MISCS --> WXAPI

    Handlers --> MYSQL
    Services --> MYSQL
    MW -->|会话校验| REDIS
    TASKS <-->|BLPOP / RPUSH| REDIS
    Backend --> LEVELDB
```

---

## 2. 后端内部分层与启动流程

```mermaid
flowchart LR
    subgraph Boot["main.go 启动流程 (fx 容器)"]
        direction TB
        S1["1.加载 config.toml<br/>→ AppConfig"] --> S2["2.建立连接<br/>MySQL / Redis / LevelDB"]
        S2 --> S3["3.MigrationService<br/>数据库自动迁移"]
        S3 --> S4["4.fx.Provide<br/>构造全部 Handler / Service"]
        S4 --> S5["5.fx.Invoke<br/>RegisterRoutes 注册路由"]
        S5 --> S6["6.fx.Invoke<br/>启动 6 个任务消费者 goroutine"]
        S6 --> S7["7.AppServer.Run<br/>监听 :5678"]
    end

    subgraph Layers["请求处理分层"]
        direction TB
        L1["middleware<br/>ParameterHandler → errorHandler(Recover)<br/>→ Static → UserAuth/AdminAuth(JWT+Redis)"]
        L2["handler<br/>参数绑定校验 · 组装业务 · resp.SUCCESS/ERROR"]
        L3["service<br/>计费 / 任务队列 / 支付 / 审核 / 存储 / 短信"]
        L4["store<br/>model(GORM实体) · vo(视图对象) · RedisQueue"]
        L1 --> L2 --> L3 --> L4
    end

    Boot -.-> Layers
```

---

## 3. 核心链路①：SSE 流式聊天时序图

```mermaid
sequenceDiagram
    autonumber
    participant U as 前端 ChatPlus.vue
    participant H as ChatHandler
    participant DB as MySQL
    participant R as Redis
    participant T as Tika
    participant AI as 大模型 API
    participant M as 审核服务
    participant P as 插件 API

    U->>H: POST /api/chat/message (SSE 长连接)
    H->>H: 用户级并发锁 userLocks.TryLock
    H->>DB: 查会话/角色/模型/用户/API-Key池
    H->>H: 校验算力·账号状态·Token 上限
    H->>DB: 加载历史上下文 (ContextDeep 条, token 预算裁剪)
    opt 带附件
        H->>T: 提取文档文本拼入 Prompt (图片走 vision 格式)
    end
    H->>AI: POST /v1/chat/completions (stream=true, Key轮换+代理)
    loop 流式响应
        AI-->>H: data: {delta}
        H-->>U: SSEvent "text" (打字机效果)
    end
    opt 模型触发 tool_calls
        H->>P: POST function.Action (拼装参数)
        P-->>H: 插件执行结果
        H-->>U: SSEvent 工具结果
    end
    opt 开启内容审核
        H->>M: Moderate(回复全文)
        M-->>H: 违规 → 记录+提示, 不保存
    end
    H->>DB: 保存 ChatMessage ×2 (问+答)
    H->>DB: DecreasePower 扣算力 + PowerLog 流水
    H->>DB: 首次对话创建 ChatItem (截断标题)
    H-->>U: SSEvent "complete" + "end"

    Note over U,H: 停止生成: GET /api/chat/stop?session_id=x<br/>→ 从 LMap 取 context.CancelFunc 掐断上游请求
```

---

## 4. 核心链路②：异步 AI 任务（MJ/SD/DALL-E/Suno/视频/即梦 通用模式）

```mermaid
flowchart TB
    A["前端提交任务<br/>POST /api/mj/image"] --> B["Handler<br/>校验算力 → 创建 Job 记录"]
    B --> DB[("MySQL<br/>geekai_mj_jobs<br/>progress=0")]
    B --> Q[("Redis List<br/>MidJourney_Task_Queue")]
    B -->|立即返回 jobId| A

    subgraph Workers["Service 后台 goroutine ×3 (启动时由 fx.Invoke 拉起)"]
        C["① Run 消费者<br/>BLPOP 阻塞取任务"]
        D["② SyncTaskProgress<br/>轮询进度, 10min 超时标失败"]
        E["③ DownloadImages<br/>下载成品图转存 OSS"]
    end

    Q -- BLPOP --> C
    C -->|"Imagine / Upscale / Variation<br/>Blend / SwapFace"| API["MidJourney 代理 API"]
    C -->|回写 task_id / channel_id| DB
    D -->|QueryTask| API
    D -->|更新 progress / 失败原因| DB
    E -->|img_url 为空且 progress=100| DB
    E --> OSS["OSS 存储<br/>Local/MinIO/七牛/阿里"]

    F["前端轮询任务列表<br/>GET /api/mj/jobs"] --> DB

    R["重启恢复: Run 启动时把 DB 中<br/>task_id='' 的任务重新入队"] -.-> Q
```

---

## 5. 认证与计费体系

```mermaid
flowchart LR
    subgraph Auth["JWT + Redis 双重认证"]
        direction TB
        A1["登录: 验证码校验 → 密码加盐比对"] --> A2["签发 JWT (HMAC, Session.SecretKey)"]
        A2 --> A3["写 Redis 会话键 users/{id}"]
        A4["每次请求: 解析 JWT → 校验 expired<br/>→ 检查 Redis 键存在 → c.Set(user_id)"]
        A5["登出 = 删 Redis 键<br/>JWT 未过期也立即失效"]
        A3 --> A4
        A4 -.-> A5
    end

    subgraph Power["算力 (Power) 计费闭环"]
        direction TB
        P1["获取: 注册赠送 / 邀请奖励<br/>充值订单 / 兑换码"] --> P2["UserService<br/>Increase/DecreasePower<br/>互斥锁 + DB事务"]
        P3["消耗: 对话(按模型power)<br/>绘画/音乐/视频任务"] --> P2
        P2 --> P4[("geekai_power_logs<br/>每笔流水: 类型/数量/余额/备注")]
    end

    subgraph Pay["支付链路"]
        direction TB
        Y1["选套餐 Product → 创建 Order"] --> Y2["PayService 生成支付链接/二维码<br/>Alipay / WxPay / GeekPay"]
        Y2 --> Y3["支付网关异步回调 notify"]
        Y3 --> Y4["验签 → 更新订单 → IncreasePower 发货"]
        Y5["StartSyncOrders 定时兜底<br/>同步未支付订单状态"] -.-> Y4
    end

    Auth ~~~ Power ~~~ Pay
```

---

## 6. 部署拓扑图（docker-compose）

```mermaid
flowchart TB
    USER(("用户浏览器"))

    subgraph Host["Docker 宿主机"]
        WEB["geekai-web<br/>Nginx :8080<br/>Vue 静态资源 + /api 反代"]
        API["geekai-api<br/>Go 单体 :5678<br/>挂载 config.toml / static / leveldb"]
        MY["geekai-mysql<br/>MySQL 8.0 (3307→3306)"]
        RD["geekai-redis<br/>Redis 6.0 (6380→6379)"]
        TK["geekai-tika<br/>Apache Tika (9999→9998)"]
    end

    EXT["外部云服务<br/>大模型 API / MJ·SD·Suno / 支付网关<br/>短信 / 审核 / OSS / 微信登录"]

    USER -->|HTTP/HTTPS| WEB
    WEB -->|proxy_pass| API
    API --> MY
    API --> RD
    API --> TK
    API <--> EXT

    MY -.健康检查通过后启动.-> API
    RD -.健康检查通过后启动.-> API
```

---

## 7. 目录 ↔ 架构对照表

| 架构层 | 目录 | 关键文件 |
|--------|------|----------|
| 启动/组装 | `api/main.go` | fx 容器：Provide 依赖 + Invoke 路由/消费者 |
| 核心 | `api/core/` | `app_server.go`(Gin封装) · `config.go` · `middleware/auth.go` |
| HTTP 层 | `api/handler/` | `chat_handler.go`+`chat_openai_handler.go`(SSE对话核心) · `admin/`(后台) |
| 业务层 | `api/service/` | `user_service.go`(计费) · `mj/ sd/ dalle/ suno/ video/ jimeng/`(异步任务) · `payment/ sms/ moderation/ oss/` |
| 数据层 | `api/store/` | `mysql.go` · `redis_queue.go` · `model/`(25+实体) · `vo/`(视图对象) |
| 前端 | `web/src/` | `views/`(页面) · `store/`(Pinia) · `components/` |
| 部署 | `docker/` | `docker-compose.yaml` · `conf/nginx/` |
