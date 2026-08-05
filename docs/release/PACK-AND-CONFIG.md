# 打包与 YAML 配置

## 后端产物

```text
api-platform/  -> release/pte-live-ecrm-api-platform/bin/api-platform
api-business/  -> release/pte-live-ecrm-api-business/bin/api-business
api-merchant/  -> release/pte-live-ecrm-api-merchant/bin/api-merchant
job/           -> release/pte-live-ecrm-job/bin/job
```

每个服务都是独立 Go module，在本机构建静态 Linux 二进制；Compose 仅以只读卷挂载产物运行。禁止在 Docker 容器或服务器上传源码后构建。禁止回退为 `api/cmd/*` 单 module 构建。

`job` 的未支付关单任务已迁至 `qixi_crm_business`，仅连接 `databases.business`；启用前仍须完成真实 MySQL 并发回归与运行机 YAML 配置。

## 配置唯一来源

运行配置均为 YAML：

| 文件 | 用途 |
| --- | --- |
| `release/config.yaml` | 共享基础设施标识；不保存 MySQL 根密码 |
| `release/config/api-platform/app.yaml` | 统一后台 API |
| `release/config/api-business/app.yaml` | PC、小程序&H5、iOS、Android、鸿蒙 C 端 API |
| `release/config/api-merchant/app.yaml` | 店铺系统 API |
| `release/config/job/app.yaml` | 异步任务（仅 `databases.business`） |

所有运行 YAML 被 Git 忽略。仓库只保留不含真实值的 [`release/config.yaml.example`](../../release/config.yaml.example) 与各独立服务 `conf/app.yaml` 模板。COS、平台支付、小程序密钥不写入运行 YAML，也不写入 Git：统一保存于被忽略的 `sql/admin/init_key.sql`，`make local-db-init` / `make test-db-init` 自动导入到统一后台配置表；`sql/business/init_key.sql` 与 `sql/merchant/init_key.sql` 只保留各自数据库边界说明。local 与 test 必须使用同一份三文件密钥 SQL。配置加载器不再提供内置 JWT 或支付密钥，`jwt.secret` 缺失会拒绝启动。

三套目标 API 只允许读取各自所属库的 `databases.<scope>.dsn`：`api-platform → admin`、`api-business → business`、`api-merchant → merchant`。跨库数据只能通过受控 API 或 NATS 事件读模型同步，禁止共享单一 `mysql.dsn`。

JWT 规则：

- PC、小程序、H5、iOS、Android、鸿蒙使用同一 C 端 JWT。
- 平台、商户、区域、客服、运营使用同一统一后台 JWT。
- 店铺系统使用独立 JWT。
- local 与 test 在不同宿主机使用完全相同的上述配置值。

七禧 API 复用 `pte_live_net` 的 MySQL、Redis、etcd、NATS 容器；但三套 API 只读写自己的 `qixi_crm_*` 数据库，绝不访问 pte-live-im 的数据库或表。MySQL 初始化通过 `pte_live_mysql` 容器内已有根凭证执行，七禧不复制该密码。

## 时区（强制 Asia/Shanghai）

管理后台时间展示统一为 **`yyyy-MM-dd HH:mm:ss`**（上海时区）。部署与运行约定：

| 层 | 要求 |
| --- | --- |
| Compose API / job | `environment.TZ=Asia/Shanghai`（见根目录 `docker-compose.yaml`） |
| Go 进程 | `OpenMySQL` 启动时 `time.Local = Asia/Shanghai`（内嵌 `time/tzdata`，alpine 无需系统 zoneinfo） |
| MySQL DSN | `parseTime=True&loc=Asia%2FShanghai&time_zone=%27%2B08%3A00%27`；`scripts/qixi-crm.sh` 同步 DSN 时写入；运行时 `NormalizeShanghaiDSN` 再次强制 |
| MySQL 容器 | 共享 `pte_live_mysql` 须 `TZ=Asia/Shanghai`（或等价 `+08:00`）；`local-db-init` 会尽量 `SET GLOBAL time_zone='+08:00'` |
| 管理端前端 | `formatShanghaiDateTime`（`admin-platform` / `admin-merchant` 的 `#/utils/date-time`）；禁止直接 `toLocaleString` 或裸显 API ISO 串 |

Element 日期控件 `value-format` 使用 `YYYY-MM-DD HH:mm:ss`。
