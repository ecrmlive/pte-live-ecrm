import type { Component } from 'vue';

import type { ComponentRecordType } from '@vben/types';

import { NATIVE_PAGE_OVERRIDES } from '#/views/native/registry';

/**
 * 商户后台 pageMap：仅 registry 原生页。
 * init_merchant_access 全量 path 已注册；未命中走 not-found（不再 glob legacy/views）。
 */
export function buildMerchantPageMap(): ComponentRecordType {
  const pageMap: ComponentRecordType = {};

  for (const [key, loader] of Object.entries(NATIVE_PAGE_OVERRIDES)) {
    pageMap[key] = loader as () => Promise<Component>;
  }

  return pageMap;
}
