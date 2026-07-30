# 服务、容器与网络矩阵

| 类别 | 服务 | 容器 | 固定 IP | 数据边界 |
| --- | --- | --- | --- | --- |
| 基础设施 | MySQL | `qixi_mergers_mysql` | `172.31.24.10` | `qixi_crm_admin`、`qixi_crm_business`、`qixi_crm_merchant` |
| 基础设施 | Redis | `qixi_mergers_redis` | `172.31.24.11` | 七禧独立缓存 |
| 基础设施 | etcd | `qixi_mergers_etcd` | `172.31.24.12` | 七禧独立服务发现 |
| 基础设施 | NATS | `qixi_mergers_nats` | `172.31.24.13` | 七禧独立事件总线 |
| 初始化 | SQL 导入 | `qixi_mergers_db_init` | `172.31.24.14` | 顺序导入三个库的五类 SQL |
| API | 统一后台 | `qixi_mergers_api_platform` | `172.31.24.21` | admin 身份、权限、后台配置与客服入口 |
| API | C 端业务 | `qixi_mergers_api_business` | `172.31.24.22` | C 用户、交易、营销、资产与内容消费 |
| API | 店铺 | `qixi_mergers_api_merchant` | `172.31.24.23` | 店铺、员工、商品经营、库存与履约 |
| 异步 | job | `qixi_mergers_job` | `172.31.24.24` | profile 保留，待旧单库迁移完成后启用 |

Compose project 与前缀固定为 `qixi_mergers`，网络固定为 `qixi_mergers_net (172.31.24.0/24)`。local/test 仅是不同执行宿主机，矩阵不得变更。

宿主机独立映射：MySQL `23306`、Redis `26379`、etcd `22379`、NATS `24222`；不得复用 pte-live 的 `13306`、`16379`、`12379` 或 `14222`。

`app-web`、`app-mp`（H5/小程序）、`admin-platform` 与 `admin-merchant` 是本机构建的前端产物，不作为开发期 Nginx 容器启动。原生端 `app-ios`、`app-adnroid`、`app-harmony` 使用各自原生构建链路。pte-live-im 完全独立，不在此矩阵中。
