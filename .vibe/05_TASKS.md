# review-c — 任务清单（Scope 真相源）

> **迁移说明**：本文件自 umbrella `.vibe/genesis/v1/05_TASKS_bff_first.md` §System 3 抽取。
>
> - **以代码为准（2026-04-23 起）**：实现锚点见 §3；**完整**任务描述/验收长文仍在 umbrella 同 TaskID。
> - 只在本文件勾进度；**不要**手改 umbrella（Orchestrator 与 scope 双向往返时例外）。

## 1. 状态概览

- ⏳ Wave 1 Foundation (S-BFF-C): **3.5/4** — T3.1.1~T3.1.3 已落地；T3.1.4 **部分完成**：已有 `Makefile`、`Dockerfile`；**`.golangci.yml` 未入库**（可与 `review-service` 的 `T1.4.7` 根目录统一拉齐，同 `review-o` 策略）。
- ⏳ Wave 2 RPC Wiring (S-BFF-C): **6/7** — T3.2.1~T3.2.6 已合入；**T3.2.7**（`internal/service` 单测 + 覆盖率）未做（仓库内仍无 `*_test.go`）。
- ⏳ Wave 3 Hardening: **0/2** + `INT-S-BFF-C` 未做（无 `metrics`/`health` 专用文件；与 `review-o` 同型缺口）。

**与 WBS 字面差异（以目录为准）**

| 原 WBS 表述 | 代码实际 |
|------------|----------|
| `internal/client/review.go` | 下游 client：`internal/data/data.go` 中 `NewReviewServiceClient` + `wire` |
| `internal/service/review.go` | `internal/service/consumer.go`（`ConsumerService` 实现 `consumer/v1` 四方法） |
| `internal/service/list_query_c.go` | 未单独拆分；`ListReviewsParam` + `data/consumer.go` 对 `ListReviews` 的字段映射 + `data` 层对 `REVIEW_STATUS_APPROVED` 的显式注入 |
| 仓库根 `.golangci.yml` | **不存在**（同刻 `review-o` 亦无，属跨服务 T1.4.7 范围） |

**剩余任务（粗算 5 项）**：T3.1.4 中 golangci 收尾、T3.2.7、T3.3.1、T3.3.2、`INT-S-BFF-C`。

## 2. Active Backlog

### 2.1 Wave 1 收尾

- [~] **T3.1.4** [基础]: Makefile + Dockerfile + `.golangci.yml`
  - **已满足**：`Makefile`、`Dockerfile`；`go build` / `docker build` 可做为发布基线前再跑一遍体积约束。
  - **未满足**：`.golangci.yml` + 与 umbrella 一致的 `make lint` 全绿；完成后改 `[x]`。

### 2.2 Wave 2 收尾

- [ ] **T3.2.7** [REQ-001..003]: `internal/service` table-driven 单测（mock 下游 gRPC），覆盖率 ≥ 70%  
  - **完整规格**：`../../.vibe/genesis/v1/05_TASKS_bff_first.md` System 3 · T3.2.7。
  - **当前缺口**：`internal/service/` 下**无** `*_test.go`。

### 2.3 Wave 3 Hardening

- [ ] **T3.3.1** Prometheus `/metrics` + 浅 `/health`（Kratos HTTP `:8002`；依赖 T3.2.7 的代码稳定面）
- [ ] **T3.3.2** 生产型 Dockerfile（多阶段）+ `deploy/k8s` + 根 `docker-compose` 与 **T3.3.1** 的 `/health` 挂钩

### 2.4 Milestone

- [ ] **INT-S-BFF-C**：端到端脚本 + `.vibe/artifacts/logs/int-s-bff-c.md`（**依赖** T3.3.2；完整条目见 umbrella）

## 3. Completed Log

### Wave 1 — Foundation (S-BFF-C)

- [x] **T3.1.1** 骨架：`cmd/review-c`、`internal/{conf,server,service,biz,data}`、`go.mod`、`third_party/reviewapis` 挂载
- [x] **T3.1.2** `internal/conf/conf.proto` + `configs/config.yaml`（`make config` 生成 `conf.pb.go`）
- [x] **T3.1.3** 共享契约：`github.com/huicod/reviewapis` + `replace` 至工作区；对外 **`consumer/v1`**

### Wave 2 — RPC Wiring (S-BFF-C)

- [x] **T3.2.1** 下游 gRPC + Consul discovery：`internal/data/data.go`（`NewReviewServiceClient`、recovery/tracing 中间件、`wire`）
- [x] **T3.2.2** `internal/server/middleware.go`：`headerExtract` + `userIDGuard` + `conf.Security` 注入
- [x] **T3.2.3~T3.2.6** `internal/service/consumer.go` + `internal/biz/consumer.go` + `internal/data/consumer.go`：`CreateReview` / `GetReview` / `ListReviewByUserID` / `ListReviews` 全路径打通至 `rv1.ReviewClient`，读路径统一 `CallerRole_C`

## 4. 源文档溯源

- `../../.vibe/genesis/v1/05_TASKS_bff_first.md` §System 3 — 本 scope 全部任务原文
- `../../.vibe/genesis/v1/04_SYSTEM_DESIGN/review-c.md` — L0 设计
- `../../.vibe/genesis/v1/04_SYSTEM_DESIGN/review-service.md` — 下游 RPC 语义（必读，尤其 `filterByRole` / `normalizeFilter` 对 C 端的行为）
- `../../.vibe/genesis/v1/04_SYSTEM_DESIGN/_research/review-c-research.md` — 前期调研
- `../../.vibe/genesis/v1/02_ARCHITECTURE_OVERVIEW.md` §SYS-02
- `../../.vibe/artifacts/error_journal.md` — Prevention Rules
- `../review-o/.vibe/` — 同构参考（只读）

## 5. 状态映射约定

- `[ ]` — 待开工
- `[x]` — 已完成 & 已 commit
- `[~]` — 部分完成 / 延后（注明原因与未完部分）
- `[!]` — 在飞但被阻塞（注明阻塞点）
