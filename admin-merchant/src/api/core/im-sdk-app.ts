import { requestClient } from '#/api/request';

export type IMSDKAppStatus = 'disabled' | 'enabled';

export interface MerchantIMSDKApp {
  id: number;
  merchant_id: number;
  sdk_app_id: string;
  name: string;
  status: IMSDKAppStatus;
  is_active: boolean;
  api_public_url: string;
  ws_public_url: string;
  pte_profile_id: string;
}

export type MerchantIMSDKAppInput = Pick<
  MerchantIMSDKApp,
  | 'api_public_url'
  | 'name'
  | 'pte_profile_id'
  | 'sdk_app_id'
  | 'status'
  | 'ws_public_url'
>;

export function getMerchantIMSDKAppsApi() {
  return requestClient.get<{ list: MerchantIMSDKApp[] }>('/settings/im-sdk-apps');
}

export function createMerchantIMSDKAppApi(data: MerchantIMSDKAppInput) {
  return requestClient.post<MerchantIMSDKApp>('/settings/im-sdk-apps', data);
}

export function updateMerchantIMSDKAppApi(
  id: number,
  data: MerchantIMSDKAppInput,
) {
  return requestClient.put<{ id: number }>(`/settings/im-sdk-apps/${id}`, data);
}

export function activateMerchantIMSDKAppApi(id: number) {
  return requestClient.post<{ id: number; is_active: boolean }>(
    `/settings/im-sdk-apps/${id}/activate`,
  );
}

export function disableMerchantIMSDKAppApi(id: number) {
  return requestClient.post<{ id: number; is_active: boolean }>(
    `/settings/im-sdk-apps/${id}/disable`,
  );
}
