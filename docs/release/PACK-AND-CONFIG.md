# Pack and Config

> 占位文档。与全局 Skill `unified-docker-release` 对齐；代码落地后填写本仓库实际路径。

## 铁律

1. 环境只有 `local` 与 `prod`：`config/local/`、`config/prod/`。
2. 配置源在源码树维护一份；pack 时同步到 local/prod。
3. 上传 rsync 整个 `release/<service>/`，排除 `config/local/**`。
4. 禁止服务器改 prod 业务配置；本机改完再 deploy 覆盖。
5. 后端挂载 `bin/` + `config/${QX_RELEASE_ENV}/app.yaml`；前端挂载 `dist/`。
6. 配置文件使用 `.yaml`（或 release 约定的 `app.yaml` / `compose.env`）。

## 预期 release 树

```text
release/
  qixi-mergers-db/
  qixi-mergers-api/
  qixi-mergers-admin/
  qixi-mergers-merchant-admin/
  qixi-mergers-h5/
  opts/
```

服务前缀以落地时 Makefile 为准。
