import { http } from '@/api/http';

export interface RefundOrder {
  refund_order_id: number;
  refund_order_sn: string;
  order_id: number;
  uid: number;
  mer_id: number;
  refund_type: number;
  refund_message: string;
  refund_price: number;
  refund_num: number;
  status: number;
  fail_message?: string;
  create_time: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchRefunds(params: Record<string, unknown>) {
  return http.get<PageResult<RefundOrder>>('/refunds', { params });
}

export function fetchRefund(id: number) {
  return http.get<RefundOrder>(`/refunds/${id}`);
}

export function approveRefund(id: number) {
  return http.post<{ ok: boolean }>(`/refunds/${id}/approve`, {});
}

export function rejectRefund(id: number, data: { fail_message: string }) {
  return http.post<{ ok: boolean }>(`/refunds/${id}/reject`, data);
}
