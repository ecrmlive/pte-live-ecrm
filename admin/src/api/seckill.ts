import { http } from '@/api/http';

export interface SeckillActive {
  seckill_active_id: number;
  name: string;
  mer_id: number;
  product_id: number;
  seckill_price: number;
  start_day: string;
  end_day: string;
  seckill_time_ids: string;
  status: number;
  store_name?: string;
  mer_name?: string;
  price?: number;
  in_window?: boolean;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchSeckillActives(params: Record<string, unknown>) {
  return http.get<PageResult<SeckillActive>>('/seckill/actives', { params });
}

export function updateSeckillActive(id: number, data: Partial<SeckillActive>) {
  return http.put<SeckillActive>(`/seckill/actives/${id}`, data);
}

export function deleteSeckillActive(id: number) {
  return http.delete<{ ok: boolean }>(`/seckill/actives/${id}`);
}
