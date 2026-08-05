import { requestClient } from '#/api/request';

export interface MerchantIntegralPolicy {
  enabled: boolean;
  max_deduction_bps: number;
  points_per_yuan: number;
  store_id: number;
}

export function getMerchantIntegralPolicyApi() {
  return requestClient.get<MerchantIntegralPolicy>('/setting/integral-policy');
}

export function saveMerchantIntegralPolicyApi(data: Pick<MerchantIntegralPolicy, 'enabled' | 'max_deduction_bps' | 'points_per_yuan'>) {
  return requestClient.put<MerchantIntegralPolicy>('/setting/integral-policy', data);
}
