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
  status_code: string;
  status_time: string;
  return_shipment?: {
    carrier_name: string;
    remark?: string;
    submitted_at: string;
    tracking_no: string;
  };
  uid: number;
}

export interface PlatformRefundPage {
  limit: number;
  list: PlatformRefundOrder[];
  page: number;
  total: number;
}

export interface PlatformRefundEvent {
  actor_id: number;
  actor_type: 'merchant' | 'platform' | 'system' | 'user';
  created_at: string;
  from_status: string;
  id: number;
  reason: string;
  to_status: string;
}

export interface PlatformRefundEventPage {
  limit: number;
  list: PlatformRefundEvent[];
  page: number;
  total: number;
}

export interface PlatformRefundExport {
  content: string;
  file_name: string;
  row_count: number;
  truncated: boolean;
}

export function listPlatformRefundsApi(params: {
  limit: number;
  page: number;
  status?: string;
}) {
  return requestClient.get<PlatformRefundPage>('/refunds', { params });
}

export function getPlatformRefundApi(id: number) {
  return requestClient.get<PlatformRefundOrder>(`/refunds/${id}`);
}

export function listPlatformRefundEventsApi(id: number, params = { page: 1, limit: 100 }) {
  return requestClient.get<PlatformRefundEventPage>(`/refunds/${id}/events`, { params });
}

export function exportPlatformRefundsApi(input: { reason: string; status?: string }) {
  return requestClient.post<PlatformRefundExport>('/refunds/export', input);
}

export function approvePlatformRefundApi(id: number) {
  return requestClient.post(`/refunds/${id}/approve`);
}

export function rejectPlatformRefundApi(id: number, failMessage: string) {
  return requestClient.post(`/refunds/${id}/reject`, { fail_message: failMessage });
}
