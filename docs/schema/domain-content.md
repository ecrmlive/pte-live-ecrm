# 内容/社区/圈子 — 表字段（`qixi_`）

## `qixi_article`

> 文章管理表

原表：`eb_article`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `article_id` | `int(10)` | NO | — | 文章管理ID |
| `cid` | `int(10)` | YES | '0' | 分类id |
| `title` | `varchar(64)` | NO | — | 文章标题 |
| `author` | `varchar(32)` | YES | NULL | 文章作者 |
| `image_input` | `varchar(128)` | NO | — | 文章图片 |
| `synopsis` | `varchar(128)` | YES | NULL | 文章简介 |
| `visit` | `varchar(255)` | YES | NULL | 浏览次数 |
| `sort` | `int(10)` | YES | '0' | 排序 |
| `url` | `varchar(128)` | YES | NULL | 原文链接 |
| `admin_id` | `int(10)` | NO | '0' | 管理员id |
| `mer_id` | `int(10)` | YES | '0' | 商户id |
| `is_hot` | `tinyint(3)` | YES | '0' | 是否热门(小程序) |
| `is_banner` | `tinyint(3)` | YES | '0' | 是否轮播图(小程序) |
| `status` | `tinyint(3)` | YES | NULL | 状态 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `wechat_news_id` | `int(11)` | YES | '0' | 微信图文id |

## `qixi_article_category`

> 文章分类表

原表：`eb_article_category`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `article_category_id` | `int(10)` | NO | — | 文章分类id |
| `pid` | `int(11)` | NO | '0' | 父级ID |
| `mer_id` | `int(10)` | NO | '0' | 商户 id |
| `title` | `varchar(32)` | NO | — | 文章分类标题 |
| `info` | `varchar(255)` | YES | NULL | 文章分类简介 |
| `image` | `varchar(128)` | NO | — | 文章分类图片 |
| `status` | `tinyint(3)` | NO | — | 状态 |
| `sort` | `int(10)` | NO | '0' | 排序 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |

## `qixi_article_content`

> 文章内容表

原表：`eb_article_content`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `article_content_id` | `int(10)` | NO | — | 文章id |
| `content` | `text` | NO | — | 文章内容 |

## `qixi_circle`

原表：`eb_circle`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `circle_id` | `int(10)` | NO | — | 商圈id |
| `pid` | `int(10)` | NO | '0' | 上级商圈id |
| `path` | `varchar(50)` | NO | '' | 路径 |
| `name` | `varchar(64)` | NO | '' | 商圈名称 |
| `circle_agent_id` | `int(10)` | NO | '0' | 商圈代理id |
| `commission_type` | `tinyint(3)` | NO | '0' | 商圈提成类型(0:按默认设置，1:单独设置) |
| `commission_rate` | `decimal(8,2)` | NO | '0.00' | 商圈提成比例(0~100%)：commission_type为0时取系统设置等级提成比例，为1时自定义设置) |
| `level` | `tinyint(3)` | NO | '0' | 等级:0一级商圈 1二级商圈 2三级商圈 |
| `remark` | `varchar(255)` | NO | '' | 说明信息 |
| `sort` | `int(10)` | NO | '0' | 排序(数字越大越靠前) |
| `status` | `tinyint(3)` | NO | '1' | 状态：0禁用 1启用 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `update_time` | `datetime` | NO | CURRENT_TIMESTAMP | 修改时间 |
| `type` | `tinyint(1)` | NO | '0' | 类型：0区域，1商户 |
| `role_id` | `int(10)` | NO | '0' | 角色权限id |
| `business_store_category` | `int(10)` | NO | '0' | 商户店铺分类 |
| `business_store_type` | `int(10)` | NO | '0' | 商户店铺类型 |

## `qixi_circle_agent`

> 商圈代理表

