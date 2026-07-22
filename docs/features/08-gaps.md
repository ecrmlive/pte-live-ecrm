# 功能点缺口清单（已结案 · 基线已锁定）

> 基于菜单按钮 + 路由补全后的结果，并已对照 `crmeb-system-menu.tsv` 分类。
> **2026-07-21 源码核对**：原「53 个无操作页面」→ 39 目录/配置/DIY + **14 叶子已结案**（路径错挂/共享 API/协议页）。**无操作空页真缺口 = 0**；**CRUD 矩阵缺口真缺口 = 0**（均为只读/配置/模块壳/路由漏抽）。  
> **功能基线已锁定**（同日用户确认）。

## 结论摘要

| 分类 | 数量 | 是否阻塞基线 |
| --- | ---: | --- |
| 空页（目录/配置/DIY/叶子结案） | 53 | 否（真缺口 0） |
| CRUD 不全 — 已判定非缺口/假缺口 | 169 | 否 |
| CRUD 不全 — 仍待核对 | 0 | 否 |
| **原 CRUD 不全合计** | **169** | **已全部结案** |

## A. 已判定非缺口（不必再补按钮）

| 端 | 模块 | 页面路由 | 类型 | 说明 |
| --- | --- | --- | --- | --- |
| 商户后台 | 营销 / 助力 | `/assist` | 目录壳 | 子菜单 2 个 |
| 商户后台 | 设置 / 快递配送 | `/city` | 目录壳 | 子菜单 2 个 |
| 商户后台 | 员工 / 配送人员 / 配送统计 | `/delivery/delivice_statistic` | 拼写/历史入口 | 源码菜单即此路径；重建可改前端路由名 |
| 商户后台 | 营销 / 秒杀 | `/marketing/seckill/list` | 目录壳 | 子菜单 2 个 |
| 商户后台 | 设置 / 一号通 | `/one_setting` | 目录壳 | 子菜单 3 个 |
| 商户后台 | 员工 | `/server` | 目录壳 | 子菜单 3 个 |
| 商户后台 | 员工 / 店员管理 | `/server_manage` | 目录壳 | 子菜单 3 个 |
| 商户后台 | 员工 / 店员管理 / 店员配置 | `/systemForm/Basics/mer_service_switch` | 配置表单页 | systemForm 配置保存，非菜单按钮 CRUD |
| 商户后台 | 设置 / 打印配置 / 打印配置 | `/systemForm/Basics/printer_tabs` | 配置表单页 | systemForm 配置保存，非菜单按钮 CRUD |
| 商户后台 | 设置 / 付费会员 | `/systemForm/Basics/svip` | 配置表单页 | systemForm 配置保存，非菜单按钮 CRUD |
| 商户后台 | 用户 / 标签管理 | `/user/_label` | 目录壳 | 子菜单 2 个 |
| 平台后台 | 财务 / 发票管理 | `/accounts/accounts` | 财务入口 | 操作在 systemFinancial* / 子菜单 |
| 平台后台 | 财务 / 用户结算 | `/accounts/record` | 财务入口 | 操作在 systemFinancial* / 子菜单 |
| 平台后台 | 设置 / 应用配置 | `/app_config` | 目录壳 | 子菜单 5 个 |
| 平台后台 | 应用 | `/apploction` | 拼写/历史入口 | 源码菜单即此路径；重建可改前端路由名 |
| 平台后台 | 营销 / 助力 | `/assist` | 目录壳 | 子菜单 2 个 |
| 平台后台 | 营销 / 余额充值 | `/banlace` | 拼写/历史入口 | 源码菜单即此路径；重建可改前端路由名 |
| 平台后台 | 店铺 / 区域代理 | `/business-zones/manage` | 目录壳 | 子菜单 4 个 |
| 平台后台 | 内容 | `/content` | 目录壳 | 子菜单 2 个 |
| 平台后台 | 设置 / 商城设置 / 配送配置 | `/delivery_config` | 目录壳 | 子菜单 3 个 |
| 平台后台 | 营销 / 直播 | `/marketing2` | 拼写/历史入口 | 源码菜单即此路径；重建可改前端路由名 |
| 平台后台 | 财务 / 店铺结算 | `/mer/accounts` | 财务入口 | 操作在 systemFinancial* / 子菜单 |
| 平台后台 | 店铺 / 店铺管理 | `/mer/mer` | 目录壳 | 子菜单 6 个 |
| 平台后台 | 店铺 / 店铺设置 | `/mer/store` | 目录壳 | 子菜单 4 个 |
| 平台后台 | 设置 / 消息管理 | `/notice` | 目录壳 | 子菜单 2 个 |
| 平台后台 | 商品 / 品牌管理 | `/product/brand` | 目录壳 | 子菜单 2 个 |
| 平台后台 | 维护 | `/safe` | 目录壳 | 子菜单 7 个 |
| 平台后台 | 维护 / 开发配置 | `/safe/exploit` | 目录壳 | 子菜单 1 个 |
| 平台后台 | 维护 / 页面链接 | `/safe/pageLinks` | 页面链接选择器 | 工具页 |
| 平台后台 | 设置 / 系统设置 / 一号通 | `/serve` | 目录壳 | 子菜单 5 个 |
| 平台后台 | 装修 / 个人中心 | `/setting/diy/personal` | 装修 DIY | 设计器能力在 DIY API |
| 平台后台 | 装修 / 商品详情 | `/setting/diy/product_detail` | 装修 DIY | 设计器能力在 DIY API |
| 平台后台 | 装修 / 店铺街 | `/setting/diy/store` | 装修 DIY | 设计器能力在 DIY API |
| 平台后台 | 装修 / 页面链接 | `/setting/page` | 页面链接选择器 | 工具页 |
| 平台后台 | 设置 | `/settings` | 目录壳 | 子菜单 5 个 |
| 平台后台 | 设置 / 商城设置 | `/shop` | 目录壳 | 子菜单 5 个 |
| 平台后台 | 装修 | `/theme` | 目录壳 | 子菜单 13 个 |
| 平台后台 | 用户 / 付费会员 | `/user/svip` | 目录壳 | 子菜单 5 个 |
| 平台后台 | 分销 / 分销等级 | `brokerage` | 目录壳 | 子菜单 2 个 |

