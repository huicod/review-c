# review-c — Scope 锚点文档

> 本文件是绑定到本 scope 的 Worker Chat 的单服务锚点。umbrella 级真相源仍在 `../AGENTS.md` 与 `../.vibe/genesis/v1/`，这里只放"当前服务此刻要做什么"。

## 🎯 服务档案

| 字段 | 值 |
|------|-----|
| **Scope name** | `review-c` |
| **Kind** | `service`（C 端薄 BFF / HTTP 网关） |
| **Umbrella SYS-ID** | SYS-02 · WBS System 3 |
| **Registry 路径** | `./review-c`（见 `../.vibe/coord/registry.yaml`） |
| **Git 仓库** | 独立 git repo（与 umbrella sibling） |
| **Path 在 umbrella** | `micro-service/review-c/`（被 umbrella `.gitignore` 忽略） |
| **Framework** | Kratos v2 + Wire + Consul resolver + gRPC client |
| **HTTP port** | 8002（默认） |
| **下游 RPC** | `review-service` v1 `ReviewClient`（经 Consul discovery） |

## 🧭 定位与职责（从 L0 `02_ARCHITECTURE_OVERVIEW.md` SYS-02 + `04_SYSTEM_DESIGN/review-c.md` 摘录）

- **薄 BFF**：暴露 4 条 C 端 HTTP 路由 —— `CreateReview` / `GetReview` / `ListReviewByUserID` / `ListReviews`
- 读路径 `caller_role=CALLER_ROLE_C`；`ListReviewByUserID` 隐式等价 role=C（proto 无 `caller_role` 字段）
- 写路径强制校验 `path/body.user_id == header.X-User-Id`（轻量越权防护），可通过配置 `security.enforce_user_id_consistency` 热关
- **不做**二次缓存、**不做**媒体上传
- 匿名脱敏由 review-service 的 `filterByRole` 保证（C 端看到 `user_id=0`），本层**不**再处理
- `ListReviews` 对 C 端只开放白名单筛选字段：`store_id / spu_id / sku_id / score / has_media / sort / page / page_size`；**禁止** `user_id / status / has_reply`（出现返回 400）
- 依赖 `../reviewapis/` 的 `consumer/v1`（HTTP 契约）与 `review/v1`（下游 gRPC 契约）；**只读**

## 📍 当前 Wave / 进度状态

- **Wave 0（进场与骨架验收）**：Kratos 骨架存在 `review-c/` 下（由业务代码初始化），但 VibeCoding 视角下所有任务仍在 `[ ]` 待开工状态
- **Wave 1 Foundation (S-BFF-C)**：⏳ 0/4
- **Wave 2 RPC Wiring (S-BFF-C)**：⏳ 0/7
- **Wave 3 Hardening (S-BFF-C)**：⏳ 0/2 + INT-S-BFF-C

**剩余总数**：15 项（整套 System 3 全量）

> **与 review-o 的对照**：本 scope 的 Phase 1/2/3 结构与 review-o 高度同构，首个任务建议直接把 review-o 同类实现作为模板（例如 `internal/client/` 布局、header 中间件骨架）；但**禁止共用代码**——v1 为了降低耦合风险，两 BFF 独立实现，只共用 pattern。

详见 [05_TASKS.md](./05_TASKS.md)。

## 📚 必读文档加载顺序（Bootstrap 时）

1. **本 scope（可写）**
   - [AGENTS.service.md](./AGENTS.service.md)（本文件）
   - [05_TASKS.md](./05_TASKS.md)
   - [artifacts/error_journal.md](./artifacts/error_journal.md)（若存在）
2. **Umbrella 强制只读**
   - `../../.vibe/genesis/v1/02_ARCHITECTURE_OVERVIEW.md` §SYS-02
   - `../../.vibe/genesis/v1/04_SYSTEM_DESIGN/review-c.md`（L0）
   - `../../.vibe/genesis/v1/08_CODING_STANDARDS.md`
   - `../../.vibe/genesis/v1/07_ARCHITECTURE_CHEATSHEET.md`
   - `../../.vibe/artifacts/error_journal.md`（重点关注 BFF / gRPC client / 越权校验 / 匿名脱敏相关规则）
3. **Umbrella 强烈推荐**
   - `../../.vibe/genesis/v1/04_SYSTEM_DESIGN/_research/review-c-research.md`
   - `../../.vibe/genesis/v1/04_SYSTEM_DESIGN/review-service.md` — **必读**：本 scope 是 review-service 的消费方，需理解下游 RPC 语义（尤其 `filterByRole` 与 `normalizeFilter` 对 C 端的行为）
   - `../../.vibe/genesis/v1/01_PRD.md` — REQ-001..003 是 C 端核心需求
   - `../../.vibe/genesis/v1/03_ADR/`
   - `../../.vibe/genesis/v1/05_TASKS_bff_first.md` §System 3 — 本 scope 全部任务的原始 acceptance criteria
   - `../../.vibe/coord/dependencies.yaml`
   - `../review-o/.vibe/` — **同构参考**（只读）

## 🔗 跨 scope 依赖（只读参考）

- **被依赖**：（C 端 App / H5 / Mini Program 前端 / API Gateway，超出本工作区）
- **依赖于**：
  - `review-service`（gRPC，Consul discovery） — 下游核心服务
  - `../reviewapis/`（契约，只读）
- **基础设施**：Consul（服务发现）；无本地存储

## 🌳 Git 规范

- 分支：`feat/<task-id>-<slug>`，`<task-id>` 使用本 scope 05_TASKS.md 中的编号（如 `T3.1.1`、`INT-S-BFF-C`）
- 每个 Level-3 Task 独立 commit
- 完工后 PR 到**本 scope 的** `main` 分支；umbrella 不做任何 commit

## ⚠️ 特别提醒

- **越权防护是 C 端的命门**：T3.2.2 `user-id-guard` middleware 必须严格校验；T3.2.7 单测要有 header-spoof 反向用例
- **白名单筛选**：T3.2.6 `ListReviews` 对 `user_id / status / has_reply` 一律 400；不要静默丢弃
- **幂等**：T3.2.3 `CreateReview` 同 `order_id` 重复提交应返回原 review（由下游 review-service 的 `FindByOrderID` 保证）
- **健康检查浅探**：`/health` 不探下游
- **日志脱敏**：绝不记录评价 content 原文；匿名评价 user_id 在日志中也应保持 0
