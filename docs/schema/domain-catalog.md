# 商品/类目/品牌 — 表字段（`qixi_`）

## `qixi_cdkey_library`

原表：`eb_cdkey_library`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `int(10)` | NO | — | id |
| `mer_id` | `int(11)` | NO | '0' | 商户ID |
| `name` | `varchar(50)` | NO | '' | 卡密库名称 |
| `remark` | `varchar(255)` | NO | '' | 备注 |
| `product_id` | `int(11)` | NO | '0' | 关联商品ID |
| `product_attr_value_id` | `int(11)` | NO | '0' | 关联商品skuID |
| `used_num` | `int(10)` | NO | '0' | 卡密已使用数量 |
| `total_num` | `int(10)` | NO | '0' | 卡密总数量 |
| `is_del` | `tinyint(3)` | NO | '0' | 是否删除 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 创建时间 |

## `qixi_guarantee`

> 保障服务选项

原表：`eb_guarantee`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `guarantee_id` | `int(10)` | NO | — | 主键 |
| `guarantee_name` | `varchar(255)` | YES | NULL | 保障服务名称 |
| `guarantee_info` | `varchar(500)` | YES | NULL | 保障服务简介 |
| `status` | `int(11)` | YES | '1' | 0.关闭，1开启 |
| `image` | `varchar(255)` | YES | NULL | 图标 |
| `sort` | `int(11)` | YES | NULL | 排序 |
| `mer_count` | `int(11)` | YES | '0' | 使用的商户数 |
| `product_cout` | `int(11)` | YES | '0' | 使用的商品数 |
| `is_del` | `int(11)` | YES | '0' | — |
| `create_time` | `timestamp` | YES | NULL | — |
| `update_time` | `timestamp` | YES | NULL | — |

## `qixi_guarantee_template`

> 保障服务模板

原表：`eb_guarantee_template`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `guarantee_template_id` | `int(10)` | NO | — | — |
| `template_name` | `varchar(255)` | YES | NULL | — |
| `mer_id` | `int(11)` | YES | NULL | — |
| `status` | `int(10)` | YES | '1' | — |
| `sort` | `int(11)` | YES | NULL | — |
| `create_time` | `timestamp` | YES | NULL | — |
| `is_del` | `int(10)` | YES | '0' | — |

## `qixi_guarantee_value`

> 保障服务模板条款

原表：`eb_guarantee_value`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `guarantee_value_id` | `int(10)` | NO | — | — |
| `guarantee_id` | `int(10)` | YES | NULL | — |
| `guarantee_template_id` | `int(10)` | YES | NULL | — |
| `mer_id` | `int(11)` | YES | NULL | — |
| `status` | `int(11)` | YES | '1' | — |

## `qixi_parameter`

原表：`eb_parameter`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `parameter_id` | `int(10)` | NO | — | — |
| `template_id` | `int(10)` | NO | — | 模板 id |
| `mer_id` | `int(10)` | YES | '0' | 商户 id |
| `name` | `varchar(32)` | NO | — | 参数名称 |
| `value` | `varchar(255)` | YES | NULL | 参数值 |
| `required` | `tinyint(4)` | NO | '0' | 是否必填 |
| `create_time` | `datetime` | YES | CURRENT_TIMESTAMP | — |
| `sort` | `int(10)` | NO | '0' | — |

## `qixi_parameter_product`

原表：`eb_parameter_product`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `int(11)` | NO | — | — |
| `parameter_value_id` | `int(11)` | NO | '0' | 参数ID |
| `product_id` | `int(11)` | NO | '0' | 商品ID |

## `qixi_parameter_template`

原表：`eb_parameter_template`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `template_id` | `int(10)` | NO | — | — |
| `mer_id` | `int(10)` | NO | '0' | 商户id |
| `template_name` | `varchar(64)` | NO | — | 模板名称 |
| `sort` | `int(10)` | YES | '0' | — |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | — |

## `qixi_parameter_value`

原表：`eb_parameter_value`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `parameter_value_id` | `int(10)` | NO | — | — |
| `parameter_id` | `int(10)` | NO | '0' | 关联参数 id |
| `product_id` | `int(10)` | NO | — | 商品 id |
| `name` | `varchar(64)` | YES | NULL | 参数名称 |
| `value` | `varchar(64)` | YES | NULL | 参数值 |
| `sort` | `int(10)` | YES | '0' | 排序 |
| `create_time` | `datetime` | YES | NULL | — |
| `mer_id` | `int(10)` | NO | '0' | — |

## `qixi_price_rule`

