# merchant-admin — 商户后台

Vue 3 + Vite + Ant Design Vue。对接 `api-admin` 前缀 `/api/merchant/v1`；JWT 强制 `mer_id`。

## 已落地

| 阶段 | 能力 |
| --- | --- |
| 1 | 登录 / 菜单壳 / 改密 |
| 2 | 商品列表、发布/编辑、上下架、库存（单规格默认 SKU） |

## 本地

```bash
pnpm install
pnpm dev   # :5174 → :18080
```

种子：`meradmin` / `admin123`
