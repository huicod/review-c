# review-c — .vibe/artifacts/

本目录存放本 scope 开发过程中产生的**临时产物**与**可沉淀经验**。

## 目录结构

```
artifacts/
├── README.md            # 本文件（版本追踪）
├── error_journal.md     # 本 scope 专属 Prevention Rules（append-only，版本追踪）
├── logs/                # 构建/测试/压测日志（.gitignore 忽略）
├── plan_*.md            # /plan 产出的临时计划（.gitignore 忽略）
└── prp_*.md             # /propose 产出的临时提案（.gitignore 忽略）
```

## Git 追踪规则

由本 scope 的 `.gitignore` 决定：

- **追踪**：`README.md`、`error_journal.md`
- **忽略**：`logs/`、`plan_*.md`、`prp_*.md`

## 与 umbrella 的关系

- 本 scope 的 `error_journal.md` 是**一级沉淀**
- 跨 scope 普适规则由 Orchestrator 合并到 umbrella `../../.vibe/artifacts/error_journal.md`

## INT-S-BFF-C 集成测试报告

INT-S-BFF-C 的验证输出固定写入 `logs/int-s-bff-c.md`（被 git 忽略，本地留档）。
