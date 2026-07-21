# 优惠券 — 表字段（`qixi_`）

## `qixi_store_coupon`

> 优惠券表

原表：`eb_store_coupon`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `coupon_id` | `int(10)` | NO | — | 优惠券表ID |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `is_timeout` | `tinyint(3)` | NO | '0' | 是否限时 |
| `start_time` | `timestamp` | YES | NULL | 优惠券领取开启时间 |
| `end_time` | `timestamp` | YES | NULL | 优惠券领取结束时间 |
| `is_limited` | `tinyint(3)` | NO | '0' | 是否限量 |
| `total_count` | `int(10)` | NO | '0' | 优惠券领取数量 |
| `remain_count` | `int(10)` | NO | '0' | 优惠券剩余领取数量 |
| `send_type` | `tinyint(3)` | NO | '0' | 0=领取 1=消费满赠 2=新人 3=买增 4=首单赠送 |
| `full_reduction` | `decimal(8,2)` | NO | '0.00' | 消费满多少赠送优惠券 |
| `title` | `varchar(64)` | NO | — | 优惠券名称 |
| `coupon_price` | `decimal(8,2)` | NO | '0.00' | 优惠券面值 |
| `use_min_price` | `int(11)` | NO | '0' | 最低消费多少金额可用优惠券 |
| `coupon_type` | `tinyint(3)` | NO | '0' | 优惠券类型 0=有效天数 1=固定时间段 |
| `coupon_time` | `int(10)` | NO | '0' | 优惠券有效期限（单位：天） |
| `use_start_time` | `timestamp` | YES | NULL | 开始时间 |
| `use_end_time` | `timestamp` | YES | NULL | 到期时间 |
| `sort` | `int(10)` | NO | '1' | 排序 |
| `status` | `tinyint(1)` | NO | '0' | 状态（0：关闭，1：开启 -1: 失效） |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `is_del` | `tinyint(3)` | NO | '0' | 是否删除 |
| `type` | `tinyint(4)` | NO | '0' | 优惠券类型 0-店铺 1-商品券 10 平台通用券 11平台品类券 12 平台跨店券 |

## `qixi_store_coupon_issue_user`

> 优惠券前台用户领取记录表

原表：`eb_store_coupon_issue_user`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `uid` | `int(11)` | NO | '0' | 领取优惠券用户ID |
| `coupon_id` | `int(11)` | NO | '0' | 优惠券ID |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 领取时间 |

## `qixi_store_coupon_product`

> 优惠卷关联商品辅助表

原表：`eb_store_coupon_product`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_id` | `int(11)` | NO | '0' | 产品id |
| `coupon_id` | `int(11)` | NO | '0' | 优惠卷id |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |

## `qixi_store_coupon_send`

> 优惠券发送记录

原表：`eb_store_coupon_send`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `coupon_send_id` | `int(10)` | NO | — | — |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `coupon_id` | `int(10)` | NO | — | 优惠券 id |
| `coupon_num` | `int(10)` | NO | '0' | 发送数量 |
| `mark` | `varchar(512)` | NO | — | 发送群体 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `status` | `tinyint(1)` | NO | '0' | 0:发送中 1:全部发送 |

## `qixi_store_coupon_user`

> 优惠券发放记录表

原表：`eb_store_coupon_user`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `coupon_user_id` | `int(11)` | NO | — | 优惠券发放记录id |
| `coupon_id` | `int(10)` | NO | '0' | 兑换的项目id |
| `mer_id` | `int(10)` | NO | '0' | 商户 id |
| `uid` | `int(10)` | NO | '0' | 优惠券所属用户 |
| `coupon_title` | `varchar(32)` | NO | — | 优惠券名称 |
| `coupon_price` | `decimal(8,2)` | NO | '0.00' | 优惠券的面值 |
| `use_min_price` | `int(11)` | NO | '0' | 最低消费多少金额可用优惠券 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 优惠券创建时间 |
| `start_time` | `timestamp` | YES | NULL | 优惠券开启时间 |
| `end_time` | `timestamp` | YES | NULL | 优惠券结束时间 |
| `use_time` | `timestamp` | YES | NULL | 使用时间 |
| `type` | `varchar(16)` | NO | 'send' | 获取方式(receive:自己领取 send:后台发送  give:满赠  new:新人 buy:买赠送) |
| `send_id` | `int(10)` | YES | '0' | 批量发送 id |
| `status` | `tinyint(1)` | NO | '0' | 状态（0：未使用，1：已使用, 2:已过期） |
| `is_fail` | `tinyint(3)` | NO | '0' | 是否有效 |
