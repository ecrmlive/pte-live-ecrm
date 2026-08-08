import type { VxeGridPropTypes } from 'vxe-table';

import { PLATFORM_LIST_PAGER_DEFAULTS } from './platform-list-pager';

/**
 * 列表布局强制标准（100%）：以店铺列表为准。
 * @see admin-platform/src/views/ecrm/merchant/list.vue
 * @see docs/acceptance/LAYOUT-FIDELITY-CHECKLIST.md
 */

/** 列表页操作列：统一右侧固定 */
export function platformListActionColumn(
  overrides: Partial<VxeGridPropTypes.Column> = {},
): VxeGridPropTypes.Column {
  return {
    align: 'center',
    field: 'action',
    fixed: 'right',
    showOverflow: false,
    slots: { default: 'action' },
    title: '操作',
    ...overrides,
  };
}

/** platform-admin 列表页 VxeGrid 分页（与 Vben useVbenVxeGrid 默认一致） */
export function platformListPagerConfig(
  overrides: Partial<VxeGridPropTypes.PagerConfig> = {},
): VxeGridPropTypes.PagerConfig {
  return {
    autoHidden: false,
    enabled: true,
    background: PLATFORM_LIST_PAGER_DEFAULTS.background,
    layouts: [...PLATFORM_LIST_PAGER_DEFAULTS.layouts],
    pageSize: PLATFORM_LIST_PAGER_DEFAULTS.pageSize,
    pageSizes: [...PLATFORM_LIST_PAGER_DEFAULTS.pageSizes],
    size: PLATFORM_LIST_PAGER_DEFAULTS.size,
    ...overrides,
  };
}

export const PLATFORM_LIST_GRID_CLASS = 'platform-vxe-grid';

/**
 * 列表 Grid 通用布局（金标准：店铺列表 merchant/list.vue）。
 * - 不传 height：vxe 的 height:'auto' 会映射为 inline height:100%，表体与 pager 之间留空洞；
 *   列表页由 use-vxe-grid 外层 h-auto + platform-list-page.scss 做内容高度 + 整页滚动。
 * - showOverflow: false：行高随单元格内容增高（多行 slot）。
 * - 备注等自由文本：列级 showOverflow: 'tooltip' + className `col--remark`
 *   （单行 nowrap、合理 minWidth + width/maxWidth 上限、表头居中/内容靠左、ellipsis；禁止 width:'auto'）；
 *   其它需单行省略的列仍可用 showOverflow: true|'tooltip' + `.col--ellipsis`。
 * - 操作列 fixed:right 须不透明底 + 高 z-index（见 platform-list-page.scss），双列窄内容区也必须始终可见可点。
 * - cellConfig.verticalAlign: 'top'：多行 slot 顶对齐。
 *
 * 所有 useVbenVxeGrid 调用会自动合并本常量（见 adapter/vxe-table.ts）。
 * 遵循项目 Vben 列表页约定。
 */
export const PLATFORM_LIST_GRID_LAYOUT = {
  gridClass: PLATFORM_LIST_GRID_CLASS,
  minHeight: 0,
  showOverflow: false,
  cellConfig: {
    verticalAlign: 'top' as const,
  },
  /**
   * 有 fixed 列时需开启横向虚拟滚动阈值（vxe scrollX / virtualXConfig）。
   * gt:0 = 始终启用 x 向滚动布局，固定列由 vxe 叠在可视区右侧；
   * 外层勿 overflow-x:auto（会卷走 fixed），宽度约束见 platform-list-page.scss。
   */
  scrollX: { enabled: true, gt: 0 },
  /**
   * 产品标准：禁止 Vxe 工具栏圆形蓝色「搜索」按钮（magnifying glass）。
   * 表单区「搜索」提交按钮保留；adapter 会强制 search:false，勿再开。
   */
  toolbarConfig: {
    search: false,
  },
};
