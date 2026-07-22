# release/opts

宿主机运维配置（**唯一 Nginx 层**）。本仓库**不使用** Docker Nginx / `qixi-mergers-gateway`。

| 路径 | 说明 |
| --- | --- |
| `nginx/qixi-mergers.local.conf` | 本机：后台站 → `:18080`（api-admin）；H5 `:18083` / PC `:18086` → `:18085`（api-app） |
| `nginx/qixi-mergers.prod.conf.example` | 生产模板：域名 / TLS / root 指向 `RELEASE_BASE_DIR` 下 dist |

## 约定

1. 前端只 pack `dist/` 到 `release/qixi-mergers-{admin,merchant-admin,h5,pc,service-web}/`，**无** frontend compose。
2. 反向代理与静态托管一律宿主机 Nginx；配置本机改完再 `make update-nginx`（rsync opts）覆盖服务器。
3. 证书目录可放 `opts/<domain>/`（勿提交私钥到 git）。

详见 `docs/release/SERVICE-MATRIX.md`、`PACK-AND-CONFIG.md`。
