# 店铺后台验收台账（stub）

> 本台账记录 **admin-merchant** 侧已实施证据；不以页面数量或 mock 假成功替代验收。  
> 关闭标准同 [ADMIN-PLATFORM-ROLE-ACCEPTANCE.md](./ADMIN-PLATFORM-ROLE-ACCEPTANCE.md)：真实 Vben 页面、真实 API、店铺 JWT + `X-AppId` 隔离、状态机/幂等测试、同状态截图。

**总状态：迁移中 / 待运行时闭环**

---

## 示范：退款订单（Wave 4 布局）

| 项目 | 状态 | 说明 |
| --- | --- | --- |
| 页面 | 迁移中 | `admin-merchant/src/views/ecrm/order/refund.vue` 已改用 `EcrmListPage` |
| API | 部分 closed | 同意/拒绝/确认收货/备注/隐藏/导出 — 见平台运行时「商户退款备注与受控隐藏」「退货退款同次跨服务闭环」 |
| 布局截图 | 待补 | `docs/acceptance/screenshots/merchant/order-refund/{crmeb,qixi}.png` |
| FIELD-CHECKLIST | 待填 | 复制 [FIELD-CHECKLIST-TEMPLATE.md](./FIELD-CHECKLIST-TEMPLATE.md) |

---

## 店铺结算

| 项目 | 状态 | 证据 / 缺口 |
| --- | --- | --- |
| 页面 | 迁移中 | `admin-merchant/src/views/ecrm/finance/settlement.vue` |
| API / 状态机 | 部分 closed | 平台运行时「商户结算审核与打款」「商户结算订单事实链」；商户侧 `api-merchant/internal/merchant/nativesettlement/`, `internal/event/merchantsettlement/` |
| 店铺 JWT 隔离 | 待运行时闭环 | 需 `mer_id` / `store_id` 越权拒绝测试 + Vben 截图 |
| 布局 / 截图 | 待补 | slug: `finance-settlement` |

---

## 好友助力

| 项目 | 状态 | 证据 / 缺口 |
| --- | --- | --- |
| 页面 | 迁移中 | `admin-merchant/src/views/ecrm/marketing/assist.vue`（registry: `/marketing/assist`） |
| API | 待运行时闭环 | 平台侧 `api-platform/internal/platform/assist/handler_integration_test.go` 仅覆盖监管；商户写操作待逐项 HTTP + 中文夹具 |
| 布局 / 截图 | 待补 | slug: `marketing-assist` |

---

## 拼团

| 项目 | 状态 | 证据 / 缺口 |
| --- | --- | --- |
| 页面 | 迁移中 | `admin-merchant/src/views/ecrm/marketing/combination.vue` |
| API | 待运行时闭环 | C 端/商户拼团状态机见 FUNCTIONAL-TRUTH §8；专用 integration test 待补 |
| 布局 / 截图 | 待补 | slug: `marketing-combination` |

---

## 核销（到店/虚拟履约）

| 项目 | 状态 | 证据 / 缺口 |
| --- | --- | --- |
| 页面 | 迁移中 | `admin-merchant/src/views/ecrm/order/verify.vue` |
| API | 待运行时闭环 | 店员核销权限与订单状态边界待 isolation HTTP + Vben |
| 布局 / 截图 | 待补 | slug: `order-verify` |

---

## 积分（integral / 积分商城）

| 项目 | 状态 | 证据 / 缺口 |
| --- | --- | --- |
| ecrm 独立页 | 迁移中 | `admin-merchant/src/views/ecrm/setting/integral.vue`（registry: `/setting/integral-policy`、`/marketing/integral/config`） |
| 平台监管参考 | 部分 closed | `api-platform/internal/platform/points/handler_integration_test.go`；C 端 `api-business/internal/business/points/handler_integration_test.go` |
| 商户积分抵扣策略 | 待运行时闭环 | `admin-merchant/src/api/core/merchant-integral.ts` + 店铺 JWT 隔离 HTTP 验收 |
| 布局 / 截图 | 待补 | slug: `setting-integral` |

---

## 关联文档

- 布局保真：[LAYOUT-FIDELITY-CHECKLIST.md](./LAYOUT-FIDELITY-CHECKLIST.md)
- 生产门禁：[HIGH-RISK-PRODUCTION-GATE.md](./HIGH-RISK-PRODUCTION-GATE.md)
- 平台五角色：[ADMIN-PLATFORM-ROLE-ACCEPTANCE.md](./ADMIN-PLATFORM-ROLE-ACCEPTANCE.md)
