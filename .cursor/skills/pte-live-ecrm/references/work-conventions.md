# 工作约定（原文）

以下约定与仓库 `AGENTS.md` 一致，执行时以 `AGENTS.md` 为准。

## 通用工作方式

- 制造的测试数据需要注意中文乱码问题
- 收到开发、修改的任务只做制定方案、开发功能、查找原因、提供解决方案、修复问题，验证和部署我会独立说
- 在项目工作区内自主完成常规开发、依赖安装、构建、测试、静态检查和诊断；不需要用户手动打开系统终端或 Codex 终端。
- 可以执行：go mod download、go test、go build、go vet、node、npm、pnpm、corepack、./gradlew、gradle、java、xcodebuild、swift、swift package、hvigor、ohpm、hpm。
- 构建缓存优先使用本机默认缓存；若受限环境无法写入系统缓存，可使用项目或临时目录缓存完成校验。不得因此要求用户手动开终端。
- 如系统弹出权限申请，优先申请最小且可复用的构建命令类别权限；不要申请无限制全盘访问。
- 可以执行 Docker Compose 的检查、日志、构建产物验证；部署、服务器操作必须由用户明确授权。
- 不允许删除系统文件，不使用破坏性命令。
- 未经用户明确授权，不允许 Git commit、push、创建 PR、重置分支或覆盖用户已有修改。
- 不输出、不写入、不提交密码、Token、密钥、证书或其他敏感信息到Git仓库。
- 修改代码前先读取相关项目文档、AGENTS.md 与对应 Skill；遵守项目 AGENTS.md 的优先规则。
- 配置文件使用 .yaml。

## 构建与部署

- 必须在本机构建打包产物，再上传服务器并使用 Docker Compose 挂载模式运行。
- 禁止把源代码上传到服务器打包。
- 禁止用 Dockerfile 在服务器构建应用；Dockerfile 仅可用于复制本机构建产物到镜像。
- 发布、打包、Docker、Nginx、服务器部署相关改动前，必须先读：
  - docs/release/COMMANDS.md
  - docs/release/PACK-AND-CONFIG.md
- 每次实际部署完成后，必须反馈部署服务器 IP。

## 固定技术版本

- Swift：6
- Go：1.26.5
- Node.js：24.18.0
- npm：11.16.0
- Corepack：0.35.0
- Docker 镜像：
  - alpine:3.24.1
  - nginx: 宿主机（不用 Docker Nginx）
  - node:24.18.0-alpine3.24
  - apachepulsar/pulsar:4.0.12
  - mysql:8.4.10
  - redis:8.8.0
  - mongo:8.0.26
  - gcr.io/etcd-development/etcd:v3.7.0
  - nats:2.12.0-alpine
  - 对象存储：腾讯云 COS（本仓不启 MinIO）
- API 技术栈：Gin、GORM、Swagger、MySQL、etcd、NATS。

## 多项目与客户端技术规范

- 新 iOS 项目、新 Android 项目、新鸿蒙项目、uni-app x：统一设计导航栏、高度44、紧贴状态栏底部、按钮44px
- 新 iOS 项目：iOS 16+、Swift 6、Swift/UIKit（禁止 SwiftUI）、MVVM/Clean Architecture、URLSession/Alamofire、URLSession WebSocketTask 或 Network Framework、Core Data、async/await/Actor、Factory/Swinject、OSLog、友盟 U-Push。
- 新 Android 项目：Android 12/API 31+、Kotlin、Jetpack Compose（兼容 View）、MVVM/Clean Architecture、OkHttp/Retrofit、OkHttp WebSocket、Room、Coroutines/Flow、Hilt/Koin、Timber、友盟 U-Push。
- 新鸿蒙项目：OpenHarmony API 23、ArkTS、ArkUI、MVVM/Repository、ohos.net.http、WebSocket API、RDB、Promise/TaskPool、Provider/手动 DI、HiLog、友盟 U-Push。
- uni-app x 使用 UTS/HBuilderX 5.0+。
- 管理后台使用 Vben 5+。
- 管理后台列表/抽屉布局：**强制**对齐店铺列表金标准 `admin-platform/src/views/ecrm/merchant/list.vue`（见 `docs/acceptance/LAYOUT-FIDELITY-CHECKLIST.md`），禁止新建 `EcrmListPage` 列表骨架。
- 管理后台时间：**强制** `yyyy-MM-dd HH:mm:ss`（Asia/Shanghai）；API/MySQL 时区见 `docs/release/PACK-AND-CONFIG.md`；前端用 `formatShanghaiDateTime`。
