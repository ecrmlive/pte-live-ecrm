import { requestClient } from '#/api/request';

export interface PlatformPointsProduct {
  cate_id?: number;
  cover_url?: string;
  create_time?: string;
  merchant_id: number;
  merchant_name: string;
  original_price: number;
  points_required: number;
  product_id: number;
  sale_status: number;
  sales?: number;
  sort?: number;
  source_product_id?: number;
  stock: number;
  store_id: number;
  store_name: string;
  title: string;
  updated_at?: string;
  version: number;
}

export interface PlatformPointsOrder {
  created_at: string;
  id: number;
  order_no: string;
  pay_status: string;
  points_amount: number;
  total_quantity: number;
  user_id: number;
}

export interface PlatformPointsExchange {
  created_at: string;
  order_id: number;
  order_no: string;
  pay_status: string;
  points_amount: number;
  quantity: number;
  title_snapshot: string;
  user_id: number;
}

export type PlatformPointsProductSaveInput = {
  cate_id?: number;
  cover_url?: string;
  merchant_id?: number;
  merchant_name?: string;
  original_price?: number;
  points_required?: number;
  sale_status?: number;
  sales?: number;
  sort?: number;
  source_product_id?: number;
  stock?: number;
  store_id?: number;
  store_name?: string;
  title?: string;
  version?: number;
};

export function listPlatformPointsProductsApi(params: {
  keyword?: string;
  limit: number;
  merchant_id?: number;
  page: number;
  sale_status?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<{
    list: PlatformPointsProduct[];
    total: number;
    page: number;
    limit: number;
  }>('/points/products', { params });
}

export function getPlatformPointsSummaryApi() {
  return requestClient.get<{ total: number; on_sale: number; stock: number }>(
    '/points/products/summary',
  );
}

export function getPlatformPointsProductApi(id: number) {
  return requestClient.get<PlatformPointsProduct>(`/points/products/${id}`);
}

export function createPlatformPointsProductApi(
  payload: PlatformPointsProductSaveInput,
) {
  return requestClient.post<PlatformPointsProduct>('/points/products', payload);
}

export function quickAddPlatformPointsProductApi(payload: {
  cate_id?: number;
  cover_url?: string;
  original_price?: number;
  points_required: number;
  sale_status?: number;
  sort?: number;
  source_product_id: number;
  stock: number;
  title?: string;
}) {
  return requestClient.post<PlatformPointsProduct>(
    '/points/products/quick',
    payload,
  );
}

export function updatePlatformPointsProductApi(
  id: number,
  payload: PlatformPointsProductSaveInput,
) {
  return requestClient.put<PlatformPointsProduct>(
    `/points/products/${id}`,
    payload,
  );
}

export function updatePlatformPointsProductStatusApi(
  id: number,
  status: 0 | 1,
) {
  return requestClient.put<PlatformPointsProduct>(
    `/points/products/${id}/status`,
    { status },
  );
}

export function copyPlatformPointsProductApi(id: number) {
  return requestClient.post<PlatformPointsProduct>(
    `/points/products/${id}/copy`,
  );
}

export function deletePlatformPointsProductApi(id: number) {
  return requestClient.delete<{ product_id: number }>(`/points/products/${id}`);
}

export function listPlatformPointsExchangesApi(
  id: number,
  params: { limit: number; page: number },
) {
  return requestClient.get<{
    list: PlatformPointsExchange[];
    total: number;
    page: number;
    limit: number;
  }>(`/points/products/${id}/exchanges`, { params });
}

export function listPlatformPointsOrdersApi(params: {
  limit: number;
  page: number;
  pay_status?: string;
  product_id?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<{
    list: PlatformPointsOrder[];
    total: number;
    page: number;
    limit: number;
  }>('/points/orders', { params });
}