## B. 原「叶子无按钮」14 页 — 已按源码结案（非真缺口）

> 对照菜单 SQL、`route/**/*.php` 的 `_path`、协议 Cache 键。结论：**不必再补菜单按钮**；重建时按下列真实能力对接。

| 端 | 菜单路由 | 结案 | 真实能力位置 |
| --- | --- | --- | --- |
| 商户后台 | `/marketing/broadcast/list` | 菜单路径过时 | 直播商品 API `_path=/marketing/studio/list`（`BroadcastGoods` CRUD） |
| 商户后台 | `/marketing/integral/config` | 配置挂在日志页 | `merchantConfigForm/Save` append 在 `_path=/marketing/integral/log`；本页无独立按钮属预期 |
| 商户后台 | `/product/reservation` | 菜单路径不一致 | 预约日历 API `_path=/order/reservation`（`ReservationService/list`） |
| 商户后台 | `/setting/sms/sms_account/index` | 入口/历史路径 | 一号通能力在 `route/merchant/yihaotong.php`，主 `_path=/setting/sms/dumpConfig` |
| 平台后台 | `/promoter/retail` | 协议富文本页 | `CacheRepository::PROMOTER_EXPLAIN`；协议 API `/sys/agreement/:key`（`_path=/setting/agreements`） |
| 平台后台 | `/product/activityLabel` | **假缺口** | 完整 CRUD：`activity/cate` + `activity/label`，`_path=product/activityLabel`（无前导 `/`，匹配器易漏） |
| 平台后台 | `/app/wechat/reply/keyword` | 共享 API | 三页共用 `wechat/reply`，`_path=/admin/app/wechat/reply`（keyword/subscribe/default 为 UI 模式） |
| 平台后台 | `/app/wechat/reply/follow/subscribe` | 同上 | 同上 |
| 平台后台 | `/app/wechat/reply/index/default` | 同上 | 同上 |
| 平台后台 | `/business-zones/settings` | 配置/结算入口 | 商圈能力在 `route/admin/circle.php`（agents/review/settlement*）；设置页无独立按钮属预期 |
| 平台后台 | `/merchant/apply-setting` | 入驻配置/协议 | 入驻审核 `/merchant/application`；协议键 `business_entry_agree` 等走 agreement API |
| 平台后台 | `/user/member/vipAgreement` | 协议富文本页 | `CacheRepository::SYS_SVIP`；协议 API `/sys/agreement/:key` |
| 平台后台 | `/freight/city/list` | `_path` 挂快递页 | `CityArea` CRUD `_path=/freight/express`（城市数据与快递公司同权限组） |
| 平台后台 | `/sms/user` | 历史路径 | 一号通/短信账户 API 在 `route/admin/yihaotong.php`，`_path=/setting/sms/sms_config/index` |

