# 用户/资产/分销 — 表字段（`qixi_`）

## `qixi_member_interests`

原表：`eb_member_interests`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `interests_id` | `int(11)` | NO | — | — |
| `name` | `varchar(20)` | YES | NULL | 名称 |
| `info` | `varchar(200)` | YES | NULL | 介绍 |
| `brokerage_level` | `tinyint(3)` | YES | NULL | 关联等级 |
| `pic` | `varchar(128)` | YES | NULL | 图标 |
| `type` | `tinyint(4)` | YES | '0' | 类型1.免费会员 2.付费会员 |
| `link` | `varchar(500)` | YES | NULL | 跳转 链接 |
| `has_type` | `int(11)` | YES | NULL | 特权类型 |
| `value` | `varchar(500)` | YES | NULL | 特权值 |
| `on_pic` | `varchar(128)` | YES | NULL | — |
| `status` | `tinyint(4)` | YES | NULL | — |

## `qixi_user`

> 用户表

原表：`eb_user`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `uid` | `int(10)` | NO | — | 用户id |
| `wechat_user_id` | `int(10)` | NO | '0' | 微信用户 id |
| `account` | `varchar(32)` | NO | — | 用户账号 |
| `pwd` | `varchar(128)` | NO | — | 用户密码 |
| `real_name` | `varchar(25)` | NO | '' | 真实姓名 |
| `sex` | `tinyint(3)` | NO | '0' | 性别 |
| `birthday` | `date` | YES | NULL | 生日 |
| `card_id` | `varchar(20)` | NO | '' | 身份证号码 |
| `mark` | `varchar(255)` | NO | '' | 用户备注 |
| `label_id` | `varchar(64)` | YES | NULL | 用户标签 id |
| `group_id` | `int(10)` | NO | '0' | 用户分组id |
| `nickname` | `varchar(16)` | NO | — | 用户昵称 |
| `avatar` | `varchar(256)` | NO | — | 用户头像 |
| `phone` | `char(15)` | YES | NULL | 手机号码 |
| `addres` | `varchar(128)` | YES | NULL | 地址 |
| `cancel_time` | `timestamp` | YES | NULL | 注销时间 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `last_time` | `timestamp` | YES | NULL | 最后一次登录时间 |
| `last_ip` | `varchar(16)` | NO | — | 最后一次登录ip |
| `now_money` | `decimal(8,2)` | NO | '0.00' | 用户余额 |
| `brokerage_price` | `decimal(8,2)` | NO | '0.00' | 佣金金额 |
| `status` | `tinyint(1)` | NO | '1' | 1为正常，0为禁止 |
| `spread_uid` | `int(10)` | NO | '0' | 推广员id |
| `spread_time` | `timestamp` | YES | NULL | 推广员关联时间 |
| `spread_limit` | `timestamp` | YES | NULL | 推广员到期时间 |
| `brokerage_level` | `int(10)` | YES | '0' | 推广员等级 |
| `user_type` | `varchar(32)` | NO | — | 用户类型 |
| `promoter_time` | `timestamp` | YES | NULL | 成功推广时间 |
| `is_promoter` | `tinyint(3)` | NO | '0' | 是否为推广员 |
| `main_uid` | `int(10)` | YES | '0' | 主账号 |
| `pay_count` | `int(10)` | NO | '0' | 用户购买次数 |
| `pay_price` | `decimal(10,2)` | NO | '0.00' | 用户消费金额 |
| `spread_count` | `int(10)` | NO | '0' | 下级人数 |
| `spread_pay_count` | `int(10)` | YES | '0' | 下级订单数 |
| `spread_pay_price` | `decimal(10,2)` | YES | '0.00' | 下级订单金额 |
| `integral` | `int(10)` | YES | '0' | 积分 |
| `member_level` | `int(10)` | YES | '0' | 免费会员等级 |
| `member_value` | `int(10)` | YES | '0' | 免费会员成长值 |
| `count_start` | `int(10)` | YES | '0' | 用户获赞数 |
| `count_fans` | `int(10)` | YES | '0' | 用户粉丝数 |
| `count_content` | `int(10)` | YES | '0' | 用户内容数量 |
| `is_svip` | `tinyint(1)` | NO | '-1' | 是否为付费会员 -1未开通过 0到期 1体验卡 2 有效期 3 永久 |
| `svip_endtime` | `timestamp` | YES | NULL | 会员结束时间 |
| `svip_save_money` | `decimal(10,2)` | NO | '0.00' | 会员节省金额 |
| `promoter_switch` | `tinyint(1)` | NO | '1' | 分销资格 0无 1有 |

