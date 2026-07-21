# CRMEB 原项目接口约定（通用）

> 从源码 `crmeb/services/ApiResponseService.php` 与路由前缀配置提取。  
> 本文档描述**原项目**接口，不是本仓库重建后的 API。

## 基础 URL 前缀

| 端 | 路由前缀 | 配置项 |
| --- | --- | --- |
| 平台后台 | `/sys/` | `admin.api_admin_prefix` |
| 商户后台 | `/mer/` | `admin.api_merchant_prefix` |
| 用户端 | `/api/` | `route/api.php` 分组 |
| 客服端 | `/ser/` | `admin.api_service_prefix` |
| 商户开放接口 | `/openapi/` | `admin.api_openapi_prefix` |

## 统一响应结构

```json
{
  "status": 200,
  "message": "success",
  "data": {}
}
```

| 字段 | 说明 |
| --- | --- |
| status | 业务状态码，成功多为 `200`，失败多为 `400`（另有 `status()` 中间态） |
| message | 提示文案 |
| data | 业务数据；失败时可能无此字段 |

成功调用常见写法：`return app('json')->success($data);` 或 `success('消息', $data)`。

## 鉴权（概要）

| 端 | 方式（源码） |
| --- | --- |
| 平台后台 | AdminToken + AdminAuth 中间件 |
| 商户后台 | MerchantToken + MerchantAuth |
| 用户端 | UserToken（部分接口可游客） |
| 客服端 | ServiceToken |
| 开放接口 | OpenApiAuth（access_key/signature/expiration/unique）+ JWT token |

## 文档生成说明

- 路径：由 `route/**/*.php` 解析，并拼上端前缀。
- 请求参数：从控制器方法内 `params` / `param` / `getPage` / Validate 规则 / 注释提取。
- 返回参数：多数接口 `data` 为动态结构；文档给出可确定字段，其余标注「见 Repository 返回」。
- 若某接口未解析到控制器，仍保留路由行，参数可能只有路径参数。

分端明细见同目录其它文件；机器可读总表：`crmeb-api-all.jsonl` / `crmeb-api-all.tsv`。