**空页阻塞项：已清零（53→0 真缺口）。** 仍待业务判断的是下方 C 类「CRUD 不完整」169 页（只读页可无 C/D）。

## C. CRUD 不完整页面分类（原 169）

> 自动规则对照 `route/**/*.php` 的 `_path` 别名 + 页面名称语义。
> **非缺口/假缺口不阻塞基线**；C.3 为仍建议人工扫一眼的少量页面。

### C.1 假缺口（路由/按钮已有，矩阵漏抽）（18）

| 端 | 模块 | 页面路由 | 矩阵缺口 | 已有 | 操作数 | 说明 |
| --- | --- | --- | --- | --- | ---: | --- |
| 商户后台 | 员工 / 服务人员 / 服务人员 | `/config/service_staff` | CUD | R | 2 | 同 _path 路由别名已含 CDRU，功能矩阵未抽全按钮 |
| 商户后台 | 营销 / 秒杀 / 秒杀活动 | `/marketing/seckill/store_seckill/list` | UD | CR | 9 | 操作数=9，路由含 R |
| 商户后台 | 设置 / 一号通 / 平台一号通 | `/setting/sms/sms_config/index` | CUD | R | 5 | 同 _path 路由别名已含 CDRU，功能矩阵未抽全按钮 |
| 商户后台 | 公告 | `/station/notice` | CU | RD | 5 | 同 _path 路由别名已含 CDRU，功能矩阵未抽全按钮 |
| 平台后台 | 维护 / 安全维护 / 数据备份 | `/maintain/dataBackup` | CU | RD | 8 | 操作数=8，路由含 DR |
| 平台后台 | 营销 / 助力 / 活动商品 | `/marketing/assist/goods_list` | CD | RU | 21 | 操作数=21，路由含 RU |
| 平台后台 | 营销 / 助力 / 助力活动 | `/marketing/assist/list` | CUD | R | 6 | 同 _path 路由别名已含 CDRU，功能矩阵未抽全按钮 |
| 平台后台 | 营销 / 直播 / 直播商品管理 | `/marketing/broadcast/list` | C | RUD | 6 | 操作数=6，路由含 DRU |
| 平台后台 | 营销 / 拼团 / 拼团商品列表 | `/marketing/combination/combination_goods` | CD | RU | 24 | 同 _path 路由别名已含 CDRU，功能矩阵未抽全按钮 |
| 平台后台 | 营销 / 拼团 / 拼团活动列表 | `/marketing/combination/combination_list` | CUD | R | 6 | 操作数=6，路由含 RU |
| 平台后台 | 营销 / 商户优惠券 / 优惠券列表 | `/marketing/coupon/list` | CUD | R | 3 | 同 _path 路由别名已含 CDRU，功能矩阵未抽全按钮 |
| 平台后台 | 营销 / 积分 / 积分订单 | `/marketing/integral/orderList` | C | RUD | 9 | 操作数=9，路由含 DRU |
| 平台后台 | 营销 / 预售 / 预售商品 | `/marketing/presell/list` | CD | RU | 21 | 操作数=21，路由含 DRU |
| 平台后台 | 营销 / 秒杀 / 秒杀管理 | `/marketing/seckill/list` | C | RUD | 29 | 操作数=29，路由含 DU |
| 平台后台 | 店铺 / 店铺管理 / 店铺入驻申请 | `/merchant/application` | C | RUD | 8 | 同 _path 路由别名已含 CDRU，功能矩阵未抽全按钮 |
| 平台后台 | 订单 / 退款订单 | `/order/refund` | CD | RU | 21 | 同 _path 路由别名已含 CDRU，功能矩阵未抽全按钮 |
| 平台后台 | 分销 / 分销礼包 | `/promoter/gift` | CD | RU | 6 | 操作数=6，路由含 RU |
| 平台后台 | 分销 / 分销订单 | `/promoter/orderList` | CD | RU | 8 | 操作数=8，路由含 RU |

### C.2a 配置表单 systemForm（21）

