# 技术栈速查

完整说明见 `docs/architecture-target.md`。

## 后端

- Go 1.26.5
- Gin + GORM + Swagger
- MySQL 8.4.10、Redis 8.8.0
- etcd v3.7.0、NATS 2.12.0-alpine
- 配置：`.yaml`

## 前端

- 平台 / 商户管理后台：Vben 5+
- 用户端 H5/小程序：uni-app x（UTS / HBuilderX 5.0+）
- 用户端 PC：Vue 3 + Vite + TypeScript（`app-web/`）
- 可选原生端规范见 `AGENTS.md`

## 部署

- 本机 pack 到 `release/<service>/`；后端通过同名 `docker-compose.yaml` 运行
- `config/local` 与 `config/test` 分离；两环境共用 `qixi_mergers` 容器名和固定 IP，不能并行运行
- 禁止服务器源码构建；Dockerfile 只复制产物
- 遵循全局 Skill `unified-docker-release`
- **MySQL / Redis / NATS / etcd**：由 `pte-live-im` 的 `db/` + `mq/` 启动（`pte_live_net`）
- **对象存储**：腾讯云 COS（`api/conf` 的 `cos:`）；本仓无 MinIO

## 建议目录（落地后）

```text
api-platform/  api-business/  api-merchant/  job/  admin-platform/  admin-merchant/  app-web/  app-mp/  app-ios/  app-adnroid/  app-harmony/  sql/  release/  docs/
```
