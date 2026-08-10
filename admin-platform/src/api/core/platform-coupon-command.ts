import { requestClient } from '#/api/request';

export interface CouponCommandRecord {
  id: number;
  user_id: number;
  coupon_id: number;
  coupon_user_id: number;
  coupon_name: string;
  store_id: number;
  action: 'issue' | 'revoke';
  from_status: string;
  to_status: string;
  reason: string;
  operator_admin_id: number;
  created_at: string;
}

export interface CouponCommandFilter {
  page: number;
  limit: number;
  user_id?: number;
  coupon_id?: number;
  action?: CouponCommandRecord['action'];
  date_from?: string;
  date_to?: string;
}

export function listCouponCommands(params: CouponCommandFilter) {
  return requestClient.get<{ list: CouponCommandRecord[]; total: number }>(
    '/user-list/coupon-commands',
    { params },
  );
}

export interface CouponReceiptRecord {
  id: number;
  user_id: number;
  coupon_id: number;
  coupon_name: string;
  store_id: number;
  store_name?: string;
  source: string;
  source_name: string;
  status: 'unused' | 'locked' | 'used' | 'expired';
  obtained_at: string;
  used_at?: string | null;
  use_time?: string | null;
  used_order_id?: number;
  recipient: string;
  nickname?: string;
  avatar_url?: string;
  coupon_type_name: string;
  coupon_price: number;
  use_min_price: number;
  use_start_time?: string | null;
  use_end_time?: string | null;
  available: boolean;
}

export interface CouponReceiptFilter {
  page: number;
  limit: number;
  user_id?: number;
  coupon_id?: number;
  mer_id?: number;
  coupon_name?: string;
  status?: CouponReceiptRecord['status'];
  source?: string;
  coupon_scope?: 'platform' | 'store';
  recipient?: string;
  date_from?: string;
  date_to?: string;
}

export function listCouponReceiptRecords(params: CouponReceiptFilter) {
  return requestClient.get<{ list: CouponReceiptRecord[]; total: number }>(
    '/user-list/coupon-records',
    { params },
  );
}

/** 平台优惠券发送批次（对齐 CRMEB systemCouponSendLst） */
export interface PlatformCouponSendRecord {
  coupon_send_id: number;
  coupon_id: number;
  title: string;
  type: number;
  coupon_type_name: string;
  create_time: string;
  coupon_type: number;
  coupon_time: number;
  use_start_time?: string | null;
  use_end_time?: string | null;
  validity_text: string;
  mark: string;
  filter_text: string;
  coupon_num: number;
  use_count: number;
  send_status: number;
}

export interface PlatformCouponSendDetail extends PlatformCouponSendRecord {
  coupon_price: number;
  use_min_price: number;
  is_timeout: number;
  start_time?: string | null;
  end_time?: string | null;
  send_type: number;
  send_type_name: string;
  is_limited: number;
  total_count: number;
  remain_count: number;
  status: number;
  sort: number;
  sent_total: number;
  used_total: number;
}

export interface PlatformCouponSendUser {
  user_id: number;
  nickname: string;
  avatar_url: string;
  source: string;
  source_name: string;
  status: string;
  status_name: string;
}

export function listPlatformCouponSendsApi(params: {
  page: number;
  limit: number;
  date_from?: string;
  date_to?: string;
  coupon_type?: number;
  coupon_name?: string;
  status?: number;
}) {
  return requestClient.get<{
    list: PlatformCouponSendRecord[];
    total: number;
    page: number;
    limit: number;
  }>('/coupons/sends', { params });
}

export function getPlatformCouponSendDetailApi(id: number) {
  return requestClient.get<PlatformCouponSendDetail>(`/coupons/sends/${id}`);
}

export function listPlatformCouponSendUsersApi(
  id: number,
  params: { page: number; limit: number },
) {
  return requestClient.get<{
    list: PlatformCouponSendUser[];
    total: number;
    page: number;
    limit: number;
  }>(`/coupons/sends/${id}/users`, { params });
}
