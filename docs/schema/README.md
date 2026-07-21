# 数据表设计（前缀 `qixi_`）

> 来源：CRMEB MER v4.0 `install/crmeb_merchant.sql`  
> 规则：**表名前缀 `eb_` → `qixi_`**；字段名暂保持与参考库一致，便于对照。  
> 本目录为**重建参考**，落地时按 Go/GORM 规范可再规范化字段命名，但前缀固定为 `qixi_`。

## 文件

| 文件 | 说明 |
| --- | --- |
| [qixi_schema_reference.sql](./qixi_schema_reference.sql) | 全量建表参考 DDL（已改前缀） |
| [table-prefix-map.tsv](./table-prefix-map.tsv) | `eb_*` ↔ `qixi_*` 对照 |
| [columns.tsv](./columns.tsv) | 全字段清单（2121 行） |
| [domains.md](./domains.md) | 按领域分组的表说明 |
| `domain-*.md` | 各领域字段明细 |

## 统计

| 项 | 数量 |
| --- | ---: |
| 表 | 165 |
| 字段行 | 2121 |

## 前缀约定（强制）

```text
CRMEB:  eb_store_order
本仓库: qixi_store_order
```

- 配置、文档、GORM `TableName()`、SQL 迁移一律使用 `qixi_`。
- **禁止**在新代码、新迁移、新文档示例中使用 `eb_` 表前缀。
- 对照 CRMEB 源码时，完成 `eb_X` → `qixi_X` 映射。
- 不把 CRMEB 安装数据当生产数据；只吸收表结构与业务语义。
- 业务主链路见 [../product-understanding.md](../product-understanding.md)。

## 领域文件

| 文件 | 领域 |
| --- | --- |
| domain-merchant.md | 商户 |
| domain-catalog.md | 商品 |
| domain-cart_order.md | 购物车/订单 |
| domain-aftersale.md | 售后 |
| domain-coupon.md | 优惠券 |
| domain-user.md | 用户/分销 |
| domain-finance.md | 财务 |
| domain-fulfillment.md | 配送 |
| domain-marketing_product.md | 营销商品 |
| domain-cs.md | 客服 |
| domain-live.md | 直播 |
| domain-content.md | 内容/社区 |
| domain-system.md | 系统 |
