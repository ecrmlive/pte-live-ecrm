# app-ios — 七禧 iOS 客户端

原生 C 端商城工程，目标为 **iOS 16+ / Swift 6 / UIKit**。与 PC、H5、小程序、Android、鸿蒙共用 `api-business` 的 C 端契约、用户主体和 JWT；认证来源固定标记为 `ios`。

## 工程入口

```text
ECRM.xcworkspace        # 生成后用 Xcode 打开
ECRM.xcodeproj          # XcodeGen 生成项目
project.yml             # 项目定义的唯一来源
ECRM/
  App/                  # 应用装配、生命周期、DI
  Core/                 # 配置、网络、日志、持久化、安全存储
  Domain/                # 领域实体与仓储协议
  Data/                  # 接口与本地存储实现
  Features/              # UIKit 页面（MVVM）
  Resources/             # 图片、颜色等资源
ECRMTests/              # 基础单元测试
```

`project.yml` 是 Xcode 工程的可重建来源。更新工程定义后，在本目录执行：

```bash
xcodegen generate
```

## 已搭建能力

- UIKit 生命周期与 5 栏用户端导航：首页、分类、购物车、订单、我的；导航栏使用系统标准 44pt 高度并紧贴状态栏。
- Clean Architecture 基础分层、构造器注入的 `AppContainer`、`async/await` 网络客户端与 `OSLog`。
- Keychain 令牌存储、Core Data 本地持久化基座、可替换的认证仓储。
- HTTP 契约：令牌仅写入 `Authori-zation: Bearer <token>`；店铺上下文使用 `X-AppId`；登录来源为 `ios`。

## 配置和安全

- 当前 `ECRMAPIBaseURL` 是不可访问的示例地址。请在本地 Xcode Scheme 或未提交的 `.xcconfig` 中覆写为实际 `api-business` HTTPS 地址。
- 令牌仅存 Keychain；不得将 JWT、支付参数、IM UserSig、密钥或真实用户数据写入源码、`Info.plist` 或日志。
- 商户店铺业务必须先调用 `/auth/store-context` 换取上下文令牌；`X-AppId` 不是 IM SDK AppId。

## Release 打包

`project.yml` 已定义 `MARKETING_VERSION=0.1.0` 与 `CURRENT_PROJECT_VERSION=1`。每次发布前通过 Xcode Scheme 或 `xcodebuild` 覆写为递增构建编号，并在 Apple Developer 受控账号中完成签名；证书、Provisioning Profile、私钥和 App Store Connect 凭据均不得写入项目或平台后台。

后台“应用 → App → iOS”维护 Bundle ID、发布版本、构建编号、App Store/企业分发下载地址、Universal Link 与更新说明。

## 后续实现顺序

1. 认证与用户资料；2. 商品/店铺消费视图；3. 多商户购物车；4. 订单与支付；5. 售后、营销、会员与客服。

功能验收以 [`docs/features/03-user-app.md`](../docs/features/03-user-app.md) 为准；交易等高风险域必须先对照 `docs/api/FUNCTIONAL-TRUTH.md`。
