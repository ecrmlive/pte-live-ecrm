import { http } from './http';

export interface StoreOrder {
  order_id: number;
  order_sn: string;
  group_order_id: number;
  uid: number;
  mer_id: number;
  pay_price: number;
  total_num: number;
  paid: number;
  status: number;
  real_name: string;
  user_phone: string;
  user_address: string;
  delivery_name?: string;
  delivery_id?: string;
  create_time: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchOrders(params: Record<string, unknown>) {
  return http.get<PageResult<StoreOrder>>('/orders', { params });
}

export function fetchOrder(id: number) {
  return http.get<StoreOrder>(`/orders/${id}`);
}

export function deliverOrder(
  id: number,
  body: { delivery_name: string; delivery_id: string; delivery_type?: string },
) {
  return http.post(`/orders/${id}/delivery`, body);
}

export function verifyOrder(id: number) {
  return http.post(`/orders/${id}/verify`);
}
