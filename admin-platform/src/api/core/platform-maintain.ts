import { requestClient } from '#/api/request';

export interface MaintainCacheListItem {
  id: string;
  name: string;
  enabled: boolean;
  remark: string;
}

export function clearMaintainCacheApi() {
  return requestClient.post<{ note: string; ok: boolean }>('/maintain/cache/clear');
}

export function fetchMaintainBackups() {
  return requestClient.get<{ list: MaintainCacheListItem[] }>('/maintain/backups');
}

export function saveMaintainBackups(list: MaintainCacheListItem[]) {
  return requestClient.put<{ list: MaintainCacheListItem[] }>('/maintain/backups', { list });
}

export function fetchMaintainGroupData() {
  return requestClient.get<{ list: MaintainCacheListItem[] }>('/maintain/group-data');
}

export function saveMaintainGroupData(list: MaintainCacheListItem[]) {
  return requestClient.put<{ list: MaintainCacheListItem[] }>('/maintain/group-data', { list });
}

export function fetchMaintainHotSearch() {
  return requestClient.get<{ list: MaintainCacheListItem[] }>('/maintain/hot-search');
}

export function saveMaintainHotSearch(list: MaintainCacheListItem[]) {
  return requestClient.put<{ list: MaintainCacheListItem[] }>('/maintain/hot-search', { list });
}

export function fetchDiySystemForms() {
  return requestClient.get<{ list: MaintainCacheListItem[] }>('/diy/system-forms');
}

export function saveDiySystemForms(list: MaintainCacheListItem[]) {
  return requestClient.put<{ list: MaintainCacheListItem[] }>('/diy/system-forms', { list });
}

export function fetchTransferRecords() {
  return requestClient.get<{ list: MaintainCacheListItem[] }>('/finance/transfer-records');
}

export function saveTransferRecords(list: MaintainCacheListItem[]) {
  return requestClient.put<{ list: MaintainCacheListItem[] }>('/finance/transfer-records', { list });
}

export function getPlatformUserAssetSummaryApi() {
  return requestClient.get<{ list: Array<{ asset_type: string; count: number; expense: number; income: number }> }>(
    '/finance/user-assets/summary',
  );
}
