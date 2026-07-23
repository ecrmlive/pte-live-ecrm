# 全端开发计划（权威）

基线已锁定（2026-07-21）。本计划把 **产品主链路 × 领域模块 × 六端入口 × release 形态** 收成可执行竖切节奏。  
验收以 [`features/`](./features/README.md) 为准；高风险状态机以 [`api/FUNCTIONAL-TRUTH.md`](./api/FUNCTIONAL-TRUTH.md) 为准；表前缀 `qixi_`；商户域强制 `mer_id`。

配套：[`release/SERVICE-MATRIX.md`](./release/SERVICE-MATRIX.md) · [`architecture-target.md`](./architecture-target.md) · [`domain-modules.md`](./domain-modules.md)

---

## 0. 已消化基线（摘要）

### 0.1 产品一句话

多商户入驻式商城：平台审规则与钱，商户经营与履约，用户购买与售后；资金可追溯。

### 0.2 端与 API（本仓库契约前缀）

| 端 | 源码目录 | 发布 | API |
| --- | --- | --- | --- |
| **后台 API** | `api/cmd/api-admin` | Docker `.20` `:18080` | platform / merchant / manager / service / open |
| **C 端 API** | `api/cmd/api-app` | Docker `.22` `:18085` | app / callback |
| Job | `api/cmd/job` | Docker `.21` | NATS |
| 平台后台 | `admin-platform/`（pack key `admin`） | dist + 宿主机 Nginx `:18081` | → api-admin |
| 商户后台 | `admin-merchant/`（pack key `merchant-admin`） | Nginx `:18082` | → api-admin |
| 用户端 H5 | `app-uni/` | Nginx `:18083` | → **api-app** |
| 用户端 PC | `app-pc/` | Nginx `:18086` | → **api-app** |
| 客服工作台 | `service-web/` | Nginx `:18084`（P1） | → api-admin |

领域共享 `api/internal/domain`。CRMEB 对照在 `docs/api/`；本仓库契约另建 `docs/openapi/`。

### 0.3 必须默写的主链路

```text
入驻审核 → 商品(SPU/SKU)上架
  → 加购(按 mer_id 分堆) → v2 check/create → pay → callback/paySuccess
  → 发货/核销 → 售后状态机 → 商户余额/提现
```

普通单：`POST …/v2/order/check|create` + `…/order/pay/:id`  
积分单：`POST …/order/v3/check|create`（≠ `/api/v3/order/*`）  
支付成功统一等价 `paySuccess`（主单 paid、子单履约、记账/分账钩子）。

### 0.4 工程约束

- Go：**api-admin / api-app 进程分立** + `job`（NATS）；不用 Docker Nginx。
- 配置仅 `.yaml`；本机 pack → rsync；验证/部署仅用户明确要求时做。
- 高风险域：下单、支付回调、退款、库存、券/积分计价、佣金、结算——先状态机后代码。

### 0.5 阶段 0 现状

已完成：release 树、固定网、`local-api-admin` / `local-api-app` 健康检查、前端 dist 骨架、宿主机 Nginx 模板。  
**本计划从阶段 1 起为业务竖切。**

---

## 1. 总原则

1. **竖切主链路**：每阶段三端（或相关端）同时可演示，不做「先堆满平台再写用户端」。
2. **每阶段可 `local-*` 跑通**：API + 必要前端页 + 迁移。
3. **领域内核优先**：`internal/domain/*`；portal 只做适配（JWT、DTO、菜单）。
4. **对照顺序**（功能清单建议）：平台商品/订单/商户/财务/营销 → 商户同序 → 用户登录/地址/车/v2 支付/售后/分销。
5. **P0 里程碑 = 阶段 3**：多商户下单支付闭环（可用 mock 支付）。

---

## 2. 阶段总览

| 阶段 | 周期（约） | 主题 | 里程碑 |
| --- | --- | --- | --- |
| 0 | ✅ | 基建与骨架 | `healthz` / pack / 宿主机 Nginx |
| 1 | 1–1.5 周 | 身份三端 | 平台/商户/C 端能登录进壳 |
| 2 | 2 周 | 商户 + 商品 | 有店有货 C 端可见 |
| **3** | **2–3 周** | **交易主链 P0** | **多店下单→支付→发货** |
| 4 | 2 周 | 售后 + 财务 | 退款闭环 + 提现审核 |
| 5 | 2–3 周 | 营销 / 积分 / 分销 | 计价 §8 + 积分单 + 佣金 |
| 6 | 按需 | 周边 | DIY / 客服 / OpenAPI / 直播 / 预约 / manager |
| **7** | **打磨** | **验收与收口** | **主链路可演示、文档/路由一致、已知缺口清单** |

