# qixi-mergers-manager

店员端（核销 / 发货）静态前端，挂载进 compose project `qixi_mergers`。

- 源码：`app-manager/`
- 宿主端口：`18087`（IP `172.30.80.35`）
- 反代：`/api/` → `api-admin`（`.20:8080`）

```bash
make pack SVC=manager
make local-manager
```
