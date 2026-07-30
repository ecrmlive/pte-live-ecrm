import { requestClient } from '#/api/request';

export interface MerchantShopProfile {
  category_id: number;
  create_time: string;
  mer_address: string;
  mer_id: number;
  mer_info: string;
  mer_name: string;
  mer_phone: string;
  mer_state: number;
  real_name: string;
  status: number;
}

export type MerchantShopProfileInput = Pick<
  MerchantShopProfile,
  'mer_address' | 'mer_info' | 'mer_name' | 'mer_phone' | 'real_name'
>;

export function getMerchantShopProfileApi() {
  return requestClient.get<MerchantShopProfile>('/setting/shop');
}

export function updateMerchantShopProfileApi(data: MerchantShopProfileInput) {
  return requestClient.put<MerchantShopProfile>('/setting/shop', data);
}