合计至 P0（阶段 3）约 **5.5–6.5 周**；至财务闭环约 **7.5–8.5 周**；营销完整约 **10–12 周**；阶段 7 为验收打磨（不新开大功能）。

---

## 3. 分阶段明细（全端）

每阶段列出：领域 · API · 各端 UI · 数据 · 验收。  
「对照 features」列指向验收切片，实现时按按钮/路由勾选。

### 阶段 1 — 身份三端

| 维 | 内容 |
| --- | --- |
| 领域 | `identity`：平台管理员、商户管理员、C 端用户；JWT；平台 RBAC + 商户 RBAC |
| API | `POST/GET` 登录、刷新、当前用户、菜单树子集、改密；中间件注入 `admin_id` / `mer_id` / `uid` |
| 平台 `admin` | Vben 登录 + 空壳布局 + 真实菜单树子集（系统/权限入口） |
| 商户 `merchant-admin` | 同上（商户账号）；所有请求带商户上下文 |
| 用户 `app-uni` | 登录/注册页（手机号或账号）；token 持久化；「我的」空壳 |
| 用户 `app-pc` | PC 布局壳 + 登录/注册 + 首页/个人中心壳（token 键 `qixi_pc_*`，与 H5 隔离） |
| 数据 | `qixi_system_admin*`、`qixi_merchant_admin*`、`qixi_user`、角色菜单相关表（从 schema 迁出） |
| 验收 | 三角色登录成功；商户 token 无法打平台接口；平台可看菜单；PC/H5 均可调 C 端登录；`features` 登录/权限相关 R/U 可勾一批 |

**不做**：完整 1295/600 按钮权限一次性导入——先导入主链路菜单子集，阶段 2+ 按模块扩权。

---

### 阶段 2 — 商户 + 商品

| 维 | 内容 |
| --- | --- |
| 领域 | `merchant` + `catalog`（+ 运费模板基础 `fulfillment`） |
| API | 入驻申请/审核；店铺资料；平台类目品牌；SPU/SKU CRUD；库存；上下架；平台商品审核 |
| 平台 | 商户审核列表/详情；商品审核；类目品牌 |
| 商户 | 商品列表/编辑/SKU/库存/上架；运费模板基础 |
| 用户 H5 | 分类、商品列表、商详（只读可售）；店铺首页基础 |
| 用户 PC | 分类广场、商品列表筛选壳、商详、店铺街/店铺页基础（功能表 4） |
| 数据 | `qixi_merchant*`、`qixi_store_product*`、`qixi_store_product_attr*`、分类品牌相关 |
| 验收 | 入驻通过 → 商户建多规格商品 → 平台审核（若需）→ C 端 H5/PC 可见可加购位（加购可阶段 3） |

商品字段语义保持：`type` / `product_type` / `delivery_way` / `spec_type`（见 product-understanding）。

---

### 阶段 3 — 交易主链 P0【里程碑】

| 维 | 内容 |
| --- | --- |
| 领域 | `cart` + `trade` + `notify` 基础 + `job` 支付后续 |
| API | 购物车 CRUD；**v2 check/create**；支付意图；`/api/callback/v1`（及兼容对照的通知路由设计）；订单列表/详情；发货/核销基础；mock 支付直调 `paySuccess` 等价逻辑 |
| 平台 | 订单监管列表/详情（只读为主） |
| 商户 | 待发货列表、发货、核销 |
| 用户 H5/PC | 地址 CRUD；购物车；结算；订单列表/详情；支付结果页 |
| Job | 支付成功后通知、超时关单（可先最小集） |
| 数据 | `qixi_store_cart`、`qixi_store_group_order`、`qixi_store_order*`、地址表 |
| 验收 | **两商户商品一单支付成功**（mock 可）；库存扣减；子单可发货；H5 或 PC 任一端可演示闭环；禁止实现 stale 下单路由 |

