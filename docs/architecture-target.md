# 目标技术架构

本仓库重建栈，不以 CRMEB PHP/Swoole 为运行时。

**服务命名 / Docker 网络 / IP / 分阶段计划**（权威）：[`docs/release/SERVICE-MATRIX.md`](./release/SERVICE-MATRIX.md)。

## 1. 固定版本

| 组件 | 版本 |
| --- | --- |
| Go | 1.26.5 |
| Node.js | 24.18.0 |
| npm | 11.16.0 |
| Corepack | 0.35.0 |
| alpine | 3.24.1 |
| nginx | 宿主机安装（本仓库不用 Docker Nginx） |
| node 镜像 | 24.18.0-alpine3.24 |
| MySQL | 8.4.10 |
| Redis | 8.8.0 |
| etcd | v3.7.0 |
| NATS | 2.12.0-alpine |
| 对象存储 | 腾讯云 COS（非本仓容器） |

配置文件统一使用 `.yaml`。

## 2. 仓库形态

```text
qixi-live-mergers/
├── AGENTS.md
├── docs/
├── api/                    # Go 同仓；进程分立
│   ├── cmd/api-admin/      # 后台：platform/merchant/manager/service/open
│   ├── cmd/api-app/        # C 端：app + callback
│   ├── cmd/job/            # NATS 消费
│   ├── internal/domain/    # 共享领域
│   ├── internal/{platform,merchant,app,...}/
│   └── conf/{admin,app,job}.yaml
├── admin-platform/         # Vben 5 平台源码 → pack key `admin` → qixi-mergers-admin
├── admin-merchant/         # Vben 5 商户源码 → pack key `merchant-admin` → qixi-mergers-merchant-admin
├── app-uni/                # uni-app x → qixi-mergers-h5
├── app-pc/                 # Vue3 PC 商城 → qixi-mergers-pc
├── sql/
├── release/                # 见 SERVICE-MATRIX
└── scripts/
```

**API 形态**：**后台 API 与 C 端 API 独立进程**（`api-admin` / `api-app`）+ `job`；领域代码共享，不按业务再拆微服务。

## 3. API 技术选型

| 层 | 选型 |
| --- | --- |
| HTTP | Gin |
| ORM | GORM |
| 文档 | Swagger / OpenAPI（本仓库契约在 `docs/openapi/`，CRMEB `docs/api/` 仅对照） |
| DB | MySQL 8.4 |
| 缓存/锁 | Redis |
| 配置/服务发现 | etcd |
| 消息 | NATS |
| 认证 | JWT；平台与商户 RBAC 隔离 |

## 4. 客户端规范

| 端 | 规范 |
| --- | --- |
| 管理后台 | Vben 5+ |
| 用户端 H5/小程序 | uni-app x（UTS / HBuilderX 5.0+） |
| 用户端 PC | Vue 3 + Vite + TypeScript（`app-pc/`，功能表 4） |
| 导航栏（新 App） | 高度 44，紧贴状态栏底部，按钮 44px |

## 5. Docker 与部署

- 网络：本仓 `qixi_mergers_net` / `172.30.80.0/24`（业务）；共享基建挂 `pte_live_net`（MySQL/Redis/NATS/etcd 由 **pte-live-im** 启动）。对象存储用腾讯云 COS。
- 业务容器：`api_admin .20`、`job .21`、`api_app .22`（均另挂 `pte_live_net`）。
- 本机构建产物 → `release/<svc>/` → rsync → 远程 Compose **挂载**运行（后端）。
- 前端仅 `dist/`；**Nginx 统一宿主机**；后台反代 api-admin，C 端反代 api-app。
- `config/local` · `config/prod`；禁止服务器源码构建。
- 必读：[`docs/release/COMMANDS.md`](./release/COMMANDS.md)、[`PACK-AND-CONFIG.md`](./release/PACK-AND-CONFIG.md)。
- 部署须用户授权；完成后反馈服务器 IP。

## 6. 领域横切要求

- 多租户：商户域带 `merchant_id`，商户接口强制隔离。
- 订单：平台主单 + 商户子单。
- 资金/退款/结算：状态机见 `docs/api/FUNCTIONAL-TRUTH.md`。
- 营销：价格计算单一入口。
- 中文测试数据：utf8mb4。

## 7. 数据库前缀

业务表前缀 `qixi_`。详见 `docs/schema/`。
