import { requestClient } from '#/api/request';

export interface UserSearchRecord {
  id: number;
  product_id: number;
  product_title: string;
  user_id: number;
  viewed_at: string;
}

export interface UserSearchRecordPage {
  limit: number;
  list: UserSearchRecord[];
  page: number;
  total: number;
}

export function listUserSearchRecordsApi(params: { limit: number; page: number }) {
  return requestClient.get<UserSearchRecordPage>('/user/search-records', { params });
}
