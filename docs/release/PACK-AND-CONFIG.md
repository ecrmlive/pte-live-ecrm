# Pack and Config

权威表：[SERVICE-MATRIX.md](./SERVICE-MATRIX.md)。

## Compose 分组

| project | 内容 |
| --- | --- |
| （无 compose） | `qixi-mergers-db` 仅同步 `sql/`；网络由脚本创建 |
| `qixi_mergers` | API + Job + 前端 |
| （外部）IM db/mq | MySQL、Redis、NATS、etcd — 由 **pte-live-im** 启动 |
| （外部）腾讯云 COS | 对象存储 — `cos:` 配置 / `QIXI_COS_*` |

## Release 树

```text
release/
  qixi-mergers-db/           # 仅 sql/（无容器；历史 nats/ 可忽略）
  qixi-mergers-api-admin/
  qixi-mergers-api-app/
  qixi-mergers-job/
  qixi-mergers-admin/        # pack key admin ← 源码 admin-platform/
  qixi-mergers-merchant-admin/  # pack key merchant-admin ← 源码 admin-merchant/
  qixi-mergers-h5/
  qixi-mergers-pc/
  qixi-mergers-service-web/
  qixi-mergers-manager/      # 店员端 app-manager
  opts/nginx/
```

## 铁律

1. `local` / `prod`；`QX_RELEASE_ENV`。
2. Go：`api/conf/{admin,app,job}.yaml` → 各 release `config/*/app.yaml`。DSN 指向 `pte_live_*`；对象存储为 `cos:`（密钥用 `QIXI_COS_*`，勿提交）。
3. rsync 整目录，排除 `config/local/**`。
4. 共享 db/mq 不进本仓 `deploy-all`；先起 IM，再 `deploy-db-reload`（sql+网络）与业务。
5. 全量菜单：`sql/043_crmeb_system_menu_full.sql`（可由 `python3 scripts/gen_crmeb_menu_sql.py` 从 CRMEB 安装包重生成）。
6. IM remote：`im.mode=remote` + `QIXI_IM_INTEGRATION_TOKEN`；对端 `PTE_MALL_INTEGRATION_*`（见 `docs/integration-pte-live-im.md`）。
7. 后台源码目录是 `admin-platform/`、`admin-merchant/`；pack key 仍为 `admin` / `merchant-admin`（≠ 源码目录名）。根目录不要再放 `admin`、`merchant-admin` 软链。