原表：`eb_price_rule`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `rule_id` | `int(10)` | NO | — | — |
| `rule_name` | `varchar(64)` | NO | — | 名称 |
| `sort` | `int(10)` | YES | '0' | 排序 |
| `is_show` | `tinyint(3)` | YES | '1' | 是否显示 |
| `is_default` | `tinyint(3)` | YES | '0' | 是否默认 |
| `content` | `longtext` | YES | — | 内容 |
| `create_time` | `datetime` | YES | CURRENT_TIMESTAMP | — |
| `update_time` | `datetime` | YES | CURRENT_TIMESTAMP | — |

## `qixi_store_attr_template`

> 商品规则值(规格)表

原表：`eb_store_attr_template`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `attr_template_id` | `int(11)` | NO | — | — |
| `template_name` | `varchar(32)` | NO | — | 规格名称 |
| `template_value` | `text` | NO | — | 规格值 |
| `mer_id` | `int(11)` | NO | — | 商户 id |

## `qixi_store_brand`

> 商品品牌表

原表：`eb_store_brand`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `brand_id` | `mediumint(9)` | NO | — | 商品品牌表ID |
| `brand_category_id` | `mediumint(9)` | NO | — | 父id |
| `brand_name` | `varchar(100)` | NO | — | 品牌名称 |
| `sort` | `mediumint(9)` | NO | — | 排序 |
| `pic` | `varchar(128)` | NO | '' | 图标 |
| `is_show` | `tinyint(1)` | NO | '1' | 是否显示 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |

## `qixi_store_brand_category`

> 品牌分类表

原表：`eb_store_brand_category`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `store_brand_category_id` | `mediumint(9)` | NO | — | 品牌分类表ID |
| `pid` | `mediumint(9)` | NO | — | 父id |
| `cate_name` | `varchar(100)` | NO | — | 分类名称 |
| `path` | `varchar(255)` | NO | '' | 路径 |
| `sort` | `mediumint(9)` | NO | — | 排序 |
| `is_show` | `tinyint(1)` | NO | '1' | 是否显示 |
| `level` | `int(10)` | NO | '0' | 等级 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 添加时间 |

## `qixi_store_category`

> 商品分类表

原表：`eb_store_category`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `store_category_id` | `mediumint(9)` | NO | — | 商品分类表ID |
| `pid` | `mediumint(9)` | NO | — | 父id |
| `cate_name` | `varchar(100)` | NO | — | 分类名称 |
| `path` | `varchar(255)` | NO | '' | 路径 |
| `sort` | `mediumint(9)` | NO | — | 排序 |
| `pic` | `varchar(128)` | NO | '' | 图标 |
| `is_show` | `tinyint(1)` | NO | '1' | 是否显示 |
| `level` | `int(10)` | NO | '0' | 等级 |
| `mer_id` | `int(10)` | YES | '0' | 商户id |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 添加时间 |
| `is_hot` | `tinyint(1)` | YES | '0' | 是否推荐 |
| `type` | `tinyint(1)` | NO | '0' | 0 商品，1 积分商品 |

## `qixi_store_product`

> 商品表

