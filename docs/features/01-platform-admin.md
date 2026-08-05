# 平台后台 — 功能点清单（页面 → 按钮/操作，含增删改查）

> 主来源：`eb_system_menu` 按钮权限。
> 若页面无按钮权限节点，则用同页 `route` 的 `_path` 接口补全（标注来源=路由）。
> CRUD：`C`创建 · `R`读取 · `U`更新 · `D`删除 · `O`其他。
>
> 本页的页面覆盖与操作统计为生成时的展示快照。全量逐项验收以 [`../generated/features-master.tsv`](../generated/features-master.tsv) 为唯一口径，其中平台后台为 **1333** 项；总验收入口见 [`../CRMEB-FULL-FUNCTION-CHECKLIST.md`](../CRMEB-FULL-FUNCTION-CHECKLIST.md)。

## 统计

| 项 | 数量 |
| --- | ---: |
| 页面 | 238 |
| 操作（按钮/接口） | 1295 |
| 其中由路由补全 | 148 |
| C | 240 |
| R | 582 |
| U | 415 |
| D | 114 |

## 页面 CRUD 覆盖一览

| 模块路径 | 页面路由 | C | R | U | D | 操作数 |
| --- | --- | --- | --- | --- | --- | ---: |
| 内容 | `/content` | — | ✓ | — | — | 0 |
| 内容 / 文章 | `/cms` | ✓ | ✓ | ✓ | ✓ | 5 |
| 内容 / 文章 / 文章分类 | `/cms/articleCategory` | ✓ | ✓ | ✓ | ✓ | 6 |
| 内容 / 文章 / 文章管理 | `/cms/article` | ✓ | ✓ | ✓ | ✓ | 8 |
| 内容 / 社区 | `/community` | ✓ | ✓ | ✓ | ✓ | 9 |
| 内容 / 社区 / 社区内容 | `/community/list` | — | ✓ | ✓ | ✓ | 7 |
| 内容 / 社区 / 社区分类 | `/community/category` | ✓ | ✓ | ✓ | ✓ | 6 |
| 内容 / 社区 / 社区评论 | `/community/reply` | — | ✓ | ✓ | ✓ | 4 |
| 内容 / 社区 / 社区话题 | `/community/topic` | ✓ | ✓ | ✓ | ✓ | 9 |
| 内容 / 社区 / 社区配置 | `/systemForm/Basics/community` | ✓ | — | ✓ | — | 1 |
| 分销 | `/promoter` | — | ✓ | — | — | 8 |
| 分销 / 佣金说明 | `/promoter/commission` | ✓ | — | — | — | 1 |
| 分销 / 分销员列表 | `/promoter/user` | ✓ | ✓ | ✓ | ✓ | 6 |
| 分销 / 分销海报 | `/group/config/68` | ✓ | ✓ | ✓ | ✓ | 8 |
| 分销 / 分销特权 | `/group/config/75` | ✓ | ✓ | ✓ | ✓ | 8 |
| 分销 / 分销礼包 | `/promoter/gift` | — | ✓ | ✓ | — | 6 |
| 分销 / 分销等级 | `brokerage` | — | ✓ | — | — | 0 |
| 分销 / 分销等级 / 分销员等级 | `/promoter/membership_level` | ✓ | ✓ | ✓ | ✓ | 6 |
| 分销 / 分销等级 / 等级规则 | `/promoter/distribution` | ✓ | — | — | — | 1 |
| 分销 / 分销订单 | `/promoter/orderList` | — | ✓ | ✓ | — | 8 |
| 分销 / 分销说明 | `/promoter/retail` | — | ✓ | — | — | 0 |
| 分销 / 分销配置 | `/systemForm/Basics/distribution_tabs` | ✓ | — | ✓ | — | 1 |
| 分销 / 提现银行 | `/group/config/76` | ✓ | ✓ | ✓ | ✓ | 8 |
| 商品 | `/product` | ✓ | ✓ | — | ✓ | 5 |
| 商品 / 价格说明 | `/product/priceDescription` | ✓ | ✓ | ✓ | ✓ | 7 |
| 商品 / 保障服务 | `/product/guarantee` | ✓ | ✓ | ✓ | ✓ | 9 |
| 商品 / 品牌管理 | `/product/brand` | — | ✓ | — | — | 0 |
| 商品 / 品牌管理 / 品牌分类 | `/product/band/brandClassify` | ✓ | ✓ | ✓ | ✓ | 6 |
| 商品 / 品牌管理 / 品牌列表 | `/product/band/brandList` | ✓ | ✓ | ✓ | ✓ | 5 |
| 商品 / 商品分类 | `/product/classify` | ✓ | ✓ | ✓ | ✓ | 9 |
| 商品 / 商品参数 | `/product/specsMain` | ✓ | ✓ | ✓ | ✓ | 11 |
| 商品 / 商品参数 / 平台商品参数 | `/product/specs` | ✓ | ✓ | ✓ | ✓ | 6 |
| 商品 / 商品参数 / 店铺商品参数 | `/product/merSpecs` | — | ✓ | — | — | 1 |
| 商品 / 商品标签 | `/product/label` | ✓ | ✓ | ✓ | ✓ | 6 |
| 商品 / 商品管理 | `/product/examine` | ✓ | ✓ | ✓ | — | 48 |
| 商品 / 活动标签 | `/product/activityLabel` | — | ✓ | — | — | 0 |
| 商品 / 评论管理 | `/product/comment` | ✓ | ✓ | ✓ | ✓ | 6 |
| 客服 | `/service` | ✓ | ✓ | ✓ | — | 4 |
| 客服 / 客服列表 | `/service/customer/list` | ✓ | ✓ | ✓ | ✓ | 10 |
| 客服 / 客服自动回复 | `/systemForm/customer_keyword` | ✓ | ✓ | ✓ | ✓ | 5 |
| 客服 / 客服设置 | `/systemForm/Basics/service` | ✓ | — | ✓ | — | 1 |
| 应用 | `/apploction` | — | ✓ | — | — | 0 |
| 应用 / 公众号 | `/app/wechat` | ✓ | ✓ | ✓ | ✓ | 9 |
| 应用 / 公众号 / 图文管理 | `/app/wechat/newsCategory` | ✓ | ✓ | ✓ | ✓ | 7 |
| 应用 / 公众号 / 微信模板消息 | `/app/wechat/template` | ✓ | ✓ | ✓ | ✓ | 6 |
| 应用 / 公众号 / 微信菜单 | `/app/wechat/menus` | ✓ | — | ✓ | — | 2 |
| 应用 / 公众号 / 自动回复 | `/admin/app/wechat/reply` | ✓ | ✓ | ✓ | ✓ | 9 |
| 应用 / 公众号 / 自动回复 / 关键字回复 | `/app/wechat/reply/keyword` | — | ✓ | — | — | 0 |
| 应用 / 公众号 / 自动回复 / 微信关注回复 | `/app/wechat/reply/follow/subscribe` | — | ✓ | — | — | 0 |
| 应用 / 公众号 / 自动回复 / 无效关键词回复 | `/app/wechat/reply/index/default` | — | ✓ | — | — | 0 |
| 应用 / 小程序 | `/app/routine` | — | ✓ | — | — | 2 |
| 应用 / 小程序 / 小程序下载 | `/app/routine/download` | — | ✓ | — | — | 1 |
| 应用 / 小程序 / 小程序订阅消息 | `/app/routine/template` | ✓ | ✓ | ✓ | ✓ | 6 |
| 店铺 | `/mer` | ✓ | ✓ | — | — | 2 |
| 店铺 / 区域代理 | `/business-zones/manage` | — | ✓ | — | — | 0 |
| 店铺 / 区域代理 / 代理人员 | `/business-zones/agents` | ✓ | ✓ | ✓ | ✓ | 6 |
| 店铺 / 区域代理 / 代理审核 | `/business-zones/agent-review` | — | — | ✓ | — | 1 |
| 店铺 / 区域代理 / 代理设置 | `/business-zones/settings` | — | ✓ | — | — | 0 |
| 店铺 / 区域代理 / 区域列表 | `/business-zones/index` | ✓ | ✓ | ✓ | ✓ | 7 |
| 店铺 / 商户管理 | `/merchant` | ✓ | ✓ | — | — | 2 |
| 店铺 / 商户管理 / 商户入驻审核 | `/merchant/review` | — | — | ✓ | — | 1 |
| 店铺 / 商户管理 / 商户列表 | `/merchant/index` | ✓ | ✓ | ✓ | ✓ | 7 |
| 店铺 / 商户管理 / 商户管理员 | `/merchant/admin-list` | ✓ | ✓ | ✓ | ✓ | 6 |
| 店铺 / 商户管理 / 商户设置 | `/merchant/apply-setting` | — | ✓ | — | — | 0 |
| 店铺 / 店铺管理 | `/mer/mer` | — | ✓ | — | — | 0 |
| 店铺 / 店铺管理 / 商户列表 | `/merchant/list` | ✓ | ✓ | ✓ | ✓ | 17 |
| 店铺 / 店铺管理 / 店铺入驻申请 | `/merchant/application` | — | ✓ | ✓ | ✓ | 8 |
| 店铺 / 店铺管理 / 店铺分类 | `/merchant/classify` | ✓ | ✓ | ✓ | ✓ | 4 |
| 店铺 / 店铺管理 / 店铺分组 | `/merchant/grouping` | ✓ | ✓ | ✓ | ✓ | 12 |
| 店铺 / 店铺管理 / 店铺分账申请 | `/merchant/applyments` | ✓ | ✓ | ✓ | — | 4 |
| 店铺 / 店铺管理 / 店铺列表 | `/merchant/list` | ✓ | ✓ | ✓ | ✓ | 34 |
| 店铺 / 店铺管理 / 店铺类型 | `/merchant/type` | ✓ | ✓ | ✓ | ✓ | 6 |
| 店铺 / 店铺设置 | `/mer/store` | — | ✓ | — | — | 0 |
| 店铺 / 店铺设置 / 保证金配置 | `/systemForm/Basics/margin` | ✓ | — | ✓ | — | 1 |
| 店铺 / 店铺设置 / 店铺保证金 | `/merchant/deposit_list` | ✓ | ✓ | ✓ | ✓ | 8 |
| 店铺 / 店铺设置 / 店铺菜单 | `/merchant/system` | ✓ | ✓ | ✓ | ✓ | 4 |
| 店铺 / 店铺设置 / 说明提示 | `/merchant/type/description` | ✓ | — | — | — | 1 |
| 用户 | `/user` | ✓ | ✓ | — | — | 7 |
| 用户 / 付费会员 | `/user/svip` | — | ✓ | — | — | 0 |
| 用户 / 付费会员 / 付费会员配置 | `/systemForm/Basics/svip` | ✓ | — | ✓ | — | 1 |
| 用户 / 付费会员 / 会员协议 | `/user/member/vipAgreement` | — | ✓ | — | — | 0 |
| 用户 / 付费会员 / 会员权益 | `/user/member/equity` | ✓ | ✓ | ✓ | — | 5 |
| 用户 / 付费会员 / 会员类型 | `/user/member/type` | ✓ | ✓ | ✓ | ✓ | 8 |
| 用户 / 付费会员 / 会员记录 | `/user/member/record` | — | ✓ | ✓ | — | 4 |
| 用户 / 搜索记录 | `/user/searchRecord` | — | ✓ | — | ✓ | 5 |
| 用户 / 用户分组 | `/user/group` | ✓ | ✓ | ✓ | ✓ | 4 |
| 用户 / 用户列表 | `/user/list` | ✓ | ✓ | ✓ | — | 28 |
| 用户 / 用户协议 | `/user/agreement` | ✓ | — | — | — | 1 |
| 用户 / 用户反馈 | `/feedback` | — | ✓ | ✓ | ✓ | 5 |
| 用户 / 用户反馈 / 反馈分类 | `/feedback/classify` | ✓ | ✓ | ✓ | ✓ | 5 |
| 用户 / 用户反馈 / 反馈列表 | `/feedback/list` | — | ✓ | ✓ | ✓ | 4 |
| 用户 / 用户标签 | `/user/label` | ✓ | ✓ | ✓ | ✓ | 4 |
| 用户 / 用户等级 | `/user/member` | — | — | ✓ | ✓ | 2 |
| 用户 / 用户等级 / 等级权益 | `/user/member/interests` | ✓ | ✓ | ✓ | ✓ | 6 |
| 用户 / 用户等级 / 等级管理 | `/user/member/list` | ✓ | ✓ | ✓ | ✓ | 6 |
| 用户 / 用户等级 / 等级说明 | `/user/member/description` | ✓ | — | — | — | 1 |
| 用户 / 用户等级 / 等级配置 | `/systemForm/Basics/members` | ✓ | — | ✓ | — | 1 |
| 用户 / 用户设置 | `/user/setup_user` | ✓ | ✓ | ✓ | ✓ | 8 |
| 维护 | `/safe` | — | ✓ | — | — | 0 |
| 维护 / 安全维护 | `/maintain` | — | ✓ | — | ✓ | 2 |
| 维护 / 安全维护 / 商业授权 | `/setting/system/maintain/auth` | ✓ | ✓ | — | — | 3 |
| 维护 / 安全维护 / 数据备份 | `/maintain/dataBackup` | — | ✓ | — | ✓ | 8 |
| 维护 / 安全维护 / 缓存清除 | `/maintain/cache` | — | — | — | ✓ | 2 |
| 维护 / 导出记录 | `/group/exportList` | — | ✓ | — | — | 2 |
| 维护 / 开发配置 | `/safe/exploit` | — | ✓ | — | — | 0 |
| 维护 / 开发配置 / 组合数据 | `/group/list` | ✓ | ✓ | ✓ | ✓ | 11 |
| 维护 / 操作日志 | `/setting/systemLog` | — | ✓ | — | — | 1 |
| 维护 / 配置分类 | `/config/classify` | ✓ | ✓ | ✓ | ✓ | 5 |
| 维护 / 配置管理 | `/config/setting` | ✓ | ✓ | ✓ | ✓ | 5 |
| 维护 / 页面链接 | `/safe/pageLinks` | — | ✓ | — | — | 0 |
| 营销 | `/marketing` | — | ✓ | — | — | 2 |
| 营销 / 专场列表 | `/group/topic/94` | ✓ | ✓ | ✓ | ✓ | 8 |
| 营销 / 优惠套餐 | `/marketing/discounts/list` | ✓ | ✓ | ✓ | ✓ | 10 |
| 营销 / 余额充值 | `/banlace` | — | ✓ | — | — | 0 |
| 营销 / 余额充值 / 余额充值配置 | `/group/config/69` | ✓ | ✓ | ✓ | ✓ | 8 |
| 营销 / 余额充值 / 余额设置 | `/systemForm/Basics/balance` | ✓ | — | ✓ | — | 1 |
| 营销 / 助力 | `/assist` | — | ✓ | — | — | 0 |
| 营销 / 助力 / 助力活动 | `/marketing/assist/list` | — | ✓ | — | — | 6 |
| 营销 / 助力 / 活动商品 | `/marketing/assist/goods_list` | — | ✓ | ✓ | — | 21 |
| 营销 / 商户优惠券 | `/marketing/coupon` | — | ✓ | — | — | 2 |
| 营销 / 商户优惠券 / 优惠券列表 | `/marketing/coupon/list` | — | ✓ | — | — | 3 |
| 营销 / 商户优惠券 / 领取记录 | `/marketing/coupon/user` | — | ✓ | — | — | 1 |
| 营销 / 平台优惠券 | `/marketing/platform_coupon` | ✓ | ✓ | — | — | 3 |
| 营销 / 平台优惠券 / 优惠券列表 | `/marketing/platform_coupon/list` | ✓ | ✓ | ✓ | ✓ | 9 |
| 营销 / 平台优惠券 / 使用说明 | `/marketing/platform_coupon/instructions` | ✓ | — | — | — | 1 |
| 营销 / 平台优惠券 / 发送记录 | `/marketing/platform_coupon/couponSend` | ✓ | ✓ | — | — | 1 |
| 营销 / 平台优惠券 / 领取记录 | `/marketing/platform_coupon/couponRecord` | — | ✓ | — | — | 1 |
| 营销 / 报名活动 | `/marketing/application/list` | ✓ | ✓ | ✓ | ✓ | 10 |
| 营销 / 拼团 | `/marketing/combination` | ✓ | — | ✓ | — | 2 |
| 营销 / 拼团 / 拼团商品列表 | `/marketing/combination/combination_goods` | — | ✓ | ✓ | — | 24 |
| 营销 / 拼团 / 拼团活动列表 | `/marketing/combination/combination_list` | — | ✓ | — | — | 6 |
| 营销 / 拼团 / 拼团设置 | `/marketing/combination/combination_set` | ✓ | — | ✓ | — | 1 |
| 营销 / 活动氛围图 | `/marketing/atmosphere/list` | ✓ | ✓ | ✓ | ✓ | 8 |
| 营销 / 活动边框图 | `/marketing/border/list` | ✓ | ✓ | ✓ | ✓ | 8 |
| 营销 / 直播 | `/marketing2` | — | ✓ | — | — | 0 |
| 营销 / 直播 / 直播商品管理 | `/marketing/broadcast/list` | — | ✓ | ✓ | ✓ | 6 |
| 营销 / 直播 / 直播间管理 | `/marketing/studio/list` | ✓ | ✓ | ✓ | ✓ | 11 |
| 营销 / 秒杀 | `/marketing/seckill` | ✓ | ✓ | ✓ | ✓ | 24 |
| 营销 / 秒杀 / 秒杀活动 | `/marketing/seckill/store_seckill/list` | ✓ | ✓ | ✓ | ✓ | 30 |
| 营销 / 秒杀 / 秒杀管理 | `/marketing/seckill/list` | — | ✓ | ✓ | ✓ | 29 |
| 营销 / 秒杀 / 秒杀配置 | `/marketing/seckill/seckillConfig` | ✓ | ✓ | ✓ | ✓ | 7 |
| 营销 / 积分 | `/marketing/integral` | ✓ | ✓ | — | — | 2 |
| 营销 / 积分 / 商品分类 | `/marketing/integral/classify` | ✓ | ✓ | ✓ | — | 12 |
| 营销 / 积分 / 商品列表 | `/marketing/integral/proList` | ✓ | ✓ | ✓ | — | 7 |
| 营销 / 积分 / 积分日志 | `/marketing/integral/log` | — | ✓ | — | — | 5 |
| 营销 / 积分 / 积分订单 | `/marketing/integral/orderList` | — | ✓ | ✓ | ✓ | 9 |
| 营销 / 积分 / 积分配置 | `/marketing/integral/config` | ✓ | — | ✓ | — | 1 |
| 营销 / 预售 | `/marketing/presell` | ✓ | ✓ | ✓ | — | 7 |
| 营销 / 预售 / 预售协议 | `/marketing/presell/agreement` | ✓ | — | — | — | 1 |
| 营销 / 预售 / 预售商品 | `/marketing/presell/list` | — | ✓ | ✓ | — | 21 |
| 装修 | `/theme` | — | ✓ | — | — | 0 |
| 装修 / 个人中心 | `/setting/diy/personal` | — | ✓ | — | — | 0 |
| 装修 / 主题风格 | `/setting/theme_style` | ✓ | — | — | — | 1 |
| 装修 / 商品分类 | `/setting/product_category` | ✓ | ✓ | — | — | 2 |
| 装修 / 商品详情 | `/setting/diy/product_detail` | — | ✓ | — | — | 0 |
| 装修 / 店铺模板 | `/setting/merchant/diyList` | ✓ | ✓ | ✓ | ✓ | 9 |
| 装修 / 店铺街 | `/setting/diy/store` | — | ✓ | — | — | 0 |
| 装修 / 微页面 | `/setting/micro/list` | ✓ | ✓ | ✓ | ✓ | 7 |
| 装修 / 悬浮菜单 | `/setting/fab` | ✓ | ✓ | — | — | 4 |
| 装修 / 系统表单 | `/systemForm/form_list` | ✓ | ✓ | ✓ | ✓ | 8 |
| 装修 / 素材管理 | `/config/picture` | ✓ | ✓ | ✓ | ✓ | 12 |
| 装修 / 页面装修 | `/setting/diy/list` | ✓ | ✓ | ✓ | ✓ | 15 |
| 装修 / 页面配置 | `/setting/system_visualization_data` | ✓ | ✓ | — | — | 3 |
| 装修 / 页面链接 | `/setting/page` | — | ✓ | — | — | 0 |
| 装修 / 页面链接 / 商户页面分类 | `/setting/diy/merchant/category/list` | ✓ | ✓ | ✓ | ✓ | 5 |
| 装修 / 页面链接 / 商户页面链接 | `/setting/diy/merLink/list` | ✓ | ✓ | ✓ | ✓ | 5 |
| 装修 / 页面链接 / 平台页面分类 | `/setting/diy/plantform/category/list` | ✓ | ✓ | ✓ | ✓ | 5 |
| 装修 / 页面链接 / 平台页面链接 | `/setting/diy/links/list` | ✓ | ✓ | ✓ | ✓ | 5 |
| 订单 | `/order` | — | ✓ | — | — | 3 |
| 订单 / 核销记录 | `/order/cancellation` | — | ✓ | ✓ | — | 9 |
| 订单 / 订单列表 | `/order/list` | — | ✓ | ✓ | — | 30 |
| 订单 / 退款订单 | `/order/refund` | — | ✓ | ✓ | — | 21 |
| 设置 | `/settings` | — | ✓ | — | — | 0 |
| 设置 / 商城设置 | `/shop` | — | ✓ | — | — | 0 |
| 设置 / 商城设置 / 协议规则 | `/setting/agreements` | ✓ | ✓ | — | — | 2 |
| 设置 / 商城设置 / 商城设置 | `/systemForm/Basics/shop_tabs` | ✓ | — | ✓ | — | 1 |
| 设置 / 商城设置 / 支付设置 | `/systemForm/Basics/pay_tabs` | ✓ | — | ✓ | — | 1 |
| 设置 / 商城设置 / 热门搜索 | `/group/config/67` | ✓ | ✓ | ✓ | ✓ | 8 |
| 设置 / 商城设置 / 配送配置 | `/delivery_config` | — | ✓ | — | — | 0 |
| 设置 / 商城设置 / 配送配置 / 城市数据 | `/freight/city/list` | — | ✓ | — | — | 0 |
| 设置 / 商城设置 / 配送配置 / 物流公司 | `/freight/express` | ✓ | ✓ | ✓ | ✓ | 10 |
| 设置 / 商城设置 / 配送配置 / 第三方送 | `/delivery` | — | ✓ | — | — | 3 |
| 设置 / 商城设置 / 配送配置 / 第三方送 / 充值记录 | `/delivery/recharge_record` | ✓ | ✓ | — | — | 3 |
| 设置 / 商城设置 / 配送配置 / 第三方送 / 配送记录 | `/delivery/usage_record` | — | ✓ | — | — | 2 |
| 设置 / 商城设置 / 配送配置 / 第三方送 / 配送配置 | `/systemForm/delivery` | ✓ | — | ✓ | — | 1 |
| 设置 / 应用配置 | `/app_config` | — | ✓ | — | — | 0 |
| 设置 / 应用配置 / APP升级配置 | `/systemForm/Basics/app_version` | ✓ | — | ✓ | — | 1 |
| 设置 / 应用配置 / APP配置 | `/systemForm/Basics/wechat_open_app` | ✓ | — | ✓ | — | 1 |
| 设置 / 应用配置 / 上传校验文件 | `/app/wechat/file` | ✓ | — | ✓ | — | 4 |
| 设置 / 应用配置 / 公众号配置 | `/systemForm/Basics/wechat` | ✓ | — | ✓ | — | 1 |
| 设置 / 应用配置 / 小程序配置 | `/systemForm/Basics/smallapp` | ✓ | — | ✓ | — | 1 |
| 设置 / 权限管理 | `/setting` | ✓ | ✓ | ✓ | — | 7 |
| 设置 / 权限管理 / 管理员管理 | `/setting/systemAdmin` | ✓ | ✓ | ✓ | ✓ | 24 |
| 设置 / 权限管理 / 菜单管理 | `/setting/menu` | ✓ | ✓ | ✓ | ✓ | 5 |
| 设置 / 权限管理 / 角色权限 | `/setting/systemRole` | ✓ | ✓ | ✓ | ✓ | 27 |
| 设置 / 消息管理 | `/notice` | — | ✓ | — | — | 0 |
| 设置 / 消息管理 / 公告管理 | `/station/notice` | ✓ | ✓ | ✓ | ✓ | 6 |
| 设置 / 消息管理 / 消息管理 | `/setting/notification/index` | ✓ | ✓ | ✓ | ✓ | 8 |
| 设置 / 系统设置 | `/sys` | — | ✓ | ✓ | — | 2 |
| 设置 / 系统设置 / 一号通 | `/serve` | — | ✓ | — | — | 0 |
| 设置 / 系统设置 / 一号通 / 商户结余 | `/service/balance_record` | — | ✓ | — | — | 1 |
| 设置 / 系统设置 / 一号通 / 服务配置 | `/service/settings` | ✓ | ✓ | ✓ | ✓ | 6 |
| 设置 / 系统设置 / 一号通 / 登陆入口 | `/setting/sms/sms_config/index` | ✓ | ✓ | ✓ | — | 17 |
| 设置 / 系统设置 / 一号通 / 短信设置 | `/sms` | ✓ | ✓ | — | — | 2 |
| 设置 / 系统设置 / 一号通 / 短信设置 / 申请记录 | `/sms/applyList` | ✓ | ✓ | — | — | 1 |
| 设置 / 系统设置 / 一号通 / 短信设置 / 短信模板 | `/sms/template` | ✓ | — | — | — | 2 |
| 设置 / 系统设置 / 一号通 / 短信设置 / 短信账户 | `/sms/user` | — | ✓ | — | — | 0 |
| 设置 / 系统设置 / 一号通 / 短信设置 / 短信配置 | `/systemForm/Basics/message` | ✓ | — | ✓ | — | 1 |
| 设置 / 系统设置 / 一号通 / 购买记录 | `/service/purchase` | — | ✓ | — | — | 2 |
| 设置 / 系统设置 / 存储管理 | `/setting/storage` | ✓ | ✓ | ✓ | ✓ | 9 |
| 设置 / 系统设置 / 接口配置 | `/systemForm/Basics/extend_tabs` | ✓ | — | ✓ | — | 1 |
| 设置 / 系统设置 / 系统设置 | `/systemForm/Basics/system_tabs` | ✓ | — | ✓ | — | 1 |
| 财务 | `/accounts` | ✓ | ✓ | ✓ | — | 5 |
| 财务 / 发票管理 | `/accounts/accounts` | — | ✓ | — | — | 0 |
| 财务 / 发票管理 / 发票列表 | `/accounts/receipt` | — | ✓ | — | — | 2 |
| 财务 / 发票管理 / 发票说明 | `/accounts/invoiceDesc` | ✓ | — | — | — | 1 |
| 财务 / 商圈代理 | `/accounts/zoneAgent` | ✓ | ✓ | ✓ | — | 5 |
| 财务 / 商圈代理 / 提成流水 | `/accounts/zoneAgent/commissionRecord` | — | ✓ | — | — | 2 |
| 财务 / 商圈代理 / 申请结算 | `/accounts/zoneAgent/settlementApply` | ✓ | ✓ | ✓ | — | 4 |
| 财务 / 商圈代理 / 结算审核 | `/accounts/zoneAgent/settlementReview` | — | — | ✓ | — | 3 |
| 财务 / 商圈代理 / 结算账号 | `/accounts/zoneAgent/settlementAccount` | — | ✓ | — | — | 2 |
| 财务 / 店铺结算 | `/mer/accounts` | — | ✓ | — | — | 0 |
| 财务 / 店铺结算 / 分账管理 | `/merchant/applyList` | — | ✓ | — | — | 5 |
| 财务 / 店铺结算 / 平台账单 | `/accounts/statement` | — | ✓ | — | — | 6 |
| 财务 / 店铺结算 / 店铺账单 | `/accounts/merchantBill` | — | ✓ | — | — | 10 |
| 财务 / 店铺结算 / 转账记录 | `/accounts/transferRecord` | — | ✓ | ✓ | — | 9 |
| 财务 / 店铺结算 / 转账设置 | `/accounts/settings` | ✓ | ✓ | ✓ | — | 2 |
| 财务 / 用户结算 | `/accounts/record` | — | ✓ | — | — | 0 |
| 财务 / 用户结算 / 充值记录 | `/accounts/bill` | — | ✓ | — | — | 2 |
| 财务 / 用户结算 / 提现管理 | `/accounts/extract` | ✓ | ✓ | ✓ | — | 6 |
| 财务 / 用户结算 / 资金记录 | `/accounts/capital` | — | ✓ | — | — | 4 |
| 财务 / 资金流水 | `/accounts/capitalFlow` | — | ✓ | — | — | 5 |
| 首页 | `/` | ✓ | ✓ | — | — | 2 |
| 首页 / 商品统计 | `/statistic/product` | — | ✓ | — | — | 9 |
| 首页 / 控制台 | `/dashboard` | — | ✓ | ✓ | — | 13 |
| 首页 / 数据大屏 | `/data-screen/index` | — | — | — | — | 1 |
| 首页 / 用户统计 | `/statistic/member` | — | ✓ | — | — | 3 |
| 首页 / 订单统计 | `/statistic/order` | — | ✓ | — | — | 9 |

