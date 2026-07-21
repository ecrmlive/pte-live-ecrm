# 功能点缺口清单（需人工补全/确认）

> 基于菜单按钮 + 路由补全后的结果。
> 「无任何操作」必须补全；「缺少某类 CRUD」需结合业务判断（只读配置页可能无 C/D）。

## 无任何操作的页面（53）

| 端 | 模块路径 | 页面路由 |
| --- | --- | --- |
| 商户后台 | 员工 | `/server` |
| 商户后台 | 员工 / 店员管理 | `/server_manage` |
| 商户后台 | 员工 / 店员管理 / 店员配置 | `/systemForm/Basics/mer_service_switch` |
| 商户后台 | 员工 / 配送人员 / 配送统计 | `/delivery/delivice_statistic` |
| 商户后台 | 用户 / 标签管理 | `/user/_label` |
| 商户后台 | 营销 / 助力 | `/assist` |
| 商户后台 | 营销 / 直播 / 直播商品管理 | `/marketing/broadcast/list` |
| 商户后台 | 营销 / 秒杀 | `/marketing/seckill/list` |
| 商户后台 | 营销 / 积分 / 积分配置 | `/marketing/integral/config` |
| 商户后台 | 订单 / 预约设置 | `/product/reservation` |
| 商户后台 | 设置 / 一号通 | `/one_setting` |
| 商户后台 | 设置 / 一号通 / 自有一号通 | `/setting/sms/sms_account/index` |
| 商户后台 | 设置 / 付费会员 | `/systemForm/Basics/svip` |
| 商户后台 | 设置 / 快递配送 | `/city` |
| 商户后台 | 设置 / 打印配置 / 打印配置 | `/systemForm/Basics/printer_tabs` |
| 平台后台 | 内容 | `/content` |
| 平台后台 | 分销 / 分销等级 | `brokerage` |
| 平台后台 | 分销 / 分销说明 | `/promoter/retail` |
| 平台后台 | 商品 / 品牌管理 | `/product/brand` |
| 平台后台 | 商品 / 活动标签 | `/product/activityLabel` |
| 平台后台 | 应用 | `/apploction` |
| 平台后台 | 应用 / 公众号 / 自动回复 / 关键字回复 | `/app/wechat/reply/keyword` |
| 平台后台 | 应用 / 公众号 / 自动回复 / 微信关注回复 | `/app/wechat/reply/follow/subscribe` |
| 平台后台 | 应用 / 公众号 / 自动回复 / 无效关键词回复 | `/app/wechat/reply/index/default` |
| 平台后台 | 店铺 / 区域代理 | `/business-zones/manage` |
| 平台后台 | 店铺 / 区域代理 / 代理设置 | `/business-zones/settings` |
| 平台后台 | 店铺 / 商户管理 / 商户设置 | `/merchant/apply-setting` |
| 平台后台 | 店铺 / 店铺管理 | `/mer/mer` |
| 平台后台 | 店铺 / 店铺设置 | `/mer/store` |
| 平台后台 | 用户 / 付费会员 | `/user/svip` |
| 平台后台 | 用户 / 付费会员 / 会员协议 | `/user/member/vipAgreement` |
| 平台后台 | 维护 | `/safe` |
| 平台后台 | 维护 / 开发配置 | `/safe/exploit` |
| 平台后台 | 维护 / 页面链接 | `/safe/pageLinks` |
| 平台后台 | 营销 / 余额充值 | `/banlace` |
| 平台后台 | 营销 / 助力 | `/assist` |
| 平台后台 | 营销 / 直播 | `/marketing2` |
| 平台后台 | 装修 | `/theme` |
| 平台后台 | 装修 / 个人中心 | `/setting/diy/personal` |
| 平台后台 | 装修 / 商品详情 | `/setting/diy/product_detail` |
| 平台后台 | 装修 / 店铺街 | `/setting/diy/store` |
| 平台后台 | 装修 / 页面链接 | `/setting/page` |
| 平台后台 | 设置 | `/settings` |
| 平台后台 | 设置 / 商城设置 | `/shop` |
| 平台后台 | 设置 / 商城设置 / 配送配置 | `/delivery_config` |
| 平台后台 | 设置 / 商城设置 / 配送配置 / 城市数据 | `/freight/city/list` |
| 平台后台 | 设置 / 应用配置 | `/app_config` |
| 平台后台 | 设置 / 消息管理 | `/notice` |
| 平台后台 | 设置 / 系统设置 / 一号通 | `/serve` |
| 平台后台 | 设置 / 系统设置 / 一号通 / 短信设置 / 短信账户 | `/sms/user` |
| 平台后台 | 财务 / 发票管理 | `/accounts/accounts` |
| 平台后台 | 财务 / 店铺结算 | `/mer/accounts` |
| 平台后台 | 财务 / 用户结算 | `/accounts/record` |

