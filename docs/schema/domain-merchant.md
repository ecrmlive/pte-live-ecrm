# 商户 — 表字段（`qixi_`）

## `qixi_merchant`

> 商户表

原表：`eb_merchant`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `mer_id` | `int(10)` | NO | — | 商户id |
| `category_id` | `int(10)` | NO | '0' | 商户分类 id |
| `type_id` | `int(10)` | YES | '0' | 店铺类型 id |
| `mer_name` | `varchar(32)` | NO | '' | 商户名称 |
| `real_name` | `varchar(32)` | NO | '' | 商户姓名 |
| `mer_phone` | `varchar(13)` | NO | '' | 商户手机号 |
| `mer_address` | `varchar(64)` | NO | '' | 商户地址 |
| `mer_keyword` | `varchar(64)` | NO | '' | 商户关键字 |
| `mer_avatar` | `varchar(128)` | YES | NULL | 商户头像 |
| `mer_banner` | `varchar(128)` | YES | NULL | 商户banner图片 |
| `mini_banner` | `varchar(128)` | YES | NULL | 商户店店铺街图片 |
| `sales` | `int(10)` | YES | '0' | 销量 |
| `product_score` | `decimal(11,1)` | YES | '5.0' | 商品描述评分 |
| `service_score` | `decimal(11,1)` | YES | '5.0' | 服务评分 |
| `postage_score` | `decimal(11,1)` | YES | '5.0' | 物流评分 |
| `mark` | `varchar(256)` | NO | — | 商户备注 |
| `reg_admin_id` | `int(10)` | NO | '0' | 总后台管理员ID |
| `sort` | `int(10)` | NO | '0' | — |
| `status` | `tinyint(1)` | NO | '0' | 商户是否禁用0锁定,1正常 |
| `commission_rate` | `decimal(6,2)` | YES | NULL | 提成比例 |
| `commission_switch` | `int(11)` | YES | '0' | 商户手续费单独设置 0 关闭 1 开启 |
| `long` | `varchar(16)` | YES | NULL | 经度 |
| `lat` | `varchar(16)` | YES | NULL | 纬度 |
| `is_del` | `tinyint(3)` | NO | '0' | 0未删除1删除 |
| `is_audit` | `tinyint(3)` | NO | '0' | 添加的产品是否审核0不审核1审核 |
| `is_bro_room` | `tinyint(3)` | NO | '1' | 是否审核直播间0不审核1审核 |
| `is_bro_goods` | `tinyint(3)` | NO | '1' | 是否审核直播商品0不审核1审核 |
| `is_best` | `tinyint(3)` | NO | '0' | 是否推荐 |
| `is_trader` | `tinyint(3)` | NO | '0' | 是否自营 |
| `mer_state` | `tinyint(3)` | NO | '0' | 商户是否1开启0关闭 |
| `mer_info` | `varchar(256)` | NO | '' | 店铺简介 |
| `service_phone` | `varchar(13)` | NO | '' | 店铺电话 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | — |
| `update_time` | `datetime` | NO | CURRENT_TIMESTAMP | — |
| `care_count` | `int(10)` | YES | '0' | 关注总数 |
| `copy_product_num` | `int(10)` | YES | '0' | 剩余复制商品次数 |
| `export_dump_num` | `int(10)` | YES | '0' | 电子面单剩余次数 |
| `mer_money` | `decimal(12,2)` | NO | '0.00' | 商户余额 |
| `financial_bank` | `varchar(255)` | YES | NULL | 银行卡转账信息 |
| `financial_wechat` | `varchar(255)` | YES | NULL | 微信转账信息 |
| `financial_alipay` | `varchar(255)` | YES | NULL | 支付宝转账信息 |
| `financial_type` | `tinyint(3)` | YES | '1' | 默认使用类型 |
| `sub_mchid` | `varchar(16)` | NO | '' | 微信支付分配的分账号 |
| `delivery_way` | `varchar(50)` | YES | '' | 配送方式 |
| `delivery_balance` | `decimal(8,2)` | NO | '0.00' | 配送余额 |
| `margin` | `decimal(8,2)` | NO | '0.00' | 保证金 |
| `margin_remind_time` | `varchar(255)` | YES | NULL | 保证金补缴提醒结束时间，时间点到了就自动关闭店铺 |
| `ot_margin` | `decimal(8,2)` | NO | '0.00' | 保证金额度 |
| `is_margin` | `tinyint(4)` | NO | '0' | 是否有保证金（0无，1有未支付，10已支付，-1 申请退款, -10 拒绝退款） |
| `offline_switch` | `tinyint(4)` | YES | '0' | 线下支付功能开关 |
| `care_ficti` | `int(11)` | YES | '0' | 虚拟关注量 |
| `region_id` | `int(11)` | NO | '0' | 商户所属分组 |
| `applyment_id` | `varchar(50)` | NO | '' | 特约商户ID |
| `business_id` | `int(10)` | NO | '0' | 店铺所属商户id |
| `applyment_switch` | `int(11)` | NO | '1' | 特约商户设置的分账比例是否合理 |

