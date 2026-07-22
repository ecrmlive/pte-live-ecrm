import { http } from './http';

export function createRefund(order_id: number, refund_message: string) {
  return http.post('/refunds', {
    order_id,
    refund_type: 1,
    refund_message,
  });
}

export function approveRefund(id: number) {
  return http.post(`/refunds/${id}/approve`, {});
}

export function fetchRefunds(params: Record<string, unknown>) {
  return http.get('/refunds', { params });
}
