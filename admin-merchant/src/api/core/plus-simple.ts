import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export async function getCollectionSettingApi() {
  return requestClient.get<{ app_id?: number; vars: { values: { status: number | string } } }>(
    '/shop/plus.collection/index',
  );
}

export async function saveCollectionSettingApi(payload: { status: number }) {
  return requestClient.post('/shop/plus.collection/index', payload);
}

export async function getOfficiaSettingApi() {
  return requestClient.get<{ vars: { values: { status: number | string } } }>(
    '/shop/plus.officia/index',
  );
}

export async function saveOfficiaSettingApi(payload: { status: number }) {
  return requestClient.post('/shop/plus.officia/index', payload);
}

export async function getFullfreeSettingApi() {
  return requestClient.get<{ vars: { values: { is_open: string; money: string } } }>(
    '/shop/plus.fullfree/index',
  );
}

export async function saveFullfreeSettingApi(payload: { is_open: string; money: string }) {
  return requestClient.post('/shop/plus.fullfree/index', payload);
}

export interface RecommendProductItem {
  product_id: number;
  product_image: string;
  product_name: string;
  sort: number;
}

export interface RecommendFormValues {
  choice: number;
  is_recommend: number;
  location: number[];
  name: string;
  num?: number;
  product?: RecommendProductItem[];
  type?: number;
}

export async function getRecommendSettingApi() {
  return requestClient.get<{
    product_arr?: number[];
    vars: { values: RecommendFormValues };
  }>('/shop/plus.recommend/index');
}

export async function saveRecommendSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.recommend/index', payload);
}

export interface SignListItem {
  continuousDays: number;
  coupon?: Array<{ coupon_num: number; name: string }>;
  endDate: string;
  lastDate: string;
  minDate: string;
  nickName: string;
  points: number;
  totalDay: number;
  user_id: number;
}

export interface SignCouponItem {
  coupon_id: number;
  coupon_num: number;
  name: string;
}

export interface SignRewardRule {
  coupon: SignCouponItem[];
  day: number | string;
  is_coupon: boolean;
  is_point: boolean;
  point: number | string;
}

export interface SignFormValues {
  content?: string;
  coupon?: SignCouponItem[];
  ever_sign?: number | string;
  increase_reward?: number | string;
  is_coupon?: boolean | number | string;
  is_increase?: boolean | number | string;
  is_open?: boolean | number | string;
  no_increase?: number;
  reward_data?: SignRewardRule[];
  sign_type?: number | string;
}

export async function getSignSettingApi() {
  return requestClient.get<{ vars: { values: SignFormValues } }>(
    '/shop/plus.sign/index',
  );
}

export async function saveSignSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.sign/index', payload);
}

export async function getSignRecordListApi(params: Record<string, unknown>) {
  return requestClient.post<{ list: PaginatedList<SignListItem> }>(
    '/shop/plus.sign/lists',
    params,
  );
}

export interface RegisterCouponItem {
  coupon_id: number;
  coupon_num: number;
  name: string;
}

export interface RegisterFormValues {
  balance?: number | string;
  coupon?: RegisterCouponItem[] | string;
  is_balance?: boolean | number | string;
  is_coupon?: boolean | number | string;
  is_open?: boolean | number | string;
  is_point?: boolean | number | string;
  login_push?: number | string;
  points?: number | string;
  push_type?: number | string;
  time?: string[];
}

export async function getRegisterSettingApi() {
  return requestClient.get<{ vars: { values: RegisterFormValues } }>(
    '/shop/plus.register/index',
  );
}

export async function saveRegisterSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.register/index', payload);
}

export interface TaskCenterItem {
  is_open?: number | string;
  name?: string;
  points?: number | string;
  rule?: string;
}

export interface TaskCenterFormValues {
  back_image?: string;
  day_task?: TaskCenterItem[];
  grow_task?: TaskCenterItem[];
}

export async function getTaskCenterSettingApi() {
  return requestClient.get<{ vars: { values: TaskCenterFormValues } }>(
    '/shop/plus.task/index',
  );
}

export async function saveTaskCenterSettingApi(payload: TaskCenterFormValues) {
  return requestClient.post('/shop/plus.task/index', payload);
}
