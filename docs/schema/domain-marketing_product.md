# 营销商品(秒杀/拼团/预售/助力) — 表字段（`qixi_`）

## `qixi_presell_order`

原表：`eb_presell_order`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `presell_order_id` | `int(10)` | NO | — | 预售尾款订单id |
| `presell_order_sn` | `varchar(32)` | NO | — | 预售订单号 |
| `uid` | `int(10)` | NO | — | 用户 id |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `order_id` | `int(10)` | NO | — | 订单id |
| `transaction_id` | `varchar(60)` | YES | NULL | 微信支付订单号(分账时有效) |
| `final_start_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 支付开始时间 |
| `final_end_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 尾款支付结时间 |
| `paid` | `tinyint(3)` | NO | '0' | 0:未支付 1:已支付 |
| `status` | `tinyint(3)` | YES | '1' | 0:无效 1:有效 |
| `pay_type` | `tinyint(3)` | NO | '0' | 支付方式 0余额 1微信 2小程序 3,4支付宝 |
| `pay_price` | `decimal(8,2)` | NO | — | 尾款 |
| `refun_price` | `decimal(8,2)` | NO | '0.00' | 退款金额 |
| `pay_time` | `timestamp` | YES | NULL | 支付时间 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `is_combine` | `tinyint(3)` | YES | '0' | 是否为合并支付 |

## `qixi_store_activity`

原表：`eb_store_activity`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `activity_id` | `int(11)` | NO | — | ID |
| `activity_name` | `varchar(128)` | NO | '' | 活动名称 |
| `start_time` | `timestamp` | YES | NULL | 开始时间 |
| `end_time` | `timestamp` | YES | NULL | 结束时间 |
| `pic` | `varchar(128)` | YES | '' | 图片 |
| `is_show` | `tinyint(1)` | YES | '0' | 是否显示 |
| `status` | `tinyint(1)` | NO | '0' | 状态0未开始，1进行中，2已结束 |
| `activity_type` | `tinyint(1)` | NO | '1' | 1.氛围图 2.边框 |
| `create_time` | `datetime` | YES | CURRENT_TIMESTAMP | 添加时间 |
| `is_del` | `tinyint(1)` | NO | '0' | 是否删除 |
| `scope_type` | `tinyint(1)` | NO | '0' | 指定类型：0全部商品1指定商品2指定分类3指定商户4秒杀活动 |
| `images` | `varchar(500)` | YES | '' | 多图 |
| `info` | `varchar(500)` | YES | '' | 简介 |
| `color` | `varchar(128)` | YES | '' | 背景色 |
| `sort` | `int(11)` | YES | '0' | 排序 |
| `mer_id` | `int(11)` | YES | '0' | 商户ID |
| `link_id` | `int(11)` | YES | NULL | 关联ID |
| `update_time` | `timestamp` | YES | NULL | — |
| `count` | `int(10)` | YES | '0' | 需要的总数 |
| `total` | `int(10)` | YES | '0' | 已有的总数 |
| `is_display` | `tinyint(1)` | YES | '1' | 是否在活动列表中显示 |

## `qixi_store_activity_cate`

> 活动标签分类表

原表：`eb_store_activity_cate`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `int(11)` | NO | — | — |
| `name` | `varchar(50)` | NO | '' | — |
| `pid` | `int(11)` | NO | '0' | — |
| `pic` | `varchar(255)` | NO | '' | — |
| `sort` | `int(11)` | NO | '0' | — |
| `status` | `tinyint(2)` | NO | '0' | — |
| `mer_id` | `int(11)` | NO | '0' | — |
| `type` | `tinyint(2)` | NO | '0' | — |
| `path` | `varchar(255)` | NO | '/' | — |
| `lv` | `int(11)` | NO | '0' | — |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |

## `qixi_store_activity_label`

> 活动标签表

