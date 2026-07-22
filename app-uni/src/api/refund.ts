import { http } from "@/utils/request";

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
  fail_message?: string;
  create_time: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function applyRefund(data: {
  order_id: number;
  refund_type?: number;
  refund_message: string;
  order_product_ids?: number[];
}) {
  return http.post<RefundOrder>("/refund/apply", {
    refund_type: 1,
    ...data,
  });
}

export function fetchRefunds(page = 1, limit = 20) {
  return http.get<PageResult<RefundOrder>>(`/refunds?page=${page}&limit=${limit}`);
}

export function fetchRefund(id: number) {
  return http.get<RefundOrder>(`/refunds/${id}`);
}

export function cancelRefund(id: number) {
  return http.post<{ ok: boolean }>(`/refunds/${id}/cancel`, {});
}

export function requestPlatformRefund(id: number) {
  return http.post<{ ok: boolean }>(`/refunds/${id}/platform`, {});
}

export function refundStatusText(status: number) {
  switch (status) {
    case 0:
      return "待审核";
    case -1:
      return "已拒绝";
    case 1:
      return "待退货";
    case 2:
      return "待收货";
    case 3:
      return "已退款";
    case 4:
      return "平台介入";
    case -2:
      return "已取消";
    default:
      return `状态${status}`;
  }
}
