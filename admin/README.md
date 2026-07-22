# admin — 平台后台

Vue 3 + Vite + Ant Design Vue。对接 `api-admin` 前缀 `/api/platform/v1`。

## 已落地（按开发计划）

| 阶段 | 能力 |
| --- | --- |
| 1 | 登录 / JWT / 菜单壳 / 改密 |
| 2 | 商户列表与启停、入驻审核、平台分类、品牌、商品审核 |

## 本地开发

```bash
# DB：000 + 001 + 002
# API：go run ./cmd/api-admin -config conf/admin.local.yaml
pnpm install
pnpm dev   # :5173 → 代理 :18080
```

种子：`admin` / `admin123`

## 打包

```bash
pnpm build
# 或 make local-admin
```
