import { useAccessStore } from '@vben/stores';

import { getDecryptedToken } from '#/utils/qixi-live-token';

import { QIXI_PLATFORM_APP_ID } from '#/utils/qixi-live-api';

/** 与 request.ts 一致：优先 Vben accessStore，避免 platform-user 残留旧 token */
export function getAdminToken() {
  try {
    const accessStore = useAccessStore();
    return accessStore.accessToken || getDecryptedToken() || '';
  } catch {
    return getDecryptedToken() || '';
  }
}

export function getAdminAuthorizationHeader() {
  const token = getAdminToken();
  return token ? `Bearer ${token}` : '';
}

export function applyAdminAuthorization(headers: Record<string, unknown> = {}) {
  const next = { ...headers };
  const authorization = getAdminAuthorizationHeader();
  if (authorization) {
    next['authori-zation'] = authorization;
  }
  // live-api AppContext：平台 admin 路由需租户 app_id（与 PHP /admin 一致为 10000）
  next.AppID = String(QIXI_PLATFORM_APP_ID);
  return next as Record<string, string>;
}
