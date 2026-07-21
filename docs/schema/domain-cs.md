# 客服 — 表字段（`qixi_`）

## `qixi_store_service`

> 客服表

原表：`eb_store_service`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `service_id` | `int(11)` | NO | — | 客服id |
| `mer_id` | `int(11)` | NO | '0' | 商户id |
| `uid` | `int(11)` | NO | — | 客服uid |
| `avatar` | `varchar(250)` | NO | — | 客服头像 |
| `nickname` | `varchar(50)` | NO | — | 客服名称 |
| `account` | `varchar(32)` | YES | NULL | 客服账号 |
| `pwd` | `varchar(64)` | YES | NULL | 客服密码 |
| `is_open` | `tinyint(3)` | NO | '0' | 开启 pc 登录 |
| `status` | `tinyint(3)` | NO | '0' | 0隐藏1显示 |
| `notify` | `int(11)` | YES | '0' | 订单通知1开启0关闭 |
| `phone` | `varchar(18)` | YES | '' | 电话 |
| `customer` | `tinyint(1)` | NO | '0' | 是否展示统计管理 |
| `is_verify` | `tinyint(3)` | NO | '0' | 是否有核销权限 |
| `is_goods` | `tinyint(3)` | YES | '0' | 是否有商品管理权限 |
| `sort` | `tinyint(3)` | NO | '0' | 排序 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `is_del` | `tinyint(3)` | NO | '0' | 是否删除 |

## `qixi_store_service_log`

> 客服用户对话记录表

原表：`eb_store_service_log`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `service_log_id` | `int(11)` | NO | — | 客服用户对话记录表ID |
| `mer_id` | `int(11)` | NO | '0' | 商户id |
| `msn` | `varchar(200)` | NO | — | 消息内容 |
| `uid` | `int(11)` | NO | — | 发送人uid |
| `service_id` | `int(10)` | NO | '0' | 客服 id |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 发送时间 |
| `type` | `tinyint(1)` | NO | '0' | 是否已读（0：否；1：是；） |
| `service_type` | `tinyint(1)` | NO | '0' | 客服是否已读（0：否；1：是；） |
| `remind` | `tinyint(1)` | NO | '0' | 是否提醒过（0：否；1：是；） |
| `send_type` | `tinyint(3)` | NO | '0' | 0:用户发送 1:客服回复 |
| `msn_type` | `tinyint(3)` | NO | '1' | 消息类型 1=文字 2=表情 3=图片 4=商品 5=订单 6=退款单 |

## `qixi_store_service_reply`

原表：`eb_store_service_reply`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `service_reply_id` | `int(10)` | NO | — | — |
| `mer_id` | `int(10)` | NO | '0' | 商户 id |
| `type` | `tinyint(3)` | NO | '1' | 1:文字 2:图片 |
| `keyword` | `varchar(64)` | NO | — | 回复的关键字 |
| `content` | `varchar(512)` | NO | — | 回复内容 |
| `status` | `tinyint(3)` | YES | '1' | 是否开启 |
| `create_time` | `datetime` | YES | CURRENT_TIMESTAMP | — |

## `qixi_store_service_user`

原表：`eb_store_service_user`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `service_user_id` | `int(10)` | NO | — | 聊天用户 id |
| `service_id` | `int(10)` | NO | — | 客服 id |
| `uid` | `int(11)` | NO | — | 用户 id |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `is_online` | `tinyint(3)` | YES | '0' | 是否在线 |
| `service_unread` | `smallint(5)` | YES | '0' | 客服未读数 |
| `user_unread` | `smallint(5)` | YES | '0' | 用户未读数 |
| `last_log_id` | `int(10)` | NO | '0' | 最后一条记录 id |
| `last_time` | `datetime` | NO | — | 最后发送时间 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 创建时间 |
