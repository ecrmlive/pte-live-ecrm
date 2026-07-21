# 表领域分组（`qixi_`）

## 优惠券 (`coupon`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_store_coupon` | `eb_store_coupon` | 优惠券表 | 22 |
| `qixi_store_coupon_issue_user` | `eb_store_coupon_issue_user` | 优惠券前台用户领取记录表 | 3 |
| `qixi_store_coupon_product` | `eb_store_coupon_product` | 优惠卷关联商品辅助表 | 3 |
| `qixi_store_coupon_send` | `eb_store_coupon_send` | 优惠券发送记录 | 7 |
| `qixi_store_coupon_user` | `eb_store_coupon_user` | 优惠券发放记录表 | 15 |

## 其他 (`other`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_store_printer` | `eb_store_printer` | — | 13 |
| `qixi_wechat_qrcode` | `eb_wechat_qrcode` | 微信二维码管理表 | 10 |
| `qixi_wechat_reply` | `eb_wechat_reply` | 微信关键字回复表 | 7 |
| `qixi_wechat_user` | `eb_wechat_user` | 微信用户表 | 19 |

## 内容/社区/圈子 (`content`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_article` | `eb_article` | 文章管理表 | 16 |
| `qixi_article_category` | `eb_article_category` | 文章分类表 | 9 |
| `qixi_article_content` | `eb_article_content` | 文章内容表 | 2 |
| `qixi_circle` | `eb_circle` | — | 17 |
| `qixi_circle_agent` | `eb_circle_agent` | 商圈代理表 | 23 |
| `qixi_circle_brokerage_checkout` | `eb_circle_brokerage_checkout` | 商圈佣金结算表 | 21 |
| `qixi_community` | `eb_community` | 社区图文表信息 | 24 |
| `qixi_community_category` | `eb_community_category` | 社区分类 | 7 |
| `qixi_community_reply` | `eb_community_reply` | 社区评论 | 12 |
| `qixi_community_topic` | `eb_community_topic` | 社区话题 | 11 |
| `qixi_feedback` | `eb_feedback` | 用户反馈表 | 13 |
| `qixi_feedback_category` | `eb_feedback_category` | 用户反馈分类表 | 10 |
| `qixi_wechat_news` | `eb_wechat_news` | 图文消息管理表 | 4 |

## 售后 (`aftersale`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_store_refund_order` | `eb_store_refund_order` | 订单退款表 | 35 |
| `qixi_store_refund_product` | `eb_store_refund_product` | 退款单产品表 | 9 |
| `qixi_store_refund_status` | `eb_store_refund_status` | 订单操作记录表 | 4 |

