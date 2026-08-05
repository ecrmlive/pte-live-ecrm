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
}

export function listCouponCommands(params: CouponCommandFilter) {
  return requestClient.get<{ list: CouponCommandRecord[]; total: number }>('/user-list/coupon-commands', { params });
}

export interface CouponReceiptRecord {
  id: number;
  user_id: number;
  coupon_id: number;
  coupon_name: string;
  store_id: number;
  source: string;
  status: 'unused' | 'locked' | 'used' | 'expired';
  obtained_at: string;
  used_order_id?: number;
}

export interface CouponReceiptFilter {
  page: number;
  limit: number;
  user_id?: number;
  coupon_id?: number;
  status?: CouponReceiptRecord['status'];
}

export function listCouponReceiptRecords(params: CouponReceiptFilter) {
  return requestClient.get<{ list: CouponReceiptRecord[]; total: number }>('/user-list/coupon-records', { params });
}