支付：`PAY_TYPE` 先支持 `balance` + mock 渠道；微信/支付宝真实回调阶段 3 末或阶段 4 初接 1 条通道即可。

---

### 阶段 4 — 售后 + 财务

| 维 | 内容 |
| --- | --- |
| 领域 | `aftersale` + `finance` |
| API | 用户退款申请/寄回/平台介入；商户审核/主动退；平台监管；提现申请/审核/打款凭证；商户余额与流水 |
| 平台 | 退款单、提现审核、打款凭证 |
| 商户 | 售后处理、余额/流水、提现申请 |
| 用户 | 申请售后、进度、退货物流 |
| 数据 | `qixi_store_refund_order*`、`qixi_financial*`、`qixi_user_bill`、分账表钩子 |
| 验收 | 按 FUNCTIONAL-TRUTH §6–7 走通「仅退款」与「提现拒绝退回余额」；库存/优惠回滚基本正确 |

---

### 阶段 5 — 营销 / 积分 / 分销

| 维 | 内容 |
| --- | --- |
| 领域 | `promotion` + `loyalty` + `distribution` |
| API | 券领取/核销进计价；SVIP/积分抵扣（§8 互斥）；积分商城 `/order/v3` 等价；分销关系与佣金入账/提现 |
| 平台 | 平台券、积分规则、分销配置、活动监管 |
| 商户 | 店铺券、秒杀/拼团等（可先券+一种活动） |
| 用户 | 领券、结算选券、积分商城、分销中心 |
| 验收 | 同一结算路径计价可单测；积分单与普通单入口不混；佣金流水可查 |

活动类（秒杀/拼团/预售/助力）可在本阶段做 **券 + 1 种活动**，其余顺延阶段 6。

---

### 阶段 6 — 周边（按需并行）

| 包 | 端 | 内容 |
| --- | --- | --- |
| 6a DIY/内容 | 平台 + C 端 | 装修页、素材、公告 |
| 6b 客服 | `service-web` + service API | 会话、查单、快捷回复 |
| 6c OpenAPI | open 前缀 | 已核对的 6 条能力 |
| 6d manager | uni 店员端 | 核销/代退/配送基础（features/04） |
| 6e 直播/预约/社区 | 预约✅ 直播 stub✅ 社区✅ | P2 |

---

## 4. 各端「横切」节奏（避免一端空转）

```text
周序（示意）     API/domain     平台 admin      商户 mer-admin     用户 H5/PC
阶段1            identity       登录+壳+菜单    登录+壳            登录+我的（PC 壳已起）
阶段2            merchant/cat   审商户/审商品   商品CRUD           分类/商详
阶段3            cart/trade     订单监管        发货核销           车/结算/支付
阶段4            after/fin      退款/提现审     售后/提现          售后进度
阶段5            promo/loy/dis  规则/监管       券/活动            券/积分/分销
阶段6            周边           DIY/…           …                  … + service-web
```

客服 / OpenAPI / manager **不阻塞 P0**；阶段 3 前不排主力。

---

## 5. 仓库落地约定（与阶段同步）

```text
api/
  cmd/api-admin, cmd/api-app, cmd/job
  internal/domain/{identity,merchant,catalog,cart,trade,...}
  internal/{platform,merchant,app,open,callback,manager,service}/
docs/openapi/
sql/00N_*.sql
admin-platform/ · admin-merchant/ · app-uni/ · app-pc/ · service-web/
```

每阶段结束更新：

- 本文件勾选表（或 `docs/generated/` 进度备注）
- `docs/openapi/` 对应 YAML
- features 相关行的实现备注（可选）

---

## 6. 风险与明确不做（阶段内）

| 项 | 处理 |
| --- | --- |
| 死路由 `/api/v3/order/*`、旧 create/check | **永不实现** |
| 响应 `data` 全量静态建模 | 不做；对照 Repository 按接口落地 |
| 一次导入全部按钮权限 | 按模块扩；阶段 1 子集 |
| Docker Nginx | 已否决；只用宿主机 `opts/nginx` |
| 后台与 C 端塞同一 API 进程 | 已否决；必须 api-admin / api-app 分立 |
| 按每个小域再拆微服务 | 不做；同仓共享 domain |
| 服务器源码构建 | 禁止 |

---

## 7. 阶段 1 进度