| 端 | 模块 | 页面路由 | 矩阵缺口 | 已有 | 操作数 | 说明 |
| --- | --- | --- | --- | --- | ---: | --- |
| 商户后台 | 设置 / 店铺配置 | `/systemForm/Basics/mer_base` | CRD | U | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 商户后台 | 财务 / 分账管理 | `/systemForm/applyList` | CUD | R | 4 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 商户后台 | 财务 / 申请分账商户 | `/systemForm/applyments` | D | CRU | 5 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 商户后台 | 设置 / 店铺信息 | `/systemForm/modifyStoreInfo` | RD | CU | 2 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 设置 / 应用配置 / APP升级配置 | `/systemForm/Basics/app_version` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 营销 / 余额充值 / 余额设置 | `/systemForm/Basics/balance` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 内容 / 社区 / 社区配置 | `/systemForm/Basics/community` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 分销 / 分销配置 | `/systemForm/Basics/distribution_tabs` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 设置 / 系统设置 / 接口配置 | `/systemForm/Basics/extend_tabs` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 店铺 / 店铺设置 / 保证金配置 | `/systemForm/Basics/margin` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 用户 / 用户等级 / 等级配置 | `/systemForm/Basics/members` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 设置 / 系统设置 / 一号通 / 短信设置 / 短信配置 | `/systemForm/Basics/message` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 设置 / 商城设置 / 支付设置 | `/systemForm/Basics/pay_tabs` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 客服 / 客服设置 | `/systemForm/Basics/service` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 设置 / 商城设置 / 商城设置 | `/systemForm/Basics/shop_tabs` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 设置 / 应用配置 / 小程序配置 | `/systemForm/Basics/smallapp` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 用户 / 付费会员 / 付费会员配置 | `/systemForm/Basics/svip` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 设置 / 系统设置 / 系统设置 | `/systemForm/Basics/system_tabs` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 设置 / 应用配置 / 公众号配置 | `/systemForm/Basics/wechat` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 设置 / 应用配置 / APP配置 | `/systemForm/Basics/wechat_open_app` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |
| 平台后台 | 设置 / 商城设置 / 配送配置 / 第三方送 / 配送配置 | `/systemForm/delivery` | RD | CU | 1 | systemForm 配置保存（常见 CU，无独立列表/删除） |

### C.2b 其它配置页（9）

| 端 | 模块 | 页面路由 | 矩阵缺口 | 已有 | 操作数 | 说明 |
| --- | --- | --- | --- | --- | ---: | --- |
| 商户后台 | 设置 / 快递配送 / 物流公司 | `/config/freight/express` | CD | RU | 3 | 配置页常见 CU 无列表/删除 |
| 商户后台 | 设置 / 同城配送 / 配送设置 | `/setting/delivery` | CRD | U | 2 | 配置页常见 CU 无列表/删除 |
| 商户后台 | 设置 / 一号通 / 配置管理 | `/setting/sms/dumpConfig` | RD | CU | 1 | 配置页常见 CU 无列表/删除 |
| 平台后台 | 财务 / 店铺结算 / 转账设置 | `/accounts/settings` | D | CRU | 2 | 配置页常见 CU 无列表/删除 |
| 平台后台 | 设置 / 应用配置 / 上传校验文件 | `/app/wechat/file` | RD | CU | 4 | 配置页常见 CU 无列表/删除 |
| 平台后台 | 营销 / 拼团 / 拼团设置 | `/marketing/combination/combination_set` | RD | CU | 1 | 配置页常见 CU 无列表/删除 |
| 平台后台 | 营销 / 积分 / 积分配置 | `/marketing/integral/config` | RD | CU | 1 | 配置页常见 CU 无列表/删除 |
| 平台后台 | 装修 / 页面配置 | `/setting/system_visualization_data` | UD | CR | 3 | 配置页常见 CU 无列表/删除 |
| 平台后台 | 设置 / 系统设置 / 一号通 / 短信设置 / 短信模板 | `/sms/template` | RUD | C | 2 | 配置页常见 CU 无列表/删除 |

### C.2c 只读报表/日志/流水（36）

