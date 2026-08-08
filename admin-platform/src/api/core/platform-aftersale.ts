import { requestClient } from '#/api/request';

export interface PlatformRefundProduct {
  order_product_id?: number;
  pay_price?: number;
  product_id: number;
  product_image?: string;
  product_info?: string;
  product_num?: number;
  product_price?: number;
  product_sku?: string;
  refund_num: number;
  refund_price: number;
  refund_product_id?: number;
}

export interface PlatformRefundOrder {
  create_time: string;
  fail_message?: string;
  is_trader?: number;
  mer_id: number;
  mer_name?: string;
  merchant_remark?: string;
  nickname?: string;
  order_id: number;
  order_sn?: string;
  order_user_remark?: string;
  product_total_price?: number;
  products?: PlatformRefundProduct[];
  refund_evidence?: string[];
  refund_initiator?: string;
  refund_message: string;
  refund_method?: string;
  refund_num: number;
  refund_order_id: number;
  refund_order_sn: string;
  refund_price: number;
  refund_type: number;
  refund_type_label?: string;
  status: number;
  status_code: string;
  status_label?: string;
  status_time: string;
  store_category_name?: string;
  store_name?: string;
  user_deleted?: boolean;
  user_phone_mask?: string;
  user_remark?: string;
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

export interface PlatformRefundTabCounts {
  all: number;
  applied: number;
  approved: number;
  awaiting_receipt: number;
  completed: number;
  dispute: number;
  rejected: number;
}

export interface PlatformRefundEvent {
  actor_id: number;
  actor_type: 'merchant' | 'platform' | 'system' | 'user';
  content?: string;
  created_at: string;
  from_status: string;
  id: number;
  operate_time?: string;
  operator?: string;
  order_sn?: string;
  reason: string;
  role?: string;
  terminal?: string;
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

export interface PlatformRefundListParams {
  date_from?: string;
  date_to?: string;
  is_trader?: number | string;
  limit: number;
  order_sn?: string;
  page: number;
  phone?: string;
  real_name?: string;
  refund_order_sn?: string;
  refund_type?: number;
  status?: string;
  tab_status?: string;
  user_search_keyword?: string;
  user_search_type?: string;
}

export function listPlatformRefundsApi(params: PlatformRefundListParams) {
  return requestClient.get<PlatformRefundPage>('/refunds', { params });
}

export function getPlatformRefundTabCountsApi(
  params?: Omit<PlatformRefundListParams, 'page' | 'limit' | 'tab_status' | 'status'>,
) {
  return requestClient.get<PlatformRefundTabCounts>('/refunds/tab-counts', {
    params,
  });
}

export function getPlatformRefundApi(id: number) {
  return requestClient.get<PlatformRefundOrder>(`/refunds/${id}`);
}

export function listPlatformRefundEventsApi(
  id: number,
  params: {
    date_from?: string;
    date_to?: string;
    limit?: number;
    page?: number;
    terminal?: string;
  } = { page: 1, limit: 100 },
) {
  return requestClient.get<PlatformRefundEventPage>(`/refunds/${id}/events`, {
    params,
  });
}

export function exportPlatformRefundsApi(
  input: Pick<
    PlatformRefundListParams,
    | 'status'
    | 'tab_status'
    | 'refund_order_sn'
    | 'refund_type'
    | 'order_sn'
    | 'phone'
    | 'real_name'
    | 'date_from'
    | 'date_to'
    | 'is_trader'
    | 'user_search_type'
    | 'user_search_keyword'
  > & { reason: string },
) {
  return requestClient.post<PlatformRefundExport>('/refunds/export', input);
}

export function approvePlatformRefundApi(id: number) {
  return requestClient.post(`/refunds/${id}/approve`);
}

export function rejectPlatformRefundApi(id: number, failMessage: string) {
  return requestClient.post(`/refunds/${id}/reject`, {
    reason: failMessage,
    fail_message: failMessage,
  });
}
