# CRMEB 原项目接口文档（本仓库文档）

从原项目源码解析，供 **qixi-live-mergers** 重建对照。

## 正确性（先读）

**[ACCURACY.md](./ACCURACY.md)** — 文档分级与下单死路由说明。文档错了开发就会错。

## OpenAPI 3.0（正式 YAML）

| 文件 | 说明 |
| --- | --- |
| **[openapi.yaml](./openapi.yaml)** | **总览（推荐）**，约 2000 paths |
| [openapi-platform-sys.yaml](./openapi-platform-sys.yaml) | 平台 `/sys/` |
| [openapi-merchant-mer.yaml](./openapi-merchant-mer.yaml) | 商户后台 `/mer/` |
| [openapi-user-api.yaml](./openapi-user-api.yaml) | 用户端 `/api/` |
| [openapi-openapi.yaml](./openapi-openapi.yaml) | 商户开放接口 `/openapi/` |
| [openapi-service-ser.yaml](./openapi-service-ser.yaml) | 客服 `/ser/` |

可导入 Swagger UI / Apifox / Postman（OpenAPI 3.0）。

统一响应组件：`components.schemas.CrmebResponse` → `{ status, message, data }`。

## Markdown 可读版

| 文件 | 说明 |
| --- | --- |
| [00-common.md](./00-common.md) | 前缀 / 响应 / 鉴权 |
| [01-platform-sys.md](./01-platform-sys.md) | 平台 972 |
| [02-merchant-mer.md](./02-merchant-mer.md) | 商户后台 595 |
| [03-user-api.md](./03-user-api.md) | 用户端 436 |
| [04-openapi.md](./04-openapi.md) | 开放接口 6（已校对） |
| [05-service-ser.md](./05-service-ser.md) | 客服 18 |

## 机器源数据

| 文件 | 说明 |
| --- | --- |
| `crmeb-api-all.jsonl` | 生成 YAML 的源 |
| `crmeb-api-all.tsv` | 表格 |

## 统计

| 项 | 数量 |
| --- | ---: |
| 接口合计 | 2027 |
| OpenAPI paths（合并去重后） | ~2000 |

## 注意

- 描述的是 **CRMEB 原项目**接口，不是本仓库已实现的 Go API。
- 重建后的本仓库 API 应另维护 `docs/openapi/`（前缀与路径按 qixi 设计），可对照本目录迁移。
- `data` 内部字段多为动态结构；YAML 中用统一 `CrmebResponse`，细节见各 operation 的 description。
- OpenAPI 里 `deprecated: true` = `stale` 死路由，**禁止实现**。
- 用户下单对照：`POST /api/v2/order/create` / `check`（不要用无方法的 v1/v3）。

源码：`~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0`
