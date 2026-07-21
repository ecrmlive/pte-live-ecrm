# Release Commands

> 占位文档。业务代码与 `release/` 目录落地后，按 `unified-docker-release` 补全真实命令。

## 预期命令（落地后）

```bash
make init-env
make init-env-prod
make local-compose-check
make deploy-db-reload
make deploy-backend-all
make deploy-frontend-all
make deploy-all
```

语义约定：

- `local-*`：本机 pack + compose，不上传
- `deploy-*`：pack + rsync 整目录 + 远程 compose
- `deploy-all`：仅业务服务，不含独立 DB/MQ 基建

在命令未实现前，不要假装已可部署；先补 `release/<service>/` 与脚本。
