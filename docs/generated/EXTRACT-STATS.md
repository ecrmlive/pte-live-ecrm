# CRMEB 功能抽取统计（机器生成）

> 生成自 `~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0`  
> 用途：功能基线清单，不是业务实现代码。

| 指标 | 数量 |
| --- | ---: |
| 菜单/权限节点总计 | 2419 |
| 平台侧节点 | 1726 |
| 商户侧节点 | 693 |
| 平台可见导航菜单 | 292 |
| 商户可见导航菜单 | 100 |
| 平台按钮级权限 | 1426 |
| 商户按钮级权限 | 591 |
| 命名路由 | 1402 |
| Repository 文件 | 186 |
| 小程序页面 | 183 |
| 数据表 | 165 |

## 生成文件索引

| 文件 | 说明 |
| --- | --- |
| `crmeb-menu-platform-nav.md` | 平台可见菜单树 |
| `crmeb-menu-merchant-nav.md` | 商户可见菜单树 |
| `crmeb-menu-platform-tree.md` | 平台全量（含按钮权限） |
| `crmeb-menu-merchant-tree.md` | 商户全量（含按钮权限） |
| `crmeb-system-menu.tsv` | 全量菜单表 |
| `crmeb-route-aliases.tsv` | 路由 name / alias / path |
| `crmeb-repositories.txt` | Repository 清单 |
| `crmeb-mp-pages.txt` | 小程序页面路径 |

## 平台后台模块摘要（可见菜单）

### 首页 `/`
- 控制台 `/dashboard`
- 数据大屏 `/data-screen/index`
- 商品统计 `/statistic/product`
- 订单统计 `/statistic/order`
- 用户统计 `/statistic/member`

### 店铺 `/mer`
- **店铺管理** `/mer/mer`
  - 店铺列表 `/merchant/list`
  - 店铺分类 `/merchant/classify`
  - 店铺分组 `/merchant/grouping`
  - 店铺类型 `/merchant/type`
  - 店铺入驻申请 `/merchant/application`
  - 店铺分账申请 `/merchant/applyments`
- **店铺设置** `/mer/store`
  - 店铺保证金 `/merchant/deposit_list`
  - 保证金配置 `/systemForm/Basics/margin`
  - 店铺菜单 `/merchant/system`
  - 说明提示 `/merchant/type/description`
- **商户管理** `/merchant`
  - 商户列表 `/merchant/index`
  - 商户入驻审核 `/merchant/review`
  - 商户管理员 `/merchant/admin-list`
  - 商户设置 `/merchant/apply-setting`
- **区域代理** `/business-zones/manage`
  - 区域列表 `/business-zones/index`
  - 代理人员 `/business-zones/agents`
  - 代理审核 `/business-zones/agent-review`
  - 代理设置 `/business-zones/settings`

### 商品 `/product`
- 商品管理 `/product/examine`
- 商品分类 `/product/classify`
- **品牌管理** `/product/brand`
  - 品牌分类 `/product/band/brandClassify`
  - 品牌列表 `/product/band/brandList`
- 评论管理 `/product/comment`
- 保障服务 `/product/guarantee`
- 商品标签 `/product/label`
- **商品参数** `/product/specsMain`
  - 店铺商品参数 `/product/merSpecs`
  - 平台商品参数 `/product/specs`
- 价格说明 `/product/priceDescription`
- 活动标签 `/product/activityLabel`

### 订单 `/order`
- 订单列表 `/order/list`
- 退款订单 `/order/refund`
- 核销记录 `/order/cancellation`

### 分销 `/promoter`
- 分销员列表 `/promoter/user`
- 分销配置 `/systemForm/Basics/distribution_tabs`
- **分销等级** `brokerage`
  - 分销员等级 `/promoter/membership_level`
  - 等级规则 `/promoter/distribution`
- 提现银行 `/group/config/76`
- 分销特权 `/group/config/75`
- 分销海报 `/group/config/68`
- 分销礼包 `/promoter/gift`
- 佣金说明 `/promoter/commission`
- 分销订单 `/promoter/orderList`
- 分销说明 `/promoter/retail`

