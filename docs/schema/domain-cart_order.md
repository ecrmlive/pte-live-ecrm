# 购物车/订单 — 表字段（`qixi_`）

## `qixi_store_cart`

> 购物车表

原表：`eb_store_cart`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `cart_id` | `bigint(20)` | NO | — | 购物车表ID |
| `uid` | `int(10)` | NO | — | 用户ID |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `product_type` | `tinyint(4)` | NO | '0' | 类型 0=普通产品，2.预售商品 |
| `product_id` | `int(10)` | NO | — | 商品ID |
| `product_attr_unique` | `varchar(16)` | NO | '' | 商品属性 |
| `cart_num` | `smallint(5)` | NO | '0' | 商品数量 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `source` | `tinyint(3)` | NO | '0' | 来源 1.直播间,2.预售商品,3.助力商品 |
| `source_id` | `int(10)` | NO | '0' | 来源关联 id |
| `is_pay` | `tinyint(1)` | NO | '0' | 0 = 未购买 1 = 已购买 |
| `is_del` | `tinyint(1)` | NO | '0' | 是否删除 |
| `is_new` | `tinyint(1)` | NO | '0' | 是否为立即购买 |
| `is_fail` | `tinyint(3)` | NO | '0' | 是否失效 |
| `spread_id` | `int(10)` | YES | '0' | 推广人 |
| `tourist_unique_key` | `varchar(20)` | NO | '' | 游客唯一标识 |
| `reservation_date` | `varchar(20)` | YES | '' | 预约商品的预约日期 |
| `reservation_id` | `int(10)` | YES | '0' | 预约商品的预约时间段ID |

## `qixi_store_cart_price`

