import { requestClient } from '#/api/request';

export interface PlatformSeckillActive {
  active_status: number;
  end_day: string;
  mer_id: number;
  mer_name?: string;
  name: string;
  product_id: number;
  seckill_active_id: number;
  seckill_price: number;
  start_day: string;
  status: number;
  store_name?: string;
}

export interface PlatformSeckillPage {
  limit: number;
  list: PlatformSeckillActive[];
  page: number;
  total: number;
}

export interface PlatformSeckillInput {
  name?: string;
  seckill_time_ids?: string;
  start_day?: string;
  end_day?: string;
  seckill_price?: number;
  once_pay_count?: number;
  status?: number;
}

export function listPlatformSeckillApi(params: {
  limit: number;
  page: number;
  mer_id?: number;
  keyword?: string;
  status?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<PlatformSeckillPage>('/seckill/actives', { params });
}

export function getPlatformSeckillApi(id: number) {
  return requestClient.get<
    PlatformSeckillActive & { seckill_time_ids: string; once_pay_count: number }
  >('/seckill/actives/' + id);
}

export function updatePlatformSeckillApi(id: number, payload: PlatformSeckillInput) {
  return requestClient.put(`/seckill/actives/${id}`, payload);
}

export function deletePlatformSeckillApi(id: number) {
  return requestClient.delete<{ ok: true }>('/seckill/actives/' + id);
}
