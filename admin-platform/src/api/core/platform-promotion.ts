import { requestClient } from '#/api/request';

export interface PlatformCoupon {
  coupon_id: number;
  coupon_price: number;
  coupon_time: number;
  create_time: string;
  is_limited: number;
  remain_count: number;
  sort: number;
  status: number;
  title: string;
  total_count: number;
  use_min_price: number;
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
  is_limited: number;
  sort: number;
  status?: number;
  title: string;
  total_count: number;
  use_min_price: number;
}

export function listPlatformCouponsApi(params: {
  limit: number;
  page: number;
  keyword?: string;
  status?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<PlatformCouponPage>('/coupons', { params });
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
