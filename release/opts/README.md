# release/opts

可选公网 / TLS 入口。业务静态站已在 Docker：compose project `qixi_mergers` 下的 admin / h5 / … nginx 容器。

| 路径 | 说明 |
| --- | --- |
| `nginx/qixi-mergers.local.conf` | 本机可选（默认直接用容器宿主口 18081–18086） |
| `nginx/qixi-mergers.prod.conf.example` | 生产 TLS：反代到容器 IP 或宿主端口 |

证书目录可放 `opts/<domain>/`（勿提交私钥）。详见 `docs/release/SERVICE-MATRIX.md`。
