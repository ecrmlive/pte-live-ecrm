import type {
  ConfigClassificationItem,
  ConfigClassificationItemInput,
  ConfigClassificationPage,
} from './platform-config-classification';

import { requestClient } from '#/api/request';

export interface PlatformConfigSetting extends ConfigClassificationItem {
  classification_name: string;
}

export interface PlatformConfigSettingInput extends ConfigClassificationItemInput {
  classification_id: number;
}

const BASE_PATH = '/maintain/config-settings';

export function listPlatformConfigSettingsApi(params: {
  backend_type?: 0 | 1;
  keyword?: string;
  limit: number;
  page: number;
}) {
  return requestClient.get<ConfigClassificationPage<PlatformConfigSetting>>(BASE_PATH, { params });
}

export function createPlatformConfigSettingApi(body: PlatformConfigSettingInput) {
  return requestClient.post<PlatformConfigSetting>(BASE_PATH, body);
}

export function updatePlatformConfigSettingApi(id: number, body: PlatformConfigSettingInput) {
  return requestClient.put<PlatformConfigSetting>(`${BASE_PATH}/${id}`, body);
}

export function setPlatformConfigSettingStatusApi(id: number, status: 0 | 1) {
  return requestClient.put<{ id: number; status: number }>(`${BASE_PATH}/${id}/status`, { status });
}

export function deletePlatformConfigSettingApi(id: number) {
  return requestClient.delete<{ ok: boolean }>(`${BASE_PATH}/${id}`);
}
