/** Select 下拉 teleport，避免弹层被裁剪 */
export const PLATFORM_SEARCH_SELECT_PROPS = {
  popperClass: 'platform-search-select-popper',
  popperOptions: { strategy: 'fixed' as const },
  teleported: true,
};
