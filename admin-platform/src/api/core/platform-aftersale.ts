import { requestClient } from '#/api/request';

export interface PlatformRefundProduct {
  create_time: string;
  order_product_id: number;
  refund_price: number;
  refund_product_id: number;
  refund_num: number;
}

export interface PlatformRefundOrder {
  create_time: string;
  fail_message?: string;
  mer_id: number;
  order_id: number;
  products?: PlatformRefundProduct[];
  refund_message: string;
  refund_num: number;
  refund_order_id: number;
  refund_order_sn: string;
  refund_price: number;
  refund_type: number;
  status: number;
  status_time: string;
  uid: number;
}

export interface PlatformRefundPage {
  limit: number;
  list: PlatformRefundOrder[];
  page: number;
  total: number;
}

export function listPlatformRefundsApi(params: {
  limit: number;
  page: number;
  status?: number;
}) {
  return requestClient.get<PlatformRefundPage>('/refunds', { params });
}

export function getPlatformRefundApi(id: number) {
  return requestClient.get<PlatformRefundOrder>(`/refunds/${id}`);
}

export function approvePlatformRefundApi(id: number) {
  return requestClient.post(`/refunds/${id}/approve`);
}

export function rejectPlatformRefundApi(id: number, failMessage: string) {
  return requestClient.post(`/refunds/${id}/reject`, { fail_message: failMessage });
}