| 项 | 状态 |
| --- | --- |
| `sql/001_identity.sql` 管理员/角色菜单/`qixi_user` 最小集 + 种子 | ✅ |
| `internal/domain/identity` + platform/merchant login/me/menus/password | ✅ |
| JWT 中间件与 `mer_id` / `uid` 强制校验 | ✅ |
| 双后台：登录页 + 布局 + 菜单 API（Vue3/Ant Design Vue 壳） | ✅ |
| `docs/openapi/platform-auth.yaml` · `merchant-auth.yaml` · `app-auth.yaml` | ✅ |
| api-app：`/api/app/v1/auth/{login,register,refresh,me}` | ✅ |
| uni-app：登录/注册页 +「我的」空壳 + token 持久化 | ✅ |

阶段 1 完成。

### 阶段 2 进度

| 项 | 状态 |
| --- | --- |
| `sql/002_platform_merchant_catalog.sql` | ✅ |
| `sql/003_merchant_product_app_catalog.sql`（SKU + 可售种子） | ✅ |
| 平台 API：商户列表/启停、入驻审核、商户分类 | ✅ |
| 平台 API：商品分类/品牌 CRUD、商品审核 | ✅ |
| 平台 admin：商户列表/入驻审核/类目/品牌/商品审核页 | ✅ |
| 商户 API：商品 CRUD / 上下架 / 库存 + 默认 SKU | ✅ |
| 商户后台：商品列表 / 发布编辑页 | ✅ |
| C 端 API：catalog 落库（home/categories/products/stores） | ✅ |
| H5：分类/列表/商详/店铺页注册；PC：店铺页 | ✅ |
| OpenAPI：`merchant-catalog.yaml` · `app-catalog.yaml` | ✅ |

阶段 2 主链路完成。

### 阶段 3 进度

| 项 | 状态 |
| --- | --- |
| `sql/004_trade_cart_order.sql`（地址/车/主子单 + 第二商户种子） | ✅ |
| domain `cart` + `trade`（check/create/PaySuccess/发货，幂等） | ✅ |
| C 端 API：`/cart` `/address` `/v2/order/*` `/order/pay/:id` `/orders` | ✅ |
| 商户 API：`/orders` + delivery（paid/status）；平台 API：订单监管只读 | ✅ |
| PC/H5：加购→车→结算→mock/余额支付→订单（路径对齐） | ✅ |
| 平台 admin 订单列表；商户后台订单/待发货发货 | ✅ |
| OpenAPI：`app-trade` · `merchant-trade` · `platform-trade` | ✅ |

验收路径（本地）：执行 `sql/004` → 登录 C 端 `demo/admin123` → 加购商户1+商户2商品 → v2 创建 → mock/余额支付 → 两子单待发货 → `meradmin`/`mer2` 分别发货。

### 阶段 4 进度

| 项 | 状态 |
| --- | --- |
| `sql/005_aftersale_finance.sql`（退款/提现表 + 菜单 + mer_money） | ✅ |
| domain `aftersale` + `finance`（仅退款 0→3；提现拒绝退余额） | ✅ |
| C 端：`/refund/apply` `/refunds*`；H5 申请/列表/详情 | ✅ |
| 商户：售后审核 + 余额/提现申请页 | ✅ |
| 平台：退款监管只读 + 提现审核（通过/拒绝退余额） | ✅ |
| OpenAPI：`app-aftersale` · `merchant-aftersale-finance` · `platform-aftersale-finance` | ✅ |

验收路径（本地）：执行 `sql/005` → 已支付子单申请仅退款 → `meradmin` 同意 → 用户余额回补/库存恢复；商户提现 → 平台拒绝 → `mer_money` 退回。

### 阶段 5 进度

| 项 | 状态 |
| --- | --- |
| `sql/006_promotion_coupon_spread.sql`（券/用户券/分销日志/账单 + 种子） | ✅ |
| domain `promotion`：领券、店铺券→平台券 Quote 计价、下单核销 | ✅ |
| trade：`coupon_user_ids` 进 v2 check/create；PaySuccess → CreditSpreadOnPay | ✅ |
| C 端 API + PC/H5：领券中心/我的券/结算选券；分销中心 | ✅ |
| 商户后台：店铺券 CRUD；平台后台：平台券 + 分销日志/佣金流水 | ✅ |
| OpenAPI：`app-promotion` · `merchant-promotion` · `platform-promotion` | ✅ |
| 积分商城 `/order/v3`、SVIP/积分抵扣互斥、秒杀拼团预售等 | ✅ 积分+秒杀+拼团+SVIP+全款预售已落地 |

