import { http } from "@/utils/request";

export interface Coupon {
  coupon_id: number;
  mer_id: number;
  title: string;
  coupon_price: number;
  discount_type: "amount" | "rate";
  discount_value: number;
  use_min_price: number;
  coupon_time: number;
  remain_count: number;
  is_limited: number;
  type: number;
  received?: boolean;
}

export interface CouponUser {
  coupon_user_id: number;
  coupon_id: number;
  mer_id: number;
  coupon_title: string;
  coupon_price: number;
  discount_type: "amount" | "rate";
  discount_value: number;
  use_min_price: number;
  status: number;
  coupon_kind?: number;
  starts_at?: string | null;
  ends_at?: string | null;
}

export function fetchCouponCenter(mer_id?: number) {
  const q = mer_id ? `?mer_id=${mer_id}` : "";
  return http.get<{ list: Coupon[] }>(`/coupons/center${q}`);
}

export function receiveCoupon(id: number) {
  return http.post<CouponUser>(`/coupons/${id}/receive`, {});
}

export function fetchMyCoupons(status?: "unused" | "history", page = 1) {
  const params = new URLSearchParams({ page: String(page), limit: "20" });
  if (status !== undefined) params.set("status", String(status));
  return http.get<{ list: CouponUser[]; total: number }>(`/coupons/mine?${params}`);
}

export function fetchUsableCoupons(cartIds: number[]) {
  return http.get<{ list: CouponUser[]; total_price: number }>(
    `/coupons/usable?cart_ids=${cartIds.join(",")}`,
  );
}

export function bindSpread(spread_uid: number) {
  return http.post<{ ok: boolean }>("/spread/bind", { spread_uid });
}
