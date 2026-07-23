import { requestClient } from '#/api/request';

export interface HomepushLink {
  name?: string;
  type?: string;
  url?: string;
  [key: string]: unknown;
}

export interface HomepushCouponItem {
  coupon_id?: number;
  name?: string;
  type?: string;
}

export interface HomepushFormValues {
  coupon?: HomepushCouponItem[];
  des?: string;
  file_path?: string;
  image_id?: number | string;
  is_open?: boolean | string;
  link?: HomepushLink | string;
  name?: string;
  remark?: string;
  title?: string;
  type?: number | string;
}

export interface HomepushGetResult {
  vars: {
    values: HomepushFormValues;
  };
}

export async function getHomepushApi() {
  return requestClient.get<HomepushGetResult>('/shop/plus.homepush/index');
}

export async function saveHomepushApi(data: Record<string, unknown>) {
  return requestClient.post('/shop/plus.homepush/index', data);
}
