import { http } from '@/api/http';

export interface StoreOrder {
  order_id: number;
  order_sn: string;
  group_order_id: number;
  mer_id: number;
  mer_name?: string;
  uid: number;
  real_name: string;
  user_phone: string;
  user_address: string;
  pay_price: number;
  total_num: number;
  paid: number;
  status: number;
  create_time: string;
  products?: Array<{
    product_id: number;
    product_num: number;
    total_price: number;
    product_info?: string;
  }>;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchOrders(params: { page?: number; limit?: number; paid?: number }) {
  return http.get<PageResult<StoreOrder>>('/orders', { params });
}

export function fetchOrder(id: number) {
  return http.get<StoreOrder>(`/orders/${id}`);
}