## 分页面操作明细

### 内容

#### 内容

- 页面路由：`/content`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 内容 / 文章

- 页面路由：`/cms`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 文章分类详情 | `R` | `systemArticleCategoryDetail` | 路由 |
| 文章列表 | `R` | `systemArticlArticleLst` | 路由 |
| 文章添加 | `C` | `systemArticleArticleCreate` | 路由 |
| 文章编辑 | `U` | `systemArticArticleleUpdate` | 路由 |
| 文章删除 | `D` | `systemArticArticleleDelete` | 路由 |

#### 内容 / 文章 / 文章分类

- 页面路由：`/cms/articleCategory`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 文章分类列表 | `R` | `systemArticleCategoryLst` | 菜单权限 |
| 文章分类添加 | `C` | `systemArticleCategoryCreate` | 菜单权限 |
| 文章分类编辑 | `U` | `systemArticleCategoryUpdate` | 菜单权限 |
| 文章分类修改状态 | `U` | `systemArticleCategoryStatus` | 菜单权限 |
| 文章分类删除 | `D` | `systemArticleCategoryDelete` | 菜单权限 |
| 文章分类详情 | `R` | `systemArticleCategoryDetail` | 菜单权限 |

#### 内容 / 文章 / 文章管理

- 页面路由：`/cms/article`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 文章列表 | `R` | `systemArticlArticleLst` | 菜单权限 |
| 文章添加 | `C` | `systemArticleArticleCreate` | 菜单权限 |
| 文章编辑 | `U` | `systemArticArticleleUpdate` | 菜单权限 |
| 文章删除 | `D` | `systemArticArticleleDelete` | 菜单权限 |
| 文章详情 | `R` | `systemArticArticleleDetail` | 菜单权限 |
| 文章修改状态 | `U` | `systemArticlArticlStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 素材列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 内容 / 社区

- 页面路由：`/community`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 文章详情 | `R` | `systemCommunityShow` | 路由 |
| 统计 | `R` | `systemCommunityTitle` | 路由 |
| 社区分类状态 | `U` | `systemCommunityCategoryLst` | 路由 |
| 社区分类添加表单 | `C` | `systemCommunityCategoryCreateForm` | 路由 |
| 社区分类添加 | `C` | `systemCommunityCategoryCreate` | 路由 |
| 社区分类编辑表单 | `U` | `systemCommunityCategoryUpdateForm` | 路由 |
| 社区分类编辑 | `U` | `systemCommunityCategoryUpdate` | 路由 |
| 社区分类详情 | `R` | `systemCommunityCategoryDetail` | 路由 |
| 社区分类删除 | `D` | `systemCommunityCategoryDelete` | 路由 |

#### 内容 / 社区 / 社区内容

- 页面路由：`/community/list`
- CRUD：C=— R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 文章列表 | `R` | `systemCommunityLst` | 菜单权限 |
| 文章详情 | `R` | `systemCommunityDetail` | 菜单权限 |
| 文章编辑 | `U` | `systemCommunityUpdate` | 菜单权限 |
| 文章删除 | `D` | `systemCommunityDelete` | 菜单权限 |
| 修改状态 | `U` | `systemCommunityStatus` | 菜单权限 |
| 文章详情 | `RU` | `systemCommunityShow` | 菜单权限 |
| 统计 | `R` | `systemCommunityTitle` | 菜单权限 |

#### 内容 / 社区 / 社区分类

- 页面路由：`/community/category`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 社区分类状态 | `RU` | `systemCommunityCategoryLst` | 菜单权限 |
| 社区分类添加 | `C` | `systemCommunityCategoryCreate` | 菜单权限 |
| 社区分类编辑 | `U` | `systemCommunityCategoryUpdate` | 菜单权限 |
| 社区分类详情 | `R` | `systemCommunityCategoryDetail` | 菜单权限 |
| 社区分类删除 | `D` | `systemCommunityCategoryDelete` | 菜单权限 |
| 社区分类修改状态 | `U` | `systemCommunityCategoryStatus` | 菜单权限 |

#### 内容 / 社区 / 社区评论

- 页面路由：`/community/reply`
- CRUD：C=— R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 社区评论列表 | `R` | `systemCommunityReplyLst` | 菜单权限 |
| 社区评论删除 | `D` | `systemCommunityReplyDelete` | 菜单权限 |
| 社区评论审核 | `U` | `systemCommunityReplyStatus` | 菜单权限 |
| 内容评论列表 | `R` | `systemCommunityReply` | 菜单权限 |

#### 内容 / 社区 / 社区话题

- 页面路由：`/community/topic`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 社区话题 | `R` | `systemCommunityTopicLst` | 菜单权限 |
| 社区话题添加 | `C` | `systemCommunityTopicCreate` | 菜单权限 |
| 社区话题编辑 | `U` | `systemCommunityTopicUpdate` | 菜单权限 |
| 社区话题详情  | `R` | `systemCommunityTopicDetail` | 菜单权限 |
| 社区话题删除 | `D` | `systemCommunityTopicDelete` | 菜单权限 |
| 社区话题修改状态 | `U` | `systemCommunityTopicStatus` | 菜单权限 |
| 社区话题推荐 | `U` | `systemCommunityTopicHot` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 内容 / 社区 / 社区配置

- 页面路由：`/systemForm/Basics/community`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

### 分销

#### 分销

- 页面路由：`/promoter`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 记录 | `R` | `systemSpreadOrderStatus` | 路由 |
| 关联订单 | `R` | `systemSpreadOrderChildrenList` | 路由 |
| 列表 | `R` | `systemOrderLst` | 路由 |
| 金额统计 | `R` | `systemOrderStat` | 路由 |
| 快递查询 | `R` | `systemOrderExpress` | 路由 |
| 头部统计 | `R` | `systemOrderTitle` | 路由 |
| 详情 | `R` | `systemOrderDetail` | 路由 |
| 导出 | `R` | `systemOrderExcel` | 路由 |

#### 分销 / 佣金说明

- 页面路由：`/promoter/commission`
- CRUD：C=✓ R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 佣金说明 | `C` | `systemAgreeSave` | 菜单权限 |

#### 分销 / 分销员列表

- 页面路由：`/promoter/user`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 分销员列表 | `R` | `systemPromoterUserLst` | 菜单权限 |
| 分销员统计 | `R` | `systemPromoterUserCount` | 菜单权限 |
| 修改分销员等级 | `CRU` | `systemUserSpreadSave` | 菜单权限 |
| 推广人列表 | `R` | `systemUserSpreadLst` | 菜单权限 |
| 推广人订单 | `R` | `systemUserSpreadOrder` | 菜单权限 |
| 清除推广人 | `RD` | `systemUserSpreadClear` | 菜单权限 |

#### 分销 / 分销海报

- 页面路由：`/group/config/68`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 详情 | `R` | `groupDetail` | 菜单权限 |
| 列表 | `R` | `groupDataLst` | 菜单权限 |
| 添加 | `C` | `groupDataCreate` | 菜单权限 |
| 编辑 | `U` | `groupDataUpdate` | 菜单权限 |
| 删除 | `D` | `groupDataDelete` | 菜单权限 |
| 修改状态 | `U` | `groupDataChangeStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 分销 / 分销特权

- 页面路由：`/group/config/75`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 详情 | `R` | `groupDetail` | 菜单权限 |
| 列表 | `R` | `groupDataLst` | 菜单权限 |
| 添加 | `C` | `groupDataCreate` | 菜单权限 |
| 编辑 | `U` | `groupDataUpdate` | 菜单权限 |
| 删除 | `D` | `groupDataDelete` | 菜单权限 |
| 修改状态 | `U` | `groupDataChangeStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 分销 / 分销礼包

- 页面路由：`/promoter/gift`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 统计 | `R` | `systemStoreBagLstFilter` | 菜单权限 |
| 列表 | `R` | `systemStoreBagLst` | 菜单权限 |
| 详情 | `R` | `systemStoreBagDetail` | 菜单权限 |
| 编辑 | `U` | `systemStoreBagUpdate` | 菜单权限 |
| 修改状态 | `U` | `systemStoreBagSwitchStatus` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreBagChangeUsed` | 菜单权限 |

#### 分销 / 分销等级

- 页面路由：`brokerage`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 分销 / 分销等级 / 分销员等级

