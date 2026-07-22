# app-uni — 用户端（H5 / 微信小程序）

一套 **uni-app Vue3 + Vite + TypeScript** 工程，同时产出：

| 端 | 命令 | 产物 |
| --- | --- | --- |
| H5 | `npm run dev:h5` / `npm run build:h5` | `dist/build/h5` → `make local-h5` → Nginx `:18083` |
| 微信小程序 | `npm run dev:mp-weixin` / `npm run build:mp-weixin` | `dist/build/mp-weixin` → 导入微信开发者工具 |

API 前缀：`/api/app/v1`（反代 `api-app` `:18085`）。

## 阶段进度（对齐 `docs/dev-plan-full.md`）

- **阶段 1**：登录 / 注册、token 持久化、请求封装、「我的」空壳
- **阶段 2**：首页 / 分类 / 商品列表 / 商详 / 店铺首页（只读，联调 stub）
- **阶段 3**：购物车 / 地址 / v2 结算支付（页已占位）

登录走 `api-app` identity（JWT：`access_token` / `refresh_token`）。商品首页/分类暂为内存 stub。  
演示账号需库内已有 C 端用户；可在登录页「注册并登录」。

## 本地开发

```bash
# 终端 1：C 端 API
make local-api-app

# 终端 2：H5
cd app-uni && npm install && npm run dev:h5

# 微信小程序
cd app-uni && npm run dev:mp-weixin
# 用微信开发者工具打开 dist/dev/mp-weixin
# 详情 → 不校验合法域名（本地联调）
```

小程序默认请求 `http://127.0.0.1:18085`（见 `src/config/env.ts`）。真机调试请改成可访问的局域网 IP，并在公众平台配置 request 合法域名。

`manifest.json` → `mp-weixin.appid` 填入正式小程序 AppID 后再上传。

## 发布 H5

```bash
make local-h5   # 内部执行 npm run build:h5，同步到 release/qixi-mergers-h5/dist
```

## 说明

目标栈文档写的是 uni-app x（UTS / HBuilderX）。当前仓库用 **CLI 可构建的 uni-app Vue3**，保证 `make local-h5` 与微信开发者工具闭环；后续若切 uvue，页面 IA 与 API 封装可复用。
