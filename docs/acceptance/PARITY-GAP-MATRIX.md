# CRMEB 双后台对照缺口矩阵

> 自动生成：`scripts/generate-parity-gap-matrix.py` · 2026-08-05
>
> 状态：`missing` 无页面 · `shell` SQL 有菜单但 registry 无组件 · `partial` 已注册待按钮级闭环 · `done` 需人工在验收台账标关闭后回写。
>
> 基线路径来自 `docs/generated/features-0{1,2}-*.tsv` 的 `page_route`；完成以 TSV 操作行 + `crmeb-vben-parity.md` 六条为准。

## 平台后台（platform）

| 指标 | 数量 |
| --- | ---: |
| 基线操作数 | 1333 |
| 特征表页面路由 | 237 |
| SQL 叶子 page | 97 |
| registry 路径 | 288 |
| 特征表 ops（已挂 registry） | 1333 / 1333 |
| missing | 0 |
| shell | 0 |
| partial | 288 |

### 按模块状态计数

| 模块 | missing | shell | partial |
| --- | ---: | ---: | ---: |
| accounts | 0 | 0 | 4 |
| content | 0 | 0 | 4 |
| marketing | 0 | 0 | 7 |
| merchant | 0 | 0 | 4 |
| operations | 0 | 0 | 1 |
| product | 0 | 0 | 2 |
| setting | 0 | 0 | 10 |
| statistic | 0 | 0 | 1 |
| user | 0 | 0 | 16 |
| 其他 | 0 | 0 | 2 |
| 内容 | 0 | 0 | 10 |
| 分销 | 0 | 0 | 13 |
| 商品 | 0 | 0 | 14 |
| 客服 | 0 | 0 | 4 |
| 应用 | 0 | 0 | 12 |
| 店铺 | 0 | 0 | 23 |
| 用户 | 0 | 0 | 21 |
| 维护 | 0 | 0 | 12 |
| 营销 | 0 | 0 | 40 |
| 装修 | 0 | 0 | 18 |
| 订单 | 0 | 0 | 4 |
| 设置 | 0 | 0 | 40 |
| 财务 | 0 | 0 | 20 |
| 首页 | 0 | 0 | 6 |

### missing（优先补齐）

_无_

### shell（有菜单无组件）

_无_

### partial（已注册，待按钮级 / 布局闭环）