原表：`eb_store_product`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_id` | `int(10)` | NO | — | 商品id |
| `mer_id` | `int(10)` | NO | '0' | 商户Id |
| `store_name` | `varchar(128)` | NO | — | 商品名称 |
| `store_info` | `varchar(256)` | YES | NULL | 商品简介 |
| `keyword` | `varchar(128)` | NO | — | 关键字 |
| `bar_code` | `varchar(15)` | NO | '' | 产品条码（一维码） |
| `brand_id` | `int(11)` | YES | NULL | 品牌 id |
| `is_show` | `tinyint(3)` | NO | '1' | 商户 状态（0:未上架，1:上架，2:定时上架） |
| `status` | `tinyint(1)` | NO | '0' | 管理员 状态（0：审核中，1：审核通过 -1: 未通过 -2: 下架） |
| `is_del` | `tinyint(3)` | NO | '0' | 是否删除 |
| `mer_status` | `tinyint(1)` | YES | '1' | 商铺状态是否 1.正常 0. 非正常 |
| `cate_id` | `int(11)` | NO | — | 分类id |
| `unit_name` | `varchar(16)` | NO | — | 单位名 |
| `sort` | `smallint(6)` | NO | '0' | 排序 |
| `rank` | `smallint(6)` | NO | '0' | 总后台排序 |
| `sales` | `mediumint(8)` | NO | '0' | 销量 |
| `price` | `decimal(10,2)` | YES | '0.00' | 最低价格 |
| `cost` | `decimal(10,2)` | YES | '0.00' | 成本价 |
| `ot_price` | `decimal(10,2)` | YES | '0.00' | 原价 |
| `stock` | `int(10)` | YES | '0' | 总库存 |
| `is_hot` | `tinyint(3)` | NO | '0' | 是否热卖 |
| `is_benefit` | `tinyint(3)` | NO | '0' | 促销单品 |
| `is_best` | `tinyint(3)` | NO | '0' | 是否精品 |
| `is_new` | `tinyint(3)` | NO | '0' | 是否新品 |
| `is_good` | `tinyint(1)` | NO | '0' | 是否优品推荐 |
| `product_type` | `tinyint(3)` | NO | '0' | 0.普通商品 1.秒杀商品,2.预售商品，3.助力商品，4.拼团商品 |
| `ficti` | `mediumint(9)` | YES | '0' | 虚拟销量 |
| `browse` | `int(11)` | YES | '0' | 浏览量 |
| `code_path` | `varchar(64)` | NO | '' | 产品二维码地址(用户小程序海报) |
| `video_link` | `varchar(200)` | NO | '' | 主图视频链接 |
| `temp_id` | `int(11)` | NO | '1' | 运费模板ID |
| `spec_type` | `tinyint(1)` | NO | '0' | 规格 0单 1多 |
| `extension_type` | `tinyint(1)` | YES | '0' | 佣金比例 0.系统，1.自定义 |
| `refusal` | `varchar(255)` | YES | NULL | 审核拒绝理由 |
| `rate` | `decimal(2,1)` | YES | '5.0' | 评价分数 |
| `reply_count` | `int(10)` | YES | '0' | 评论数 |
| `give_coupon_ids` | `varchar(500)` | YES | NULL | 赠送优惠券 |
| `is_gift_bag` | `tinyint(1)` | YES | '0' | 是否为礼包 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `care_count` | `int(11)` | NO | '0' | 收藏数 |
| `is_used` | `int(11)` | YES | '1' | 显示/隐藏 |
| `old_product_id` | `int(11)` | YES | '0' | 原商品ID |
| `image` | `varchar(256)` | NO | '' | 商品图片 |
| `slider_image` | `varchar(2000)` | NO | — | 轮播图 |
| `guarantee_template_id` | `int(11)` | YES | '0' | 保障服务模板 |
| `once_max_count` | `int(11)` | YES | '0' | 订单单次购买数量最大限制 |
| `once_min_count` | `int(11)` | NO | '0' | 单次购买最低限购 |
| `integral_rate` | `int(11)` | NO | '-1' | 积分抵扣比例 |
| `integral_total` | `int(10)` | YES | '0' | 使用积分抵扣总数 |
| `integral_price_total` | `decimal(8,2)` | YES | '0.00' | 使用积分抵扣金额总数 |
| `labels` | `varchar(255)` | YES | '' | 标签id |
| `delivery_way` | `varchar(100)` | YES | NULL | 1.仅到店自提2快递计价配送3全国包邮 |
| `delivery_free` | `int(11)` | YES | '0' | 全国包邮 |
| `type` | `tinyint(3)` | NO | '0' | 0.实体商品，1.虚拟商品，2 网盘，3 卡密 |
| `extend` | `varchar(1000)` | NO | '' | 扩展信息 |
| `pay_limit` | `tinyint(4)` | NO | '0' | 购买总数限制 0:不限购，1单次限购 2 长期限购 |
| `svip_price_type` | `tinyint(3)` | NO | '0' | 0不参加，1默认比例，2自定义 |
| `svip_price` | `decimal(10,2)` | NO | '0.00' | 会员价 |
| `mer_svip_status` | `tinyint(3)` | NO | '1' | 商户会员状态 |
| `param_temp_id` | `varchar(255)` | YES | NULL | 参数模板ID |
| `refund_switch` | `tinyint(4)` | YES | '1' | 是否支持退款 |
| `delete` | `tinyint(4)` | YES | '0' | — |
| `mer_form_id` | `int(11)` | YES | '0' | 系统表单ID |
| `good_ids` | `varchar(2000)` | YES | '' | 推荐商品 |
| `auto_on_time` | `int(11)` | YES | NULL | 自动上架时间 |
| `auto_off_time` | `int(11)` | YES | NULL | 自动下架时间 |
| `active_id` | `int(11)` | YES | NULL | 秒杀活动ID |
| `cate_hot` | `tinyint(4)` | YES | '0' | 分类大图推荐 1 推荐 |
| `bar_code_number` | `varchar(255)` | YES | '' | 商品条码 |
| `custom_temp_id` | `varchar(255)` | NO | '' | 自定义参数模版id |
| `activity_label_ids` | `varchar(255)` | NO | '' | 活动标签 |

## `qixi_store_product_assist`

> 商品助力活动表

原表：`eb_store_product_assist`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_assist_id` | `int(10)` | NO | — | — |
| `start_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 开始时间 |
| `end_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 结束时间 |
| `status` | `int(10)` | NO | '1' | 平台控制状态：1开启，0.结束 |
| `pay_count` | `int(10)` | YES | '0' | 限购数量，0为不限制 |
| `assist_count` | `int(10)` | YES | '0' | 助力总需人数 |
| `assist_user_count` | `int(10)` | YES | '0' | 单人可助力次数 |
| `product_id` | `int(10)` | NO | '0' | 商品ID |
| `is_show` | `tinyint(3)` | NO | '1' | 商户控制状态 0.下架；1.上架 |
| `store_name` | `varchar(128)` | NO | — | 商品活动标题 |
| `mer_id` | `int(10)` | NO | '0' | — |
| `store_info` | `varchar(255)` | YES | NULL | 商品简介 |
| `is_del` | `int(10)` | NO | '0' | — |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |
| `product_status` | `int(11)` | YES | '0' | 审核状态；0.待审核，1审核通过，-1 审核失败，-2 强制下架 |
| `refusal` | `varchar(255)` | YES | NULL | — |
| `action_status` | `int(11)` | YES | '1' | 活动状态1开启，-1 结束 |