| 端 | 模块 | 页面路由 | 矩阵缺口 | 已有 | 操作数 | 说明 |
| --- | --- | --- | --- | --- | ---: | --- |
| 商户后台 | 财务 / 资金流水 | `/accounts/capitalFlow` | CUD | R | 5 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 商户后台 | 财务 / 账单管理 | `/accounts/statement` | CUD | R | 6 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 商户后台 | 员工 / 服务人员 / 服务统计 | `/config/service_statistic` | D | CRU | 5 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 商户后台 | 设置 / 同城配送 / 充值记录 | `/delivery/recharge_record` | UD | CR | 1 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 商户后台 | 设置 / 同城配送 / 配送记录 | `/delivery/usage_record` | CU | RD | 3 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 商户后台 | 营销 / 优惠券 / 发送记录 | `/marketing/coupon/send` | UD | CR | 1 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 商户后台 | 营销 / 优惠券 / 领取记录 | `/marketing/coupon/user` | CUD | R | 1 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 商户后台 | 营销 / 积分 / 积分日志 | `/marketing/integral/log` | D | CRU | 4 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 商户后台 | 订单 / 核销记录 | `/order/cancellation` | CUD | R | 2 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 商户后台 | 设置 / 权限管理 / 操作日志 | `/setting/systemLog` | CD | RU | 4 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 商户后台 | 首页 / 订单统计 | `/statistic/order` | CUD | R | 3 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 商户后台 | 首页 / 商品统计 | `/statistic/product` | CUD | R | 3 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 商户后台 | 用户 / 搜索记录 | `/user/searchRecord` | CUD | R | 1 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 财务 / 用户结算 / 充值记录 | `/accounts/bill` | CUD | R | 2 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 财务 / 用户结算 / 资金记录 | `/accounts/capital` | CUD | R | 4 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 财务 / 资金流水 | `/accounts/capitalFlow` | CUD | R | 5 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 财务 / 店铺结算 / 店铺账单 | `/accounts/merchantBill` | CUD | R | 10 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 财务 / 店铺结算 / 平台账单 | `/accounts/statement` | CUD | R | 6 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 财务 / 店铺结算 / 转账记录 | `/accounts/transferRecord` | CD | RU | 9 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 财务 / 商圈代理 / 提成流水 | `/accounts/zoneAgent/commissionRecord` | CUD | R | 2 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 设置 / 商城设置 / 配送配置 / 第三方送 / 充值记录 | `/delivery/recharge_record` | UD | CR | 3 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 设置 / 商城设置 / 配送配置 / 第三方送 / 配送记录 | `/delivery/usage_record` | CUD | R | 2 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 维护 / 导出记录 | `/group/exportList` | CUD | R | 2 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 营销 / 商户优惠券 / 领取记录 | `/marketing/coupon/user` | CUD | R | 1 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 营销 / 积分 / 积分日志 | `/marketing/integral/log` | CUD | R | 5 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 营销 / 平台优惠券 / 领取记录 | `/marketing/platform_coupon/couponRecord` | CUD | R | 1 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 营销 / 平台优惠券 / 发送记录 | `/marketing/platform_coupon/couponSend` | UD | CR | 1 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 订单 / 核销记录 | `/order/cancellation` | CD | RU | 9 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 设置 / 系统设置 / 一号通 / 购买记录 | `/service/purchase` | CUD | R | 2 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 维护 / 操作日志 | `/setting/systemLog` | CUD | R | 1 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 设置 / 系统设置 / 一号通 / 短信设置 / 申请记录 | `/sms/applyList` | UD | CR | 1 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 首页 / 用户统计 | `/statistic/member` | CUD | R | 3 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 首页 / 订单统计 | `/statistic/order` | CUD | R | 9 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 首页 / 商品统计 | `/statistic/product` | CUD | R | 9 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 用户 / 付费会员 / 会员记录 | `/user/member/record` | CD | RU | 4 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |
| 平台后台 | 用户 / 搜索记录 | `/user/searchRecord` | CU | RD | 5 | 日志/统计/流水页通常只需 R（+导出），缺 CUD 属预期 |

### C.2d 协议/说明富文本（9）

