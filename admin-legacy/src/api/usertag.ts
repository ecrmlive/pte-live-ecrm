import { http } from '@/api/http';

export interface UserLabel {
  label_id: number;
  label_name: string;
  sort: number;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchUserLabels(params: Record<string, unknown>) {
  return http.get<PageResult<UserLabel>>('/user/labels', { params });
}

export function createUserLabel(data: { label_name: string; sort?: number }) {
  return http.post<UserLabel>('/user/labels', data);
}

export function updateUserLabel(id: number, data: { label_name: string; sort?: number }) {
  return http.put<UserLabel>(`/user/labels/${id}`, data);
}

export function deleteUserLabel(id: number) {
  return http.delete<{ ok: boolean }>(`/user/labels/${id}`);
}

export function fetchUserMarkedLabels(uid: number) {
  return http.get<{ list: UserLabel[] }>(`/user/${uid}/labels`);
}

export function markUserLabels(uid: number, label_ids: number[]) {
  return http.put<{ ok: boolean }>(`/user/${uid}/labels`, { label_ids });
}
