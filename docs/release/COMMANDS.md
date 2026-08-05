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

仅在这些 YAML 中填写密码、JWT 和 pte-live-im 受控 API 凭证。COS、平台支付、小程序密钥统一写入被 Git 忽略的 `sql/admin/init_key.sql`；`sql/business/init_key.sql`、`sql/merchant/init_key.sql` 只保留各库边界说明。`make local-db-init` / `make test-db-init` 会按三库固定顺序自动导入；local 与 test 必须同步同一份三文件密钥 SQL。`local` 与 `test` 的 C 端 JWT 必须一致，统一后台 JWT 必须一致；店铺 JWT 独立。不得使用 `.env`、`jwt.env`、环境变量注入或把真实值提交到 Git。

## 构建与启动

```bash
make pack-backend       # 本机构建 Linux 二进制到 release/；不在 Docker/服务器编译
make local-infra        # 校验 pte-live-im 共享 MySQL、Redis、etcd、NATS 已运行
make local-db-init      # 直接在 pte_live_mysql 导入 admin/business/merchant 三库 SQL（可重复执行）
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
| Compose project | `pte_live_ecrm` |
| 共享网络 | `pte_live_net` / `172.30.0.0/24` |
| 共享 etcd | Compose 项目 `pte_live_mq`：`pte_live_etcd1`、`pte_live_etcd2`、`pte_live_etcd3` / `.12`、`.14`、`.15` |
| 共享 NATS | Compose 项目 `pte_live_mq`：`pte_live_nats1`、`pte_live_nats2`、`pte_live_nats3` / `.13`、`.16`、`.17` |
| 共享 MySQL / Redis | Compose 项目 `pte_live_db`：`pte_live_mysql` / `.10`，`pte_live_redis` / `.11` |
| SQL 初始化 | 不创建容器；`make local-db-init` 直接执行到 `pte_live_mysql` |
| API | `pte_live_ecrm_api_platform`、`pte_live_ecrm_api_business`、`pte_live_ecrm_api_merchant` / `.61`–`.63` |

七禧只映射三个 API 到宿主机 `18081`、`18082`、`18083`；不映射也不启动任何重复基础设施。容器内使用共享 `pte_live_*` 名称。七禧业务仍只使用独立的 `qixi_crm_admin`、`qixi_crm_business`、`qixi_crm_merchant` 三个库，绝不读写 pte-live-im 的数据库或表。

七禧 API 加入 `pte_live_net`，但绝不复用 pte 容器名、IM 数据库或 IM 表；pte-live-im 的 IM 规则仍由其自身仓库维护。

`make local-compose-check` / `make test-compose-check` 只做 Compose 结构校验；无需在七禧 YAML 中复制 MySQL 根密码。
