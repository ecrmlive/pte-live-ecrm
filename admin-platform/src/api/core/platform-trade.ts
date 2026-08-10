import { requestClient } from '#/api/request';

export interface PlatformOrderProduct {
  activity_id?: number;
  cost_price?: number;
  pay_price?: number;
  product_id: number;
  product_image?: string;
  product_info: string;
  product_num: number;
  product_price: number;
  product_sku?: string;
  total_price: number;
}

export interface PlatformOrder {
  activity_type?: number;
  activity_type_label?: string;
  commission_total?: number;
  create_time: string;
  delivery_id: string;
  delivery_name: string;
  delivery_type: string;
  delivery_type_label?: string;
  discount_amount?: number;
  first_brokerage?: number;
  freight_price?: number;
  group_order_id?: number;
  group_order_sn?: string;
  integral_price?: number;
  member_discount?: number;
  mer_category_name?: string;
  mer_id: number;
  mer_name?: string;
  merchant_coupon?: number;
  merchant_remark?: string;
  nickname?: string;
  order_id: number;
  order_sn: string;
  order_type_label?: string;
  paid: number;
  pay_price: number;
  pay_status_label?: string;
  pay_time?: string;
  pay_type: number;
  pay_type_label?: string;
  platform_coupon?: number;
  points_deduction?: number;
  product?: PlatformOrderProduct;
  product_type_label?: string;
  products?: PlatformOrderProduct[];
  real_name?: string;
  second_brokerage?: number;
  spread_name?: string;
  status: number;
  status_label?: string;
  store_category_name?: string;
  store_id?: number;
  store_name?: string;
  store_type_name?: string;
  top_spread_name?: string;
  total_num: number;
  total_price: number;
  uid?: number;
  user_address: string;
  user_deleted?: boolean;
  user_phone: string;
  user_phone_mask?: string;
  user_remark?: string;
}

export interface PlatformOrderPage {
  limit: number;
  list: PlatformOrder[];
  page: number;
  total: number;
}

export interface PlatformOrderTabCounts {
  all: number;
  completed: number;
  deleted: number;
  refunded: number;
  unpaid: number;
  unreceived: number;
  unshipped: number;
  unevaluated: number;
}

export interface PlatformOrderLogRow {
  content: string;
  operate_time: string;
  operator: string;
  order_sn: string;
  role: string;
  terminal?: string;
}

export function listPlatformOrdersApi(params: {
  activity_type?: number;
  date_from?: string;
  date_to?: string;
  delivery_type?: string;
  is_spread?: 0 | 1;
  keyword?: string;
  limit: number;
  mer_category_id?: number;
  mer_id?: number;
  mer_type_id?: number;
  order_search_keyword?: string;
  order_search_type?: string;
  order_sn?: string;
  page: number;
  paid?: number;
  pay_type?: number;
  phone?: string;
  product_name?: string;
  product_type?: number;
  real_name?: string;
  spread_keyword?: string;
  status?: number;
  store_id?: number;
  tab_status?: string;
  top_spread_keyword?: string;
  user_search_keyword?: string;
  user_search_type?: string;
}) {
  return requestClient.get<PlatformOrderPage>('/orders', { params });
}

export function getPlatformOrderTabCountsApi(params?: {
  activity_type?: number;
  date_from?: string;
  date_to?: string;
  delivery_type?: string;
  is_spread?: 0 | 1;
  keyword?: string;
  mer_category_id?: number;
  mer_id?: number;
  mer_type_id?: number;
  order_search_keyword?: string;
  order_search_type?: string;
  pay_type?: number;
  product_name?: string;
  product_type?: number;
  spread_keyword?: string;
  store_id?: number;
  top_spread_keyword?: string;
  user_search_keyword?: string;
  user_search_type?: string;
}) {
  return requestClient.get<PlatformOrderTabCounts>('/orders/tab-counts', {
    params,
  });
}

export function getPlatformOrderApi(id: number) {
  return requestClient.get<PlatformOrder>(`/orders/${id}`);
}

export function listPlatformOrderLogsApi(
  id: number,
  params: {
    date_from?: string;
    date_to?: string;
    limit: number;
    page: number;
    terminal?: string;
  },
) {
  return requestClient.get<{ list: PlatformOrderLogRow[]; total: number }>(
    `/orders/${id}/logs`,
    { params },
  );
}

/** 核销记录（已核销订单只读监管） */
export interface PlatformVerifyRecord extends PlatformOrder {
  verifier_account_id?: number;
  verifier_name?: string;
  verify_status?: string;
  verify_status_label?: string;
  verify_time?: string;
}

export interface PlatformVerifySummary {
  alipay_amount: number;
  balance_amount: number;
  paid_count: number;
  pay_amount: number;
  refund_amount: number;
  wechat_amount: number;
}

export type PlatformVerifyListParams = {
  date_from?: string;
  date_to?: string;
  is_trader?: number;
  limit: number;
  order_keyword?: string;
  page: number;
  pay_type?: number;
  user_keyword?: string;
};

export function listPlatformVerifyRecordsApi(params: PlatformVerifyListParams) {
  return requestClient.get<{
    limit: number;
    list: PlatformVerifyRecord[];
    page: number;
    total: number;
  }>('/orders/verify-records', { params });
}

export function getPlatformVerifySummaryApi(
  params?: Omit<PlatformVerifyListParams, 'page' | 'limit'>,
) {
  return requestClient.get<PlatformVerifySummary>(
    '/orders/verify-records/summary',
    { params },
  );
}
