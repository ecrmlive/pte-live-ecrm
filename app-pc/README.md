# app-pc — 用户端 PC 商城

Vue 3 + Vite + TypeScript。验收对齐功能表 4；API 走 **api-app** `/api/app/v1`。

## 开发

```bash
cd app-pc
pnpm install
pnpm dev          # http://127.0.0.1:5183 ，/api 反代 :18085
```

## 打包

```bash
make local-pc     # → release/qixi-mergers-pc/dist
```

宿主机 Nginx：`18086`（见 `release/opts/nginx/`）。

## 阶段

| 阶段 | 内容 |
| --- | --- |
| PC-0 | 工程 + release + Nginx |
| PC-1 | 登录 / 首页壳 / 个人中心壳 |
| PC-2+ | 随 `docs/dev-plan-full.md` 阶段 2–5 竖切 |