- 页面路由：`/promoter/membership_level`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 分销员等级列表 | `R` | `systemUserBrokerageLst` | 菜单权限 |
| 分销员等级添加 | `C` | `systemUserBrokerageCreate` | 菜单权限 |
| 分销员等级编辑 | `U` | `systemUserBrokerageUpdate` | 菜单权限 |
| 分销员等级删除 | `D` | `systemUserBrokerageDelete` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 分销 / 分销等级 / 等级规则

- 页面路由：`/promoter/distribution`
- CRUD：C=✓ R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 等级规则 | `C` | `systemAgreeSave` | 菜单权限 |

#### 分销 / 分销订单

- 页面路由：`/promoter/orderList`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 金额统计 | `R` | `systemSpreadOrderStat` | 菜单权限 |
| 快递查询 | `R` | `systemSpreadOrderExpress` | 菜单权限 |
| 头部统计 | `R` | `systemSpreadOrderTitle` | 菜单权限 |
| 详情 | `R` | `systemSpreadOrderDetail` | 菜单权限 |
| 导出 | `R` | `systemSpreadOrderExcel` | 菜单权限 |
| 记录 | `RU` | `systemSpreadOrderStatus` | 菜单权限 |
| 关联订单 | `R` | `systemSpreadOrderChildrenList` | 菜单权限 |
| 列表 | `R` | `systemSpreadOrderLst` | 菜单权限 |

#### 分销 / 分销说明

- 页面路由：`/promoter/retail`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 分销 / 分销配置

- 页面路由：`/systemForm/Basics/distribution_tabs`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 配置保存 | `CU` | `configOthersSettingUpdate` | 菜单权限 |

#### 分销 / 提现银行

- 页面路由：`/group/config/76`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 详情 | `R` | `groupDetail` | 菜单权限 |
| 列表 | `R` | `groupDataLst` | 菜单权限 |
| 添加 | `C` | `groupDataCreate` | 菜单权限 |
| 编辑 | `U` | `groupDataUpdate` | 菜单权限 |
| 删除 | `D` | `groupDataDelete` | 菜单权限 |
| 修改状态 | `U` | `groupDataChangeStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

### 商品

#### 商品

- 页面路由：`/product`
- CRUD：C=✓ R=✓ U=— D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 平台参数列表 | `R` | `systemStoreParameterTemplateLst` | 路由 |
| 商户参数模板 | `R` | `systemStoreParameterTemplateMerLst` | 路由 |
| 详情 | `R` | `systemStoreParameterTemplateDetail` | 路由 |
| 删除 | `D` | `systemStoreParameterTemplateDelete` | 路由 |
| 添加 | `C` | `systemStoreParameterTemplateCreate` | 路由 |

#### 商品 / 价格说明

- 页面路由：`/product/priceDescription`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 价格说明列表 | `R` | `systemPriceRuleLst` | 菜单权限 |
| 添加价格说明 | `C` | `systemPriceRuleCreate` | 菜单权限 |
| 修改价格说明 | `U` | `systemPriceRuleUpdate` | 菜单权限 |
| 价格说明修改状态 | `U` | `systemPriceRuleStatus` | 菜单权限 |
| 删除价格说明 | `D` | `systemPriceRuleDelete` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 商品 / 保障服务

- 页面路由：`/product/guarantee`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemGuaranteeLst` | 菜单权限 |
| 添加 | `C` | `systemGuaranteeCreate` | 菜单权限 |
| 编辑 | `U` | `systemGuaranteeUpdate` | 菜单权限 |
| 详情 | `R` | `systemGuaranteeDetail` | 菜单权限 |
| 删除 | `D` | `systemGuaranteeDelete` | 菜单权限 |
| 排序 | `U` | `systemGuaranteeSort` | 菜单权限 |
| 修改状态 | `U` | `systemGuaranteeStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 商品 / 品牌管理

- 页面路由：`/product/brand`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 商品 / 品牌管理 / 品牌分类

- 页面路由：`/product/band/brandClassify`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑 | `U` | `systemStoreBrandCategoryUpdate` | 菜单权限 |
| 列表 | `R` | `systemStoreBrandCategoryLst` | 菜单权限 |
| 详情 | `R` | `systemStoreBrandCategoryDtailt` | 菜单权限 |
| 添加 | `C` | `systemStoreBrandCategoryCreate` | 菜单权限 |
| 删除 | `D` | `systemStoreBrandCategoryDelete` | 菜单权限 |
| 修改状态 | `U` | `systemStoreBrandCategorySwitchStatus` | 菜单权限 |

#### 商品 / 品牌管理 / 品牌列表

- 页面路由：`/product/band/brandList`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemStoreBrandLst` | 菜单权限 |
| 修改状态 | `U` | `systemStoreBrandSwithStatus` | 菜单权限 |
| 添加 | `C` | `systemStoreBrandCreate` | 菜单权限 |
| 编辑 | `U` | `systemStoreBrandUpdate` | 菜单权限 |
| 删除 | `D` | `systemStoreBrandDelete` | 菜单权限 |

#### 商品 / 商品分类

- 页面路由：`/product/classify`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑 | `U` | `systemStoreCategoryUpdate` | 菜单权限 |
| 列表 | `R` | `systemStoreCategoryLst` | 菜单权限 |
| 详情 | `R` | `systemStoreCategoryDtailt` | 菜单权限 |
| 添加 | `C` | `systemStoreCategoryCreate` | 菜单权限 |
| 删除 | `D` | `systemStoreCategoryDelete` | 菜单权限 |
| 修改状态 | `U` | `systemStoreCategorySwitchStatus` | 菜单权限 |
| 修改推荐 | `U` | `systemStoreCategorySwitchIsHot` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 商品 / 商品参数

- 页面路由：`/product/specsMain`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑 | `U` | `systemStoreParameterTemplateUpdate` | 路由 |
| 删除属性 | `D` | `systemStoreParameterTemplateDeleteValue` | 路由 |
| 添加表单 | `C` | `systemStoreCategoryCreateForm` | 路由 |
| 编辑表单 | `U` | `systemStoreCategoryUpdateForm` | 路由 |
| 编辑 | `U` | `systemStoreCategoryUpdate` | 路由 |
| 列表 | `R` | `systemStoreCategoryLst` | 路由 |
| 详情 | `R` | `systemStoreCategoryDtailt` | 路由 |
| 添加 | `C` | `systemStoreCategoryCreate` | 路由 |
| 删除 | `D` | `systemStoreCategoryDelete` | 路由 |
| 修改状态 | `U` | `systemStoreCategorySwitchStatus` | 路由 |
| 修改推荐 | `U` | `systemStoreCategorySwitchIsHot` | 路由 |

#### 商品 / 商品参数 / 平台商品参数

- 页面路由：`/product/specs`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 平台参数列表 | `R` | `systemStoreParameterTemplateLst` | 菜单权限 |
| 详情 | `R` | `systemStoreParameterTemplateDetail` | 菜单权限 |
| 删除 | `D` | `systemStoreParameterTemplateDelete` | 菜单权限 |
| 添加 | `C` | `systemStoreParameterTemplateCreate` | 菜单权限 |
| 编辑 | `U` | `systemStoreParameterTemplateUpdate` | 菜单权限 |
| 删除属性 | `D` | `systemStoreParameterTemplateDeleteValue` | 菜单权限 |

#### 商品 / 商品参数 / 店铺商品参数

- 页面路由：`/product/merSpecs`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商户参数模板 | `R` | `systemStoreParameterTemplateMerLst` | 菜单权限 |

#### 商品 / 商品标签

- 页面路由：`/product/label`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemStoreProductLabelLst` | 菜单权限 |
| 添加 | `C` | `systemStoreProductLabelCreate` | 菜单权限 |
| 编辑 | `U` | `systemStoreProductLabelUpdate` | 菜单权限 |
| 详情 | `R` | `systemStoreProductLabelDetail` | 菜单权限 |
| 删除 | `D` | `systemStoreProductLabelDelete` | 菜单权限 |
| 修改状态 | `U` | `systemStoreProductLabelStatus` | 菜单权限 |

#### 商品 / 商品管理

- 页面路由：`/product/examine`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 统计 | `R` | `systemStoreProductLstFilter` | 菜单权限 |
| 列表 | `R` | `systemStoreProductLst` | 菜单权限 |
| 详情 | `R` | `systemStoreProductDetail` | 菜单权限 |
| 编辑 | `U` | `systemStoreProductUpdate` | 菜单权限 |
| 上下架 | `U` | `systemStoreProductSwitchStatus` | 菜单权限 |
| 分销状态变更商品检测 | `U` | `systemStoreProductCheck` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreProductChangeUsed` | 菜单权限 |
| 虚拟销量 | `C` | `systemStoreProductAddFicti` | 菜单权限 |
| 设置标签 | `U` | `systemStoreProductLabels` | 菜单权限 |
| 批量上下架 | `U` | `systemStoreProductSwitchBatchStatus` | 菜单权限 |
| 批量设置标签 | `U` | `systemStoreProductSwitchBatchLabels` | 菜单权限 |
| 批量设置推荐 | `U` | `systemStoreProductSwitchBatchHot` | 菜单权限 |
| 获取商品操作记录 | `R` | `systemStoreProductGetOperateList` | 菜单权限 |
| 获取自营商品列表 | `R` | `systemStoreProductGetSelfProductList` | 菜单权限 |
| 批量设置分类推荐 | `U` | `systemStoreProductSwitchBatchCateHot` | 菜单权限 |
| 批量复制商品到店铺 | `C` | `systemStoreProductBatchCopy` | 菜单权限 |
| 统计 | `R` | `systemStoreProductLstFilter` | 菜单权限 |
| 列表 | `R` | `systemStoreProductLst` | 菜单权限 |
| 详情 | `R` | `systemStoreProductDetail` | 菜单权限 |
| 编辑 | `U` | `systemStoreProductUpdate` | 菜单权限 |
| 上下架 | `U` | `systemStoreProductSwitchStatus` | 菜单权限 |
| 批量上下架 | `U` | `systemStoreProductSwitchBatchStatus` | 菜单权限 |
| 批量设置标签 | `U` | `systemStoreProductSwitchBatchLabels` | 菜单权限 |
| 批量设置推荐 | `U` | `systemStoreProductSwitchBatchHot` | 菜单权限 |
| 批量设置分类推荐 | `U` | `systemStoreProductSwitchBatchCateHot` | 菜单权限 |
| 分销状态变更商品检测 | `U` | `systemStoreProductCheck` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreProductChangeUsed` | 菜单权限 |
| 虚拟销量 | `C` | `systemStoreProductAddFicti` | 菜单权限 |
| 设置标签 | `U` | `systemStoreProductLabels` | 菜单权限 |
| 获取商品操作记录 | `R` | `systemStoreProductGetOperateList` | 菜单权限 |
| 获取自营商品列表 | `R` | `systemStoreProductGetSelfProductList` | 菜单权限 |
| 批量复制商品到店铺 | `C` | `systemStoreProductBatchCopy` | 菜单权限 |
| 统计 | `R` | `systemStoreProductLstFilter` | 菜单权限 |
| 列表 | `R` | `systemStoreProductLst` | 菜单权限 |
| 详情 | `R` | `systemStoreProductDetail` | 菜单权限 |
| 编辑 | `U` | `systemStoreProductUpdate` | 菜单权限 |
| 上下架 | `U` | `systemStoreProductSwitchStatus` | 菜单权限 |
| 批量上下架 | `U` | `systemStoreProductSwitchBatchStatus` | 菜单权限 |
| 批量设置标签 | `U` | `systemStoreProductSwitchBatchLabels` | 菜单权限 |
| 批量设置推荐 | `U` | `systemStoreProductSwitchBatchHot` | 菜单权限 |
| 批量设置分类推荐 | `U` | `systemStoreProductSwitchBatchCateHot` | 菜单权限 |
| 分销状态变更商品检测 | `U` | `systemStoreProductCheck` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreProductChangeUsed` | 菜单权限 |
| 虚拟销量 | `C` | `systemStoreProductAddFicti` | 菜单权限 |
| 设置标签 | `U` | `systemStoreProductLabels` | 菜单权限 |
| 获取商品操作记录 | `R` | `systemStoreProductGetOperateList` | 菜单权限 |
| 获取自营商品列表 | `R` | `systemStoreProductGetSelfProductList` | 菜单权限 |
| 批量复制商品到店铺 | `C` | `systemStoreProductBatchCopy` | 菜单权限 |

#### 商品 / 活动标签

- 页面路由：`/product/activityLabel`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 商品 / 评论管理

- 页面路由：`/product/comment`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemProductReplyLst` | 菜单权限 |
| 添加虚拟评论 | `C` | `systemProductReplyCreate` | 菜单权限 |
| 排序 | `U` | `systemProductReplySort` | 菜单权限 |
| 删除 | `D` | `systemProductReplyDelete` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

### 客服

#### 客服

- 页面路由：`/service`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 用户与商户聊天记录 | `R` | `adminServiceMerchantUserLogLst` | 路由 |
| 列表 | `R` | `adminServiceReplyLst` | 路由 |
| 添加 | `C` | `adminServiceReplyCreate` | 路由 |
| 编辑 | `U` | `adminServiceReplyUpdate` | 路由 |

#### 客服 / 客服列表

- 页面路由：`/service/customer/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `adminServiceLst` | 菜单权限 |
| 登录 | `O` | `adminServiceLogin` | 菜单权限 |
| 添加 | `C` | `adminServiceCreate` | 菜单权限 |
| 编辑 | `U` | `adminServiceUpdate` | 菜单权限 |
| 修改状态 | `U` | `adminServiceSwitchStatus` | 菜单权限 |
| 删除 | `D` | `adminServiceDelete` | 菜单权限 |
| 客服的全部用户  | `O` | `adminServiceServiceUserList` | 菜单权限 |
| 用户与客服聊天记录 | `R` | `adminServiceServiceUserLogLst` | 菜单权限 |
| 客服的聊天用户列表 | `R` | `adminServiceServiceMerchantUserList` | 菜单权限 |
| 用户与商户聊天记录 | `R` | `adminServiceMerchantUserLogLst` | 菜单权限 |

#### 客服 / 客服自动回复

- 页面路由：`/systemForm/customer_keyword`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `adminServiceReplyLst` | 菜单权限 |
| 添加 | `C` | `adminServiceReplyCreate` | 菜单权限 |
| 编辑 | `U` | `adminServiceReplyUpdate` | 菜单权限 |
| 切换状态 | `U` | `adminServiceReplyStatus` | 菜单权限 |
| 删除 | `D` | `adminServiceReplyDelete` | 菜单权限 |

#### 客服 / 客服设置

- 页面路由：`/systemForm/Basics/service`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

### 应用

#### 应用

- 页面路由：`/apploction`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 应用 / 公众号

- 页面路由：`/app/wechat`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 微信菜单配置 | `R` | `wechatMenu` | 路由 |
| 保存微信菜单配置 | `C` | `saveWechatMenu` | 路由 |
| 详情 | `R` | `wechatReplyInfo` | 路由 |
| 编辑 | `CU` | `saveWechatReply` | 路由 |
| 添加 | `C` | `createWechatReply` | 路由 |
| 修改 | `U` | `updateWechatReply` | 路由 |
| 列表 | `R` | `wechatReplyLst` | 路由 |
| 删除 | `D` | `wechatReplyDelete` | 路由 |
| 修改状态 | `U` | `wechatReplyStatus` | 路由 |

#### 应用 / 公众号 / 图文管理

- 页面路由：`/app/wechat/newsCategory`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 添加 | `C` | `systemWechatNewsCreate` | 菜单权限 |
| 编辑 | `U` | `systemWechatNewsUpdate` | 菜单权限 |
| 删除 | `D` | `systemWechatNewsDelete` | 菜单权限 |
| 列表 | `R` | `systemWechatNewsLst` | 菜单权限 |
| 详情 | `R` | `systemWechatNewsDetail` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 应用 / 公众号 / 微信模板消息

- 页面路由：`/app/wechat/template`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 同步 | `U` | `systemTemplateMessageSync` | 菜单权限 |
| 列表 | `R` | `systemTemplateMessageLst` | 菜单权限 |
| 添加 | `C` | `systemTemplateMessageCreate` | 菜单权限 |
| 编辑 | `U` | `systemTemplateMessageUpdate` | 菜单权限 |
| 删除 | `D` | `systemTemplateMessageDelete` | 菜单权限 |
| 修改状态 | `U` | `systemTemplateMessageSwitchStatus` | 菜单权限 |

#### 应用 / 公众号 / 微信菜单

- 页面路由：`/app/wechat/menus`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 微信菜单配置 | `U` | `wechatMenu` | 菜单权限 |
| 保存微信菜单配置 | `CU` | `saveWechatMenu` | 菜单权限 |

#### 应用 / 公众号 / 自动回复

- 页面路由：`/admin/app/wechat/reply`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 详情 | `R` | `wechatReplyInfo` | 菜单权限 |
| 编辑 | `CU` | `saveWechatReply` | 菜单权限 |
| 添加 | `C` | `createWechatReply` | 菜单权限 |
| 修改 | `U` | `updateWechatReply` | 菜单权限 |
| 列表 | `R` | `wechatReplyLst` | 菜单权限 |
| 删除 | `D` | `wechatReplyDelete` | 菜单权限 |
| 修改状态 | `U` | `wechatReplyStatus` | 菜单权限 |
| 上传图片 | `C` | `wechatUploadImage` | 菜单权限 |
| 上传语音 | `C` | `wechatUploadVoice` | 菜单权限 |

#### 应用 / 公众号 / 自动回复 / 关键字回复

- 页面路由：`/app/wechat/reply/keyword`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 应用 / 公众号 / 自动回复 / 微信关注回复

- 页面路由：`/app/wechat/reply/follow/subscribe`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 应用 / 公众号 / 自动回复 / 无效关键词回复

- 页面路由：`/app/wechat/reply/index/default`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 应用 / 小程序

- 页面路由：`/app/routine`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 小程序配置 | `R` | `configRoutineConfig` | 路由 |
| 小程序下载 | `R` | `configRoutineDownload` | 路由 |

#### 应用 / 小程序 / 小程序下载

- 页面路由：`/app/routine/download`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 小程序下载 | `R` | `configRoutineDownload` | 菜单权限 |

#### 应用 / 小程序 / 小程序订阅消息

- 页面路由：`/app/routine/template`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 同步 | `U` | `systemTemplateMessageMinSync` | 菜单权限 |
| 列表  | `R` | `systemTemplateMessageMinList` | 菜单权限 |
| 添加 | `C` | `systemTemplateMessageMinCreate` | 菜单权限 |
| 编辑 | `U` | `systemTemplateMessageMinUpdate` | 菜单权限 |
| 删除 | `D` | `systemTemplateMessageMinDelete` | 菜单权限 |
| 修改状态 | `U` | `systemTemplateMessageMinSwitchStatus` | 菜单权限 |

