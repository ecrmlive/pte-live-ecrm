# 售后 — 表字段（`qixi_`）

## `qixi_store_refund_order`

> 订单退款表

原表：`eb_store_refund_order`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `refund_order_id` | `int(10)` | NO | — | 退款单id |
| `refund_order_sn` | `varchar(32)` | NO | — | 退款单号 |
| `order_id` | `int(10)` | NO | — | 订单id |
| `uid` | `int(10)` | NO | '0' | 用户 id |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `extension_one` | `decimal(8,2)` | NO | '0.00' | 退还一级佣金 |
| `extension_two` | `decimal(8,2)` | NO | '0.00' | 退还二级佣金 |
| `integral` | `int(10)` | YES | '0' | 退还积分 |
| `delivery_type` | `varchar(32)` | YES | NULL | 快递公司 |
| `delivery_id` | `varchar(32)` | YES | NULL | 快递单号 |
| `delivery_mark` | `varchar(200)` | YES | NULL | 快递备注 |
| `delivery_pics` | `varchar(255)` | YES | NULL | 快递凭证 |
| `delivery_phone` | `varchar(18)` | YES | NULL | 联系电话 |
| `mer_delivery_user` | `varchar(32)` | YES | NULL | 收货人 |
| `mer_delivery_address` | `varchar(32)` | YES | NULL | 收货地址 |
| `phone` | `varchar(18)` | YES | NULL | 联系电话 |
| `mark` | `varchar(200)` | YES | '' | 备注 |
| `mer_mark` | `varchar(255)` | YES | '' | 商户备注 |
| `admin_mark` | `varchar(255)` | YES | '' | 平台备注 |
| `pics` | `text` | YES | — | 图片 |
| `refund_type` | `tinyint(1)` | NO | — | 退款类型 1:退款 2:退款退货 |
| `refund_message` | `varchar(128)` | NO | — | 退款原因 |
| `refund_price` | `decimal(8,2)` | NO | '0.00' | 退款金额 |
| `platform_refund_price` | `decimal(8,2)` | YES | '0.00' | 退款平台优惠券金额 |
| `refund_postage` | `decimal(8,2)` | YES | '0.00' | 退的运费 |
| `refund_num` | `int(10)` | NO | '0' | 退款数 |
| `fail_message` | `varchar(200)` | YES | NULL | 未通过原因 |
| `status` | `tinyint(1)` | NO | '0' | 状态 0:待审核 -1:审核未通过 1:待退货 2:待收货 3:已退款 |
| `status_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 状态改变时间 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |
| `reconciliation_id` | `int(10)` | YES | '0' | 对账id |
| `is_del` | `tinyint(3)` | NO | '0' | — |
| `is_system_del` | `tinyint(1)` | YES | '0' | 商户删除 |
| `admin_id` | `int(10)` | YES | '0' | 管理/客服ID |
| `user_type` | `tinyint(1)` | YES | '1' | 用户类型 1 用户 2 平台 3 商户 4 客服  |

## `qixi_store_refund_product`

> 退款单产品表

原表：`eb_store_refund_product`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `refund_product_id` | `int(10)` | NO | — | 订单产品 id |
| `refund_order_id` | `int(10)` | NO | — | 退款单 |
| `order_product_id` | `int(10)` | NO | — | 订单产品id |
| `refund_price` | `decimal(8,2)` | NO | '0.00' | 退款金额 |
| `platform_refund_price` | `decimal(8,2)` | YES | NULL | 平台券退款金额 |
| `refund_postage` | `decimal(8,2)` | YES | NULL | 退邮费金额 |
| `refund_integral` | `int(10)` | YES | '0' | 退金额 |
| `refund_num` | `int(10)` | NO | '0' | 退货数 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |

## `qixi_store_refund_status`

> 订单操作记录表

原表：`eb_store_refund_status`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `refund_order_id` | `int(10)` | NO | — | 退款单订单id |
| `change_type` | `varchar(32)` | NO | — | 操作类型 |
| `change_message` | `varchar(256)` | NO | — | 操作备注 |
| `change_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 操作时间 |