原表：`eb_store_cart_price`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `cart_price_id` | `bigint(20)` | NO | — | 购物车价格表ID |
| `cart_id` | `int(10)` | NO | — | 购物车id |
| `old_price` | `decimal(8,2)` | NO | '0.00' | 改价前价格 |
| `type` | `int(10)` | NO | '0' | 改价类型(0:一口价，1:立减金额，2:折扣率) |
| `reduce_price` | `decimal(8,2)` | NO | '0.00' | 立减金额（type为1时） |
| `discount_rate` | `int(10)` | NO | '0' | 折扣率（type为2时） |
| `new_price` | `decimal(8,2)` | NO | '0.00' | 改价后价格(type为0时一口价，或type不为0计算后的金额) |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 添加时间 |
| `update_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 修改时间 |
| `is_batch` | `tinyint(1)` | NO | '0' | 是否为批量改价(0:否，1:是) |

## `qixi_store_group_order`

> 用户订单表

原表：`eb_store_group_order`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `group_order_id` | `int(10)` | NO | — | — |
| `group_order_sn` | `varchar(32)` | NO | — | 订单号 |
| `uid` | `int(10)` | NO | — | 用户 ID |
| `total_postage` | `decimal(8,2)` | NO | '0.00' | 邮费 |
| `total_price` | `decimal(8,2)` | NO | '0.00' | 订单总额 |
| `total_num` | `int(10)` | NO | '0' | 商品数 |
| `integral` | `int(10)` | YES | '0' | 使用积分数量 |
| `integral_price` | `decimal(10,2)` | YES | '0.00' | 积分抵扣金额 |
| `give_integral` | `int(10)` | YES | '0' | 赠送积分 |
| `coupon_price` | `decimal(8,2)` | NO | '0.00' | 优惠金额 |
| `real_name` | `varchar(32)` | NO | — | 联系人 |
| `user_phone` | `varchar(18)` | NO | — | 联系电话 |
| `user_address` | `varchar(128)` | NO | — | 收货地址 |
| `pay_price` | `decimal(8,2)` | NO | — | 支付金额 |
| `pay_postage` | `decimal(8,2)` | NO | '0.00' | 支付邮费 |
| `cost` | `decimal(8,2)` | NO | — | 成本价 |
| `coupon_id` | `varchar(128)` | YES | NULL |  平台优惠券 |
| `give_coupon_ids` | `varchar(500)` | YES | '' | 赠送优惠券 |
| `paid` | `tinyint(3)` | NO | '0' | 是否支付 |
| `pay_time` | `timestamp` | YES | NULL | 支付时间 |
| `pay_type` | `tinyint(1)` | NO | — | 支付方式 0=余额 1=微信 2=小程序 3=h5 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `is_remind` | `tinyint(3)` | NO | '0' | 是否提醒 |
| `is_del` | `tinyint(3)` | NO | '0' | — |
| `is_combine` | `tinyint(3)` | YES | '0' | 是否为合并支付  |
| `activity_type` | `tinyint(3)` | YES | '0' | — |
| `is_first` | `tinyint(1)` | YES | '0' | 是否为用户首单 |
| `is_behalf` | `tinyint(1)` | NO | '0' | 是否为代客下单订单 |

## `qixi_store_import`

> 导入批次记录

原表：`eb_store_import`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `import_id` | `int(10)` | NO | — | — |
| `import_type` | `varchar(20)` | YES | NULL | delivery发货单 |
| `type` | `int(11)` | YES | '1' | 类型：1发货，2送货，3虚拟，4电子面单 |
| `count` | `int(11)` | YES | NULL | 总数 |
| `success` | `int(11)` | YES | NULL | 成功数 |
| `status` | `int(11)` | YES | '0' | 0.处理中，1成功，10部分完成，-1失败 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `mer_id` | `int(11)` | YES | NULL | — |

## `qixi_store_import_delivery`

> 导入发货单详细记录

原表：`eb_store_import_delivery`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `import_delivery_id` | `int(10)` | NO | — | — |
| `import_id` | `int(10)` | NO | — | — |
| `order_sn` | `varchar(32)` | YES | NULL | 订单sn |
| `delivery_type` | `int(11)` | YES | '1' | 类型：1发货，2送货，3虚拟，4电子面单 |
| `delivery_name` | `varchar(64)` | YES | NULL | 快递公司 |
| `delivery_id` | `varchar(64)` | YES | NULL | 快递单号 |
| `status` | `tinyint(4)` | YES | NULL | 状态 |
| `mark` | `varchar(255)` | YES | NULL | 备注 |
| `create_time` | `timestamp` | YES | CURRENT_TIMESTAMP | — |
| `mer_id` | `int(11)` | YES | NULL | — |

## `qixi_store_order`

> 订单表

原表：`eb_store_order`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `order_id` | `int(10)` | NO | — | 订单ID |
| `main_id` | `int(10)` | YES | '0' | 拆单前 id |
| `group_order_id` | `int(10)` | NO | '0' | 订单组 id |
| `order_sn` | `varchar(32)` | NO | — | 订单号 |
| `uid` | `int(10)` | NO | — | 用户id |
| `spread_uid` | `int(10)` | YES | '0' | 推荐人id |
| `top_uid` | `int(10)` | YES | '0' | 二级推荐人 id |
| `real_name` | `varchar(32)` | NO | — | 用户姓名 |
| `user_phone` | `varchar(18)` | NO | — | 用户电话 |
| `user_address` | `varchar(128)` | NO | — | 详细地址 |
| `cart_id` | `varchar(256)` | NO | — | 购物车id |
| `total_num` | `int(10)` | NO | '0' | 订单商品总数 |
| `total_price` | `decimal(8,2)` | NO | '0.00' | 订单总价 |
| `total_postage` | `decimal(8,2)` | NO | '0.00' | 邮费 |
| `pay_price` | `decimal(8,2)` | NO | '0.00' | 实际支付金额 |
| `pay_postage` | `decimal(8,2)` | NO | '0.00' | 支付邮费 |
| `is_selfbuy` | `tinyint(3)` | NO | '0' | 是否为自购 |
| `extension_one` | `decimal(8,2)` | NO | '0.00' | 一级佣金 |
| `extension_two` | `decimal(8,2)` | NO | '0.00' | 二级佣金 |
| `commission_rate` | `decimal(7,4)` | NO | '0.0000' | 平台手续费 |
| `integral` | `int(10)` | YES | '0' | 使用积分数量 |
| `integral_price` | `decimal(8,2)` | YES | '0.00' | 积分抵扣金额 |
| `give_integral` | `int(10)` | YES | '0' | 赠送积分 |
| `coupon_id` | `varchar(128)` | NO | '' | 优惠券id |
| `coupon_price` | `decimal(8,2)` | NO | '0.00' | 优惠券金额 |
| `platform_coupon_price` | `decimal(8,2)` | YES | '0.00' | 平台优惠券金额 |
| `svip_discount` | `decimal(8,2)` | YES | '0.00' | svip优惠金额 |
| `order_type` | `tinyint(3)` | YES | '0' | 0普通1自提 |
| `paid` | `tinyint(3)` | NO | '0' | 支付状态 |
| `pay_time` | `timestamp` | YES | NULL | 支付时间 |
| `pay_type` | `tinyint(1)` | NO | — | 支付方式 0余额 1微信 2小程序 3 h5  4支付宝 5 支付宝扫码 6 微信扫码 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `status` | `tinyint(1)` | NO | '0' | 订单状态（0：待发货；1：待收货；2：待评价；3：已完成； 9: 拼团中 10:  待付尾款 11:尾款超时未付 -1：已退款） |
| `delivery_type` | `varchar(32)` | YES | NULL | 发货类型(1:发货 2: 送货 3: 虚拟,4电子面单，5同城 6 卡密自动发货) |
| `is_virtual` | `tinyint(3)` | YES | '0' | 0:实物订单 1:虚拟订单 |
| `delivery_name` | `varchar(50)` | YES | NULL | 快递名称/送货人姓名 |
| `delivery_id` | `varchar(255)` | YES | NULL | 快递单号/手机号 |
| `mark` | `varchar(512)` | NO | — | 备注 |
| `remark` | `varchar(512)` | YES | NULL | 管理员备注 |
| `admin_mark` | `varchar(512)` | YES | NULL | 总后台备注 |
| `verify_code` | `char(16)` | YES | NULL | 核销码 |
| `verify_time` | `timestamp` | YES | NULL | 核销时间/收货时间 |
| `verify_service_id` | `int(10)` | YES | NULL | 核销客服 id |
| `transaction_id` | `varchar(60)` | YES | NULL | 微信支付订单号(分账时有效) |
| `activity_type` | `tinyint(3)` | NO | '0' | 1:秒杀 2:预售 3:助力 10:套餐 |
| `order_extend` | `varchar(1024)` | YES | NULL | 自定义表单数据 |
| `mer_id` | `int(10)` | NO | '0' | 商户ID |
| `reconciliation_id` | `tinyint(3)` | NO | '0' | 对账id |
| `cost` | `decimal(8,2)` | NO | — | 成本价 |
| `is_del` | `tinyint(3)` | NO | '0' | 是否删除 |
| `is_system_del` | `tinyint(1)` | YES | '0' | 后台是否删除 |
| `verify_status` | `tinyint(3)` | NO | '0' | 核销订单状态0 默认 1 部分核销 2 全部核销 |
| `refund_switch` | `tinyint(3)` | YES | '1' | 是否支持退款 |
| `kuaidi_label` | `varchar(255)` | YES | NULL | 快递单号图片 |
| `task_id` | `varchar(255)` | NO | '' | 快递单号任务ID |
| `is_behalf` | `tinyint(1)` | NO | '0' | 是否为代客下单订单 |
| `behalf_no_verify` | `tinyint(1)` | NO | '0' | 代客下单订单是否无需核销（0:需要核销，1:无需核销） |
| `enable_assigned` | `tinyint(1)` | YES | '0' | 预约商品订单是否派单 0 领取 1 指派 |
| `staffs_id` | `int(10)` | YES | '0' | 领取/指派的员工ID |
| `is_cancel` | `int(11)` | YES | '0' | 是否可取消预约(0:不可取消,1:可取消) |
| `reservation_service_voucher` | `text` | YES | — | 预约订单服务凭证 |
| `clock_in_info` | `text` | YES | — | 预约订单服务打卡信息 |
| `kuaidi_order_id` | `varchar(255)` | NO | '' | 商家寄件的快递单号 |
| `is_stock_up` | `int(11)` | NO | '0' | 商家寄件快递是否发出 |
| `merchant_take_id` | `int(11)` | NO | '0' | 自提点id |
| `merchant_take_info` | `text` | NO | — | 同城配送信息 |

## `qixi_store_order_product`

> 订单购物详情表

原表：`eb_store_order_product`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `order_product_id` | `int(10)` | NO | — | 订单产品 id |
| `order_id` | `int(10)` | NO | — | 订单id |
| `uid` | `int(10)` | NO | '0' | 用户 id |
| `cart_id` | `int(10)` | NO | '0' | 购物车id |
| `product_id` | `int(10)` | NO | '0' | 商品ID |
| `extension_one` | `decimal(8,2)` | NO | '0.00' | 一级佣金 |
| `extension_two` | `decimal(8,2)` | NO | '0.00' | 二级佣金 |
| `integral` | `int(10)` | YES | '0' | 使用积分(单数) |
| `integral_price` | `decimal(10,2)` | YES | '0.00' | 积分抵扣金额 |
| `integral_total` | `int(10)` | YES | '0' | 使用积分(总数) |
| `coupon_price` | `decimal(8,2)` | YES | '0.00' | 优惠金额 |
| `platform_coupon_price` | `decimal(8,2)` | YES | '0.00' | 平台优惠金额 |
| `svip_discount` | `decimal(8,2)` | YES | '0.00' | svip优惠金额 |
| `postage_price` | `decimal(8,2)` | YES | '0.00' | 运费 |
| `product_sku` | `char(12)` | NO | — | 商品 sku |
| `is_refund` | `tinyint(3)` | NO | '0' | 是否退款 0:未退款 1:退款中 2:部分退款 3=全退 |
| `product_num` | `int(10)` | NO | '0' | 购买数量 |
| `product_type` | `int(11)` | NO | '0' | 0.普通商品 1.秒杀商品,2.预售商品 |
| `activity_id` | `int(10)` | NO | '0' | 活动关联 id |
| `refund_num` | `int(10)` | NO | '0' | 可申请退货数量 |
| `is_reply` | `tinyint(3)` | NO | '0' | 是否评价 |
| `cost` | `decimal(10,2)` | NO | — | 商品成本价 |
| `product_price` | `decimal(10,2)` | NO | — | 商品金额 |
| `total_price` | `decimal(10,2)` | NO | — | 商品售价 |
| `cart_info` | `text` | NO | — | 购买东西的详细信息 |
| `refund_switch` | `tinyint(3)` | YES | '1' | 是否支持退款 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | — |
| `reservation_date` | `varchar(20)` | YES | '' | 预约商品的预约日期 |
| `reservation_id` | `int(10)` | YES | '0' | 预约商品的预约时间段ID |
| `reservation_time_part` | `varchar(20)` | NO | '' | 预约时间段 |
| `settlement_price` | `decimal(8,2)` | NO | '0.00' | 员工结算价格 |

## `qixi_store_order_profitsharing`

> 分账表

原表：`eb_store_order_profitsharing`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `profitsharing_id` | `int(10)` | NO | — | — |
| `profitsharing_sn` | `varchar(32)` | NO | — | 分账 id |
| `order_id` | `int(10)` | NO | — | 订单 id |
| `mer_id` | `int(10)` | NO | — | 商户 id |
| `transaction_id` | `varchar(60)` | NO | — | 微信支付订单号 |
| `profitsharing_price` | `decimal(10,2)` | NO | '0.00' | 分账金额 |
| `profitsharing_refund` | `decimal(10,2)` | NO | '0.00' | 退款金额 |
| `profitsharing_mer_price` | `decimal(10,2)` | NO | '0.00' | 分账分出去的金额即给平台的手续费 |
| `type` | `varchar(32)` | NO | — | 分类 |
| `status` | `tinyint(1)` | NO | '0' | 0:未分账 1:已分账 -1已退款 -2失败 |
| `error_msg` | `varchar(255)` | YES | NULL | 失败原因 |
| `profitsharing_time` | `timestamp` | YES | NULL | 分账时间 |
| `create_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 创建时间 |
| `is_combine` | `tinyint(1)` | YES | '1' | 分账类型：1 平台收付通 2 服务商 |

