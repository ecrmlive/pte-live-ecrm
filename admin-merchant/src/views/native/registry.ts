import type { Component } from 'vue';

/**
 * shop 遗留原生页覆盖表。mergers 走 views/mergers/registry + import.meta.glob。
 */
export const NATIVE_PAGE_OVERRIDES: Record<
  string,
  () => Promise<Component>
> = {};
