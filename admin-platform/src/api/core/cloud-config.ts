import { requestClient } from '#/api/request';

export interface CloudConfigField {
  hint?: string;
  key: string;
  label: string;
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
