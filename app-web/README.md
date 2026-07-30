# app-web — 用户端 PC Web 商城

Vue 3 + Vite + TypeScript。C 端 API 走 **api-business** `/api/app/v1`，与 H5、小程序及原生端使用同一 C 端 JWT 契约。

## 开发

```bash
cd app-web
pnpm install
pnpm dev          # http://127.0.0.1:5183 ，/api 反代 api-business :18082
```

## 打包

```bash
pnpm build        # 生成 dist/
```

本地开发不安装或启动 Nginx；通过 Vite 的 `127.0.0.1:5183` 访问。运行服务使用根目录唯一的 `docker-compose.yaml`。

## 阶段

| 阶段 | 内容 |
| --- | --- |
| 运行边界 | PC 仅通过 `api-business` 访问 C 端数据，不访问后台或店铺 API |
| 鉴权 | 仅发送 `Authori-zation: Bearer <token>` |
| 验收基线 | `docs/features/03-user-app.md` 与 `docs/CRMEB-FULL-FUNCTION-CHECKLIST.md` |
