import { requestClient } from '#/api/request';

export interface PlatformPresell {
  action_status?: number;
  delivery_day?: number;
  delivery_type?: number;
  down_price?: number;
  end_time: string;
  final_end_time?: string;
  final_price?: number;
  final_start_time?: string;
  is_show?: number;
  mer_id: number;
  mer_name?: string;
  pay_count?: number;
  presell_type?: number;
  price: number;
  product_id?: number;
  product_presell_id: number;
  product_status?: number;
  seles?: number;
  start_time: string;
  status: number;
  stock?: number;
  store_info?: string;
  store_name: string;
}

export interface PlatformPresellPage {
  limit: number;
  list: PlatformPresell[];
  page: number;
  total: number;
}

export function listPlatformPresellsApi(params: {
  limit: number;
  page: number;
  mer_id?: number;
  keyword?: string;
  status?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<PlatformPresellPage>('/presell/actives', { params });
}

export function getPlatformPresellApi(id: number) {
  return requestClient.get<PlatformPresell>(`/presell/actives/${id}`);
}

export function updatePlatformPresellApi(id: number, payload: { status: number }) {
  return requestClient.put(`/presell/actives/${id}`, payload);
}
