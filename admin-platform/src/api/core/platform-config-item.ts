import { requestClient } from '#/api/request';

export type ConfigItemType = 'backup' | 'hot_search' | 'system_form';

export interface ConfigItem {
  code: string;
  enabled: boolean;
  id: number;
  item_type: ConfigItemType;
  name: string;
  payload: Record<string, unknown>;
  remark: string;
  sort: number;
  status: number;
  created_at: string;
  updated_at: string;
}

export interface ConfigItemPage {
  limit: number;
  list: ConfigItem[];
  page: number;
  total: number;
}

export interface ConfigItemInput {
  code?: string;
  name: string;
  payload?: Record<string, unknown>;
  remark?: string;
  sort?: number;
  status?: number;
}

const pathByType: Record<ConfigItemType, string> = {
  backup: '/maintain/backups',
  hot_search: '/maintain/hot-search',
  system_form: '/diy/system-forms',
};

export function listConfigItemsApi(
  type: ConfigItemType,
  params: { keyword?: string; limit: number; page: number; status?: number },
) {
  return requestClient.get<ConfigItemPage>(pathByType[type], { params });
}

export function createConfigItemApi(type: ConfigItemType, body: ConfigItemInput) {
  return requestClient.post<ConfigItem>(pathByType[type], body);
}

export function updateConfigItemApi(type: ConfigItemType, id: number, body: ConfigItemInput) {
  return requestClient.put<ConfigItem>(`${pathByType[type]}/${id}`, body);
}

export function setConfigItemStatusApi(type: ConfigItemType, id: number, status: 0 | 1) {
  return requestClient.put<{ id: number; status: number }>(`${pathByType[type]}/${id}/status`, {
    status,
  });
}

export function deleteConfigItemApi(type: ConfigItemType, id: number) {
  return requestClient.delete<{ ok: boolean }>(`${pathByType[type]}/${id}`);
}
