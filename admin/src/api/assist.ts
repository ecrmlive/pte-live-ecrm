import { http } from '@/api/http';

export interface ProductAssist {
  product_assist_id: number;
  product_id: number;
  mer_id: number;
  store_name: string;
  assist_price: number;
  assist_count: number;
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

export function fetchAssists(params: Record<string, unknown>) {
  return http.get<PageResult<ProductAssist>>('/assist/actives', { params });
}

export function updateAssist(id: number, data: Record<string, unknown>) {
  return http.put<ProductAssist>(`/assist/actives/${id}`, data);
}

export function deleteAssist(id: number) {
  return http.delete<{ ok: boolean }>(`/assist/actives/${id}`);
}
