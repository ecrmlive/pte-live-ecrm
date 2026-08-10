import { requestClient } from '#/api/request';

export interface PlatformStorageConfig {
  provider: string;
  region: string;
  bucket_name: string;
  enabled: boolean;
  remark: string;
}

/** 对齐 CRMEB UserInfo 字段类型 */
export type PlatformUserSetupFieldType =
  | 'input'
  | 'int'
  | 'phone'
  | 'date'
  | 'radio'
  | 'address'
  | 'id_card'
  | 'email';

export interface PlatformUserSetupField {
  id: number;
  field: string;
  title: string;
  is_used: number;
  is_require: number;
  is_show: number;
  type: PlatformUserSetupFieldType;
  msg: string;
  content?: string[];
  is_default: number;
  sort: number;
}

export interface PlatformUserSetupCoupon {
  coupon_id: number;
  title: string;
  coupon_type: number;
  coupon_price: number;
  use_min_price: number;
  coupon_time: number;
  is_timeout: number;
  use_start_time?: string;
  use_end_time?: string;
}

/** 对齐 CRMEB 用户设置（基础信息 + 登录注册 / 注册有礼） */
export interface PlatformUserSetupConfig {
  user_default_avatar: string;
  fields: PlatformUserSetupField[];
  is_phone_login: number;
  first_avatar_switch: number;
  open_update_info: number;
  wechat_phone_switch: number;
  newcomer_status: number;
  register_popup_pic: string;
  register_money_status: number;
  register_give_money: number;
  register_integral_status: number;
  register_give_integral: number;
  register_coupon_status: number;
  register_give_coupon: PlatformUserSetupCoupon[];
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
    .then((data) => ({
      note: data.note,
      config: parseConfig<PlatformStorageConfig>(data.config),
    }));
}

export function savePlatformStorageConfigApi(config: PlatformStorageConfig) {
  return requestClient
    .put<{ config: string }>('/setting/storage', {
      config: stringifyConfig(config),
    })
    .then((data) => parseConfig<PlatformStorageConfig>(data.config));
}

export function getPlatformUserSetupConfigApi() {
  return requestClient
    .get<{ config: string; note: string }>('/setting/user-setup')
    .then((data) => ({
      note: data.note,
      config: parseConfig<PlatformUserSetupConfig>(data.config),
    }));
}

export function savePlatformUserSetupConfigApi(config: PlatformUserSetupConfig) {
  return requestClient
    .put<{ config: string }>('/setting/user-setup', {
      config: stringifyConfig(config),
    })
    .then((data) => parseConfig<PlatformUserSetupConfig>(data.config));
}

export function getPlatformTransferSettingsConfigApi() {
  return requestClient
    .get<{ config: string; note: string }>('/setting/transfer-settings')
    .then((data) => ({
      note: data.note,
      config: parseConfig<PlatformTransferSettingsConfig>(data.config),
    }));
}

export function savePlatformTransferSettingsConfigApi(
  config: PlatformTransferSettingsConfig,
) {
  return requestClient
    .put<{ config: string }>('/setting/transfer-settings', {
      config: stringifyConfig(config),
    })
    .then((data) => parseConfig<PlatformTransferSettingsConfig>(data.config));
}
