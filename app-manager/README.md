# app-manager · 店员端（阶段 6d）

核销 / 发货 / 代退最小竖切。API 前缀：`/api/manager/v1` → api-admin。

## 本地

```bash
cd app-manager
pnpm install
pnpm dev   # :5175，代理到 18080
```

演示账号：`staff1` / `admin123`（需执行 `sql/009_manager_service.sql`）。

## 能力

- 登录（`qixi_store_service`，需 `is_verify=1` 才能核销）
- 待核销列表 / 核销码查询 / 确认核销
- 待发货列表 / 填写物流发货
- 代退申请与同意（仅退款）