## `qixi_merchant_admin`

> 商户管理员表

原表：`eb_merchant_admin`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `merchant_admin_id` | `smallint(5)` | NO | — | 商户管理员表ID |
| `mer_id` | `int(10)` | NO | — | 商户ID(属于哪一个商户) |
| `account` | `varchar(32)` | NO | — | 商户管理员账号 |
| `pwd` | `char(64)` | NO | — | 商户管理员密码 |
| `real_name` | `varchar(16)` | NO | — | 商户管理员姓名 |
| `phone` | `varchar(13)` | YES | NULL | 商户管理员手机号 |
| `last_ip` | `varchar(16)` | YES | NULL | 商户管理员最后一次登录IP地址 |
| `last_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 商户管理员最后一次登录时间 |
| `roles` | `varchar(128)` | YES | '' | — |
| `login_count` | `int(10)` | NO | '0' | 商户管理员登录次数 |
| `level` | `tinyint(3)` | NO | '1' | 商户管理员等级(管理员添加的为0, 商户添加的为1) |
| `is_del` | `tinyint(3)` | NO | '0' | 是否删除 |
| `status` | `tinyint(3)` | NO | '1' | 是否有效 1有效 0无效  |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 商户管理员添加时间 |

## `qixi_merchant_applyments`

> 商户申请分账商户号表

原表：`eb_merchant_applyments`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `mer_applyments_id` | `int(10)` | NO | — | — |
| `out_request_no` | `varchar(128)` | YES | NULL | 业务申请编号 |
| `applyment_id` | `varchar(100)` | YES | NULL | 微信支付分配的申请单号 |
| `mer_id` | `int(11)` | YES | '0' | 商户ID |
| `sub_mchid` | `varchar(100)` | YES | NULL | 二级商户号/特约商户号 |
| `mer_name` | `varchar(50)` | YES | NULL | 商户名 |
| `info` | `text` | YES | — | 申请资料 |
| `status` | `int(11)` | YES | '0' | 申请状态: 0.平台未提交，-1.平台驳回，10.平台提交审核中，11.需用户操作 ，20.已完成，30.已冻结，40.驳回 |
| `message` | `varchar(1000)` | YES | NULL | 返回信息 |
| `mark` | `varchar(255)` | YES | NULL | 备注 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `is_del` | `int(11)` | YES | '0' | 删除 |
| `type` | `tinyint(2)` | NO | '0' | 0 平台收付通 1 服务商分账 |

## `qixi_merchant_category`

> 商户分类表

原表：`eb_merchant_category`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `merchant_category_id` | `int(10)` | NO | — | 商户分类 id |
| `commission_rate` | `decimal(6,4)` | NO | '0.0000' | 手续费 |
| `category_name` | `varchar(32)` | NO | — | 商户分类名称 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |

## `qixi_merchant_intention`

> 商户申请表

原表：`eb_merchant_intention`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `mer_intention_id` | `int(10)` | NO | — | ID |
| `uid` | `int(10)` | YES | '0' | 用户ID |
| `phone` | `varchar(11)` | YES | NULL | 手机号 |
| `mer_name` | `varchar(30)` | YES | NULL | 商户名称 |
| `name` | `varchar(30)` | YES | NULL | 客户姓名 |
| `create_time` | `datetime` | YES | CURRENT_TIMESTAMP | 提交时间 |
| `status` | `tinyint(4)` | YES | '0' | 处理状态 1通过 ，2未通过 |
| `fail_msg` | `varchar(255)` | YES | NULL | 未通过原因 |
| `is_del` | `tinyint(4)` | YES | '0' | 删除状态 1删除 ，0未删除 |
| `mark` | `varchar(255)` | YES | NULL | 备注 |
| `mer_id` | `int(10)` | YES | '0' | 关联商户 |
| `images` | `varchar(2000)` | YES | NULL | 多图 |
| `merchant_category_id` | `int(10)` | YES | '0' | 商户分类 |
| `mer_type_id` | `int(10)` | YES | '0' | 店铺类型 |
| `circle_id` | `int(10)` | NO | '0' | 商圈id |
| `business_id` | `int(10)` | NO | '0' | 店铺所属商户id |

## `qixi_merchant_region`

原表：`eb_merchant_region`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `region_id` | `int(11)` | NO | — | — |
| `name` | `varchar(255)` | NO | '' | 名称 |
| `info` | `varchar(255)` | NO | '' | 简介 |
| `pid` | `int(11)` | NO | '0' | 父级ID |
| `path` | `varchar(255)` | NO | '/' | 父级路径 |
| `pic` | `varchar(255)` | NO | '' | 图片 |
| `lv` | `int(11)` | NO | '0' | 等级 |
| `sort` | `int(11)` | NO | '0' | 排序 |
| `status` | `int(11)` | NO | '0' | 状态 0 关闭, 1 开启 |
| `type` | `int(11)` | NO | '0' | 类型:0 |
| `mer_id` | `int(11)` | NO | '0' | 商户ID |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 创建时间 |
| `update_time` | `timestamp` | YES | NULL | 更新时间 |

## `qixi_merchant_type`

> 商户类型表

原表：`eb_merchant_type`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `mer_type_id` | `int(10)` | NO | — | 商户类型 id |
| `type_name` | `varchar(16)` | NO | — | 类型名称 |
| `type_info` | `varchar(512)` | YES | NULL | 类型要求 |
| `description` | `varchar(512)` | YES | NULL | 类型说明 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `margin` | `decimal(8,2)` | NO | '0.00' | 保证金 |
| `is_margin` | `tinyint(3)` | NO | '0' | 是否有保证金（0无，1有） |
| `mark` | `varchar(255)` | NO | '' | 备注 |
| `update_time` | `datetime` | YES | NULL | 更新时间 |

## `qixi_store_group`

> 分组表

原表：`eb_store_group`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `store_group_id` | `int(10)` | NO | — | 分组id |
| `pid` | `int(10)` | NO | '0' | 上级商圈id |
| `path` | `varchar(50)` | NO | '' | 路径 |
| `name` | `varchar(64)` | NO | '' | 名称 |
| `level` | `tinyint(3)` | NO | '0' | 等级:0一级 1二级 2三级 |
| `positioning_status` | `tinyint(3)` | NO | '1' | 是否开启定位：0否 1是 |
| `longitude` | `varchar(16)` | NO | '' | 经度，positioning_status 为1时有效 |
| `latitude` | `varchar(16)` | NO | '' | 纬度，positioning_status 为1时有效 |
| `address` | `varchar(100)` | NO | '' | 中心点地址，positioning_status 为1时有效 |
| `diy_temp_id` | `int(10)` | NO | '0' | 首页模板id |
| `remark` | `varchar(255)` | NO | '' | 说明信息 |
| `sort` | `int(10)` | NO | '0' | 排序(数字越大越靠前) |
| `status` | `tinyint(3)` | NO | '1' | 状态：0关闭 1开启 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `update_time` | `datetime` | NO | CURRENT_TIMESTAMP | 修改时间 |
