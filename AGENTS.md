# AGENTS.md

本文件约束在 **qixi-live-ecrm** 内工作的 Agent（Cursor / Codex）。修改代码前必须先读本文档与对应 Skill。

## 项目是什么

多商户商城管理系统。功能基线对齐 CRMEB Merchant v4.0；技术实现按本仓库目标栈重建（Go + Vben + uni-app x），不把 PHP/Swoole 当作运行时。

文档入口：`docs/README.md`  
Skill：`.cursor/skills/qixi-live-ecrm/SKILL.md`（Codex 镜像：`codex-skills/qixi-live-ecrm/`）

外部只读参考源码（勿提交进本仓库）：

```text
~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0
```

## 优先级规则

1. 本文件 `AGENTS.md` 优先于通用习惯中的冲突项。
2. 仓库 Skill 与 `docs/` 次之。
3. 发布/打包改动前必须先读 `docs/release/COMMANDS.md` 与 `docs/release/PACK-AND-CONFIG.md`（文件齐全前先补文档再改脚本）。
4. 全局 `unified-docker-release` Skill 约束 release 形态。

## 任务边界

收到开发、修改的任务只做：制定方案、开发功能、查找原因、提供解决方案、修复问题。  
**验证和部署**由用户独立明确要求后再做。

## 通用工作方式

- 制造的测试数据需要注意中文乱码问题
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
  - nginx: 宿主机安装（本仓库不使用 Docker Nginx / nginx 镜像跑业务入口）
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

## 本仓库业务约束

- 功能验收以 `docs/features/` 为准；产品全景以 `docs/product-understanding.md` 为准。
- 对照 CRMEB 时只读外部源码：`~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0`。
- **数据库与表前缀固定为**：`qixi_crm_admin.qixi_crm_a_`（统一后台）、`qixi_crm_business.qixi_crm_b_`（C 端业务）、`qixi_crm_merchant.qixi_crm_m_`（店铺）。IM 表规则严格以 pte-live-im 仓库 `sql/init_im_schema.sql` 的实际定义为准，七禧不得重定义。禁止新代码使用 `qixi_m_*`、裸 `qixi_` 或 `eb_` 表前缀。见 `docs/SYSTEM-ARCHITECTURE.md`。
- 订单、支付、退款、库存、优惠、积分、佣金、商户结算为高风险域；改前先理清状态机与幂等。
- 商户域接口必须做 `mer_id` / `merchant_id` 隔离。
- 禁止把密钥、证书、真实手机号/身份证样例提交进库；测试数据用明显假数据且保持 utf8mb4 中文可读。

## 功能基线与写代码

先读 `docs/analysis-completeness.md`。

- **功能基线已锁定**（2026-07-21）：允许技术方案与业务编码。
- 验收以 `docs/features/` 为准；接口/状态机对照 `docs/api/FUNCTIONAL-TRUTH.md`；系统、JWT、数据库与目录边界以 `docs/SYSTEM-ARCHITECTURE.md` 为准。
- 高风险域（订单/支付/退款/库存/券/积分/佣金/结算）改前先对齐状态机与幂等。
- 验证与部署仅在用户明确要求时执行。

## 建议工作流

1. 读 `docs/analysis-completeness.md` 与 `docs/overview.md`。
2. 读本 Skill 的 references（功能图 / 技术栈 / 工作约定）。
3. 需要对照原系统时再打开 `~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0` 对应 Repository / 菜单 SQL。
4. 出方案 → 改代码 →（用户要求时）测试/部署。
