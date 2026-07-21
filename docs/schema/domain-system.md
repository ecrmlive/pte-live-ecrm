# 系统/权限/配置/DIY/素材 — 表字段（`qixi_`）

## `qixi_cache`

> 微信缓存表

原表：`eb_cache`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `key` | `varchar(32)` | NO | — | — |
| `expire_time` | `int(11)` | NO | '0' | 0=永久 |
| `result` | `longtext` | NO | — | 缓存数据 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 缓存时间 |

## `qixi_city_area`

> 省市区县数据

原表：`eb_city_area`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `int(10)` | NO | — | — |
| `path` | `varchar(128)` | NO | '/' | 省市级别 |
| `parent_id` | `int(11)` | NO | '0' | 父级id |
| `type` | `varchar(32)` | NO | — | 类型 |
| `name` | `varchar(100)` | NO | '' | 名称 |
| `level` | `tinyint(3)` | NO | '0' | 级别 |
| `code` | `varchar(100)` | NO | '' | 城市编码 |
| `snum` | `int(10)` | NO | '0' | 子级个数 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |

## `qixi_diy`

原表：`eb_diy`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `int(11)` | NO | — | — |
| `version` | `varchar(255)` | NO | '' | 版本号 |
| `name` | `varchar(255)` | NO | '' | 页面名称 |
| `title` | `varchar(100)` | NO | '' | 网站标题 |
| `cover_image` | `varchar(255)` | NO | '' | 封面图 |
| `template_name` | `varchar(255)` | NO | '' | 模板名称 |
| `default_value` | `longtext` | YES | — | 默认数据 |
| `add_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `update_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 更新时间 |
| `status` | `tinyint(1)` | NO | '0' | 是否使用 |
| `type` | `tinyint(1)` | NO | '0' | 页面类型 |
| `is_show` | `tinyint(1)` | NO | '0' | 显示首页 |
| `is_bg_color` | `tinyint(1)` | NO | '0' | 颜色是否选中 |
| `is_bg_pic` | `tinyint(1)` | NO | '0' | 背景图是否选中 |
| `is_diy` | `tinyint(1)` | NO | '0' | 是否是diy数据 |
| `color_picker` | `varchar(50)` | NO | '' | 背景颜色 |
| `bg_pic` | `varchar(256)` | NO | '' | 背景图 |
| `bg_tab_val` | `tinyint(1)` | NO | '0' | 背景图图片样式 |
| `is_del` | `tinyint(1)` | NO | '0' | 是否删除 |
| `mer_id` | `int(11)` | NO | '0' | 商户ID |
| `is_default` | `tinyint(1)` | NO | '0' | 默认模板(1.平台默认 2.商户默认） |
| `scope_type` | `tinyint(1)` | YES | '4' | 适用范围类型：0.全部店铺、1. 指定店铺、2. 指定商户分类、3. 指定店铺类型、4. 指定商户类别 |
| `value` | `longtext` | YES | — | 页面数据 |

## `qixi_excel`

> 导出文件记录表

原表：`eb_excel`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `excel_id` | `int(10)` | NO | — | id |
| `name` | `varchar(255)` | YES | NULL | 文件名 |
| `status` | `int(11)` | YES | '0' | 0.默认，1.完成，2.失败 |
| `type` | `varchar(255)` | YES | NULL | 类型 |
| `path` | `varchar(255)` | YES | NULL | 文件路径 |
| `mer_id` | `int(11)` | YES | '0' | 商户id |
| `admin_id` | `int(11)` | YES | NULL | 操作者id |
| `is_del` | `int(11)` | YES | '0' | 是否删除 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `message` | `varchar(255)` | YES | NULL | — |

## `qixi_extend`

原表：`eb_extend`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `extend_id` | `int(10)` | NO | — | 商户 id |
| `link_id` | `int(10)` | NO | — | 关联字段 |
| `mer_id` | `int(10)` | NO | '0' | 商户 id |
| `extend_type` | `varchar(32)` | NO | — | 扩展字段 |
| `extend_value` | `varchar(255)` | NO | — | 扩展值 |
| `update_time` | `datetime` | NO | — | 更新时间 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 创建时间 |

## `qixi_label_rule`

> 自定标签规则

原表：`eb_label_rule`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `label_rule_id` | `int(10)` | NO | — | — |
| `mer_id` | `int(10)` | NO | '0' | 商户 id |
| `label_id` | `int(10)` | NO | '0' | 标签 id |
| `type` | `tinyint(3)` | YES | '0' | 0=订单数 1=订单金额 |
| `min` | `decimal(8,2)` | NO | '0.00' | 最小值 |
| `max` | `decimal(8,2)` | NO | '0.00' |  最大值 |
| `user_num` | `int(10)` | NO | '0' | 用户数 |
| `update_time` | `timestamp` | YES | NULL | 更新时间 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |

## `qixi_open_auth`

原表：`eb_open_auth`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `int(11)` | NO | — | — |
| `title` | `varchar(50)` | YES | NULL | 标题 |
| `access_key` | `varchar(50)` | YES | NULL | — |
| `secret_key` | `varchar(255)` | YES | NULL | — |
| `status` | `tinyint(4)` | YES | NULL | 状态 |
| `mark` | `varchar(255)` | YES | NULL | 备注 |
| `mer_id` | `int(11)` | YES | NULL | 商户ID |
| `auth` | `varchar(255)` | YES | NULL | 权限 |
| `sort` | `int(11)` | YES | NULL | — |
| `is_del` | `tinyint(4)` | YES | '0' | 是否删除 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 创建时间 |
| `update_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 更新时间 |
| `delete_time` | `timestamp` | YES | NULL | 删除时间 |
| `last_ip` | `varchar(50)` | YES | NULL | 最后登录的IP |
| `last_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 最后登录的时间 |

## `qixi_operate_log`

原表：`eb_operate_log`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `operate_log_id` | `int(11)` | NO | — | id |
| `mer_id` | `int(11)` | NO | '0' | 商户id |
| `title` | `varchar(255)` | YES | NULL | 标题 |
| `relevance_id` | `int(11)` | NO | '0' | 关联id |
| `relevance_title` | `varchar(255)` | YES | NULL | 关联标题 |
| `relevance_type` | `varchar(255)` | NO | — | 关联类型 |
| `type` | `enum('1','2')` | NO | '1' | 1|平台 2|商户 |
| `category` | `varchar(255)` | NO | — | 类别 |
| `action` | `varchar(255)` | NO | '' | 操作类型 |
| `operator_role_id` | `int(11)` | YES | NULL | 操作角色id |
| `operator_role_nickname` | `varchar(255)` | YES | NULL | 操作角色昵称 |
| `operator_uid` | `varchar(255)` | NO | — | 操作用户id |
| `operator_nickname` | `varchar(255)` | YES | NULL | 操作用户昵称 |
| `mark` | `varchar(2000)` | YES | NULL | 备注 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 添加时间 |

## `qixi_page_category`

> 页面链接分类

原表：`eb_page_category`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `int(10)` | NO | — | — |
| `pid` | `int(10)` | NO | '0' | 父类id |
| `type` | `varchar(50)` | NO | 'link' | 类型:link、special、product、product_category、custom |
| `name` | `varchar(50)` | NO | '' | 分类名称 |
| `sort` | `smallint(5)` | NO | '0' | 排序 |
| `status` | `tinyint(1)` | NO | '1' | 状态 |
| `add_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `level` | `tinyint(1)` | NO | '0' | — |
| `is_mer` | `int(10)` | NO | '0' | — |

## `qixi_page_link`

> 页面链接

原表：`eb_page_link`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `int(10)` | NO | — | — |
| `cate_id` | `int(10)` | NO | '0' | 分类id |
| `type` | `tinyint(1)` | NO | '1' | 分组1:基础2:分销3:个人中心 |
| `name` | `varchar(50)` | NO | '' | 页面名称 |
| `url` | `varchar(255)` | NO | '' | 页面链接 |
| `param` | `varchar(255)` | NO | '' | 参数 |
| `example` | `varchar(255)` | NO | '' | 事例 |
| `status` | `tinyint(1)` | NO | '1' | 状态 |
| `sort` | `smallint(5)` | NO | '0' | 排序 |
| `add_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `is_mer` | `tinyint(1)` | NO | '0' | 1是商户的链接 |

## `qixi_record`

原表：`eb_record`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `bigint(20)` | NO | — | — |
| `type` | `varchar(20)` | NO | '' | — |
| `uid` | `int(11)` | NO | '0' | — |
| `link_id` | `int(11)` | NO | '0' | — |
| `num` | `int(11)` | NO | '0' | — |
| `title` | `varchar(100)` | NO | '' | — |

## `qixi_relevance`

原表：`eb_relevance`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `relevance_id` | `int(10)` | NO | — | — |
| `left_id` | `int(10)` | NO | — | — |
| `right_id` | `int(10)` | NO | — | — |
| `type` | `varchar(32)` | NO | '' | — |

## `qixi_routine_qrcode`

> 小程序二维码管理表

原表：`eb_routine_qrcode`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `routine_qrcode_id` | `int(10)` | NO | — | 微信二维码ID |
| `third_type` | `varchar(32)` | NO | — | 二维码类型 spread(用户推广) product_spread(商品推广) |
| `third_id` | `int(10)` | NO | — | 用户id |
| `status` | `tinyint(3)` | NO | '1' | 状态 0不可用 1可用 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `page` | `varchar(255)` | YES | NULL | 小程序页面路径带参数 |
| `qrcode_url` | `varchar(255)` | YES | NULL | 小程序二维码路径 |
| `url_time` | `timestamp` | YES | NULL | 二维码添加时间 |

## `qixi_sms_record`

> 短信发送记录表

原表：`eb_sms_record`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `sms_record_id` | `int(10)` | NO | — | 短信发送记录编号 |
| `uid` | `varchar(255)` | NO | — | 短信平台账号 |
| `phone` | `char(11)` | NO | — | 接受短信的手机号 |
| `content` | `text` | NO | — | 短信内容 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 发送短信时间 |
| `ip` | `varchar(16)` | NO | '' | 添加记录ip |
| `template` | `varchar(255)` | NO | — | 短信模板ID |
| `resultcode` | `int(10)` | YES | NULL | 状态码 100=成功,130=失败,131=空号,132=停机,133=关机,134=无状态 |
| `record_id` | `int(10)` | NO | — | 发送记录id |

## `qixi_staffs`

原表：`eb_staffs`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `staffs_id` | `int(11)` | NO | — | 员工id |
| `mer_id` | `int(11)` | NO | '0' | 商户id |
| `uid` | `int(11)` | NO | — | 关联用户uid |
| `photo` | `varchar(250)` | NO | — | 证件照 |
| `name` | `varchar(50)` | NO | — | 姓名 |
| `phone` | `varchar(18)` | YES | '' | 电话 |
| `remark` | `varchar(250)` | YES | '' | 备注 |
| `status` | `tinyint(3)` | NO | '0' | 状态：0、关闭；1、开启； |
| `sort` | `int(10)` | NO | '0' | 排序 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `delete_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 删除时间 |

## `qixi_system_admin`

> 后台管理员表

原表：`eb_system_admin`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `admin_id` | `int(10)` | NO | — | 后台管理员表ID |
| `account` | `varchar(32)` | NO | — | 后台管理员账号 |
| `pwd` | `varchar(64)` | NO | — | 后台管理员密码 |
| `real_name` | `varchar(16)` | NO | — | 后台管理员姓名 |
| `phone` | `varchar(12)` | YES | NULL | 联系电话 |
| `roles` | `varchar(128)` | NO | — | 后台管理员权限(role_id), 多个逗号分隔 |
| `last_ip` | `varchar(16)` | YES | NULL | 后台管理员最后一次登录ip |
| `last_time` | `timestamp` | YES | NULL | 后台管理员最后一次登录时间 |
| `login_count` | `int(10)` | NO | '0' | 登录次数 |
| `status` | `tinyint(3)` | NO | '1' | 后台管理员状态 1有效0无效 |
| `level` | `tinyint(3)` | NO | '1' | — |
| `is_del` | `tinyint(3)` | NO | '0' | — |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 后台管理员添加时间 |
| `update_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 后台管理员编辑时间 |
| `region_ids` | `varchar(200)` | NO | '' | 分组ID |
| `is_agent` | `tinyint(1)` | NO | '0' | 是否为区域代理：0否1是 |
| `circle_agent_id` | `int(10)` | NO | '0' | 商圈代理id |

## `qixi_system_attachment`

> 附件管理表

原表：`eb_system_attachment`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `attachment_id` | `int(11)` | NO | — | — |
| `attachment_category_id` | `int(10)` | NO | '0' | 分类ID 0编辑器,1产品图片,2拼团图片,3砍价图片,4秒杀图片,5文章图片,6组合数据图 |
| `attachment_name` | `varchar(100)` | NO | — | 附件名称 |
| `attachment_src` | `varchar(200)` | NO | — | 附件路径 |
| `upload_type` | `tinyint(3)` | NO | '1' | 图片上传类型 1本地 2七牛云 3OSS 4COS  |
| `user_type` | `int(11)` | NO | '0' | 图片上传模块类型 0总后台后台  >0商户后台  -1用户生成 |
| `user_id` | `int(10)` | NO | '0' | 上传用户的 id |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 上传时间 |
| `attachment_type` | `tinyint(1)` | NO | '0' | 素材类型 |

## `qixi_system_attachment_category`

> 附件分类表

原表：`eb_system_attachment_category`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `attachment_category_id` | `int(10)` | NO | — | — |
| `pid` | `int(10)` | NO | '0' | 父级ID |
| `path` | `varchar(512)` | NO | '' | 路径 |
| `attachment_category_name` | `varchar(32)` | NO | — | 分类名称 |
| `attachment_category_enname` | `varchar(16)` | NO | — | 分类目录 |
| `sort` | `smallint(5)` | NO | '0' | 排序 |
| `mer_id` | `int(10)` | NO | '0' | 商户 id |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |

## `qixi_system_config`

> 配置表

原表：`eb_system_config`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `config_id` | `int(10)` | NO | — | 配置id |
| `config_classify_id` | `int(10)` | NO | — | 配置分类id |
| `config_name` | `varchar(64)` | NO | — | 字段名称 |
| `config_key` | `varchar(64)` | NO | — | 字段 key |
| `config_type` | `varchar(20)` | NO | 'input' | 配置类型 |
| `config_rule` | `varchar(255)` | YES | NULL | 规则 |
| `config_props` | `varchar(255)` | YES | '' | 配置 |
| `required` | `tinyint(3)` | NO | '0' | 必填 |
| `info` | `varchar(255)` | YES | '' | 配置说明 |
| `sort` | `smallint(5)` | NO | '0' | 排序 |
| `user_type` | `tinyint(3)` | NO | '0' | 0=总后台配置 1=商户后台配置 |
| `status` | `tinyint(3)` | NO | '0' | 是否显示 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |
| `linked_status` | `tinyint(3)` | YES | '0' | 是否开启联动显示 0/1 |
| `linked_id` | `int(11)` | YES | '0' | 联动显示的id 信息 |
| `linked_value` | `int(11)` | YES | '0' | 联动显示的值 |

## `qixi_system_config_classify`

> 配置分类表

原表：`eb_system_config_classify`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `config_classify_id` | `int(10)` | NO | — | 配置分类id |
| `pid` | `int(11)` | YES | '0' | 父级ID |
| `classify_name` | `varchar(255)` | NO | — | 配置分类名称 |
| `classify_key` | `varchar(255)` | NO | — | 配置分类英文名称 |
| `info` | `varchar(30)` | YES | NULL | 配置分类说明 |
| `sort` | `smallint(5)` | NO | '0' | 排序 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `icon` | `varchar(30)` | YES | NULL | 图标 |
| `status` | `tinyint(3)` | NO | '1' | 配置分类状态 |

## `qixi_system_config_value`

> 配置表

原表：`eb_system_config_value`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `config_value_id` | `int(10)` | NO | — | 配置id |
| `config_key` | `varchar(32)` | NO | — | 配置分类key |
| `value` | `text` | NO | — | 值 |
| `mer_id` | `int(10)` | NO | '0' | 商户 id |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |

## `qixi_system_form`

原表：`eb_system_form`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `form_id` | `int(11)` | NO | — | — |
| `name` | `varchar(255)` | NO | '' | 表单标题 |
| `form_keys` | `text` | YES | — | 表单所有的key |
| `value` | `text` | YES | — | 表单内容 |
| `status` | `tinyint(4)` | NO | '1' | 状态 |
| `is_del` | `tinyint(4)` | NO | '0' | 是否删除 |
| `mer_id` | `int(11)` | NO | '0' | 商户ID |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `update_time` | `timestamp` | YES | NULL | — |

## `qixi_system_group`

> 组合数据表

原表：`eb_system_group`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `group_id` | `int(11)` | NO | — | 组合数据ID |
| `group_name` | `varchar(50)` | NO | — | 数据组名称 |
| `group_info` | `varchar(256)` | NO | — | 数据提示 |
| `group_key` | `varchar(50)` | NO | — | 数据字段 |
| `fields` | `text` | YES | — | 数据组字段以及类型（json数据） |
| `user_type` | `tinyint(3)` | NO | '0' | 0=总后台配置 1=商户后台配置 |
| `sort` | `smallint(5)` | NO | '0' | 排序 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |

## `qixi_system_group_data`

> 组合数据详情表

原表：`eb_system_group_data`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `group_data_id` | `int(10)` | NO | — | 组合数据详情ID |
| `group_id` | `int(10)` | NO | '0' | 对应的数据组id |
| `value` | `text` | NO | — | 数据组对应的数据值（json数据） |
| `sort` | `int(11)` | NO | '0' | 数据排序 |
| `status` | `tinyint(3)` | NO | '1' | 状态（1：开启；0：关闭；） |
| `mer_id` | `int(10)` | NO | '0' | 商户 id |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加数据时间 |

## `qixi_system_log`

> 管理员操作记录表

原表：`eb_system_log`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `log_id` | `int(10)` | NO | — | 管理员操作记录ID |
| `admin_id` | `int(10)` | NO | — | 管理员id |
| `admin_name` | `varchar(64)` | NO | — | 管理员姓名 |
| `route` | `varchar(128)` | NO | — | 路由 |
| `method` | `varchar(12)` | NO | — | 方式 |
| `url` | `varchar(256)` | NO | — | 链接 |
| `ip` | `varchar(16)` | NO | — | 登录IP |
| `mer_id` | `int(10)` | NO | '0' | 商户id |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |

## `qixi_system_menu`

> 菜单表

原表：`eb_system_menu`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `menu_id` | `int(10)` | NO | — | 菜单ID |
| `pid` | `int(10)` | NO | '0' | 父级id |
| `path` | `varchar(512)` | NO | — | 路径 |
| `icon` | `varchar(32)` | YES | '' | 图标 |
| `menu_name` | `varchar(128)` | NO | '' | 按钮名 |
| `route` | `varchar(64)` | NO | — | 路由名称 |
| `params` | `varchar(128)` | NO | '' | 参数 |
| `sort` | `tinyint(4)` | NO | '1' | 排序 |
| `is_show` | `tinyint(3)` | NO | '1' | 是否显示 |
| `is_mer` | `tinyint(3)` | NO | '1' | 模块，1 平台， 2商户 |
| `is_menu` | `tinyint(3)` | NO | '1' | 类型，1菜单 2 权限 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `update_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `is_agent` | `tinyint(1)` | NO | '0' | 0:平台,1:区域,2:商户 |

## `qixi_system_notice`

> 商户公告

原表：`eb_system_notice`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `notice_id` | `int(10)` | NO | — | — |
| `admin_id` | `int(10)` | NO | — | 管理员 id |
| `notice_title` | `varchar(128)` | NO | — | 通知标题 |
| `notice_content` | `text` | NO | — | 通知内容 |
| `type` | `tinyint(3)` | NO | — | 通知类型 |
| `type_str` | `varchar(512)` | NO | — | 通知说明 |
| `is_del` | `tinyint(3)` | NO | '0' | 0:正常 1:删除 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |
| `status` | `int(10)` | NO | '1' | 状态（0:关;1:开） |

## `qixi_system_notice_config`

原表：`eb_system_notice_config`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `notice_config_id` | `bigint(20)` | NO | — | — |
| `notice_title` | `varchar(20)` | YES | NULL | 消息名称 |
| `const_key` | `varchar(100)` | YES | NULL | 通知标识 |
| `notice_info` | `varchar(50)` | YES | NULL | 消息说明 |
| `notice_sys` | `tinyint(4)` | YES | '-1' | 站内消息 |
| `notice_sms` | `tinyint(4)` | YES | '-1' | 短信消息 |
| `notice_wechat` | `tinyint(4)` | YES | '-1' | 公众号模板消息 |
| `wechat_tempkey` | `varchar(100)` | YES | NULL | 微信模板关联ID |
| `wechat_content` | `varchar(255)` | NO | '' | 微信模板内容 |
| `wechat_tempid` | `varchar(255)` | YES | NULL | 微信模板ID |
| `notice_routine` | `tinyint(4)` | YES | '-1' | 小程序订阅消息 |
| `routine_tempkey` | `varchar(100)` | YES | NULL | 订阅消息关联ID |
| `routine_content` | `varchar(255)` | NO | '' | 小程序订阅消息内容 |
| `routine_tempid` | `varchar(255)` | YES | NULL | 小程序消息ID |
| `type` | `int(10)` | YES | '0' | 1商户通知， 0用户通知 |
| `sms_tempid` | `varchar(50)` | YES | NULL | 一号通短信模板ID |
| `sms_ali_tempid` | `varchar(50)` | YES | NULL | 阿里云短信模板ID |
| `sms_content` | `varchar(100)` | YES | NULL | 阿里云短信模板内容 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 创建时间  |
| `update_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 更新时间 |
| `kid` | `char(100)` | YES | '0' | — |

## `qixi_system_notice_log`

原表：`eb_system_notice_log`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `notice_log_id` | `int(10)` | NO | — | — |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `notice_id` | `int(10)` | NO | — | 公告 id |
| `is_read` | `tinyint(3)` | NO | '0' | 是否已读 |
| `read_time` | `timestamp` | YES | NULL | 读取时间 |
| `is_del` | `tinyint(3)` | NO | '0' | 是否删除 |

## `qixi_system_role`

> 身份管理表

原表：`eb_system_role`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `role_id` | `int(10)` | NO | — | 身份管理id |
| `role_name` | `varchar(32)` | NO | — | 身份管理名称 |
| `rules` | `text` | NO | — | 身份管理权限(menus_id) |
| `status` | `tinyint(3)` | NO | '1' | 状态 |
| `mer_id` | `int(10)` | NO | '0' | 商户 id |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `update_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 更新时间 |
| `is_agent` | `tinyint(1)` | NO | '0' | 角色类型：0:平台,1:区域,2:商户 |
| `circle_id` | `int(10)` | NO | '0' | 商圈ID |
| `is_default` | `tinyint(1)` | NO | '0' | 是否默认角色 |

## `qixi_system_storage`

原表：`eb_system_storage`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `bigint(20)` | NO | — | — |
| `type` | `tinyint(4)` | NO | '1' | — |
| `access_key` | `varchar(100)` | NO | '' | — |
| `name` | `varchar(100)` | NO | '' | — |
| `region` | `varchar(100)` | NO | '' | — |
| `acl` | `varchar(100)` | NO | '' | — |
| `domain` | `varchar(255)` | NO | '' | — |
| `cdn` | `varchar(255)` | YES | NULL | — |
| `status` | `int(11)` | NO | '0' | — |
| `is_del` | `tinyint(4)` | NO | '0' | — |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `update_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |

## `qixi_template_message`

> 微信模板

原表：`eb_template_message`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `template_id` | `int(10)` | NO | — | 模板id |
| `type` | `tinyint(1)` | NO | '0' | 0=订阅消息,1=微信模板消息 |
| `tempkey` | `char(50)` | NO | — | 模板编号 |
| `name` | `char(100)` | NO | — | 模板名 |
| `content` | `varchar(1000)` | NO | — | 回复内容 |
| `tempid` | `char(100)` | YES | NULL | 模板ID |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `status` | `tinyint(4)` | NO | '0' | 状态 |
| `kid` | `char(100)` | YES | '0' | — |
