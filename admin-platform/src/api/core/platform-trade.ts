import { requestClient } from '#/api/request';

export interface PlatformOrderProduct {
  activity_id: number;
  product_id: number;
  product_info: string;
  product_num: number;
  product_price: number;
  total_price: number;
}

export interface PlatformOrder {
  create_time: string;
  delivery_id: string;
  delivery_name: string;
  delivery_type: string;
  integral_price: number;
  mer_id: number;
  mer_name?: string;
  order_id: number;
  order_sn: string;
  paid: number;
  pay_price: number;
  pay_time?: string;
  pay_type: number;
  products?: PlatformOrderProduct[];
  status: number;
  total_num: number;
  total_price: number;
  user_address: string;
  user_phone: string;
}

export interface PlatformOrderPage {
  limit: number;
  list: PlatformOrder[];
  page: number;
  total: number;
}

export function listPlatformOrdersApi(params: { limit: number; page: number; paid?: number; status?: number }) {
  return requestClient.get<PlatformOrderPage>('/orders', { params });
}

export function getPlatformOrderApi(id: number) {
  return requestClient.get<PlatformOrder>(`/orders/${id}`);
}