### 营销 `/marketing`
- **平台优惠券** `/marketing/platform_coupon`
  - 优惠券列表 `/marketing/platform_coupon/list`
  - 领取记录 `/marketing/platform_coupon/couponRecord`
  - 发送记录 `/marketing/platform_coupon/couponSend`
  - 使用说明 `/marketing/platform_coupon/instructions`
- **商户优惠券** `/marketing/coupon`
  - 优惠券列表 `/marketing/coupon/list`
  - 领取记录 `/marketing/coupon/user`
- **秒杀** `/marketing/seckill`
  - 秒杀配置 `/marketing/seckill/seckillConfig`
  - 秒杀管理 `/marketing/seckill/list`
  - 秒杀活动 `/marketing/seckill/store_seckill/list`
- **直播** `/marketing2`
  - 直播间管理 `/marketing/studio/list`
  - 直播商品管理 `/marketing/broadcast/list`
- **预售** `/marketing/presell`
  - 预售商品 `/marketing/presell/list`
  - 预售协议 `/marketing/presell/agreement`
- **助力** `/assist`
  - 活动商品 `/marketing/assist/goods_list`
  - 助力活动 `/marketing/assist/list`
- **拼团** `/marketing/combination`
  - 拼团设置 `/marketing/combination/combination_set`
  - 拼团商品列表 `/marketing/combination/combination_goods`
  - 拼团活动列表 `/marketing/combination/combination_list`
- **积分** `/marketing/integral`
  - 积分配置 `/marketing/integral/config`
  - 积分日志 `/marketing/integral/log`
  - 商品分类 `/marketing/integral/classify`
  - 商品列表 `/marketing/integral/proList`
  - 积分订单 `/marketing/integral/orderList`
- 活动氛围图 `/marketing/atmosphere/list`
- 活动边框图 `/marketing/border/list`
- 专场列表 `/group/topic/94`
- 优惠套餐 `/marketing/discounts/list`
- **余额充值** `/banlace`
  - 余额设置 `/systemForm/Basics/balance`
  - 余额充值配置 `/group/config/69`
- 报名活动 `/marketing/application/list`

### 用户 `/user`
- 用户列表 `/user/list`
- 用户分组 `/user/group`
- 用户标签 `/user/label`
- **用户反馈** `/feedback`
  - 反馈分类 `/feedback/classify`
  - 反馈列表 `/feedback/list`
- 用户协议 `/user/agreement`
- 搜索记录 `/user/searchRecord`
- **用户等级** `/user/member`
  - 等级管理 `/user/member/list`
  - 等级配置 `/systemForm/Basics/members`
  - 等级说明 `/user/member/description`
- 用户设置 `/user/setup_user`
- **付费会员** `/user/svip`
  - 付费会员配置 `/systemForm/Basics/svip`
  - 会员类型 `/user/member/type`
  - 会员权益 `/user/member/equity`
  - 会员记录 `/user/member/record`
  - 会员协议 `/user/member/vipAgreement`

### 内容 `/content`
- **文章** `/cms`
  - 文章管理 `/cms/article`
  - 文章分类 `/cms/articleCategory`
- **社区** `/community`
  - 社区分类 `/community/category`
  - 社区话题 `/community/topic`
  - 社区内容 `/community/list`
  - 社区评论 `/community/reply`
  - 社区配置 `/systemForm/Basics/community`

### 财务 `/accounts`
- **商圈代理** `/accounts/zoneAgent`
  - 结算审核 `/accounts/zoneAgent/settlementReview`
  - 提成流水 `/accounts/zoneAgent/commissionRecord`
- **店铺结算** `/mer/accounts`
  - 平台账单 `/accounts/statement`
  - 店铺账单 `/accounts/merchantBill`
  - 转账记录 `/accounts/transferRecord`
  - 转账设置 `/accounts/settings`
  - 分账管理 `/merchant/applyList`
