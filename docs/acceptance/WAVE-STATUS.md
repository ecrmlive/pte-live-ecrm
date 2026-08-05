# CRMEB 双后台对照 — Wave 落地进度（2026-08-05）

刷新缺口：`python3 scripts/generate-parity-gap-matrix.py` → [PARITY-GAP-MATRIX.md](./PARITY-GAP-MATRIX.md)

## 优先级（用户确认）

1. **先把管理系统（`admin-platform` + `admin-merchant`）做到按钮级 / 真实 API 闭环**，目标是双后台相对 CRMEB 可宣称可用。
2. **H5 / 小程序（`app-uni` 一套代码）质量不足，后续单独对齐开发**；本阶段不把 C 端页面当作关闭条件。
3. 原生 App / PC 店面同理，不阻塞后台 Wave。

## 路径挂载量化（特征表 page_route × registry）

| 端 | registry 路径 | 特征表 ops 已挂载 | missing 页 |
| --- | ---: | ---: | ---: |
| 平台 | ~288 | **1333 / 1333（100%）** | **0** |
| 商户 | ~115 | **615 / 615（100%）** | **0** |

说明：

- **100% 指「每个 CRMEB page_route 都有 Vben 组件映射」**，不是「每个按钮已按 parity 六条验收关闭」。
- 矩阵中多数仍为 `partial`：需继续加深按钮权限、布局字段对照与截图后，才能标 `done`。
- 部分路径为**别名**或 **stub 配置页**；高风险生产闸门（支付密钥等）未关。

## Wave 状态

| Wave | 状态 |
| --- | --- |
| 0 对照工程化 | 已完成 |
| 1 加深已有页 | 持续（订单/售后/核销权限与 EcrmListPage 加深中） |
| 2 平台补齐 | 路径挂载完成；取消订单页已接 EcrmListPage |
| 3 商户补齐 | 路径挂载完成；代客下单 / 核销 / **售后同意拒绝权限码**已落地 |
| 4 布局还原 | **金标准已锁定为店铺列表**（`LAYOUT-FIDELITY-CHECKLIST.md`）；一级菜单「店铺功能」；全量页迁移与截图未完成 |
| 5 高风险闸门 | 文档已有；生产支付未关 |

## 本轮续作（2026-08-05）

1. **商户代客下单**：真实写库 + `order.proxy`
2. **品牌拆分** + `platform_brand_category`
3. **订单页加深**（平台/商户列表 EcrmListPage）
4. **商户核销原生化** → `qixi_crm_b_order_verification`
5. **支付成功签发 unused 核销码**（pickup/service；C 端接口已备，H5 UI 后置）
6. **SQL 收敛**为 `init_table/config/data/key/test_data`（`init_key.sql` 本地）
7. **商户售后 RBAC**：同意 / 拒绝 / 确认收货分别受 `refund.approve`、`refund.reject` 约束
8. **平台取消订单页**：`cancellation.vue` 接入 `EcrmListPage`，与订单监管布局对齐
9. **平台售后页**：`order/refund.vue` 迁 `EcrmListPage`
10. **商户发票原生化**：`nativeinvoice` 读 `qixi_crm_b_order_invoice`（`store_id` 隔离）+ `invoice.audit`；`invoice.vue` 接入 EcrmListPage 与状态筛选；测试数据含待审夹具 `9900202`
11. **商户提现页**：`finance/withdraw.vue` 迁 `EcrmListPage`
12. **商户结算页**：`finance/settlement.vue` 迁 `EcrmListPage`（权限码已有）
13. **平台转账记录**：去掉 `setting_cache` stub，`GET /finance/transfer-records` 读 `qixi_crm_a_merchant_settlement_view` 打款链路（approved/paid/rejected）
14. **平台财务加深**：`merchant-settlement` / `invoices` / `user-assets` / `statement` 迁 `EcrmListPage`
15. **商户财务加深**：`capital-flow` / `statement` / `transfer` 迁 `EcrmListPage`
16. **优惠套餐原生化**：去掉 `setting_cache` stub；商户 CRUD `qixi_crm_m_marketing_activity`（discount）并同步 `qixi_crm_b_marketing_activity_view`；平台只读监管 + 投影上下架
17. **营销装饰原生化**：氛围/边框/专题/报名 → `qixi_crm_a_marketing_decor`（`nativemarketingdecor`），去掉 setting_cache stub
18. **维护配置原生化**：热搜/组合数据/备份/系统表单 → `qixi_crm_a_config_item`（`nativeconfigitem`）；菜单 ID 冲突已修正（营销 207–211、系统表单 212）

## 下一迭代（管理系统 partial → done，不含 H5）

1. 按 [FIELD-CHECKLIST-TEMPLATE.md](./FIELD-CHECKLIST-TEMPLATE.md) 从订单/财务域逐页对照字段与按钮
2. 继续 stub 配置页接真实业务表（价格说明 / 活动标签 / 应用配置等仍可能走 setting_cache）
3. [HIGH-RISK-PRODUCTION-GATE.md](./HIGH-RISK-PRODUCTION-GATE.md) 关闭支付/退款/结算生产闸门
4. 补 `docs/acceptance/screenshots/{platform|merchant}/...`
5. （后置）`app-uni` H5/小程序单独对齐：`delivery_type` 选择与整体体验重做