### 店铺

#### 店铺

- 页面路由：`/mer`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 分账商户申请备注 | `C` | `systemMerchantApplymentsMarrkSave` | 路由 |
| 列表 | `R` | `systemOrderProfitsharingLst` | 路由 |

#### 店铺 / 区域代理

- 页面路由：`/business-zones/manage`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 店铺 / 区域代理 / 代理人员

- 页面路由：`/business-zones/agents`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商圈代理添加 | `C` | `systemCircleAgentCreate` | 菜单权限 |
| 商圈代理编辑 | `U` | `systemCircleAgentUpdate` | 菜单权限 |
| 商圈代理删除 | `D` | `systemCircleAgentDelete` | 菜单权限 |
| 关联商户列表 | `R` | `systemCircleAgentMerchant` | 菜单权限 |
| 代理选项 | `O` | `systemCircleAgentOptions` | 菜单权限 |
| 重置密码 | `O` | `systemCircleAgentResetPassword` | 菜单权限 |

#### 店铺 / 区域代理 / 代理审核

- 页面路由：`/business-zones/agent-review`
- CRUD：C=— R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商圈代理审核 | `U` | `systemCircleAgentAudit` | 菜单权限 |

#### 店铺 / 区域代理 / 代理设置

- 页面路由：`/business-zones/settings`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 店铺 / 区域代理 / 区域列表

- 页面路由：`/business-zones/index`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商圈列表 | `R` | `systemCircleList` | 菜单权限 |
| 商圈详情 | `R` | `systemCircleDetail` | 菜单权限 |
| 商圈添加 | `C` | `systemCircleCreate` | 菜单权限 |
| 商圈编辑 | `U` | `systemCircleUpdate` | 菜单权限 |
| 商圈删除 | `D` | `systemCircleDelete` | 菜单权限 |
| 商圈状态切换 | `U` | `systemCircleSwitch` | 菜单权限 |
| 关联商户列表 | `R` | `systemCircleMerchantList` | 菜单权限 |

#### 店铺 / 商户管理

- 页面路由：`/merchant`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 分账商户申请备注 | `C` | `systemMerchantApplymentsMarrkSave` | 路由 |
| 列表 | `R` | `systemOrderProfitsharingLst` | 路由 |

#### 店铺 / 商户管理 / 商户入驻审核

- 页面路由：`/merchant/review`
- CRUD：C=— R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商圈代理审核 | `U` | `systemCircleAgentAudit` | 菜单权限 |

#### 店铺 / 商户管理 / 商户列表

- 页面路由：`/merchant/index`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商圈列表 | `R` | `systemCircleList` | 菜单权限 |
| 商圈详情 | `R` | `systemCircleDetail` | 菜单权限 |
| 商圈添加 | `C` | `systemCircleCreate` | 菜单权限 |
| 商圈编辑 | `U` | `systemCircleUpdate` | 菜单权限 |
| 商圈删除 | `D` | `systemCircleDelete` | 菜单权限 |
| 商圈状态切换 | `U` | `systemCircleSwitch` | 菜单权限 |
| 关联商户列表 | `R` | `systemCircleMerchantList` | 菜单权限 |

#### 店铺 / 商户管理 / 商户管理员

- 页面路由：`/merchant/admin-list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商圈代理添加 | `C` | `systemCircleAgentCreate` | 菜单权限 |
| 商圈代理编辑 | `U` | `systemCircleAgentUpdate` | 菜单权限 |
| 商圈代理删除 | `D` | `systemCircleAgentDelete` | 菜单权限 |
| 关联商户列表 | `R` | `systemCircleAgentMerchant` | 菜单权限 |
| 代理选项 | `O` | `systemCircleAgentOptions` | 菜单权限 |
| 重置密码 | `O` | `systemCircleAgentResetPassword` | 菜单权限 |

#### 店铺 / 商户管理 / 商户设置

- 页面路由：`/merchant/apply-setting`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 店铺 / 店铺管理

- 页面路由：`/mer/mer`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 店铺 / 店铺管理 / 商户列表

- 页面路由：`/merchant/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 虚拟关注量 | `O` | `systemMerchantCareFicti` | 菜单权限 |
| 商户列表 | `CR` | `systemMerchantCreateForm` | 菜单权限 |
| 商户列表统计 | `R` | `systemMerchantCount` | 菜单权限 |
| 商户列表 | `R` | `systemMerchantLst` | 菜单权限 |
| 商户添加 | `C` | `systemMerchantCreate` | 菜单权限 |
| 商户编辑 | `U` | `systemMerchantUpdate` | 菜单权限 |
| 商户修改推荐 | `U` | `systemMerchantStatus` | 菜单权限 |
| 商户开启/关闭 | `O` | `systemMerchantClose` | 菜单权限 |
| 商户删除 | `D` | `systemMerchantDelete` | 菜单权限 |
| 商户修改密码 | `U` | `systemMerchantAdminPassword` | 菜单权限 |
| 商户登录 | `O` | `systemMerchantLogin` | 菜单权限 |
| 修改采集商品次数 | `CU` | `systemMerchantChangeCopy` | 菜单权限 |
| 详情 | `R` | `systemMerchantDetail` | 菜单权限 |
| 操作日志 | `R` | `systemMerchantOperateList` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |
| 商户添加店铺 | `C` | `systemMerchantBusinessCreate` | 菜单权限 |

#### 店铺 / 店铺管理 / 店铺入驻申请

- 页面路由：`/merchant/application`
- CRUD：C=— R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemMerchantIntentionLst` | 菜单权限 |
| 审核 | `U` | `systemMerchantIntentionStatus` | 菜单权限 |
| 删除 | `D` | `systemMerchantIntentionDelete` | 菜单权限 |
| 备注 | `U` | `systemMerchantIntentionMark` | 菜单权限 |
| 列表 | `R` | `systemMerchantIntentionLst` | 菜单权限 |
| 审核 | `U` | `systemMerchantIntentionStatus` | 菜单权限 |
| 删除 | `D` | `systemMerchantIntentionDelete` | 菜单权限 |
| 备注 | `U` | `systemMerchantIntentionMark` | 菜单权限 |

#### 店铺 / 店铺管理 / 店铺分类

- 页面路由：`/merchant/classify`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商户分类列表 | `R` | `systemMerchantCategoryLst` | 菜单权限 |
| 商户分类添加 | `C` | `systemMerchantCategoryCreate` | 菜单权限 |
| 商户分类删除 | `D` | `systemMerchantCategoryDelete` | 菜单权限 |
| 商户分类编辑 | `U` | `systemMerchantCategoryUpdate` | 菜单权限 |

#### 店铺 / 店铺管理 / 店铺分组

- 页面路由：`/merchant/grouping`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 添加 | `C` | `systemMerchantRegionCreate` | 菜单权限 |
| 编辑 | `U` | `systemMerchantRegionUpdate` | 菜单权限 |
| 列表 | `R` | `systemMerchantRegionLst` | 菜单权限 |
| 删除 | `D` | `systemMerchantRegionDelete` | 菜单权限 |
| 列表 | `R` | `systemStoreGroupLst` | 菜单权限 |
| 详情 | `R` | `systemStoreGroupDetail` | 菜单权限 |
| 添加 | `C` | `systemStoreGroupCreate` | 菜单权限 |
| 编辑 | `U` | `systemStoreGroupUpdate` | 菜单权限 |
| 删除 | `D` | `systemStoreGroupDelete` | 菜单权限 |
| 状态切换 | `U` | `systemStoreGroupSwitchStatus` | 菜单权限 |
| 设置店铺分组模板 | `U` | `systemStoreGroupSetTemplate` | 菜单权限 |
| 关联店铺列表 | `R` | `systemStoreGroupStores` | 菜单权限 |

实现闭环（统一后台）：平台角色的 `merchant.group.manage` 按钮权限覆盖列表、详情、新增、编辑、删除、启停、模板绑定和关联店铺查看；数据存入 `qixi_crm_a_store_group` 与 `qixi_crm_a_store_group_merchant`。服务端在同一事务内校验商户投影存在性、父子循环、最多三级、移动子树层级与路径、删除非叶节点和状态向子分组级联。`sql/admin/init_test_data.sql` 提供不含真实个人信息的中文树与店铺关联夹具；运行态验收须在隔离数据库应用夹具后执行。

#### 店铺 / 店铺管理 / 店铺分账申请

- 页面路由：`/merchant/applyments`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 分账商户申请列表 | `CR` | `systemMerchantApplymentsLst` | 菜单权限 |
| 分账商户申请详情 | `CR` | `systemMerchantApplymentsDetail` | 菜单权限 |
| 分账商户申请审核 | `CU` | `systemMerchantApplymentsStatus` | 菜单权限 |
| 分账商户申请备注 | `CU` | `systemMerchantApplymentsMarrkSave` | 菜单权限 |

#### 店铺 / 店铺管理 / 店铺列表

- 页面路由：`/merchant/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商户列表 | `CR` | `systemMerchantCreateForm` | 菜单权限 |
| 商户列表统计 | `R` | `systemMerchantCount` | 菜单权限 |
| 商户列表 | `R` | `systemMerchantLst` | 菜单权限 |
| 商户添加 | `C` | `systemMerchantCreate` | 菜单权限 |
| 商户编辑 | `U` | `systemMerchantUpdate` | 菜单权限 |
| 商户修改推荐 | `U` | `systemMerchantStatus` | 菜单权限 |
| 商户开启/关闭 | `O` | `systemMerchantClose` | 菜单权限 |
| 商户删除 | `D` | `systemMerchantDelete` | 菜单权限 |
| 商户修改密码 | `U` | `systemMerchantAdminPassword` | 菜单权限 |
| 商户登录 | `O` | `systemMerchantLogin` | 菜单权限 |
| 修改采集商品次数 | `CU` | `systemMerchantChangeCopy` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |
| 详情 | `R` | `systemMerchantDetail` | 菜单权限 |
| 操作日志 | `R` | `systemMerchantOperateList` | 菜单权限 |
| 虚拟关注量 | `O` | `systemMerchantCareFicti` | 菜单权限 |
| 商户添加店铺 | `C` | `systemMerchantBusinessCreate` | 菜单权限 |
| 虚拟关注量 | `O` | `systemMerchantCareFicti` | 菜单权限 |
| 店铺列表 | `CR` | `systemMerchantCreateForm` | 菜单权限 |
| 店铺列表统计 | `R` | `systemMerchantCount` | 菜单权限 |
| 店铺列表 | `R` | `systemMerchantLst` | 菜单权限 |
| 店铺添加 | `C` | `systemMerchantCreate` | 菜单权限 |
| 店铺编辑 | `U` | `systemMerchantUpdate` | 菜单权限 |
| 店铺修改推荐 | `U` | `systemMerchantStatus` | 菜单权限 |
| 店铺开启/关闭 | `O` | `systemMerchantClose` | 菜单权限 |
| 店铺删除 | `D` | `systemMerchantDelete` | 菜单权限 |
| 店铺修改密码 | `U` | `systemMerchantAdminPassword` | 菜单权限 |
| 店铺登录 | `O` | `systemMerchantLogin` | 菜单权限 |
| 修改采集商品次数 | `CU` | `systemMerchantChangeCopy` | 菜单权限 |
| 详情 | `R` | `systemMerchantDetail` | 菜单权限 |
| 操作日志 | `R` | `systemMerchantOperateList` | 菜单权限 |
| 商户添加店铺 | `C` | `systemMerchantBusinessCreate` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 店铺 / 店铺管理 / 店铺类型

- 页面路由：`/merchant/type`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemMerchantTypeLst` | 菜单权限 |
| 添加 | `C` | `systemMerchantTypeCreate` | 菜单权限 |
| 编辑 | `U` | `systemMerchantTypeUpdate` | 菜单权限 |
| 删除 | `D` | `systemMerchantTypeDelete` | 菜单权限 |
| 备注 | `U` | `systemMerchantTypeMark` | 菜单权限 |
| 备注 | `RU` | `systemMerchantTypeDetail` | 菜单权限 |

实现闭环（统一后台）：平台角色以 `merchant.type.manage` 管理类型、保证金规则、类型说明、备注、状态及店铺菜单代码授权。数据使用 `qixi_crm_a_merchant_type` 与 `qixi_crm_a_merchant_type_menu`，保证金启用时必须大于零，未启用时金额归零；中文夹具仅用于隔离本地验收。

#### 店铺 / 店铺设置

- 页面路由：`/mer/store`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 店铺 / 店铺设置 / 保证金配置

- 页面路由：`/systemForm/Basics/margin`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

#### 店铺 / 店铺设置 / 店铺保证金

- 页面路由：`/merchant/deposit_list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 缴纳记录 | `R` | `systemMerchantMarginLst` | 菜单权限 |
| 扣费记录 | `R` | `systemMarginList` | 菜单权限 |
| 扣除保证金 | `O` | `systemMarginSet` | 菜单权限 |
| 退款申请列表 | `CRD` | `systemMarginRefundList` | 菜单权限 |
| 退款申请详情 | `CRUD` | `systemMarginRefundShow` | 菜单权限 |
| 审核 | `U` | `systemMarginRefundSwitchStatus` | 菜单权限 |
| 备注 | `U` | `systemMarginRefundMark` | 菜单权限 |
| 待缴列表 | `R` | `systemMarginMakeUpMarginLst` | 菜单权限 |

#### 店铺 / 店铺设置 / 店铺菜单

- 页面路由：`/merchant/system`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商户菜单/权限列表 | `R` | `systemMerchantMenuGetLst` | 菜单权限 |
| 商户菜单/权限添加 | `C` | `systemMerchantMenuCreate` | 菜单权限 |
| 商户菜单/权限编辑 | `U` | `systemMerchantMenuUpdate` | 菜单权限 |
| 商户菜单/权限删除 | `D` | `systemMerchantMenuDelete` | 菜单权限 |

#### 店铺 / 店铺设置 / 说明提示

- 页面路由：`/merchant/type/description`
- CRUD：C=✓ R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 店铺类型说明  | `C` | `systemAgreeSave` | 菜单权限 |

统一后台等价入口：`/merchant/type/description` 映射至协议设置页中的保证金说明键 `sys_deposit_agree`；`/merchant/deposit_list` 映射至店铺保证金工作台，避免旧路径落入空白页。

### 用户

#### 用户

- 页面路由：`/user`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 详情关联列表 | `R` | `systemCouponShowLst` | 路由 |
| 发送优惠券 | `C` | `systemCouponSend` | 路由 |
| 批量设置分销员 | `R` | `getMemberLevelBatchSpreadForm` | 路由 |
| 批量设置分销员 | `R` | `getMemberLevelBatchSpread` | 路由 |
| 用户搜索记录 | `R` | `systemUserSearchLog` | 路由 |
| 清除用户搜索记录 | `R` | `systemUserClearSearchLog` | 路由 |
| 用户搜索记录导出 | `R` | `systemUserExportSearchLog` | 路由 |

#### 用户 / 付费会员

- 页面路由：`/user/svip`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 用户 / 付费会员 / 付费会员配置

- 页面路由：`/systemForm/Basics/svip`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

#### 用户 / 付费会员 / 会员协议

- 页面路由：`/user/member/vipAgreement`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 用户 / 付费会员 / 会员权益

- 页面路由：`/user/member/equity`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemUserSvipInterestsLst` | 菜单权限 |
| 编辑 | `U` | `systemUserSvipInterestsUpdate` | 菜单权限 |
| 编辑状态 | `U` | `systemUserSvipInterestsStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 用户 / 付费会员 / 会员类型

- 页面路由：`/user/member/type`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemUserSvipLst` | 菜单权限 |
| 添加 | `C` | `systemUserSvipCreate` | 菜单权限 |
| 编辑表单 | `U` | `systemUserSvipUpdateForm` | 菜单权限 |
| 编辑 | `U` | `systemUserSvipTypeUpdate` | 菜单权限 |
| 删除 | `D` | `systemUserSvipDelete` | 菜单权限 |
| 修改状态 | `U` | `systemUserSvipStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 用户 / 付费会员 / 会员记录

- 页面路由：`/user/member/record`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemUserSvipPayLst` | 路由 |
| 统计 | `R` | `systemUserSvipCountInfo` | 路由 |
| 列表 | `R` | `systemUserSvipInterestsLst` | 路由 |
| 编辑 | `U` | `systemUserSvipInterestsUpdateForm` | 路由 |

#### 用户 / 搜索记录

- 页面路由：`/user/searchRecord`
- CRUD：C=— R=✓ U=— D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 用户搜索记录 | `R` | `systemUserSearchLog` | 菜单权限 |
| 用户搜索记录导出 | `R` | `systemUserExportSearchLog` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `systemStoreExcelDownload` | 菜单权限 |
| 清除用户搜索记录 | `RD` | `systemUserClearSearchLog` | 菜单权限 |

#### 用户 / 用户分组

- 页面路由：`/user/group`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 用户分组列表 | `R` | `systemUserGroupLst` | 菜单权限 |
| 用户分组添加 | `C` | `systemUserGroupCreate` | 菜单权限 |
| 用户分组删除 | `D` | `systemUserGroupDelete` | 菜单权限 |
| 用户分组编辑 | `U` | `systemUserGroupUpdate` | 菜单权限 |

#### 用户 / 用户列表

- 页面路由：`/user/list`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 发送优惠券 | `C` | `systemCouponSend` | 菜单权限 |
| 用户列表 | `R` | `systemUserLst` | 菜单权限 |
| 用户编辑 | `U` | `systemUserUpdate` | 菜单权限 |
| 用户修改余额 | `U` | `systemUserChangeNowMoney` | 菜单权限 |
| 用户修改积分 | `U` | `systemUserChangeIntegral` | 菜单权限 |
| 用户发送图文 | `C` | `systemWechatUserSendNews` | 菜单权限 |
| 用户详情 | `R` | `systemUserDetail` | 菜单权限 |
| 用户消费记录 | `R` | `systemUserOrder` | 菜单权限 |
| 用户持有优惠券 | `O` | `systemUserCoupon` | 菜单权限 |
| 用户余额变动列表 | `R` | `systemUserBill` | 菜单权限 |
| 推荐人修改记录 | `RU` | `systemUserSpreadLog` | 菜单权限 |
| 修改推荐人 | `RU` | `systemUserSpreadChange` | 菜单权限 |
| 用户修改会员等级 | `CU` | `systemUserMemberSave` | 菜单权限 |
| 用户添加 | `C` | `systemUserCreate` | 菜单权限 |
| 用户修改密码 | `U` | `systemUserChangePassword` | 菜单权限 |
| 用户分组编辑 | `U` | `systemUserChangeGroup` | 菜单权限 |
| 用户分组批量编辑 | `U` | `systemUserBatchChangeGroup` | 菜单权限 |
| 用户标签编辑 | `U` | `systemUserChangeLabel` | 菜单权限 |
| 用户标签批量编辑 | `U` | `systemUserBatchChangeLabel` | 菜单权限 |
| 优惠券列表 | `R` | `systemCouponList` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |
| 用户标签编辑 | `U` | `systemUserSvipUpdate` | 菜单权限 |
| 积分记录 | `R` | `systemUserIntegralList` | 菜单权限 |
| 签到记录 | `R` | `systemUserSginLog` | 菜单权限 |
| 浏览记录 | `R` | `systemUserHistory` | 菜单权限 |
| 用户信息导出 | `R` | `systemUserExcel` | 菜单权限 |
| 批量设置分销员 | `RU` | `getMemberLevelBatchSpread` | 菜单权限 |

#### 用户 / 用户协议

- 页面路由：`/user/agreement`
- CRUD：C=✓ R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 用户协议 | `C` | `systemAgreeSave` | 菜单权限 |

#### 用户 / 用户反馈

- 页面路由：`/feedback`
- CRUD：C=— R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 修改状态 | `U` | `systemUserFeedBackCategorySwitchStatus` | 路由 |
| 删除 | `D` | `systemUserFeedBackCategoryDelete` | 路由 |
| 列表 | `R` | `systemUserFeedBackLst` | 路由 |
| 详情 | `R` | `systemUserFeedBackDetail` | 路由 |
| 回复表单 | `R` | `systemUserFeedBackReplyForm` | 路由 |

#### 用户 / 用户反馈 / 反馈分类

- 页面路由：`/feedback/classify`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemUserFeedBackCategoryLst` | 菜单权限 |
| 添加 | `C` | `systemUserFeedBackCategoryCreate` | 菜单权限 |
| 编辑 | `U` | `systemUserFeedBackCategoryUpdate` | 菜单权限 |
| 修改状态 | `U` | `systemUserFeedBackCategorySwitchStatus` | 菜单权限 |
| 删除 | `D` | `systemUserFeedBackCategoryDelete` | 菜单权限 |

