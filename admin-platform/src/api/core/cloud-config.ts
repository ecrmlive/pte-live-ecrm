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
