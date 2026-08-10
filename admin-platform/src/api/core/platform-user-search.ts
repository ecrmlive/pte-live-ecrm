import { requestClient } from '#/api/request';

export type UserSearchChannel =
  | ''
  | 'wechat'
  | 'mini_program'
  | 'h5'
  | 'app'
  | 'pc';

export interface UserSearchRecord {
  id: number;
  user_id: number;
  nickname: string;
  avatar_url: string;
  user_type: string;
  keyword: string;
  created_at: string;
}

export interface UserSearchFilter {
  keyword?: string;
  nickname?: string;
  user_type?: Exclude<UserSearchChannel, ''>;
  start_date?: string;
  end_date?: string;
}

export function listUserSearchRecords(
  params: UserSearchFilter & { page: number; limit: number },
) {
  return requestClient.get<{ list: UserSearchRecord[]; total: number }>(
    '/user-search-records',
    { params },
  );
}

export function clearUserSearchRecords(data?: {
  reason?: string;
  idempotency_key?: string;
}) {
  return requestClient.post<{ cleared_count: number }>(
    '/user-search-records/clear',
    {
      reason: data?.reason || '一键清空搜索记录',
      idempotency_key:
        data?.idempotency_key || `search-clear-all-${crypto.randomUUID()}`,
    },
  );
}

export function exportUserSearchRecords(
  data: UserSearchFilter & { reason?: string },
) {
  return requestClient.post<{
    file_name: string;
    content: string;
    row_count: number;
    truncated: boolean;
  }>('/user-search-records/export', {
    ...data,
    reason: data.reason || '平台后台导出搜索记录',
  });
}