#### 用户 / 用户反馈 / 反馈列表

- 页面路由：`/feedback/list`
- CRUD：C=— R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemUserFeedBackLst` | 菜单权限 |
| 详情 | `R` | `systemUserFeedBackDetail` | 菜单权限 |
| 回复 | `U` | `systemUserFeedBackReply` | 菜单权限 |
| 删除 | `D` | `systemUserFeedBackDelete` | 菜单权限 |

#### 用户 / 用户标签

- 页面路由：`/user/label`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 用户标签列表 | `R` | `systemUserLabelLst` | 菜单权限 |
| 用户标签添加 | `C` | `systemUserLabelCreate` | 菜单权限 |
| 用户标签删除 | `D` | `systemUserLabelDelete` | 菜单权限 |
| 用户标签编辑 | `U` | `systemUserLabelUpdate` | 菜单权限 |

#### 用户 / 用户等级

- 页面路由：`/user/member`
- CRUD：C=— R=— U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 删除 | `D` | `systemUserSvipDelete` | 路由 |
| 修改状态 | `U` | `systemUserSvipStatus` | 路由 |

#### 用户 / 用户等级 / 等级权益

- 页面路由：`/user/member/interests`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 会员权益 | `R` | `systemUserMemberInterestsLst` | 菜单权限 |
| 会员权益添加 | `C` | `systemUserMemberInterestsCreate` | 菜单权限 |
| 会员权益编辑 | `U` | `systemUserMemberInterestsUpdate` | 菜单权限 |
| 会员权益删除 | `D` | `systemUserMemberInterestsDelete` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 用户 / 用户等级 / 等级管理

- 页面路由：`/user/member/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 普通会员等级列表 | `R` | `systemUserMemberLst` | 菜单权限 |
| 普通会员等级添加 | `C` | `systemUserMemberCreate` | 菜单权限 |
| 普通会员等级编辑 | `U` | `systemUserMemberUpdate` | 菜单权限 |
| 普通会员等级删除 | `D` | `systemUserMemberDelete` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 用户 / 用户等级 / 等级说明

- 页面路由：`/user/member/description`
- CRUD：C=✓ R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 会员等级规则 | `C` | `systemAgreeSave` | 菜单权限 |

#### 用户 / 用户等级 / 等级配置

- 页面路由：`/systemForm/Basics/members`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

#### 用户 / 用户设置

- 页面路由：`/user/setup_user`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemUserInfolst` | 菜单权限 |
| 添加 | `CR` | `systemUserInfoCreate` | 菜单权限 |
| 保存信息 | `CR` | `systemUserInfoSaveAll` | 菜单权限 |
| 删除 | `RD` | `systemUserInfoDelete` | 菜单权限 |
| 保存注册配置 | `CU` | `systemUserRegisterConfig` | 菜单权限 |
| 新人礼优惠券列表 | `R` | `systemUserRegisterCoupon` | 菜单权限 |
| 扩展信息表单 | `C` | `systemUserFieldSaveForm` | 菜单权限 |
| 添加或编辑 | `CRU` | `systemUserInfoFieldSave` | 菜单权限 |

### 维护

#### 维护

- 页面路由：`/safe`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 维护 / 安全维护

- 页面路由：`/maintain`
- CRUD：C=— R=✓ U=— D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 数据库备份下载 | `R` | `systemSafetyDatabaseDownloadFile` | 路由 |
| 数据库备份删除 | `D` | `systemSafetyDatabaseDeleteFile` | 路由 |

#### 维护 / 安全维护 / 商业授权

- 页面路由：`/setting/system/maintain/auth`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 获取去版权信息 | `CR` | `systemCopyright` | 菜单权限 |
| 获取授权信息 | `CR` | `systemAuthCopyright` | 菜单权限 |
| 保存去版权信息 | `C` | `systemSaveCopyright` | 菜单权限 |

#### 维护 / 安全维护 / 数据备份

- 页面路由：`/maintain/dataBackup`
- CRUD：C=— R=✓ U=— D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 数据库列表 | `R` | `systemSafetyDatabaseLst` | 菜单权限 |
| 数据库备份列表 | `R` | `systemSafetyDatabaseFileList` | 菜单权限 |
| 数据库备份详情 | `R` | `systemSafetyDatabaseDetail` | 菜单权限 |
| 备份 | `O` | `systemSafetyDatabaseBackups` | 菜单权限 |
| 数据库优化 | `O` | `systemSafetyDatabaseOptimize` | 菜单权限 |
| 数据库维护 | `O` | `systemSafetyDatabaseRepair` | 菜单权限 |
| 数据库备份下载 | `R` | `systemSafetyDatabaseDownloadFile` | 菜单权限 |
| 数据库备份删除 | `D` | `systemSafetyDatabaseDeleteFile` | 菜单权限 |

#### 维护 / 安全维护 / 缓存清除

- 页面路由：`/maintain/cache`
- CRUD：C=— R=— U=— D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 清除缓存 | `D` | `systemClearCache` | 菜单权限 |
| 替换素材域名 | `O` | `systemAttachmentReplaceHost` | 菜单权限 |

#### 维护 / 导出记录

- 页面路由：`/group/exportList`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 下载 | `R` | `systemStoreExcelDownload` | 菜单权限 |

#### 维护 / 开发配置

- 页面路由：`/safe/exploit`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 维护 / 开发配置 / 组合数据

- 页面路由：`/group/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 组合数据配置列表 | `RU` | `groupLst` | 菜单权限 |
| 组合数据配置添加 | `CU` | `groupCreate` | 菜单权限 |
| 组合数据配置编辑 | `U` | `groupUpdate` | 菜单权限 |
| 详情 | `R` | `groupDetail` | 菜单权限 |
| 列表 | `R` | `groupDataLst` | 菜单权限 |
| 添加 | `C` | `groupDataCreate` | 菜单权限 |
| 编辑 | `U` | `groupDataUpdate` | 菜单权限 |
| 删除 | `D` | `groupDataDelete` | 菜单权限 |
| 修改状态 | `U` | `groupDataChangeStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 维护 / 操作日志

- 页面路由：`/setting/systemLog`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 操作日志 | `R` | `systemAdminLog` | 菜单权限 |

#### 维护 / 配置分类

- 页面路由：`/config/classify`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 配置分类添加 | `CU` | `configClassifyCreate` | 菜单权限 |
| 配置分类删除 | `UD` | `configClassifyDelete` | 菜单权限 |
| 配置分类编辑 | `U` | `configClassifyUpdate` | 菜单权限 |
| 配置分类修改状态 | `U` | `configClassifySwitchStatus` | 菜单权限 |
| 配置分类列表 | `RU` | `configClassifyLst` | 菜单权限 |

#### 维护 / 配置管理

- 页面路由：`/config/setting`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 配置添加 | `CU` | `configSettingCreate` | 菜单权限 |
| 配置编辑 | `U` | `configSettingUpdate` | 菜单权限 |
| 配置修改状态 | `U` | `configSettingSwitchStatus` | 菜单权限 |
| 配置列表 | `RU` | `configSettingLst` | 菜单权限 |
| 配置删除 | `UD` | `configSettingDelete` | 菜单权限 |

#### 维护 / 页面链接

- 页面路由：`/safe/pageLinks`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

### 营销

#### 营销

- 页面路由：`/marketing`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemStoreCouponLst` | 路由 |
| 使用记录 | `R` | `systemCouponIssue` | 路由 |

#### 营销 / 专场列表

- 页面路由：`/group/topic/94`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 详情 | `R` | `groupDetail` | 菜单权限 |
| 列表 | `R` | `groupDataLst` | 菜单权限 |
| 添加 | `C` | `groupDataCreate` | 菜单权限 |
| 编辑 | `U` | `groupDataUpdate` | 菜单权限 |
| 删除 | `D` | `groupDataDelete` | 菜单权限 |
| 修改状态 | `U` | `groupDataChangeStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 营销 / 优惠套餐

- 页面路由：`/marketing/discounts/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 优惠套餐列表 | `R` | `systemStoreDiscountsLst` | 菜单权限 |
| 优惠套餐详情 | `R` | `systemStoreDiscountsDetail` | 菜单权限 |
| 优惠套餐修改状态 | `U` | `systemStoreDiscountsStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |
| 价格说明列表 | `R` | `systemPriceRuleLst` | 菜单权限 |
| 添加价格说明 | `C` | `systemPriceRuleCreate` | 菜单权限 |
| 修改价格说明 | `U` | `systemPriceRuleUpdate` | 菜单权限 |
| 价格说明修改状态 | `U` | `systemPriceRuleStatus` | 菜单权限 |
| 删除价格说明 | `D` | `systemPriceRuleDelete` | 菜单权限 |

#### 营销 / 余额充值

- 页面路由：`/banlace`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 营销 / 余额充值 / 余额充值配置

- 页面路由：`/group/config/69`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 详情 | `R` | `groupDetail` | 菜单权限 |
| 列表 | `R` | `groupDataLst` | 菜单权限 |
| 添加 | `C` | `groupDataCreate` | 菜单权限 |
| 编辑 | `U` | `groupDataUpdate` | 菜单权限 |
| 删除 | `D` | `groupDataDelete` | 菜单权限 |
| 修改状态 | `U` | `groupDataChangeStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 营销 / 余额充值 / 余额设置

- 页面路由：`/systemForm/Basics/balance`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

#### 营销 / 助力

- 页面路由：`/assist`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 营销 / 助力 / 助力活动

- 页面路由：`/marketing/assist/list`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemStoreProductAssistSetLst` | 菜单权限 |
| 详情 | `R` | `systemStoreProductAssistSetDetail` | 菜单权限 |
| 列表 | `R` | `systemStoreProductAssistSetLst` | 菜单权限 |
| 详情 | `R` | `systemStoreProductAssistSetDetail` | 菜单权限 |
| 列表 | `R` | `systemStoreProductAssistSetLst` | 菜单权限 |
| 详情 | `R` | `systemStoreProductAssistSetDetail` | 菜单权限 |

#### 营销 / 助力 / 活动商品

- 页面路由：`/marketing/assist/goods_list`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemStoreProductAssistLst` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreProductAssistShow` | 菜单权限 |
| 详情 | `R` | `systemStoreProductAssistDetail` | 菜单权限 |
| 编辑 | `U` | `systemStoreProductAssistProductUpdate` | 菜单权限 |
| 审核 | `U` | `systemStoreProductAssistStatus` | 菜单权限 |
| 编辑数据 | `RU` | `systemStoreProductAssistGet` | 菜单权限 |
| 设置标签 | `U` | `systemStoreProductAssistLabels` | 菜单权限 |
| 列表 | `R` | `systemStoreProductAssistLst` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreProductAssistShow` | 菜单权限 |
| 详情 | `R` | `systemStoreProductAssistDetail` | 菜单权限 |
| 编辑 | `U` | `systemStoreProductAssistProductUpdate` | 菜单权限 |
| 审核 | `U` | `systemStoreProductAssistStatus` | 菜单权限 |
| 编辑数据 | `RU` | `systemStoreProductAssistGet` | 菜单权限 |
| 设置标签 | `U` | `systemStoreProductAssistLabels` | 菜单权限 |
| 列表 | `R` | `systemStoreProductAssistLst` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreProductAssistShow` | 菜单权限 |
| 详情 | `R` | `systemStoreProductAssistDetail` | 菜单权限 |
| 编辑 | `U` | `systemStoreProductAssistProductUpdate` | 菜单权限 |
| 审核 | `U` | `systemStoreProductAssistStatus` | 菜单权限 |
| 编辑数据 | `RU` | `systemStoreProductAssistGet` | 菜单权限 |
| 设置标签 | `U` | `systemStoreProductAssistLabels` | 菜单权限 |

#### 营销 / 商户优惠券

- 页面路由：`/marketing/coupon`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemStoreCouponLst` | 路由 |
| 使用记录 | `R` | `systemCouponIssue` | 路由 |

#### 营销 / 商户优惠券 / 优惠券列表

- 页面路由：`/marketing/coupon/list`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemStoreCouponLst` | 菜单权限 |
| 详情 | `R` | `systemCouponDetail` | 菜单权限 |
| 商品列表 | `R` | `systemCouponProduct` | 菜单权限 |

#### 营销 / 商户优惠券 / 领取记录

- 页面路由：`/marketing/coupon/user`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 使用记录 | `R` | `systemCouponIssue` | 菜单权限 |

#### 营销 / 平台优惠券

- 页面路由：`/marketing/platform_coupon`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 详情 | `R` | `systemCouponShow` | 路由 |
| 使用记录 | `R` | `systemCouponIssue` | 路由 |
| 复制表单 | `C` | `systemCouponCloneForm` | 路由 |

#### 营销 / 平台优惠券 / 优惠券列表

- 页面路由：`/marketing/platform_coupon/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 添加 | `C` | `systemCouponCreate` | 菜单权限 |
| 编辑 | `U` | `systemCouponUpdate` | 菜单权限 |
| 删除 | `D` | `systemCouponDelete` | 菜单权限 |
| 修改状态 | `U` | `systemCouponStatus` | 菜单权限 |
| 列表 | `R` | `systemCouponList` | 菜单权限 |
| 详情 | `RU` | `systemCouponShow` | 菜单权限 |
| 详情关联列表 | `RU` | `systemCouponShowLst` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 营销 / 平台优惠券 / 使用说明

- 页面路由：`/marketing/platform_coupon/instructions`
- CRUD：C=✓ R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 使用说明 | `C` | `systemAgreeSave` | 菜单权限 |

#### 营销 / 平台优惠券 / 发送记录

- 页面路由：`/marketing/platform_coupon/couponSend`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 发送记录 | `CR` | `systemCouponSendLst` | 菜单权限 |

#### 营销 / 平台优惠券 / 领取记录

- 页面路由：`/marketing/platform_coupon/couponRecord`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 使用记录 | `R` | `systemCouponIssue` | 菜单权限 |

#### 营销 / 报名活动

- 页面路由：`/marketing/application/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 添加 | `C` | `systemActivityFormCreate` | 菜单权限 |
| 列表 | `R` | `systemActivityFormLst` | 菜单权限 |
| 编辑 | `U` | `systemActivityFormUpdate` | 菜单权限 |
| 详情 | `R` | `systemActivityFormDetail` | 菜单权限 |
| 删除 | `D` | `systemActivityFormDelete` | 菜单权限 |
| 修改状态 | `U` | `systemActivityFormStatus` | 菜单权限 |
| 活动记录 | `R` | `systemFormActivUserLst` | 菜单权限 |
| 活动记录导出 | `R` | `systemFormActivUserExcel` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 营销 / 拼团

- 页面路由：`/marketing/combination`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 排序 | `U` | `systemStoreProductGroupSort` | 路由 |
| 设置标签 | `C` | `systemStoreProductGroupLabels` | 路由 |

#### 营销 / 拼团 / 拼团商品列表

- 页面路由：`/marketing/combination/combination_goods`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemStoreProductGroupLst` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreProductGroupShow` | 菜单权限 |
| 详情 | `R` | `systemStoreProductGroupDetail` | 菜单权限 |
| 编辑 | `U` | `systemStoreProductGroupProductUpdate` | 菜单权限 |
| 审核 | `U` | `systemStoreProductGroupStatus` | 菜单权限 |
| 编辑数据 | `RU` | `systemStoreProductGroupGet` | 菜单权限 |
| 排序 | `U` | `systemStoreProductGroupSort` | 菜单权限 |
| 设置标签 | `U` | `systemStoreProductGroupLabels` | 菜单权限 |
| 列表 | `R` | `systemStoreProductGroupLst` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreProductGroupShow` | 菜单权限 |
| 详情 | `R` | `systemStoreProductGroupDetail` | 菜单权限 |
| 编辑 | `U` | `systemStoreProductGroupProductUpdate` | 菜单权限 |
| 审核 | `U` | `systemStoreProductGroupStatus` | 菜单权限 |
| 编辑数据 | `RU` | `systemStoreProductGroupGet` | 菜单权限 |
| 排序 | `U` | `systemStoreProductGroupSort` | 菜单权限 |
| 设置标签 | `U` | `systemStoreProductGroupLabels` | 菜单权限 |
| 列表 | `R` | `systemStoreProductGroupLst` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreProductGroupShow` | 菜单权限 |
| 详情 | `R` | `systemStoreProductGroupDetail` | 菜单权限 |
| 编辑 | `U` | `systemStoreProductGroupProductUpdate` | 菜单权限 |
| 审核 | `U` | `systemStoreProductGroupStatus` | 菜单权限 |
| 编辑数据 | `RU` | `systemStoreProductGroupGet` | 菜单权限 |
| 排序 | `U` | `systemStoreProductGroupSort` | 菜单权限 |
| 设置标签 | `U` | `systemStoreProductGroupLabels` | 菜单权限 |

#### 营销 / 拼团 / 拼团活动列表

- 页面路由：`/marketing/combination/combination_list`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemStoreProductGroupBuyingLst` | 菜单权限 |
| 详情 | `R` | `systemStoreProductGroupBuyingDetail` | 菜单权限 |
| 列表 | `R` | `systemStoreProductGroupBuyingLst` | 菜单权限 |
| 详情 | `R` | `systemStoreProductGroupBuyingDetail` | 菜单权限 |
| 列表 | `R` | `systemStoreProductGroupBuyingLst` | 菜单权限 |
| 详情 | `R` | `systemStoreProductGroupBuyingDetail` | 菜单权限 |

