import { http } from '@/api/http';

export interface AppUserRow {
  uid: number;
  account: string;
  nickname: string;
  phone?: string;
  is_svip: number;
  svip_endtime?: string | null;
  is_svip_active: boolean;
  integral?: number;
  now_money?: number;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchUsers(params: Record<string, unknown>) {
  return http.get<PageResult<AppUserRow>>('/users', { params });
}

export function setUserSvip(uid: number, data: { is_svip: number; svip_endtime?: string }) {
  return http.put<AppUserRow>(`/users/${uid}/svip`, data);
}
