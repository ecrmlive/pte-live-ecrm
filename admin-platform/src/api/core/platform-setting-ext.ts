import { requestClient } from '#/api/request';

export interface PlatformStorageConfig {
  provider: string;
  region: string;
  bucket_name: string;
  enabled: boolean;
  remark: string;
}

export interface PlatformUserSetupConfig {
  register_enabled: boolean;
  mobile_required: boolean;
  invite_required: boolean;
  remark: string;
}

export interface PlatformTransferSettingsConfig {
  enabled: boolean;
  min_amount: number;
  remark: string;
}

function parseConfig<T>(raw: string): T {
  return JSON.parse(raw) as T;
}

function stringifyConfig<T>(config: T): string {
  return JSON.stringify(config);
}

export function getPlatformStorageConfigApi() {
  return requestClient
    .get<{ config: string; note: string }>('/setting/storage')
    .then((data) => ({ note: data.note, config: parseConfig<PlatformStorageConfig>(data.config) }));
}

export function savePlatformStorageConfigApi(config: PlatformStorageConfig) {
  return requestClient
    .put<{ config: string }>('/setting/storage', { config: stringifyConfig(config) })
    .then((data) => parseConfig<PlatformStorageConfig>(data.config));
}

export function getPlatformUserSetupConfigApi() {
  return requestClient
    .get<{ config: string; note: string }>('/setting/user-setup')
    .then((data) => ({ note: data.note, config: parseConfig<PlatformUserSetupConfig>(data.config) }));
}

export function savePlatformUserSetupConfigApi(config: PlatformUserSetupConfig) {
  return requestClient
    .put<{ config: string }>('/setting/user-setup', { config: stringifyConfig(config) })
    .then((data) => parseConfig<PlatformUserSetupConfig>(data.config));
}

export function getPlatformTransferSettingsConfigApi() {
  return requestClient
    .get<{ config: string; note: string }>('/setting/transfer-settings')
    .then((data) => ({ note: data.note, config: parseConfig<PlatformTransferSettingsConfig>(data.config) }));
}

export function savePlatformTransferSettingsConfigApi(config: PlatformTransferSettingsConfig) {
  return requestClient
    .put<{ config: string }>('/setting/transfer-settings', { config: stringifyConfig(config) })
    .then((data) => parseConfig<PlatformTransferSettingsConfig>(data.config));
}
