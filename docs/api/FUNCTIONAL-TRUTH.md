# 源码真实功能对照（高风险域）

对照源码：`~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0`。与接口清单交叉验证后的**功能真相**，防止按死路由或错端点开发。

## 1. 用户下单（务必分清三条线）

| 业务 | 真实路径 | 控制器 | 文档置信度 |
| --- | --- | --- | --- |
| 普通订单核对 | `POST /api/v2/order/check` | `StoreOrder::v2CheckOrder` | high |
| 普通订单创建 | `POST /api/v2/order/create` | `StoreOrder::v2CreateOrder` | high |
| 普通订单支付 | `POST /api/order/pay/:id` | `StoreOrder::groupOrderPay` | high |
| 支付回调（本仓库） | `POST /api/callback/v1/pay/{mock\|wechat\|alipay}` | 沙箱 HMAC → paySuccess；非真实 SDK | high |
| 积分订单核对 | `POST /api/order/v3/check` | `PointsOrder::beforCheck` | high |
| 积分订单创建 | `POST /api/order/v3/create` | `PointsOrder::createOrder` | high |
| 积分订单再支付 | `POST /api/order/points/pay/:id` | `groupOrderPay` + `is_points=1` | high |
| ~~旧版创建~~ | `POST /api/order/create` | 方法不存在 | **stale** |
| ~~旧版核对~~ | `POST /api/order/check` | 方法不存在 | **stale** |
| ~~StoreOrder v3~~ | `POST /api/v3/order/*` | 方法不存在 | **stale** |
| ~~积分独立 pay~~ | `POST /api/order/v3/pay/:id` | `orderPay` 不存在 | **stale** |

**注意路径差异：**

- `/api/v3/order/*` ≠ `/api/order/v3/*`
- 前者是死路由；后者是**积分商城订单**（`activity_type=20`）。

## 2. 购物车 / 地址（参数曾缺失，已按源码补全）

| 接口 | 关键请求字段（源码 checkParams） |
| --- | --- |
| `POST /api/user/cart/create` | `product_id`, `product_attr_unique`, `cart_num`, `is_new`, `product_type`, `group_buying_id`, `spread_id`, `reservation_id`, `reservation_date` |
| `POST /api/user/cart/check/:id` | 方法名是 `StoreCart::check`；路由文件误写 `checkCerate`（拼写错误） |
| `POST /api/user/address/create` | `real_name`, `phone`, `area[]`, `detail`, `post_code`, `is_default`；带 `address_id` 则编辑 |

## 3. 退款

| 端 | 真实能力 | 路径示例 |
| --- | --- | --- |
| 用户 | 申请退款/退货 | `POST /api/refund/apply/:id` |
| 用户 | 退货寄回 | `POST /api/refund/back_goods/:id` |
| 用户 | 平台介入 | `POST /api/refund/platform/:id` |
| 商户 | 主动退款/审核 | `/mer/store/refundorder/*` |
| 平台 | 退款监管 | `/sys/order/refund/*` |
| 店员(manager) | 代退 | `/api/server/:merId/refund/*` |

## 4. 商户开放接口（仅 6 条，已核对）

| 方法 | 路径 | 控制器 |
| --- | --- | --- |
| POST | `/openapi/auth` | `openapi/Auth::auth` |
| GET | `/openapi/order/list` | `openapi/store/StoreOrder::lst` |
| GET | `/openapi/order/detail/:id` | `openapi/store/StoreOrder::detail` |
| POST | `/openapi/product/create` | `openapi/store/StoreProduct::create`（字段=CREATE_PARAMS） |
| GET | `/openapi/product/list` | `openapi/store/StoreProduct::lst` |
| GET | `/openapi/product/detail/:id` | `openapi/store/StoreProduct::detail` |

## 5. 支付方式与回调（源码）

`StoreOrderRepository::PAY_TYPE`：

`balance` · `weixin` · `routine` · `h5` · `alipay` · `alipayQr` · `weixinQr` · `offline` · `weixinBarCode` · `alipayBarCode`

支付成功统一进入 `StoreOrderRepository::paySuccess`（`group_order.paid=1`，拆商户子单、记账、分账等）。

| 回调路由（前缀 `/api/`） | 控制器 | 渠道 |
| --- | --- | --- |
| `ANY /api/notice/:type` | `Common::notify` | `routine` 小程序 / `wechat` 公众号 / `app` / `partner` 服务商 / `combine` 合单 |
| `ANY /api/notice/mchNotify/:type` | `Common::mchNotify` | 微信商家转账/收付通转账回调 |
| `ANY /api/notice/pay/alipay` | `Common::alipayNotify` | 支付宝 `AlipayService::handleNotify` |
| `ANY /api/notice/callback` | `Common::deliveryNotify` | 同城配送 |
| `ANY /api/order_call_back` | `Common::callBackNotify` | 商家寄件（serve_token 解密） |

余额/线下支付不走第三方回调，在下单/支付接口内直接 `paySuccess` 或改状态。

## 6. 退款单状态机

源码注释与列表统计（`StoreRefundOrderRepository`）为准：

