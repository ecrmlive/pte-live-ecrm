# qixi-live-mergers

多商户商城管理系统。

功能基线对齐 **CRMEB Merchant v4.0**；技术实现按本仓库目标栈重建（Go Gin/GORM、Vben 5+、uni-app x），不以 PHP/Swoole 为运行时。

## 文档

从这里开始：[docs/README.md](docs/README.md)

| 文档 | 说明 |
| --- | --- |
| [docs/product-understanding.md](docs/product-understanding.md) | 我们要做什么（主链路） |
| [docs/features/README.md](docs/features/README.md) | 各端→按钮→CRUD |
| [docs/schema/README.md](docs/schema/README.md) | 表设计，前缀 `qixi_` |
| [docs/overview.md](docs/overview.md) | 项目定位与三角色 |
| [docs/architecture-target.md](docs/architecture-target.md) | 目标技术架构 |

表前缀：`eb_*`（CRMEB）→ **`qixi_*`（本仓库）**。

## Agent

- [AGENTS.md](AGENTS.md) — 协作与工程约定（含个人工作习惯）
- Cursor Skill：`.cursor/skills/qixi-live-mergers/`
- Codex Skill 镜像：`codex-skills/qixi-live-mergers/`

## 外部参考（不入库）

```text
~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0
```

## 当前状态

已完成需求分析文档与 Skill。业务源码（API / 管理后台 / 用户端 / release）待后续任务落地。