#### 营销 / 拼团 / 拼团设置

- 页面路由：`/marketing/combination/combination_set`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 配置保存 | `CU` | `configOthersGroupBuyingUpdate` | 菜单权限 |

#### 营销 / 活动氛围图

- 页面路由：`/marketing/atmosphere/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 添加 | `C` | `systemActivityAtmosphereCreate` | 菜单权限 |
| 列表 | `R` | `systemActivityAtmosphereLst` | 菜单权限 |
| 编辑 | `U` | `systemActivityAtmosphereUpdate` | 菜单权限 |
| 详情 | `R` | `systemActivityAtmosphereDetail` | 菜单权限 |
| 删除 | `D` | `systemActivityAtmosphereDelete` | 菜单权限 |
| 修改状态 | `U` | `systemActivityAtmosphereStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 营销 / 活动边框图

- 页面路由：`/marketing/border/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 添加 | `C` | `systemActivityBorderCreate` | 菜单权限 |
| 列表 | `R` | `systemActivityBorderLst` | 菜单权限 |
| 编辑 | `U` | `systemActivityBorderUpdate` | 菜单权限 |
| 详情 | `R` | `systemActivityBorderDetail` | 菜单权限 |
| 删除 | `D` | `systemActivityBorderDelete` | 菜单权限 |
| 修改状态 | `U` | `systemActivityBorderStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 营销 / 直播

- 页面路由：`/marketing2`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 营销 / 直播 / 直播商品管理

- 页面路由：`/marketing/broadcast/list`
- CRUD：C=— R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemBroadcastGoodsLst` | 菜单权限 |
| 详情 | `R` | `systemBroadcastGoodsDetail` | 菜单权限 |
| 审核 | `U` | `systemBroadcastGoodsApply` | 菜单权限 |
| 修改状态 | `U` | `systemBroadcastGoodsChangeStatus` | 菜单权限 |
| 排序 | `U` | `systemBroadcastGoodsSort` | 菜单权限 |
| 删除 | `D` | `systemBroadcastGoodsDelete` | 菜单权限 |

#### 营销 / 直播 / 直播间管理

- 页面路由：`/marketing/studio/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemBroadcastRoomLst` | 菜单权限 |
| 详情 | `R` | `systemBroadcastRoomDetail` | 菜单权限 |
| 申请 | `C` | `systemBroadcastRoomApply` | 菜单权限 |
| 修改状态 | `U` | `systemBroadcastRoomChangeStatus` | 菜单权限 |
| 排序 | `U` | `systemBroadcastRoomSort` | 菜单权限 |
| 修改状态 | `U` | `systemBroadcastRoomChangeLiveStatus` | 菜单权限 |
| 删除 | `D` | `systemBroadcastRoomDelete` | 菜单权限 |
| 商品列表 | `R` | `systemBroadcastRoomGoods` | 菜单权限 |
| 客服开关 | `U` | `systemBroadcastRoomCloseKf` | 菜单权限 |
| 禁言开关 | `U` | `systemBroadcastRoomCloseComment` | 菜单权限 |
| 收录开关 | `U` | `systemBroadcastRoomClosesFeeds` | 菜单权限 |

#### 营销 / 秒杀

- 页面路由：`/marketing/seckill`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 排序 | `U` | `systemSeckillConfigStatus` | 路由 |
| 删除 | `D` | `systemSeckillConfigDelete` | 路由 |
| 列表 | `R` | `systemSeckillActiveGetActiveList` | 路由 |
| 详情 | `R` | `systemSeckillActiveGetActiveInfo` | 路由 |
| 列表 | `R` | `systemSeckillActiveGetActiveAll` | 路由 |
| 创建 | `C` | `systemSeckillActiveCreateActive` | 路由 |
| 编辑 | `U` | `systemSeckillActiveUpdateActive` | 路由 |
| 编辑状态 | `U` | `systemSeckillActiveUpdateActiveStatus` | 路由 |
| 删除 | `D` | `systemSeckillActiveDeleteActive` | 路由 |
| 活动统计数据面板 | `R` | `systemSeckillActiveChartPanel` | 路由 |
| 活动参与人统计列表 | `R` | `systemSeckillActiveChartPeople` | 路由 |
| 活动订单统计列表 | `R` | `systemSeckillActiveChartOrder` | 路由 |
| 活动商品统计列表 | `R` | `systemSeckillActiveChartProduct` | 路由 |
| 统计 | `R` | `systemStoreSeckillProductLstFilter` | 路由 |
| 列表 | `R` | `systemStoreSeckillProductPageLst` | 路由 |
| 列表 | `R` | `systemStoreSeckillProductLst` | 路由 |
| 详情 | `R` | `systemStoreSeckillProductDetail` | 路由 |
| 编辑 | `U` | `systemStoreSeckillProductUpdate` | 路由 |
| 审核 | `U` | `systemStoreSeckillProductSwitchStatus` | 路由 |
| 审核表单 | `U` | `systemStoreSeckillProductSwitchStatusForm` | 路由 |
| 显示/隐藏 | `U` | `systemStoreSeckillProductChangeUsed` | 路由 |
| 设置标签 | `C` | `systemStoreSeckillProductLabels` | 路由 |
| 加入回收站 | `C` | `systemStoreSeckillProductDelete` | 路由 |
| 删除 | `D` | `systemStoreSeckillProductDestory` | 路由 |

#### 营销 / 秒杀 / 秒杀活动

- 页面路由：`/marketing/seckill/store_seckill/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemSeckillActiveGetActiveList` | 菜单权限 |
| 详情 | `R` | `systemSeckillActiveGetActiveInfo` | 菜单权限 |
| 创建 | `C` | `systemSeckillActiveCreateActive` | 菜单权限 |
| 编辑 | `U` | `systemSeckillActiveUpdateActive` | 菜单权限 |
| 编辑状态 | `U` | `systemSeckillActiveUpdateActiveStatus` | 菜单权限 |
| 删除 | `D` | `systemSeckillActiveDeleteActive` | 菜单权限 |
| 活动统计数据面板 | `R` | `systemSeckillActiveChartPanel` | 菜单权限 |
| 活动参与人统计列表 | `R` | `systemSeckillActiveChartPeople` | 菜单权限 |
| 活动订单统计列表 | `R` | `systemSeckillActiveChartOrder` | 菜单权限 |
| 活动商品统计列表 | `R` | `systemSeckillActiveChartProduct` | 菜单权限 |
| 列表 | `R` | `systemSeckillActiveGetActiveList` | 菜单权限 |
| 详情 | `R` | `systemSeckillActiveGetActiveInfo` | 菜单权限 |
| 创建 | `C` | `systemSeckillActiveCreateActive` | 菜单权限 |
| 编辑 | `U` | `systemSeckillActiveUpdateActive` | 菜单权限 |
| 编辑状态 | `U` | `systemSeckillActiveUpdateActiveStatus` | 菜单权限 |
| 删除 | `D` | `systemSeckillActiveDeleteActive` | 菜单权限 |
| 活动统计数据面板 | `R` | `systemSeckillActiveChartPanel` | 菜单权限 |
| 活动参与人统计列表 | `R` | `systemSeckillActiveChartPeople` | 菜单权限 |
| 活动订单统计列表 | `R` | `systemSeckillActiveChartOrder` | 菜单权限 |
| 活动商品统计列表 | `R` | `systemSeckillActiveChartProduct` | 菜单权限 |
| 列表 | `R` | `systemSeckillActiveGetActiveList` | 菜单权限 |
| 详情 | `R` | `systemSeckillActiveGetActiveInfo` | 菜单权限 |
| 创建 | `C` | `systemSeckillActiveCreateActive` | 菜单权限 |
| 编辑 | `U` | `systemSeckillActiveUpdateActive` | 菜单权限 |
| 编辑状态 | `U` | `systemSeckillActiveUpdateActiveStatus` | 菜单权限 |
| 删除 | `D` | `systemSeckillActiveDeleteActive` | 菜单权限 |
| 活动统计数据面板 | `R` | `systemSeckillActiveChartPanel` | 菜单权限 |
| 活动参与人统计列表 | `R` | `systemSeckillActiveChartPeople` | 菜单权限 |
| 活动订单统计列表 | `R` | `systemSeckillActiveChartOrder` | 菜单权限 |
| 活动商品统计列表 | `R` | `systemSeckillActiveChartProduct` | 菜单权限 |

#### 营销 / 秒杀 / 秒杀管理

- 页面路由：`/marketing/seckill/list`
- CRUD：C=— R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 统计 | `R` | `systemStoreSeckillProductLstFilter` | 菜单权限 |
| 列表 | `R` | `systemStoreSeckillProductLst` | 菜单权限 |
| 编辑 | `U` | `systemStoreSeckillProductUpdate` | 菜单权限 |
| 审核 | `U` | `systemStoreSeckillProductSwitchStatus` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreSeckillProductChangeUsed` | 菜单权限 |
| 设置标签 | `U` | `systemStoreSeckillProductLabels` | 菜单权限 |
| 列表 | `R` | `systemStoreSeckillProductPageLst` | 菜单权限 |
| 加入回收站 | `D` | `systemStoreSeckillProductDelete` | 菜单权限 |
| 删除 | `D` | `systemStoreSeckillProductDestory` | 菜单权限 |
| 统计 | `R` | `systemStoreSeckillProductLstFilter` | 菜单权限 |
| 列表 | `R` | `systemStoreSeckillProductPageLst` | 菜单权限 |
| 列表 | `R` | `systemStoreSeckillProductLst` | 菜单权限 |
| 详情 | `R` | `systemStoreSeckillProductDetail` | 菜单权限 |
| 编辑 | `U` | `systemStoreSeckillProductUpdate` | 菜单权限 |
| 审核 | `U` | `systemStoreSeckillProductSwitchStatus` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreSeckillProductChangeUsed` | 菜单权限 |
| 设置标签 | `U` | `systemStoreSeckillProductLabels` | 菜单权限 |
| 加入回收站 | `D` | `systemStoreSeckillProductDelete` | 菜单权限 |
| 删除 | `D` | `systemStoreSeckillProductDestory` | 菜单权限 |
| 统计 | `R` | `systemStoreSeckillProductLstFilter` | 菜单权限 |
| 列表 | `R` | `systemStoreSeckillProductPageLst` | 菜单权限 |
| 列表 | `R` | `systemStoreSeckillProductLst` | 菜单权限 |
| 详情 | `R` | `systemStoreSeckillProductDetail` | 菜单权限 |
| 编辑 | `U` | `systemStoreSeckillProductUpdate` | 菜单权限 |
| 审核 | `U` | `systemStoreSeckillProductSwitchStatus` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreSeckillProductChangeUsed` | 菜单权限 |
| 设置标签 | `U` | `systemStoreSeckillProductLabels` | 菜单权限 |
| 加入回收站 | `D` | `systemStoreSeckillProductDelete` | 菜单权限 |
| 删除 | `D` | `systemStoreSeckillProductDestory` | 菜单权限 |

#### 营销 / 秒杀 / 秒杀配置

- 页面路由：`/marketing/seckill/seckillConfig`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemSeckillConfigLst` | 菜单权限 |
| 添加 | `C` | `systemSeckillConfigCreate` | 菜单权限 |
| 编辑 | `U` | `systemSeckillConfigUpdate` | 菜单权限 |
| 排序 | `U` | `systemSeckillConfigStatus` | 菜单权限 |
| 删除 | `D` | `systemSeckillConfigDelete` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 营销 / 积分

- 页面路由：`/marketing/integral`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 积分配置保存 | `C` | `systemUserIntegralConfigSave` | 路由 |
| 积分统计 | `R` | `systemUserIntegralTitle` | 路由 |

#### 营销 / 积分 / 商品分类

- 页面路由：`/marketing/integral/classify`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `pointsCateLst` | 菜单权限 |
| 详情 | `R` | `pointsCateDetail` | 菜单权限 |
| 添加 | `C` | `pointsCateCreate` | 菜单权限 |
| 编辑 | `U` | `pointsCateUpdate` | 菜单权限 |
| 修改状态 | `U` | `pointsCateStatus` | 菜单权限 |
| 列表 | `R` | `pointsProductLst` | 菜单权限 |
| 获取规格 | `R` | `pointsCateFormatAttr` | 菜单权限 |
| 编辑 | `RU` | `pointsProductDetail` | 菜单权限 |
| 添加 | `C` | `pointsProductCreate` | 菜单权限 |
| 编辑 | `U` | `pointsProductUpdate` | 菜单权限 |
| 修改状态 | `U` | `pointsProductStatus` | 菜单权限 |
| 预览 | `R` | `pointsProductPreview` | 菜单权限 |

#### 营销 / 积分 / 商品列表

- 页面路由：`/marketing/integral/proList`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `pointsProductLst` | 菜单权限 |
| 获取规格 | `R` | `pointsCateFormatAttr` | 菜单权限 |
| 编辑 | `RU` | `pointsProductDetail` | 菜单权限 |
| 添加 | `C` | `pointsProductCreate` | 菜单权限 |
| 编辑 | `U` | `pointsProductUpdate` | 菜单权限 |
| 修改状态 | `U` | `pointsProductStatus` | 菜单权限 |
| 预览 | `R` | `pointsProductPreview` | 菜单权限 |

#### 营销 / 积分 / 积分日志

- 页面路由：`/marketing/integral/log`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 积分统计 | `R` | `systemUserIntegralTitle` | 菜单权限 |
| 积分日志 | `R` | `systemUserIntegralLst` | 菜单权限 |
| 积分导出 | `R` | `systemUserIntegralExcel` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `systemStoreExcelDownload` | 菜单权限 |

#### 营销 / 积分 / 积分订单

- 页面路由：`/marketing/integral/orderList`
- CRUD：C=— R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `pointsOrderLst` | 菜单权限 |
| 编辑 | `RU` | `pointsOrderDetail` | 菜单权限 |
| 发货 | `U` | `pointsOrderDelivery` | 菜单权限 |
| 批量发货 | `U` | `pointsOrderBatchDelivery` | 菜单权限 |
| 快递查询 | `R` | `pointsOrderExpress` | 菜单权限 |
| 导出 | `R` | `pointsOrderExcel` | 菜单权限 |
| 备注 | `U` | `pointsOrderMark` | 菜单权限 |
| 记录 | `RU` | `pointsOrderStatus` | 菜单权限 |
| 删除 | `D` | `pointsOrderDelete` | 菜单权限 |

#### 营销 / 积分 / 积分配置

- 页面路由：`/marketing/integral/config`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 积分配置保存 | `CU` | `systemUserIntegralConfigSave` | 菜单权限 |

#### 营销 / 预售

- 页面路由：`/marketing/presell`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 审核 | `U` | `systemStoreProductPresellSwitchStatus` | 路由 |
| 设置标签 | `C` | `systemStoreProductPresellLabels` | 路由 |
| 列表 | `R` | `systemStoreProductAssistLst` | 路由 |
| 显示/隐藏 | `U` | `systemStoreProductAssistShow` | 路由 |
| 详情 | `R` | `systemStoreProductAssistDetail` | 路由 |
| 编辑 | `U` | `systemStoreProductAssistProductUpdate` | 路由 |
| 审核 | `U` | `systemStoreProductAssistStatus` | 路由 |

#### 营销 / 预售 / 预售协议

- 页面路由：`/marketing/presell/agreement`
- CRUD：C=✓ R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 预售协议 | `C` | `systemAgreeSave` | 菜单权限 |

#### 营销 / 预售 / 预售商品

- 页面路由：`/marketing/presell/list`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemStoreProductPresellLst` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreProductPresellShow` | 菜单权限 |
| 详情 | `R` | `systemStoreProductPresellDetail` | 菜单权限 |
| 编辑数据 | `RU` | `systemStoreProductPresellGet` | 菜单权限 |
| 编辑 | `U` | `systemStoreProductPresellUpdate` | 菜单权限 |
| 审核 | `U` | `systemStoreProductPresellSwitchStatus` | 菜单权限 |
| 设置标签 | `U` | `systemStoreProductPresellLabels` | 菜单权限 |
| 列表 | `R` | `systemStoreProductPresellLst` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreProductPresellShow` | 菜单权限 |
| 详情 | `R` | `systemStoreProductPresellDetail` | 菜单权限 |
| 编辑数据 | `RU` | `systemStoreProductPresellGet` | 菜单权限 |
| 编辑 | `U` | `systemStoreProductPresellUpdate` | 菜单权限 |
| 审核 | `U` | `systemStoreProductPresellSwitchStatus` | 菜单权限 |
| 设置标签 | `U` | `systemStoreProductPresellLabels` | 菜单权限 |
| 列表 | `R` | `systemStoreProductPresellLst` | 菜单权限 |
| 显示/隐藏 | `U` | `systemStoreProductPresellShow` | 菜单权限 |
| 详情 | `R` | `systemStoreProductPresellDetail` | 菜单权限 |
| 编辑数据 | `RU` | `systemStoreProductPresellGet` | 菜单权限 |
| 编辑 | `U` | `systemStoreProductPresellUpdate` | 菜单权限 |
| 审核 | `U` | `systemStoreProductPresellSwitchStatus` | 菜单权限 |
| 设置标签 | `U` | `systemStoreProductPresellLabels` | 菜单权限 |

### 装修

#### 装修

- 页面路由：`/theme`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 装修 / 个人中心

- 页面路由：`/setting/diy/personal`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 装修 / 主题风格

- 页面路由：`/setting/theme_style`
- CRUD：C=✓ R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 一键换色保存 | `C` | `systemSetChangeColor` | 菜单权限 |

#### 装修 / 商品分类

- 页面路由：`/setting/product_category`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商品分类信息 | `R` | `systemDiyProductCategoryInfo` | 菜单权限 |
| 保存商品分类 | `C` | `systemDiyProductCategorySave` | 菜单权限 |

#### 装修 / 商品详情

- 页面路由：`/setting/diy/product_detail`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 装修 / 店铺模板

- 页面路由：`/setting/merchant/diyList`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表  | `R` | `systemMerDiyLst` | 菜单权限 |
| 详情  | `R` | `systemMerDiyDetail` | 菜单权限 |
| 添加/编辑 | `CU` | `systemMerDiyCreate` | 菜单权限 |
| 设置默认 | `U` | `systemMerDiySetDefault` | 菜单权限 |
| 重置 | `O` | `systemMerDiyRecovery` | 菜单权限 |
| 删除 | `D` | `systemMerDiyDelete` | 菜单权限 |
| 复制 | `C` | `systemMerDiyCopy` | 菜单权限 |
| 保存适用范围 | `CR` | `systemMerDiyGetScope` | 菜单权限 |
| 保存适用范围 | `C` | `systemMerDiySetScope` | 菜单权限 |

