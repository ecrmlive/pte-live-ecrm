# 七禧 CRM 系统架构总则

> 本文件是系统边界、账号、JWT、数据库、目录与迁移顺序的唯一口径。与旧文档、旧目录、旧 SQL、旧 Make 目标冲突时，以本文件为准；冲突项必须删除或重写，不能兼容保留。

## 1. 正式系统与入口

| 系统 | 入口形态 | 技术 | 账号/权限 |
| --- | --- | --- | --- |
| PC 用户端 | 独立 Web 商城 | Vue 3 | C 端用户 |
| 小程序 & H5 用户端 | 同一套程序 | uni-app x | C 端用户 |
| iOS 用户端 | 原生 App | Swift 6 / UIKit | C 端用户 |
| Android 用户端 | 原生 App | Kotlin / Jetpack Compose | C 端用户 |
| 鸿蒙用户端 | 原生 App | ArkTS / ArkUI | C 端用户 |
| 统一后台管理系统 | `/admin/` | Vben 5.7+ | 平台、商户、区域、客服、运营角色 |
| 店铺管理系统 | `/merchant/` | Vben 5.7+ | 店铺主账号及店铺员工角色 |

平台、商户、区域、客服、运营是**统一后台管理系统内的五种角色菜单和数据范围**，不是五个前端项目。店铺管理系统是唯一独立后台前端。店员、配送员、服务人员属于店铺管理系统角色；客服属于统一后台管理系统角色。

PC、小程序&H5、iOS、Android、鸿蒙是独立用户端交付物：`app-web` 负责 PC Web 页面和交互，`app-mp` 是小程序与 H5 的唯一共享程序，原生端分别使用 `app-ios`、`app-adnroid`、`app-harmony`。六端只通过 `api-business` 访问同一份 `qixi_crm_business` 业务数据，不能以后台页面或接口替代用户端页面。

## 2. 账号与 JWT

### C 端

一个用户主体可绑定多个登录来源，必须区分并记录 `wechat`、`mini_program`、`h5`、`pc`、`ios`、`android`、`harmony`。渠道只表示认证来源和风控审计，不复制用户、订单或资产。

- PC、小程序、H5 使用同一套 C 端 JWT 配置与签名密钥。
- iOS、Android、鸿蒙与 PC、小程序、H5 使用同一套 C 端 JWT 配置与签名密钥。
- JWT 的 `subject` 为统一 C 端用户 ID，包含 `scope=c_user` 与登录渠道；六端均可验证和刷新同一用户 token。
- 所有 JWT HTTP 请求只使用 `Authori-zation: Bearer <token>`；店铺范围上下文只使用 `X-AppId`。完整规则见 [JWT、店铺 AppId 与商户 IM SDK AppId 契约](./auth-store-appid-im-contract.md)。

### 统一后台

- 平台、商户、区域、客服、运营共用一套后台 JWT 配置与签名密钥。
- JWT 包含后台账号 ID、`scope=admin_console`、角色代码和数据范围版本；菜单由服务端按角色返回，接口再次做 RBAC 与数据范围校验。
- 角色代码固定：`platform`、`merchant`、`region`、`customer_service`、`operations`。

店铺管理系统是独立安全边界，使用店铺 JWT；不得复用 C 端 JWT 或统一后台 JWT。任何 JWT 密钥只允许写入对应环境的 `app.yaml`，不得提交到仓库。

## 3. 数据库与表前缀

| 数据库 | 表前缀 | 所有者与内容 |
| --- | --- | --- |
| `qixi_crm_admin` | `qixi_crm_a_` | 统一后台账号、角色、菜单、数据范围、区域、运营配置、后台审计与后台配置 |
| `qixi_crm_business` | `qixi_crm_b_` | C 端用户及多端身份绑定、商品消费视图、购物车、订单、支付、售后、会员、营销、资金与业务审计 |
| `qixi_crm_merchant` | `qixi_crm_m_` | 店铺主体、店铺账号/员工、商品经营资料、库存、履约、店铺装修、物流、店铺财务视图与店铺操作日志 |
| pte-live-im | 以 pte-live-im 实际 SQL 为准 | IM 的后台管理、会话、消息、成员与场景表 |

禁止再创建或引用 `qixi_m_admin_`、`qixi_m_app_`、裸 `qixi_` 业务表。qixi CRM 与 pte-live-im 不共享数据库表，也不得在跨库关系上创建外键；只交换稳定 ID、事件和受控 API。

业务数据跨库规则：每张表只归属一个库；跨库只保存对方 ID；读模型通过 API/事件同步；下单、支付、退款、库存和结算必须有幂等键与本库事务，不允许跨库事务。

### 支付主体边界

