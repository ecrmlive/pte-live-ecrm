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
| [docs/release/SERVICE-MATRIX.md](docs/release/SERVICE-MATRIX.md) | 各端命名 / Docker 网络与固定 IP |
| [docs/release/COMMANDS.md](docs/release/COMMANDS.md) | `make local-*` / `deploy-*` |

表前缀：`eb_*`（CRMEB）→ **`qixi_*`（本仓库）**。

## 仓库骨架（阶段 0）

```text
api/                 # Go：cmd/api-admin、cmd/api-app、cmd/job（进程分立）
admin/               # 平台后台（Vben）
merchant-admin/      # 商户后台
app-uni/             # 用户端 uni-app x
sql/
release/qixi-mergers-*/
```

本机：

```bash
make init-env
make local-db && make local-mq
make local-api-admin   # 后台 http://127.0.0.1:18080/healthz
make local-api-app     # C 端 http://127.0.0.1:18085/healthz
make local-admin       # pack dist；宿主机 Nginx 见 release/opts/nginx
```

## Agent

- [AGENTS.md](AGENTS.md) — 协作与工程约定（含个人工作习惯）
- Cursor Skill：`.cursor/skills/qixi-live-mergers/`
- Codex Skill 镜像：`codex-skills/qixi-live-mergers/`

## 外部参考（不入库）

```text
~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0
```

## 当前状态

功能基线已锁定；阶段 0 基建与骨架已落地（release 网络/IP、Go 健康检查、前端占位、Make/脚本）。业务竖切从阶段 1 身份三端开始。
