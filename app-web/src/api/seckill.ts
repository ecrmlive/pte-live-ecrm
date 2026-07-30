import { http } from "@/utils/request";

export interface SeckillActive {
  seckill_active_id: number;
  name: string;
  product_id: number;
  seckill_price: number;
  price?: number;
  store_name?: string;
  mer_name?: string;
  in_window: boolean;
  start_day: string;
  end_day: string;
}

export function fetchSeckillList(page = 1, limit = 20) {
  return http.get<{ list: SeckillActive[]; total: number }>(
    `/seckill/actives?page=${page}&limit=${limit}`,
    false,
  );
}
