# Wave 5 — 高风险域生产门禁

> 隔离运行时验收见 [ADMIN-PLATFORM-RUNTIME-2026-08-04.md](./ADMIN-PLATFORM-RUNTIME-2026-08-04.md)。**本门禁**在上线前额外要求：状态机与 [FUNCTIONAL-TRUTH.md](../api/FUNCTIONAL-TRUTH.md) 一致、幂等/并发测试可复核、支付与 IM 凭据受控、跨库投影经 outbox/NATS 闭环。

**Wave 5 总状态：`open`（未关闭）** — 下列任一门禁项未在生产等价环境验证，不得宣称高风险域「生产就绪」。

---

## 1. 状态机对照（FUNCTIONAL-TRUTH）

| 域 | 规范章节 | 关键路径 / 约束 |
| --- | --- | --- |
| 订单创建/支付 | §1、§5 | `POST /api/v2/order/check|create`；支付回调 `POST /api/callback/v1/pay/wechat`；sandbox 仅 `pay/mock` 显式开启 |
| 退款 | §6 | `0→1→2→3` / 仅退款 `refunding→refunded`；微信终态 `POST /api/callback/v1/refund/wechat`；禁止审核即标已退款 |
| 用户分销提现 | §7.0 | `applied→approved→paid`；`payout_idempotency_key` + 内部凭证；不存收款账户 |
| 商户结算 | §7.1 | `bill_pending→bill_frozen→withdraw_applied→approved→paid`；平台只读投影 + NATS command |
| 订单应计/冲销 | §7.2 | 确认收货 `accrue`；退款终态 `reverse`；冻结账期不可回写 |
| 库存 | §6 退款联动、活动库存 | 退款终态触发商户库存账本回补；活动库存见 stock restore 测试 |
| 优惠券/计价 | §8 | 店铺券/平台券/SVIP/积分互斥；v2 check/create 为准 |
| 积分商城 | §1 积分线 | `/api/order/v3/*` 与 `activity_type=20` 快照不可被监管改价篡改 |
| 佣金/分销 | 运行时 + §7 | 推广员资格变更不重算历史佣金账本 |

---

## 2. 幂等与并发测试（仓库内可 grep 路径）

下列为**已存在**的自动化测试包/用例（`go test` 可复核）。集成测试多数需 `ECRM_*_TEST_DSN` 或临时 MySQL/NATS；未列出的域表示仍缺 dedicated integration test。

### 2.1 订单 / 支付 / 幂等

| 用例 | 路径 |
| --- | --- |
| 创建幂等键校验 | `api-business/internal/business/order/idempotency_test.go` — `TestValidIdempotencyKey`, `TestCreateRejectsMissingIdempotencyKeyBeforeDatabaseAccess` |
| Mock 支付生产门禁 | `api-business/internal/business/order/mock_guard_test.go` — `TestProductionHandlerRejectsMockPayment`, `TestProductionCallbackDoesNotExposeMockPayment`, `TestProductionCallbackDoesNotExposeMockRefund` |
| 确认收货 → 结算应计 outbox | `api-business/internal/business/order/settlement_integration_test.go` — `TestConfirmReceiptIntegrationEnqueuesFactAnchoredSettlementAccrual` |
| 微信支付签名/回调 | `api-business/internal/pkg/wechatpayv3/client_test.go`, `refund_test.go` — `TestNativePrepaySignsAndReturnsCodeURL`, `TestRefundSignsServerSideRequest`, `TestVerifyAndDecryptRefundCallback` |

### 2.2 退款 / 库存 / 结算冲销

| 用例 | 路径 |
| --- | --- |
| Sandbox 退款回调 + 库存/结算 outbox | `api-business/internal/business/order/refund_callback_integration_test.go` — `TestSandboxRefundCallbackIntegrationRecordsServerDerivedReversal` |
| 用户退货寄回 HTTP | `api-business/internal/business/refund/return_shipment_integration_test.go` — `TestReturnShipmentHTTPIntegration` |
| 商户确认收货 HTTP | `api-merchant/internal/merchant/nativerefund/confirm_return_integration_test.go` — `TestConfirmReturnHTTPIntegration` |
| 退款状态/金额单元 | `api-business/internal/business/refund/handler_test.go`, `api-business/internal/business/refundprocessor/processor_test.go` |
| 平台退款 RBAC/区域 | `api-platform/internal/platform/nativerefund/handler_integration_test.go` — `TestRefundHTTPRBACAndMerchantRegionScopes` |
| 活动库存回补 | `api-business/internal/domain/trade/stock_restore_test.go` — `TestRestoreActivityStock_PresellAndAssist`（`api-platform`、`api-merchant`、`job` 同路径镜像） |
| 商户库存命令校验 | `api-merchant/internal/event/merchantstock/command_test.go` — `TestStockCommandValidationAndFailClosedDatabase` |
| 结算命令幂等键 | `api-business/internal/event/merchantledger/outbox_test.go` — `TestSettlementCommandRequiresFactAnchoredIdempotencyKey` |

