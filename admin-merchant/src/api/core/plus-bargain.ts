import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface BargainProductItem {
  bargain_num?: number;
  bargain_price?: string | number;
  bargain_product_id: number;
  create_time?: string;
  product_id?: number;
  image?: { file_path?: string };
  product?: { product_price?: string | number };
  product_name?: string;
  sort?: number;
  status?: number;
  status_text?: string;
  stock?: number;
  totalNum?: number;
  totalPerson?: number;
  helpNum?: number;
  successNum?: number;
  start_time_text?: string;
  end_time_text?: string;
}

export async function getBargainProductListApi(params: Record<string, unknown>) {
  return requestClient.post<{
    exclude_ids?: number[];
    list: PaginatedList<BargainProductItem>;
  }>('/shop/plus.bargain.product/index', params);
}

export async function deleteBargainProductApi(bargainProductId: number) {
  return requestClient.post('/shop/plus.bargain.product/delete', {
    bargain_product_id: bargainProductId,
  });
}

export async function setBargainProductStateApi(payload: {
  bargain_product_id: number;
  status: number;
}) {
  return requestClient.post('/shop/plus.bargain.product/state', payload);
}

export async function getBargainSettingApi() {
  return requestClient.get<{ vars: { values: Record<string, unknown> } }>(
    '/shop/plus.bargain.setting/index',
  );
}

export async function saveBargainSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.bargain.setting/index', payload);
}

export interface BargainProductFormMeta {
  delivery: Array<Record<string, unknown>>;
  logistics: Array<{ value: number | string }>;
  model: Record<string, unknown>;
}

export async function getBargainProductAddMetaApi(productId: number) {
  return requestClient.get<BargainProductFormMeta>('/shop/plus.bargain.product/add', {
    params: { product_id: productId },
  });
}

export async function getBargainProductEditMetaApi(bargainProductId: number) {
  return requestClient.get<BargainProductFormMeta>('/shop/plus.bargain.product/edit', {
    params: { bargain_product_id: bargainProductId },
  });
}

export async function addBargainProductApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.bargain.product/add', payload);
}

export async function editBargainProductApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.bargain.product/edit', payload);
}

export interface BargainStatisticsDetail {
  createNum?: number | string;
  helpNum?: number | string;
  payNum?: number | string;
  successNum?: number | string;
  totalMoney?: number | string;
  totalNum?: number | string;
}

export async function getBargainStatisticsDetailApi(params: {
  bargain_product_id?: number | string;
  product_id?: number | string;
}) {
  return requestClient.post<{ data: BargainStatisticsDetail }>(
    '/shop/plus.bargain.product/statistics',
    params,
  );
}

export async function getBargainJoinListApi(params: Record<string, unknown>) {
  return requestClient.post<{ list: PaginatedList<Record<string, unknown>> }>(
    '/shop/plus.bargain.product/join',
    params,
  );
}

export async function getBargainOrderListApi(params: Record<string, unknown>) {
  return requestClient.post<{ list: PaginatedList<Record<string, unknown>> }>(
    '/shop/plus.bargain.product/order',
    params,
  );
}

export async function getBargainTaskListApi(params: Record<string, unknown>) {
  return requestClient.get<{ list: PaginatedList<Record<string, unknown>> }>(
    '/shop/plus.bargain.task/index',
    { params },
  );
}

export async function getBargainTaskDetailApi(params: {
  bargain_task_id: number;
  list_rows?: number;
  page?: number;
}) {
  return requestClient.get<{ list: PaginatedList<Record<string, unknown>> }>(
    '/shop/plus.bargain.task/detail',
    { params },
  );
}

export async function getBargainProductDetailApi(params: {
  bargain_task_id: number;
  list_rows?: number;
  page?: number;
}) {
  return requestClient.post<{ list: PaginatedList<Record<string, unknown>> }>(
    '/shop/plus.bargain.product/detail',
    params,
  );
}