- `/` · 2 ops
- `/accounts` · 5 ops
- `/accounts/accounts` · 1 ops
- `/accounts/bill` · 2 ops
- `/accounts/capital` · 4 ops
- `/accounts/capitalFlow` · 5 ops
- `/accounts/extract` · 6 ops
- `/accounts/invoiceDesc` · 1 ops
- `/accounts/invoices`
- `/accounts/merchant-settlement`
- `/accounts/merchantBill` · 10 ops
- `/accounts/receipt` · 2 ops
- `/accounts/record` · 1 ops
- `/accounts/settings` · 2 ops
- `/accounts/statement` · 6 ops
- `/accounts/transferRecord` · 9 ops
- `/accounts/user-assets`
- `/accounts/withdraw`
- `/accounts/zoneAgent` · 5 ops
- `/accounts/zoneAgent/commissionRecord` · 2 ops
- `/accounts/zoneAgent/settlementAccount` · 2 ops
- `/accounts/zoneAgent/settlementApply` · 4 ops
- `/accounts/zoneAgent/settlementReview` · 3 ops
- `/admin/app/wechat/reply` · 9 ops
- `/app/routine` · 2 ops
- `/app/routine/download` · 1 ops
- `/app/routine/template` · 6 ops
- `/app/wechat` · 9 ops
- `/app/wechat/file` · 4 ops
- `/app/wechat/menus` · 2 ops
- `/app/wechat/newsCategory` · 7 ops
- `/app/wechat/reply/follow/subscribe` · 1 ops
- `/app/wechat/reply/index/default` · 1 ops
- `/app/wechat/reply/keyword` · 1 ops
- `/app/wechat/template` · 6 ops
- `/app_config` · 1 ops
- `/apploction` · 1 ops
- `/assist` · 1 ops
- `/banlace` · 1 ops
- `/brokerage`
- `/business-zones/agent-review` · 1 ops
- `/business-zones/agents` · 6 ops
- `/business-zones/index` · 7 ops
- `/business-zones/manage` · 1 ops
- `/business-zones/settings` · 1 ops
- `/cms` · 5 ops
- `/cms/article` · 8 ops
- `/cms/articleCategory` · 6 ops
- `/community` · 9 ops
- `/community/category` · 6 ops
- `/community/list` · 7 ops
- `/community/reply` · 4 ops
- `/community/topic` · 9 ops
- `/config/classify` · 5 ops
- `/config/picture` · 12 ops
- `/config/setting` · 5 ops
- `/content` · 1 ops
- `/content/attachment`
- `/content/community`
- `/content/diy`
- `/content/notice`
- `/dashboard` · 13 ops
- `/data-screen/index` · 1 ops
- `/delivery` · 3 ops
- `/delivery/recharge_record` · 3 ops
- `/delivery/usage_record` · 2 ops
- `/delivery_config` · 1 ops
- `/feedback` · 5 ops
- `/feedback/classify` · 5 ops
- `/feedback/list` · 4 ops
- `/freight/city/list` · 1 ops
- `/freight/express` · 10 ops
- `/group/config/67` · 8 ops
- `/group/config/68` · 8 ops
- `/group/config/69` · 8 ops
- `/group/config/75` · 8 ops
- `/group/config/76` · 8 ops
- `/group/exportList` · 2 ops
- `/group/list` · 11 ops
- `/group/topic/94` · 8 ops
- `/maintain` · 2 ops
- `/maintain/cache` · 2 ops
- `/maintain/dataBackup` · 8 ops
- `/marketing` · 2 ops
- `/marketing/application/list` · 10 ops
- `/marketing/assist`
- `/marketing/assist/goods_list` · 21 ops
- `/marketing/assist/list` · 6 ops
- `/marketing/atmosphere/list` · 8 ops
- `/marketing/border/list` · 8 ops
- `/marketing/broadcast`
- `/marketing/broadcast/list` · 6 ops
- `/marketing/combination` · 2 ops
- `/marketing/combination/combination_goods` · 24 ops
- `/marketing/combination/combination_list` · 6 ops
- `/marketing/combination/combination_set` · 1 ops
- `/marketing/coupon` · 2 ops
- `/marketing/coupon/list` · 3 ops
- `/marketing/coupon/receipt-records`
- `/marketing/coupon/send-records`
- … 另有 188 条 registry 路径

## 商户后台（merchant）

| 指标 | 数量 |
| --- | ---: |
| 基线操作数 | 615 |
| 特征表页面路由 | 100 |
| SQL 叶子 page | 49 |
| registry 路径 | 116 |
| 特征表 ops（已挂 registry） | 615 / 615 |
| missing | 0 |
| shell | 0 |
| partial | 116 |

### 按模块状态计数

| 模块 | missing | shell | partial |
| --- | ---: | ---: | ---: |
| devise | 0 | 0 | 1 |
| finance | 0 | 0 | 3 |
| marketing | 0 | 0 | 3 |
| order | 0 | 0 | 1 |
| setting | 0 | 0 | 8 |
| 公告 | 0 | 0 | 1 |
| 员工 | 0 | 0 | 10 |
| 商品 | 0 | 0 | 9 |
| 用户 | 0 | 0 | 6 |
| 营销 | 0 | 0 | 25 |
| 装修 | 0 | 0 | 5 |
| 订单 | 0 | 0 | 8 |
| 设置 | 0 | 0 | 25 |
| 财务 | 0 | 0 | 8 |
| 首页 | 0 | 0 | 3 |

### missing（优先补齐）

_无_

### shell（有菜单无组件）

_无_

### partial（已注册，待按钮级 / 布局闭环）