## 商品/类目/品牌 (`catalog`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_cdkey_library` | `eb_cdkey_library` | — | 10 |
| `qixi_guarantee` | `eb_guarantee` | 保障服务选项 | 11 |
| `qixi_guarantee_template` | `eb_guarantee_template` | 保障服务模板 | 7 |
| `qixi_guarantee_value` | `eb_guarantee_value` | 保障服务模板条款 | 5 |
| `qixi_parameter` | `eb_parameter` | — | 8 |
| `qixi_parameter_product` | `eb_parameter_product` | — | 3 |
| `qixi_parameter_template` | `eb_parameter_template` | — | 5 |
| `qixi_parameter_value` | `eb_parameter_value` | — | 8 |
| `qixi_price_rule` | `eb_price_rule` | — | 8 |
| `qixi_store_attr_template` | `eb_store_attr_template` | 商品规则值(规格)表 | 4 |
| `qixi_store_brand` | `eb_store_brand` | 商品品牌表 | 7 |
| `qixi_store_brand_category` | `eb_store_brand_category` | 品牌分类表 | 8 |
| `qixi_store_category` | `eb_store_category` | 商品分类表 | 12 |
| `qixi_store_product` | `eb_store_product` | 商品表 | 71 |
| `qixi_store_product_assist` | `eb_store_product_assist` | 商品助力活动表 | 17 |
| `qixi_store_product_assist_set` | `eb_store_product_assist_set` | 助力发起列表 | 13 |
| `qixi_store_product_assist_sku` | `eb_store_product_assist_sku` | — | 6 |
| `qixi_store_product_assist_user` | `eb_store_product_assist_user` | 助力记录表 | 7 |
| `qixi_store_product_attr` | `eb_store_product_attr` | 商品属性表 | 4 |
| `qixi_store_product_attr_reservation` | `eb_store_product_attr_reservation` | — | 10 |
| `qixi_store_product_attr_result` | `eb_store_product_attr_result` | — | 5 |
| `qixi_store_product_attr_value` | `eb_store_product_attr_value` | 商品属性值表 | 23 |
| `qixi_store_product_cate` | `eb_store_product_cate` | 商品商户分类关联表 | 3 |
| `qixi_store_product_cdkey` | `eb_store_product_cdkey` | — | 11 |
| `qixi_store_product_content` | `eb_store_product_content` | 商品详情表 | 3 |
| `qixi_store_product_copy` | `eb_store_product_copy` | — | 8 |
| `qixi_store_product_group` | `eb_store_product_group` | 拼团商品信息表 | 23 |
| `qixi_store_product_group_buying` | `eb_store_product_group_buying` | 拼团活动表 | 13 |
| `qixi_store_product_group_sku` | `eb_store_product_group_sku` | — | 6 |
| `qixi_store_product_group_user` | `eb_store_product_group_user` | 拼团成员表 | 11 |
| `qixi_store_product_label` | `eb_store_product_label` | — | 9 |
| `qixi_store_product_presell` | `eb_store_product_presell` | 商品预售活动表 | 21 |
| `qixi_store_product_presell_sku` | `eb_store_product_presell_sku` | — | 12 |
| `qixi_store_product_reply` | `eb_store_product_reply` | 商品评论表 | 22 |
| `qixi_store_product_reservation` | `eb_store_product_reservation` | — | 22 |
| `qixi_store_product_sku` | `eb_store_product_sku` | — | 10 |
| `qixi_store_product_take` | `eb_store_product_take` | 用户到货通知记录 | 7 |
| `qixi_store_product_unit` | `eb_store_product_unit` | — | 7 |
| `qixi_store_spu` | `eb_store_spu` | 商品搜索信息表 | 19 |

## 商户 (`merchant`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_merchant` | `eb_merchant` | 商户表 | 55 |
| `qixi_merchant_admin` | `eb_merchant_admin` | 商户管理员表 | 14 |
| `qixi_merchant_applyments` | `eb_merchant_applyments` | 商户申请分账商户号表 | 13 |
| `qixi_merchant_category` | `eb_merchant_category` | 商户分类表 | 4 |
| `qixi_merchant_intention` | `eb_merchant_intention` | 商户申请表 | 16 |
| `qixi_merchant_region` | `eb_merchant_region` | — | 13 |
| `qixi_merchant_type` | `eb_merchant_type` | 商户类型表 | 9 |
| `qixi_store_group` | `eb_store_group` | 分组表 | 15 |

## 客服 (`cs`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_store_service` | `eb_store_service` | 客服表 | 17 |
| `qixi_store_service_log` | `eb_store_service_log` | 客服用户对话记录表 | 11 |
| `qixi_store_service_reply` | `eb_store_service_reply` | — | 7 |
| `qixi_store_service_user` | `eb_store_service_user` | — | 10 |

