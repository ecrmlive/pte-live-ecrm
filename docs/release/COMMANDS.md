# Release Commands

服务矩阵：[SERVICE-MATRIX.md](./SERVICE-MATRIX.md)。  
配置约定：[PACK-AND-CONFIG.md](./PACK-AND-CONFIG.md)。

## 语义

| 前缀 | 含义 |
| --- | --- |
| `local-*` | 本机 pack +（若有）compose，**不**上传 |
| `deploy-*` | pack + rsync +（若有）远程 compose |
| `deploy-all` | api-admin + api-app + job + 前端 dist，**不含** DB/MQ |
| `update-nginx` | rsync `release/opts` |

**API 分立**：`api-admin`（后台）与 `api-app`（C 端）为两个容器；无统一 `local-api`。

## 本机

```bash
make init-env
make local-db
make local-mq
make local-api-admin          # :18080
make local-api-app            # :18085
make local-job
make local-admin              # pack dist → 宿主机 Nginx :18081 → api-admin
make local-merchant-admin
make local-h5                 # → api-app :18083
make local-pc                 # → api-app :18086
make local-service-web        # → api-admin :18084
make local-backend-all
make local-frontend-all       # admin + merchant-admin + h5 + pc + service-web
make local-compose-check
```

## 部署（需授权）

```bash
make deploy-db-reload
make deploy-mq-reload
make deploy-api-admin
make deploy-api-app
make deploy-job
make deploy-admin
make deploy-merchant-admin
make deploy-h5
make deploy-pc
make deploy-service-web
make deploy-backend-all
make deploy-frontend-all
make deploy-all
make update-nginx
```

## 脚本

```text
scripts/qixi-release.sh
scripts/release/{pack,config,compose,bundle,lib}.sh
```
