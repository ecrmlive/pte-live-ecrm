import { requestClient } from '#/api/request';

export interface MerchantCoupon {
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

export interface MerchantCouponPage {
  limit: number;
  list: MerchantCoupon[];
  page: number;
  total: number;
}

export interface MerchantCouponSaveInput {
  coupon_price: number;
  coupon_time: number;
  is_limited: number;
  sort: number;
  status?: number;
  title: string;
  total_count: number;
  use_min_price: number;
}

export interface MerchantCouponRecord {
  coupon_id: number;
  coupon_price: number;
  coupon_title: string;
  coupon_user_id: number;
  create_time: string;
  end_time?: string;
  mer_id: number;
  send_id: number;
  start_time?: string;
  status: number;
  type: string;
  uid: number;
  use_min_price: number;
  use_time?: string;
}

export interface MerchantCouponSend {
  coupon_id: number;
  coupon_num: number;
  coupon_send_id: number;
  create_time: string;
  mark: string;
  status: number;
}

export interface MerchantCouponPageResult<T> {
  limit: number;
  list: T[];
  page: number;
  total: number;
}

export interface MerchantCouponListParams {
  date_from?: string;
  date_to?: string;
  keyword?: string;
  limit: number;
  page: number;
  status?: 0 | 1;
}

export function listMerchantCouponsApi(params: MerchantCouponListParams) {
  return requestClient.get<MerchantCouponPage>('/coupons', { params });
}

export function createMerchantCouponApi(body: MerchantCouponSaveInput) {
  return requestClient.post<MerchantCoupon>('/coupons', body);
}

export function updateMerchantCouponApi(id: number, body: MerchantCouponSaveInput) {
  return requestClient.put<MerchantCoupon>(`/coupons/${id}`, body);
}

export function setMerchantCouponStatusApi(id: number, status: number) {
  return requestClient.post(`/coupons/${id}/status`, { status });
}

export function deleteMerchantCouponApi(id: number) {
  return requestClient.delete(`/coupons/${id}`);
}

export function sendMerchantCouponApi(id: number, body: { mark: string; uids: number[] }) {
  return requestClient.post<MerchantCouponSend>(`/coupons/${id}/send`, body);
}

export interface MerchantCouponRecordListParams {
  coupon_id?: number;
  date_from?: string;
  date_to?: string;
  keyword?: string;
  limit: number;
  page: number;
  status?: 0 | 1 | 2;
  uid?: number;
}

export interface MerchantCouponSendListParams {
  coupon_id?: number;
  date_from?: string;
  date_to?: string;
  limit: number;
  page: number;
  status?: 0 | 1;
}

export function listMerchantCouponRecordsApi(params: MerchantCouponRecordListParams) {
  return requestClient.get<MerchantCouponPageResult<MerchantCouponRecord>>('/coupons/records', { params });
}

export function listMerchantCouponSendsApi(params: MerchantCouponSendListParams) {
  return requestClient.get<MerchantCouponPageResult<MerchantCouponSend>>('/coupons/sends', { params });
}
