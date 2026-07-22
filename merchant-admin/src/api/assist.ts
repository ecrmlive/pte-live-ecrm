import { http } from '@/api/http';

export interface ProductAssist {
  product_assist_id: number;
  product_id: number;
  store_name: string;
  assist_price: number;
  assist_count: number;
  assist_user_count?: number;
  stock: number;
  is_show: number;
  start_time?: string;
  end_time?: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchAssists(params: Record<string, unknown>) {
  return http.get<PageResult<ProductAssist>>('/assist/actives', { params });
}

export function createAssist(data: Record<string, unknown>) {
  return http.post<ProductAssist>('/assist/actives', data);
}

export function updateAssist(id: number, data: Record<string, unknown>) {
  return http.put<ProductAssist>(`/assist/actives/${id}`, data);
}

export function setAssistShow(id: number, is_show: number) {
  return http.put<ProductAssist>(`/assist/actives/${id}/show`, { is_show });
}

export function deleteAssist(id: number) {
  return http.delete<{ ok: boolean }>(`/assist/actives/${id}`);
}
