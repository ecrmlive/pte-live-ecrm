# ecrm-harmo — 七禧鸿蒙客户端

独立 OpenHarmony API 23 原生工程，使用 ArkTS、ArkUI、MVVM 与 Repository 分层；不以 `app-uni` WebView 替代。

## 当前工程范围

- 已提供标准 OpenHarmony 工程配置、`EntryAbility` 和首页 Tab 骨架。
- 首页包含商城主入口、登录入口、商品浏览、购物车、订单和“我的”占位页面，方便按 C 端功能清单继续实现。
- 网络层使用 `@ohos.net.http`（`@kit.NetworkKit`）并预置统一的 `Authori-zation` 与 `X-AppId` 请求头约定。
- 认证 Repository 固定提交 `channel: harmony`；与 `app-pc`、`app-uni`、`app-ios`、`app-adnroid` 共用 `api-business` 的 C 端 JWT。

尚未实现真实登录、商品、订单和支付能力，也没有内置任何环境地址或凭据。

## 目录

```text
AppScope/                         应用级配置
entry/src/main/ets/
  entryability/                   UIAbility 入口
  common/                         常量、主题、网络客户端
  data/repository/                API Repository
  domain/model/                   领域模型
  features/                       按页面组织的 UI 与 ViewModel
  pages/                          应用根页面
config/app.example.yaml           非敏感环境配置示例
```

## 开发前配置

1. 在 DevEco Studio 打开本目录，安装 OpenHarmony API 23 SDK。
2. 复制 `config/app.example.yaml` 为本地未跟踪配置，并填入受控环境提供的 API 地址；不要提交真实地址中的凭据、Token 或证书。
3. 将 `common/constants/AppConfig.ets` 的占位 API 地址替换为构建时注入的环境配置。当前占位值不会连接任何生产服务。

页面顶部导航固定为 44vp，并放在系统状态栏安全区之后；交互按钮最小高度为 44vp。
