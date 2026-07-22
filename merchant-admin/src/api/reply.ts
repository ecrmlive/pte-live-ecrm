import { http } from '@/api/http';

export interface ServiceReply {
  service_reply_id: number;
  mer_id: number;
  type: number;
  keyword: string;
  content: string;
  status: number;
  create_time?: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchReplies(params: Record<string, unknown>) {
  return http.get<PageResult<ServiceReply>>('/setting/replies', { params });
}

export function createReply(data: Record<string, unknown>) {
  return http.post<ServiceReply>('/setting/replies', data);
}

export function updateReply(id: number, data: Record<string, unknown>) {
  return http.put<ServiceReply>(`/setting/replies/${id}`, data);
}

export function deleteReply(id: number) {
  return http.delete<{ ok: boolean }>(`/setting/replies/${id}`);
}
