# 其他 — 表字段（`qixi_`）

## `qixi_store_printer`

原表：`eb_store_printer`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `printer_id` | `bigint(20)` | NO | — | — |
| `printer_name` | `varchar(50)` | NO | '' | 名称 |
| `printer_appkey` | `varchar(50)` | NO | '' | 打印机的应用ID |
| `printer_terminal` | `varchar(50)` | NO | '' | 打印机终端号 |
| `printer_appid` | `varchar(50)` | NO | '' | 打印机应用用户ID |
| `printer_secret` | `varchar(50)` | NO | '' | 打印机应用密匙 |
| `status` | `tinyint(4)` | NO | '0' | — |
| `mer_id` | `int(11)` | NO | '0' | — |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |
| `type` | `tinyint(4)` | YES | '0' | 0 易联云 1 飞鹅云 |
| `times` | `int(11)` | NO | '1' | 打印联数 |
| `print_content` | `varchar(2000)` | NO | '' | 打印内容 |
| `print_type` | `int(11)` | NO | '1' | 打印时机1支付后，2下单后 |

## `qixi_wechat_qrcode`

> 微信二维码管理表

原表：`eb_wechat_qrcode`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `wechat_qrcode_id` | `int(10)` | NO | — | 微信二维码ID |
| `third_type` | `varchar(32)` | NO | — | 二维码类型 |
| `third_id` | `int(10)` | NO | — | 类型id |
| `ticket` | `varchar(255)` | NO | — | 二维码参数 |
| `expire_seconds` | `int(10)` | NO | '0' | 二维码有效时间 |
| `status` | `tinyint(3)` | NO | '1' | 状态 |
| `url` | `varchar(255)` | NO | — | 微信访问url |
| `qrcode_url` | `varchar(255)` | NO | — | 微信二维码url |
| `scan` | `int(10)` | NO | '0' | 被扫的次数 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |

## `qixi_wechat_reply`

> 微信关键字回复表

原表：`eb_wechat_reply`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `wechat_reply_id` | `mediumint(8)` | NO | — | 微信关键字回复id |
| `key` | `varchar(64)` | NO | — | 关键字 |
| `type` | `varchar(32)` | NO | — | 回复类型 |
| `data` | `text` | NO | — | 回复数据 |
| `status` | `tinyint(3)` | NO | '1' | 0=不可用  1 =可用 |
| `hidden` | `tinyint(3)` | NO | '0' | 是否显示 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | — |

## `qixi_wechat_user`

> 微信用户表

原表：`eb_wechat_user`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `wechat_user_id` | `int(10)` | NO | — | 微信用户id |
| `unionid` | `varchar(60)` | YES | NULL | 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段 |
| `openid` | `varchar(30)` | YES | NULL | 用户的标识，对当前公众号唯一 |
| `routine_openid` | `varchar(32)` | YES | NULL | 小程序唯一身份ID |
| `nickname` | `varchar(64)` | NO | — | 用户的昵称 |
| `headimgurl` | `varchar(256)` | NO | — | 用户头像 |
| `sex` | `tinyint(3)` | NO | '0' | 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知 |
| `city` | `varchar(32)` | NO | — | 用户所在城市 |
| `language` | `varchar(32)` | NO | — | 用户的语言，简体中文为zh_CN |
| `province` | `varchar(32)` | NO | — | 用户所在省份 |
| `country` | `varchar(32)` | NO | — | 用户所在国家 |
| `remark` | `varchar(256)` | YES | NULL | 公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注 |
| `groupid` | `smallint(5)` | YES | '0' | 用户所在的分组ID（兼容旧的用户分组接口） |
| `tagid_list` | `varchar(256)` | YES | NULL | 用户被打上的标签ID列表 |
| `subscribe` | `tinyint(3)` | YES | '0' | 用户是否订阅该公众号标识 |
| `subscribe_time` | `int(10)` | YES | NULL | 关注公众号时间 |
| `session_key` | `varchar(32)` | YES | NULL | 小程序用户会话密匙 |
| `user_type` | `varchar(32)` | YES | 'wechat' | 用户类型 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 创建时间 |
