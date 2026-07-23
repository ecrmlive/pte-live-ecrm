import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

/** 微信带货库 — 全部走 api-platform Go（/shop/plus.live.product/*），不经过 PHP :11500 */

export interface LiveWxProductItem {
  audit_status?: number;
  audit_status_text?: string;
  cover_img?: string;
  create_time?: string;
  goods_id?: number | string;
  name: string;
  price?: number | string;
  price2?: number | string;
  price_text?: string;
  price_type?: number;
  price_type_text?: string;
  product_id?: number;
  wx_product_id: number;
  [key: string]: unknown;
}

export interface LiveWxProductForm {
  cover_img: string;
  name: string;
  price: number | string;
  price2?: number | string;
  price_type: number;
  product_id: number | string;
  shop_supplier_id?: number;
  wx_product_id?: number;
}

export async function getLiveWxProductListApi(params: {
  list_rows?: number;
  page?: number;
  status?: number;
}) {
  return requestClient.post<{ list: PaginatedList<LiveWxProductItem> }>(
    '/shop/plus.live.product/index',
    params,
  );
}

export async function addLiveWxProductApi(payload: LiveWxProductForm) {
  return requestClient.post<{ msg?: string }>('/shop/plus.live.product/add', payload);
}

export async function editLiveWxProductApi(payload: LiveWxProductForm) {
  return requestClient.post<{ msg?: string }>('/shop/plus.live.product/edit', payload);
}

export async function deleteLiveWxProductApi(wxProductId: number) {
  return requestClient.post<{ msg?: string }>('/shop/plus.live.product/delete', {
    wx_product_id: wxProductId,
  });
}

export async function getLiveWxProductPickListApi(params: {
  list_rows?: number;
  page?: number;
  room_id: number | string;
}) {
  return requestClient.post<{
    excludeIds?: Array<number | string>;
    list: PaginatedList<LiveWxProductItem>;
  }>('/shop/plus.live.product/list', params);
}