| 端 | 模块 | 页面路由 | 矩阵缺口 | 已有 | 操作数 | 说明 |
| --- | --- | --- | --- | --- | ---: | --- |
| 平台后台 | 财务 / 发票管理 / 发票说明 | `/accounts/invoiceDesc` | RUD | C | 1 | 协议/说明页：保存即 C/U，列表/删除非必要；API 多为 agreement/:key |
| 平台后台 | 营销 / 平台优惠券 / 使用说明 | `/marketing/platform_coupon/instructions` | RUD | C | 1 | 协议/说明页：保存即 C/U，列表/删除非必要；API 多为 agreement/:key |
| 平台后台 | 营销 / 预售 / 预售协议 | `/marketing/presell/agreement` | RUD | C | 1 | 协议/说明页：保存即 C/U，列表/删除非必要；API 多为 agreement/:key |
| 平台后台 | 店铺 / 店铺设置 / 说明提示 | `/merchant/type/description` | RUD | C | 1 | 协议/说明页：保存即 C/U，列表/删除非必要；API 多为 agreement/:key |
| 平台后台 | 分销 / 佣金说明 | `/promoter/commission` | RUD | C | 1 | 协议/说明页：保存即 C/U，列表/删除非必要；API 多为 agreement/:key |
| 平台后台 | 分销 / 分销等级 / 等级规则 | `/promoter/distribution` | RUD | C | 1 | 协议/说明页：保存即 C/U，列表/删除非必要；API 多为 agreement/:key |
| 平台后台 | 设置 / 商城设置 / 协议规则 | `/setting/agreements` | UD | CR | 2 | 协议/说明页：保存即 C/U，列表/删除非必要；API 多为 agreement/:key |
| 平台后台 | 用户 / 用户协议 | `/user/agreement` | RUD | C | 1 | 协议/说明页：保存即 C/U，列表/删除非必要；API 多为 agreement/:key |
| 平台后台 | 用户 / 用户等级 / 等级说明 | `/user/member/description` | RUD | C | 1 | 协议/说明页：保存即 C/U，列表/删除非必要；API 多为 agreement/:key |

### C.2e 模块壳/父菜单（35）

| 端 | 模块 | 页面路由 | 矩阵缺口 | 已有 | 操作数 | 说明 |
| --- | --- | --- | --- | --- | ---: | --- |
| 商户后台 | 营销 / 直播 | `/` | UD | CR | 5 | 一级/模块入口，操作在子页面 |
| 商户后台 | 财务 | `/accounts` | U | CRD | 8 | 一级/模块入口，操作在子页面 |
| 商户后台 | 首页 | `/dashboard` | CD | RU | 10 | 一级/模块入口，操作在子页面 |
| 商户后台 | 首页 / 控制台 | `/dashboard` | U | CRD | 5 | 一级/模块入口，操作在子页面 |
| 商户后台 | 设置 / 同城配送 | `/delivery` | CUD | R | 3 | 一级/模块入口，操作在子页面 |
| 商户后台 | 营销 | `/marketing` | UD | CR | 4 | 一级/模块入口，操作在子页面 |
| 商户后台 | 营销 / 拼团 | `/marketing/combination` | RUD | C | 2 | 父级菜单，CRUD 在子路由 |
| 商户后台 | 营销 / 优惠券 | `/marketing/coupon` | UD | CR | 4 | 父级菜单，CRUD 在子路由 |
| 商户后台 | 营销 / 积分 | `/marketing/integral` | CUD | R | 2 | 父级菜单，CRUD 在子路由 |
| 商户后台 | 订单 | `/order` | RD | CU | 2 | 一级/模块入口，操作在子页面 |
| 商户后台 | 商品 | `/product` | UD | CR | 2 | 一级/模块入口，操作在子页面 |
| 商户后台 | 设置 / 权限管理 | `/setting` | CUD | R | 1 | 一级/模块入口，操作在子页面 |
| 商户后台 | 用户 | `/user` | D | CRU | 6 | 一级/模块入口，操作在子页面 |
| 平台后台 | 首页 | `/` | UD | CR | 2 | 一级/模块入口，操作在子页面 |
| 平台后台 | 财务 | `/accounts` | D | CRU | 5 | 一级/模块入口，操作在子页面 |
| 平台后台 | 应用 / 小程序 | `/app/routine` | CUD | R | 2 | 父级菜单，CRUD 在子路由 |
| 平台后台 | 首页 / 控制台 | `/dashboard` | CD | RU | 13 | 一级/模块入口，操作在子页面 |
| 平台后台 | 设置 / 商城设置 / 配送配置 / 第三方送 | `/delivery` | CUD | R | 3 | 一级/模块入口，操作在子页面 |
| 平台后台 | 用户 / 用户反馈 | `/feedback` | C | RUD | 5 | 一级/模块入口，操作在子页面 |
| 平台后台 | 维护 / 安全维护 | `/maintain` | CU | RD | 2 | 一级/模块入口，操作在子页面 |
| 平台后台 | 营销 | `/marketing` | CUD | R | 2 | 一级/模块入口，操作在子页面 |
| 平台后台 | 营销 / 拼团 | `/marketing/combination` | RD | CU | 2 | 父级菜单，CRUD 在子路由 |
| 平台后台 | 营销 / 商户优惠券 | `/marketing/coupon` | CUD | R | 2 | 父级菜单，CRUD 在子路由 |
| 平台后台 | 营销 / 积分 | `/marketing/integral` | UD | CR | 2 | 父级菜单，CRUD 在子路由 |
| 平台后台 | 营销 / 平台优惠券 | `/marketing/platform_coupon` | UD | CR | 3 | 父级菜单，CRUD 在子路由 |
| 平台后台 | 店铺 | `/mer` | UD | CR | 2 | 一级/模块入口，操作在子页面 |
| 平台后台 | 店铺 / 商户管理 | `/merchant` | UD | CR | 2 | 一级/模块入口，操作在子页面 |
| 平台后台 | 订单 | `/order` | CUD | R | 3 | 一级/模块入口，操作在子页面 |
| 平台后台 | 商品 | `/product` | U | CRD | 5 | 一级/模块入口，操作在子页面 |
| 平台后台 | 分销 | `/promoter` | CUD | R | 8 | 一级/模块入口，操作在子页面 |
| 平台后台 | 客服 | `/service` | D | CRU | 4 | 一级/模块入口，操作在子页面 |
| 平台后台 | 设置 / 权限管理 | `/setting` | D | CRU | 7 | 一级/模块入口，操作在子页面 |
| 平台后台 | 设置 / 系统设置 / 一号通 / 短信设置 | `/sms` | UD | CR | 2 | 一级/模块入口，操作在子页面 |
| 平台后台 | 设置 / 系统设置 | `/sys` | CD | RU | 2 | 一级/模块入口，操作在子页面 |
| 平台后台 | 用户 | `/user` | UD | CR | 7 | 一级/模块入口，操作在子页面 |

