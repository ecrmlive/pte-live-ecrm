# 接口文档正确性说明（必读）

## 一句话

**不能假装 100% 正确。** 已按源码校验分级；`high` 可对照开发，`stale`/`unresolved` 禁止照抄实现。

### 重大发现（下单）

| 路由 | 路由声明方法 | 源码是否存在 | 文档结论 |
| --- | --- | --- | --- |
| `POST /api/v2/order/create` | v2CreateOrder | ✅ 存在 | **high** — 以此为准 |
| `POST /api/v2/order/check` | v2CheckOrder | ✅ 存在 | **high** — 以此为准 |
| `POST /api/order/create` | createOrder | ❌ 不存在 | **stale** |
| `POST /api/order/check` | checkOrder | ❌ 不存在 | **stale** |
| `POST /api/v3/order/create` | v3CreateOrder | ❌ 不存在 | **stale** |
| `POST /api/v3/order/check` | v3CheckOrder | ❌ 不存在 | **stale** |

当前这份 CRMEB 源码里，用户下单实现落在 **v2**；v1/v3 路由是死路由。重建时不要被错误文档带去实现不存在的方法。

## 结论（诚实）

| 等级 | 含义 | 数量（约） | 开发能否直接用 |
| --- | --- | --- | --- |
| **high** | 路径正确且控制器方法真实存在；参数从方法体抽取 | 1999 | 可作主对照，实现前仍打开源码看 data |
| **stale** | 路由有、方法无（原项目死路由） | 26 | **禁止**按文档实现 |
| **unresolved** | 未可靠映射 | 2 | 先查 route 源码 |
| 其它 | 待定 | 0 | 先核实 |

**已人工/脚本强校验通过的关键项：**

- 路径端前缀：`sys/` `mer/` `api/` `openapi/` `ser/`（全量检查 0 错误）
- 商户开放接口 6 条：路径、控制器、鉴权参数、与源码一致
- 关键用户接口：`POST /api/auth/login`、`POST /api/v2/order/create|check`、`GET /api/user` 为 **high**；`/api/order/create|check` 与 `/api/v3/order/create|check` 为 **stale**（控制器无对应方法）
- 表前缀参考 DDL：165 张均为 `qixi_`，无残留 `eb_` 表名

**不能保证 100% 的部分（开发必须知道）：**

1. **`data` 响应体字段**：多数接口动态组装，OpenAPI 用统一 `CrmebResponse`；详细字段需看 Repository 返回或抓包。
2. **请求参数**：只从控制器可见的 `params`/`param`/`getPage`/常量抽取；`form-builder` 动态表单、中间件注入参数可能不全。
3. **stale 路由**：原项目就有「路由有、方法无」的情况，文档已标注，**不是我们编造的接口**。
4. **Header 鉴权名**：以原项目前端/中间件为准，YAML 中 securitySchemes 为示意。

## 开发铁律

1. 只实现 `doc_confidence=high` 的接口，或你亲自在源码核对过的接口。  
2. 实现前打开 `controller_file` 与 `route` 再对一次路径和参数。  
3. 本仓库重建 API 的正式契约以**将来的 qixi OpenAPI**为准；本目录是 CRMEB **对照文档**。  
4. 表结构以 `docs/schema/` 的 `qixi_` 为准。

## 校验产物

- `VALIDATION-REPORT.md` — 自动校验明细  
- `crmeb-api-all.jsonl` 字段 `doc_confidence` / `doc_note`  
- OpenAPI YAML operation.description 含置信度说明  

源码根：`~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0`
