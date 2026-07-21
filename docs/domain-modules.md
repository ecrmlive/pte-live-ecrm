# 领域模块拆分

按业务边界拆模块，便于 Go 包或后续微服务演进。命名可映射到 `internal/domain/<name>`。

## 模块一览

| 模块 | 职责 | P0 |
| --- | --- | --- |
| `identity` | 平台管理员、商户管理员、C 端用户、登录、RBAC | ✅ |
| `merchant` | 入驻、审核、店铺资料、类型/分类、关注 | ✅ |
| `catalog` | 平台/店铺分类、品牌、SPU/SKU、库存、评价、标签 | ✅ |
| `cart` | 多商户购物车 | ✅ |
| `trade` | 下单、支付意图、订单状态、发货、核销、发票 | ✅ |
| `aftersale` | 退款/退货状态机 | ✅ |
| `promotion` | 优惠券、秒杀、拼团、预售、助力、满减 | P1 |
| `loyalty` | 积分、会员等级/付费会员 | P1 |
| `distribution` | 分销关系、佣金、提现 | P1 |
| `finance` | 商户账户、平台流水、结算、分账 | ✅ |
| `fulfillment` | 运费模板、快递、同城配送员 | P0 基础运费 / P2 配送员 |
| `cs` | 客服会话、快捷回复 | P1 |
| `content` | 资讯、公告、社区种草、DIY 装修、素材 | P1 DIY / P2 社区 |
| `live` | 直播间与带货 | P2 |
| `service-order` | 预约服务工单、服务人员 | P2 |
| `notify` | 短信、模板消息、站内信（NATS 消费） | ✅ 基础 |
| `analytics` | 经营统计报表 | P1 |

## 关键依赖

```text
identity
   └── merchant
         └── catalog ── promotion / loyalty
                └── cart ── trade ── aftersale
                              ├── finance
                              ├── fulfillment
                              └── distribution
content / diy ── 读写 catalog 展示配置
cs / notify ── 依附 trade / identity
live / service-order / community ── 依附 catalog + trade
```

## 状态机（必须文档化后实现）

### 订单（商户子单）

```text
pending_pay → paid → delivering / awaiting_verify → completed
                 ↘ cancelled
paid / delivering → aftersale_* （与售后单关联）
```

### 商户入驻

```text
draft → pending_review → approved | rejected
approved → enabled | disabled
```

### 商户结算 / 提现

```text
bill_pending → bill_frozen → withdraw_applied → approved → paid
                                      ↘ rejected
```

### 售后

```text
applied → merchant_handling → platform_intervene? → refunding → refunded | rejected
```

## 与 CRMEB Repository 对照（实现时可打开源码）

| 本模块 | 主要对照 |
| --- | --- |
| `trade` | `store/order/*` |
| `catalog` | `store/product/*` |
| `merchant` | `system/merchant/*` |
| `finance` | `system/financial/*`, reconciliation / profitsharing |
| `promotion` | coupon、seckill、presell、group、assist、discounts |
| `distribution` | `user` brokerage / spread |
| `cs` | `store/service/*` |

## 包内建议结构（Go）

```text
internal/domain/<module>/
  model.go          # 实体
  repo.go           # 接口
  service.go        # 用例
  errors.go
internal/infra/persist/<module>/
  *.go              # GORM 实现
```

跨模块调用走 application/service 层，禁止 infra 互相直连。