验收路径（本地）：执行 `sql/006` → 领平台券/店铺券 → v2 check 带 `coupon_user_ids` 见优惠 → create 核销 → 支付；推广员 `demo` 下级支付后账单可查。

### 阶段 6 进度（首刀）

| 项 | 状态 |
| --- | --- |
| `sql/007_loyalty_points.sql`（积分字段/积分商品/公告种子） | ✅ |
| 积分商城：`POST /order/v3/check|create`（`product_id`）+ `/order/points/pay/:id`（≠ `/v3/order/*`） | ✅ |
| 普通单 v2：`use_integral` 抵扣（§8 券后、最多应付 20%） | ✅ |
| 入口隔离：积分商品禁止进 v2；普通商品禁止进 v3；积分单单独购买 | ✅ |
| C 端：积分商品列表 + PC/H5 兑换结算；公告只读 | ✅ |
| 平台：公告 CRUD 页（6a） | ✅ |
| OpenAPI：`app-loyalty` · `platform-content` | ✅ |
| 6b 客服：JWT + 查单 + 快捷回复 + 本地会话桥（`/api/service/v1` + service-web） | ✅ |
| 6c OpenAPI 6 条（签名换票→JWT，`sql/008` + `/api/open/v1`） | ✅ |
| 6d manager：店员登录/核销/代退（`sql/009` + `app-manager` + `/api/manager/v1` + release `:18087`） | ✅ |
| 商户后台订单核销按钮 | ✅ |
| 6a DIY：`sql/010` + 平台 CRUD/启用 + C 端 `/diy/home`（banners/menus） | ✅ |
| DIY 可视化：`sql/046` + `{page,items[]}` + Vben 三栏编辑器（platform/merchant）+ uni `DiyRenderer` | ✅ |
| 秒杀：场次/活动 + 三端 UI；v2 改价、`product_type=1`/`activity_id`、限购、清店铺券 | ✅ |
| OpenAPI：`app-diy-seckill` · `platform-diy-seckill` · `merchant-seckill` | ✅ |
| 拼团：`sql/011` + 商户 CRUD + 平台监管 + C 端开团/参团；`/order/group/*`；支付计人满员 status 9→0 | ✅ |
| OpenAPI：`app-combination` · `platform-combination` · `merchant-combination` | ✅ |
| SVIP：`sql/012` + v2 会员价、店铺券互斥、活动禁积分；平台用户/商户配置/商品会员价 | ✅ |
| OpenAPI：`app-svip` · `platform-svip` · `merchant-svip` | ✅ |
| 6e 预约：`sql/013` + `/order/reservation/*` + 商户时段；H5/PC 下单 + 商户配置页 | ✅ |
| OpenAPI：`app-reservation` · `merchant-reservation` | ✅ |
| 全款预售：`sql/014` + `activity_type=2`；`/order/presell/check|create`；H5/PC/商户/平台；付后待发货 | ✅ |
| OpenAPI：`app-presell` · `platform-presell` · `merchant-presell` | ✅ |
| 定金预售：`sql/015` + 尾款 `status=10/11` + `/presell/pay/:id`；H5 尾款页/商户建定金活动 | ✅ |
| OpenAPI：`app-presell`（含 finals/pay）· `merchant-presell` | ✅ |
| 6e 直播：`sql/016` + 直播间 CRUD/挂货/审核；`live_status` 101/102/103；H5 列表详情 | ✅ |
| OpenAPI：`app-live` · `merchant-live` · `platform-live` | ✅ |
| 6e 社区：`sql/017` + 发帖/评论/挂货 + 平台审核；H5/PC/商户只读 | ✅ |
| OpenAPI：`app-community` · `platform-community` · `merchant-community` | ✅ |
| 助力：`sql/018` + `activity_type=3`；发起/帮砍 + `/order/assist/*`；H5/PC/商户/平台 | ✅ |
| OpenAPI：`app-assist` · `platform-assist` · `merchant-assist` | ✅ |

