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