#### 装修 / 店铺街

- 页面路由：`/setting/diy/store`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 装修 / 微页面

- 页面路由：`/setting/micro/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表  | `R` | `systemDiyMicroLst` | 菜单权限 |
| 详情  | `R` | `systemDiyMicroDetail` | 菜单权限 |
| 添加/编辑 | `CU` | `systemDiyMicroCreate` | 菜单权限 |
| 重置 | `O` | `systemDiyMicroRecovery` | 菜单权限 |
| 删除 | `D` | `systemDiyMicroDelete` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 装修 / 悬浮菜单

- 页面路由：`/setting/fab`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 悬浮按钮信息 | `R` | `systemDiyFabInfo` | 菜单权限 |
| 保存悬浮按钮 | `C` | `systemDiyFabSave` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 装修 / 系统表单

- 页面路由：`/systemForm/form_list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 添加 | `C` | `systemFormCreate` | 菜单权限 |
| 编辑 | `U` | `systemFormUpdate` | 菜单权限 |
| 编辑状态 | `U` | `systemFormStatusSwitch` | 菜单权限 |
| 删除 | `D` | `systemFormDelete` | 菜单权限 |
| 详情 | `R` | `systemFormDetail` | 菜单权限 |
| 列表 | `R` | `systemFormLst` | 菜单权限 |
| 导出 | `R` | `systemFormExcel` | 菜单权限 |
| 表单提交记录 | `CR` | `systemFormUserLst` | 菜单权限 |

#### 装修 / 素材管理

- 页面路由：`/config/picture`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 素材分类列表 | `R` | `systemAttachmentCategoryGetFormatList` | 菜单权限 |
| 素材分类添加 | `C` | `systemAttachmentCategoryCreate` | 菜单权限 |
| 素材编辑 | `U` | `systemAttachmentCategoryUpdate` | 菜单权限 |
| 素材删除 | `D` | `systemAttachmentCategoryDelete` | 菜单权限 |
| 素材列表 | `R` | `systemAttachmentLst` | 菜单权限 |
| 素材删除 | `D` | `systemAttachmentDelete` | 菜单权限 |
| 批量移动 | `O` | `systemAttachmentBatchChangeCategory` | 菜单权限 |
| 素材编辑 | `U` | `systemAttachmentUpdate` | 菜单权限 |
| 上传二维码 | `C` | `systemAttachmentScanQrcode` | 菜单权限 |
| 扫码上传图片 | `C` | `systemAttachmentScanImage` | 菜单权限 |
| 扫码上传保存 | `C` | `systemAttachmentScanImageSave` | 菜单权限 |
| 在线图片 | `O` | `systemAttachmentOnline` | 菜单权限 |

#### 装修 / 页面装修

- 页面路由：`/setting/diy/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表  | `R` | `systemDiyLst` | 菜单权限 |
| 添加/编辑 | `CU` | `systemDiyCreate` | 菜单权限 |
| 使用模板 | `U` | `systemDiyStatus` | 菜单权限 |
| 设置默认 | `U` | `systemDiySetDefault` | 菜单权限 |
| 重置 | `O` | `systemDiyRecovery` | 菜单权限 |
| 删除 | `D` | `systemDiyDelete` | 菜单权限 |
| 商品列表 | `R` | `systemDiyProductLst` | 菜单权限 |
| 复制 | `C` | `systemDiyCopy` | 菜单权限 |
| 个人中心装修 | `R` | `systemVisualUserInfo` | 菜单权限 |
| 店铺街装修 | `O` | `systemVisualStoreStreet` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |
| 详情  | `R` | `systemDiyDetail` | 菜单权限 |
| 商品详情  | `R` | `systemDiyGetProductDetail` | 菜单权限 |
| 商品详情保存  | `CR` | `systemDiySaveProductDetail` | 菜单权限 |

#### 装修 / 页面配置

- 页面路由：`/setting/system_visualization_data`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 可视化列表 | `R` | `systemVisualStoreGetThemeKey` | 菜单权限 |
| 可视化详情 | `R` | `systemVisualStoreGetTheme` | 菜单权限 |
| 可视化保存 | `C` | `systemVisualSetTheme` | 菜单权限 |

#### 装修 / 页面链接

- 页面路由：`/setting/page`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 装修 / 页面链接 / 商户页面分类

- 页面路由：`/setting/diy/merchant/category/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表  | `R` | `systemDiyPageMerCategroyLst` | 菜单权限 |
| 添加 | `C` | `systemDiyPageMerCategroyCreate` | 菜单权限 |
| 编辑 | `U` | `systemDiyPageMerCategroyUpdate` | 菜单权限 |
| 编辑状态 | `U` | `systemDiyPageMerCategroyStatus` | 菜单权限 |
| 删除 | `D` | `systemDiyPageMerCategroyDelete` | 菜单权限 |

#### 装修 / 页面链接 / 商户页面链接

- 页面路由：`/setting/diy/merLink/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemDiyPageLinkMerLst` | 菜单权限 |
| 添加 | `C` | `systemDiyPageLinkMerCreate` | 菜单权限 |
| 编辑 | `U` | `systemDiyPageLinkMerUpdate` | 菜单权限 |
| 删除 | `D` | `systemDiyPageLinkMerDelete` | 菜单权限 |
| 修改状态 | `U` | `systemDiyPageLinkMerStatus` | 菜单权限 |

#### 装修 / 页面链接 / 平台页面分类

- 页面路由：`/setting/diy/plantform/category/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表  | `R` | `systemDiyPageCategroyLst` | 菜单权限 |
| 添加 | `C` | `systemDiyPageCategroyCreate` | 菜单权限 |
| 编辑 | `U` | `systemDiyPageCategroyUpdate` | 菜单权限 |
| 编辑状态 | `U` | `systemDiyPageCategroyStatus` | 菜单权限 |
| 删除 | `D` | `systemDiyPageCategroyDelete` | 菜单权限 |

#### 装修 / 页面链接 / 平台页面链接

- 页面路由：`/setting/diy/links/list`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemDiyPageLinkLst` | 菜单权限 |
| 添加 | `C` | `systemDiyPageLinkCreate` | 菜单权限 |
| 编辑 | `U` | `systemDiyPageLinkUpdate` | 菜单权限 |
| 删除 | `D` | `systemDiyPageLinkDelete` | 菜单权限 |
| 修改状态 | `U` | `systemDiyPageLinkStatus` | 菜单权限 |

### 订单

#### 订单

- 页面路由：`/order`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 记录 | `R` | `systemOrderStatus` | 路由 |
| 关联订单 | `R` | `systemOrderChildrenList` | 路由 |
| 核销 | `R` | `systemOrderTakeStat` | 路由 |

#### 订单 / 核销记录

- 页面路由：`/order/cancellation`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 核销 | `U` | `systemOrderTakeStat` | 菜单权限 |
| 核销订单 | `RU` | `systemTakeOrderLst` | 菜单权限 |
| 头部统计 | `R` | `systemTakeOrderTitle` | 菜单权限 |
| 核销 | `U` | `systemOrderTakeStat` | 菜单权限 |
| 核销订单 | `RU` | `systemTakeOrderLst` | 菜单权限 |
| 头部统计 | `R` | `systemTakeOrderTitle` | 菜单权限 |
| 核销 | `U` | `systemOrderTakeStat` | 菜单权限 |
| 核销订单 | `RU` | `systemTakeOrderLst` | 菜单权限 |
| 头部统计 | `R` | `systemTakeOrderTitle` | 菜单权限 |

#### 订单 / 订单列表

- 页面路由：`/order/list`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemOrderLst` | 菜单权限 |
| 金额统计 | `R` | `systemOrderStat` | 菜单权限 |
| 快递查询 | `R` | `systemOrderExpress` | 菜单权限 |
| 头部统计 | `R` | `systemOrderTitle` | 菜单权限 |
| 详情 | `R` | `systemOrderDetail` | 菜单权限 |
| 导出 | `R` | `systemOrderExcel` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelDownload` | 菜单权限 |
| 记录 | `RU` | `systemOrderStatus` | 菜单权限 |
| 关联订单 | `O` | `systemOrderChildrenList` | 菜单权限 |
| 列表 | `R` | `systemOrderLst` | 菜单权限 |
| 金额统计 | `R` | `systemOrderStat` | 菜单权限 |
| 快递查询 | `R` | `systemOrderExpress` | 菜单权限 |
| 头部统计 | `R` | `systemOrderTitle` | 菜单权限 |
| 详情 | `R` | `systemOrderDetail` | 菜单权限 |
| 导出 | `R` | `systemOrderExcel` | 菜单权限 |
| 记录 | `RU` | `systemOrderStatus` | 菜单权限 |
| 关联订单 | `O` | `systemOrderChildrenList` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelDownload` | 菜单权限 |
| 列表 | `R` | `systemOrderLst` | 菜单权限 |
| 金额统计 | `R` | `systemOrderStat` | 菜单权限 |
| 快递查询 | `R` | `systemOrderExpress` | 菜单权限 |
| 头部统计 | `R` | `systemOrderTitle` | 菜单权限 |
| 详情 | `R` | `systemOrderDetail` | 菜单权限 |
| 导出 | `R` | `systemOrderExcel` | 菜单权限 |
| 记录 | `RU` | `systemOrderStatus` | 菜单权限 |
| 关联订单 | `O` | `systemOrderChildrenList` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelDownload` | 菜单权限 |

#### 订单 / 退款订单

- 页面路由：`/order/refund`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemRefundOrderLst` | 菜单权限 |
| 导出 | `R` | `systemRefundOrderExcel` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `systemStoreExcelDownload` | 菜单权限 |
| 详情 | `R` | `systemRefundOrderDetail` | 菜单权限 |
| 日志 | `R` | `systemRefundOrderLog` | 菜单权限 |
| 审核 | `U` | `systemRefundOrderApprove` | 菜单权限 |
| 列表 | `R` | `systemRefundOrderLst` | 菜单权限 |
| 详情 | `R` | `systemRefundOrderDetail` | 菜单权限 |
| 日志 | `R` | `systemRefundOrderLog` | 菜单权限 |
| 审核 | `U` | `systemRefundOrderApprove` | 菜单权限 |
| 导出 | `R` | `systemRefundOrderExcel` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `systemStoreExcelDownload` | 菜单权限 |
| 列表 | `R` | `systemRefundOrderLst` | 菜单权限 |
| 详情 | `R` | `systemRefundOrderDetail` | 菜单权限 |
| 日志 | `R` | `systemRefundOrderLog` | 菜单权限 |
| 审核 | `U` | `systemRefundOrderApprove` | 菜单权限 |
| 导出 | `R` | `systemRefundOrderExcel` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `systemStoreExcelDownload` | 菜单权限 |

### 设置

#### 设置

- 页面路由：`/settings`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 设置 / 商城设置

- 页面路由：`/shop`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 设置 / 商城设置 / 协议规则

- 页面路由：`/setting/agreements`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 协议列表 | `R` | `systemAgreeKeyLsy` | 菜单权限 |
| 商户入住申请协议 | `C` | `systemAgreeSave` | 菜单权限 |

#### 设置 / 商城设置 / 商城设置

- 页面路由：`/systemForm/Basics/shop_tabs`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

#### 设置 / 商城设置 / 支付设置

- 页面路由：`/systemForm/Basics/pay_tabs`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

#### 设置 / 商城设置 / 热门搜索

- 页面路由：`/group/config/67`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 详情 | `R` | `groupDetail` | 菜单权限 |
| 列表 | `R` | `groupDataLst` | 菜单权限 |
| 添加 | `C` | `groupDataCreate` | 菜单权限 |
| 编辑 | `U` | `groupDataUpdate` | 菜单权限 |
| 删除 | `D` | `groupDataDelete` | 菜单权限 |
| 修改状态 | `U` | `groupDataChangeStatus` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 设置 / 商城设置 / 配送配置

- 页面路由：`/delivery_config`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 设置 / 商城设置 / 配送配置 / 城市数据

- 页面路由：`/freight/city/list`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 设置 / 商城设置 / 配送配置 / 物流公司

- 页面路由：`/freight/express`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemExpressLst` | 菜单权限 |
| 修改状态 | `U` | `systemExpressSwitchStatus` | 菜单权限 |
| 编辑 | `U` | `systemExpressUpdate` | 菜单权限 |
| 删除 | `D` | `systemExpressDelete` | 菜单权限 |
| 同步 | `U` | `systemExpressSync` | 菜单权限 |
| 列表 | `R` | `systemServeExportLst` | 菜单权限 |
| 列表 | `R` | `systemCityAreaLst` | 菜单权限 |
| 编辑 | `CU` | `systemCityAreaCreate` | 菜单权限 |
| 编辑 | `U` | `systemCityAreaUpdate` | 菜单权限 |
| 删除 | `D` | `systemCityAreaDelete` | 菜单权限 |

#### 设置 / 商城设置 / 配送配置 / 第三方送

- 页面路由：`/delivery`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 门店列表 | `R` | `systemDeliveryStationlst` | 路由 |
| 门店详情 | `R` | `systemDeliveryStationDetail` | 路由 |
| 门店筛选 | `R` | `systemStoreDeliveryOptions` | 路由 |

#### 设置 / 商城设置 / 配送配置 / 第三方送 / 充值记录

- 页面路由：`/delivery/recharge_record`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 充值记录 | `CR` | `systemDeliveryStationPaayyLst` | 菜单权限 |
| 统计 | `R` | `systemDeliveryOrderTitle` | 菜单权限 |
| 余额 | `R` | `systemDeliveryStationGetBalance` | 菜单权限 |

#### 设置 / 商城设置 / 配送配置 / 第三方送 / 配送记录

- 页面路由：`/delivery/usage_record`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 配送记录 | `R` | `systemDeliveryOrderLst` | 菜单权限 |
| 配送详情 | `R` | `systemDeliveryOrderDetail` | 菜单权限 |

#### 设置 / 商城设置 / 配送配置 / 第三方送 / 配送配置

- 页面路由：`/systemForm/delivery`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置 | `CU` | `systemDeliveryConfigSave` | 菜单权限 |

#### 设置 / 应用配置

- 页面路由：`/app_config`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 设置 / 应用配置 / APP升级配置

- 页面路由：`/systemForm/Basics/app_version`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

#### 设置 / 应用配置 / APP配置

- 页面路由：`/systemForm/Basics/wechat_open_app`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

#### 设置 / 应用配置 / 上传校验文件

- 页面路由：`/app/wechat/file`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 上传文件 | `C` | `configUpload` | 菜单权限 |
| 上传原名文件 | `C` | `configUploadName` | 菜单权限 |
| 微信校验文件上传 | `C` | `configWechatUploadSet` | 菜单权限 |
| 小程序配置 | `U` | `configRoutineConfig` | 菜单权限 |

#### 设置 / 应用配置 / 公众号配置

- 页面路由：`/systemForm/Basics/wechat`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

#### 设置 / 应用配置 / 小程序配置

- 页面路由：`/systemForm/Basics/smallapp`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

#### 设置 / 权限管理

- 页面路由：`/setting`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 一键换色 | `R` | `systemGetChangeColor` | 路由 |
| 一键换色保存 | `C` | `systemSetChangeColor` | 路由 |
| 列表  | `R` | `systemDiyPageCategroyLst` | 路由 |
| 添加表单 | `C` | `systemDiyPageCategroyCreateForm` | 路由 |
| 添加 | `C` | `systemDiyPageCategroyCreate` | 路由 |
| 编辑表单 | `U` | `systemDiyPageCategroyUpdateForm` | 路由 |
| 编辑 | `U` | `systemDiyPageCategroyUpdate` | 路由 |

#### 设置 / 权限管理 / 管理员管理

- 页面路由：`/setting/systemAdmin`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 管理员列表 | `R` | `systemAdminLst` | 菜单权限 |
| 管理员修改状态 | `U` | `systemAdminStatus` | 菜单权限 |
| 管理员添加 | `C` | `systemAdminCreate` | 菜单权限 |
| 管理员编辑 | `U` | `systemAdminUpdate` | 菜单权限 |
| 管理员修改密码 | `U` | `systemAdminPassword` | 菜单权限 |
| 管理员删除 | `D` | `systemAdminDelete` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |
| 管理员列表 | `R` | `systemAdminLst` | 菜单权限 |
| 管理员修改状态 | `U` | `systemAdminStatus` | 菜单权限 |
| 管理员添加 | `C` | `systemAdminCreate` | 菜单权限 |
| 管理员编辑 | `U` | `systemAdminUpdate` | 菜单权限 |
| 管理员修改密码 | `U` | `systemAdminPassword` | 菜单权限 |
| 管理员删除 | `D` | `systemAdminDelete` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |
| 管理员列表 | `R` | `systemAdminLst` | 菜单权限 |
| 管理员修改状态 | `U` | `systemAdminStatus` | 菜单权限 |
| 管理员添加 | `C` | `systemAdminCreate` | 菜单权限 |
| 管理员编辑 | `U` | `systemAdminUpdate` | 菜单权限 |
| 管理员修改密码 | `U` | `systemAdminPassword` | 菜单权限 |
| 管理员删除 | `D` | `systemAdminDelete` | 菜单权限 |
| 上传图片 | `C` | `uploadImage` | 菜单权限 |
| 图片列表 | `R` | `systemAttachmentLst` | 菜单权限 |

#### 设置 / 权限管理 / 菜单管理

- 页面路由：`/setting/menu`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 平台菜单/权限列表 | `R` | `systemMenuGetLst` | 菜单权限 |
| 平台菜单/权限添加 | `C` | `systemMenuCreate` | 菜单权限 |
| 平台菜单/权限编辑 | `U` | `systemMenuUpdate` | 菜单权限 |
| 平台菜单/权限删除 | `D` | `systemMenuDelete` | 菜单权限 |
| 搜索获取菜单 | `R` | `getMenusList` | 菜单权限 |

#### 设置 / 权限管理 / 角色权限

