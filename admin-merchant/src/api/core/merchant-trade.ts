import { requestClient } from '#/api/request';

export interface MerchantOrderProduct {
  order_product_id: number;
  product_id: number;
  product_info: string;
  product_num: number;
  product_price: number;
  product_sku: string;
  total_price: number;
}

export interface MerchantOrder {
  can_verify?: boolean;
  create_time: string;
  delivery_id: string;
  delivery_name: string;
  delivery_type: string;
  has_verify_code?: boolean;
  mark: string;
  order_id: number;
  order_sn: string;
  paid: number;
  pay_price: number;
  pay_type: number;
  products?: MerchantOrderProduct[];
  real_name: string;
  status: number;
  total_num: number;
  user_address: string;
  user_phone: string;
  verify_code?: string;
  verify_status?: string;
}

export interface MerchantOrderPage {
  limit: number;
  list: MerchantOrder[];
  page: number;
  total: number;
}

export function listMerchantOrdersApi(params: {
  limit: number;
  page: number;
  paid?: number;
  status?: number;
  verify_tab?: 'pending' | 'verified';
  keyword?: string;
  order_sn?: string;
  pay_type?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<MerchantOrderPage>('/orders', { params });
}

export function getMerchantOrderApi(id: number) {
  return requestClient.get<MerchantOrder>(`/orders/${id}`);
}

export function deliverMerchantOrderApi(
  id: number,
  body: { delivery_id: string; delivery_name: string; delivery_type: string },
) {
  return requestClient.post(`/orders/${id}/delivery`, body);
}

export function verifyMerchantOrderApi(id: number, verifyCode?: string) {
  return requestClient.post(`/orders/${id}/verify`, verifyCode ? { verify_code: verifyCode } : {});
}
