# 接入 pte-live-im（客服 IM）

状态：**remote + Web SDK E2EE 发消息**。本仓 `/cs` 不再传文本正文，仅保留会话元数据（订单卡片、`mer_id` 等）。

对端：`pte-live-im`（mall S2S，默认 `app_id=30001`）· Web SDK：`pte-live-im-sdk/packages/im-web-sdk`（service-web 内置 bundle `pte-im-web-sdk.js`）。

---

## 1. 怎么接

| 项 | 做法 |
| --- | --- |
| 角色 | **本仓 = 业务真相**（会话/订单卡片/`mer_id`）；**pte-live-im = UserSig / C2C / E2EE** |
| 场景 | `scene=chat` C2C；**不要** `scene=shop` |
| S2S | `POST /api/v1/integrations/mall/usersig` · `…/conversation/open-single` |
| Header | `X-Pte-Mall-Integration-Token`（仅服务端；env 注入） |
| 数值 ID | app：`1_000_000_000+uid`；坐席：`2_000_000_000+service_id` |
| 前端 | `PteLiveIMWebClient`：UserSig → `start()` → `sendText(im_conversation_id, text)` |

---

## 2. 本仓配置

`api/conf/{admin,app}.yaml`：

```yaml
im:
  mode: remote                 # 仅 remote；local 假 UserSig 已移除
  api_base: "http://pte_live_api_im:11504"
  ws_public_url: "ws://127.0.0.1:11510/ws"
  app_id: "30001"
  integration_token: ""        # 必须配置，建议环境变量
```

环境变量：`QIXI_IM_MODE` · `QIXI_IM_API_BASE` · `QIXI_IM_WS_PUBLIC_URL` · `QIXI_IM_INTEGRATION_TOKEN`。

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
3. service-web「连接 IM」→ SDK `start`；发送走 `sendText`，不 POST 文本到 `/cs`  
4. 商户隔离：不同 `mer_id` 不串  

---

## 4. 边界

| 不做 | 原因 |
| --- | --- |
| IM 直连 `qixi_*` | 边界 |
| 客服用 `scene=shop` | 直播弹幕语义 |
| local 假 UserSig | 运营前只保留 remote |
