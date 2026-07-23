# 接入 pte-live-im（客服 IM）

状态：**本仓 remote 桥 + WS 建连已编码**；消息正文仍走本仓 `/cs`（IM 侧强制 E2EE，完整密聊需 SDK）。  
对端改动在 `~/Documents/GitHub/pte-live-im`（mall S2S，默认 `app_id=30001`）。

---

## 1. 怎么接

| 项 | 做法 |
| --- | --- |
| 角色 | **本仓 = 业务真相**（会话/订单卡片/`mer_id`）；**pte-live-im = UserSig / C2C / WS** |
| 场景 | `scene=chat` C2C；**不要** `scene=shop` |
| S2S | `POST /api/v1/integrations/mall/usersig` · `…/conversation/open-single` |
| Header | `X-Pte-Mall-Integration-Token`（仅服务端；env 注入） |
| 数值 ID | app：`1_000_000_000+uid`；坐席：`2_000_000_000+service_id` |
| 前端 | service-web「连接 IM WS」：`/ws?sdkAppID&identifier&userSig` + `login` envelope |

---

## 2. 本仓配置

`api/conf/{admin,app}.yaml`：

```yaml
im:
  mode: remote                 # local | remote
  api_base: "http://pte_live_api_im:11504"   # 容器内；本机进程可用 http://127.0.0.1:11504
  ws_public_url: "ws://127.0.0.1:11510/ws" # 浏览器连宿主机发布口
  app_id: "30001"
  integration_token: ""        # 建议用环境变量，勿提交真实值
```

环境变量覆盖：`QIXI_IM_MODE` · `QIXI_IM_API_BASE` · `QIXI_IM_WS_PUBLIC_URL` · `QIXI_IM_INTEGRATION_TOKEN`。

IM 侧（pte-live-im）：

| Env | 说明 |
| --- | --- |
| `PTE_MALL_INTEGRATION_ENABLED` | `true` |
| `PTE_MALL_INTEGRATION_APP_ID` | 默认 `30001` |
| `PTE_MALL_INTEGRATION_TOKEN` | 与本仓 `integration_token` 一致 |

SQL：`042`（会话表）+ `044`（`im_conversation_id` / `im_user_num`）。

Compose：`api-admin` / `api-app` / `job` 已附加外部网络 `pte_live_net`（需先起 pte-live-im）。

共享基建（必须先起 IM）：MySQL=`pte_live_mysql`、Redis=`pte_live_redis`、NATS=`pte_live_nats1..3`、etcd=`pte_live_etcd1..3`。  
本仓不再启动独立 mysql/redis/nats/etcd/MinIO；对象存储用腾讯云 COS。建库脚本：`sql/000_shared_im_mysql_bootstrap.sql`。

---

## 3. BFF 接口

| 端 | 路径 | 说明 |
| --- | --- | --- |
| app | `POST /api/app/v1/cs/threads` | 开会话；remote 时 S2S `open-single` |
| app/service | `GET …/im/credential?thread_id=` | UserSig + `ws_url` + `sdk_app_id` + `im_conversation_id` |
| service | 原有 threads/messages/查单/快捷回复 | 文本仍本仓 REST |

验收（remote）：

1. IM 与商城 token 对齐，`mode=remote`。  
2. C 端开线程 → `im_conversation_id>0`。  
3. service-web 点「连接 IM WS」→ `login_ack`。  
4. 商户隔离：不同 `mer_id` 会话不串。

---

## 4. 明确边界

| 不做（本刀） | 原因 |
| --- | --- |
| 浏览器内完整 E2EE 发消息 | 需 pte-im-sdk；IM 强制 E2EE |
| IM 直连 `qixi_*` | 边界 |
| 客服用 `scene=shop` | 直播弹幕语义 |

P1：H5/uni 接入 SDK 密聊；订单卡片自定义 payload；坐席在线路由。
