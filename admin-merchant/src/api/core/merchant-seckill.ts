import { requestClient } from '#/api/request';

export interface MerchantSeckillTime {
  end_time: number;
  seckill_time_id: number;
  start_time: number;
  status: number;
  title: string;
}

export interface MerchantSeckillActive {
  active_status: number;
  end_day: string;
  image?: string;
  in_window: boolean;
  name: string;
  once_pay_count: number;
  price?: number;
  product_id: number;
  seckill_active_id: number;
  seckill_price: number;
  seckill_time_ids: string;
  start_day: string;
  status: number;
  store_name?: string;
}

export interface MerchantSeckillPage {
  limit: number;
  list: MerchantSeckillActive[];
  page: number;
  total: number;
}

export interface MerchantSeckillSaveInput {
  end_day: string;
  name: string;
  once_pay_count: number;
  product_id: number;
  seckill_price: number;
  seckill_time_ids: string;
  start_day: string;
  status?: number;
}

export function listMerchantSeckillTimesApi() {
  return requestClient.get<{ list: MerchantSeckillTime[] }>('/seckill/times');
}

export function listMerchantSeckillActivesApi(params: { limit: number; page: number }) {
  return requestClient.get<MerchantSeckillPage>('/seckill/actives', { params });
}

export function createMerchantSeckillActiveApi(body: MerchantSeckillSaveInput) {
  return requestClient.post<MerchantSeckillActive>('/seckill/actives', body);
}

export function updateMerchantSeckillActiveApi(id: number, body: MerchantSeckillSaveInput) {
  return requestClient.put<MerchantSeckillActive>(`/seckill/actives/${id}`, body);
}

export function setMerchantSeckillStatusApi(id: number, status: number) {
  return requestClient.put<MerchantSeckillActive>(`/seckill/actives/${id}/status`, { status });
}

export function deleteMerchantSeckillActiveApi(id: number) {
  return requestClient.delete(`/seckill/actives/${id}`);
}