### C.2f 审核页（3）

| 端 | 模块 | 页面路由 | 矩阵缺口 | 已有 | 操作数 | 说明 |
| --- | --- | --- | --- | --- | ---: | --- |
| 平台后台 | 财务 / 商圈代理 / 结算审核 | `/accounts/zoneAgent/settlementReview` | CRD | U | 3 | 审核页以同意/拒绝(U)为主，无前台创建(C) |
| 平台后台 | 店铺 / 区域代理 / 代理审核 | `/business-zones/agent-review` | CRD | U | 1 | 审核页以同意/拒绝(U)为主，无前台创建(C) |
| 平台后台 | 店铺 / 商户管理 / 商户入驻审核 | `/merchant/review` | CRD | U | 1 | 审核页以同意/拒绝(U)为主，无前台创建(C) |

### C.2g UGC 监管（无需后台创建）（3）

| 端 | 模块 | 页面路由 | 矩阵缺口 | 已有 | 操作数 | 说明 |
| --- | --- | --- | --- | --- | ---: | --- |
| 平台后台 | 内容 / 社区 / 社区内容 | `/community/list` | C | RUD | 7 | 内容由用户产生，后台只需审/删/改，无需 C |
| 平台后台 | 内容 / 社区 / 社区评论 | `/community/reply` | C | RUD | 4 | 内容由用户产生，后台只需审/删/改，无需 C |
| 平台后台 | 用户 / 用户反馈 / 反馈列表 | `/feedback/list` | C | RUD | 4 | 内容由用户产生，后台只需审/删/改，无需 C |

### C.2h 禁删/软删合理（18）