## `qixi_store_product_assist_set`

> 助力发起列表

原表：`eb_store_product_assist_set`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_assist_set_id` | `int(10)` | NO | — | — |
| `product_assist_id` | `int(10)` | NO | — | — |
| `product_id` | `int(10)` | NO | — | — |
| `uid` | `int(10)` | NO | '0' | — |
| `status` | `int(11)` | NO | '1' | 状态：-1 未完成 ，1 进行中， 10 已完成，20.已支付 |
| `assist_count` | `int(10)` | NO | '0' | 需助力总人数 |
| `assist_user_count` | `int(10)` | NO | '0' | 单人可助力次数 |
| `yet_assist_count` | `int(10)` | NO | '0' | 已助力人数 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 时间 |
| `mer_id` | `int(10)` | YES | '0' | — |
| `share_num` | `int(10)` | YES | '0' | — |
| `view_num` | `int(10)` | YES | '0' | — |
| `is_del` | `int(10)` | YES | '0' | — |

## `qixi_store_product_assist_sku`

原表：`eb_store_product_assist_sku`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_assist_id` | `int(10)` | NO | '0' | — |
| `product_id` | `int(10)` | NO | — | — |
| `unique` | `char(12)` | NO | — | — |
| `assist_price` | `decimal(10,2)` | NO | '0.00' | 助力售价 |
| `stock` | `int(10)` | NO | '0' | — |
| `stock_count` | `int(10)` | YES | '0' | 总限购 |

## `qixi_store_product_assist_user`

> 助力记录表

原表：`eb_store_product_assist_user`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_assist_user_id` | `int(10)` | NO | — | — |
| `product_assist_set_id` | `int(10)` | NO | — | — |
| `product_assist_id` | `int(10)` | NO | — | — |
| `uid` | `int(10)` | NO | '0' | id |
| `avatar_img` | `varchar(256)` | YES | NULL | 头像 |
| `nickname` | `varchar(50)` | YES | NULL | 昵称 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |

## `qixi_store_product_attr`

> 商品属性表

原表：`eb_store_product_attr`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_id` | `int(10)` | NO | '0' | 商品ID |
| `attr_name` | `varchar(32)` | NO | — | 属性名 |
| `attr_values` | `varchar(2000)` | NO | — | 属性值 |
| `type` | `tinyint(1)` | YES | '0' | 活动类型 0=商品 |

## `qixi_store_product_attr_reservation`

原表：`eb_store_product_attr_reservation`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `attr_reservation_id` | `int(10)` | NO | — | attr_reservation_id |
| `attr_value_id` | `int(10)` | NO | '0' | 商品属性值表ID |
| `unique` | `char(12)` | NO | '' | 唯一值 |
| `product_id` | `int(10)` | NO | '0' | 商品ID |
| `start_time` | `varchar(10)` | NO | '' | 开始时间段 |
| `end_time` | `varchar(10)` | NO | '' | 结束时间段 |
| `stock` | `int(11)` | NO | '0' | 可约数量 |
| `use_num` | `int(11)` | NO | '0' | 使用数量 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `update_time` | `datetime` | NO | CURRENT_TIMESTAMP | 修改时间 |

## `qixi_store_product_attr_result`

原表：`eb_store_product_attr_result`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `id` | `int(11)` | NO | — | 主键id |
| `product_id` | `int(10)` | NO | — | 商品ID |
| `result` | `longtext` | NO | — | 商品属性参数 |
| `change_time` | `int(10)` | YES | NULL | 上次修改时间 |
| `type` | `tinyint(1)` | YES | NULL | 活动类型 0=商品，1=秒杀，2=预售，3=助力, 4=拼团 |

## `qixi_store_product_attr_value`

