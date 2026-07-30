import { http } from "@/utils/request";

export interface Coupon {
  coupon_id: number;
  mer_id: number;
  title: string;
  coupon_price: number;
  use_min_price: number;
  received?: boolean;
}

export interface CouponUser {
  coupon_user_id: number;
  coupon_id: number;
  mer_id: number;
  coupon_title: string;
  coupon_price: number;
  use_min_price: number;
  status: number;
  coupon_kind?: number;
}

export function fetchCouponCenter(mer_id?: number) {
  const q = mer_id ? `?mer_id=${mer_id}` : "";
  return http.get<{ list: Coupon[] }>(`/coupons/center${q}`);
}

export function receiveCoupon(id: number) {
  return http.post<CouponUser>(`/coupons/${id}/receive`, {});
}

export function fetchMyCoupons(status?: number, page = 1) {
  const params = new URLSearchParams({ page: String(page), limit: "20" });
  if (status !== undefined) params.set("status", String(status));
  return http.get<{ list: CouponUser[]; total: number }>(`/coupons/mine?${params}`);
}

export function fetchUsableCoupons(cartIds: number[]) {
  return http.get<{ list: CouponUser[]; total_price: number }>(
    `/coupons/usable?cart_ids=${cartIds.join(",")}`,
  );
}
