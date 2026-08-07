/** 与 @vben/plugins/vxe-table use-vxe-grid 默认 pager 对齐（platform-admin 同源） */
export const MERCHANT_LIST_PAGE_SIZES = [10, 15, 20, 30, 50, 100] as const;

export const MERCHANT_LIST_PAGER_LAYOUTS = [
  'Total',
  'Sizes',
  'Home',
  'PrevJump',
  'PrevPage',
  'Number',
  'NextPage',
  'NextJump',
  'End',
] as const;

export const MERCHANT_LIST_PAGER_DEFAULTS = {
  background: true,
  layouts: [...MERCHANT_LIST_PAGER_LAYOUTS],
  /** 与 platform-admin 列表金标准一致：默认每页 10 条 */
  pageSize: 10,
  pageSizes: [...MERCHANT_LIST_PAGE_SIZES],
  size: 'mini' as const,
};