> 商品属性值表

原表：`eb_store_product_attr_value`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `value_id` | `int(10)` | NO | — | 主键 |
| `product_id` | `int(10)` | NO | — | 商品ID |
| `detail` | `varchar(1000)` | NO | '' | — |
| `sku` | `varchar(128)` | NO | — | 商品属性索引值 (attr_value|attr_value[|....]) |
| `stock` | `int(10)` | NO | — | 属性对应的库存 |
| `sales` | `int(10)` | NO | '0' | 销量 |
| `image` | `varchar(128)` | YES | NULL | 图片 |
| `bar_code` | `varchar(50)` | NO | '' | 产品条码 |
| `cost` | `decimal(8,2)` | NO | — | 成本价 |
| `ot_price` | `decimal(8,2)` | NO | '0.00' | 原价 |
| `price` | `decimal(8,2)` | NO | — | 价格 |
| `volume` | `decimal(8,2)` | NO | '0.00' | 体积 |
| `weight` | `decimal(8,2)` | NO | '0.00' | 重量 |
| `type` | `tinyint(1)` | YES | '0' | 活动类型 0=商品; 20 积分商品 |
| `extension_one` | `decimal(8,2)` | YES | '0.00' | 一级佣金 |
| `extension_two` | `decimal(8,2)` | YES | '0.00' | 二级佣金 |
| `unique` | `char(12)` | NO | '' | 唯一值 |
| `svip_price` | `decimal(10,2)` | NO | '0.00' | 会员价 |
| `library_id` | `int(11)` | YES | '0' | 一次性卡密关联卡密库地 |
| `bar_code_number` | `varchar(50)` | NO | '' | 规格条码 |
| `is_default_select` | `tinyint(1)` | YES | '0' | 默认显示 |
| `is_show` | `tinyint(1)` | YES | '1' | — |
| `settlement_price` | `decimal(8,2)` | NO | '0.00' | 员工结算价格 |

## `qixi_store_product_cate`

> 商品商户分类关联表

原表：`eb_store_product_cate`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_id` | `int(11)` | YES | NULL | — |
| `mer_cate_id` | `int(11)` | YES | NULL | — |
| `mer_id` | `int(11)` | YES | NULL | — |

## `qixi_store_product_cdkey`

原表：`eb_store_product_cdkey`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `cdkey_id` | `int(11)` | NO | — | — |
| `is_type` | `tinyint(4)` | YES | '0' | 卡密类型： 0 固定卡密， 1 一次性卡密 |
| `value_id` | `int(11)` | YES | NULL | 商品规格ID |
| `key` | `varchar(255)` | YES | NULL | 卡密内容 |
| `pwd` | `varchar(255)` | YES | NULL | 卡密密码 |
| `status` | `tinyint(4)` | YES | '1' | 状态： 1 可用  -1 不可用/已使用 |
| `product_id` | `int(11)` | YES | NULL | 商品ID |
| `library_id` | `int(11)` | YES | '0' | 卡密库ID |
| `is_use` | `tinyint(4)` | YES | '0' | 是否已使用 |
| `mer_id` | `int(11)` | YES | NULL | 商户ID |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | 添加时间 |

## `qixi_store_product_content`

> 商品详情表

原表：`eb_store_product_content`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_id` | `int(10)` | NO | — | 商品id |
| `content` | `longtext` | NO | — | 商品详情 |
| `type` | `tinyint(1)` | NO | '0' | 商品类型 0=普通 |

## `qixi_store_product_copy`

原表：`eb_store_product_copy`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `store_product_copy_id` | `int(10)` | NO | — | — |
| `type` | `varchar(255)` | YES | NULL | — |
| `mer_id` | `int(11)` | YES | NULL | 商户id |
| `num` | `int(11)` | YES | NULL | 数量 |
| `number` | `int(11)` | YES | '1' | 剩余数量 |
| `message` | `varchar(255)` | YES | NULL | — |
| `info` | `text` | YES | — | 信息 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |

## `qixi_store_product_group`

> 拼团商品信息表

