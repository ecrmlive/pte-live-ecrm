import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface BuyActivityItem {
  buy_id: number;
  create_time?: number | string;
  end_time_text: string;
  name: string;
  sort: number;
  start_time_text: string;
  status_text: string;
}

export interface BuyActivityListQuery {
  list_rows?: number;
  name?: string;
  page?: number;
  status?: number;
}

export async function getBuyActivityListApi(params: BuyActivityListQuery) {
  return requestClient.post<{ list: PaginatedList<BuyActivityItem> }>(
    '/shop/plus.buyactivity/index',
    params,
  );
}

export async function deleteBuyActivityApi(buyId: number) {
  return requestClient.post('/shop/plus.buyactivity/delete', { buy_id: buyId });
}

export async function getBuyActivityEditMetaApi(buyId: number) {
  return requestClient.get<BuyActivityEditMetaResult>('/shop/plus.buyactivity/edit', {
    params: { buy_id: buyId },
  });
}

export async function addBuyActivityApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.buyactivity/add', payload);
}

export async function editBuyActivityApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.buyactivity/edit', payload);
}

export interface BuyActivityProductRow {
  product_attr?: string;
  product_id: number;
  product_name: string;
  product_num: number | string;
  product_sku_id?: number;
  spec_sku_id?: number | string;
  spec_type?: number;
}

export interface BuyActivityLimitProductRow {
  product_id: number;
  product_name: string;
  product_num: number | string;
}

export interface BuyActivityFormModel {
  buy_id?: number;
  end_time: string;
  max_times: number;
  name: string;
  send_type: number;
  sort: number | string;
  start_time: string;
  status: number;
}

export interface BuyActivityEditMetaResult {
  model: BuyActivityFormModel & {
    end_time_text?: string;
    limit_product?: BuyActivityLimitProductRow[];
    product_ids?: BuyActivityProductRow[];
    start_time_text?: string;
  };
}
