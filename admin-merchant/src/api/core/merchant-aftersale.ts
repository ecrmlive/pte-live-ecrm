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
  return_shipment?: {
    carrier_name: string;
    remark?: string;
    submitted_at: string;
    tracking_no: string;
  };
  uid: number;
}

export interface MerchantRefundPage {
  limit: number;
  list: MerchantRefundOrder[];
  page: number;
  total: number;
}

export interface MerchantRefundEvent { actor_id: number; actor_type: 'merchant' | 'platform' | 'system' | 'user'; created_at: string; from_status: string; id: number; reason: string; to_status: string; }
export interface MerchantRefundEventPage { limit: number; list: MerchantRefundEvent[]; page: number; total: number; }
export interface MerchantRefundExport { content: string; file_name: string; row_count: number; truncated: boolean; }
export interface MerchantReturnShipment { carrier_name: string; remark?: string; submitted_at: string; tracking_no: string; }

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

export function listMerchantRefundEventsApi(id: number, params = { page: 1, limit: 100 }) { return requestClient.get<MerchantRefundEventPage>(`/refunds/${id}/events`, { params }); }
export function getMerchantRefundExpressApi(id: number) { return requestClient.get<MerchantReturnShipment>(`/refunds/${id}/express`); }
export function exportMerchantRefundsApi(params?: { status?: number }) { return requestClient.get<MerchantRefundExport>('/refunds/export', { params }); }
export function addMerchantRefundRemarkApi(id: number, input: { idempotency_key: string; note: string }) { return requestClient.post(`/refunds/${id}/remark`, input); }
export function hideMerchantRefundApi(id: number, input: { idempotency_key: string; reason: string }) { return requestClient.delete(`/refunds/${id}`, { data: input }); }

export function approveMerchantRefundApi(id: number) {
  return requestClient.post(`/refunds/${id}/approve`);
}

export function rejectMerchantRefundApi(id: number, failMessage: string) {
  return requestClient.post(`/refunds/${id}/reject`, { fail_message: failMessage });
}

export function confirmMerchantRefundReturnApi(id: number) {
  return requestClient.post(`/refunds/${id}/confirm-return`);
}