- **用户结算** `/accounts/record`
  - 提现管理 `/accounts/extract`
  - 充值记录 `/accounts/bill`
  - 资金记录 `/accounts/capital`
- 资金流水 `/accounts/capitalFlow`
- **发票管理** `/accounts/accounts`
  - 发票列表 `/accounts/receipt`
  - 发票说明 `/accounts/invoiceDesc`

### 应用 `/apploction`
- **公众号** `/app/wechat`
  - 微信菜单 `/app/wechat/menus`
  - 自动回复 `/admin/app/wechat/reply`
  - 图文管理 `/app/wechat/newsCategory`
- **小程序** `/app/routine`
  - 小程序下载 `/app/routine/download`

### 装修 `/theme`
- 主题风格 `/setting/theme_style`
- 页面装修 `/setting/diy/list`
- 店铺街 `/setting/diy/store`
- 个人中心 `/setting/diy/personal`
- 店铺模板 `/setting/merchant/diyList`
- 商品详情 `/setting/diy/product_detail`
- 页面配置 `/setting/system_visualization_data`
- **页面链接** `/setting/page`
  - 平台页面分类 `/setting/diy/plantform/category/list`
  - 平台页面链接 `/setting/diy/links/list`
  - 商户页面分类 `/setting/diy/merchant/category/list`
  - 商户页面链接 `/setting/diy/merLink/list`
- 素材管理 `/config/picture`
- 系统表单 `/systemForm/form_list`
- 悬浮菜单 `/setting/fab`
- 商品分类 `/setting/product_category`

### 客服 `/service`
- 客服自动回复 `/systemForm/customer_keyword`
- 客服列表 `/service/customer/list`
- 客服设置 `/systemForm/Basics/service`

### 设置 `/settings`
- **系统设置** `/sys`
  - 系统设置 `/systemForm/Basics/system_tabs`
  - 接口配置 `/systemForm/Basics/extend_tabs`
  - 存储管理 `/setting/storage`
  - 一号通 `/serve`
- **商城设置** `/shop`
  - 商城设置 `/systemForm/Basics/shop_tabs`
  - 支付设置 `/systemForm/Basics/pay_tabs`
  - 热门搜索 `/group/config/67`
  - 协议规则 `/setting/agreements`
  - 配送配置 `/delivery_config`
- **权限管理** `/setting`
  - 角色权限 `/setting/systemRole`
  - 管理员管理 `/setting/systemAdmin`
  - 菜单管理 `/setting/menu`
- **消息管理** `/notice`
  - 公告管理 `/station/notice`
  - 消息管理 `/setting/notification/index`
- **应用配置** `/app_config`
  - 公众号配置 `/systemForm/Basics/wechat`
  - 小程序配置 `/systemForm/Basics/smallapp`
  - APP配置 `/systemForm/Basics/wechat_open_app`
  - 上传校验文件 `/app/wechat/file`
  - APP升级配置 `/systemForm/Basics/app_version`

### 维护 `/safe`
- **开发配置** `/safe/exploit`
  - 组合数据 `/group/list`
- **安全维护** `/maintain`
  - 数据备份 `/maintain/dataBackup`
  - 商业授权 `/setting/system/maintain/auth`
  - 缓存清除 `/maintain/cache`
- 配置分类 `/config/classify`
- 配置管理 `/config/setting`
- 操作日志 `/setting/systemLog`
- 导出记录 `/group/exportList`
- 页面链接 `/safe/pageLinks`

### 短信模板 `/sms/template`

### 申请记录 `/sms/applyList`


## 商户后台模块摘要（可见菜单）

### 首页 `/dashboard`
- 控制台 `/dashboard`
- 商品统计 `/statistic/product`
- 订单统计 `/statistic/order`

### 商品 `/product`
- 商品列表 `/product/list`
- 卡密列表 `/product/cdkey`
- 商品分类 `/product/classify`
- 商品规格 `/product/attr`
- 商品参数 `/product/specs`
- 商品单位 `/product/unit`
- 商品标签 `/product/label`
- 服务模板 `/config/guarantee`