## `qixi_store_order_receipt`

> 订单发票信息

原表：`eb_store_order_receipt`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `order_receipt_id` | `int(10)` | NO | — | — |
| `order_id` | `varchar(255)` | NO | '0' | 订单ID |
| `uid` | `int(11)` | NO | '0' | 用户ID |
| `receipt_info` | `varchar(500)` | YES | '' | 发票类型：1.普通发票，2.增值税发票 |
| `status` | `tinyint(4)` | YES | '0' | 开票状态：1.已出票,10.已寄出 |
| `receipt_sn` | `varchar(255)` | YES | '' | 发票单号 |
| `receipt_no` | `varchar(255)` | YES | NULL | 发票编号 |
| `delivery_info` | `varchar(255)` | YES | NULL | 收票联系信息 |
| `pic` | `varchar(500)` | NO | '' | 发票文件地址 |
| `mark` | `varchar(255)` | YES | NULL | 用户备注 |
| `receipt_price` | `decimal(10,2)` | YES | NULL | 开票金额 |
| `order_price` | `decimal(10,2)` | YES | NULL | 订单金额 |
| `status_time` | `datetime` | NO | — | 状态变更时间 |
| `is_del` | `tinyint(1)` | YES | '0' | — |
| `create_time` | `timestamp` | YES | NULL | — |
| `mer_id` | `int(11)` | YES | '0' | — |
| `mer_mark` | `varchar(255)` | YES | NULL | 备注 |

