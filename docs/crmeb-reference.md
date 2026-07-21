# CRMEB 参考源码分析

> 路径：`~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0`  
> 版本：`CRMEB-MER-v4.0`（`.version`）  
> 用途：功能与领域对照。本仓库不运行、不提交该 PHP 源码。

## 1. 技术栈（参考系统）

| 项 | 实际 |
| --- | --- |
| 语言 | PHP >= 8.0 |
| 框架 | ThinkPHP 8 + ThinkORM |
| 运行时 | think-swoole 4.1（Swoole 常驻） |
| 队列 | think-queue |
| 微信 | EasyWeChat 6 |
| 其他 | JWT、支付 SDK、OSS/七牛/OBS、表单构建、Spreadsheet |

README 仍写 TP6/PHP7.x，以 `composer.json` 与根目录 `AGENTS.md` 为准。

## 2. 目录与分层

```text
route → controller/{admin,merchant,api,openapi,service}
      → validate
      → repositories  （主业务编排）
      → dao / model
      → crmeb/services | jobs | listens
      → app('json')->success/fail
```

| 目录 | 职责 |
| --- | --- |
| `app/controller/admin` | 平台后台 |
| `app/controller/merchant` | 商户后台 |
| `app/controller/api` | C 端 API |
| `app/controller/openapi` | OpenAPI |
| `app/controller/service` | 服务端入口（客服等） |
| `app/common/repositories` | 业务核心 |
| `crmeb/services` | 支付、微信、短信、上传等 |
| `crmeb/jobs` / `listens` | 队列与事件/定时器 |
| `public/system` `public/mer` | 平台/商户后台静态资源 |
| `extend/mp-weixin` | 微信小程序产物参考 |
| `install/crmeb_merchant.sql` | 初始化库（约 165 张 `eb_*` 表） |

## 3. 路由入口

| 文件 | 端 |
| --- | --- |
| `route/admin/*` | 平台：merchant、product、order、marketing、user、finance 等 |
| `route/merchant/*` | 商户：product、order、coupon、marketing、staffs 等 |
| `route/api/*` | 用户：login、业务接口 |
| `route/openapi/*` | 开放接口 |
| `route/service.php` | 服务相关 |

## 4. 核心数据表分组（`eb_` 前缀）

| 分组 | 代表表 |
| --- | --- |
| 商户 | `eb_merchant`, `eb_merchant_admin`, `eb_merchant_intention`, `eb_merchant_category`, `eb_merchant_type`, `eb_merchant_applyments` |
| 商品 | `eb_store_product`, `eb_store_product_attr_value`, `eb_store_spu`, `eb_store_category`, `eb_store_brand` |
| 活动商品 | `eb_store_seckill_*`, `eb_store_product_group*`, `eb_store_product_presell*`, `eb_store_product_assist*` |
| 购物车订单 | `eb_store_cart`, `eb_store_group_order`, `eb_store_order`, `eb_store_order_product`, `eb_store_order_status` |
| 售后 | `eb_store_refund_order`, `eb_store_refund_product`, `eb_store_refund_status` |
| 优惠券 | `eb_store_coupon*` |
| 用户资产 | `eb_user`, `eb_user_bill`, `eb_user_brokerage`, `eb_user_extract`, `eb_user_recharge`, `eb_user_address` |
| 财务结算 | `eb_financial`, `eb_financial_record`, `eb_store_order_profitsharing` |
| 配送 | `eb_delivery_*`, `eb_shipping_template*` |
| 客服 | `eb_store_service*` |
| 直播 | `eb_broadcast_*` |
| 社区/圈子 | `eb_community*`, `eb_circle*` |
| 系统 | `eb_system_admin`, `eb_system_role`, `eb_system_menu`, `eb_system_config*`, `eb_diy` |

完整表名列表可从 `install/crmeb_merchant.sql` 的 `CREATE TABLE` 提取。

## 5. 高风险链路（改写时必须先对齐）

参考源码中改动需特别谨慎的文件（见其 `AGENTS.md`）：

- `StoreOrderCreateRepository` / `StoreOrderRepository` / `StoreRefundOrderRepository`
- `crmeb/listens/pay/`、`crmeb/services/PayService.php`、微信服务
- 库存、活动商品、优惠券、积分、佣金、商户结算相关 Repository

本仓库重建时对应能力必须具备：

1. 幂等下单与支付回调
2. 库存扣减与回滚
3. 多商户拆单与分账字段
4. 退款与结算状态机互斥
5. 优惠/积分/佣金可追溯账单

## 6. 原项目已有 Agent 资料

| 路径 | 用途 |
| --- | --- |
| `CRMEB_MER_v4.0/AGENTS.md` | PHP 项目协作约定 |
| `codex-skills/crmeb-merchant-extension-guide/` | PHP 二开指南（architecture / extension / risk / troubleshooting） |

本仓库 Skill **不替代**上述 PHP 二开文档；仅在需要对照原实现时读取外部路径。

## 7. 前端参考

| 端 | 参考位置 |
| --- | --- |
| 平台后台 | `public/system` + `system.html` |
| 商户后台 | `public/mer` + `mer.html` |
| 客服端 | `public/ser` |
| 数据大屏 | `public/screen` |
| 小程序页面树 | `extend/mp-weixin/v4.0/pages/*`（首页、商品、订单、店铺、分销、社区、配送、员工等） |

目标管理后台使用 **Vben 5+**，C 端优先 **uni-app x**，信息架构可对照上述页面树与脑图。