- 平台自营商品（消费视图 `merchant_id=0`、`store_id=0`）只使用 `qixi_crm_admin` 中的平台微信/支付宝配置。
- 入驻商户商品只使用该店铺在 `qixi_crm_merchant` 中独立维护的微信/支付宝配置；平台配置绝不作为商户支付的兜底。
- 平台配置与店铺配置均以加密值落库；仅将加密运行时读模型投影至 `qixi_crm_business`。客户端和后台列表均不返回私钥、证书或 API 密钥。
- 不同支付主体的商品不得生成同一个合并支付单，购物车必须分别结算。

## 4. IM 边界

IM 必须使用 `pte-live-im` 与 `pte-live-im-sdk`。消息正文、会话、UserSig 和 IM 管理数据必须直接遵循 pte-live-im `sql/init_im_schema.sql` 的实际表规则、鉴权规则和 SDK 机制；当前实际前缀包括 `pte_live_im_admin_*`、`pte_live_im_*`、`pte_live_chat_*` 与 `pte_live_scene_*`。禁止七禧另起一套 IM 规则、协议或消息表。pte-live-im 是 IM 的唯一实现与数据所有者，七禧仓库不得修改其源码、初始化 SQL、表或容器；七禧只保存业务关联（如订单卡片、商户/店铺授权、会话业务索引）。

七禧容器网络、IP、数据库与 pte-live 项目独立。调用 IM 通过明确配置的服务地址和服务端凭证完成，禁止加入 `pte_live_net` 或直接使用 pte-live 容器名。

商户可以绑定多个 pte-live-im SDK AppId，但每个商户只能启用一个；SDK AppId 不等于店铺 `X-AppId`。选择与会话快照规则见 [JWT、店铺 AppId 与商户 IM SDK AppId 契约](./auth-store-appid-im-contract.md)。

## 5. 目标目录与服务边界

后续目录调整必须收敛为：

```text
admin-platform/      # 一个 Vben 5.7+：平台/商户/区域/客服/运营
admin-merchant/      # 一个 Vben 5.7+：店铺管理
app-web/             # PC Web 商城
app-mp/              # 一个 uni-app x：小程序与 H5
app-ios/             # iOS 16+，Swift 6 / UIKit
app-adnroid/         # Android 12+，Kotlin / Jetpack Compose（目录名按项目约定）
app-harmony/         # OpenHarmony API 23，ArkTS / ArkUI
api-platform/        # 独立 Go module；统一后台 API，qixi_crm_admin
api-business/        # 独立 Go module；C 端业务 API，qixi_crm_business
api-merchant/        # 独立 Go module；店铺 API，qixi_crm_merchant
job/                 # 独立 Go module；事件、异步任务和对账
contracts/           # 仅 OpenAPI/事件 Schema；不得存放可执行业务实现
sql/                 # 三个库的集中初始化 SQL
```

`admin-platform` 是统一后台的唯一 Vben 目录，`admin-merchant` 是店铺系统的唯一 Vben 目录。店员能力属于 `admin-merchant` 的角色菜单，客服能力属于 `admin-platform` 的角色菜单；不保留独立店员或客服前端目录。后端收敛为 `api-platform`、`api-business`、`api-merchant` 三个服务契约。

四个后端目录都是**独立项目**：各自必须拥有 `go.mod`、`cmd/`、`internal/`、`conf/app.yaml`、测试和独立发布产物；不得再保留 `api/cmd/*` 这种单一 Go module 下的多服务目录。服务之间禁止引用对方的 `internal/` 包、禁止共享数据库连接或跨库事务。可共享的仅是 `contracts/` 中版本化的 OpenAPI 与 NATS 事件 Schema；认证、领域、持久化、迁移和运行配置必须在所属服务项目内。

旧 `api/`、`api/cmd/api-admin`、`api/cmd/api-app` 已清理；`api/build/` 不是合法缓存位置。所有正式构建只能从四个独立服务目录发起。

## 6. 迁移顺序

1. 文档：删除旧双前缀、多个独立后台、共享 pte 网络/数据库等描述，锁定本总则。
2. SQL：删除旧运行初始化文件，按三个 qixi CRM 库和 pte IM 边界重新设计 DDL、基础数据、配置、密钥占位和测试数据。
3. Make / Compose / 目录：统一 `qixi_mergers` Compose 身份与独立 `qixi_mergers_net`，建立三个 API、两个 Vben 与两个用户端的构建入口。
4. 代码合并：先合并统一后台五角色菜单，再独立店铺系统；随后迁移 API、鉴权、数据访问与 OpenAPI。
5. 功能开发：按 CRMEB 功能总表逐项实现、截图和闭环验收。

当前 `api/` 单 module 与旧后台目录仍在迁移中；在四个独立 Go 项目、两套 Vben 和发布脚本全部收敛并通过验证前，不得把“文档已定版”表述为目录、代码或功能已经迁移完成。
