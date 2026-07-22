import { http } from './http';

export interface StoreOrder {
  order_id: number;
  order_sn: string;
  mer_id: number;
  uid: number;
  pay_price: number;
  total_num: number;
  paid: number;
  status: number;
  real_name: string;
  user_phone: string;
  user_address: string;
  verify_code?: string;
  create_time?: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchAwaitVerify(page = 1, limit = 20) {
  return http.get<PageResult<StoreOrder>>('/orders', {
    params: { page, limit, await_verify: 1 },
  });
}

export function fetchAwaitShip(page = 1, limit = 20) {
  return http.get<PageResult<StoreOrder>>('/orders', {
    params: { page, limit, paid: 1, status: 0 },
  });
}

export function deliverOrder(
  id: number,
  body: { delivery_name: string; delivery_id: string; delivery_type?: string },
) {
  return http.post(`/orders/${id}/delivery`, body);
}

export function fetchOrder(id: number) {
  return http.get<StoreOrder>(`/orders/${id}`);
}

export function fetchOrderByCode(code: string) {
  return http.get<StoreOrder>(`/orders/code/${encodeURIComponent(code)}`);
}

export function verifyOrder(id: number, verify_code?: string) {
  return http.post(`/orders/${id}/verify`, { verify_code: verify_code || '' });
}
