import { http } from '@/api/http';

export interface Coupon {
  coupon_id: number;
  mer_id: number;
  title: string;
  coupon_price: number;
  use_min_price: number;
  coupon_time: number;
  total_count: number;
  remain_count: number;
  is_limited: number;
  status: number;
  type: number;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchCoupons(params: Record<string, unknown>) {
  return http.get<PageResult<Coupon>>('/coupons', { params });
}

export function createCoupon(data: {
  title: string;
  coupon_price: number;
  use_min_price: number;
  coupon_time?: number;
  total_count?: number;
  is_limited?: number;
}) {
  return http.post<Coupon>('/coupons', data);
}

export function setCouponStatus(id: number, status: number) {
  return http.post(`/coupons/${id}/status`, { status });
}

export function deleteCoupon(id: number) {
  return http.delete<{ ok: boolean }>(`/coupons/${id}`);
}

export function deleteCoupon(id: number) {
  return http.delete(`/coupons/${id}`);
}