## 用户/资产/分销 (`user`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_member_interests` | `eb_member_interests` | — | 11 |
| `qixi_user` | `eb_user` | 用户表 | 45 |
| `qixi_user_address` | `eb_user_address` | 用户地址表 | 20 |
| `qixi_user_bill` | `eb_user_bill` | 用户账单表 | 13 |
| `qixi_user_brokerage` | `eb_user_brokerage` | — | 10 |
| `qixi_user_extract` | `eb_user_extract` | 用户提现表 | 24 |
| `qixi_user_fields` | `eb_user_fields` | — | 2 |
| `qixi_user_group` | `eb_user_group` | 用户分组表 | 3 |
| `qixi_user_history` | `eb_user_history` | 浏览记录表 | 6 |
| `qixi_user_info` | `eb_user_info` | — | 13 |
| `qixi_user_label` | `eb_user_label` | 用户标签表 | 5 |
| `qixi_user_merchant` | `eb_user_merchant` | 商户用户表 | 11 |
| `qixi_user_receipt` | `eb_user_receipt` | 用户发票信息 | 14 |
| `qixi_user_recharge` | `eb_user_recharge` | 用户充值表 | 10 |
| `qixi_user_relation` | `eb_user_relation` | 用户记录表 | 4 |
| `qixi_user_sign` | `eb_user_sign` | 签到记录表 | 7 |
| `qixi_user_spread_log` | `eb_user_spread_log` | — | 6 |
| `qixi_user_visit` | `eb_user_visit` | 商品浏览分析表 | 6 |

## 直播 (`live`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_broadcast_assistant` | `eb_broadcast_assistant` | 直播助手信息 | 6 |
| `qixi_broadcast_goods` | `eb_broadcast_goods` | 直播商品表 | 19 |
| `qixi_broadcast_room` | `eb_broadcast_room` | 直播间表 | 34 |
| `qixi_broadcast_room_goods` | `eb_broadcast_room_goods` | 直播间导入商品表 | 3 |

## 系统/权限/配置/DIY/素材 (`system`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_cache` | `eb_cache` | 微信缓存表 | 4 |
| `qixi_city_area` | `eb_city_area` | 省市区县数据 | 9 |
| `qixi_diy` | `eb_diy` | — | 23 |
| `qixi_excel` | `eb_excel` | 导出文件记录表 | 10 |
| `qixi_extend` | `eb_extend` | — | 7 |
| `qixi_label_rule` | `eb_label_rule` | 自定标签规则 | 9 |
| `qixi_open_auth` | `eb_open_auth` | — | 15 |
| `qixi_operate_log` | `eb_operate_log` | — | 15 |
| `qixi_page_category` | `eb_page_category` | 页面链接分类 | 9 |
| `qixi_page_link` | `eb_page_link` | 页面链接 | 11 |
| `qixi_record` | `eb_record` | — | 6 |
| `qixi_relevance` | `eb_relevance` | — | 4 |
| `qixi_routine_qrcode` | `eb_routine_qrcode` | 小程序二维码管理表 | 8 |
| `qixi_sms_record` | `eb_sms_record` | 短信发送记录表 | 9 |
| `qixi_staffs` | `eb_staffs` | — | 11 |
| `qixi_system_admin` | `eb_system_admin` | 后台管理员表 | 17 |
| `qixi_system_attachment` | `eb_system_attachment` | 附件管理表 | 9 |
| `qixi_system_attachment_category` | `eb_system_attachment_category` | 附件分类表 | 8 |
| `qixi_system_config` | `eb_system_config` | 配置表 | 16 |
| `qixi_system_config_classify` | `eb_system_config_classify` | 配置分类表 | 9 |
| `qixi_system_config_value` | `eb_system_config_value` | 配置表 | 5 |
| `qixi_system_form` | `eb_system_form` | — | 9 |
| `qixi_system_group` | `eb_system_group` | 组合数据表 | 8 |
| `qixi_system_group_data` | `eb_system_group_data` | 组合数据详情表 | 7 |
| `qixi_system_log` | `eb_system_log` | 管理员操作记录表 | 9 |
| `qixi_system_menu` | `eb_system_menu` | 菜单表 | 14 |
| `qixi_system_notice` | `eb_system_notice` | 商户公告 | 9 |
| `qixi_system_notice_config` | `eb_system_notice_config` | — | 21 |
| `qixi_system_notice_log` | `eb_system_notice_log` | — | 6 |
| `qixi_system_role` | `eb_system_role` | 身份管理表 | 10 |
| `qixi_system_storage` | `eb_system_storage` | — | 12 |
| `qixi_template_message` | `eb_template_message` | 微信模板 | 9 |

