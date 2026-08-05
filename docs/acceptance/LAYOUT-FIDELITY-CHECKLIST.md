# 管理后台布局强制标准（100%）

> **本文件是 platform / merchant 管理后台列表与抽屉页的强制布局规范，不是可选建议。**  
> **金标准实现**：`admin-platform/src/views/ecrm/merchant/list.vue`（店铺列表）。  
> 目标：统一后台布局分区与交互对齐 [CRMEB Merchant](https://mer.crmeb.net/admin/dashboard)（控制台 / 列表 / 添加编辑详情 / 确认弹窗）。  
> **颜色不纳入本清单**：主题色由 Vben 主题配置决定，不要求复刻 CRMEB 珊瑚红。  
> **组件强制**：必须使用 Vben Common UI / 适配层组件，禁止以自研 `Ecrm*` Element 壳作为列表/表单主路径。

## 0. 金标准（必须遵守）

| 项 | 强制要求 |
| --- | --- |
| 参考页 | **仅以** `admin-platform/src/views/ecrm/merchant/list.vue` 为 100% 布局样板 |
| 适用范围 | `admin-platform`、`admin-merchant` 所有列表页、筛选、工具栏、分页、添加/编辑/详情 Drawer |
| 共享工具 | `platformListActionColumn` / `platformListPagerConfig` / `PLATFORM_LIST_GRID_CLASS`（`#/constants/platform-list-grid`）及对应 SCSS |
| 禁止 | 新建 `EcrmListPage` + 手写 `el-table` / `el-pagination` 作为列表骨架；Modal 代替本应侧滑的宽表单/详情（除非 CRMEB 本身是居中弹窗） |
| Agent / 开发 | 改任何列表页前先打开店铺列表对照；偏离金标准视为未完成 |

店铺列表已落地、且后续页必须复刻的要点：

1. `Page` + `auto-content-height`，默认不传页面 `title` / `description`（由菜单/面包屑承担）。
2. `useVbenVxeGrid`：`formOptions` 筛选 → 状态 Tab / 工具栏 → 表格 → 底部分页。
3. 操作列 `fixed: 'right'`（`platformListActionColumn`）；单元格可用 `ElButton` link / `ElSwitch` / `ElTag`。
4. 添加 / 编辑 / 详情使用 `useVbenDrawer`（宽抽屉 ~1000px）；详情只读描述 + 可切编辑；确认用 `confirm()`。
5. 表格行高自适应（`showOverflow: false` + platform grid 布局），整页滚动而非表体内滚。
6. 无「刷新列表」类无意义按钮；主操作（如添加）放在工具栏。

## 1. 强制组件栈（Vben）

对照官方文档与本地 skill（`~/Downloads/vben-skills`）：

| 场景 | 必须使用 | 文档 |
| --- | --- | --- |
| 页面壳 | `Page`（`auto-content-height`） | `@vben/common-ui` |
| 列表 + 筛选 + 分页 | `useVbenVxeGrid`（`#/adapter/vxe-table`）+ `formOptions` | [vben-vxe-table](https://doc.vben.pro/components/common-ui/vben-vxe-table.html) |
| 表单字段 | `useVbenForm` / Grid 内 `formOptions.schema` | [vben-form](https://doc.vben.pro/components/common-ui/vben-form.html) |
| 居中确认 / 轻量弹窗 | `useVbenModal` | [vben-modal](https://doc.vben.pro/components/common-ui/vben-modal.html) |
| 详情 / 编辑侧滑 | `useVbenDrawer` | [vben-drawer](https://doc.vben.pro/components/common-ui/vben-drawer.html) |
| 只读详情 | 双列描述布局（对齐店铺详情）或 `VbenDescriptions` | [vben-descriptions](https://doc.vben.pro/components/common-ui/vben-descriptions.html) |
| 确认框 | `confirm` / `alert`（`@vben/common-ui`） | skill `alert.md` |
| 裁剪 / 富文本 | `Cropper` / `Tiptap`（按页需要） | [cropper](https://doc.vben.pro/components/common-ui/vben-cropper.html) / [tiptap](https://doc.vben.pro/components/common-ui/vben-tiptap.html) |

Element Plus 仅允许作为 **单元格内控件**（如 `ElButton` link、`ElSwitch`、`ElTag`），不得再包一层 `EcrmListPage` + 手写 `el-table` / `el-pagination` / `el-form` 作为列表页骨架。

### 1.1 列表页结构（对齐店铺列表 / CRMEB 分区）

```vue
<Page auto-content-height>
  <Grid><!-- useVbenVxeGrid：formOptions 即筛选；proxy + pagerConfig 即分页 -->
    <template #toolbar-actions>
      <!-- 状态 Tab（可选）→ 主按钮（添加） -->
    </template>
    <template #action="{ row }"><!-- fixed: 'right' 操作列 --></template>
  </Grid>
  <ShopDrawer /><!-- useVbenDrawer：添加 / 编辑 / 详情 -->
</Page>
```

- 外层 `Page` 用 `auto-content-height`，**默认不传** `title` / `description`。
- 分区顺序固定：**筛选（Grid form）→ Tab / 工具栏 → 表格 → 分页**。
- 操作列：`fixed: 'right'`。
- 确认类交互用 `confirm({ title: '提示', icon: 'warning', ... })`。
- 权限：`v-if` / 按钮码不得因布局迁移删除或放宽。

### 1.2 参考实现

- **金标准（强制）**：`admin-platform/src/views/ecrm/merchant/list.vue`
- 共享列/分页：`admin-platform/src/constants/platform-list-grid.ts`
- 既有 Vxe 范例（迁移中）：`admin-platform/src/views/role/index.vue`、`admin-platform/src/views/ecrm/user/label.vue`

### 1.3 店铺侧栏

种子：`sql/admin/init_data.sql`。本地需重跑 `make local-db-init`（或更新菜单行）后重新登录。

```text
店铺功能 (lucide:store)
├── 店铺管理 (lucide:layout-grid)
│   ├── 店铺列表 / 店铺分类 / 店铺分组 / 店铺类型
│   ├── 店铺入驻申请 / 店铺分账申请
└── 店铺设置 (lucide:settings)
    ├── 店铺保证金 / 保证金配置 / 店铺菜单 / 说明提示
商户管理 (lucide:users)
├── 商户列表 / 商户入驻审核 / 商户管理员 / 商户设置
```

图标统一 `lucide:*`，且必须落在 `admin-platform` 离线白名单（`platform-lucide-icons.ts` / `build-iconify-offline.mjs`）。

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

## 4. 完成状态

| 项 | 状态 |
| --- | --- |
| 金标准锁定为店铺列表 | **已完成** |
| Vben 组件栈约定落地（文档 + 店铺列表） | **已完成** |
| 全量页面迁移至 `useVbenVxeGrid` / Drawer | **进行中** |
| 逐页 FIELD-CHECKLIST 填完 | **进行中** |
| 逐页 `crmeb.png` + `qixi.png` | **未开始**（目录已预留） |
| 删除 `Ecrm*` 包装组件 | **待迁移完成后** |

**判定**：仅在「截图齐全 **且** 字段/列/筛选项顺序与 CRMEB 一致 **且** 布局 100% 对齐店铺列表金标准」后，该页方可标为 `aligned`；否则保持 `partial`。
