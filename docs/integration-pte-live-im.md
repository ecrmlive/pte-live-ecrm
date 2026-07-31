# 接入 pte-live-im（客服 IM）

状态：**remote + Web SDK E2EE 发消息**。本仓 `/cs` 不再传文本正文，仅保留会话元数据（订单卡片、`mer_id` 等）。

对端：`pte-live-im`（mall S2S，默认 `app_id=30001`）· Web SDK：`pte-live-im-sdk/packages/im-web-sdk`（由 `app-uni`、`app-pc` 与 Vben 系统按需集成）。

---

## 1. 怎么接

| 项 | 做法 |
| --- | --- |
| 角色 | **本仓 = 业务真相**（会话/订单卡片/`mer_id`）；**pte-live-im = UserSig / C2C / E2EE** |
| 场景 | `scene=chat` C2C；**不要** `scene=shop` |
| S2S | `POST /api/v1/integrations/mall/usersig` · `…/conversation/open-single` |
| Header | `X-Pte-Mall-Integration-Token`（仅服务端；env 注入） |
| 数值 ID | app：`1_000_000_000+uid`；坐席：`2_000_000_000+service_id` |
| 前端 | `pte-live-im-sdk` uni-app x UTS：短期 UserSig → `createPteIMSDK(...).login()` → `start()`；由 `pte-im-uikit` 发送 E2EE 文本 |

---

## 2. 本仓配置

`api-platform/conf/app.yaml`、`api-business/conf/app.yaml` 与 `api-merchant/conf/app.yaml`：

```yaml
im:
  mode: remote                 # 仅 remote；local 假 UserSig 已移除
  api_base: "http://pte_live_api_im:11504" # 仅本仓服务端 S2S，禁止返回给浏览器
  api_public_url: "http://127.0.0.1:11504" # H5/小程序实际可访问的 IM API 地址
  ws_public_url: "ws://127.0.0.1:11510/ws"
  app_id: "30001"
  integration_token: ""        # 必须配置，建议环境变量
```

环境变量：`QIXI_IM_MODE` · `QIXI_IM_API_BASE` · `QIXI_IM_API_PUBLIC_URL` · `QIXI_IM_WS_PUBLIC_URL` · `QIXI_IM_INTEGRATION_TOKEN`。

`api_base` 是容器内 S2S 地址，`api_public_url` 与 `ws_public_url` 是返回到 H5/小程序的真实公网地址；两者不得使用 `pte_live_api_im` 等 Docker 服务名。未配置客户端公网地址时，`/cs/im/credential` 会明确拒绝发放凭证，避免前端拿到不可访问的内网地址。

JWT 以 [SYSTEM-ARCHITECTURE.md](./SYSTEM-ARCHITECTURE.md) 与 [JWT、店铺 AppId 与商户 IM SDK AppId 契约](./auth-store-appid-im-contract.md) 为准：PC/小程序/H5 共用 C 端 JWT；平台/商户/区域/客服/运营共用统一后台 JWT；店铺管理系统独立 JWT。所有 JWT HTTP 请求只使用 `Authori-zation: Bearer <token>`。它们都与 pte-live-im 的 UserSig / JWT 是不同协议，禁止互相复用。商户 IM 必须按该商户当前启用的 SDK AppId 签发，不能继续假设全局单 `app_id=30001`。

IM 侧（pte-live-im）：`PTE_MALL_INTEGRATION_ENABLED`、`PTE_MALL_INTEGRATION_APP_ID`、`PTE_MALL_INTEGRATION_TOKEN`。

---

## 3. BFF 接口

| 端 | 路径 | 说明 |
| --- | --- | --- |
| app | `POST /api/app/v1/cs/threads` | 开会话；remote 时 S2S `open-single` |
| app/service | `GET …/im/credential?thread_id=` | UserSig + `api_url` + `ws_url` + `im_conversation_id` |
| service | `POST …/threads/:id/messages` | **仅** `msg_type=order|system`；`text` 返回错误 |

验收：

1. `im.mode=remote` 且 token 对齐  
2. C 端开线程 → `im_conversation_id>0`  
3. C 端/客服端「连接 IM」→ 共享 SDK `start`；发送走 SDK E2EE，不 POST 文本到 `/cs`
4. 商户隔离：不同 `mer_id` 不串  
5. 凭证隔离：C 端只能用自己的 `thread_id` 换取凭证；客服必须先认领该会话

---

## 4. 边界

| 不做 | 原因 |
| --- | --- |
| IM 直连 `qixi_*` | 边界 |
| 客服用 `scene=shop` | 直播弹幕语义 |
| local 假 UserSig | 运营前只保留 remote |
