import { http } from '@/api/http';

export interface Express {
  express_id: number;
  name: string;
  code: string;
  sort: number;
  is_show: number;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchExpress(params: Record<string, unknown>) {
  return http.get<PageResult<Express>>('/express', { params });
}

export function createExpress(data: { name: string; code?: string; sort?: number; is_show?: number }) {
  return http.post<Express>('/express', data);
}

export function updateExpress(
  id: number,
  data: { name: string; code?: string; sort?: number; is_show?: number },
) {
  return http.put<Express>(`/express/${id}`, data);
}

export function deleteExpress(id: number) {
  return http.delete<{ ok: boolean }>(`/express/${id}`);
}
