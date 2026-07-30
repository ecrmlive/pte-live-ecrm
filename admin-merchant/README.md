# admin-merchant — 七禧多商户商户后台（Vben 5.7）

Vben workspace 成员（根在 `admin-platform/`）。

```bash
cd admin-platform
pnpm install
pnpm -r --filter "./internal/**" --filter "./packages/**" --filter "./scripts/**" run stub
pnpm dev:merchant
# http://localhost:5174/#/auth/login
```

API：`http://127.0.0.1:18083/api/merchant/v1`。

演示账号：`meradmin` / `admin123`。