- 页面路由：`/setting/systemRole`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 身份列表 | `R` | `systemRoleGetList` | 菜单权限 |
| 身份添加 | `C` | `systemRoleCreate` | 菜单权限 |
| 身份编辑 | `U` | `systemRoleUpdate` | 菜单权限 |
| 身份修改状态 | `U` | `systemRoleStatus` | 菜单权限 |
| 身份删除 | `D` | `systemRoleDelete` | 菜单权限 |
| 身份列表 | `R` | `systemOrganizationRoleGetList` | 菜单权限 |
| 身份添加 | `C` | `systemOrganizationRoleCreate` | 菜单权限 |
| 身份编辑 | `U` | `systemOrganizationRoleUpdate` | 菜单权限 |
| 身份修改状态 | `U` | `systemOrganizationRoleStatus` | 菜单权限 |
| 身份删除 | `D` | `systemOrganizationRoleDelete` | 菜单权限 |
| 身份选项 | `O` | `systemOrganizationRoleOptions` | 菜单权限 |
| 身份列表 | `R` | `systemRoleGetList` | 菜单权限 |
| 身份添加 | `C` | `systemRoleCreate` | 菜单权限 |
| 身份编辑 | `U` | `systemRoleUpdate` | 菜单权限 |
| 身份修改状态 | `U` | `systemRoleStatus` | 菜单权限 |
| 身份删除 | `D` | `systemRoleDelete` | 菜单权限 |
| 身份列表 | `R` | `systemOrganizationRoleGetList` | 菜单权限 |
| 身份添加 | `C` | `systemOrganizationRoleCreate` | 菜单权限 |
| 身份编辑 | `U` | `systemOrganizationRoleUpdate` | 菜单权限 |
| 身份修改状态 | `U` | `systemOrganizationRoleStatus` | 菜单权限 |
| 身份删除 | `D` | `systemOrganizationRoleDelete` | 菜单权限 |
| 身份选项 | `O` | `systemOrganizationRoleOptions` | 菜单权限 |
| 身份列表 | `R` | `systemRoleGetList` | 菜单权限 |
| 身份添加 | `C` | `systemRoleCreate` | 菜单权限 |
| 身份编辑 | `U` | `systemRoleUpdate` | 菜单权限 |
| 身份修改状态 | `U` | `systemRoleStatus` | 菜单权限 |
| 身份删除 | `D` | `systemRoleDelete` | 菜单权限 |

#### 设置 / 消息管理

- 页面路由：`/notice`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 设置 / 消息管理 / 公告管理

- 页面路由：`/station/notice`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 系统公告列表 | `R` | `systemNoticeList` | 菜单权限 |
| 系统公告发布 | `C` | `systemNoticeCreate` | 菜单权限 |
| 系统公告编辑 | `U` | `systemNoticeUpdate` | 菜单权限 |
| 系统公告修改状态 | `U` | `systemNoticeSwitchStatus` | 菜单权限 |
| 系统公告详情 | `R` | `systemNoticeDetail` | 菜单权限 |
| 系统公告删除 | `D` | `systemNoticeDelete` | 菜单权限 |

#### 设置 / 消息管理 / 消息管理

- 页面路由：`/setting/notification/index`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 消息配置列表 | `RU` | `systemNoticeConfigLst` | 菜单权限 |
| 消息配置添加 | `CU` | `systemNoticeConfigCreate` | 菜单权限 |
| 消息配置编辑 | `U` | `systemNoticeConfigUpdate` | 菜单权限 |
| 消息配置详情 | `RU` | `systemNoticeConfigDetail` | 菜单权限 |
| 消息配置删除 | `UD` | `systemNoticeConfigDelete` | 菜单权限 |
| 消息配置修改状态 | `U` | `systemNoticeConfigStatus` | 菜单权限 |
| 消息配置修改模板ID | `RU` | `systemNoticeConfigGetChangeTempId` | 菜单权限 |
| 消息配置修改模板ID | `U` | `systemNoticeConfigSetChangeTempId` | 菜单权限 |

#### 设置 / 系统设置

- 页面路由：`/sys`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 配置表单 | `R` | `systemDeliveryConfigForm` | 路由 |
| 编辑配置 | `U` | `systemDeliveryConfigSave` | 路由 |

#### 设置 / 系统设置 / 一号通

- 页面路由：`/serve`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 设置 / 系统设置 / 一号通 / 商户结余

- 页面路由：`/service/balance_record`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商户结余 | `R` | `systemServeMerLst` | 菜单权限 |

#### 设置 / 系统设置 / 一号通 / 服务配置

- 页面路由：`/service/settings`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemServeMerMealLst` | 菜单权限 |
| 详情 | `R` | `systemServeMealDetail` | 菜单权限 |
| 添加 | `C` | `systemServeMealCreate` | 菜单权限 |
| 编辑 | `U` | `systemServeMealUpdate` | 菜单权限 |
| 删除 | `D` | `systemServeMealDelete` | 菜单权限 |
| 修改状态 | `U` | `systemServeMealStatus` | 菜单权限 |

#### 设置 / 系统设置 / 一号通 / 登陆入口

- 页面路由：`/setting/sms/sms_config/index`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 使用记录 | `CR` | `systemStoreProductCopyLst` | 菜单权限 |
| 短信发送记录 | `CR` | `smsRecord` | 菜单权限 |
| 退出登录 | `O` | `smsLogout` | 菜单权限 |
| 获取验证码 | `R` | `systemServeCaptcha` | 菜单权限 |
| 验证码校验 | `O` | `systemServeCaptchaCheck` | 菜单权限 |
| 注册 | `O` | `systemServeRegister` | 菜单权限 |
| 登录 | `O` | `systemServeLogin` | 菜单权限 |
| 修改密码 | `U` | `systemServeChangePassword` | 菜单权限 |
| 修改手机号 | `U` | `systemServeChangePhone` | 菜单权限 |
| 检测登录状态 | `U` | `systemServeIsLogin` | 菜单权限 |
| 使用记录 | `R` | `systemServeRecordLst` | 菜单权限 |
| 套餐列表 | `R` | `systemServeMealLst` | 菜单权限 |
| 购买套餐 | `O` | `systemServePayMeal` | 菜单权限 |
| 开通服务 | `C` | `systemServeOpenServe` | 菜单权限 |
| 修改签名 | `U` | `systemServeChangeSign` | 菜单权限 |
| 模板 | `R` | `systemServeExportTemps` | 菜单权限 |
| 使用记录 | `R` | `systemServeExportDumpLst` | 菜单权限 |

#### 设置 / 系统设置 / 一号通 / 短信设置

- 页面路由：`/sms`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 短信模板 | `R` | `systemServeSmsTemps` | 路由 |
| 申请模板 | `C` | `systemServeSmsApply` | 路由 |

#### 设置 / 系统设置 / 一号通 / 短信设置 / 申请记录

- 页面路由：`/sms/applyList`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 申请记录 | `CR` | `systemServeSmsApplyRecord` | 菜单权限 |

#### 设置 / 系统设置 / 一号通 / 短信设置 / 短信模板

- 页面路由：`/sms/template`
- CRUD：C=✓ R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 短信模板 | `O` | `systemServeSmsTemps` | 菜单权限 |
| 申请模板 | `C` | `systemServeSmsApply` | 菜单权限 |

#### 设置 / 系统设置 / 一号通 / 短信设置 / 短信账户

- 页面路由：`/sms/user`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 设置 / 系统设置 / 一号通 / 短信设置 / 短信配置

- 页面路由：`/systemForm/Basics/message`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

#### 设置 / 系统设置 / 一号通 / 购买记录

- 页面路由：`/service/purchase`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 购买记录 | `R` | `systemServePayLst` | 菜单权限 |
| 商户购买记录 | `R` | `systemServeMerPayLst` | 菜单权限 |

#### 设置 / 系统设置 / 存储管理

- 页面路由：`/setting/storage`
- CRUD：C=✓ R=✓ U=✓ D=✓

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 配置信息 | `RU` | `systemStorageGetConfig` | 菜单权限 |
| 提交配置 | `CU` | `systemStorageSetConfig` | 菜单权限 |
| 保存云存储配置 | `CU` | `systemStorageUpdate` | 菜单权限 |
| 同步存储空间 | `U` | `systemStorageSync` | 菜单权限 |
| 存储空间列表 | `R` | `systemStorageLstRegion` | 菜单权限 |
| 添加存储空间 | `C` | `systemStorageCreateRegion` | 菜单权限 |
| 删除存储空间 | `D` | `systemStorageDeleteRegion` | 菜单权限 |
| 使用存储空间 | `U` | `systemStorageRegionSwtichStatus` | 菜单权限 |
| 修改存储空间名称 | `U` | `systemStorageUpdateDomain` | 菜单权限 |

#### 设置 / 系统设置 / 接口配置

- 页面路由：`/systemForm/Basics/extend_tabs`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

#### 设置 / 系统设置 / 系统设置

- 页面路由：`/systemForm/Basics/system_tabs`
- CRUD：C=✓ R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 编辑配置信息 | `CU` | `configSave` | 菜单权限 |

### 财务

#### 财务

- 页面路由：`/accounts`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 配置信息 | `R` | `systemOrderProfitsharingGetConfig` | 路由 |
| 配置保存 | `C` | `systemOrderProfitsharingSetConfig` | 路由 |
| 申请列表 | `R` | `systemUserExtractLst` | 路由 |
| 审核表单 | `U` | `systemUserExtractSwitchStatusForm` | 路由 |
| 审核 | `U` | `systemUserExtractSwitchStatus` | 路由 |

#### 财务 / 发票管理

- 页面路由：`/accounts/accounts`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 财务 / 发票管理 / 发票列表

- 页面路由：`/accounts/receipt`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemOrderReceiptList` | 菜单权限 |
| 详情 | `R` | `systemOrderReceiptDetail` | 菜单权限 |

#### 财务 / 发票管理 / 发票说明

- 页面路由：`/accounts/invoiceDesc`
- CRUD：C=✓ R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 发票说明  | `C` | `systemAgreeSave` | 菜单权限 |

#### 财务 / 商圈代理

- 页面路由：`/accounts/zoneAgent`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 结算方式get | `R` | `systemCircleAgentGetSettlementMethod` | 路由 |
| 结算方式post | `C` | `systemCircleAgentSetSettlementMethod` | 路由 |
| 结算记录列表 | `R` | `systemCircleCheckoutList` | 路由 |
| 结算记录详情 | `R` | `systemCircleCheckoutDetail` | 路由 |
| 平台结算审核 | `U` | `systemCircleCheckoutAudit` | 路由 |

#### 财务 / 商圈代理 / 提成流水

- 页面路由：`/accounts/zoneAgent/commissionRecord`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商圈提成流水列表 | `R` | `systemCircleFinancialRecordList` | 菜单权限 |
| 商圈提成流水列表 | `R` | `systemCircleFinancialRecordList` | 菜单权限 |

#### 财务 / 商圈代理 / 申请结算

- 页面路由：`/accounts/zoneAgent/settlementApply`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商圈申请结算获取余额 | `CR` | `systemCircleCheckoutCreate` | 菜单权限 |
| 商圈申请结算提交 | `C` | `systemCircleCheckoutSave` | 菜单权限 |
| 商圈撤销结算 | `O` | `systemCircleCheckoutRevoke` | 菜单权限 |
| 商圈备注 | `U` | `systemCircleCheckoutRemark` | 菜单权限 |

#### 财务 / 商圈代理 / 结算审核

- 页面路由：`/accounts/zoneAgent/settlementReview`
- CRUD：C=— R=— U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 平台结算审核 | `U` | `systemCircleCheckoutAudit` | 菜单权限 |
| 平台转账 | `O` | `systemCircleCheckoutTransfer` | 菜单权限 |
| 平台备注 | `U` | `systemCircleCheckoutPlatformRemark` | 菜单权限 |

#### 财务 / 商圈代理 / 结算账号

- 页面路由：`/accounts/zoneAgent/settlementAccount`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 结算方式get | `R` | `systemCircleAgentGetSettlementMethod` | 菜单权限 |
| 结算方式post | `O` | `systemCircleAgentSetSettlementMethod` | 菜单权限 |

#### 财务 / 店铺结算

- 页面路由：`/mer/accounts`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 财务 / 店铺结算 / 分账管理

- 页面路由：`/merchant/applyList`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemOrderProfitsharingLst` | 菜单权限 |
| 重新分账 | `O` | `systemOrderProfitsharingAgain` | 菜单权限 |
| 导出 | `R` | `systemOrderProfitsharingExport` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `systemStoreExcelDownload` | 菜单权限 |

#### 财务 / 店铺结算 / 平台账单

- 页面路由：`/accounts/statement`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemFinancialRecordLst` | 菜单权限 |
| 统计 | `R` | `systemFinancialRecordTitle` | 菜单权限 |
| 详情 | `R` | `systemFinancialRecordDetail` | 菜单权限 |
| 导出 | `R` | `systemFinancialRecordDetailExport` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `systemStoreExcelDownload` | 菜单权限 |

#### 财务 / 店铺结算 / 店铺账单

- 页面路由：`/accounts/merchantBill`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 商户列表 | `R` | `systemFinancialRecordMerLst` | 菜单权限 |
| 商户统计 | `R` | `systemFinancialRecordMerAcountsLst` | 菜单权限 |
| 商户财务头部统计 | `R` | `systemFinancialRecordMerTitle` | 菜单权限 |
| 商户财务详情 | `R` | `systemFinancialRecordMerDetail` | 菜单权限 |
| 商户财务导出 | `R` | `systemFinancialRecordMerExcel` | 菜单权限 |
| 商户列表 | `R` | `systemFinancialRecordMerLst` | 菜单权限 |
| 商户统计 | `R` | `systemFinancialRecordMerAcountsLst` | 菜单权限 |
| 商户财务头部统计 | `R` | `systemFinancialRecordMerTitle` | 菜单权限 |
| 商户财务详情 | `R` | `systemFinancialRecordMerDetail` | 菜单权限 |
| 商户财务导出 | `R` | `systemFinancialRecordMerExcel` | 菜单权限 |

#### 财务 / 店铺结算 / 转账记录

- 页面路由：`/accounts/transferRecord`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemFinancialList` | 菜单权限 |
| 详情 | `R` | `systemFinancialDetail` | 菜单权限 |
| 编辑 | `U` | `systemFinancialUpdate` | 菜单权限 |
| 修改状态 | `U` | `systemFinancialSwitchStatus` | 菜单权限 |
| 备注 | `U` | `systemFinancialMark` | 菜单权限 |
| 统计 | `R` | `systemFinancialTitle` | 菜单权限 |
| 导出 | `R` | `systemFinancialExport` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `systemStoreExcelDownload` | 菜单权限 |

#### 财务 / 店铺结算 / 转账设置

- 页面路由：`/accounts/settings`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 配置信息 | `RU` | `systemOrderProfitsharingGetConfig` | 菜单权限 |
| 配置保存 | `CU` | `systemOrderProfitsharingSetConfig` | 菜单权限 |

#### 财务 / 用户结算

- 页面路由：`/accounts/record`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| （仅页面入口，未解析到菜单按钮或同路径路由） | `R` | — | 待补 |

#### 财务 / 用户结算 / 充值记录

- 页面路由：`/accounts/bill`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemUserRechargeList` | 菜单权限 |
| 统计 | `R` | `systemUserRechargeTotal` | 菜单权限 |

#### 财务 / 用户结算 / 提现管理

- 页面路由：`/accounts/extract`
- CRUD：C=✓ R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 申请列表 | `CR` | `systemUserExtractLst` | 菜单权限 |
| 审核 | `U` | `systemUserExtractSwitchStatus` | 菜单权限 |
| 导出 | `R` | `systemUserExtractExport` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `systemStoreExcelDownload` | 菜单权限 |
| 详情 | `R` | `systemUserExtractDetail` | 菜单权限 |

#### 财务 / 用户结算 / 资金记录

- 页面路由：`/accounts/capital`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemUserBillList` | 菜单权限 |
| 导出 | `R` | `systemUserBillExport` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `systemStoreExcelDownload` | 菜单权限 |

#### 财务 / 资金流水

- 页面路由：`/accounts/capitalFlow`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 列表 | `R` | `systemFinancialRecordList` | 菜单权限 |
| 导出 | `R` | `systemFinancialRecordExport` | 菜单权限 |
| 统计 | `R` | `systemFinancialCount` | 菜单权限 |
| 导出列表 | `R` | `systemStoreExcelLst` | 菜单权限 |
| 导出下载 | `R` | `systemStoreExcelDownload` | 菜单权限 |

### 首页

#### 首页

- 页面路由：`/`
- CRUD：C=✓ R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 分账商户申请备注 | `C` | `systemMerchantApplymentsMarrkSave` | 路由 |
| 列表 | `R` | `systemOrderProfitsharingLst` | 路由 |

#### 首页 / 商品统计

- 页面路由：`/statistic/product`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 顶部统计 | `R` | `systemAnalyticsProductTop` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsProductLineChart` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsProductTypePieChart` | 菜单权限 |
| 顶部统计 | `R` | `systemAnalyticsProductTop` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsProductLineChart` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsProductTypePieChart` | 菜单权限 |
| 顶部统计 | `R` | `systemAnalyticsProductTop` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsProductLineChart` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsProductTypePieChart` | 菜单权限 |

#### 首页 / 控制台

- 页面路由：`/dashboard`
- CRUD：C=— R=✓ U=✓ D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 主要数据 | `O` | `systemStatisticsMain` | 菜单权限 |
| 当日订单 | `O` | `systemStatisticsOrder` | 菜单权限 |
| 当日订单数 | `O` | `systemStatisticsOrderNum` | 菜单权限 |
| 当日支付人数 | `O` | `systemStatisticsOrderUser` | 菜单权限 |
| 商户销量 | `O` | `systemStatisticsMerchantStock` | 菜单权限 |
| 商户访问量 | `O` | `systemStatisticsMerchantRate` | 菜单权限 |
| 商户销售额 | `O` | `systemStatisticsMerchantVisit` | 菜单权限 |
| 用户数据 | `O` | `systemStatisticsUserData` | 菜单权限 |
| 成交用户 | `O` | `systemStatisticsUser` | 菜单权限 |
| 成交用户占比 | `O` | `systemStatisticsUserRate` | 菜单权限 |
| 未处理业务统计 | `RU` | `systemStatisticsAdminCount` | 菜单权限 |
| 待办事项 | `O` | `systemStatisticsAdminTodo` | 菜单权限 |
| 商户销量排行 | `O` | `systemStatisticsMerchantTop` | 菜单权限 |

#### 首页 / 数据大屏

- 页面路由：`/data-screen/index`
- CRUD：C=— R=— U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 数据大屏 | `O` | `systemDataScreen` | 菜单权限 |

#### 首页 / 用户统计

- 页面路由：`/statistic/member`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 顶部统计 | `R` | `systemAnalyticsUserTop` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsUserLineChart` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsUserTypePieChart` | 菜单权限 |

#### 首页 / 订单统计

- 页面路由：`/statistic/order`
- CRUD：C=— R=✓ U=— D=—

| 操作 | CRUD | 标识 | 来源 |
| --- | --- | --- | --- |
| 顶部统计 | `R` | `systemAnalyticsOrderTop` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsOrderLineChart` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsOrderTypePieChart` | 菜单权限 |
| 顶部统计 | `R` | `systemAnalyticsOrderTop` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsOrderLineChart` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsOrderTypePieChart` | 菜单权限 |
| 顶部统计 | `R` | `systemAnalyticsOrderTop` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsOrderLineChart` | 菜单权限 |
| 折线图统计 | `R` | `systemAnalyticsOrderTypePieChart` | 菜单权限 |