## `qixi_user_address`

> 用户地址表

原表：`eb_user_address`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `address_id` | `mediumint(8)` | NO | — | 用户地址id |
| `uid` | `int(10)` | NO | — | 用户id |
| `real_name` | `varchar(32)` | NO | '' | 收货人姓名 |
| `phone` | `varchar(16)` | NO | '' | 收货人电话 |
| `province` | `varchar(64)` | NO | '' | 收货人所在省 |
| `province_id` | `int(10)` | YES | '0' | 省 id |
| `city` | `varchar(64)` | NO | '' | 收货人所在市 |
| `city_id` | `int(11)` | NO | '0' | 城市id |
| `district` | `varchar(64)` | NO | '' | 收货人所在区 |
| `district_id` | `int(10)` | YES | '0' | 区域 id |
| `street` | `varchar(64)` | YES | NULL | 街/镇 |
| `street_id` | `int(10)` | YES | '0' | 街镇 id |
| `detail` | `varchar(256)` | NO | '' | 收货人详细地址 |
| `post_code` | `int(10)` | NO | — | 邮编 |
| `longitude` | `varchar(16)` | NO | '0' | 经度 |
| `latitude` | `varchar(16)` | NO | '0' | 纬度 |
| `is_default` | `tinyint(3)` | NO | '0' | 是否默认 |
| `is_del` | `tinyint(3)` | NO | '0' | 是否删除 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `tourist_unique_key` | `varchar(20)` | NO | '' | 游客唯一标识 |

## `qixi_user_bill`

> 用户账单表

原表：`eb_user_bill`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `bill_id` | `int(10)` | NO | — | 用户账单id |
| `uid` | `int(10)` | NO | '0' | 用户uid |
| `link_id` | `varchar(32)` | NO | '0' | 关联id |
| `pm` | `tinyint(3)` | NO | '0' | 0 = 支出 1 = 获得 |
| `title` | `varchar(64)` | NO | — | 账单标题 |
| `category` | `varchar(64)` | NO | — | 明细种类 |
| `type` | `varchar(64)` | NO | '' | 明细类型 |
| `number` | `decimal(11,2)` | NO | '0.00' | 明细数字 |
| `balance` | `decimal(11,2)` | NO | '0.00' | 剩余 |
| `mark` | `varchar(512)` | NO | — | 备注 |
| `mer_id` | `int(10)` | YES | '0' | 商户 id |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `status` | `tinyint(1)` | NO | '1' | 0 = 待确定 1 = 有效 -1 = 无效 |

## `qixi_user_brokerage`

原表：`eb_user_brokerage`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `user_brokerage_id` | `int(10)` | NO | — | — |
| `brokerage_level` | `tinyint(3)` | NO | — | 等级 |
| `brokerage_name` | `varchar(32)` | NO | — | vip 名称 |
| `brokerage_icon` | `varchar(128)` | NO | — | vip 图标 |
| `brokerage_rule` | `varchar(1500)` | NO | — | 升级规则 |
| `user_num` | `int(10)` | NO | '0' | vip 人数 |
| `extension_one` | `decimal(8,2)` | NO | '0.00' | 一级佣金 |
| `extension_two` | `decimal(8,2)` | NO | — | 二级佣金 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `type` | `tinyint(3)` | NO | '0' | 默认0分销会员等级，1 免费会员等级 |

## `qixi_user_extract`

> 用户提现表