原表：`eb_circle_agent`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `circle_agent_id` | `int(10)` | NO | — | 商圈代理id |
| `uid` | `int(10)` | NO | '0' | 用户id |
| `name` | `varchar(64)` | NO | '' | 代理名称 |
| `phone` | `varchar(16)` | NO | '' | 联系电话 |
| `qualification` | `text` | YES | — | 身份资质 |
| `remark` | `varchar(255)` | NO | '' | 说明信息 |
| `extend` | `text` | YES | — | 扩展信息 |
| `audit_admin_id` | `int(10)` | NO | '0' | 审核管理员 |
| `audit_reason` | `varchar(255)` | NO | '' | 审核拒绝原因 |
| `audit_time` | `datetime` | YES | NULL | 审核时间 |
| `status` | `tinyint(3)` | NO | '0' | 状态：0待审核 1审核通过 -1审核拒绝 -2撤销 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `payment_method` | `tinyint(3)` | NO | '0' | 结算方式:0:银行卡，1:微信，2:支付宝 |
| `payment_name` | `varchar(16)` | NO | '' | 结算名称 |
| `payment_account` | `varchar(30)` | NO | '' | 结算账号 |
| `payment_bank` | `varchar(30)` | NO | '' | 开户行 |
| `payment_qr_img` | `varchar(200)` | NO | '' | 收款二维码图片 |
| `balance` | `decimal(8,2)` | NO | '0.00' | 佣金余额 |
| `update_time` | `datetime` | NO | CURRENT_TIMESTAMP | 修改时间 |
| `type` | `tinyint(1)` | NO | '0' | 类型：0区域，1商户 |
| `business_name` | `varchar(64)` | NO | '' | 商户名称 |
| `business_store_category` | `int(10)` | NO | '0' | 商户店铺分类 |
| `business_store_type` | `int(10)` | NO | '0' | 商户店铺类型 |

## `qixi_circle_brokerage_checkout`

> 商圈佣金结算表

原表：`eb_circle_brokerage_checkout`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `checkout_id` | `int(10)` | NO | — | — |
| `agent_id` | `int(10)` | NO | — | 商圈代理id |
| `agent_name` | `varchar(64)` | NO | '' | 商圈代理名称 |
| `agent_phone` | `varchar(64)` | NO | '' | 商圈代理电话 |
| `withdrawal_sn` | `varchar(32)` | NO | '' | 提现流水号 |
| `withdrawal_amount` | `decimal(8,2)` | NO | '0.00' | 提现金额 |
| `withdrawal_type` | `int(10)` | NO | '0' | 提现方式 0:银行卡，1:微信，2:支付宝 |
| `transfer_voucher` | `text` | YES | — | 转账凭证 |
| `remark` | `varchar(200)` | NO | '' | 备注 |
| `platform_remark` | `varchar(200)` | NO | '' | 平台备注 |
| `transfer_remark` | `varchar(200)` | NO | '' | 转账备注 |
| `audit_time` | `datetime` | YES | NULL | 审核时间 |
| `audit_status` | `tinyint(3)` | NO | '0' | 审核状态：0待审核 1审核通过 -1审核拒绝 -2撤销 |
| `audit_admin_id` | `int(10)` | NO | '0' | 审核管理员 |
| `audit_reason` | `varchar(255)` | NO | '' | 审核拒绝原因 |
| `status` | `tinyint(3)` | NO | '0' | 状态：0未到帐 1已到帐 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `update_time` | `datetime` | NO | CURRENT_TIMESTAMP | 修改时间 |
| `withdrawal_qr_img` | `varchar(200)` | NO | '' | 收款二维码图片 |
| `withdrawal_name` | `varchar(16)` | NO | '' | 结算名称 |
| `withdrawal_account` | `varchar(30)` | NO | '' | 结算账号 |

## `qixi_community`

> 社区图文表信息

原表：`eb_community`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `community_id` | `int(10)` | NO | — | — |
| `title` | `varchar(255)` | YES | NULL | 标题 |
| `image` | `varchar(1000)` | YES | NULL | 图片 |
| `category_id` | `int(10)` | YES | '0' | — |
| `topic_id` | `int(10)` | YES | '0' | 话题 |
| `uid` | `int(10)` | YES | '0' | 用户 |
| `count_start` | `int(10)` | YES | '0' | 点赞数 |
| `count_reply` | `int(10)` | YES | '0' | 评论数 |
| `count_share` | `int(10)` | YES | '0' | 分享数 |
| `status` | `tinyint(4)` | YES | '0' | 审核状态 |
| `is_show` | `tinyint(4)` | YES | '0' | 显示状态 |
| `start` | `tinyint(1)` | YES | '1' | 星级排序 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `is_del` | `tinyint(1)` | YES | '0' | — |
| `content` | `varchar(1000)` | YES | NULL | — |
| `refusal` | `varchar(255)` | YES | NULL | 拒绝理由 |
| `is_hot` | `tinyint(4)` | YES | '0' | 是否推荐 |
| `order_id` | `int(10)` | YES | '0' | 关联订单ID |
| `is_type` | `tinyint(1)` | NO | '1' | 1 图文 2 视频 |
| `video_link` | `varchar(255)` | YES | NULL | 视频链接 |
| `pv` | `int(11)` | YES | '0' | 浏览量 |
| `update_time` | `timestamp` | YES | NULL | 更新时间 |
| `status_time` | `timestamp` | YES | NULL | 审核时间 |
| `mer_id` | `int(10)` | YES | '0' | 商户ID |

