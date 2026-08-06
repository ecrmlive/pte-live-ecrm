import { requestClient } from '#/api/request';

export function clearMaintainCacheApi() {
  return requestClient.post<{ note: string; ok: boolean }>('/maintain/cache/clear');
}

export function getPlatformUserAssetSummaryApi(params?: {
  asset_type?: string;
  user_id?: number;
  keyword?: string;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<{ list: Array<{ asset_type: string; count: number; expense: number; income: number }> }>(
    '/finance/user-assets/summary',
    { params },
  );
}
