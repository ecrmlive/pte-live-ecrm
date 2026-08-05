import { requestClient } from '#/api/request';

export interface PlatformDiscount {
  activity_id: number;
  ends_at: string;
  free_shipping: boolean;
  name: string;
  package_price: number;
  product_ids: number[];
  remark: string;
  starts_at: string;
  status: number;
  store_id: number;
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
  status?: number | string;
  store_id?: number;
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
