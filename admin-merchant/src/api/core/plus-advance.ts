import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface AdvanceProductItem {
  advance_product_id: number;
  create_time?: number | string;
  limit_num: number;
  money: number | string;
  product: {
    image?: Array<{ file_path: string }>;
    product_name: string;
    spec_type: number;
  };
  sku: Array<{
    advance_price: number | string;
    product_price: number | string;
  }>;
  sort: number;
  status: number;
  stock: number;
}

export interface AdvanceSettingValues {
  end_time: number | string;
  image: Array<{ file_id?: number; file_path: string }>;
  is_agent: boolean | number;
  is_coupon: boolean | number;
  is_point: boolean | number;
  is_user_grade: boolean | number;
  money_return: boolean | number;
  pay_time: number | string;
}

export async function getAdvanceProductListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{
    exclude_ids: number[];
    list: PaginatedList<AdvanceProductItem>;
  }>('/shop/plus.advance.Product/index', params);
}

export async function deleteAdvanceProductApi(id: number) {
  return requestClient.post('/shop/plus.advance.Product/delete', { id });
}

export async function getAdvanceSettingApi() {
  return requestClient.get<{ vars: { values: AdvanceSettingValues } }>(
    '/shop/plus.advance.Setting/index',
  );
}

export async function saveAdvanceSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.advance.Setting/index', payload);
}

export interface AdvanceProductAddMetaResult {
  model: Record<string, unknown>;
  specList: Array<{ spec_name: string; spec_sku_id: number | string }>;
}

export interface AdvanceProductEditMetaResult {
  model: Record<string, unknown>;
}

export async function getAdvanceProductAddMetaApi(productId: number) {
  return requestClient.get<AdvanceProductAddMetaResult>('/shop/plus.advance.Product/add', {
    params: { product_id: productId },
  });
}

export async function getAdvanceProductEditMetaApi(advanceProductId: number) {
  return requestClient.get<AdvanceProductEditMetaResult>('/shop/plus.advance.Product/edit', {
    params: { advance_product_id: advanceProductId },
  });
}

export async function addAdvanceProductApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.advance.Product/add', payload);
}

export async function editAdvanceProductApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.advance.Product/edit', payload);
}
