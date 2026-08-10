import { requestClient } from '#/api/request';

export interface PlatformShopConfig {
  site_name: string;
  site_url: string;
  order_auto_cancel_minutes: number;
  order_auto_receive_days: number;
  enabled: boolean;
  remark: string;
}

export interface PlatformPayConfig {
  wechat_enabled: boolean;
  alipay_enabled: boolean;
  balance_enabled: boolean;
  remark: string;
}

export interface PlatformWechatAppConfig {
  app_name: string;
  enabled: boolean;
  remark: string;
}

/** 对齐 CRMEB margin_remind_switch / margin_remind_day */
export interface PlatformMarginConfig {
  margin_remind_switch: boolean;
  margin_remind_day: number;
}

function parseConfig<T>(raw: string): T {
  return JSON.parse(raw) as T;
}

function stringifyConfig<T>(config: T): string {
  return JSON.stringify(config);
}

export function getPlatformShopConfigApi() {
  return requestClient
    .get<{ config: string; note: string }>('/setting/shop')
    .then((data) => ({
      note: data.note,
      config: parseConfig<PlatformShopConfig>(data.config),
    }));
}

export function savePlatformShopConfigApi(config: PlatformShopConfig) {
  return requestClient
    .put<{ config: string }>('/setting/shop', { config: stringifyConfig(config) })
    .then((data) => parseConfig<PlatformShopConfig>(data.config));
}

export function getPlatformPayConfigApi() {
  return requestClient
    .get<{ config: string; note: string }>('/setting/pay')
    .then((data) => ({
      note: data.note,
      config: parseConfig<PlatformPayConfig>(data.config),
    }));
}

export function savePlatformPayConfigApi(config: PlatformPayConfig) {
  return requestClient
    .put<{ config: string }>('/setting/pay', { config: stringifyConfig(config) })
    .then((data) => parseConfig<PlatformPayConfig>(data.config));
}

export function getPlatformWechatAppConfigApi() {
  return requestClient
    .get<{ config: string; note: string }>('/setting/wechat-app')
    .then((data) => ({
      note: data.note,
      config: parseConfig<PlatformWechatAppConfig>(data.config),
    }));
}

export function savePlatformWechatAppConfigApi(config: PlatformWechatAppConfig) {
  return requestClient
    .put<{ config: string }>('/setting/wechat-app', { config: stringifyConfig(config) })
    .then((data) => parseConfig<PlatformWechatAppConfig>(data.config));
}

export function getPlatformMarginConfigApi() {
  return requestClient
    .get<{ config: string; note: string }>('/setting/margin')
    .then((data) => ({
      note: data.note,
      config: parseConfig<PlatformMarginConfig>(data.config),
    }));
}

export function savePlatformMarginConfigApi(config: PlatformMarginConfig) {
  return requestClient
    .put<{ config: string }>('/setting/margin', { config: stringifyConfig(config) })
    .then((data) => parseConfig<PlatformMarginConfig>(data.config));
}

/** 对齐 CRMEB ConfigOthers::update（分销配置 distribution_tabs） */
export interface PlatformDistributionConfig {
  extension_status: boolean;
  extension_self: boolean;
  extension_limit: boolean;
  extension_limit_day: number;
  /** 0礼包 1手动 2人人 3满额 */
  promoter_type: number;
  promoter_low_money: number;
  /** 0全部 1推广员 2非推广员 3关闭 */
  extension_pop: number;
  /** 0～1，例 0.15 = 15% */
  extension_one_rate: number;
  extension_two_rate: number;
  user_extract_min: number;
  lock_brokerage_timer: number;
  /** 0线下 1企业付款到零钱 2商家转账到零钱 */
  sys_extension_type: number;
  /** 0银行卡 1微信 2支付宝 4余额 */
  withdraw_type: string[];
  /** 1线下转账 2自动转账 */
  extract_switch: number;
  transfer_scene_id: number;
  max_bag_number: number;
}

export function getPlatformDistributionConfigApi() {
  return requestClient
    .get<{ config: string; note: string }>('/setting/distribution')
    .then((data) => ({
      note: data.note,
      config: parseConfig<PlatformDistributionConfig>(data.config),
    }));
}

export function savePlatformDistributionConfigApi(
  config: PlatformDistributionConfig,
) {
  return requestClient
    .put<{ config: string }>('/setting/distribution', {
      config: stringifyConfig(config),
    })
    .then((data) => parseConfig<PlatformDistributionConfig>(data.config));
}

/** 对齐 CRMEB 商户设置：入驻页背景 + 自定义表单字段 */
export type MerchantApplyFieldType =
  | 'checkbox'
  | 'city'
  | 'date'
  | 'daterange'
  | 'radio'
  | 'select'
  | 'text'
  | 'textarea'
  | 'time'
  | 'timerange'
  | 'image';

export type MerchantApplyContentType =
  | 'text'
  | 'number'
  | 'mobile'
  | 'idcard'
  | 'email';

export type MerchantApplyCityLevel =
  | 'province_city'
  | 'province_city_district'
  | 'province_city_district_street';

export type MerchantApplyDefaultVisible = 'show' | 'hide';

export type MerchantApplyDefaultMode = 'current' | 'specify';

export interface MerchantApplyFormField {
  id: string;
  type: MerchantApplyFieldType;
  title: string;
  content_type: MerchantApplyContentType;
  default_value: string;
  placeholder: string;
  required: boolean;
  options?: string[];
  /** image：最多上传 */
  max_upload?: number;
  /** city：省市级别 */
  city_level?: MerchantApplyCityLevel;
  /** date/time 类：默认值显示/隐藏 */
  default_visible?: MerchantApplyDefaultVisible;
  /** date/time 类：当前 / 指定 */
  default_mode?: MerchantApplyDefaultMode;
  /** date/time 类：指定值 */
  specify_value?: string;
}

export interface PlatformMerchantApplyConfig {
  background_image: string;
  form_fields: MerchantApplyFormField[];
}

export function getPlatformMerchantApplyConfigApi() {
  return requestClient
    .get<{ config: string; note: string }>('/setting/merchant-apply')
    .then((data) => ({
      note: data.note,
      config: parseConfig<PlatformMerchantApplyConfig>(data.config),
    }));
}

export function savePlatformMerchantApplyConfigApi(
  config: PlatformMerchantApplyConfig,
) {
  return requestClient
    .put<{ config: string }>('/setting/merchant-apply', {
      config: stringifyConfig(config),
    })
    .then((data) => parseConfig<PlatformMerchantApplyConfig>(data.config));
}

/** 对齐 CRMEB circle_config：默认三级提成 + 代理申请表单字段 */
export interface PlatformAgentZoneConfig {
  one_agent_commission: number;
  two_agent_commission: number;
  three_agent_commission: number;
  form_fields: MerchantApplyFormField[];
}

export function getPlatformAgentZoneConfigApi() {
  return requestClient
    .get<{ config: string; note: string }>('/setting/agent-zone')
    .then((data) => ({
      note: data.note,
      config: parseConfig<PlatformAgentZoneConfig>(data.config),
    }));
}

export function savePlatformAgentZoneConfigApi(config: PlatformAgentZoneConfig) {
  return requestClient
    .put<{ config: string }>('/setting/agent-zone', {
      config: stringifyConfig(config),
    })
    .then((data) => parseConfig<PlatformAgentZoneConfig>(data.config));
}
