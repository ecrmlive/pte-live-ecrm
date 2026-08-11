import { requestClient } from '#/api/request';

export interface CloudConfigField {
  hint?: string;
  input_type?: 'number' | 'select' | 'switch' | 'textarea' | 'url';
  key: string;
  label: string;
  options?: string[];
  required: boolean;
  secret: boolean;
}

export interface CloudConfigGroup {
  configured: boolean;
  fields: CloudConfigField[];
  group_key: string;
  label: string;
  updated_at?: string;
  values: Record<string, string>;
}

export async function fetchCloudConfigs() {
  const response = await requestClient.get<{ list: CloudConfigGroup[] }>(
    '/setting/cloud-configs',
  );
  return response.list;
}

export async function saveCloudConfig(group: string, values: Record<string, string>) {
  return requestClient.put<CloudConfigGroup>(
    `/setting/cloud-configs/${group}`,
    { values },
  );
}

/** 小程序页面专用：AppSecret 仅返回“已配置”掩码，留空保存不会覆盖已有密钥。 */
export function getRoutineConfigApi() {
  return requestClient.get<CloudConfigGroup>('/setting/routine-config');
}

export function saveRoutineConfigApi(values: Record<string, string>) {
  return requestClient.put<CloudConfigGroup>('/setting/routine-config', {
    values,
  });
}

export type NativePlatform = 'android' | 'harmony' | 'ios';

export function getMobileAppConfigApi(platform: NativePlatform) {
  return requestClient.get<CloudConfigGroup>(
    `/setting/mobile-app-config/${platform}`,
  );
}

export function saveMobileAppConfigApi(
  platform: NativePlatform,
  values: Record<string, string>,
) {
  return requestClient.put<CloudConfigGroup>(
    `/setting/mobile-app-config/${platform}`,
    { values },
  );
}

export function getPushConfigApi(platform: NativePlatform) {
  return requestClient.get<CloudConfigGroup>(
    `/setting/push-config/${platform}`,
  );
}

export function savePushConfigApi(
  platform: NativePlatform,
  values: Record<string, string>,
) {
  return requestClient.put<CloudConfigGroup>(
    `/setting/push-config/${platform}`,
    { values },
  );
}

export interface MapClientConfig {
  amap_web_js_key: string;
  amap_web_js_security_code: string;
  configured: boolean;
  provider: string;
}

/** 已鉴权：仅返回 Web JS Key/安全密钥（不下发 Web 服务 Key 与各端 Key）。 */
export async function fetchMapClientConfig() {
  return requestClient.get<MapClientConfig>('/setting/map-client-config');
}
