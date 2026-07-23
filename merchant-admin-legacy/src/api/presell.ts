import { http } from '@/api/http';

export interface PresellActive {
  product_presell_id: number;
  product_id: number;
  store_name: string;
  price: number;
  down_price?: number;
  final_price?: number;
  presell_type?: number;
  stock: number;
  is_show: number;
  pay_count?: number;
  delivery_day?: number;
  start_time?: string;
  end_time?: string;
  final_start_time?: string;
  final_end_time?: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchPresells(params: Record<string, unknown>) {
  return http.get<PageResult<PresellActive>>('/presell/actives', { params });
}

export function createPresell(data: Record<string, unknown>) {
  return http.post<PresellActive>('/presell/actives', data);
}

export function updatePresell(id: number, data: Record<string, unknown>) {
  return http.put<PresellActive>(`/presell/actives/${id}`, data);
}

export function setPresellShow(id: number, is_show: number) {
  return http.put<PresellActive>(`/presell/actives/${id}/show`, { is_show });
}

export function deletePresell(id: number) {
  return http.delete<{ ok: boolean }>(`/presell/actives/${id}`);
}
