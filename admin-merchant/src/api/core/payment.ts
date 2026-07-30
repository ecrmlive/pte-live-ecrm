import { requestClient } from '#/api/request';

export type MerchantPaymentChannelCode = 'alipay' | 'wechat';

export interface MerchantPaymentChannel {
  channel: MerchantPaymentChannelCode;
  enabled: boolean;
	configured: boolean;
}

export function getMerchantPaymentChannelsApi() {
  return requestClient.get<{ list: MerchantPaymentChannel[] }>('/payment/channels');
}

export function updateMerchantPaymentChannelApi(
  channel: MerchantPaymentChannelCode,
  values: Record<string, string>,
) {
  return requestClient.put<MerchantPaymentChannel>(`/payment/channels/${channel}`, { values });
}