## CRUD 不完整的页面（169，供核对）

| 端 | 模块路径 | 页面路由 | 缺口 | 已有 | 操作数 |
| --- | --- | --- | --- | --- | ---: |
| 商户后台 | 公告 | `/station/notice` | 缺少:CU | 已有:RD | 5 |
| 商户后台 | 员工 / 服务人员 / 服务人员 | `/config/service_staff` | 缺少:CUD | 已有:R | 2 |
| 商户后台 | 员工 / 服务人员 / 服务统计 | `/config/service_statistic` | 缺少:D | 已有:CRU | 5 |
| 商户后台 | 商品 | `/product` | 缺少:UD | 已有:CR | 2 |
| 商户后台 | 用户 | `/user` | 缺少:D | 已有:CRU | 6 |
| 商户后台 | 用户 / 搜索记录 | `/user/searchRecord` | 缺少:CUD | 已有:R | 1 |
| 商户后台 | 用户 / 用户管理 | `/user/list` | 缺少:D | 已有:CRU | 6 |
| 商户后台 | 营销 | `/marketing` | 缺少:UD | 已有:CR | 4 |
| 商户后台 | 营销 / 优惠券 | `/marketing/coupon` | 缺少:UD | 已有:CR | 4 |
| 商户后台 | 营销 / 优惠券 / 发送记录 | `/marketing/coupon/send` | 缺少:UD | 已有:CR | 1 |
| 商户后台 | 营销 / 优惠券 / 领取记录 | `/marketing/coupon/user` | 缺少:CUD | 已有:R | 1 |
| 商户后台 | 营销 / 助力 / 助力活动 | `/marketing/assist/assist_set` | 缺少:CUD | 已有:R | 2 |
| 商户后台 | 营销 / 拼团 | `/marketing/combination` | 缺少:RUD | 已有:C | 2 |
| 商户后台 | 营销 / 拼团 / 拼团活动列表 | `/marketing/combination/combination_list` | 缺少:UD | 已有:CR | 4 |
| 商户后台 | 营销 / 直播 | `/` | 缺少:UD | 已有:CR | 5 |
| 商户后台 | 营销 / 秒杀 / 秒杀活动 | `/marketing/seckill/store_seckill/list` | 缺少:UD | 已有:CR | 9 |
| 商户后台 | 营销 / 积分 | `/marketing/integral` | 缺少:CUD | 已有:R | 2 |
| 商户后台 | 营销 / 积分 / 积分日志 | `/marketing/integral/log` | 缺少:D | 已有:CRU | 4 |
| 商户后台 | 装修 / 商品分类 | `/devise/diy/product_category` | 缺少:UD | 已有:CR | 2 |
| 商户后台 | 订单 | `/order` | 缺少:RD | 已有:CU | 2 |
| 商户后台 | 订单 / 商品评价 | `/product/reviews` | 缺少:D | 已有:CRU | 6 |
| 商户后台 | 订单 / 核销记录 | `/order/cancellation` | 缺少:CUD | 已有:R | 2 |
| 商户后台 | 订单 / 预约服务 | `/order/reservation` | 缺少:CRUD | 已有: | 1 |
| 商户后台 | 设置 / 一号通 / 平台一号通 | `/setting/sms/sms_config/index` | 缺少:CUD | 已有:R | 5 |
| 商户后台 | 设置 / 一号通 / 配置管理 | `/setting/sms/dumpConfig` | 缺少:RD | 已有:CU | 1 |
| 商户后台 | 设置 / 同城配送 | `/delivery` | 缺少:CUD | 已有:R | 3 |
| 商户后台 | 设置 / 同城配送 / 充值记录 | `/delivery/recharge_record` | 缺少:UD | 已有:CR | 1 |
| 商户后台 | 设置 / 同城配送 / 配送记录 | `/delivery/usage_record` | 缺少:CU | 已有:RD | 3 |
| 商户后台 | 设置 / 同城配送 / 配送设置 | `/setting/delivery` | 缺少:CRD | 已有:U | 2 |
| 商户后台 | 设置 / 店铺信息 | `/systemForm/modifyStoreInfo` | 缺少:RD | 已有:CU | 2 |
| 商户后台 | 设置 / 店铺配置 | `/systemForm/Basics/mer_base` | 缺少:CRD | 已有:U | 1 |
| 商户后台 | 设置 / 快递配送 / 物流公司 | `/config/freight/express` | 缺少:CD | 已有:RU | 3 |
| 商户后台 | 设置 / 打印配置 | `/setting/printer` | 缺少:CUD | 已有:R | 8 |
| 商户后台 | 设置 / 权限管理 | `/setting` | 缺少:CUD | 已有:R | 1 |
| 商户后台 | 设置 / 权限管理 / 操作日志 | `/setting/systemLog` | 缺少:CD | 已有:RU | 4 |
| 商户后台 | 财务 | `/accounts` | 缺少:U | 已有:CRD | 8 |
| 商户后台 | 财务 / 分账管理 | `/systemForm/applyList` | 缺少:CUD | 已有:R | 4 |
| 商户后台 | 财务 / 发票管理 | `/order/invoice` | 缺少:D | 已有:CRU | 6 |
| 商户后台 | 财务 / 收款方式 | `/accounts/payType` | 缺少:RUD | 已有:C | 1 |
| 商户后台 | 财务 / 申请分账商户 | `/systemForm/applyments` | 缺少:D | 已有:CRU | 5 |
| 商户后台 | 财务 / 账单管理 | `/accounts/statement` | 缺少:CUD | 已有:R | 6 |
| 商户后台 | 财务 / 资金流水 | `/accounts/capitalFlow` | 缺少:CUD | 已有:R | 5 |
| 商户后台 | 首页 | `/dashboard` | 缺少:CD | 已有:RU | 10 |
| 商户后台 | 首页 / 商品统计 | `/statistic/product` | 缺少:CUD | 已有:R | 3 |
| 商户后台 | 首页 / 控制台 | `/dashboard` | 缺少:U | 已有:CRD | 5 |
| 商户后台 | 首页 / 订单统计 | `/statistic/order` | 缺少:CUD | 已有:R | 3 |
| 平台后台 | 内容 / 社区 / 社区内容 | `/community/list` | 缺少:C | 已有:RUD | 7 |
| 平台后台 | 内容 / 社区 / 社区评论 | `/community/reply` | 缺少:C | 已有:RUD | 4 |
| 平台后台 | 内容 / 社区 / 社区配置 | `/systemForm/Basics/community` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 分销 | `/promoter` | 缺少:CUD | 已有:R | 8 |
| 平台后台 | 分销 / 佣金说明 | `/promoter/commission` | 缺少:RUD | 已有:C | 1 |
| 平台后台 | 分销 / 分销礼包 | `/promoter/gift` | 缺少:CD | 已有:RU | 6 |
| 平台后台 | 分销 / 分销等级 / 等级规则 | `/promoter/distribution` | 缺少:RUD | 已有:C | 1 |
| 平台后台 | 分销 / 分销订单 | `/promoter/orderList` | 缺少:CD | 已有:RU | 8 |
| 平台后台 | 分销 / 分销配置 | `/systemForm/Basics/distribution_tabs` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 商品 | `/product` | 缺少:U | 已有:CRD | 5 |
| 平台后台 | 商品 / 商品参数 / 店铺商品参数 | `/product/merSpecs` | 缺少:CUD | 已有:R | 1 |
| 平台后台 | 商品 / 商品管理 | `/product/examine` | 缺少:D | 已有:CRU | 48 |
| 平台后台 | 客服 | `/service` | 缺少:D | 已有:CRU | 4 |
| 平台后台 | 客服 / 客服设置 | `/systemForm/Basics/service` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 应用 / 公众号 / 微信菜单 | `/app/wechat/menus` | 缺少:RD | 已有:CU | 2 |
| 平台后台 | 应用 / 小程序 | `/app/routine` | 缺少:CUD | 已有:R | 2 |
| 平台后台 | 应用 / 小程序 / 小程序下载 | `/app/routine/download` | 缺少:CUD | 已有:R | 1 |
| 平台后台 | 店铺 | `/mer` | 缺少:UD | 已有:CR | 2 |
| 平台后台 | 店铺 / 区域代理 / 代理审核 | `/business-zones/agent-review` | 缺少:CRD | 已有:U | 1 |
| 平台后台 | 店铺 / 商户管理 | `/merchant` | 缺少:UD | 已有:CR | 2 |
| 平台后台 | 店铺 / 商户管理 / 商户入驻审核 | `/merchant/review` | 缺少:CRD | 已有:U | 1 |
| 平台后台 | 店铺 / 店铺管理 / 店铺入驻申请 | `/merchant/application` | 缺少:C | 已有:RUD | 8 |
| 平台后台 | 店铺 / 店铺管理 / 店铺分账申请 | `/merchant/applyments` | 缺少:D | 已有:CRU | 4 |
| 平台后台 | 店铺 / 店铺设置 / 保证金配置 | `/systemForm/Basics/margin` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 店铺 / 店铺设置 / 说明提示 | `/merchant/type/description` | 缺少:RUD | 已有:C | 1 |
| 平台后台 | 用户 | `/user` | 缺少:UD | 已有:CR | 7 |
| 平台后台 | 用户 / 付费会员 / 付费会员配置 | `/systemForm/Basics/svip` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 用户 / 付费会员 / 会员权益 | `/user/member/equity` | 缺少:D | 已有:CRU | 5 |
| 平台后台 | 用户 / 付费会员 / 会员记录 | `/user/member/record` | 缺少:CD | 已有:RU | 4 |
| 平台后台 | 用户 / 搜索记录 | `/user/searchRecord` | 缺少:CU | 已有:RD | 5 |
| 平台后台 | 用户 / 用户列表 | `/user/list` | 缺少:D | 已有:CRU | 28 |
| 平台后台 | 用户 / 用户协议 | `/user/agreement` | 缺少:RUD | 已有:C | 1 |
| 平台后台 | 用户 / 用户反馈 | `/feedback` | 缺少:C | 已有:RUD | 5 |
| 平台后台 | 用户 / 用户反馈 / 反馈列表 | `/feedback/list` | 缺少:C | 已有:RUD | 4 |
| 平台后台 | 用户 / 用户等级 | `/user/member` | 缺少:CR | 已有:UD | 2 |
| 平台后台 | 用户 / 用户等级 / 等级说明 | `/user/member/description` | 缺少:RUD | 已有:C | 1 |
| 平台后台 | 用户 / 用户等级 / 等级配置 | `/systemForm/Basics/members` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 维护 / 安全维护 | `/maintain` | 缺少:CU | 已有:RD | 2 |
| 平台后台 | 维护 / 安全维护 / 商业授权 | `/setting/system/maintain/auth` | 缺少:UD | 已有:CR | 3 |
| 平台后台 | 维护 / 安全维护 / 数据备份 | `/maintain/dataBackup` | 缺少:CU | 已有:RD | 8 |
| 平台后台 | 维护 / 安全维护 / 缓存清除 | `/maintain/cache` | 缺少:CRU | 已有:D | 2 |
| 平台后台 | 维护 / 导出记录 | `/group/exportList` | 缺少:CUD | 已有:R | 2 |
| 平台后台 | 维护 / 操作日志 | `/setting/systemLog` | 缺少:CUD | 已有:R | 1 |
| 平台后台 | 营销 | `/marketing` | 缺少:CUD | 已有:R | 2 |
| 平台后台 | 营销 / 余额充值 / 余额设置 | `/systemForm/Basics/balance` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 营销 / 助力 / 助力活动 | `/marketing/assist/list` | 缺少:CUD | 已有:R | 6 |
| 平台后台 | 营销 / 助力 / 活动商品 | `/marketing/assist/goods_list` | 缺少:CD | 已有:RU | 21 |
| 平台后台 | 营销 / 商户优惠券 | `/marketing/coupon` | 缺少:CUD | 已有:R | 2 |
| 平台后台 | 营销 / 商户优惠券 / 优惠券列表 | `/marketing/coupon/list` | 缺少:CUD | 已有:R | 3 |
| 平台后台 | 营销 / 商户优惠券 / 领取记录 | `/marketing/coupon/user` | 缺少:CUD | 已有:R | 1 |
| 平台后台 | 营销 / 平台优惠券 | `/marketing/platform_coupon` | 缺少:UD | 已有:CR | 3 |
| 平台后台 | 营销 / 平台优惠券 / 使用说明 | `/marketing/platform_coupon/instructions` | 缺少:RUD | 已有:C | 1 |
| 平台后台 | 营销 / 平台优惠券 / 发送记录 | `/marketing/platform_coupon/couponSend` | 缺少:UD | 已有:CR | 1 |
| 平台后台 | 营销 / 平台优惠券 / 领取记录 | `/marketing/platform_coupon/couponRecord` | 缺少:CUD | 已有:R | 1 |
| 平台后台 | 营销 / 拼团 | `/marketing/combination` | 缺少:RD | 已有:CU | 2 |
| 平台后台 | 营销 / 拼团 / 拼团商品列表 | `/marketing/combination/combination_goods` | 缺少:CD | 已有:RU | 24 |
| 平台后台 | 营销 / 拼团 / 拼团活动列表 | `/marketing/combination/combination_list` | 缺少:CUD | 已有:R | 6 |
| 平台后台 | 营销 / 拼团 / 拼团设置 | `/marketing/combination/combination_set` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 营销 / 直播 / 直播商品管理 | `/marketing/broadcast/list` | 缺少:C | 已有:RUD | 6 |
| 平台后台 | 营销 / 秒杀 / 秒杀管理 | `/marketing/seckill/list` | 缺少:C | 已有:RUD | 29 |
| 平台后台 | 营销 / 积分 | `/marketing/integral` | 缺少:UD | 已有:CR | 2 |
| 平台后台 | 营销 / 积分 / 商品分类 | `/marketing/integral/classify` | 缺少:D | 已有:CRU | 12 |
| 平台后台 | 营销 / 积分 / 商品列表 | `/marketing/integral/proList` | 缺少:D | 已有:CRU | 7 |
| 平台后台 | 营销 / 积分 / 积分日志 | `/marketing/integral/log` | 缺少:CUD | 已有:R | 5 |
| 平台后台 | 营销 / 积分 / 积分订单 | `/marketing/integral/orderList` | 缺少:C | 已有:RUD | 9 |
| 平台后台 | 营销 / 积分 / 积分配置 | `/marketing/integral/config` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 营销 / 预售 | `/marketing/presell` | 缺少:D | 已有:CRU | 7 |
| 平台后台 | 营销 / 预售 / 预售协议 | `/marketing/presell/agreement` | 缺少:RUD | 已有:C | 1 |
| 平台后台 | 营销 / 预售 / 预售商品 | `/marketing/presell/list` | 缺少:CD | 已有:RU | 21 |
| 平台后台 | 装修 / 主题风格 | `/setting/theme_style` | 缺少:RUD | 已有:C | 1 |
| 平台后台 | 装修 / 商品分类 | `/setting/product_category` | 缺少:UD | 已有:CR | 2 |
| 平台后台 | 装修 / 悬浮菜单 | `/setting/fab` | 缺少:UD | 已有:CR | 4 |
| 平台后台 | 装修 / 页面配置 | `/setting/system_visualization_data` | 缺少:UD | 已有:CR | 3 |
| 平台后台 | 订单 | `/order` | 缺少:CUD | 已有:R | 3 |
| 平台后台 | 订单 / 核销记录 | `/order/cancellation` | 缺少:CD | 已有:RU | 9 |
| 平台后台 | 订单 / 订单列表 | `/order/list` | 缺少:CD | 已有:RU | 30 |
| 平台后台 | 订单 / 退款订单 | `/order/refund` | 缺少:CD | 已有:RU | 21 |
| 平台后台 | 设置 / 商城设置 / 协议规则 | `/setting/agreements` | 缺少:UD | 已有:CR | 2 |
| 平台后台 | 设置 / 商城设置 / 商城设置 | `/systemForm/Basics/shop_tabs` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 设置 / 商城设置 / 支付设置 | `/systemForm/Basics/pay_tabs` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 设置 / 商城设置 / 配送配置 / 第三方送 | `/delivery` | 缺少:CUD | 已有:R | 3 |
| 平台后台 | 设置 / 商城设置 / 配送配置 / 第三方送 / 充值记录 | `/delivery/recharge_record` | 缺少:UD | 已有:CR | 3 |
| 平台后台 | 设置 / 商城设置 / 配送配置 / 第三方送 / 配送记录 | `/delivery/usage_record` | 缺少:CUD | 已有:R | 2 |
| 平台后台 | 设置 / 商城设置 / 配送配置 / 第三方送 / 配送配置 | `/systemForm/delivery` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 设置 / 应用配置 / APP升级配置 | `/systemForm/Basics/app_version` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 设置 / 应用配置 / APP配置 | `/systemForm/Basics/wechat_open_app` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 设置 / 应用配置 / 上传校验文件 | `/app/wechat/file` | 缺少:RD | 已有:CU | 4 |
| 平台后台 | 设置 / 应用配置 / 公众号配置 | `/systemForm/Basics/wechat` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 设置 / 应用配置 / 小程序配置 | `/systemForm/Basics/smallapp` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 设置 / 权限管理 | `/setting` | 缺少:D | 已有:CRU | 7 |
| 平台后台 | 设置 / 系统设置 | `/sys` | 缺少:CD | 已有:RU | 2 |
| 平台后台 | 设置 / 系统设置 / 一号通 / 商户结余 | `/service/balance_record` | 缺少:CUD | 已有:R | 1 |
| 平台后台 | 设置 / 系统设置 / 一号通 / 登陆入口 | `/setting/sms/sms_config/index` | 缺少:D | 已有:CRU | 17 |
| 平台后台 | 设置 / 系统设置 / 一号通 / 短信设置 | `/sms` | 缺少:UD | 已有:CR | 2 |
| 平台后台 | 设置 / 系统设置 / 一号通 / 短信设置 / 申请记录 | `/sms/applyList` | 缺少:UD | 已有:CR | 1 |
| 平台后台 | 设置 / 系统设置 / 一号通 / 短信设置 / 短信模板 | `/sms/template` | 缺少:RUD | 已有:C | 2 |
| 平台后台 | 设置 / 系统设置 / 一号通 / 短信设置 / 短信配置 | `/systemForm/Basics/message` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 设置 / 系统设置 / 一号通 / 购买记录 | `/service/purchase` | 缺少:CUD | 已有:R | 2 |
| 平台后台 | 设置 / 系统设置 / 接口配置 | `/systemForm/Basics/extend_tabs` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 设置 / 系统设置 / 系统设置 | `/systemForm/Basics/system_tabs` | 缺少:RD | 已有:CU | 1 |
| 平台后台 | 财务 | `/accounts` | 缺少:D | 已有:CRU | 5 |
| 平台后台 | 财务 / 发票管理 / 发票列表 | `/accounts/receipt` | 缺少:CUD | 已有:R | 2 |
| 平台后台 | 财务 / 发票管理 / 发票说明 | `/accounts/invoiceDesc` | 缺少:RUD | 已有:C | 1 |
| 平台后台 | 财务 / 商圈代理 | `/accounts/zoneAgent` | 缺少:D | 已有:CRU | 5 |
| 平台后台 | 财务 / 商圈代理 / 提成流水 | `/accounts/zoneAgent/commissionRecord` | 缺少:CUD | 已有:R | 2 |
| 平台后台 | 财务 / 商圈代理 / 申请结算 | `/accounts/zoneAgent/settlementApply` | 缺少:D | 已有:CRU | 4 |
| 平台后台 | 财务 / 商圈代理 / 结算审核 | `/accounts/zoneAgent/settlementReview` | 缺少:CRD | 已有:U | 3 |
| 平台后台 | 财务 / 商圈代理 / 结算账号 | `/accounts/zoneAgent/settlementAccount` | 缺少:CUD | 已有:R | 2 |
| 平台后台 | 财务 / 店铺结算 / 分账管理 | `/merchant/applyList` | 缺少:CUD | 已有:R | 5 |
| 平台后台 | 财务 / 店铺结算 / 平台账单 | `/accounts/statement` | 缺少:CUD | 已有:R | 6 |
| 平台后台 | 财务 / 店铺结算 / 店铺账单 | `/accounts/merchantBill` | 缺少:CUD | 已有:R | 10 |
| 平台后台 | 财务 / 店铺结算 / 转账记录 | `/accounts/transferRecord` | 缺少:CD | 已有:RU | 9 |
| 平台后台 | 财务 / 店铺结算 / 转账设置 | `/accounts/settings` | 缺少:D | 已有:CRU | 2 |
| 平台后台 | 财务 / 用户结算 / 充值记录 | `/accounts/bill` | 缺少:CUD | 已有:R | 2 |
| 平台后台 | 财务 / 用户结算 / 提现管理 | `/accounts/extract` | 缺少:D | 已有:CRU | 6 |
| 平台后台 | 财务 / 用户结算 / 资金记录 | `/accounts/capital` | 缺少:CUD | 已有:R | 4 |
| 平台后台 | 财务 / 资金流水 | `/accounts/capitalFlow` | 缺少:CUD | 已有:R | 5 |
| 平台后台 | 首页 | `/` | 缺少:UD | 已有:CR | 2 |
| 平台后台 | 首页 / 商品统计 | `/statistic/product` | 缺少:CUD | 已有:R | 9 |
| 平台后台 | 首页 / 控制台 | `/dashboard` | 缺少:CD | 已有:RU | 13 |
| 平台后台 | 首页 / 数据大屏 | `/data-screen/index` | 缺少:CRUD | 已有: | 1 |
| 平台后台 | 首页 / 用户统计 | `/statistic/member` | 缺少:CUD | 已有:R | 3 |
| 平台后台 | 首页 / 订单统计 | `/statistic/order` | 缺少:CUD | 已有:R | 9 |
