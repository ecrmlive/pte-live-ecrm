# pte-live-ecrm

多商户商城管理系统。

功能基线对齐 **CRMEB Merchant v4.0**；技术实现按本仓库目标栈重建（Go Gin/GORM、Vben 5+、uni-app x），不以 PHP/Swoole 为运行时。

## 文档

从这里开始：[docs/README.md](docs/README.md)

| 文档 | 说明 |
| --- | --- |
| [docs/product-understanding.md](docs/product-understanding.md) | 我们要做什么（主链路） |
| [docs/features/README.md](docs/features/README.md) | 各端→按钮→CRUD |
| [docs/SYSTEM-ARCHITECTURE.md](docs/SYSTEM-ARCHITECTURE.md) | 系统、应用、JWT 与三库边界唯一口径 |
| [docs/schema/README.md](docs/schema/README.md) | 三库表设计：`qixi_crm_a_`、`qixi_crm_b_`、`qixi_crm_m_` |
| [docs/overview.md](docs/overview.md) | 项目定位与三角色 |
| [docs/architecture-target.md](docs/architecture-target.md) | 目标技术架构 |
| [docs/release/SERVICE-MATRIX.md](docs/release/SERVICE-MATRIX.md) | 各端命名 / Docker 网络与固定 IP |
| [docs/release/COMMANDS.md](docs/release/COMMANDS.md) | `make local-*` / `make test-*` |

表前缀：后台 `qixi_crm_a_*`、C 端业务 `qixi_crm_b_*`、店铺 `qixi_crm_m_*`。`pte-live-im` 的表与数据库规则仅由其自身仓库维护。

## 仓库骨架（阶段 0）

```text
admin-platform/      # 平台/商户/区域/客服/运营统一后台（Vben 5.7+）
admin-merchant/      # 独立店铺管理系统（Vben 5.7+）
app-pc/              # 用户端 PC Web 商城
app-uni/             # 用户端 uni-app x：H5 与小程序
app-ios/             # 用户端 iOS
app-adnroid/         # 用户端 Android（目录名按项目约定）
app-harmony/         # 用户端鸿蒙
api-platform/        # 统一后台 API（独立 Go module）
api-business/        # C 端业务 API（独立 Go module）
api-merchant/        # 店铺 API（独立 Go module）
job/                 # 异步任务（独立 Go module）
sql/
release/pte-live-ecrm-*/
```

本机：

```bash
make init-env-local
make pack-backend
make local-infra        # 校验 pte-live-im 共享基础设施已运行
make local-db-init      # 在 pte_live_mysql 中初始化七禧三库
make local-backend      # api-platform :18081、api-business :18082、api-merchant :18083

# 测试宿主机执行同一组目标；容器名、共享网络、数据库和 YAML 与 local 完全一致
make init-env-test
make test-infra
make test-db-init
make test-backend
```

## Agent

- [AGENTS.md](AGENTS.md) — 协作与工程约定（含个人工作习惯）
- Cursor Skill：`.cursor/skills/pte-live-ecrm/`
- Codex Skill 镜像：`codex-skills/pte-live-ecrm/`

## 外部参考（不入库）

```text
~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0
```

## 当前状态

功能基线已锁定，但尚未达到全端 100% CRMEB 对齐或生产交付条件。实际完成度以 `docs/CRMEB-FULL-FUNCTION-CHECKLIST.md`、各端可运行接口和测试证据为准，不以页面占位或菜单数量宣称完成。
