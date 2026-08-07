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
 * 列表 Grid 通用布局。
 * - height: 'auto'：配置保留；实际表体行为由 platform-list-page.scss 覆盖（Page 整页滚动，非 tbody 内滚）。
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
  height: 'auto' as const,
  minHeight: 0,
  showOverflow: false,
  cellConfig: {
    verticalAlign: 'top' as const,
  },
  /** 有 fixed 列时开启横向滚动；固定列由 vxe 叠在可视区右侧，勿把 overflow-x 加在会卷走 fixed 的外层 */
  scrollX: { enabled: true, gt: 0 },
};
