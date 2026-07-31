import { http } from "@/utils/request";

export interface RefundItem {
  refund_product_id: number;
  order_product_id: number;
  refund_price: number;
  refund_num: number;
}

export interface RefundOrder {
  refund_order_id: number;
  refund_order_sn: string;
  order_id: number;
  mer_id: number;
  refund_type: number;
  refund_message: string;
  refund_price: number;
  refund_num: number;
  status: number;
  status_code: "applied" | "merchant_handling" | "platform_intervene" | "refunding" | "refunded" | "rejected" | "cancelled";
  create_time: string;
  status_time: string;
  products: RefundItem[];
}

export function fetchRefunds(page = 1, limit = 20) {
  return http.get<{ list: RefundOrder[]; total: number; page: number; limit: number }>(`/refunds?page=${page}&limit=${limit}`);
}

export function applyRefund(orderID: number, refundMessage: string, idempotencyKey: string) {
  return http.post<RefundOrder>("/refund/apply", {
    order_id: orderID,
    refund_type: 1,
    refund_message: refundMessage,
    idempotency_key: idempotencyKey,
  });
}

export function cancelRefund(refundID: number) {
  return http.post<{ ok: boolean }>(`/refunds/${refundID}/cancel`);
}

export function requestPlatformIntervention(refundID: number) {
  return http.post<{ ok: boolean }>(`/refunds/${refundID}/platform`);
}
