# Pack and Config

权威服务/IP 表：[SERVICE-MATRIX.md](./SERVICE-MATRIX.md)。

## 铁律

1. `local` / `prod`；`QX_RELEASE_ENV` 选择。
2. Go 配置源码主文件（按进程）：
   - `api/conf/admin.yaml` → pack **覆盖** `release/qixi-mergers-api-admin/config/local/app.yaml`
   - `api/conf/app.yaml` → 同上 → `qixi-mergers-api-app`
   - `api/conf/job.yaml` → 同上 → `qixi-mergers-job`
   - `config/prod/app.yaml`：**仅首次**从 `*.example` 种子，后续 pack **不覆盖**（本机改密后保留）；模板见 `config/prod/app.yaml.example`
3. `pack db` 将仓库 `sql/` 同步到 `release/qixi-mergers-db/sql/`；compose 挂载 `./sql`（可 rsync 部署）。
4. rsync 整目录，排除 `config/local/**`。
5. **后台 API 与 C 端 API 进程分立**；领域代码可共享同一 `api/internal/domain`。
6. 前端只 `dist/`（含 `service-web`）；Nginx 仅宿主机。
7. 禁止服务器源码构建。

## Release 树

```text
release/
  qixi-mergers-db/
  qixi-mergers-mq/
  qixi-mergers-api-admin/    # 后台 API
  qixi-mergers-api-app/      # C 端 API
  qixi-mergers-job/
  qixi-mergers-admin/        # dist
  qixi-mergers-merchant-admin/
  qixi-mergers-h5/
  qixi-mergers-pc/           # PC 商城 dist
  qixi-mergers-service-web/
  opts/nginx/
```

## 网络固定 IP（业务）

| 容器 | IP |
| --- | ---: |
| `qixi_mergers_api_admin` | `.20` |
| `qixi_mergers_job` | `.21` |
| `qixi_mergers_api_app` | `.22` |
