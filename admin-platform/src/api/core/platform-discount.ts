import { requestClient } from '#/api/request';

export interface PlatformDiscountProduct {
  image?: string;
  product_id: number;
  spec?: string;
  store_name?: string;
  type: number;
}

export interface PlatformDiscount {
  activity_id: number;
  combo_products?: PlatformDiscountProduct[];
  create_time?: string;
  created_at?: string;
  ends_at: string;
  free_shipping: boolean;
  id?: number;
  is_limit: number;
  is_time: number;
  limit_num: number;
  main_products?: PlatformDiscountProduct[];
  name: string;
  package_price: number;
  package_type: number;
  package_type_label?: string;
  product_ids: number[];
  products?: PlatformDiscountProduct[];
  qty_label?: string;
  remain_label?: string;
  remain_num?: number;
  remark: string;
  starts_at: string;
  status: number;
  status_label?: string;
  store_id: number;
  store_name?: string;
  time_label?: string;
  version: number;
}

export interface PlatformDiscountPage {
  limit: number;
  list: PlatformDiscount[];
  page: number;
  total: number;
}

export function listPlatformDiscountsApi(params: {
  keyword?: string;
  limit: number;
  page: number;
  package_type?: number | string;
  status?: number | string;
  store_id?: number;
  type?: number | string;
}) {
  return requestClient.get<PlatformDiscountPage>('/marketing/discounts', { params });
}

export function getPlatformDiscountApi(id: number) {
  return requestClient.get<PlatformDiscount>(`/marketing/discounts/${id}`);
}

export function setPlatformDiscountStatusApi(id: number, status: 0 | 1) {
  return requestClient.put<{ activity_id: number; status: number }>(
    `/marketing/discounts/${id}/status`,
    { status },
  );
}