- `/` · 5 ops
- `/accounts` · 8 ops
- `/accounts/capitalFlow` · 5 ops
- `/accounts/payType` · 1 ops
- `/accounts/statement` · 6 ops
- `/accounts/transManagement` · 9 ops
- `/assist` · 1 ops
- `/city` · 1 ops
- `/community/list` · 7 ops
- `/config` · 7 ops
- `/config/freight/express` · 3 ops
- `/config/freight/shippingTemplates` · 5 ops
- `/config/guarantee` · 7 ops
- `/config/picture` · 12 ops
- `/config/service` · 10 ops
- `/config/service_staff` · 9 ops
- `/config/service_statistic` · 5 ops
- `/dashboard` · 15 ops
- `/delivery` · 3 ops
- `/delivery/delivery_point` · 9 ops
- `/delivery/delivice_statistic` · 1 ops
- `/delivery/personnel_manage` · 7 ops
- `/delivery/personnel_manage/index` · 7 ops
- `/delivery/recharge_record` · 1 ops
- `/delivery/store_manage` · 8 ops
- `/delivery/usage_record` · 3 ops
- `/devise/` · 12 ops
- `/devise/diy/index`
- `/devise/diy/list` · 13 ops
- `/devise/diy/product_category` · 2 ops
- `/finance/balance`
- `/finance/settlement`
- `/finance/withdraw`
- `/group/topic/95` · 8 ops
- `/marketing` · 4 ops
- `/marketing/assist`
- `/marketing/assist/assist_set` · 2 ops
- `/marketing/assist/list` · 11 ops
- `/marketing/broadcast`
- `/marketing/broadcast/list` · 1 ops
- `/marketing/combination` · 2 ops
- `/marketing/combination/combination_goods` · 12 ops
- `/marketing/combination/combination_list` · 4 ops
- `/marketing/community`
- `/marketing/coupon` · 4 ops
- `/marketing/coupon/list` · 8 ops
- `/marketing/coupon/send` · 1 ops
- `/marketing/coupon/user` · 1 ops
- `/marketing/discounts/list` · 8 ops
- `/marketing/integral` · 2 ops
- `/marketing/integral/config` · 1 ops
- `/marketing/integral/log` · 4 ops
- `/marketing/presell/list` · 11 ops
- `/marketing/seckill/list` · 1 ops
- `/marketing/seckill/product/list` · 15 ops
- `/marketing/seckill/store_seckill/list` · 9 ops
- `/marketing/studio/assistant` · 7 ops
- `/marketing/studio/list` · 26 ops
- `/one_setting` · 1 ops
- `/order` · 2 ops
- `/order/cancellation` · 2 ops
- `/order/customer` · 22 ops
- `/order/invoice` · 6 ops
- `/order/list` · 41 ops
- `/order/refund` · 12 ops
- `/order/reservation` · 1 ops
- `/order/verify`
- `/product` · 2 ops
- `/product/attr` · 4 ops
- `/product/cdkey` · 13 ops
- `/product/classify` · 8 ops
- `/product/label` · 6 ops
- `/product/list` · 41 ops
- `/product/reservation` · 1 ops
- `/product/reviews` · 6 ops
- `/product/specs` · 6 ops
- `/product/unit` · 4 ops
- `/server` · 1 ops
- `/server_manage` · 1 ops
- `/setting` · 1 ops
- `/setting/admins`
- `/setting/attachment`
- `/setting/delivery` · 2 ops
- `/setting/im-sdk-app`
- `/setting/integral-policy`
- `/setting/payment`
- `/setting/printer` · 8 ops
- `/setting/printer/list` · 7 ops
- `/setting/role`
- `/setting/shop`
- `/setting/sms/dumpConfig` · 1 ops
- `/setting/sms/sms_account/index` · 1 ops
- `/setting/sms/sms_config/index` · 5 ops
- `/setting/staff`
- `/setting/systemAdmin` · 6 ops
- `/setting/systemLog` · 4 ops
- `/setting/systemRole` · 10 ops
- `/station/notice` · 5 ops
- `/statistic/order` · 3 ops
- `/statistic/product` · 3 ops
- … 另有 16 条 registry 路径

