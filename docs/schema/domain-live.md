# 直播 — 表字段（`qixi_`）

## `qixi_broadcast_assistant`

> 直播助手信息

原表：`eb_broadcast_assistant`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `assistant_id` | `bigint(20)` | NO | — | — |
| `username` | `varchar(50)` | YES | NULL | 微信号 |
| `nickname` | `varchar(100)` | YES | NULL | 微信昵称 |
| `mer_id` | `int(11)` | YES | NULL | 商户ID |
| `mark` | `varchar(255)` | YES | NULL | 备注 |
| `is_del` | `tinyint(1)` | YES | '0' | — |

## `qixi_broadcast_goods`

> 直播商品表

原表：`eb_broadcast_goods`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `broadcast_goods_id` | `int(10)` | NO | — | — |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `goods_id` | `int(10)` | NO | '0' | 微信商品ID |
| `audit_id` | `int(10)` | NO | '0' | 审核单 id |
| `cover_img` | `varchar(255)` | NO | — | 图片 |
| `name` | `varchar(64)` | NO | — | 商品名称 |
| `price` | `decimal(10,2)` | NO | '0.00' | 价格 |
| `product_type` | `tinyint(3)` | NO | '0' | 商品类型 |
| `product_id` | `int(10)` | NO | '0' | 商品 id |
| `error_msg` | `varchar(255)` | YES | NULL | 未通过原因 |
| `audit_status` | `tinyint(3)` | NO | '0' | 0：未审核，1：审核中，2:审核通过，3审核失败 |
| `status` | `tinyint(1)` | NO | '0' | 审核状态0=未审核1=微信审核2=审核通过-1=审核未通过 |
| `is_show` | `tinyint(3)` | NO | '0' | 是否显示 |
| `is_mer_show` | `tinyint(3)` | NO | '0' | 商户是否显示 |
| `mark` | `varchar(512)` | YES | NULL | 备注 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |
| `sort` | `smallint(5)` | NO | '0' | 排序 |
| `is_del` | `tinyint(3)` | NO | '0' | — |
| `is_mer_del` | `tinyint(3)` | NO | '0' | 商户是否删除 |

## `qixi_broadcast_room`

> 直播间表

原表：`eb_broadcast_room`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `broadcast_room_id` | `int(10)` | NO | — | — |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `room_id` | `int(10)` | NO | '0' | 直播间 id |
| `name` | `varchar(32)` | NO | — | 直播间名字 |
| `cover_img` | `varchar(255)` | NO | — | 背景图 |
| `share_img` | `varchar(255)` | NO | — | 分享图 |
| `start_time` | `timestamp` | YES | NULL | 直播计划开始时间 |
| `end_time` | `timestamp` | YES | NULL | 直播计划结束时间 |
| `anchor_name` | `varchar(32)` | NO | — | 主播昵称 |
| `anchor_wechat` | `varchar(32)` | NO | — | 主播微信号 |
| `phone` | `varchar(32)` | NO | — | 主播手机号 |
| `type` | `tinyint(3)` | NO | '0' | 直播间类型 【1: 推流，0：手机直播】 |
| `screen_type` | `tinyint(3)` | NO | '1' | 横屏、竖屏 【1：横屏，0：竖屏】 |
| `close_like` | `tinyint(3)` | NO | '0' | 是否关闭点赞 |
| `close_goods` | `tinyint(3)` | NO | '0' | 是否关闭货架 |
| `close_comment` | `tinyint(3)` | NO | '0' | 是否关闭评论 |
| `close_share` | `tinyint(3)` | NO | '0' | 是否关闭分享 |
| `close_kf` | `tinyint(3)` | NO | '0' | 是否关闭客服 |
| `error_msg` | `varchar(255)` | YES | NULL | 未通过原因 |
| `status` | `tinyint(1)` | NO | '0' | 审核状态0=未审核1=微信审核2=审核通过-1=审核未通过 |
| `live_status` | `smallint(5)` | NO | '102' | 直播状态101：直播中，102：未开始，103已结束，104禁播，105：暂停，106：异常，107：已过期 |
| `mark` | `varchar(512)` | YES | NULL | 备注 |
| `replay_status` | `tinyint(3)` | YES | '0' | 回放状态 |
| `is_mer_show` | `tinyint(3)` | NO | '1' | 商户是否显示 |
| `is_show` | `tinyint(3)` | NO | '1' | 是否显示 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `star` | `smallint(5)` | NO | '1' | 推荐星级 |
| `sort` | `smallint(5)` | NO | '0' | 排序 |
| `is_del` | `tinyint(3)` | NO | '0' | 是否删除 |
| `is_mer_del` | `tinyint(3)` | NO | '0' | 商户是否删除 |
| `feeds_img` | `varchar(255)` | YES | NULL | 封面图 |
| `push_url` | `varchar(255)` | YES | NULL | 推流地址 |
| `assistant_id` | `varchar(255)` | YES | NULL | 小助手ID |
| `is_feeds_public` | `tinyint(3)` | YES | '0' | 是否开启官方收录，1 开启，0 关闭 |

## `qixi_broadcast_room_goods`

> 直播间导入商品表

原表：`eb_broadcast_room_goods`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `broadcast_room_id` | `int(10)` | NO | — | — |
| `broadcast_goods_id` | `int(10)` | NO | — | — |
| `on_sale` | `tinyint(4)` | YES | '1' | 商品上下架 |