## 营销商品(秒杀/拼团/预售/助力) (`marketing_product`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_presell_order` | `eb_presell_order` | — | 16 |
| `qixi_store_activity` | `eb_store_activity` | — | 21 |
| `qixi_store_activity_cate` | `eb_store_activity_cate` | 活动标签分类表 | 11 |
| `qixi_store_activity_label` | `eb_store_activity_label` | 活动标签表 | 15 |
| `qixi_store_activity_related` | `eb_store_activity_related` | — | 13 |
| `qixi_store_discounts` | `eb_store_discounts` | — | 19 |
| `qixi_store_discounts_product` | `eb_store_discounts_product` | — | 8 |
| `qixi_store_seckill_active` | `eb_store_seckill_active` | 商户设置秒杀商品关联表 | 20 |
| `qixi_store_seckill_time` | `eb_store_seckill_time` | 秒杀时间段配置 | 7 |

## 财务/结算 (`finance`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_circle_financial_record` | `eb_circle_financial_record` | 商圈提成流水 | 16 |
| `qixi_financial` | `eb_financial` | 商户财务申请提现 | 20 |
| `qixi_financial_record` | `eb_financial_record` | 商户财务流水 | 13 |
| `qixi_serve_meal` | `eb_serve_meal` | — | 9 |
| `qixi_serve_order` | `eb_serve_order` | — | 12 |

## 购物车/订单 (`cart_order`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_store_cart` | `eb_store_cart` | 购物车表 | 18 |
| `qixi_store_cart_price` | `eb_store_cart_price` | — | 10 |
| `qixi_store_group_order` | `eb_store_group_order` | 用户订单表 | 28 |
| `qixi_store_import` | `eb_store_import` | 导入批次记录 | 8 |
| `qixi_store_import_delivery` | `eb_store_import_delivery` | 导入发货单详细记录 | 10 |
| `qixi_store_order` | `eb_store_order` | 订单表 | 66 |
| `qixi_store_order_product` | `eb_store_order_product` | 订单购物详情表 | 31 |
| `qixi_store_order_profitsharing` | `eb_store_order_profitsharing` | 分账表 | 14 |
| `qixi_store_order_receipt` | `eb_store_order_receipt` | 订单发票信息 | 17 |
| `qixi_store_order_status` | `eb_store_order_status` | 订单操作记录表 | 9 |
| `qixi_user_order` | `eb_user_order` | 支付订单信息 | 19 |

## 配送/运费 (`fulfillment`)

| qixi 表 | CRMEB 原表 | 说明 | 字段数 |
| --- | --- | --- | ---: |
| `qixi_delivery_config` | `eb_delivery_config` | 配送设置表 | 15 |
| `qixi_delivery_order` | `eb_delivery_order` | — | 28 |
| `qixi_delivery_service` | `eb_delivery_service` | — | 13 |
| `qixi_delivery_station` | `eb_delivery_station` | 同城配送门店列表 | 28 |
| `qixi_express` | `eb_express` | 快递公司表 | 13 |
| `qixi_express_partner` | `eb_express_partner` | — | 10 |
| `qixi_shipping_template` | `eb_shipping_template` | 运费表 | 10 |
| `qixi_shipping_template_free` | `eb_shipping_template_free` | 指定包邮信息表 | 5 |
| `qixi_shipping_template_region` | `eb_shipping_template_region` | 配送区域表 | 7 |
| `qixi_shipping_template_undelivery` | `eb_shipping_template_undelivery` | 指定不配送区域表 | 3 |