### 2.3 商户结算 / 提现 / 保证金

| 用例 | 路径 |
| --- | --- |
| 结算账本冻结账期不可变 | `api-merchant/internal/event/merchantledger/command_integration_test.go` — `TestSettlementLedgerIntegrationKeepsFrozenBillImmutable` |
| 结算命令重放/校验 | `api-merchant/internal/event/merchantsettlement/command_test.go` — `TestSettlementCommandValidationAndFailClosedDatabase`, `TestSettlementCommandReplayRequiresSameActionPayloadAndOperator` |
| 平台结算投影状态 | `api-platform/internal/platform/nativesettlement/handler_test.go` — `TestSettlementStatusAllowsOnlyProjectionStates`, `TestSettlementCommandValidation` |
| 平台结算投影事件 | `api-platform/internal/event/merchantsettlement/merchant_settlement_test.go` — `TestSettlementProjectionRejectsUnknownState` |
| 三库结算夹具契约 | `api-platform/internal/domain/business_sql_fixture_contract_test.go` — `TestMerchantSettlementLedgerFixtureContract`, `TestAdminMerchantSettlementProjectionFixtureContract` |
| 用户提现打款幂等 | `api-platform/internal/platform/nativewithdraw/handler_integration_test.go` — `TestWithdrawHTTPRBACAndPayoutIdempotency` |
| 商户保证金打款/扣减 | `api-platform/internal/platform/merchantdeposit/handler_integration_test.go` — `TestMarkPaidIntegrationIdempotencyAndBalance`, `TestDeductIntegrationIdempotencyFingerprintAndCentBoundary`, `TestReviewIntegrationConcurrencyAndImmutableLedger` |
| 分账申请并发审核 | `api-platform/internal/platform/profitsharing/handler_integration_test.go` — `TestProfitsharingReviewIntegrationConcurrencyAndAudit` |

### 2.4 优惠券 / 积分 / 佣金

| 用例 | 路径 |
| --- | --- |
| 下单券 ID 互斥 | `api-business/internal/business/order/coupon_test.go` — `TestUniqueCouponIDsAllowsPlatformAndStorePair`, `TestUniqueCouponIDsRejectsDuplicateOrZero` |
| 用户券状态视图 | `api-business/internal/business/coupon/handler_test.go` — `TestUserCouponViewPreservesUsedAndExpiredState` |
| 计价无券基线 | `api-business/internal/domain/trade/pricing_test.go` — `TestApplyPricing_NoCoupon` |
| 积分商品报价/下架 | `api-business/internal/business/points/handler_integration_test.go` — `TestQuoteRejectsOffSaleProduct` |
| 平台积分监管 RBAC/快照 | `api-platform/internal/platform/points/handler_integration_test.go` — `TestPointsHTTPRBACVersionAndOrderSnapshot` |
| 推广员/佣金边界 | `api-platform/internal/platform/userlist/status_integration_test.go` — `TestPromoterIntegrationRBACAtomicityAndCommissionBoundary` |
| 分销佣金标记 | `api-business/internal/business/distribution/handler_test.go` — `TestCommissionMark` |
| 平台用户资产调账校验 | `api-platform/internal/platform/userlist/handler_test.go` — `TestValidAdjustment`, `TestAdjustmentMinorUsesFixedCents` |

### 2.5 跨服务 / NATS / 入驻

| 用例 | 路径 |
| --- | --- |
| 入驻 NATS request/reply | `api-merchant/internal/event/merchantonboarding/onboarding_integration_test.go` — `TestStartIntegrationNATSRequestReplyChineseProvisioning` |
| 平台入驻客户端契约 | `api-platform/internal/event/merchantonboarding/client_integration_test.go` — `TestProvisionIntegrationNATSRequestReplyContract` |
| 入驻并发投影 | `api-platform/internal/domain/merchant/onboarding_integration_test.go` — `TestFinalizeIntentionApprovalIntegrationConcurrentProjection` |

---

## 3. 微信 / 支付宝凭据

| 规则 | 说明 |
| --- | --- |
| 测试密钥 | `wechatpayv3` 包内 `testKeys(t)` 生成临时 RSA/AES，仅用于 `client_test.go` / `refund_test.go` 的 httptest |
| 禁止入库 | 不得提交商户私钥、APIv3 密钥、支付宝应用私钥、真实商户号到 Git |
| 配置门禁 | `api-platform/internal/domain/content/mall_setting_config_test.go` — `TestPayConfigRejectsSensitiveKeys`；支付配置不得 echo 旧密钥 |
| 生产 | 真实商户凭据仅通过部署侧 secret 注入；未经明确授权不得在 CI/文档/夹具中使用 |
| Mock 渠道 | sandbox `pay/mock`、`refund/mock` 必须在生产路由不可见（见 `mock_guard_test.go`） |

---

## 4. 客服：统一后台 JWT + pte-live-im 多 AppId