原表：`eb_user_extract`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `extract_id` | `int(10)` | NO | — | — |
| `uid` | `int(10)` | NO | '0' | 用户 id |
| `extract_sn` | `varchar(255)` | YES | NULL | — |
| `real_name` | `varchar(64)` | YES | NULL | 姓名 |
| `extract_type` | `tinyint(1)` | YES | '0' | 0 银行卡 1 支付宝 2微信 3 零钱 |
| `bank_code` | `varchar(32)` | YES | '0' | 银行卡 |
| `bank_address` | `varchar(256)` | YES | '' | 开户地址 |
| `alipay_code` | `varchar(64)` | YES | '' | 支付宝账号 |
| `wechat` | `varchar(15)` | YES | NULL | 微信号 |
| `extract_pic` | `varchar(128)` | YES | NULL | 收款码 |
| `extract_price` | `decimal(8,2)` | YES | '0.00' | 提现金额 |
| `balance` | `decimal(8,2)` | YES | '0.00' | 余额 |
| `mark` | `varchar(512)` | YES | NULL | 管理员备注 |
| `admin_id` | `int(11)` | YES | '0' | 审核管理员 |
| `fail_msg` | `varchar(128)` | YES | NULL | 无效原因 |
| `status_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 无效时间 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `status` | `tinyint(4)` | YES | '0' | -1 未通过 0 审核中 1 已提现 |
| `bank_name` | `varchar(255)` | YES | NULL | 银行名称 |
| `wechat_status` | `varchar(50)` | NO | '' | 微信转账状态：ACCEPTED: 转账已受理;PROCESSING: 转账锁定资金中;WAIT_USER_CONFIRM: 待收款用户确认;TRANSFERING: 转账中;SUCCESS: 转账成功;FAIL: 转账失败;CANCELING: 撤销中;CANCELLED: 转账撤销完成 |
| `package_info` | `varchar(255)` | NO | '' | 跳转微信支付收款页的package信息 |
| `transfer_bill_no` | `varchar(100)` | NO | '' | 微信转账单号 |
| `wechat_app_id` | `varchar(100)` | NO | '' | 微信appid |
| `wechat_mch_id` | `varchar(100)` | NO | '' | 微信mchid |

## `qixi_user_fields`

原表：`eb_user_fields`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `bigint(20)` | NO | — | — |
| `uid` | `bigint(20)` | NO | '0' | 用户id |

## `qixi_user_group`

> 用户分组表

原表：`eb_user_group`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `group_id` | `smallint(5)` | NO | — | — |
| `group_name` | `varchar(64)` | NO | — | 用户分组名称 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |

## `qixi_user_history`

> 浏览记录表

原表：`eb_user_history`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `user_history_id` | `int(10)` | NO | — | — |
| `res_id` | `int(10)` | YES | NULL | 历史记录对象的ID |
| `res_type` | `int(11)` | YES | NULL | 历史记录类型 |
| `uid` | `int(11)` | YES | NULL | — |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `update_time` | `int(11)` | YES | NULL | — |

## `qixi_user_info`

原表：`eb_user_info`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `bigint(20)` | NO | — | — |
| `field` | `varchar(255)` | NO | '' | 字段 |
| `title` | `varchar(255)` | NO | '字段名' | 字段名 |
| `is_used` | `tinyint(4)` | NO | '0' | 是否使用 |
| `is_require` | `tinyint(4)` | NO | '0' | 是否必填 |
| `is_show` | `tinyint(4)` | NO | '0' | 是否在用户端展示 |
| `type` | `varchar(255)` | NO | '' | 信息格式 |
| `msg` | `varchar(255)` | NO | '' | 提示信息 |
| `content` | `varchar(255)` | YES | NULL | 配置内容 |
| `is_default` | `tinyint(4)` | NO | '0' | 是否系统默认字段 |
| `sort` | `int(11)` | NO | '0' | 排序 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `update_time` | `timestamp` | YES | NULL | — |

## `qixi_user_label`

> 用户标签表

原表：`eb_user_label`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `label_id` | `int(11)` | NO | — | — |
| `label_name` | `varchar(255)` | NO | '' | 标签名称 |
| `mer_id` | `int(10)` | NO | '0' | 商户 id |
| `type` | `tinyint(3)` | NO | '0' | 0=手动标签 1=自动标签 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |

## `qixi_user_merchant`

> 商户用户表

原表：`eb_user_merchant`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `user_merchant_id` | `int(10)` | NO | — | — |
| `uid` | `int(10)` | NO | — | 用户 id |
| `mer_id` | `int(10)` | NO | '0' | 商户 id |
| `first_pay_time` | `timestamp` | YES | NULL | 首次消费时间 |
| `last_pay_time` | `timestamp` | YES | NULL | 最后一次消费时间 |
| `pay_num` | `int(10)` | NO | '0' | 消费次数 |
| `pay_price` | `decimal(10,2)` | NO | '0.00' | 消费金额 |
| `label_id` | `varchar(256)` | YES | NULL | 用户标签 |
| `last_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 最后一次访问时间 |
| `status` | `tinyint(3)` | YES | '1' | 状态 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |

