import type { Component } from 'vue';

/**
 * shop 遗留原生页覆盖表。ecrm 走 views/ecrm/registry + import.meta.glob。
 */
export const NATIVE_PAGE_OVERRIDES: Record<
  string,
  () => Promise<Component>
> = {};
