# H5 用户端覆盖矩阵（持续验收）

> 基线：CRMEB MER v4.0 用户端 342 项路由操作，详见 `docs/features/03-user-app.md`。
> 本表按用户可见能力归并，不能用页面数量替代逐操作验收；每项只有同时具备入口、真实 API、权限/状态机、测试和同态截图后才可关闭。

| 能力域 | CRMEB 对照 | 当前入口与实现 | 状态 |
| --- | --- | --- | --- |
| 登录、资料与密码 | `auth/*`、`user/change/*` | 登录、注册、微信小程序静默登录、资料修改、密码修改、短信重置密码、手机号绑定与换绑、短信登录 | 本轮完成短信登录/绑定/换绑/重置密码 UI、一次性摘要验证码、图形安全验证与加密短信网关适配；真实短信网关配置与投递验收仍待完成 |
| 账户注销 | `POST api/user/cancel` | `pages/profile/cancel` → `POST /api/app/v1/auth/cancel` | 本轮完成：资金/积分/佣金、未完成订单、处理中售后触发二次确认；停用账户并递增 `auth_version`，不物理删除历史事实 |
| 首页、DIY、商品、分类与搜索 | `common/home`、`diy`、`product/*`、`system/city` | `pages/index`、`diy`、`goods`、`category`、`brand`、`city/select` | 本地闭环：分级选城、首页当前城市与收货地址带回省/市/区已实现；业务城市数据可由平台管理员触发重同步：平台事务 outbox 发布、业务 NATS 订阅写入投影；初始化夹具仅用于本地验收，平台城市增改删管理页仍待补齐；浏览器 GPS 反查需接入已授权地图服务后单独验收 |
| 品牌目录与筛选 | `brand/lst` | `pages/brand/list` → `GET /api/app/v1/catalog/brands`；商品列表以 `brand` 精确筛选 | 本轮完成：商户录入、平台审核投影、中文品牌夹具、用户端目录和筛选；尚待真实 MySQL 集成回归 |
| 新人权益 | `new_people` | 注册事务自动发放 `qixi_crm_b_onboarding_coupon`；`pages/coupon/newcomer` → `GET /api/app/v1/coupons/newcomer` | 本轮完成：策略、券映射、中文夹具、账号/小程序静默注册同事务发券和已到账展示；待真实 MySQL 事务回归 |
| 开屏广告 | `open_screen` | 首页 `GET /api/app/v1/open-screen`，支持运营开关、跳转、跳过及展示间隔 | 本轮完成：配置模型、中文夹具、接口、H5 展示控制；待真实 MySQL 与真机视觉回归 |
| 店铺、购物车与普通交易 | `store/*`、`v2/order/*` | 店铺街、店铺页、购物车、结算、订单、物流、发票 | 主路径已接入；真实支付回调、MySQL 并发及截图仍未验收 |
| 售后与互动 | `refund/*`、社区路由 | 售后申请/寄回/详情、社区、评价、反馈 | 已有主体实现；需要以真实状态数据回归 |
| 营销与会员 | 券、秒杀、拼团、预售、助力、积分、充值、SVIP | 对应活动页、积分、充值、SVIP | 已有主体实现；营销并发/回调仍需集成验收 |
| 分销、客服与内容 | `spread/*`、`service/*`、内容路由 | 分销、公告资讯、客服、直播、预约 | 已有主体实现；本机 CLI 的 tslib 类型诊断和同态截图尚未关闭 |

## 本轮手机号绑定验收场景

1. H5 资料页可进入绑定手机号页；已绑定号码仅展示状态，不允许在该入口直接替换。
2. 发送验证码前必须通过 `login_sms` 图形验证；验证码仅以手机号、用途、验证码三元摘要入库，五分钟过期且仅可消费一次。
3. 同手机号/用途一分钟内重复发送返回 429；绑定时以事务校验手机号唯一性。
4. 短信仅经平台加密配置的 HTTPS 网关发送，授权令牌不返回客户端、不写入夹具、日志或响应。
5. 当前未提供真实短信网关凭据；H5 构建与网关模拟单测通过；真实短信投递仍待受控环境验收。

自动验证：`go test ./internal/business/auth ./internal/pkg/smsclient ./cmd`、`UNI_INPUT_DIR=. ./node_modules/.bin/uni build -p h5`。

## 本轮城市选择验收场景

1. `GET /api/app/v1/system/city/lst/:pid` 只读取 `qixi_crm_b_city_view` 可见节点，不跨库读取平台城市配置，也不暴露写入口。
2. 本地 utf8mb4 夹具覆盖中国演示区域 → 华东/华南演示省 → 杭州/宁波/深圳/广州 → 区县四级中文层级。
3. 首页选择到城市级时持久化当前城市；再次返回首页会显示已选城市。
4. 收货地址使用同一选择器并带回省、市、区与区域编码；仍允许手动修正文本。
5. 浏览器 GPS 坐标到行政区划的反查依赖生产地图服务授权与隐私告知，当前不以伪造坐标或未经授权的第三方服务替代。

自动验证：`go test ./internal/business/city ./cmd`、`UNI_INPUT_DIR=. ./node_modules/.bin/uni build -p h5`。

## 本轮账号注销验收场景

1. 无余额、积分、佣金、未完成订单、处理中售后的用户：首次请求直接逻辑停用账户；所有既有令牌因 `auth_version` 递增而失效。
2. 存在任一受限事实的用户：首次请求只返回 10 分钟有效的确认令牌摘要匹配凭据及中文原因；确认令牌不入库明文。
3. 令牌失效、被替换或错误：服务端返回冲突，账户和历史数据不变。
4. 注销成功：仅把用户状态置为停用并写审计记录；订单、资金、售后、发票及客服业务事实均保留。
5. H5：个人中心直接进入“注销账户”；二次确认后清理本地会话并返回“我的”。

自动验证：`go test ./internal/business/auth`、`go test ./internal/business/order ./internal/business/catalog ./internal/business/funding`、`UNI_INPUT_DIR=. ./node_modules/.bin/uni build -p h5`。页面模板空值告警已清零；本机 CLI 仍有 DCloud tslib 虚拟模块 rootDir 诊断，不能据此宣称全部发布阻断项已经关闭。
