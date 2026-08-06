import type { VxeGridPropTypes } from 'vxe-table';

import { MERCHANT_LIST_PAGER_DEFAULTS } from './merchant-list-pager';

/** 列表页操作列：统一右侧固定 */
export function merchantListActionColumn(
  overrides: Partial<VxeGridPropTypes.Column> = {},
): VxeGridPropTypes.Column {
  return {
    field: 'action',
    fixed: 'right',
    showOverflow: false,
    slots: { default: 'action' },
    title: '操作',
    ...overrides,
  };
}

/** merchant-admin 列表页 VxeGrid 分页（与 Vben / platform ListPagination 一致） */
export function merchantListPagerConfig(
  overrides: Partial<VxeGridPropTypes.PagerConfig> = {},
): VxeGridPropTypes.PagerConfig {
  return {
    autoHidden: false,
    enabled: true,
    background: MERCHANT_LIST_PAGER_DEFAULTS.background,
    layouts: [...MERCHANT_LIST_PAGER_DEFAULTS.layouts],
    pageSize: MERCHANT_LIST_PAGER_DEFAULTS.pageSize,
    pageSizes: [...MERCHANT_LIST_PAGER_DEFAULTS.pageSizes],
    size: MERCHANT_LIST_PAGER_DEFAULTS.size,
    ...overrides,
  };
}

export const MERCHANT_LIST_GRID_CLASS = 'merchant-vxe-grid';

/**
 * 列表 Grid 通用布局。
 * - height: 'auto'：配置保留；实际表体行为由 merchant-list-page.scss 覆盖（Page 整页滚动，非 tbody 内滚）。
 * - showOverflow: false：允许行高随单元格内容增高（订单等多行 slot 必需）。
 * - 自由文本列（备注等）须列级 showOverflow: true，单行省略，避免撑高行挤掉 fixed 操作列。
 * - cellConfig.verticalAlign: 'top'：多行内容与操作列顶对齐。
 *
 * 所有 useVbenVxeGrid 调用会自动合并本常量（见 adapter/vxe-table.ts）。
 * 遵循项目 Vben 列表页约定。
 */
export const MERCHANT_LIST_GRID_LAYOUT = {
  gridClass: MERCHANT_LIST_GRID_CLASS,
  height: 'auto' as const,
  minHeight: 0,
  showOverflow: false,
  cellConfig: {
    verticalAlign: 'top' as const,
  },
  /** 有 fixed 列时需开启横向滚动容器 */
  scrollX: { enabled: true, gt: 0 },
};