验收路径（本地）：
- 全款：`sql/014` → 商品3 全款预售 → `/order/presell/create` → 支付 → 待发货。
- 定金：`sql/015` → 活动2（商品12）定金¥20+尾款¥79 → create → 付定金 → status=10 → `POST /presell/pay/:id` → status=0。
- 直播：`sql/016` → H5 `/pages/live/list` 见演示直播间 → 详情挂货跳商品；商户建房 → 平台审核通过。
- 社区：`sql/017` → H5 `/pages/community/list` 见演示帖 → 评论/挂货；发帖待审 → 平台通过。
- 助力：`sql/018` → 商品13 / set=1（status=10）→ `POST /order/assist/create` → 支付 → set status=20；好友帮砍用 uid=2。

阶段 6 ✅ 收口。下一刀见 **阶段 7**。

---

## 7. 阶段 7 — 打磨与验收

**目标**：不新开大块 P2；把阶段 1–6 竖切收成「可本地演示 + 文档一致 + 已知缺口显式」。

### 7.1 本刀进度

| 项 | 状态 |
| --- | --- |
| 路由去重：H5 `pages.json` 助力 list/detail 重复注册 | ✅ |
| 路由去重：PC `assist` 重复路由；保留 `assist` + `assist/:id` | ✅ |
| OpenAPI 索引：直播文件名 `*-live.yaml`；去掉重复 SVIP 行 | ✅ |
| `go build ./...`（api） | ✅ |
| 验收矩阵（主链路 + 阶段 6 周边）写入本表 | ✅ |
| 已知缺口 / 不做清单 | ✅ |
| Job：超时未支付关单（`CloseExpiredUnpaid`，默认 30min） | ✅ |
| 店员端：待发货 UI + README SQL 名修正（`009_manager_service.sql`） | ✅ |
| `sql/019` 菜单接线 + README 应用命令 | ✅ |
| 平台菜单：退款/财务/提现与品牌 `menu_id=14` 冲突 → 改用 30/31/32 + `019` 修复脚本 | ✅ |
| 社区菜单 `pid` → 内容父级（20） | ✅ |
| 死代码：未接线的 `app/spread`（分销已在 coupon handler） | ✅ |
| PC 首页领券→`/coupons`、补积分入口；个人中心补预售/助力/社区 | ✅ |
| H5「我的」补助力入口；PC 社区发帖页 | ✅ |
| 侧栏：父级目录不跳转 Placeholder；按菜单树展开 openKeys | ✅ |
| 商户社区 path→`/marketing/community` | ✅ |
| 设置最小 CRUD + 菜单可见：`021`；OpenAPI `*-setting` | ✅ |
| RBAC：菜单树勾选 / 半选父级；商户角色菜单 `022`；管理员角色多选 | ✅ |
| 商户子账号绑定角色：`/setting/admins` + `023`；角色仅本店/共享模板 | ✅ |
| `023` 菜单 path 对齐 `/setting/admins`；演示角色「商户运营」+ 子账号 `mersub` | ✅ |
| `/setting` 跳首子页（平台 admin / 商户 shop） | ✅ |
| DIY 首页补「预约」入口（`020`） | ✅ |
| 超时未支付关单：`ListExpiredUnpaidGroups` persist 补齐（job 可调用） | ✅ |
| 活动库存回滚：取消/超时关单归还预售·助力预扣库存；job 挂 `SetPresell`/`SetAssist` | ✅ |
| 拼团未支付取消：`CancelUnpaid` 软删成员；空团/团长退出关团；关单事务成功后释放席位；未支付成员占席；定金预售作废尾款单 | ✅ |
| 按钮级权限：`024` 发货/核销 + `025` 同意/拒绝退款；`RequireMerchantMenu` + FE `hasPerm`；mersub 无核销/无拒绝 | ✅ |
| DIY 素材库：`026` 分类/附件；平台+商户上传页；DIY 选图写入 banner；静态 `/uploads` | ✅ |
| 平台按钮权限：`027` 退款/提现审；`/auth/permissions`；退款监管写操作；演示 `auditor` | ✅ |
| 平台介入闭环：C 端 `POST /refunds/:id/platform`（0→4）+ H5 入口；平台可审 0/4；素材分类校验/商户分类 UI | ✅ |
| 商品按钮：`028` 上下架/改库存；演示 `merprod`（可库存不可上下架） | ✅ |
| 商品发布/删除 + 社区审帖：`029`；`merprod` 可发布不可删；`auditor` 可审不可删 | ✅ |
| 营销启停/直播审房：`030` 秒杀启停·拼团上下架·直播审房；`meract` 仅秒杀启停 | ✅ |
| 优惠券启停：`031` 平台/商户 `coupon/toggle`；`meract`/`auditor` 可启停 | ✅ |
| 预售/助力上下架：`032` `presell/toggle`·`assist/toggle`；`RequireMerchantMenu` + FE；meract 无此二权 | ✅ |
| 新建优惠券：`033` `coupon/create`；meract/auditor 可启停不可新建 | ✅ |
| 删除优惠券 + 活动创建/删除：`034` 券删 + 秒杀/拼团/预售/助力 create·delete；meract/auditor 无删除/创建活动 | ✅ |
| 直播/预约/SVIP：`035` 商户 `broadcast/*`·`reservation/config`·`svip/update` + 平台 `svip/update`；meract/auditor 无写 | ✅ |
| 开播/挂货 + 素材写：`036` `broadcast/live|goods`·`attachment/upload|delete`（平台+商户）；meract/auditor 无 | ✅ |
| DIY + 店员发货细权：`037` `diy/create|update|delete|active|pick`；`staff1.is_goods=1` | ✅ |
| 商户设置写：`038` `shop/update`·`staff/write`·`admins/write`·`roles/write` | ✅ |
| 客服快捷回复可配置：`039` `qixi_store_service_reply` + `/setting/replies`·`reply/write`；service 按 mer_id 读库 | ✅ |
| 商户社区写：`040` `community/create|update|delete`；发帖待审；meract/mersub 无写 | ✅ |
| 协议 + 公告 C 端：`041` `qixi_cache` + `/setting/agreements`·`agreement/update`；H5/PC 公告/协议页 | ✅ |
| 支付渠道沙箱：`payment` 配置；wechat/alipay → PayIntent + `/api/callback/v1/pay/*` HMAC；H5/PC 支付页 | ✅ |
| release 收口：`pack db`→`sql/`；`service-web` Make/ALL；prod Nginx CS；yaml sync 约定写清 | ✅ |

