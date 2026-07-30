import { requestClient } from '#/api/request';

export interface MerchantRefundProduct {
  create_time: string;
  order_product_id: number;
  refund_price: number;
  refund_product_id: number;
  refund_num: number;
}

export interface MerchantRefundOrder {
  create_time: string;
  fail_message?: string;
  order_id: number;
  products?: MerchantRefundProduct[];
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

export interface MerchantRefundPage {
  limit: number;
  list: MerchantRefundOrder[];
  page: number;
  total: number;
}

export function listMerchantRefundsApi(params: {
  limit: number;
  page: number;
  status?: number;
}) {
  return requestClient.get<MerchantRefundPage>('/refunds', { params });
}

export function getMerchantRefundApi(id: number) {
  return requestClient.get<MerchantRefundOrder>(`/refunds/${id}`);
}

export function approveMerchantRefundApi(id: number) {
  return requestClient.post(`/refunds/${id}/approve`);
}

export function rejectMerchantRefundApi(id: number, failMessage: string) {
  return requestClient.post(`/refunds/${id}/reject`, { fail_message: failMessage });
}