## `qixi_community_category`

> 社区分类

原表：`eb_community_category`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `category_id` | `bigint(20)` | NO | — | — |
| `cate_name` | `varchar(50)` | YES | NULL | 分类名 |
| `pid` | `int(11)` | YES | NULL | 父级ID |
| `path` | `varchar(255)` | YES | '/' | 路径  |
| `is_show` | `tinyint(4)` | YES | '1' | 状态 |
| `level` | `int(11)` | YES | '0' | 等级 |
| `sort` | `int(11)` | YES | NULL | — |

## `qixi_community_reply`

> 社区评论

原表：`eb_community_reply`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `reply_id` | `bigint(20)` | NO | — | — |
| `content` | `varchar(255)` | YES | NULL | 评论内容 |
| `pid` | `int(10)` | YES | '0' | 回复id |
| `uid` | `int(10)` | YES | '0' | 发言人 |
| `re_uid` | `int(10)` | YES | '0' | 回复人 |
| `count_start` | `int(10)` | YES | '0' | 点赞数 |
| `count_reply` | `int(10)` | YES | '0' | 评论数 |
| `status` | `tinyint(4)` | YES | '1' | 状态  |
| `community_id` | `int(10)` | YES | '0' | 文章id |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `is_del` | `tinyint(4)` | YES | '0' | — |
| `refusal` | `varchar(255)` | YES | NULL | 拒绝原因 |

## `qixi_community_topic`

> 社区话题

原表：`eb_community_topic`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `topic_id` | `bigint(20)` | NO | — | — |
| `topic_name` | `varchar(100)` | YES | NULL | 话题 |
| `status` | `tinyint(4)` | YES | '1' | 状态 |
| `is_hot` | `tinyint(4)` | YES | '0' | 推荐 |
| `category_id` | `int(10)` | YES | '0' | 分类id |
| `is_del` | `tinyint(4)` | YES | '0' | — |
| `pic` | `varchar(128)` | YES | NULL | 图标 |
| `count_use` | `int(10)` | YES | '0' | 使用次数 |
| `count_view` | `int(10)` | YES | '0' | 浏览量 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `sort` | `int(10)` | YES | '0' | — |

## `qixi_feedback`

> 用户反馈表

原表：`eb_feedback`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `feedback_id` | `int(10)` | NO | — | — |
| `uid` | `int(10)` | NO | '0' | — |
| `type` | `varchar(255)` | NO | — | — |
| `content` | `varchar(512)` | NO | — | — |
| `images` | `text` | YES | — | 反馈图片 |
| `realname` | `varchar(24)` | NO | — | 姓名 |
| `contact` | `varchar(32)` | NO | — | 联系方式 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |
| `status` | `tinyint(1)` | YES | '0' | 状态 |
| `reply` | `varchar(255)` | YES | NULL | 回复，最终给用户的回复内容 |
| `remake` | `varchar(255)` | YES | NULL | 备注，后台人员自己查看用 |
| `is_del` | `tinyint(1)` | NO | '0' | — |
| `update_time` | `timestamp` | YES | NULL | 回复时间 |

## `qixi_feedback_category`

> 用户反馈分类表

原表：`eb_feedback_category`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `feedback_category_id` | `mediumint(9)` | NO | — | 商品分类表ID |
| `pid` | `mediumint(9)` | NO | — | 父id |
| `cate_name` | `varchar(100)` | NO | — | 分类名称 |
| `path` | `varchar(255)` | NO | '' | 路径 |
| `sort` | `mediumint(9)` | NO | — | 排序 |
| `pic` | `varchar(128)` | NO | '' | 图标 |
| `is_show` | `tinyint(1)` | NO | '1' | 是否显示 |
| `level` | `int(10)` | NO | '0' | 等级 |
| `mer_id` | `int(10)` | YES | '0' | 商户id |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 添加时间 |

## `qixi_wechat_news`

> 图文消息管理表

原表：`eb_wechat_news`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `wechat_news_id` | `int(10)` | NO | — | 图文消息管理ID |
| `mer_id` | `int(11)` | YES | '0' | 商户id |
| `status` | `tinyint(3)` | NO | '1' | 状态 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
