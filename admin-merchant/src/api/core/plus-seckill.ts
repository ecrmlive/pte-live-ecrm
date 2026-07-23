import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface SeckillProductItem {
  create_time?: string;
  is_delete?: number;
  limit_num?: number;
  product?: { image?: Array<{ file_path: string }> };
  product_id?: number;
  product_name?: string;
  product_price?: string | number;
  seckill_activity_id?: number;
  seckill_price?: string | number;
  seckill_product_id: number;
  stock?: number;
  title?: string;
  total_sales?: number;
  active?: {
    end_time_text?: string;
    start_time_text?: string;
    status_text?: string;
  };
}

export interface SeckillProductListQuery {
  active_status?: number;
  list_rows?: number;
  page?: number;
  search?: string;
  seckill_activity_id?: number;
  status?: number;
}

export async function getSeckillProductListApi(params: SeckillProductListQuery) {
  return requestClient.post<{
    activeList: Array<{ seckill_activity_id: number; title: string }>;
    list: PaginatedList<SeckillProductItem>;
  }>('/shop/plus.seckill.product/index', params);
}

export async function setSeckillProductStateApi(payload: {
  is_delete: number;
  seckill_product_id: number;
}) {
  return requestClient.post('/shop/plus.seckill.product/state', payload);
}

export interface SeckillTimeItem {
  end_time?: string;
  end_time_text?: string;
  id: number;
  start_time?: string;
  start_time_text?: string;
  status?: number;
  title?: string;
}

export async function getSeckillTimeListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<SeckillTimeItem> }>(
    '/shop/plus.seckill.time/index',
    params,
  );
}

export async function addSeckillTimeApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.seckill.time/add', payload);
}

export async function editSeckillTimeApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.seckill.time/edit', payload);
}

export async function deleteSeckillTimeApi(id: number) {
  return requestClient.post('/shop/plus.seckill.time/delete', { id });
}

export async function setSeckillTimeStateApi(payload: { id: number; status: number }) {
  return requestClient.post('/shop/plus.seckill.time/state', payload);
}

export interface SeckillStatisticsDetail {
  payPerson?: number | string;
  stock?: number | string;
  totalPerson?: number | string;
  totalMoney?: number | string;
  total_stock?: number | string;
}

export async function getSeckillStatisticsDetailApi(params: {
  product_id?: number | string;
  seckill_activity_id?: number | string;
  seckill_product_id?: number | string;
}) {
  return requestClient.post<{ data: SeckillStatisticsDetail }>(
    '/shop/plus.seckill.product/statistics',
    params,
  );
}

export async function getSeckillJoinListApi(params: Record<string, unknown>) {
  return requestClient.post<{ list: PaginatedList<Record<string, unknown>> }>(
    '/shop/plus.seckill.product/join',
    params,
  );
}

export async function getSeckillOrderListApi(params: Record<string, unknown>) {
  return requestClient.post<{ list: PaginatedList<Record<string, unknown>> }>(
    '/shop/plus.seckill.product/order',
    params,
  );
}

export interface SeckillActiveListItem {
  create_time?: string;
  end_time_text?: string;
  product_num?: number;
  seckill_activity_id: number;
  sort?: number;
  start_time_text?: string;
  status?: number;
  status_text?: string;
  timeList?: string[];
  title: string;
  total_sales?: number;
}

export async function getSeckillActiveListApi(params: Record<string, unknown>) {
  return requestClient.post<{ list: PaginatedList<SeckillActiveListItem> }>(
    '/shop/plus.seckill.Active/index',
    params,
  );
}

export async function deleteSeckillActiveApi(seckillActivityId: number) {
  return requestClient.post('/shop/plus.seckill.Active/delete', {
    seckill_activity_id: seckillActivityId,
  });
}

export async function setSeckillActiveStateApi(payload: {
  seckill_activity_id: number;
  status: number;
}) {
  return requestClient.post('/shop/plus.seckill.Active/state', payload);
}

export async function getSeckillActiveAddMetaApi() {
  return requestClient.get<{ timeList: unknown[] }>('/shop/plus.seckill.Active/add');
}

export async function addSeckillActiveApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.seckill.Active/add', payload);
}

export async function getSeckillActiveEditMetaApi(seckillActivityId: number) {
  return requestClient.get<{
    detail: Record<string, unknown>;
    product_list: Array<Record<string, unknown>>;
    timeList: unknown[];
  }>('/shop/plus.seckill.Active/edit', { params: { seckill_activity_id: seckillActivityId } });
}

export async function editSeckillActiveApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.seckill.Active/edit', payload);
}
