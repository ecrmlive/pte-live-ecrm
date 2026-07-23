# admin-platform — 栖息多商户平台后台（Vben 5.7）

基于 **Vben Admin 5**（Element Plus + `Page` / `useVbenModal` / `useVbenVxeGrid`）。

## 本地开发

```bash
cd admin-platform
pnpm install
pnpm -r --filter "./internal/**" --filter "./packages/**" --filter "./scripts/**" run stub
pnpm dev:platform
# http://localhost:5173/#/auth/login
```

API：`http://127.0.0.1:18080/api/platform/v1`（`api-admin`）。

演示账号：`admin` / `admin123`。

## 业务页

已迁移叶子在 `src/views/mergers/**`，菜单 path 映射见 `src/views/mergers/registry.ts`。

## 构建

```bash
pnpm run build:platform   # → dist/
pnpm run build:merchant   # 商户应用 → ../admin-merchant/dist/
```
