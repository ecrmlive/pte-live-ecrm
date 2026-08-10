import { requestClient } from '#/api/request';

export interface PlatformIntegralOrderProduct {
  order_product_id?: number;
  product_id?: number;
  product_image?: string;
  product_info?: string;
  product_num?: number;
  product_price?: number;
  product_sku?: string;
  total_price?: number;
}

export interface PlatformIntegralOrderRow {
  can_delete?: boolean;
  can_deliver?: boolean;
  create_time: string;
  delivery_id?: string;
  delivery_name?: string;
  delivery_type?: string;
  delivery_type_label?: string;
  freight_price?: number;
  group_order_id?: number;
  mer_id?: number;
  mer_name?: string;
  merchant_remark?: string;
  nickname?: string;
  order_id: number;
  order_sn: string;
  pay_amount?: number;
  pay_price?: number;
  points_amount: number;
  product?: PlatformIntegralOrderProduct;
  products?: PlatformIntegralOrderProduct[];
  real_name?: string;
  status: number;
  status_label: string;
  store_id?: number;
  store_name?: string;
  total_num?: number;
  total_price?: number;
  uid: number;
  user_address?: string;
  user_deleted?: boolean;
  user_phone?: string;
  user_remark?: string;
}

export interface PlatformIntegralOrderPage {
  limit: number;
  list: PlatformIntegralOrderRow[];
  page: number;
  total: number;
}

export interface PlatformIntegralOrderExport {
  content: string;
  file_name: string;
  row_count: number;
  truncated: boolean;
}

export type PlatformIntegralOrderQuery = {
  date_from?: string;
  date_to?: string;
  keyword?: string;
  limit?: number;
  page?: number;
  search_type?: string;
  status?: string | number;
};

export function listPlatformIntegralOrdersApi(
  params: PlatformIntegralOrderQuery,
) {
  return requestClient.get<PlatformIntegralOrderPage>('/integral/orders', {
    params,
  });
}

export function getPlatformIntegralOrderApi(id: number) {
  return requestClient.get<PlatformIntegralOrderRow>(`/integral/orders/${id}`);
}

export function exportPlatformIntegralOrdersApi(
  input: PlatformIntegralOrderQuery,
) {
  return requestClient.post<PlatformIntegralOrderExport>(
    '/integral/orders/export',
    input,
  );
}

export function deliverPlatformIntegralOrderApi(
  id: number,
  payload: {
    delivery_id?: string;
    delivery_name?: string;
    delivery_type: string;
    remark?: string;
  },
) {
  return requestClient.post<{ order_id: number; status: string }>(
    `/integral/orders/${id}/delivery`,
    payload,
  );
}

export function deletePlatformIntegralOrderApi(id: number) {
  return requestClient.delete<{ order_id: number }>(`/integral/orders/${id}`);
}
