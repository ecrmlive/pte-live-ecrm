import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface FullreduceItem {
  active_name: string;
  create_time?: string;
  full_type: number;
  full_type_text?: string;
  full_value: number | string;
  fullreduce_id: number;
  reduce_type: number;
  reduce_type_text?: string;
  reduce_value: number | string;
}

export interface FullreduceFormPayload {
  active_name: string;
  full_type: number;
  full_value: number | string;
  fullreduce_id?: number;
  reduce_type: number;
  reduce_value: number | string;
}

export async function getFullreduceListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<FullreduceItem> }>(
    '/shop/plus.fullreduce/index',
    params,
  );
}

export async function addFullreduceApi(payload: FullreduceFormPayload) {
  return requestClient.post('/shop/plus.fullreduce/add', payload);
}

export async function editFullreduceApi(payload: FullreduceFormPayload) {
  return requestClient.post('/shop/plus.fullreduce/edit', payload);
}

export async function deleteFullreduceApi(fullreduceId: number) {
  return requestClient.post('/shop/plus.fullreduce/delete', { fullreduce_id: fullreduceId });
}

export interface FullreduceReduceRule {
  full_type?: number;
  full_value?: number | string;
  reduce_type?: number;
  reduce_value?: number | string;
}

export interface FullreduceProductItem {
  image?: Array<{ file_path?: string }>;
  product_id: number;
  product_name: string;
  product_price?: number | string;
  product_stock?: number;
  reduce_list?: FullreduceReduceRule[];
  reduce_pid?: number | null;
  sales_actual?: number;
}

export async function getFullreduceProductListApi(params: Record<string, unknown>) {
  return requestClient.post<{
    category: Array<{ category_id: number; name: string }>;
    list: PaginatedList<FullreduceProductItem>;
    product_count?: Record<string, number>;
  }>('/shop/plus.fullreduce/product', params);
}

export async function saveFullreduceProductApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.fullreduce/editProduct', payload);
}

export async function batchFullreduceProductApi(productIds: number[], isJoin = 0) {
  return requestClient.post('/shop/plus.fullreduce/batchProduct', {
    is_join: isJoin,
    productIds: productIds.join(','),
  });
}