| 项 | 要求 |
| --- | --- |
| 鉴权 | 客服队列/转接/设置走统一后台 JWT + RBAC（`api-platform/internal/admin/customerservice/`） |
| IM 边界 | 聊天正文、UserSig、WebSocket 由 **pte-live-im** 提供；后台仅订单/系统事件投影 |
| 商户 IM | 店铺 `im-sdk-app`（`admin-merchant/src/views/ecrm/setting/im-sdk-app.vue`）配置商户级 AppId，与平台 IM 隔离 |
| 测试 | `api-platform/internal/admin/customerservice/handler_test.go` — 转接幂等；`handler_integration_test.go` — `TestCustomerServiceHTTPRBACAndQueueClosure` |
| 商户 IM 事件无凭据 | `api-merchant/internal/merchant/imsdk/handler_test.go` — `TestIMEventPayloadContainsNoCredential` |

---

## 5. 跨库投影（outbox / NATS）

| 事实源 | 投影目标 | 机制 |
| --- | --- | --- |
| 商户结算账单/账本 | `qixi_crm_a_merchant_settlement_view` | 商户 `qixi_crm_m_outbox` → NATS → 平台 upsert（`source_settlement_id` 幂等） |
| 商户入驻开通 | `qixi_crm_a_merchant_view` | 入驻 command + 平台审核投影（`onboarding_integration_test.go`） |
| 商品审核 | 业务 `product_view` | 商户审核 outbox 双层重投（运行时 §商品监管） |
| 退款/库存/结算 | 商户库存账本、结算账本 | 业务 outbox → NATS → 商户 command 消费者 |

**生产门禁**：须至少一次「业务服务 + 商户服务 + 平台服务 + NATS + 三库」同进程/同 compose 的端到端演练，且与隔离 HTTP 测试结果一致。当前运行时已验证部分 NATS 闭环（结算、退款 sandbox），**完整 C 端申请 outbox 与全链路支付回调仍 open**。

---

## 6. 门禁状态表（相对 2026-08-04 运行时）

| 门禁项 | 2026-08-04 运行时 | Wave 5 生产门禁 |
| --- | --- | --- |
| 三库导入 / 中文菜单 / Vben 登录 | **closed** | closed |
| 五角色 RBAC（发票、日志、SVIP、用户命令等） | **closed**（隔离 HTTP） | 需生产等价 JWT/菜单复验 |
| 用户资产/券/关系/分组/标签/启停/密码/创建 | **closed**（隔离 + 部分 Vben） | 需生产 DSN 回归 + 截图归档 |
| 商户保证金 / 分账 / 分类 / 分组 / 类型 | **closed** | 需生产 NATS/财务联调 |
| 商户结算 `approved→paid` + 订单应计/冲销 | **closed**（隔离三库 NATS） | **open** — 生产 NATS subject、监控、重试策略 |
| Sandbox 退款 + 库存回补 + 结算负向 | **closed** | **open** — 真实微信/支付宝退款回调 |
| 平台退款监管 / 导出 / 区域范围 | **closed**（HTTP） | **open** — 生产渠道对账 |
| 退货退款全链路 HTTP | **closed** | **open** — 生产物流与渠道 |
| 用户提现审核打款 | **closed**（内部凭证） | **open** — 若接真实打款需单独门禁 |
| 积分商城监管 | **closed** | closed（监管域）；C 端积分下单生产 **open** |
| 预售 / 助力 / 分销读模型 | **closed**（监管） | 商户侧写操作与 C 端闭环 **open** |
| 客服队列 / 转接 / 设置 | **closed**（HTTP）；Vben 截图部分待补 | **open** — IM 多 AppId 生产联调 |
| 订单/控制台区域范围 | **closed**（HTTP + 平台 Vben） | **open** — 商户端订单页逐项 |
| 商品审核 outbox 重投 | **closed**（HTTP） | **open** — 生产 outbox 监控 |
| 商户入驻 NATS | **partial** — 组件验收通过，缺 C 端 outbox 全链路 | **open** |
| 真实微信支付 / 支付宝 | 未测 | **open** |
| Wave 4 布局截图全量 | 未开始 | **open** |

---

## 7. 上线前勾选（维护者）

- [ ] `docs/api/FUNCTIONAL-TRUTH.md` 与实现 diff 为空或已 documented exception
- [ ] §2 所列高风险包 `go test` 在 staging DSN 通过
- [ ] 生产配置无 mock 路由、无测试 PEM、无 echo 密钥
- [ ] NATS outbox 积压/死信告警就绪
- [ ] 客服：平台 + 各店铺 `im-sdk-app` 与 pte-live-im 联调通过
- [ ] [LAYOUT-FIDELITY-CHECKLIST.md](./LAYOUT-FIDELITY-CHECKLIST.md) 高风险页截图 + 字段顺序 closed
- [ ] [ADMIN-PLATFORM-RUNTIME-2026-08-04.md](./ADMIN-PLATFORM-RUNTIME-2026-08-04.md) 中「待补截图」项已归档或仍标注 open