| 端 | 模块 | 页面路由 | 矩阵缺口 | 已有 | 操作数 | 说明 |
| --- | --- | --- | --- | --- | ---: | --- |
| 商户后台 | 营销 / 拼团 / 拼团活动列表 | `/marketing/combination/combination_list` | UD | CR | 4 | 路由已有 RU；缺 D 多为业务禁删 |
| 商户后台 | 财务 / 发票管理 | `/order/invoice` | D | CRU | 6 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |
| 商户后台 | 订单 / 商品评价 | `/product/reviews` | D | CRU | 6 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |
| 商户后台 | 用户 / 用户管理 | `/user/list` | D | CRU | 6 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |
| 平台后台 | 财务 / 用户结算 / 提现管理 | `/accounts/extract` | D | CRU | 6 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |
| 平台后台 | 财务 / 发票管理 / 发票列表 | `/accounts/receipt` | CUD | R | 2 | 路由已有 CRU；缺 D 多为业务禁删 |
| 平台后台 | 财务 / 商圈代理 | `/accounts/zoneAgent` | D | CRU | 5 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |
| 平台后台 | 财务 / 商圈代理 / 申请结算 | `/accounts/zoneAgent/settlementApply` | D | CRU | 4 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |
| 平台后台 | 营销 / 积分 / 商品分类 | `/marketing/integral/classify` | D | CRU | 12 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |
| 平台后台 | 营销 / 积分 / 商品列表 | `/marketing/integral/proList` | D | CRU | 7 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |
| 平台后台 | 营销 / 预售 | `/marketing/presell` | D | CRU | 7 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |
| 平台后台 | 财务 / 店铺结算 / 分账管理 | `/merchant/applyList` | CUD | R | 5 | 路由已有 CRU；缺 D 多为业务禁删 |
| 平台后台 | 店铺 / 店铺管理 / 店铺分账申请 | `/merchant/applyments` | D | CRU | 4 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |
| 平台后台 | 订单 / 订单列表 | `/order/list` | CD | RU | 30 | 路由已有 CRU；缺 D 多为业务禁删 |
| 平台后台 | 商品 / 商品管理 | `/product/examine` | D | CRU | 48 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |
| 平台后台 | 设置 / 系统设置 / 一号通 / 登陆入口 | `/setting/sms/sms_config/index` | D | CRU | 17 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |
| 平台后台 | 用户 / 用户列表 | `/user/list` | D | CRU | 28 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |
| 平台后台 | 用户 / 付费会员 / 会员权益 | `/user/member/equity` | D | CRU | 5 | 业务常禁止物理删除（评价/发票/提现/商品审核等） |

### C.2i 看板/大屏（1）

| 端 | 模块 | 页面路由 | 矩阵缺口 | 已有 | 操作数 | 说明 |
| --- | --- | --- | --- | --- | ---: | --- |
| 平台后台 | 首页 / 数据大屏 | `/data-screen/index` | CRUD | — | 1 | 看板/大屏以查询展示为主 |

### C.2j 特殊形态页（1）

| 端 | 模块 | 页面路由 | 矩阵缺口 | 已有 | 操作数 | 说明 |
| --- | --- | --- | --- | --- | ---: | --- |
| 平台后台 | 应用 / 公众号 / 微信菜单 | `/app/wechat/menus` | RD | CU | 2 | 微信菜单为发布式 CU，无传统列表删除语义 |

### C.3 原「仍待核对」15 页 — 已二次结案（非缺口）

| 端 | 页面路由 | 结案 | 真实能力 |
| --- | --- | --- | --- |
| 商户后台 | `/marketing/assist/assist_set` | 假缺口 | `merchant/marketing.php`：lst/create/detail/sort/labels/preview |
| 商户后台 | `/devise/diy/product_category` | 假缺口 | DIY：info/create/delete/copy |
| 商户后台 | `/order/reservation` | 非缺口-日历 | 仅预约日历 `ReservationService/list`（非传统 CRUD 表） |
| 商户后台 | `/setting/printer` | 非缺口-模块壳 | 子页 `_path=/setting/printer/list` 与 `/content` |
| 商户后台 | `/accounts/payType` | 假缺口 | 收款方式 form + account 保存 + 转账 lst/detail |
| 平台后台 | `/product/merSpecs` | 非缺口-监管只读 | 商户参数模板列表/详情/删除（平台侧监管） |
| 平台后台 | `/app/routine/download` | 非缺口-工具 | 小程序模板下载/配置，非实体 CRUD |
| 平台后台 | `/user/member` | 非缺口-模块壳 | 子页 type/list/equity/record/interests |
| 平台后台 | `/setting/system/maintain/auth` | 非缺口-特殊 | 对应 `_path=/maintain/auth` 去版权/授权 |
| 平台后台 | `/maintain/cache` | 非缺口-运维 | 清缓存/换域名等运维动作（非业务 CRUD） |
| 平台后台 | `/setting/theme_style` | 非缺口-配置 | 一键换色 get/save |
| 平台后台 | `/setting/product_category` | 非缺口-DIY | 商品分类装修 info/create |
| 平台后台 | `/setting/fab` | 假缺口 | 悬浮按钮 create/delete/scope |
| 平台后台 | `/service/balance_record` | 非缺口-只读 | 一号通购买/使用记录 |
| 平台后台 | `/accounts/zoneAgent/settlementAccount` | 假缺口 | 结算方式 get/post + 结算记录 list/detail |

**CRUD 不全 169：已全部归入非缺口/假缺口，待核对 = 0。**