### 7.2 验收矩阵（本地演示）

| 域 | 最小路径 | SQL |
| --- | --- | --- |
| 主链 P0 | 登录 → 加购 → v2 check/create → pay → 商户发货 | `000`–`006` |
| 售后/财务 | 仅退款 → 平台审 → 商户提现审（菜单见 `019`） | 阶段 4 / `019` |
| 营销 | 领券/选券 → 计价 §8；分销关系只读 | 阶段 5 |
| 积分 | `/order/v3` 兑换；v2 `use_integral` | `007` |
| 秒杀/拼团/SVIP | 三端入口 + 对应下单路径 | `010`–`012` |
| 预约 | H5/PC 选时段下单 | `013` |
| 全款/定金预售 | 全款待发货；定金→尾款 status 10→0 | `014`–`015` |
| 直播 stub | 列表/详情挂货；商户建房→平台审 | `016` |
| 社区 | 发帖待审→平台通过；评论/挂货 | `017` |
| 助力 | set 满员 → `/order/assist/create` → 支付 → status 20 | `018` |
| 6b/6c/6d | service 查单；open 签名换票；manager 核销/发货/代退 | `008`–`009` |
| Job | 超时未支付关单（默认 30min）+ 预售/助力库存回滚 + 拼团席位释放 | job 配置 |
| 素材库 | 平台/商户上传图片 → DIY 选图写入 banner | `026` |
| 平台审 | `auditor` 可同意退款/通过提现，不可拒绝 | `027` |
| 平台介入 | H5 售后详情申请介入 → 平台审同意/拒绝 | 阶段 4 + `027` |
| 商品按钮 | `merprod` 可改库存/发布，不可上下架/删除 | `028`–`029` |
| 社区审帖 | `auditor` 可审不可删 | `029` |
| 营销/直播/券/素材 | `meract` 仅启停类；无活动/券/直播 CRUD·开播挂货、无预约/SVIP/素材写；`auditor` 审房/券启停，无券 CRUD、无设会员/素材写 | `030`–`036` |
| DIY/店员 | `admin` 可 DIY 写/选图；`auditor` 无；`staff1` 可核销+发货（`is_goods`） | `037` |
| 设置写 | `meradmin` 可改店铺/店员/子账号/角色；`meract`/`mersub` 无写 | `038` |
| 快捷回复 | 商户维护 → service-web 按 mer_id 加载开启项 | `039` |
| 商户社区 | `meradmin` 可发帖/改帖/删帖；改后待平台审；`meract`/`mersub` 无写 | `040` |
| 协议/公告 | 平台改协议 → C 端读；首页/我的进公告与用户协议 | `041` |
| 支付沙箱 | H5/PC：wechat/alipay → 意图 → 回调验签 → 已支付 | `payment` 配置 |
| 物流/辅资料/文章/标签 | 平台快递·文章·用户标签；商户运费模板·商品标签/保障/参数；C 端读文章；商户发票审 | `045` |
| release | `pack db` 含 sql；`local-service-web` / `local-manager`；frontend-all 含 CS+店员 | pack 脚本 |

