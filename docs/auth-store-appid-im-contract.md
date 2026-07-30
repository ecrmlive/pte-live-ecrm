# JWT、店铺 AppId 与商户 IM SDK AppId 统一契约

> 状态：**迁移实施中**。本文件优先于旧代码中 `Authorization`、`AppID`、`App-Id`、全局单 IM `app_id` 的兼容写法；未完成项见第 7 节，不得宣称已完全收敛。

## 1. 唯一 HTTP 请求头

### 1.1 JWT

所有需要 JWT 的 HTTP 接口只识别下面的请求头，名称和格式固定：

```http
Authori-zation: Bearer <access-token-or-refresh-token>
```

- 不接受 `Authorization`、query 参数、cookie、请求体字段或自定义 token 头作为 JWT 来源。
- 普通业务请求携带 access token；`/auth/refresh` 携带 refresh token。服务端按路由校验 token `kind`，但来源仍是同一个 `Authori-zation` 请求头。
- 登录、短信验证码、支付/IM 服务商回调等公开接口不携带业务 JWT；它们必须使用各自的签名/验签机制，不能借用 JWT。
- 令牌返回仍采用 JSON `access_token` / `refresh_token`；响应头不回传 JWT，避免浏览器/代理错误缓存。

### 1.2 跨域预检

所有 qixi API 与网关对 OPTIONS 使用同一允许头列表，禁止每个服务自行增删：

```http
Access-Control-Allow-Headers: Content-Type, X-Requested-With, Form-type, Referer, Connection, Content-Length, Host, Origin, Authori-zation, Accept, Accept-Encoding, X-AppId
```

响应头应包含 `Vary: Origin, Access-Control-Request-Headers`。`Authori-zation` 是非标准拼写但为微信端约束，后端、网关、PC、H5、uni-app x 和原生端必须完全一致。

## 2. JWT 域、主体与完整身份快照

JWT 只承载**服务端已解析的身份快照**。客户端请求体、query、cookie 不能写入或覆盖其中任一身份字段；JWT 也不保存手机号、密码、支付参数、IM UserSig、IM 密钥或数据范围明细。

所有新签发的 access / refresh JWT 都必须包含以下字段。`merchant_app_id` 和 `im_sdk_app_id` 即使当前主体没有店铺上下文也必须出现，值为空字符串，配合 `auth_context` 明确表示“无商户上下文”，绝不以伪造的 `0`、`default` 或其他商户值替代。

| Claim | 含义与规则 |
| --- | --- |
| `sub`、`principal_id`、`principal_type` | 不可变主体标识；类型仅可为 `c_user`、`admin_user`、`store_account`、`merchant_user`、`open_client`。`sub` 必须等于十进制 `principal_id`。 |
| `roles` | 当前角色代码数组；C 端为 `customer`，统一后台为实际 RBAC 角色，店铺后台为 `owner` / `manager` / `clerk` / `delivery` / `service`。不能由前端菜单反推。 |
| `client_platform` | 令牌签发端：C 端仅可为 `wechat`、`mini_program`、`h5`、`pc`、`ios`、`android`、`harmony`；统一后台为 `admin_web`；店铺后台为 `merchant_web`。 |
| `portal`、`scope`、`auth_context` | API 域、权限域和上下文类型。上下文只能是 `platform`、`admin`、`merchant`、`store`，避免同一个账号被错误解释成其他端身份。 |
| `merchant_id`、`store_id`、`merchant_app_id` | 商户/店铺身份快照。`merchant_app_id` 是 `qixi_crm_m_store.app_id`，也是 HTTP `X-AppId`；店铺后台令牌三者均必填。 |
| `im_sdk_app_id` | 该商户当前启用的 pte-live-im SDK AppId 快照；只允许服务端按 `merchant_id` 查询后写入。未配置 IM 时显式为空，不得由客户端指定。 |
| `identity_version`、`data_scope_version` | 账号身份版本和统一后台数据范围版本。密码修改、禁用、强制退出、角色/商户归属变更必须递增身份版本；权限范围变更必须递增数据范围版本。 |
| `session_id`、`jti` | 一次登录会话标识和单令牌唯一 ID；access / refresh 共用 `session_id`，但各自 `jti` 必须不同，用于审计、注销和重放检测。 |
| `iss`、`aud`、`iat`、`nbf`、`exp` | 标准 JWT 约束。`iss=qixi-mergers`，`aud=qixi-mergers:<portal>`；服务端必须校验算法 HS256、受众、有效期和未生效时间。 |

兼容期中的 `uid`、`admin_id`、`channel`、`store_app_id` 只供旧 handler 读取：新接口以 `principal_*`、`client_platform`、`merchant_app_id` 为准；店铺令牌中 `store_app_id` 必须与 `merchant_app_id` 完全相等。

JWT 只证明“谁在调用”，不承载可伪造的店铺选择：

