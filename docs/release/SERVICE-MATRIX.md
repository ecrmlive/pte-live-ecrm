# 服务、容器与网络矩阵

| 类别 | 服务 | 容器 | 固定 IP | 数据边界 |
| --- | --- | --- | --- | --- |
| 共享基础设施 | MySQL | `pte_live_mysql` | `172.30.0.10` | 七禧仅使用 `qixi_crm_admin`、`qixi_crm_business`、`qixi_crm_merchant` 三库 |
| 共享基础设施 | Redis | `pte_live_redis` | `172.30.0.11` | 七禧缓存命名空间 |
| 共享基础设施 | etcd | `pte_live_etcd1`、`pte_live_etcd2`、`pte_live_etcd3` | `.12`、`.14`、`.15` | 共享服务发现集群 |
| 共享基础设施 | NATS | `pte_live_nats1`、`pte_live_nats2`、`pte_live_nats3` | `.13`、`.16`、`.17` | 共享事件总线集群 |
| 初始化 | SQL 导入 | 无容器 | — | `make local-db-init` 直接导入共享 MySQL 中的七禧三库 |
| API | 统一后台 | `pte_live_ecrm_api_platform` | `172.30.0.61` | admin 身份、权限、后台配置与客服入口 |
| API | C 端业务 | `pte_live_ecrm_api_business` | `172.30.0.62` | C 用户、交易、营销、资产与内容消费 |
| API | 店铺 | `pte_live_ecrm_api_merchant` | `172.30.0.63` | 店铺、员工、商品经营、库存与履约 |
| 异步 | job | `pte_live_ecrm_job` | `172.30.0.64` | profile 保留，待旧单库迁移完成后启用 |

Compose project 与前缀固定为 `pte_live_ecrm`，七禧 API 固定加入 `pte_live_net (172.30.0.0/24)`。local/test 仅是不同执行宿主机，矩阵不得变更。

七禧仅映射 API 宿主机端口 `18081`、`18082`、`18083`。MySQL、Redis、etcd、NATS 端口由 pte-live-im 基础设施统一管理，七禧不新增映射。

`app-pc`、`app-uni`（H5/小程序）、`admin-platform` 与 `admin-merchant` 是本机构建的前端产物，不作为开发期 Nginx 容器启动。原生端 `app-ios`、`app-adnroid`、`app-harmony` 使用各自原生构建链路。pte-live-im 仍独立维护其 IM 数据库和表，但为七禧提供共享基础设施。
