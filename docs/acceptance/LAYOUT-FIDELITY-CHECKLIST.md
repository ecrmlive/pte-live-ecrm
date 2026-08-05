# Wave 4 — CRMEB 后台布局保真验收

> 目标：统一后台（platform / merchant）列表、表单、详情页的布局分区与字段顺序对齐 CRMEB Merchant v4.0，而非仅实现同等 API。

## 1. 共享布局组件

两端均从 `admin-*/src/components/ecrm` 引入，禁止在各页面重复手写顶筛 + 卡片 + 分页结构。

| 组件 | 路径 | 用途 |
| --- | --- | --- |
| `EcrmListPage` | `EcrmListPage.vue` | 列表页骨架：`#filters` → `#actions`（可选）→ 默认 slot 表格 → `#pager` |
| `EcrmFormDrawer` | `EcrmFormDrawer.vue` | 侧滑表单/详情；`#footer` 放提交/取消 |
| `EcrmFormDialog` | `EcrmFormDialog.vue` | 弹窗表单；`#footer` 放确认/取消 |
| `EcrmDetailSections` | `EcrmDetailSections.vue` | 详情分块：`sections[]` 驱动 `el-descriptions` |

导出：

```ts
import { EcrmDetailSections, EcrmFormDialog, EcrmFormDrawer, EcrmListPage } from '#/components/ecrm';
```

### 1.1 `EcrmListPage` 插槽约定

```vue
<Page auto-content-height>
  <EcrmListPage title="页面标题" description="可选说明">
    <template #filters><!-- 筛选项，顺序对齐 CRMEB --></template>
    <template #actions><!-- 导出、批量操作等 --></template>
    <!-- 默认 slot：el-table -->
    <template #pager><!-- el-pagination --></template>
  </EcrmListPage>
</Page>
```

- 外层保留 Vben `Page`（`auto-content-height`），标题与说明交给 `EcrmListPage`，避免重复。
- 筛选、工具栏按钮、表格列、分页器的**自上而下 / 自左而右**顺序必须与 CRMEB 同路径页面一致。
- 权限：`v-if` / 按钮码逻辑不得因布局迁移而删除或放宽。

### 1.2 表单与详情

- 新建/编辑：优先 `EcrmFormDrawer`（CRMEB 侧滑）或 `EcrmFormDialog`（CRMEB 弹窗）；字段顺序见 per-page 清单。
- 只读详情：大块用 `EcrmDetailSections`；嵌套子表（商品行等）仍可用 `el-table`。

## 2. 单页字段对照清单

每一页迁移或验收前，复制 [FIELD-CHECKLIST-TEMPLATE.md](./FIELD-CHECKLIST-TEMPLATE.md) 并填写：

- 元信息：端、CRMEB 路径、七禧路径、TSV 操作数、截图目录
- 筛选项、工具栏、表格列、表单字段、详情分块 — 各表逐项勾选「已对齐」
- 完成条件勾选（真实菜单、真实 API、RBAC、高风险域、截图）

## 3. 截图证据目录

同一路由、同一筛选/状态，分别截取 CRMEB 与七禧：

```text
docs/acceptance/screenshots/{platform|merchant}/<slug>/crmeb.png
docs/acceptance/screenshots/{platform|merchant}/<slug>/qixi.png
```

- `<slug>`：由路由 path 推导，kebab-case，例如 `/accounts/withdraw` → `accounts-withdraw`，`/order/refund` → `order-refund`。
- 两张图须为**同一业务状态**（相同筛选、相同行数据形态）；不得用空列表掩盖字段差异。
- 说明见 [screenshots/README.md](./screenshots/README.md)。

## 4. Wave 4 完成状态

| 项 | 状态 |
| --- | --- |
| 共享组件落地（platform + merchant） | 已完成 |
| 全量页面迁移至 `EcrmListPage` 等 | **进行中** |
| 逐页 FIELD-CHECKLIST 填完 | **进行中** |
| 逐页 `crmeb.png` + `qixi.png` | **未开始**（目录已预留） |

**Wave 4 总状态：`partial`（部分）** — 仅在「截图齐全 **且** 字段/列/筛选项顺序与 CRMEB 一致」后，该页方可标为 `aligned`；否则保持 `partial`。

## 5. 示范页面（本波已 refactor）

| 端 | 页面 | 组件 |
| --- | --- | --- |
| platform | `admin-platform/src/views/ecrm/accounts/withdraw.vue` | `EcrmListPage` |
| merchant | `admin-merchant/src/views/ecrm/order/refund.vue` | `EcrmListPage` |

其余列表页按同一模式迁移；详情抽屉/拒绝对话框可在后续波次改为 `EcrmFormDrawer` / `EcrmFormDialog`。

## 6. 验收步骤（单页）

1. 从 CRMEB 参考站或 `~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0` 确认字段顺序。
2. 复制 FIELD-CHECKLIST 模板并填表。
3. 页面改用 `EcrmListPage`（及表单/详情组件）。
4. 真实菜单进入，验证 RBAC 与 API（无 mock 假成功）。
5. 截取 `crmeb.png` / `qixi.png` 写入对应 slug 目录。
6. 清单全部勾选后，在 [CRMEB-PORTAL-PARITY-LEDGER.md](./CRMEB-PORTAL-PARITY-LEDGER.md) 或角色台账中更新该页布局状态。
