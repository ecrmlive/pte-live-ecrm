import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface CouponListItem {
  coupon_id: number;
  coupon_type?: { text?: string; value?: number };
  min_price: string | number;
  name: string;
  total_num: number;
}

export async function getCouponListApi(params: {
  coupon_type?: number;
  create_time?: string;
  list_rows?: number;
  name?: string;
  page?: number;
  show_center?: number;
} = {}) {
  return requestClient.post<{ list: PaginatedList<CouponListItem> }>(
    '/shop/plus.coupon.coupon/index',
    {
      show_center: -1,
      ...params,
    },
  );
}

export async function deleteCouponApi(couponId: number) {
  return requestClient.post('/shop/plus.coupon.coupon/delete', { coupon_id: couponId });
}

export async function getCouponSettingApi() {
  return requestClient.get<{ vars: { values: Record<string, unknown> } }>(
    '/shop/plus.coupon.setting/index',
  );
}

export async function saveCouponSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.coupon.setting/index', payload);
}

export async function getCouponReceiveListApi(params: Record<string, unknown>) {
  return requestClient.post<{ list: PaginatedList<Record<string, unknown>> }>(
    '/shop/plus.coupon.coupon/receive',
    params,
  );
}

export interface SeckillActiveItem {
  active_id: number;
  active_name: string;
  end_time: string;
  start_time: string;
  status?: { text?: string; value?: number };
}

export async function getSeckillActiveListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<SeckillActiveItem> }>(
    '/shop/plus.seckill.Active/index',
    params,
  );
}

export async function deleteSeckillActiveApi(activeId: number) {
  return requestClient.post('/shop/plus.seckill.Active/delete', { active_id: activeId });
}

export async function getSeckillSettingApi() {
  return requestClient.get<{ vars: { values: Record<string, unknown> } }>(
    '/shop/plus.seckill.Setting/index',
  );
}

export async function saveSeckillSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.seckill.Setting/index', payload);
}

export interface PointsProductItem {
  limit_num?: number;
  point_product_id: number;
  product?: {
    image?: Array<{ file_path: string }>;
    product_name?: string;
    spec_type?: number;
  };
  sku?: Array<{ point_money?: number | string; point_num?: number }>;
  sort?: number;
  status?: number;
  stock?: number;
}

export async function getPointsProductListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{
    exclude_ids?: number[];
    list: PaginatedList<PointsProductItem>;
  }>('/shop/plus.points.product/index', params);
}

export async function deletePointsProductApi(pointProductId: number) {
  return requestClient.post('/shop/plus.points.product/delete', { id: pointProductId });
}

export async function getPointsSettingApi() {
  return requestClient.get<{ vars: { values: Record<string, unknown> } }>(
    '/shop/plus.points/setting',
  );
}

export async function savePointsSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.points/setting', payload);
}

export async function getPointsExchangeSettingsApi() {
  return requestClient.get<{ vars: { values: Record<string, unknown> } }>(
    '/shop/plus.points.product/settings',
  );
}

export async function savePointsExchangeSettingsApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.points.product/settings', payload);
}

export interface PointsExchangeRecordItem {
  create_time?: string;
  order_id?: number;
  pay_price?: string | number;
  points_num?: number;
  state_text?: string;
  user?: { avatarUrl?: string; nickName?: string };
}

export async function getPointsExchangeRecordListApi(params: {
  list_rows?: number;
  nickName?: string;
  order_status?: number | string;
  page?: number;
}) {
  return requestClient.post<{
    list: PaginatedList<PointsExchangeRecordItem>;
    status?: Array<{ id: number; name: string }>;
  }>('/shop/plus.points.product/record', params);
}

export interface CouponFormPayload extends Record<string, unknown> {
  active_time?: string[];
  apply_range?: number;
  category_list?: Record<string, unknown>;
  content?: string;
  coupon_id?: number;
  coupon_type?: number | string;
  discount?: number | string;
  expire_day?: number | string;
  expire_type?: number | string;
  free_limit?: number;
  max_price?: number | string;
  min_price?: number | string;
  name?: string;
  product_ids?: number[];
  reduce_price?: number | string;
  show_center?: number;
  sort?: number;
  total_num?: number | string;
}

export async function getCouponEditMetaApi(couponId: number) {
  return requestClient.get<{ detail: Record<string, unknown> }>(
    '/shop/plus.coupon.coupon/edit',
    { params: { coupon_id: couponId } },
  );
}

export async function addCouponApi(payload: CouponFormPayload) {
  return requestClient.post('/shop/plus.coupon.coupon/add', payload);
}

export async function saveCouponApi(payload: CouponFormPayload) {
  return requestClient.post('/shop/plus.coupon.coupon/edit', payload);
}

export interface CouponSendGradeItem {
  grade_id: number;
  name: string;
}

export async function getCouponSendMetaApi() {
  return requestClient.get<{ list: CouponSendGradeItem[] }>(
    '/shop/plus.coupon.coupon/SendCoupon',
  );
}

export async function sendCouponApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.coupon.coupon/SendCoupon', payload);
}