原表：`eb_store_activity_label`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `int(11)` | NO | — | — |
| `type` | `tinyint(1)` | NO | '0' | 类型：0平台 2:商户 |
| `mer_id` | `int(11)` | NO | '0' | 商户id |
| `label_cate` | `int(11)` | NO | '0' | 标签分类 |
| `label_name` | `varchar(255)` | NO | '' | 标签名称 |
| `style_type` | `tinyint(1)` | NO | '1' | 样式类型1：自定义2：图片 |
| `color` | `varchar(32)` | NO | '' | 颜色 |
| `bg_color` | `varchar(32)` | NO | '' | 背景颜色 |
| `border_color` | `varchar(32)` | NO | '' | 边框颜色 |
| `icon` | `varchar(255)` | YES | '' | 图标 |
| `default_type` | `tinyint(1)` | NO | '0' | 系统默认标签分类(1:包邮,2:领券,3:上门,4:到店,5:同城,6:拼团,7:秒杀,8:助力,9:预售,10:自营) |
| `is_show` | `tinyint(1)` | NO | '1' | 移动端是否展示 |
| `status` | `tinyint(1)` | NO | '1' | 状态是否开启 |
| `sort` | `int(11)` | NO | '0' | 排序 |
| `add_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 添加时间 |

## `qixi_store_activity_related`

原表：`eb_store_activity_related`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `bigint(20)` | NO | — | — |
| `activity_id` | `int(11)` | NO | '0' | 活动ID |
| `activity_type` | `varchar(255)` | YES | NULL | 活动类型 |
| `keys` | `varchar(2000)` | YES | NULL | 主要信息 |
| `value` | `varchar(2000)` | YES | NULL | 活动值 |
| `form_value` | `text` | YES | — | 表单内容 |
| `uid` | `int(11)` | NO | '0' | 用户ID |
| `link_id` | `int(11)` | YES | NULL | — |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `nickname` | `varchar(255)` | YES | NULL | — |
| `avatar` | `varchar(255)` | YES | NULL | — |
| `phone` | `varchar(20)` | YES | NULL | — |
| `is_del` | `tinyint(1)` | NO | '0' | 是否删除 1|是 0|否 |

## `qixi_store_discounts`

原表：`eb_store_discounts`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `discount_id` | `int(11)` | NO | — | 自增ID |
| `title` | `varchar(255)` | NO | '' | 套餐名称 |
| `image` | `varchar(500)` | NO | '' | 组合套餐主图 |
| `type` | `tinyint(1)` | NO | '0' | 套餐类型0固定1搭配 |
| `is_limit` | `tinyint(1)` | NO | '0' | 是否限量0不限量1限量 |
| `limit_num` | `int(11)` | NO | '0' | 限量个数 |
| `link_ids` | `varchar(255)` | NO | '' | 关联标签 |
| `product_ids` | `varchar(255)` | YES | '' | 商品IDS |
| `is_time` | `tinyint(1)` | NO | '0' | 是否限时0不限时1限时 |
| `start_time` | `int(11)` | NO | '0' | 开始时间 |
| `stop_time` | `int(11)` | NO | '0' | 结束时间 |
| `sort` | `int(11)` | NO | '0' | 排序 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `free_shipping` | `tinyint(1)` | NO | '0' | 是否包邮0不包邮1包邮 |
| `status` | `tinyint(1)` | NO | '1' | 平台是否上架0不上架1上架 |
| `is_del` | `tinyint(1)` | NO | '0' | 是否删除0未删除1已删除 |
| `mer_id` | `int(11)` | YES | NULL | 商户ID |
| `is_show` | `tinyint(1)` | YES | '0' | 商户是否上架0不上架1上架 |
| `sales` | `int(10)` | YES | NULL | 销量 |

## `qixi_store_discounts_product`

原表：`eb_store_discounts_product`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `discount_product_id` | `int(11)` | NO | — | 自增ID |
| `discount_id` | `int(11)` | NO | — | 优惠套餐ID |
| `product_id` | `int(11)` | NO | — | 商品ID |
| `store_name` | `varchar(255)` | NO | — | 商品名称 |
| `image` | `varchar(500)` | NO | '' | 商品图 |
| `type` | `tinyint(1)` | NO | '0' | 是否搭配0不是1是 |
| `temp_id` | `int(11)` | NO | '0' | 运费模版Id |
| `mer_id` | `int(11)` | YES | NULL | — |

## `qixi_store_seckill_active`

> 商户设置秒杀商品关联表

原表：`eb_store_seckill_active`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `seckill_active_id` | `int(10)` | NO | — | ID |
| `name` | `varchar(64)` | NO | — | 活动名称 |
| `seckill_time_ids` | `varchar(255)` | YES | '' | 活动场次 |
| `start_day` | `timestamp` | NO | '0000-00-00 00:00:00' | 开始日期 |
| `end_day` | `timestamp` | NO | '0000-00-00 00:00:00' | 结束日期 |
| `start_time` | `int(10)` | NO | — | 开始时间 |
| `end_time` | `int(10)` | NO | — | 结束时间 |
| `mer_id` | `int(10)` | NO | — | 商户ID |
| `all_pay_count` | `int(10)` | NO | '0' | 活动有效期内每个用户可购买该商品总数限制 |
| `once_pay_count` | `int(10)` | NO | '0' | 单次购买最大数量限制 |
| `product_id` | `int(10)` | NO | — | 商品ID |
| `product_category_ids` | `varchar(255)` | YES | NULL | 平台一级商品分类/为空均可参与 |
| `merchant_count` | `int(11)` | NO | '0' | 商户数量 |
| `product_count` | `int(11)` | NO | '0' | 商品数量 |
| `active_status` | `enum('0','1','-1')` | YES | '0' | 活动状态: 0未开始 1 进行中 -1 已结束 |
| `sign` | `int(11)` | NO | — | 标识 1=秒杀活动 |
| `status` | `tinyint(3)` | YES | '0' | 0=未开启,1=已开启 |
| `create_time` | `int(11)` | NO | — | 创建时间 |
| `update_time` | `int(11)` | NO | — | 修改时间 |
| `delete_time` | `int(11)` | YES | NULL | 删除时间 |

## `qixi_store_seckill_time`

> 秒杀时间段配置

原表：`eb_store_seckill_time`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `seckill_time_id` | `int(10)` | NO | — | — |
| `title` | `varchar(255)` | YES | '' | — |
| `start_time` | `int(10)` | NO | — | 开始时间 |
| `end_time` | `int(10)` | NO | — | 结束时间 |
| `status` | `tinyint(3)` | NO | '0' | 1，0状态 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `pic` | `varchar(255)` | YES | NULL | 图片 |
