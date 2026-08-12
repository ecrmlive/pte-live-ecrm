import { requestClient } from '#/api/request';

export type MaintainCacheScope = 'all' | 'config' | 'replace_domain' | 'store';

export function clearMaintainCacheApi(data: {
  new_domain?: string;
  old_domain?: string;
  scope: MaintainCacheScope;
}) {
  return requestClient.post<{
    deleted_keys: number;
    note: string;
    ok: boolean;
    updated_assets?: number;
  }>('/maintain/cache/actions', data);
}

export function getPlatformUserAssetSummaryApi(params?: {
  asset_type?: string;
  date_from?: string;
  date_to?: string;
  keyword?: string;
  user_id?: number;
}) {
  return requestClient.get<{ list: Array<{ asset_type: string; count: number; expense: number; income: number }> }>(
    '/finance/user-assets/summary',
    { params },
  );
}