原表：`eb_store_product_group`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_group_id` | `int(10)` | NO | — | 主键ID |
| `product_id` | `int(10)` | YES | '0' | 商品ID |
| `start_time` | `datetime` | YES | NULL | 开始时间 |
| `end_time` | `datetime` | YES | NULL | 结束时间 |
| `time` | `int(10)` | YES | '0' | 开团时长 |
| `buying_count_num` | `int(11)` | YES | '0' | 拼团总人数 |
| `buying_num` | `int(11)` | YES | '0' | 最少真实购买人数 |
| `pay_count` | `int(10)` | YES | '0' | 活动购买总人数 |
| `once_pay_count` | `int(10)` | YES | '0' | 单次购买数量 |
| `status` | `int(11)` | YES | '0' | 平台控制状态 |
| `mer_id` | `int(10)` | YES | '0' | 商户ID |
| `ficti_status` | `int(11)` | YES | '0' | 虚拟成团状态 |
| `ficti_num` | `int(11)` | YES | '0' | 最多虚拟人数 |
| `is_show` | `int(11)` | YES | '0' | 上下架 |
| `is_del` | `int(10)` | YES | '0' | — |
| `success_num` | `int(10)` | YES | '0' | 成功团数 |
| `product_status` | `int(11)` | YES | '0' | — |
| `price` | `decimal(10,2)` | YES | '0.00' | — |
| `action_status` | `int(11)` | YES | '0' | 活动状态 |
| `create_time` | `datetime` | YES | NULL | — |
| `refusal` | `varchar(255)` | YES | NULL | — |
| `leader_extension` | `tinyint(4)` | YES | '0' | 团长分销 |
| `leader_rate` | `decimal(10,2)` | YES | '0.00' | 分销比例 |

## `qixi_store_product_group_buying`

> 拼团活动表

原表：`eb_store_product_group_buying`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `group_buying_id` | `int(10)` | NO | — | — |
| `product_group_id` | `int(10)` | YES | '0' | 活动商品ID |
| `status` | `int(11)` | YES | '0' | 状态：0。默认，进行中，10.已完成，-1 时间到未完成 |
| `ficti_status` | `int(11)` | YES | '0' | 虚拟成团状态0.未开启，1开启 |
| `ficti_num` | `int(10)` | YES | '0' | 虚拟成团人数 |
| `buying_count_num` | `int(10)` | YES | '0' | 成团总人数 |
| `buying_num` | `int(10)` | YES | '0' | 真实人数 |
| `yet_buying_num` | `int(10)` | YES | '0' | 已参团人数 |
| `is_del` | `int(11)` | YES | '0' | — |
| `mer_id` | `int(10)` | YES | '0' | — |
| `end_time` | `int(11)` | YES | NULL | 结束时间 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `is_hidde` | `tinyint(1)` | YES | '0' | 是否隐藏团信息 0 否 1 是 |

## `qixi_store_product_group_sku`

原表：`eb_store_product_group_sku`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_group_id` | `int(10)` | NO | '0' | — |
| `product_id` | `int(10)` | NO | — | — |
| `unique` | `char(12)` | NO | — | — |
| `active_price` | `decimal(10,2)` | NO | '0.00' | 活动价 |
| `stock` | `int(10)` | NO | '0' | — |
| `stock_count` | `int(10)` | YES | '0' | — |

## `qixi_store_product_group_user`

> 拼团成员表

原表：`eb_store_product_group_user`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `group_buying_id` | `int(10)` | YES | '0' | 团ID |
| `product_group_id` | `int(10)` | YES | '0' | 活动商品ID |
| `status` | `int(11)` | YES | '0' | 状态 |
| `is_initiator` | `int(10)` | YES | '0' | 是否为 团长 |
| `order_id` | `int(10)` | YES | '0' | 订单ID |
| `uid` | `int(10)` | YES | '0' | 用户ID  |
| `nickname` | `varchar(255)` | YES | NULL | 昵称 |
| `avatar` | `varchar(255)` | YES | NULL | 头像 |
| `is_del` | `int(10)` | YES | '0' | — |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `is_leader` | `tinyint(1)` | YES | '0' | 是否为创建者 |

## `qixi_store_product_label`

原表：`eb_store_product_label`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_label_id` | `int(10)` | NO | — | — |
| `label_name` | `varchar(50)` | YES | NULL | 标签名 |
| `status` | `tinyint(1)` | YES | NULL | 状态 |
| `info` | `varchar(255)` | YES | NULL | 说明 |
| `sort` | `int(11)` | YES | NULL | 排序 |
| `type` | `int(11)` | YES | '0' | 类型  |
| `mer_id` | `int(11)` | YES | '0' | 商户ID |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `is_del` | `tinyint(1)` | YES | '0' | — |

## `qixi_store_product_presell`

> 商品预售活动表

原表：`eb_store_product_presell`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_presell_id` | `int(10)` | NO | — | — |
| `start_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 预售开始时间 |
| `end_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 预售结束时间 |
| `final_start_time` | `varchar(30)` | YES | '' | 尾款支付开始时间 |
| `final_end_time` | `varchar(30)` | YES | '' | 尾款支付结时间 |
| `status` | `int(10)` | NO | '1' | 平台控制状态：1开启，0.结束 |
| `presell_type` | `int(10)` | NO | '0' | 预售类型：1.全款预售，2.定金预售 |
| `pay_count` | `int(10)` | YES | '0' | 限购数量，0为不限制 |
| `delivery_type` | `int(10)` | NO | '0' | 发货类型：1.支付成功后 ； 2. 预售结束后 |
| `delivery_day` | `int(10)` | YES | '0' | 发货时间 |
| `product_id` | `int(10)` | NO | '0' | 商品ID |
| `price` | `decimal(10,2)` | NO | '0.00' | 预售最低价 |
| `is_show` | `tinyint(3)` | YES | NULL | 商户控制状态 0.下架；1.上架 |
| `store_name` | `varchar(128)` | NO | — | 商品活动标题 |
| `mer_id` | `int(10)` | NO | '0' | — |
| `store_info` | `varchar(255)` | YES | NULL | 商品简介 |
| `is_del` | `int(10)` | NO | '0' | — |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |
| `product_status` | `int(11)` | YES | '0' | 审核状态；0.待审核，1审核通过，-1 审核失败，-2 强制下架 |
| `refusal` | `varchar(255)` | YES | NULL | — |
| `action_status` | `int(11)` | YES | '1' | 活动状态1开启，-1 结束 |

