import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface AssembleProductItem {
  assemble_num?: number;
  assemble_price?: string | number;
  assemble_product_id: number;
  create_time?: string;
  product_id?: number;
  image?: { file_path?: string };
  product?: { product_price?: string | number };
  product_name?: string;
  sort?: number;
  status?: number;
  status_text?: string;
  stock?: number;
  successNum?: number;
  totalNum?: number;
  totalPerson?: number;
  start_time_text?: string;
  end_time_text?: string;
}

export async function getAssembleProductListApi(params: Record<string, unknown>) {
  return requestClient.post<{
    exclude_ids?: number[];
    list: PaginatedList<AssembleProductItem>;
  }>('/shop/plus.assemble.product/index', params);
}

export async function deleteAssembleProductApi(assembleProductId: number) {
  return requestClient.post('/shop/plus.assemble.product/delete', {
    assemble_product_id: assembleProductId,
  });
}

export async function setAssembleProductStateApi(payload: {
  assemble_product_id: number;
  status: number;
}) {
  return requestClient.post('/shop/plus.assemble.product/state', payload);
}

export async function getAssembleSettingApi() {
  return requestClient.get<{ vars: { values: Record<string, unknown> } }>(
    '/shop/plus.assemble.setting/index',
  );
}

export async function saveAssembleSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.assemble.setting/index', payload);
}

export interface AssembleProductFormMeta {
  delivery: Array<Record<string, unknown>>;
  logistics: Array<{ value: number | string }>;
  model: Record<string, unknown>;
}

export async function getAssembleProductAddMetaApi(productId: number) {
  return requestClient.get<AssembleProductFormMeta>('/shop/plus.assemble.product/add', {
    params: { product_id: productId },
  });
}

export async function getAssembleProductEditMetaApi(assembleProductId: number) {
  return requestClient.get<AssembleProductFormMeta>('/shop/plus.assemble.product/edit', {
    params: { assemble_product_id: assembleProductId },
  });
}

export async function addAssembleProductApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.assemble.product/add', payload);
}

export async function editAssembleProductApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.assemble.product/edit', payload);
}

export interface AssembleStatisticsDetail {
  payPerson?: number | string;
  successNum?: number | string;
  totalMoney?: number | string;
  totalNum?: number | string;
  totalPerson?: number | string;
}

export async function getAssembleStatisticsDetailApi(params: {
  assemble_product_id?: number | string;
  product_id?: number | string;
}) {
  return requestClient.post<{ data: AssembleStatisticsDetail }>(
    '/shop/plus.assemble.product/statistics',
    params,
  );
}

export async function getAssembleJoinListApi(params: Record<string, unknown>) {
  return requestClient.post<{ list: PaginatedList<Record<string, unknown>> }>(
    '/shop/plus.assemble.product/join',
    params,
  );
}

export async function getAssembleOrderListApi(params: Record<string, unknown>) {
  return requestClient.post<{ list: PaginatedList<Record<string, unknown>> }>(
    '/shop/plus.assemble.product/order',
    params,
  );
}

export interface AssembleRecordSummary {
  fail?: number | string;
  joinNum?: number | string;
  success?: number | string;
}

export async function getAssembleRecordListApi(params: Record<string, unknown>) {
  return requestClient.get<{
    data: AssembleRecordSummary;
    list: PaginatedList<Record<string, unknown>>;
  }>('/shop/plus.assemble.record/index', { params });
}

export async function finishAssembleRecordApi(assembleBillId: number) {
  return requestClient.post('/shop/plus.assemble.record/finish', {
    assemble_bill_id: assembleBillId,
  });
}

export async function getAssembleRecordDetailApi(params: {
  assemble_bill_id: number;
  list_rows?: number;
  page?: number;
}) {
  return requestClient.post<{ list: PaginatedList<Record<string, unknown>> }>(
    '/shop/plus.assemble.record/detail',
    params,
  );
}