| status | 含义 | 说明 |
| ---: | --- | --- |
| 0 | 待审核 | `REFUND_STATUS_WAIT` |
| 1 | 待退货 | 仅退款同意后也可能直接到已退款（视类型） |
| 2 | 待收货 | 用户已寄回 |
| 3 | 已退款 | 终态成功 |
| 4 | 平台介入 | `REFUND_PLATFORM_INTERVENE` |
| -1 | 审核拒绝 | |
| -2 | 用户取消 | `REFUND_STATUS_CANCEL`（注释曾写 -10，以常量为准） |

**源码陷阱：** 文件内 `const REFUND_STATUS_SUCCESS = 1` 与「3=已退款」统计不一致，**重建勿照抄该常量名语义**，以列表统计与业务注释 `0/-1/1/2/3/4/-2` 为准。

主路径：

```text
申请(0) → 商户同意 → (仅退款)直接退款成功(3)
                 → (退货退款)待退货(1) → 用户寄回 → 待收货(2) → 确认收货退款(3)
       → 拒绝(-1) / 用户取消(-2) / 平台介入(4)
```

## 7. 商户提现 / 结算

| 字段 | 值 | 含义 |
| --- | --- | --- |
| `type` | 0 | 商户余额提现 |
| `type` | 1 | 保证金退回 |
| `status` | 0 | 待平台审核 |
| `status` | 1 | 审核通过 |
| `status` | -1 | 审核拒绝（提现拒绝会 `addMoney` 退回余额） |
| `financial_status` | 0 | 未转账/未打款 |
| `financial_status` | 1 | 已打款（上传凭证 `Financial::update` 设为 1；保证金同意时也会置 1） |
| `financial_type` | 1/2/3 | 银行卡 / 微信 / 支付宝（账户信息） |

主路径（提现 type=0）：

```text
商户申请(saveApply): status=0, financial_status=0, 扣 mer_money
  → 平台审核同意: status=1
  → 平台上传打款凭证(update image): financial_status=1
  → 或审核拒绝: status=-1, 退回 mer_money
```

相关 API：`/mer/financial/*`（申请）、`/sys` accounts Financial（审核/凭证）。菜单「店铺结算」`/mer/accounts` 为入口目录。


## 8. 营销计价规则（`StoreOrderCreateRepository`，普通单 v2）

源码：`app/common/repositories/store/order/StoreOrderCreateRepository.php`。

### 8.1 购物车组合约束（下单前校验）

| 规则 | 行为 |
| --- | --- |
| 多 `product_type` 混加 | 拒绝：「存在多类型商品」 |
| 活动商品（`product_type≠0` 且 ≠10） | 必须单独购买；且仅单商户 |
| 套餐商品 `product_type=10` | 可与规则内商品共存（活动单独购例外） |
| 虚拟商品（`product.type≠0`） | 必须单独购买 |
| 自定义表单商品（`mer_form_id`） | 必须单独购买 |
| 多商户 | 仅当 `sys_switch_combine=1` 允许合并支付，否则拒绝 |
| 配送方式不统一 | 同店多商品配送方式不一致则拒绝分开下单 |
| 限购 | `pay_limit`：单次/累计限购、起购数 |

### 8.2 优惠券 / SVIP / 积分

| 规则 | 行为 |
| --- | --- |
| 店铺券可用活动 | **仅普通商品与预售**（`order_type==0` 或 `2`）；秒杀/拼团/助力等清空店铺券 |
| 平台券 | `$enabledPlatformCoupon = !$order_type` → **仅普通订单**可用平台券 |
| SVIP 价与券 | 商户配置 `svip_coupon_merge`：≠`1` 且用了 SVIP 价时，**禁用店铺券** |
| 积分抵扣 | 依赖 `mer_integral_status` / `mer_integral_rate`；且 `product_type==0`、应付 > 0、用户有积分 |
| 抵扣顺序（金额构成） | 商品价 → SVIP 折扣 → 店铺券 → 平台券 → 积分抵扣 → 运费规则 |

配置字段（商户）：`mer_integral_status`, `mer_integral_rate`, `mer_store_stock`, `svip_coupon_merge`, `enable_assigned`, `enable_tostore_assigned`。

### 8.3 重建注意

- 计价以 **v2 check/create** 为准，不要实现死路由上的旧逻辑。
- 积分商城单走 `PointsOrderCreateRepository`（`/api/order/v3/*`），规则不同。
- 互斥细节若与前端展示不一致，以 Repository 抛错文案为准。

## 9. 校验统计

- API：high=1999 stale=26 unresolved=2 total=2027
- schema：165 表与 install SQL 一致
- gaps：空页 53 + CRUD 矩阵 169 均已结案；文档层真缺口 0  
- **功能基线已锁定**（2026-07-21）

## 10. 仍不能 100% 的部分

1. 响应 `data` 结构需看 Repository 返回值。
2. form-builder 动态表单字段未全部展开。
3. 营销计价主规则见 §8；活动价库存细则仍可继续深挖 Seckill/Presell Repository。
4. 活动价/库存细则仍可深挖 Seckill/Presell Repository（不阻塞基线）。

详见 [ACCURACY.md](./ACCURACY.md)、[VALIDATION-REPORT.md](./VALIDATION-REPORT.md)、[../features/08-gaps.md](../features/08-gaps.md)。
