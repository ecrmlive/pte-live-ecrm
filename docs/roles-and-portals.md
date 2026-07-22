# 角色与端入口

服务容器 / 端口 / IP 权威表：[`docs/release/SERVICE-MATRIX.md`](./release/SERVICE-MATRIX.md)。

## 角色矩阵

| 角色 | 入口 | 数据范围 | 对应 release（local 端口） |
| --- | --- | --- | --- |
| 平台超管/运营 | 平台后台（Vben） | 全平台 | `qixi-mergers-admin` dist → 宿主机 Nginx `:18081` |
| 商户主账号/店员 | 商户后台（Vben）+ 手机端 | 本 `merchant_id` | `qixi-mergers-merchant-admin` → `:18082` |
| C 端用户 | uni-app x / H5 / 小程序 | 本人订单与资产 | `qixi-mergers-h5` → `:18083` |
| C 端用户（PC） | Vue3 PC 商城 | 本人订单与资产 | `qixi-mergers-pc` → `:18086` |
| 分销员 | C 端分销中心 | 本人团队与佣金 | 同 H5 / PC |
| 客服坐席 | 客服工作台 (P1) | 授权店铺会话 | `qixi-mergers-service-web` → `:18084` |
| 配送员 / 服务人员 | 配送/服务端 (P2) | 分配单据 | 后挂 `/api/manager` 等 |

后台 API：`qixi-mergers-api-admin` → `:18080`（`.20`）— platform / merchant / manager / service / open。  
C 端 API：`qixi-mergers-api-app` → `:18085`（`.22`）— app / callback。  
前端入口：**仅宿主机 Nginx**；后台站反代 api-admin，H5/PC 反代 api-app。

## 端与 API 前缀

| 端 | 进程 | API 前缀 | 鉴权 |
| --- | --- | --- | --- |
| 平台后台 | api-admin | `/api/platform/v1` | 平台 JWT + RBAC |
| 商户后台 | api-admin | `/api/merchant/v1` | 商户 JWT + `mer_id` |
| 店员 manager (P1) | api-admin | `/api/manager/v1` | 商户员工 JWT |
| 客服 (P1) | api-admin | `/api/service/v1` | 客服 JWT |
| 开放接口 | api-admin | `/api/open/v1` | AppKey / 签名 |
| C 端 | api-app | `/api/app/v1` | 用户 JWT |
| 回调 | api-app | `/api/callback/v1` | 支付平台验签 |

## 页面信息架构（用户端）

1. 首页 / 分类 / 购物车 / 用户  
2. 商品详情 → 规格 → 下单  
3. 订单列表/详情/售后  
4. 店铺街 → 店铺首页  
5. 营销活动页（秒杀/拼团/砍价/积分）  
6. 分销中心 / 余额 / 优惠券  
7. 客服会话  
8. 商户入驻申请  

## 管理端信息架构

### 平台

数据概览 → 商户 → 商品 → 订单 → 用户 → 营销 → 财务 → 配送 → 内容/DIY → 设置/权限

### 商户

店铺概览 → 商品 → 订单 → 营销 → 员工/客服 → 财务 → 店铺设置

## 权限原则

- 平台菜单与商户菜单两套 RBAC，不共用角色表主键空间（可同库不同表或 `scope` 字段）。
- 商户接口中间件强制注入并校验 `merchant_id`。
- 资金、退款、结算类操作写操作日志，关键操作可二次确认。