| 调用端 | scope | 主体 | 允许的店铺上下文 |
| --- | --- | --- | --- |
| PC / 小程序 / H5 / iOS / Android / 鸿蒙 | `c_user` | `c_user` / `principal_id` | 平台首页令牌的店铺/IM 字段为空；进入店铺上下文后必须由 `X-AppId` 解析，订单、商品、客服会话仍以服务端归属二次校验 |
| 统一后台（平台/商户/区域/客服/运营） | `admin_console` | `admin_user` / `principal_id`、角色、数据范围版本 | 全局后台账号没有唯一商户时 AppId/IM 字段为空；不能把可管理的商户列表塞入 JWT，服务端按 RBAC/数据范围解析 |
| 店铺管理系统 | `store_console` | `store_account` / `principal_id`、`merchant_id`、`store_id`、角色 | `merchant_app_id` 与 `X-AppId` 必须映射到 token 内同一个店铺；同时携带当前 `im_sdk_app_id` 快照 |

同一域内的 JWT 密钥规则不变：C 端六端共用一套；统一后台五角色共用一套；店铺后台独立一套。JWT 与 pte-live-im UserSig 是两个协议，绝不互换。

## 3. X-AppId：店铺应用标识，不是 IM SDK AppId

`X-AppId` 是七禧店铺应用标识，固定映射到一个店铺；它不是微信 AppID、支付 AppID，也不是 pte-live-im 的 `sdk_app_id`。

```text
merchant (1) ── (1) store ── (1) qixi store_app_id  ← X-AppId
merchant (1) ── (N) IM SDK AppId                     ← 当前启用一条
```

规则：

1. 一个入驻商户对应一个店铺，`merchant_id` 在店铺表中唯一；一个店铺拥有一个全局唯一、不可复用的 `app_id`。
2. 店铺范围接口必须携带 `X-AppId`；服务端只通过数据库把它解析为 `store_id`，客户端不可直接提交或覆盖 `store_id` / `merchant_id`。
3. 店铺后台的 `X-AppId` 必须与店铺 JWT 内的 `store_id` 一致，不一致直接返回 403。
4. C 端商城的店铺页、商品页、购物车、订单、客服等店铺范围接口使用 `X-AppId`；订单/商品/会话仍以服务端归属数据二次校验，不能只凭 header 放行。
5. 平台自营使用保留主体 `merchant_id=0, store_id=0`，不伪造入驻商户 `X-AppId`；平台后台通过 RBAC 选择自营主体。

## 4. 商户 IM SDK AppId

商户可绑定多个 pte-live-im SDK AppId，但同一商户任意时刻只能有一个“当前启用”记录：

| 字段 | 说明 |
| --- | --- |
| `merchant_id` | 商户归属；一个商户多条记录 |
| `sdk_app_id` | pte-live-im 返回的 SDK AppId；与 `X-AppId` 完全不同 |
| `name` | 商户可读名称，例如“客服生产实例” |
| `status` | `enabled` / `disabled` |
| `is_active` | 仅一条可为 `1`；切换必须事务化 |
| `api_public_url` / `ws_public_url` | 仅客户端可访问的 IM 公网地址，禁止 Docker 服务名 |
| `pte_profile_id` | 七禧受控 S2S 配置引用；不保存或回显 pte 集成 Token |

选择流程：

```text
X-AppId → store_id → merchant_id → 当前启用的 merchant IM SDK AppId
→ 服务端向对应 pte-live-im 配置申请 UserSig → 返回 sdk_app_id + UserSig
```

- 客户端不能传 `sdk_app_id` 来指定租户；只能接收服务端签发的凭证。
- 切换当前 SDK AppId 后，新会话立即使用新配置；已建立会话保留创建时的 `sdk_app_id` 快照，避免会话跨 IM 租户。
- 平台自营 IM 与商户 IM 同样隔离；平台自营可以有自己的平台配置，但不能作为商户缺省 IM 配置。

## 5. 数据模型与跨库投影

所有真实密钥只在所属库加密保存，不进 Git、不进 JWT、不通过 C 端接口返回。

| 所属库 | 表/对象 | 约束 |
| --- | --- | --- |
| `qixi_crm_merchant` | `qixi_crm_m_store.app_id` | `UNIQUE(app_id)`、`UNIQUE(merchant_id)`；`X-AppId` 的唯一来源 |
| `qixi_crm_merchant` | `qixi_crm_m_im_sdk_app` | `UNIQUE(merchant_id,sdk_app_id)`；事务 + 唯一活动约束保证每商户至多一条 active |
| `qixi_crm_business` | `qixi_crm_b_store_view` | 店铺 AppId → store/merchant 的只读投影 |
| `qixi_crm_business` | `qixi_crm_b_merchant_im_sdk_app_view` | 当前启用 IM AppId 的只读投影；会话建立时读取 |
| `qixi_crm_business` | 订单/客服会话快照 | 保存 `store_id`、`merchant_id`、`store_app_id`、`im_sdk_app_id`，不可被后续切换改写 |
| `qixi_crm_admin` | 平台自营 IM/云配置 | 仅服务平台自营主体，不能覆盖商户配置 |

