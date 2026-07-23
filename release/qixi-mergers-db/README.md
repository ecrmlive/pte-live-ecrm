# qixi-mergers-db

本目录 **不再运行任何容器**（无 MySQL / Redis / NATS / etcd / MinIO）。

| 组件 | 归属 |
| --- | --- |
| MySQL / Redis / NATS / etcd | **pte-live-im**（`db/` + `mq/`，网络 `pte_live_net`） |
| 对象存储 | **腾讯云 COS**（见 `api/conf/*.yaml` 的 `cos:`） |

`pack db` 仅把仓库 `sql/` 同步到本目录，便于部署机灌库。  
业务网络 `qixi_mergers_net` 由 release 脚本按需 `docker network create`。