## `qixi_store_order_status`

> 订单操作记录表

原表：`eb_store_order_status`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `order_id` | `int(10)` | NO | — | 订单id |
| `order_sn` | `varchar(256)` | YES | NULL | 订单号 |
| `type` | `varchar(20)` | YES | NULL | 订单类型 |
| `change_type` | `varchar(32)` | NO | — | 操作类型 |
| `change_message` | `varchar(256)` | NO | — | 操作备注 |
| `nickname` | `varchar(20)` | YES | NULL | — |
| `uid` | `int(11)` | YES | NULL | 操作者ID |
| `user_type` | `tinyint(4)` | YES | NULL | 操作者类型 |
| `change_time` | `timestamp` | NO | CURRENT_TIMESTAMP | 操作时间 |

## `qixi_user_order`

> 支付订单信息

原表：`eb_user_order`

| 字段 | 类型 | 空 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `order_id` | `int(10)` | NO | — | — |
| `link_id` | `int(11)` | YES | NULL | 关联ID |
| `pay_type` | `varchar(10)` | YES | NULL | 支付方式：1微信，2支付宝 |
| `title` | `varchar(50)` | NO | — | — |
| `order_sn` | `varchar(50)` | YES | NULL | 订单ID |
| `pay_price` | `decimal(8,2)` | YES | NULL | 价格 |
| `order_info` | `varchar(255)` | YES | NULL | 订单信息信息 |
| `order_type` | `varchar(255)` | YES | NULL | 订单类型 S 付费会员  |
| `paid` | `int(11)` | NO | '0' | 支付状态 |
| `pay_time` | `timestamp` | YES | NULL | 支付时间 |
| `status` | `int(11)` | NO | '0' | 状态：默认1 |
| `mer_id` | `int(11)` | YES | NULL | 商户ID |
| `uid` | `int(11)` | NO | '0' | 用户ID |
| `create_time` | `timestamp` | YES | NULL | — |
| `is_del` | `int(11)` | NO | '0' | — |
| `admin_id` | `int(10)` | YES | NULL | 管理员ID |
| `other` | `varchar(50)` | YES | NULL | 其他参数 |
| `end_time` | `timestamp` | YES | NULL | — |
| `transaction_id` | `varchar(255)` | YES | '' | — |