商户库到业务库只能经 outbox + NATS 事件（`store.app_id.changed`、`merchant.im_sdk_app.activated`）建立读模型；`api-business` 不得直接查询 `qixi_crm_merchant`。

## 6. 统一接口要求

| 场景 | JWT | X-AppId | 服务端校验 |
| --- | --- | --- | --- |
| C 端店铺商品/购物车/下单/客服 | `Authori-zation` | 必填 | AppId→店铺；商品/订单/会话必须属于该店铺 |
| 店铺后台 | `Authori-zation` | 必填 | AppId→店铺必须等于 store JWT 中 `store_id` |
| 统一后台 | `Authori-zation` | 不作为身份来源 | RBAC 与数据范围；需要操作店铺时以服务端选择结果校验 |
| 登录/注册 | 无 JWT | 仅店铺定制登录页可传 | 不能以 header 赋予权限 |
| 支付/IM 服务商回调 | 无业务 JWT | 不使用 | 使用支付平台/IM 服务商签名与事件幂等键 |

## 7. 当前差异与迁移顺序

已完成：qixi API JWT 中间件与刷新接口只读取 `Authori-zation`；五个 API 服务的 CORS 已统一；PC、uni-app x 与两个 Vben 主请求客户端已改用 `Authori-zation`；店铺后台 JWT 已携带 `store_app_id` 并强制校验 `X-AppId`；商户库已具备店铺 `app_id` 与 IM SDK AppId 映射表、店铺后台配置入口和事务化当前启用切换。

已补齐：新 JWT 的主体类型/主体 ID、角色、客户端平台、商户 AppId、当前 IM SDK AppId、身份/数据范围版本、会话 ID、`jti`、`aud`、`nbf`；店铺登录和刷新从 `qixi_crm_m_im_sdk_app` 读取当前启用 SDK AppId 并写入令牌。初始化表已增加 `auth_version` 字段；店铺/统一后台密码修改会递增该版本，刷新和当前身份查询会拒绝旧版本。

已补齐：`api-business` 的 `POST /api/app/v1/auth/store-context`。C 端先使用全局登录令牌，再以 `X-AppId` 让业务库的 `qixi_crm_b_store_view` 和 `qixi_crm_b_merchant_im_sdk_app_view` 解析店铺及当前 IM SDK AppId，换取携带实际商户上下文的 JWT；刷新该 JWT 时会复核投影，店铺或 IM 当前配置已变更则拒绝旧上下文。

已补齐（本轮）：PC、uni-app x 在店铺页和商详读取消费视图返回的 `merchant_app_id` 后换发店铺上下文 JWT；统一请求客户端自动注入 `X-AppId`；入驻商户商品 `POST /cart` 会同时校验 JWT 的店铺快照、请求 `X-AppId` 和商品消费视图归属，不能只伪造 header。

已补齐（本轮）：店铺后台对当前启用 IM SDK AppId 的启用、停用、在用配置更新均在同一商户库事务内写入 `qixi_crm_m_outbox`；事件只含商户、SDK AppId、公开端点和受控配置引用，不含 PTE Token、UserSig 或密钥。`api-merchant` 重试发布 pending Outbox 至 NATS `qixi.merchant.im-sdk-app.v1`，`api-business` 订阅并幂等维护 `qixi_crm_b_merchant_im_sdk_app_view`，全程不跨库直查商户表。

仍未完成：C 端的 `X-AppId → qixi_crm_b_store_view` 全量路由强制与订单/客服会话二次归属校验；pte-live-im 商城 S2S 目前仍强制单个业务 AppId，尚未支持按七禧商户“当前 SDK AppId”签发 UserSig；所有受保护业务路由对 `identity_version` 的缓存化即时拒绝；存量库仍需补齐 `qixi_crm_m_store.app_id`。`auth_version` 已通过三个 `01_table.sql` 中的 MySQL 8.4 幂等 `ALTER TABLE ... ADD COLUMN IF NOT EXISTS` 纳入初始化/升级。

迁移必须按以下顺序执行：

1. 完成存量库迁移、商户 outbox/NATS 投影、C 端 `X-AppId` 解析和订单/会话快照字段。
2. 将剩余原生端/辅助入口的 JWT 请求改为 `Authori-zation`，删除所有遗留 `Authorization`、`AppID`、`App-Id` 兼容代码。
3. 扩展 pte-live-im 的受信商城 S2S 合约，使七禧服务端可按受控商户当前 SDK AppId 签发 UserSig；补齐会话快照与切换回归测试。
4. 补齐 JWT 头、跨店拒绝、IM 切换快照、CORS 预检的跨服务集成测试。
