# 服务命名 · Docker 网络 · IP 矩阵（权威）

项目：`qixi-live-mergers`。与 `unified-docker-release` 对齐。  
全端开发计划：[`docs/dev-plan-full.md`](../dev-plan-full.md)。  
打包命令：[COMMANDS.md](./COMMANDS.md)、[PACK-AND-CONFIG.md](./PACK-AND-CONFIG.md)。

## 1. Token

| 用途 | 值 |
| --- | --- |
| release 目录前缀 | `qixi-mergers-` |
| Compose project / 容器名前缀 | `qixi_mergers` |
| Docker 网络名 | `qixi_mergers_net` |
| 网段 | `172.30.80.0/24` |

网络由 `release/qixi-mergers-db` 首次创建；其余 Docker 服务 `external: true`。

## 2. Nginx 约定

- **统一宿主机 Nginx**（`release/opts/nginx/`）；不用 Docker Nginx。
- 后台前端 `/api` → **api-admin**；C 端 H5 `/api` → **api-app**。

## 3. 基建

| release 目录 | compose project | 容器名 | 固定 IP | 宿主端口 (local) |
| --- | --- | --- | ---: | --- |
| `qixi-mergers-db` | `qixi_mergers_db` | `qixi_mergers_mysql` | `.10` | `13306` |
| 同上 | 同上 | `qixi_mergers_redis` | `.11` | `16379` |
| 同上 | 同上 | `qixi_mergers_etcd` | `.12` | `12379` |
| 同上 | 同上 | `qixi_mergers_minio` | `.14` | `19000` / `19001` |
| `qixi-mergers-mq` | `qixi_mergers_mq` | `qixi_mergers_nats` | `.13` | `14222` |

## 4. 业务后端（API 已拆分）

| 职责 | release 目录 | 容器名 | 固定 IP | 产物 | 宿主端口 |
| --- | --- | --- | ---: | --- | --- |
| **后台 API** | `qixi-mergers-api-admin` | `qixi_mergers_api_admin` | `.20` | `bin/api-admin` | `18080` |
| **C 端 API** | `qixi-mergers-api-app` | `qixi_mergers_api_app` | `.22` | `bin/api-app` | `18085` |
| 异步任务 | `qixi-mergers-job` | `qixi_mergers_job` | `.21` | `bin/job` | 无 |

compose project 均为 `qixi_mergers`。源码同仓 `api/`，进程分立：`cmd/api-admin` · `cmd/api-app` · `cmd/job`；领域代码可共享 `internal/domain`。

### 4.1 前缀归属

| 进程 | 前缀 |
| --- | --- |
| **api-admin** | `/api/platform/v1` · `/api/merchant/v1` · `/api/manager/v1` · `/api/service/v1` · `/api/open/v1` |
| **api-app** | `/api/app/v1` · `/api/callback/v1` |

配置源：`api/conf/admin.yaml` → api-admin；`api/conf/app.yaml` → api-app；`api/conf/job.yaml` → job。

## 5. 前端（无容器 · 宿主机 Nginx）

| 端 | 源码 | release（dist） | listen | 反代 API |
| --- | --- | --- | ---: | --- |
| 平台后台 | `admin/` | `qixi-mergers-admin` | `18081` | api-admin `:18080` |
| 商户后台 | `merchant-admin/` | `qixi-mergers-merchant-admin` | `18082` | api-admin |
| 用户 H5 | `app-uni/` | `qixi-mergers-h5` | `18083` | **api-app** `:18085` |
| 客服 P1 | `service-web/` | `qixi-mergers-service-web` | `18084` | api-admin |
| 用户 PC | `app-pc/` | `qixi-mergers-pc` | `18086` | **api-app** `:18085` |

## 6. 本机 URL

| 用途 | URL |
| --- | --- |
| 后台 API | `http://127.0.0.1:18080/healthz` |
| C 端 API | `http://127.0.0.1:18085/healthz` |
| 平台后台 | `http://127.0.0.1:18081` |
| 商户后台 | `http://127.0.0.1:18082` |
| 用户 H5 | `http://127.0.0.1:18083` |
| 用户 PC | `http://127.0.0.1:18086` |
| MinIO Console | `http://127.0.0.1:19001` |

## 7. Prod

- 同一容器名与网段固定 IP。
- 公网：宿主机 Nginx 分别反代 api-admin / api-app；支付回调域名指向 **api-app**。

## 8. 分阶段

见 [`docs/dev-plan-full.md`](../dev-plan-full.md)。阶段 0 已含双 API 骨架。