演示账号（种子）：平台 `admin`·`auditor` / 商户 `meradmin`·`mer2`·`merprod`·`meract` / 子账号 `mersub` / C 端 `demo` / 店员 `staff1`；密码均为 `admin123`。

### 7.3 已知缺口（明确不做 / 后续按需）

- 直播：无微信推流 / 实时弹幕；房间含 `play_url`/`push_url` stub + 挂货；**非**真实推流。
- 社区：商户可发帖/改帖/删帖（平台审）；热门话题列表已加深；无复杂推荐流。
- 支付：wechat/alipay **沙箱验签闭环**已接（`payment.sandbox` + HMAC `notify_token`）；**非**真实微信/支付宝 SDK / 商户证书。
- 客服：正式 JWT + 本仓会话；`im.mode=remote` 时 S2S UserSig + open-single + service-web **WS 建连**（见 `docs/integration-pte-live-im.md`）；IM 密聊正文仍待 SDK（E2EE）。
- 菜单/按钮：`sql/043` 全量导入 CRMEB `eb_system_menu`（冲突 id +20000；`is_mer`/`is_menu` 重映射；超管/商户模板 rules 重建）。Vben 仅绑定本刀已实现叶子（物流/文章/用户标签/运费模板/商品辅资料/发票等）；其余 CRMEB `route` 仍走 Placeholder。
- 店员端：`app-manager` 已纳入 release（`:18087` / `.35`）。
- 配送员 / 发票 / 短信：最小竖切；无第三方短信/税控/骑手调度。商户发票页已接 `/api/merchant/v1/invoices`。
- **物流 / 商品辅资料 / 文章 / 用户标签**：`sql/045` + `domain/logistics|productmeta|article|usertag` 已落地（平台快递/城市、商户运费模板、标签/保障/参数模板、平台文章 CRUD、用户标签/打标、C 端文章公开读）。无真实物流轨迹/电子面单 SDK。
- 仍缺（后续刀）：用户余额充值/提现真支付、收藏/签到、代客下单、卡密库、打印机 SDK、微信自动回复、财务对账列表深化。
- 协议：常用键已扩；非法律文本审核。
- 前端：未注册叶子仍走 Placeholder。

### 7.4 下一刀建议

已有库执行：`sql/019`–`045`（含 `043` 全量菜单、`044` IM、`045` 物流/辅资料/文章/标签）。

本刀（缺口竖切）：物流（快递/城市/运费模板）+ 商品辅资料 + 平台文章/用户标签 + C 端文章公开 + 商户发票 FE 接线。发货可选 `express_id` 回填快递名。

用户明确要求时再做其一：
1. 余额充值/提现、收藏/签到、代客下单、卡密、打印机、微信自动回复；
2. H5/uni 接 pte-im-sdk 密聊；真实微信/支付宝 SDK / 推流；
3. 本机 `local-*` 联调或 `deploy-*`（须用户授权验证/部署）。

---

## 8. 文档索引关系

| 文档 | 职责 |
| --- | --- |
| 本文件 | **全端开发节奏与里程碑** |
| `features/*` | 按钮/操作级验收 |
| `FUNCTIONAL-TRUTH.md` | 高风险真相 |
| `schema/*` | 表结构 |
| `SERVICE-MATRIX.md` | 命名/网络/发布 |
| `domain-modules.md` | 包边界 |
