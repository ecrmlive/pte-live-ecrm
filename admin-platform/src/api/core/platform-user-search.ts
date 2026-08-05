import { requestClient } from '#/api/request';

export interface UserSearchRecord {
  id: number;
  user_id: number;
  keyword: string;
  source: 'h5' | 'mini' | 'pc';
  created_at: string;
}

export interface UserSearchFilter {
  user_id?: number;
  keyword?: string;
  source?: UserSearchRecord['source'];
  start_date?: string;
  end_date?: string;
}

export function listUserSearchRecords(params: UserSearchFilter & { page: number; limit: number }) {
  return requestClient.get<{ list: UserSearchRecord[]; total: number }>('/user-search-records', { params });
}
export function clearUserSearchRecords(data: { user_id: number; reason: string; idempotency_key: string }) {
  return requestClient.post<{ user_id: number; cleared_count: number }>('/user-search-records/clear', data);
}
export function exportUserSearchRecords(data: UserSearchFilter & { reason: string }) {
  return requestClient.post<{ file_name: string; content: string; row_count: number; truncated: boolean }>('/user-search-records/export', data);
}
