# Release Commands

矩阵：[SERVICE-MATRIX.md](./SERVICE-MATRIX.md)。

## 前置：IM 共享基建

MySQL / Redis / NATS / etcd **不在本仓启动**，由 `pte-live-im` 的 `db/` + `mq/` 负责（网络 `pte_live_net`）。

```bash
# 在 pte-live-im 仓库按该仓库文档启动 db + mq 后：
docker exec -i pte_live_mysql mysql -uroot -p"$MYSQL_ROOT_PASSWORD" \
  < sql/000_shared_im_mysql_bootstrap.sql
# 再按 sql/README.md 灌入业务迁移（宿主口仍为 127.0.0.1:13306）
```

## 本机

```bash
make init-env
make local-db                 # 同步 sql/ + 确保 qixi_mergers_net（无容器）
make local-api-admin
make local-api-app
make local-job
make local-admin              # … 等同 local-frontend 各端
make local-manager            # 店员端 :18087
make local-backend-all
make local-frontend-all
make local-compose-check
```

`make local-mq` / `deploy-mq-reload` 已废弃（请到 IM 仓库操作）。

## 部署（需授权）

```bash
make deploy-db-reload         # 同步 sql/ + 远程创建网络（无容器）
make deploy-api-admin
make deploy-api-app
make deploy-job
make deploy-admin / deploy-merchant-admin / deploy-h5 / deploy-pc / deploy-service-web / deploy-manager
make deploy-backend-all
make deploy-frontend-all
make deploy-all               # 业务 qixi_mergers，不含共享 db/mq
make update-nginx
```
