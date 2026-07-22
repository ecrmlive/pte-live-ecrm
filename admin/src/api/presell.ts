import { http } from '@/api/http';

export interface PresellActive {
  product_presell_id: number;
  product_id: number;
  mer_id: number;
  store_name: string;
  price: number;
  stock: number;
  is_show: number;
  status: number;
  mer_name?: string;
  ot_price?: number;
  start_time?: string;
  end_time?: string;
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

export function updatePresell(id: number, data: Record<string, unknown>) {
  return http.put<PresellActive>(`/presell/actives/${id}`, data);
}

export function deletePresell(id: number) {
  return http.delete<{ ok: boolean }>(`/presell/actives/${id}`);
}
