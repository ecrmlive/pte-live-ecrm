# 本仓库 OpenAPI 契约

与 CRMEB 对照文档（`docs/api/`）分离。目标服务边界固定为 `api-platform`、`api-business`、`api-merchant`；仍引用 `api-admin` 或 `api-app` 的内容仅是待清理的旧迁移参考，不能作为运行或验收契约。

| 文件 | 说明 |
| --- | --- |
| `platform-auth.yaml` | 平台后台身份（阶段 1） |
| `merchant-auth.yaml` | 商户后台身份（阶段 1） |
| `app-auth.yaml` | C 端用户身份（阶段 1） |
| `health.yaml` | api-platform / api-business / api-merchant 健康检查 |
| `platform-merchant-catalog.yaml` | 平台商户审核/类目品牌/商品审核（阶段 2） |
| `platform-business-zones.yaml` | 平台区域商圈、代理人员与代理审核（建设中） |
| `merchant-catalog.yaml` | 商户商品 CRUD / 库存 / 上下架（阶段 2） |
| `app-catalog.yaml` | C 端可售商品只读（阶段 2） |
| `app-trade.yaml` | C 端购物车/地址/`v2/order`下单/支付（阶段 3） |
| `app-callback-pay.yaml` | C 端微信支付 v3 / Mock 支付回调（阶段 3）；支付宝回调待接入 |
| `merchant-trade.yaml` | 商户订单发货核销（阶段 3；权威） |
| `platform-trade.yaml` | 平台订单监管只读（阶段 3；权威） |
| `merchant-order.yaml` | 同 `merchant-trade` 别名，勿再扩写 |
| `platform-order.yaml` | 同 `platform-trade` 别名，勿再扩写 |
| `app-aftersale.yaml` | C 端仅退款（阶段 4） |
| `merchant-aftersale-finance.yaml` | 商户售后/提现（阶段 4） |
| `platform-aftersale-finance.yaml` | 平台退款/提现审核（阶段 4） |
| `app-promotion.yaml` | C 端领券选券/分销（阶段 5） |
| `merchant-promotion.yaml` | 商户店铺券（阶段 5） |
| `platform-promotion.yaml` | 平台券/分销监管（阶段 5） |
| `app-loyalty.yaml` | C 端积分商城 `/order/v3` + 抵扣（阶段 6） |
| `platform-content.yaml` | 平台公告 CRUD + 协议规则（阶段 6a/7） |
| `app-content.yaml` | C 端公告/协议只读（阶段 7） |
| `app-customer-service.yaml` | C 端客服会话与 IM 凭据（阶段 6b/7） |
| `platform-diy-seckill.yaml` | 平台 DIY + 秒杀监管（阶段 6a/活动） |
| `merchant-seckill.yaml` | 商户秒杀 CRUD（阶段 6 活动） |
| `app-diy-seckill.yaml` | C 端 DIY 首页 + 秒杀列表/价覆盖（阶段 6） |
| `app-combination.yaml` | C 端拼团列表/开团参团下单（阶段 6） |
| `platform-combination.yaml` | 平台拼团监管（阶段 6） |
| `merchant-combination.yaml` | 商户拼团 CRUD（阶段 6） |
| `platform-svip.yaml` | 平台付费会员授予（阶段 6） |
| `merchant-svip.yaml` | 商户 SVIP 与店铺券叠加配置（阶段 6） |
| `app-presell.yaml` | C 端预售（全款+定金尾款，阶段 6） |
| `platform-presell.yaml` | 平台预售监管（阶段 6） |
| `merchant-presell.yaml` | 商户预售 CRUD（阶段 6） |
| `app-live.yaml` | C 端直播间（阶段 6e） |
| `platform-live.yaml` | 平台直播监管（阶段 6e） |
| `merchant-live.yaml` | 商户直播间 CRUD（阶段 6e） |
| `app-community.yaml` | C 端社区种草（阶段 6e） |
| `platform-community.yaml` | 平台社区审核（阶段 6e） |
| `merchant-community.yaml` | 商户本店种草 CRUD（阶段 6e/7；`community/create|update|delete`） |
| `app-assist.yaml` | C 端助力发起/助力/下单（阶段 6） |
| `merchant-assist.yaml` | 商户助力 CRUD（阶段 6） |
| `platform-assist.yaml` | 平台助力监管（阶段 6） |
| `service-desk.yaml` | 客服查单 + 按 mer_id 读快捷回复 `/api/service/v1`（阶段 6b/7） |
| `platform-customer-service.yaml` | 统一后台客服队列、数据范围、领取与可追溯转接；IM 消息仍由 pte-live-im 提供 |
| `merchant-service-reply.yaml` | 商户快捷回复 CRUD（阶段 7；`reply/write`） |
| `merchant-open.yaml` | 商户开放接口 6 条 `/api/open/v1`（阶段 6c） |
| `manager-trade.yaml` | 店员端登录/核销/代退/发货 `/api/manager/v1`（阶段 6d） |
| `app-svip.yaml` | C 端 SVIP 计价互斥（阶段 6） |
| `app-reservation.yaml` | C 端预约列表/下单（阶段 6e） |
| `merchant-reservation.yaml` | 商户预约配置（阶段 6e） |
| `platform-setting.yaml` | 平台管理员/角色/菜单最小 CRUD（阶段 7） |
| `merchant-setting.yaml` | 商户店铺/店员/子账号/角色（阶段 7） |
| `merchant-diy.yaml` | 商户装修页与默认页配置（阶段 6a） |
| `merchant-fulfillment.yaml` | 商户配送人员与发票审核（阶段 7） |
| `merchant-logistics.yaml` | 商户快递/城市查询与运费模板（阶段 7） |
| `merchant-product-meta.yaml` | 商户商品标签/保障/参数模板（阶段 7） |
| `platform-attachment.yaml` | 平台素材库分类/上传（阶段 7） |
| `merchant-attachment.yaml` | 商户素材库（mer_id 隔离，阶段 7） |
| `platform-logistics.yaml` | 平台快递公司与城市（阶段 7） |
| `platform-product-meta.yaml` | 平台商品标签/保障监管（阶段 7） |
| `platform-user-segmentation.yaml` | 平台用户标签、分组与打标（阶段 7） |