## `qixi_store_product_presell_sku`

原表：`eb_store_product_presell_sku`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_presell_id` | `int(10)` | NO | '0' | — |
| `product_id` | `int(10)` | NO | — | — |
| `unique` | `char(12)` | NO | — | — |
| `presell_price` | `decimal(10,2)` | NO | '0.00' | 预售价 |
| `stock` | `int(10)` | NO | '0' | — |
| `stock_count` | `int(11)` | NO | '0' | 总限购 |
| `down_price` | `decimal(10,2)` | YES | '0.00' | 订金 |
| `final_price` | `decimal(10,2)` | NO | '0.00' | 尾款金额 |
| `one_take` | `int(10)` | YES | '0' | 第一阶段参与人数 |
| `one_pay` | `int(10)` | YES | '0' | 第一阶段支付人数 |
| `two_pay` | `int(10)` | YES | '0' | 第二阶段支付人数 |
| `seles` | `int(10)` | YES | '0' | 销量 |

## `qixi_store_product_reply`

> 商品评论表

原表：`eb_store_product_reply`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `reply_id` | `int(11)` | NO | — | 评论ID |
| `uid` | `int(11)` | NO | — | 用户ID |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `order_product_id` | `int(11)` | NO | — | 订单商品ID |
| `unique` | `char(12)` | YES | NULL | 商品 sku |
| `product_id` | `int(11)` | NO | — | 商品id |
| `product_type` | `tinyint(4)` | NO | '0' | 0=普通商品 |
| `product_score` | `tinyint(1)` | NO | — | 商品分数 |
| `service_score` | `tinyint(1)` | NO | — | 服务分数 |
| `postage_score` | `tinyint(1)` | NO | — | 物流分数 |
| `rate` | `float(2,1)` | YES | '5.0' | 平均值 |
| `comment` | `varchar(512)` | NO | — | 评论内容 |
| `pics` | `text` | YES | — | 评论图片 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 评论时间 |
| `merchant_reply_content` | `varchar(300)` | YES | NULL | 管理员回复内容 |
| `merchant_reply_time` | `timestamp` | YES | NULL | 管理员回复时间 |
| `sort` | `smallint(5)` | NO | '1' | 商家排序 |
| `is_del` | `tinyint(3)` | NO | '0' | 0未删除1已删除 |
| `is_reply` | `tinyint(1)` | NO | '0' | 0未回复1已回复 |
| `is_virtual` | `tinyint(1)` | NO | '0' | 0不是虚拟评价1是虚拟评价 |
| `nickname` | `varchar(64)` | NO | — | 用户名称 |
| `avatar` | `varchar(255)` | NO | — | 用户头像 |

## `qixi_store_product_reservation`

原表：`eb_store_product_reservation`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_reservation_id` | `int(10)` | NO | — | product_reservation_id |
| `product_id` | `int(10)` | NO | '0' | 商品id |
| `reservation_time_type` | `tinyint(1)` | NO | '1' | 预约时段划分类型(1:自动划分,2:自定义划分) |
| `reservation_start_time` | `varchar(20)` | NO | '' | 预约开始时间 |
| `reservation_end_time` | `varchar(20)` | NO | '' | 预约结束时间 |
| `reservation_time_interval` | `int(11)` | NO | '10' | 时间跨度,以分钟为单位 |
| `time_period` | `text` | NO | — | 时间段 |
| `reservation_type` | `tinyint(1)` | NO | '2' | 预约类型(1:到店服务,2:上门服务,3:上门+到店服务) |
| `show_num_type` | `tinyint(1)` | NO | '0' | 是否展示可约数量(0:不展示,1:展示) |
| `sale_time_type` | `tinyint(1)` | NO | '1' | 可售日期(1:每天,2:自定义时间) |
| `sale_time_start_day` | `varchar(20)` | NO | '' | 可售日期自定义开始时间 |
| `sale_time_end_day` | `varchar(20)` | NO | '' | 可售日期自定义结束时间 |
| `sale_time_week` | `varchar(20)` | NO | '' | 可售日期周数据 |
| `show_reservation_days` | `int(11)` | NO | '1' | 显示日期范围 |
| `is_advance` | `tinyint(1)` | NO | '0' | 是否提前预约(0:不提前,1:提前) |
| `advance_time` | `int(11)` | NO | '1' | 提前预约时间,以小时为单位 |
| `is_cancel_reservation` | `tinyint(1)` | NO | '0' | 是否可取消预约(0:不可取消,1:可取消) |
| `cancel_reservation_time` | `int(11)` | NO | '1' | 取消预约提前时间,以小时为单位 |
| `reservation_form_type` | `tinyint(1)` | NO | '1' | 预约表单类型(1:每个预约提交,2:每单提交) |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `update_time` | `datetime` | NO | CURRENT_TIMESTAMP | 修改时间 |
| `is_del` | `tinyint(1)` | NO | '0' | 是否删除 |

