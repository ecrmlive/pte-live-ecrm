import { requestClient } from '#/api/request';

export interface StoreCustomer {
  last_order_at: string;
  mobile: string;
  nickname: string;
  order_count: number;
  status: number;
  total_pay: number;
  user_id: number;
}

export interface StoreCustomerPage {
  limit: number;
  list: StoreCustomer[];
  page: number;
  total: number;
}

export function listStoreCustomersApi(params: { keyword?: string; limit: number; page: number }) {
  return requestClient.get<StoreCustomerPage>('/store-customers', { params });
}
