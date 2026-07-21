# 配送/运费 — 表字段（`qixi_`）

## `qixi_delivery_config`

> 配送设置表

原表：`eb_delivery_config`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `delivery_config_id` | `int(11)` | NO | — | ID |
| `mer_id` | `int(11)` | NO | '0' | 商户ID |
| `min_delivery_amount` | `decimal(8,2)` | NO | '0.00' | 起送价 |
| `base_shipping_fee` | `decimal(8,2)` | NO | '0.00' | 基础运费 |
| `free_shipping_amount` | `decimal(8,2)` | NO | '0.00' | 包邮规则 |
| `is_premium_stack_enabled` | `int(11)` | NO | '0' | 是否开启溢价叠加(0:关 1:开) |
| `distance_premium_config` | `text` | NO | — | 距离溢价设置 |
| `weight_premium_config` | `text` | NO | — | 重量溢价设置 |
| `delivery_time_type` | `int(11)` | NO | '1' | 配送时间类型(1:可选定时送达 2:统一尽快送达) |
| `selectable_days` | `int(11)` | NO | '7' | 可选天数 |
| `delivery_prompt` | `varchar(200)` | NO | '' | 送达文案 |
| `status` | `int(11)` | NO | '0' | 配送状态(0:关 1:开) |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `update_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 修改时间 |
| `commission_rate` | `int(10)` | NO | '0' | 配送员提成比例(0~100%) |

## `qixi_delivery_order`

原表：`eb_delivery_order`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `delivery_order_id` | `int(10)` | NO | — | — |
| `station_id` | `int(10)` | NO | '0' | 门店ID |
| `order_id` | `int(10)` | NO | — | 订单ID |
| `order_code` | `varchar(255)` | NO | '' | 配送方订单号 |
| `city_code` | `varchar(20)` | NO | '' | 所属城市 |
| `order_sn` | `varchar(32)` | NO | '' | 订单sn |
| `cargo_price` | `decimal(8,2)` | NO | '0.00' | 配送订单价格 |
| `finish_code` | `varchar(255)` | NO | '' | 收货码 |
| `user_name` | `varchar(20)` | NO | '' | 用户名 |
| `status` | `int(11)` | NO | '0' | 状态 取消=-1, 待取货＝2,配送中＝3,已完成＝4,物品返回中=9,物品返回完成=10,骑士到店=100 |
| `receiver_phone` | `varchar(11)` | NO | '' | 收货人电话 |
| `from_address` | `varchar(255)` | NO | '' | 起始位置 |
| `to_address` | `varchar(255)` | NO | '' | 结束位置 |
| `distance` | `float` | NO | '0' | 配送距离 |
| `fee` | `decimal(8,2)` | NO | '0.00' | 配送费 |
| `mer_id` | `int(10)` | NO | '0' | 商户ID |
| `mark` | `varchar(255)` | NO | '' | 订单备注 |
| `station_type` | `int(10)` | NO | — | 平台类型 |
| `reason` | `varchar(255)` | NO | — | 取消原因 |
| `from_lat` | `varchar(255)` | NO | '' | — |
| `from_lng` | `varchar(255)` | NO | '' | — |
| `to_lat` | `varchar(255)` | NO | '' | — |
| `to_lng` | `varchar(255)` | NO | '' | — |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |
| `deduct_fee` | `decimal(8,2)` | NO | '0.00' | 取消订单违约金 |
| `uid` | `int(10)` | NO | '0' | — |
| `service_id` | `int(11)` | NO | '0' | 服务人员id |
| `service_fee` | `decimal(8,2)` | NO | '0.00' | 服务费 |

## `qixi_delivery_service`

原表：`eb_delivery_service`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `service_id` | `int(11)` | NO | — | id |
| `uid` | `int(11)` | NO | '0' | 配送员uid |
| `type` | `tinyint(1)` | NO | '1' | 类型：0平台1:商户 |
| `relation_id` | `int(11)` | NO | '0' | 门店、供应商id |
| `avatar` | `varchar(250)` | NO | '' | 配送员头像 |
| `name` | `varchar(50)` | NO | '' | 配送员名称 |
| `phone` | `varchar(20)` | NO | '0' | 手机号码 |
| `create_time` | `int(11)` | NO | '0' | 添加时间 |
| `is_del` | `tinyint(1)` | NO | '0' | 是否删除 |
| `status` | `tinyint(3)` | NO | '1' | 0隐藏1显示 |
| `mer_id` | `int(11)` | NO | '0' | 商户ID |
| `remark` | `varchar(255)` | NO | '' | 备注 |
| `sort` | `int(11)` | NO | '0' | 排序 |

## `qixi_delivery_station`

> 同城配送门店列表

原表：`eb_delivery_station`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `station_id` | `bigint(20)` | NO | — | — |
| `station_name` | `varchar(255)` | NO | '' | 门店名称 |
| `business` | `int(11)` | NO | '0' | 支持配送的物品品类 |
| `city_name` | `varchar(100)` | NO | '' | 门店所属市 |
| `station_address` | `varchar(255)` | NO | '' | 门店地址 |
| `lng` | `char(20)` | NO | '' | 门店经度 |
| `lat` | `char(20)` | NO | '' | 门店纬度 |
| `contact_name` | `char(10)` | NO | '' | 联系人姓名 |
| `phone` | `char(11)` | NO | '' | 联系人电话 |
| `origin_shop_id` | `varchar(255)` | NO | '' | 门店编码,可自定义,但必须唯一;若不填写,则系统自动生成 |
| `username` | `varchar(255)` | NO | '' | 达达商家app账号(若不需要登陆app,则不用设置) |
| `password` | `varchar(255)` | NO | '' | 达达商家app密码(若不需要登陆app,则不用设置)\n |
| `status` | `tinyint(2)` | NO | '1' | 状态 1启用 0关闭 |
| `mer_id` | `int(11)` | NO | '0' | 商户ID |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `mark` | `varchar(255)` | NO | '' | 备注 |
| `type` | `tinyint(2)` | NO | '0' | 类型 0 到店自提 1 达达 2 uu |
| `switch_city` | `tinyint(2)` | NO | '1' | 同城配送：1 支持 0 不支持 |
| `switch_take` | `tinyint(2)` | NO | '0' | 到店自提：1 支持 0 不支持 |
| `business_date` | `varchar(100)` | NO | '' | 营业日期 |
| `business_time_start` | `varchar(20)` | NO | '' | 营业开始时间 |
| `business_time_end` | `varchar(20)` | NO | '' | 营业结束时间 |
| `card_number` | `varchar(20)` | NO | '' | 身份证号 |
| `is_del` | `tinyint(2)` | NO | '0' | — |
| `range_type` | `int(11)` | NO | '1' | 距离设置类型(1:范围 2:行政区 3:电子围栏) |
| `radius` | `float(8,2)` | NO | '1.00' | 服务半径(km) |
| `region` | `varchar(200)` | NO | '' | 行政区域 |
| `fence` | `text` | NO | — | 电子围栏配置 |

## `qixi_express`

> 快递公司表

原表：`eb_express`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `mediumint(8)` | NO | — | 快递公司id |
| `code` | `varchar(50)` | NO | — | 快递公司简称 |
| `name` | `varchar(50)` | NO | — | 快递公司全称 |
| `mark` | `varchar(255)` | YES | '' | 备注 |
| `partner_id` | `int(11)` | YES | '0' | 月结账号 |
| `partner_key` | `int(11)` | YES | '0' | 月结密码 |
| `net` | `int(11)` | YES | '0' | 取件网点 |
| `sort` | `int(11)` | NO | — | 排序 |
| `is_show` | `int(11)` | YES | '1' | 是否显示 |
| `check_man` | `tinyint(4)` | NO | '0' | 是否需要承载快递员名称 |
| `partner_name` | `tinyint(4)` | NO | '0' | 是否需要客户账户名称 |
| `is_code` | `tinyint(4)` | NO | '0' | 是否需要承载编号 |
| `open_mer` | `varchar(255)` | NO | '' | 开启该物流公司的商户 |

## `qixi_express_partner`

原表：`eb_express_partner`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `int(10)` | NO | — | — |
| `express_id` | `int(11)` | NO | — | 快递公司id |
| `account` | `varchar(20)` | YES | NULL | 月结账号 |
| `key` | `varchar(50)` | YES | NULL | 月结密码 |
| `net_name` | `varchar(50)` | YES | NULL | 取件网点 |
| `mer_id` | `int(11)` | YES | NULL | — |
| `status` | `int(11)` | NO | '1' | — |
| `check_man` | `varchar(50)` | YES | NULL | 承载快递员名 |
| `partner_name` | `varchar(50)` | YES | NULL | 客户账户名称 |
| `code` | `varchar(50)` | YES | NULL | 承载编号 |

## `qixi_shipping_template`

> 运费表

原表：`eb_shipping_template`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `shipping_template_id` | `int(10)` | NO | — | 编号 |
| `name` | `varchar(255)` | NO | — | 模板名称 |
| `type` | `tinyint(3)` | NO | '1' | 计费方式 0=数量 1=重量 2=体积 |
| `appoint` | `tinyint(3)` | NO | '0' | 开启指定包邮 |
| `undelivery` | `tinyint(3)` | NO | '0' | 开启指定区域不配送 |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `is_default` | `tinyint(3)` | YES | '0' | 默认模板 |
| `sort` | `int(11)` | NO | '0' | 排序 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `info` | `varchar(1000)` | YES | NULL | 运费说明 |

## `qixi_shipping_template_free`

> 指定包邮信息表

原表：`eb_shipping_template_free`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `shipping_template_free_id` | `int(10)` | NO | — | 编号 |
| `temp_id` | `int(10)` | NO | '0' | 模板ID |
| `city_id` | `text` | NO | — | 城市ID /id/id/id/id/ |
| `number` | `int(10)` | NO | '0' | 包邮件数 |
| `price` | `decimal(10,2)` | NO | '0.00' | 包邮金额 |

## `qixi_shipping_template_region`

> 配送区域表

原表：`eb_shipping_template_region`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `shipping_template_region_id` | `int(10)` | NO | — | 编号 |
| `temp_id` | `int(10)` | NO | '0' | 模板ID |
| `city_id` | `text` | NO | — | 城市ID /id/id/id/ |
| `first` | `decimal(10,2)` | NO | '0.00' | 首件 |
| `first_price` | `decimal(10,2)` | NO | '0.00' | 首件运费 |
| `continue` | `decimal(10,2)` | NO | '0.00' | 续件 |
| `continue_price` | `decimal(10,2)` | NO | '0.00' | 续件运费 |

## `qixi_shipping_template_undelivery`

> 指定不配送区域表

原表：`eb_shipping_template_undelivery`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `shipping_template_undelivery_id` | `int(10)` | NO | — | 编号 |
| `temp_id` | `int(10)` | NO | '0' | 模板ID |
| `city_id` | `text` | NO | — | 城市ID /id/id/id/ |
