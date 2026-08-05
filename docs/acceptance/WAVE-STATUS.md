# CRMEB 双后台对照 — Wave 落地进度（2026-08-05）

刷新缺口：`python3 scripts/generate-parity-gap-matrix.py` → [PARITY-GAP-MATRIX.md](./PARITY-GAP-MATRIX.md)

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
| 1 加深已有页 | 持续（订单监管/商户订单/用户列表等已加深；全量按钮级未关） |
| 2 平台补齐 | 路径挂载完成；品牌分类/列表已拆分；实现深度参差 |
| 3 商户补齐 | 路径挂载完成；**代客下单 POST /orders/proxy 已真实创建** |
| 4 布局还原 | 规范 + 示范页已有；全量截图未做 |
| 5 高风险闸门 | 文档已有；生产支付未关 |

## 本轮续作（2026-08-05）

1. **商户代客下单**：`api-merchant` 写入 `group_order` + `order` + `order_item` + `stock_command_outbox(reserve)`；权限 `order.proxy`；前端幂等键；UI 去掉 501 提示。
2. **品牌拆分**：`/product/band/brandClassify` → `brand-classify.vue`；`/product/band/brandList` → `brand.vue`；新增 `qixi_crm_a_platform_brand_category` 与 `/brand-categories` API。
3. **订单页加深**：平台/商户订单列表接入 `EcrmListPage`；商户发货按钮按 `order.deliver` 权限码显隐。

## 下一迭代（从 partial → done）

1. 按 [FIELD-CHECKLIST-TEMPLATE.md](./FIELD-CHECKLIST-TEMPLATE.md) 逐页对照 CRMEB 字段与按钮
2. 把其余 stub 配置页接到真实业务表（非仅 setting_cache）
3. 按 [HIGH-RISK-PRODUCTION-GATE.md](./HIGH-RISK-PRODUCTION-GATE.md) 关闭支付/退款/结算生产闸门
4. 补 `docs/acceptance/screenshots/{platform|merchant}/<slug>/{crmeb,qixi}.png`
5. 已有库执行 SQL：`qixi_crm_a_platform_brand_category` + brand.`category_id`；商户菜单 `order.proxy`（`02_data.sql` 已含，需重跑或手工插入）
