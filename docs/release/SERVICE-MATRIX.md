# 服务命名 · Docker 网络 · IP 矩阵（权威）

项目：`qixi-live-mergers`。与 `unified-docker-release` 对齐。  
全端开发计划：[`docs/dev-plan-full.md`](../dev-plan-full.md)。  
打包命令：[COMMANDS.md](./COMMANDS.md)、[PACK-AND-CONFIG.md](./PACK-AND-CONFIG.md)。

## 1. Token

| 用途 | 值 |
| --- | --- |
| release 目录前缀 | `qixi-mergers-` |
| Compose project / 容器名前缀 | `qixi_mergers` |
| 本仓 Docker 网络 | `qixi_mergers_net` / `172.30.80.0/24` |
| 共享基建网络 | `pte_live_net`（由 **pte-live-im** 创建） |

`qixi_mergers_net` 由 release 脚本 `docker network create` 按需创建（`make local-db` / 业务 up）；业务服务 `external: true`。  
API / Job 另挂 `pte_live_net`，以访问 IM 侧 MySQL / Redis / NATS / etcd。

## 2. Compose 分组

| project `name` | release 目录 | 成员 |
| --- | --- | --- |
| （无 compose） | `qixi-mergers-db` | 仅同步 `sql/`；无容器 |
| **`qixi_mergers`** | 各业务目录 | api-admin、api-app、job、admin、merchant-admin、h5、pc、service-web、manager |

### 2.1 共享基建（pte-live-im，不在本仓启动）

| 组件 | 容器名（IM） | 说明 |
| --- | --- | --- |
| MySQL | `pte_live_mysql` | 宿主 `127.0.0.1:13306`；库名 `qixi_mergers`（见 `sql/000_shared_im_mysql_bootstrap.sql`） |
| Redis | `pte_live_redis` | 宿主 `127.0.0.1:16379` |
| NATS | `pte_live_nats1..3` | JetStream 集群 |
| etcd | `pte_live_etcd1..3` | 服务发现 |

请在 `~/Documents/GitHub/pte-live-im` 启动 `db/` + `mq/`。本仓 `make local-mq` / `deploy-mq-reload` 已废弃。

### 2.2 对象存储

| 组件 | 说明 |
| --- | --- |
| 腾讯云 COS | `api/conf/*.yaml` 的 `cos:`；密钥用 `QIXI_COS_BUCKET` / `QIXI_COS_SECRET_ID` / `QIXI_COS_SECRET_KEY` 等环境变量；`cos.enabled=false` 时本地 `upload` 目录回退 |

## 3. 业务（`qixi_mergers`）

### 3.1 后端

| 职责 | release | 容器名 | IP（`qixi_mergers_net`） | 宿主端口 |
| --- | --- | --- | ---: | --- |
| 后台 API | `qixi-mergers-api-admin` | `qixi_mergers_api_admin` | `.20` | `18080` |
| Job | `qixi-mergers-job` | `qixi_mergers_job` | `.21` | 无 |
| C 端 API | `qixi-mergers-api-app` | `qixi_mergers_api_app` | `.22` | `18085` |

### 3.2 前端（nginx 容器）

| 端 | 容器名 | IP | 宿主端口 | 反代 |
| --- | --- | ---: | ---: | --- |
| 平台 | `qixi_mergers_admin` | `.30` | `18081` | api-admin |
| 商户 | `qixi_mergers_merchant_admin` | `.31` | `18082` | api-admin |
| H5 | `qixi_mergers_h5` | `.32` | `18083` | api-app |
| 客服 | `qixi_mergers_service_web` | `.33` | `18084` | api-admin |
| PC | `qixi_mergers_pc` | `.34` | `18086` | api-app |
| 店员 | `qixi_mergers_manager` | `.35` | `18087` | api-admin |

## 4. 本机 URL

| 用途 | URL |
| --- | --- |
| 后台 API | `http://127.0.0.1:18080/healthz` |
| C 端 API | `http://127.0.0.1:18085/healthz` |
| 平台 / 商户 / H5 / 客服 / PC / 店员 | `:18081` … `:18084`、`:18086`、`:18087` |
| 共享 MySQL（IM） | `127.0.0.1:13306` → 库 `qixi_mergers` |
| COS 对外域 | 见配置 `cos.base_url`（如 `https://cos.qxkejiwl.top/qixi-mergers`） |

## 5. Prod

- IM：先部署 `pte-live-im` 的 db + mq。  
- 本仓：`deploy-db-reload`（sql + 网络）→ 灌库 → 配置 COS 密钥 → `deploy-all`（业务，不含共享 db/mq）。
