import { requestClient } from '#/api/request';

export interface PlatformCoupon {
  coupon_id: number;
  coupon_price: number;
  coupon_time: number;
  coupon_type?: number;
  mer_id?: number;
  create_time: string;
  end_time?: string | null;
  full_reduction?: number;
  is_limited: number;
  is_timeout?: number;
  remain_count: number;
  send_type?: number;
  sort: number;
  start_time?: string | null;
  status: number;
  title: string;
  total_count: number;
  type?: number;
  use_end_time?: string | null;
  use_min_price: number;
  use_start_time?: string | null;
}

export interface PlatformCouponDetail extends PlatformCoupon {
  received_total: number;
  used_total: number;
  use_type: number;
}

export interface PlatformCouponPage {
  limit: number;
  list: PlatformCoupon[];
  page: number;
  total: number;
}

export interface PlatformCouponSaveInput {
  coupon_price: number;
  coupon_time: number;
  coupon_type: number;
  end_time?: string;
  full_reduction?: number;
  is_limited: number;
  is_timeout: number;
  send_type: number;
  sort: number;
  start_time?: string;
  status?: number;
  title: string;
  total_count: number;
  use_end_time?: string;
  use_min_price: number;
  use_start_time?: string;
  use_type: number;
}

export function listPlatformCouponsApi(params: {
  limit: number;
  page: number;
  keyword?: string;
  status?: number;
  send_type?: number;
}) {
  return requestClient.get<PlatformCouponPage>('/coupons', { params });
}

export function getPlatformCouponDetailApi(id: number) {
  return requestClient.get<PlatformCouponDetail>(`/coupons/${id}`);
}

export function createPlatformCouponApi(body: PlatformCouponSaveInput) {
  return requestClient.post<PlatformCoupon>('/coupons', body);
}

export function updatePlatformCouponApi(id: number, body: PlatformCouponSaveInput) {
  return requestClient.put<PlatformCoupon>(`/coupons/${id}`, body);
}

export function setPlatformCouponStatusApi(id: number, status: number) {
  return requestClient.post(`/coupons/${id}/status`, { status });
}

export function deletePlatformCouponApi(id: number) {
  return requestClient.delete(`/coupons/${id}`);
}

export function clonePlatformCouponApi(id: number) {
  return requestClient.post<PlatformCoupon>(`/coupons/${id}/clone`);
}

/** 平台「商户优惠券」列表行 */
export interface StoreCouponListItem extends PlatformCoupon {
  mer_id: number;
  mer_name: string;
  is_trader: number;
  trader_name: string;
  coupon_type_name: string;
  claim_text: string;
  validity_text: string;
  received_total: number;
  used_total: number;
}

export interface StoreCouponDetail extends PlatformCouponDetail {
  mer_id: number;
  mer_name?: string;
  is_trader?: number;
  trader_name?: string;
  coupon_type_name?: string;
  claim_text?: string;
  validity_text?: string;
}

export function listStoreCouponsApi(params: {
  page: number;
  limit: number;
  keyword?: string;
  status?: number;
  is_trader?: number;
}) {
  return requestClient.get<PlatformCouponPage & { list: StoreCouponListItem[] }>(
    '/coupons/store',
    { params },
  );
}

export function getStoreCouponDetailApi(id: number) {
  return requestClient.get<StoreCouponDetail>(`/coupons/store/${id}`);
}
