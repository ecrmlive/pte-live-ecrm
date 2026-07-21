# 财务/结算 — 表字段（`qixi_`）

## `qixi_circle_financial_record`

> 商圈提成流水

原表：`eb_circle_financial_record`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `record_id` | `int(10)` | NO | — | — |
| `circle_id` | `int(10)` | NO | '0' | 商圈id |
| `circle_name` | `varchar(64)` | NO | '' | 商圈名称 |
| `mer_id` | `int(10)` | NO | '0' | 商户id |
| `mer_name` | `varchar(64)` | NO | '' | 商户名称 |
| `agent_id` | `int(10)` | NO | '0' | 商圈代理id |
| `agent_name` | `varchar(64)` | NO | '' | 商圈代理名称 |
| `record_sn` | `varchar(32)` | NO | '' | 流水号 |
| `order_id` | `int(10)` | NO | '0' | 订单号 |
| `order_sn` | `varchar(32)` | NO | '' | 订单编号 |
| `user_id` | `int(10)` | NO | '0' | 用户 id |
| `user_info` | `varchar(32)` | NO | '' | 用户名 |
| `amount` | `decimal(8,2)` | NO | '0.00' | 金额 |
| `status` | `tinyint(3)` | NO | '0' | 状态：0冻结中 1解冻 -1失效 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `update_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 修改时间 |

## `qixi_financial`

> 商户财务申请提现

原表：`eb_financial`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `financial_id` | `int(10)` | NO | — | — |
| `financial_sn` | `varchar(32)` | NO | — | 单号 |
| `mer_money` | `decimal(12,2)` | NO | — | 余额 |
| `extract_money` | `decimal(12,2)` | NO | '0.00' | 提现金额 |
| `financial_type` | `int(10)` | YES | '0' | 收款类型 |
| `financial_account` | `varchar(500)` | NO | — | 商户账户信息 |
| `financial_status` | `int(10)` | YES | '0' | 转账状态 |
| `status` | `int(11)` | NO | — | 审核0待审核，1通过 ，-1 未通过 |
| `refusal` | `varchar(32)` | YES | NULL | 拒绝理由 |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `image` | `varchar(1000)` | YES | NULL | 凭证 |
| `admin_id` | `int(11)` | YES | NULL | — |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 创建时间 |
| `status_time` | `timestamp` | YES | NULL | 审核时间 |
| `update_time` | `timestamp` | YES | NULL | 修改拼凭证时间 |
| `is_del` | `int(10)` | YES | '0' | — |
| `mark` | `varchar(255)` | YES | NULL | 商户备注 |
| `admin_mark` | `varchar(255)` | YES | NULL | 平台备注 |
| `mer_admin_id` | `int(11)` | YES | NULL | 商户管理员 |
| `type` | `int(10)` | YES | '0' | 申请类型 0.余额 1 保证金 |

## `qixi_financial_record`

> 商户财务流水

原表：`eb_financial_record`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `financial_record_id` | `int(10)` | NO | — | — |
| `financial_record_sn` | `varchar(32)` | NO | — | 流水号 |
| `order_id` | `int(10)` | NO | — | 订单号 |
| `order_sn` | `varchar(32)` | NO | — | 订单编号 |
| `user_info` | `varchar(32)` | NO | — | 用户名 |
| `user_id` | `int(10)` | NO | — | 用户 id |
| `financial_type` | `varchar(32)` | NO | — | 流水类型 |
| `financial_pm` | `tinyint(3)` | NO | '0' | 0 = 支出 1 = 获得 |
| `number` | `decimal(8,2)` | NO | '0.00' | 金额 |
| `type` | `tinyint(1)` | NO | '-1' | 0:商户  1:公共  2:平台 |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `pay_type` | `int(11)` | NO | — | 支付类型 |

## `qixi_serve_meal`

原表：`eb_serve_meal`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `meal_id` | `int(10)` | NO | — | — |
| `name` | `varchar(30)` | YES | NULL | 套餐名称 |
| `type` | `int(11)` | YES | '0' | 套餐类型,1复制商品，2电子面单 |
| `price` | `decimal(8,2)` | YES | '0.00' | 价格 |
| `num` | `int(11)` | YES | '1' | 数量 |
| `sort` | `int(11)` | YES | NULL | 排序 |
| `status` | `int(11)` | YES | '1' | 状态 |
| `is_del` | `int(11)` | YES | '0' | — |
| `create_time` | `timestamp` | YES | NULL | — |

## `qixi_serve_order`

原表：`eb_serve_order`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `order_id` | `int(10)` | NO | — | — |
| `meal_id` | `int(11)` | YES | NULL | 套餐ID |
| `pay_type` | `int(11)` | YES | NULL | 支付方式：1微信，2支付宝,3 平台操作 |
| `order_sn` | `varchar(50)` | YES | NULL | 订单ID |
| `pay_price` | `decimal(8,2)` | YES | NULL | 价格 |
| `order_info` | `varchar(255)` | YES | NULL | 套餐信息 |
| `type` | `int(11)` | YES | NULL | 套餐类型 1 采集 2 电子面单 10 保证金 20同城配送 30 会员充值 |
| `status` | `int(11)` | YES | '0' | 状态：默认0，支付成功 1，支付失败 -1，20已退款 |
| `mer_id` | `int(11)` | YES | NULL | 商户ID/用户ID |
| `create_time` | `timestamp` | YES | NULL | — |
| `is_del` | `int(11)` | YES | '0' | — |
| `pay_time` | `timestamp` | YES | NULL | 支付时间 |
