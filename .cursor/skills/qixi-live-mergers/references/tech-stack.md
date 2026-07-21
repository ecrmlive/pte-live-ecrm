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
- 用户端：uni-app x（UTS / HBuilderX 5.0+）
- 可选原生端规范见 `AGENTS.md`

## 部署

- 本机 pack → rsync `release/<service>/` → 远程 Docker Compose 挂载
- `config/local`（本机）与 `config/prod`（服务器）分离
- 禁止服务器源码构建；Dockerfile 只复制产物
- 遵循全局 Skill `unified-docker-release`

## 建议目录（落地后）

```text
api/  admin/  merchant-admin/  app-uni/  sql/  release/  docs/
```
