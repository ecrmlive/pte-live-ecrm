# 目标技术架构

本仓库重建栈，不以 CRMEB PHP/Swoole 为运行时。

## 1. 固定版本

| 组件 | 版本 |
| --- | --- |
| Go | 1.26.5 |
| Node.js | 24.18.0 |
| npm | 11.16.0 |
| Corepack | 0.35.0 |
| alpine | 3.24.1 |
| nginx | 1.30.3-alpine3.23 |
| node 镜像 | 24.18.0-alpine3.24 |
| MySQL | 8.4.10 |
| Redis | 8.8.0 |
| etcd | v3.7.0 |
| NATS | 2.12.0-alpine |
| MinIO（可选对象存储） | RELEASE.2025-10-15T17-29-55Z |

配置文件统一使用 `.yaml`。

## 2. 推荐仓库形态（落地时）

```text
qixi-live-mergers/
├── AGENTS.md
├── docs/
├── api/                    # Go：Gin + GORM + Swagger
│   ├── cmd/
│   ├── internal/
│   │   ├── platform/       # 平台后台 API
│   │   ├── merchant/       # 商户后台 API
│   │   ├── app/            # C 端 API
│   │   ├── domain/         # 订单/商品/结算等领域
│   │   └── pkg/
│   └── conf/
├── admin/                  # Vben 5+ 平台后台
├── merchant-admin/         # Vben 5+ 商户后台（可与 admin 同仓多入口）
├── app-uni/                # uni-app x 用户端
├── sql/
├── release/                # 本机构建产物 + Docker 挂载
└── scripts/
```

若拆多仓，以契约（OpenAPI / NATS 事件）为准，避免循环依赖。

## 3. API 技术选型

| 层 | 选型 |
| --- | --- |
| HTTP | Gin |
| ORM | GORM |
| 文档 | Swagger / OpenAPI |
| DB | MySQL 8.4 |
| 缓存/锁 | Redis |
| 配置/服务发现 | etcd |
| 消息 | NATS（异步：支付后续、通知、结算任务） |
| 认证 | JWT / Session 按端分离；平台与商户 RBAC 隔离 |

## 4. 客户端规范

| 端 | 规范 |
| --- | --- |
| 管理后台 | Vben 5+ |
| 新移动端（uni-app x） | UTS / HBuilderX 5.0+ |
| 导航栏（新 App 项目） | 高度 44，紧贴状态栏底部，按钮 44px |
| iOS（若原生） | iOS 16+、Swift 6、UIKit（禁止 SwiftUI）、MVVM/Clean |
| Android（若原生） | API 31+、Kotlin、Compose、MVVM/Clean |
| 鸿蒙（若原生） | OpenHarmony API 23、ArkTS/ArkUI |

本阶段优先：**Go API + Vben 双后台 + uni-app x 用户端**。

## 5. 部署铁律

1. **本机构建**打包产物，再上传服务器；禁止在服务器用源码构建业务。
2. Dockerfile **只复制本机产物**进镜像或挂载 `bin/` / `dist/`，禁止服务器 `docker build` 编译应用。
3. Docker Compose **挂载模式**运行；`config/local` 与 `config/prod` 分离。
4. 改动发布前必读：
   - `docs/release/COMMANDS.md`
   - `docs/release/PACK-AND-CONFIG.md`
5. 实际部署完成后反馈**部署服务器 IP**。
6. 部署/服务器操作必须由用户明确授权。

细节遵循全局 Skill `unified-docker-release` 与（落地后的）仓库 `docs/release/`。

## 6. 领域横切要求

- 多租户：所有商户域数据带 `merchant_id`，平台接口可跨租户，商户接口强制隔离。
- 订单：平台主单 + 商户子单（对齐 CRMEB `store_group_order` / `store_order` 概念）。
- 资金：余额/佣金/结算流水可追溯；提现审核状态机明确。
- 营销：价格计算单一入口，避免前台/后台两套算法。
- 中文测试数据：注意连接与文件编码（utf8mb4），避免乱码。

## 7. 数据库前缀

所有业务表使用前缀 `qixi_`（由 CRMEB `eb_` 映射）。详见 `docs/schema/`。