## `qixi_store_product_sku`

原表：`eb_store_product_sku`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_sku_id` | `int(11)` | NO | — | — |
| `active_id` | `int(10)` | NO | '0' | 活动ID |
| `active_product_id` | `int(10)` | YES | NULL | 活动商品的ID |
| `active_type` | `int(10)` | YES | '0' | 活动类型 |
| `product_id` | `int(10)` | NO | — | — |
| `unique` | `char(12)` | NO | — | — |
| `price` | `decimal(10,2)` | YES | '0.00' | 原价 |
| `active_price` | `decimal(10,2)` | NO | '0.00' | 活动售价 |
| `stock` | `int(10)` | NO | '0' | — |
| `stock_count` | `int(10)` | YES | '0' | 总限购 |

## `qixi_store_product_take`

> 用户到货通知记录

原表：`eb_store_product_take`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_take_id` | `int(10)` | NO | — | — |
| `product_id` | `int(11)` | YES | NULL | — |
| `unique` | `char(12)` | YES | NULL | — |
| `status` | `int(10)` | YES | '0' | 默认0，发送 1 |
| `uid` | `int(11)` | YES | NULL | 用户 |
| `type` | `varchar(255)` | YES | NULL | 1.PC,2.公众号,3.小程序 |
| `is_del` | `int(11)` | YES | '0' | — |

## `qixi_store_product_unit`

原表：`eb_store_product_unit`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `product_unit_id` | `int(11)` | NO | — | id |
| `mer_id` | `int(11)` | NO | '0' | 商户id |
| `value` | `varchar(255)` | NO | '' | 值 |
| `status` | `tinyint(1)` | NO | '1' | 状态 1|开启 0|关闭 |
| `is_del` | `tinyint(1)` | NO | '0' | 是否删除 0|正常 1|删除 |
| `sort` | `int(11)` | NO | '0' | 排序 |
| `create_time` | `datetime` | NO | CURRENT_TIMESTAMP | 创建时间 |

## `qixi_store_spu`

> 商品搜索信息表

原表：`eb_store_spu`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `spu_id` | `int(10)` | NO | — | — |
| `mer_id` | `int(10)` | YES | '0' | 商户ID |
| `product_id` | `int(10)` | YES | '0' | 商品ID |
| `product_type` | `int(10)` | YES | '0' | 活动类型0普通，1秒杀，2预售，3助力 |
| `activity_id` | `int(10)` | YES | '0' | 活动ID |
| `status` | `int(11)` | YES | '0' | 0.下架，1.上架 |
| `store_name` | `varchar(128)` | YES | NULL | 商品名称 |
| `ot_price` | `decimal(10,2)` | YES | '0.00' | — |
| `keyword` | `varchar(255)` | YES | NULL | 关键词 |
| `price` | `decimal(10,2)` | YES | '0.00' | 最低价格 |
| `rank` | `int(11)` | YES | NULL | 排序 |
| `create_time` | `datetime` | YES | NULL | — |
| `temp_id` | `int(10)` | YES | '0' | 运费模板 |
| `sort` | `int(10)` | YES | '0' | 商户排序 |
| `star` | `int(11)` | YES | '1' | 星级 |
| `image` | `varchar(255)` | YES | NULL | 主图 |
| `is_del` | `int(10)` | YES | '0' | — |
| `mer_labels` | `varchar(255)` | YES | '' | 标签id |
| `sys_labels` | `varchar(255)` | YES | '' | 标签id |
