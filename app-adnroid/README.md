# ecrm-Android — 七喜 Android 客户端

原生 C 端 Android 工程，目录名按项目约定固定为 `app-adnroid`。

## 已搭建

- Android 12 / API 31+，Kotlin、Jetpack Compose、Material 3。
- MVVM/Clean 的 `presentation`、`domain`、`data`、`core` 分层；Hilt 注入。
- Retrofit + OkHttp（含 WebSocket 能力）连接 `api-business`，请求令牌仅使用 `Authori-zation: Bearer <token>`。
- Room 本地缓存、DataStore 会话存储、Coroutines/Flow、Timber 日志。
- U-Push 初始化边界与 Manifest 占位符；真实 SDK、AppKey 和厂商配置后续由受控环境接入，禁止提交仓库。
- 首页、购物车、我的三个导航入口；仅首页作为骨架演示，不宣称已完成商城功能。

## 业务契约

- 用户登录来源固定为 `android`；Android 与 PC/H5/小程序/iOS/鸿蒙共享 C 端用户、JWT、订单和资产。
- C 端接口来自 `api-business`；店铺范围接口使用服务端校验的 `X-AppId`，不能由客户端伪造商户或店铺 ID。
- IM 后续仅接入 `pte-live-im` 与 `pte-live-im-sdk`，不自建会话或消息协议。

## 本地配置

仓库不保存地址、友盟 AppKey 或任何密钥。同步/构建时按需注入 Gradle 属性：

```text
ECRM_API_BASE_URL=https://your-api.example/
UMENG_APPKEY=your-app-key
```

开发环境未注入 API 地址时使用无效占位地址，避免误连真实环境。

## Release 打包

发布版本由 Gradle 属性注入，默认值仅用于开发。密钥库和密码只放在本机的 `~/.gradle/gradle.properties` 或 CI Secret，禁止写进项目、后台或 Git：

```text
ECRM_APPLICATION_ID=com.example.ecrm
ECRM_VERSION_NAME=1.0.0
ECRM_VERSION_CODE=1
ECRM_RELEASE_STORE_FILE=/absolute/path/release.jks
ECRM_RELEASE_STORE_PASSWORD=***
ECRM_RELEASE_KEY_ALIAS=release
ECRM_RELEASE_KEY_PASSWORD=***
```

执行 `./gradlew assembleRelease` 生成已签名 APK；执行 `./gradlew bundleRelease` 生成 AAB。两个 Release 任务在缺少完整签名参数时会明确失败，避免产出不能上架或安装的未签名包。后台“应用 → App → Android”只维护包名、版本、下载地址和证书 SHA-256 指纹，不能保存 keystore 或密码。

## 下一步

按 `docs/features/03-user-app.md` 和 `docs/openapi/app-*.yaml` 逐模块接入：认证 → 店铺/商品 → 购物车 → 下单支付 → 订单售后。支付、库存、优惠与售后需先对齐 `docs/api/FUNCTIONAL-TRUTH.md` 的状态机与幂等规则。