## `qixi_user_receipt`

> 用户发票信息

原表：`eb_user_receipt`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `user_receipt_id` | `int(10)` | NO | — | — |
| `receipt_type` | `tinyint(1)` | YES | '0' | 发票类型：1.普通发票，2.增值税发票 |
| `receipt_title` | `varchar(128)` | YES | '' | 发票抬头 |
| `receipt_title_type` | `varchar(255)` | YES | '0' | 发票抬头类型：1.个人，2.企业 |
| `duty_paragraph` | `varchar(255)` | YES | '' | 税号 |
| `email` | `varchar(255)` | YES | '' | 邮箱 |
| `bank_name` | `varchar(255)` | YES | '' | 开户行 |
| `bank_code` | `varchar(255)` | YES | '0' | 银行账号 |
| `address` | `varchar(255)` | YES | '' | 企业地址 |
| `tel` | `varchar(255)` | YES | '' | 企业电话 |
| `is_default` | `tinyint(4)` | NO | '0' | 是否默认 |
| `uid` | `int(11)` | NO | '0' | 用户ID |
| `is_del` | `tinyint(1)` | YES | '0' | — |
| `create_time` | `timestamp` | YES | NULL | — |

## `qixi_user_recharge`

> 用户充值表

原表：`eb_user_recharge`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `recharge_id` | `int(10)` | NO | — | — |
| `uid` | `int(10)` | NO | — | 充值用户UID |
| `order_id` | `varchar(32)` | NO | — | 订单号 |
| `price` | `decimal(8,2)` | NO | '0.00' | 充值金额 |
| `give_price` | `decimal(8,2)` | NO | '0.00' | 购买赠送金额 |
| `recharge_type` | `varchar(32)` | NO | — | 充值类型 |
| `paid` | `tinyint(3)` | NO | '0' | 是否充值 |
| `pay_time` | `timestamp` | YES | NULL | 充值支付时间 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 充值时间 |
| `refund_price` | `decimal(10,2)` | YES | '0.00' | 退款金额 |

## `qixi_user_relation`

> 用户记录表

原表：`eb_user_relation`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `uid` | `int(10)` | NO | — | 用户ID |
| `type_id` | `int(10)` | NO | — | 类型的 id |
| `type` | `tinyint(4)` | NO | '0' | 关联类型(0 普通商品、1秒杀2、预售3、助力4、拼团、10 = 店铺、12=购买过) |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |

## `qixi_user_sign`

> 签到记录表

原表：`eb_user_sign`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `sign_id` | `int(11)` | NO | — | — |
| `uid` | `int(11)` | NO | '0' | 用户uid |
| `title` | `varchar(255)` | NO | '' | 签到说明 |
| `number` | `int(11)` | NO | '0' | 获得积分 |
| `integral` | `int(11)` | NO | '0' | 剩余积分 |
| `sign_num` | `int(10)` | NO | '0' | 连续签到天数 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |

## `qixi_user_spread_log`

原表：`eb_user_spread_log`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `user_spread_log_id` | `int(10)` | NO | — | — |
| `uid` | `int(10)` | NO | — | uid |
| `old_spread_uid` | `int(10)` | NO | — | 原来的推荐人uid |
| `spread_uid` | `int(10)` | NO | — | 新的推荐人 uid |
| `admin_id` | `int(10)` | NO | — | 修改的管理员 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |

## `qixi_user_visit`

> 商品浏览分析表

原表：`eb_user_visit`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `user_visit_id` | `int(11)` | NO | — | — |
| `uid` | `int(11)` | YES | NULL | 用户ID |
| `type` | `varchar(32)` | NO | — | 记录类型 |
| `type_id` | `int(11)` | NO | '0' | 商品ID |
| `content` | `varchar(255)` | YES | NULL | 备注描述 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
