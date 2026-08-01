# 功能地图速查

完整清单见 `docs/feature-matrix.md`；原图在 `docs/assets/`。

## P0 闭环

商户入驻审核 → 商品/SKU/库存 → 多商户购物车拆单 → 支付 → 发货/核销 → 售后 → 商户结算提现 → 平台/商户 RBAC

## P1

优惠券、秒杀、拼团、积分、分销、会员、DIY、素材、客服、统计

## P2

预售、砍价助力、卡密/虚拟、预约服务、同城配送员、社区种草、直播、电子面单/打印机

## 端

| 端 | 要点 |
| --- | --- |
| 用户端 | 首页/分类/商品/购物车/订单/店铺/营销/分销/积分/客服 |
| 平台后台 | 商户/商品监管/订单/用户/营销/财务/配送/DIY/系统 |
| 商户后台 | 商品/订单/营销/员工客服/财务/店铺设置 |
| 扩展角色 | 配送员、服务人员、客服坐席 |

## 领域模块

`identity` `merchant` `catalog` `cart` `trade` `aftersale` `finance`  
→ `promotion` `loyalty` `distribution` `fulfillment` `cs` `content` `live` …
