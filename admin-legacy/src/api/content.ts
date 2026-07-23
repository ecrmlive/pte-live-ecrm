import { http } from '@/api/http';

export interface Notice {
  notice_id: number;
  title: string;
  content: string;
  is_show: number;
  sort: number;
  create_time?: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchNotices(params: Record<string, unknown>) {
  return http.get<PageResult<Notice>>('/notices', { params });
}

export function createNotice(data: { title: string; content: string; sort?: number; is_show?: number }) {
  return http.post<Notice>('/notices', data);
}

export function updateNotice(
  id: number,
  data: { title: string; content: string; sort?: number; is_show?: number },
) {
  return http.put<Notice>(`/notices/${id}`, data);
}

export function deleteNotice(id: number) {
  return http.delete<{ ok: boolean }>(`/notices/${id}`);
}

export interface Agreement {
  key: string;
  label: string;
  content: string;
}

export function fetchAgreements() {
  return http.get<{ list: Agreement[] }>('/agreements');
}

export function fetchAgreement(key: string) {
  return http.get<Agreement>(`/agreements/${key}`);
}

export function saveAgreement(key: string, content: string) {
  return http.put<Agreement>(`/agreements/${key}`, { content });
}
