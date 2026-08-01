import qs from 'qs';

import { resolveApiBaseUrl } from '#/utils/pte-live-api';
import { getDecryptedToken } from '#/utils/pte-live-token';

/** 商户后台文件下载（PHP 接口需在 query 带 token） */
export function buildShopDownloadUrl(
  shopPath: string,
  params: Record<string, number | string | undefined>,
) {
  const token = getDecryptedToken();
  if (!token) {
    return '';
  }
  const path = shopPath.startsWith('/') ? shopPath : `/${shopPath}`;
  const query = qs.stringify({ ...params, token });
  const apiBase = resolveApiBaseUrl();
  if (apiBase) {
    return `${apiBase}${path}?${query}`;
  }
  return `/api${path}?${query}`;
}
