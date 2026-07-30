# 本地与测试服务器命令

`local` 和 `test` 只表示命令所在的宿主机。二者必须运行完全相同的 Docker Compose、容器名、网络、固定 IP、数据库名和 YAML；不得在同一台主机并行运行。

## 初始化运行 YAML

```bash
make init-env-local
# 在测试服务器工作副本执行同一条：
make init-env-test
```

两条命令创建同一套未纳入 Git 的文件：

```text
release/config.yaml
release/config/api-platform/app.yaml
release/config/api-business/app.yaml
release/config/api-merchant/app.yaml
release/config/job/app.yaml
```

仅在这些 YAML 中填写密码、JWT、支付、COS 和 pte-live-im 受控 API 凭证。`local` 与 `test` 的 C 端 JWT 必须一致，统一后台 JWT 必须一致；店铺 JWT 独立。不得使用 `.env`、`jwt.env`、环境变量注入或把真实值提交到 Git。

## 构建与启动

```bash
make pack-backend       # 本机构建 Linux 二进制到 release/；不在 Docker/服务器编译
make local-infra        # MySQL、Redis、etcd、NATS
make local-db-init      # 导入 admin/business/merchant 三库 SQL（可重复执行）
make local-db-reset     # 删除并重建三库后按固定顺序导入 SQL（仅 local，破坏性操作）
make local-backend      # 启动三个 API
make local-ps
make local-down
```

测试服务器执行同名 `test-*` 目标，实际行为与 `local-*` 完全一致：

```bash
make test-infra
make test-db-init
make test-backend
```

命令不会安装或启动本机 Nginx。PC、H5、小程序本地开发使用自身 Vite/HBuilderX 服务并通过 `127.0.0.1` 访问。

## 固定 Compose 身份

唯一 Compose 文件是仓库根目录 [`docker-compose.yaml`](../../docker-compose.yaml)。

| 项 | 固定值 |
| --- | --- |
| Compose project | `qixi_mergers` |
| 网络 | `qixi_mergers_net` / `172.31.24.0/24` |
| MySQL | `qixi_mergers_mysql` / `172.31.24.10` |
| Redis | `qixi_mergers_redis` / `172.31.24.11` |
| etcd | `qixi_mergers_etcd` / `172.31.24.12` |
| NATS | `qixi_mergers_nats` / `172.31.24.13` |
| SQL 初始化 | `qixi_mergers_db_init` / `172.31.24.14` |
| API | `qixi_mergers_api_platform`、`qixi_mergers_api_business`、`qixi_mergers_api_merchant` / `.21`–`.23` |

宿主机端口为七禧独立映射：MySQL `127.0.0.1:23306`、Redis `127.0.0.1:26379`、etcd `127.0.0.1:22379`、NATS `127.0.0.1:24222`、三个 API 分别为 `18081`、`18082`、`18083`。这些端口不与 `pte-live` 共用；容器内服务仍只使用 `qixi_mergers_net` 中的固定容器名和 IP。

严禁加入 `pte_live_net`、复用 pte 容器名、数据库、Redis、NATS 或 etcd。pte-live-im 仅经 YAML 中配置的受控 API/SDK 集成。

`make local-compose-check` / `make test-compose-check` 只做 Compose 结构校验；必须先在 `release/config.yaml` 填写 MySQL 密码。
