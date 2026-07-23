import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface PackageListItem {
  code_type: number;
  end_time?: { text: string };
  gift_package_id: number;
  name: string;
  people: number;
  start_time?: { text: string };
  status?: { text: string; value: number };
  total_num: number;
}

export interface PackageOrderItem {
  create_time: string;
  order_id: number;
  order_no: string;
  pay_price: number | string;
  pay_status?: { text: string };
  pay_type?: { text: string };
  user?: { nickName: string };
}

export async function getPackageListApi(params: { page?: number; search?: string }) {
  return requestClient.post<{ list: PaginatedList<PackageListItem> }>(
    '/shop/plus.package/index',
    params,
  );
}

export async function deletePackageApi(id: number) {
  return requestClient.post('/shop/plus.package/delete', { id });
}

export async function sendPackageApi(id: number) {
  return requestClient.post('/shop/plus.package/send', { id });
}

export async function endPackageApi(id: number) {
  return requestClient.post('/shop/plus.package/end', { id });
}

export async function getPackageOrderListApi(params: {
  id?: number;
  list_rows?: number;
  page?: number;
  search?: string;
}) {
  return requestClient.post<{ list: PaginatedList<PackageOrderItem> }>(
    '/shop/plus.package/orderlist',
    params,
  );
}

export async function addPackageApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.package/add', payload);
}

export async function getPackageEditMetaApi(giftPackageId: number) {
  return requestClient.get<{ data: Record<string, unknown> }>('/shop/plus.package/edit', {
    params: { gift_package_id: giftPackageId },
  });
}

export async function editPackageApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.package/edit', payload);
}

export async function getPackageQrcodePreviewApi(params: { id: number; source: string }) {
  return requestClient.post<{ image: string }>('/shop/plus.package/qrcode/preview', params);
}