### 订单 `/order`
- 订单管理 `/order/list`
- 退款订单 `/order/refund`
- 核销记录 `/order/cancellation`
- 预约服务 `/order/reservation`
- 预约设置 `/product/reservation`
- 商品评价 `/product/reviews`
- 代客下单 `/order/customer`

### 营销 `/marketing`
- **优惠券** `/marketing/coupon`
  - 优惠券列表 `/marketing/coupon/list`
  - 领取记录 `/marketing/coupon/user`
  - 发送记录 `/marketing/coupon/send`
- **直播** `/`
  - 直播间管理 `/marketing/studio/list`
  - 直播商品管理 `/marketing/broadcast/list`
  - 直播助手 `/marketing/studio/assistant`
- **秒杀** `/marketing/seckill/list`
  - 秒杀活动 `/marketing/seckill/store_seckill/list`
  - 秒杀商品 `/marketing/seckill/product/list`
- 预售 `/marketing/presell/list`
- **助力** `/assist`
  - 助力商品 `/marketing/assist/list`
  - 助力活动 `/marketing/assist/assist_set`
- **拼团** `/marketing/combination`
  - 拼团商品列表 `/marketing/combination/combination_goods`
  - 拼团活动列表 `/marketing/combination/combination_list`
- **积分** `/marketing/integral`
  - 积分配置 `/marketing/integral/config`
  - 积分日志 `/marketing/integral/log`
- 专场列表 `/group/topic/95`
- 优惠套餐 `/marketing/discounts/list`
- 逛逛社区 `/community/list`

### 财务 `/accounts`
- 资金流水 `/accounts/capitalFlow`
- 发票管理 `/order/invoice`
- 转账记录 `/accounts/transManagement`
- 收款方式 `/accounts/payType`
- 账单管理 `/accounts/statement`
- 分账管理 `/systemForm/applyList`
- 申请分账商户 `/systemForm/applyments`

### 用户 `/user`
- **标签管理** `/user/_label`
  - 手动标签 `/user/label`
  - 自动标签 `/user/maticlabel`
- 用户管理 `/user/list`
- 搜索记录 `/user/searchRecord`

### 员工 `/server`
- **店员管理** `/server_manage`
  - 店员列表 `/config/service`
  - 自动回复 `/systemForm/customer_keyword`
  - 店员配置 `/systemForm/Basics/mer_service_switch`
- **配送人员** `/delivery/personnel_manage/index`
  - 配送员管理 `/delivery/personnel_manage`
  - 配送统计 `/delivery/delivice_statistic`
- **服务人员** `/config/service_staff`
  - 服务人员 `/config/service_staff`
  - 服务统计 `/config/service_statistic`

### 装修 `/devise/`
- 装修 `/devise/diy/list`
- 素材管理 `/config/picture`
- 系统表单 `/systemForm/form_list`
- 商品分类 `/devise/diy/product_category`

### 设置 `/config`
- 店铺信息 `/systemForm/modifyStoreInfo`
- 店铺配置 `/systemForm/Basics/mer_base`
- 付费会员 `/systemForm/Basics/svip`
- **打印配置** `/setting/printer`
  - 打印配置 `/systemForm/Basics/printer_tabs`
  - 小票打印 `/setting/printer/list`
- **一号通** `/one_setting`
  - 自有一号通 `/setting/sms/sms_account/index`
  - 平台一号通 `/setting/sms/sms_config/index`
  - 配置管理 `/setting/sms/dumpConfig`
- **快递配送** `/city`
  - 运费模板 `/config/freight/shippingTemplates`
  - 物流公司 `/config/freight/express`
- **权限管理** `/setting`
  - 身份管理 `/setting/systemRole`
  - 管理员管理 `/setting/systemAdmin`
  - 操作日志 `/setting/systemLog`
- 开放账户 `/systemForm/openAuth/list`
- **同城配送** `/delivery`
  - 配送记录 `/delivery/usage_record`
  - 配送门店 `/delivery/delivery_point`
  - 配送设置 `/setting/delivery`

### 公告 `/station/notice`

