import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface AppWxSetting {
  wxapp_id: string;
  wxapp_secret: string;
}

export interface AppMpSetting {
  mpapp_id: string;
  mpapp_secret: string;
}

export interface PayTypeOption {
  name: string;
  value: number;
}

export interface PlatformPayConfig {
  name: string;
  pay_type: number[];
  value: string;
}

export interface AppPayFormPayload {
  alipay_appid?: string;
  alipay_cert_path?: string;
  alipay_privatekey?: string;
  alipay_publickey?: string;
  apikey?: string;
  cert_pem?: string;
  key_pem?: string;
  mchid?: string;
  pay_type: Record<string, PlatformPayConfig>;
  platform_pem?: string;
  serial_no?: string;
  wx_cash_type?: number | string;
}

export interface AppPayDetailResult {
  app: AppPayFormPayload & { app_id?: number; pay_type?: Record<string, PlatformPayConfig> | null };
  pay_type: Record<number, PayTypeOption>;
  platform: Record<string, PlatformPayConfig>;
}

export interface AppOpenSetting {
  logo?: string;
  openapp_id?: string;
  openapp_secret?: string;
}

export interface AppShareSetting {
  bind_type?: number;
  down_url?: string;
  gh_id?: string;
  open_site?: string;
  type?: number;
  web_url?: string;
}

export interface AppH5AlipaySetting {
  app_id?: string;
  is_open?: boolean | number;
  privateKey?: string;
  publicKey?: string;
}

export interface AppUpdateItem {
  create_time?: string;
  pkg_url_android?: string;
  pkg_url_ios?: string;
  update_id: number;
  version_android?: string;
  version_ios?: string;
  wgt_url?: string;
}

export async function getAppWxSettingApi() {
  return requestClient.get<{ data: AppWxSetting | null }>('/shop/appsetting.appwx/index');
}

export async function saveAppWxSettingApi(payload: AppWxSetting) {
  return requestClient.post('/shop/appsetting.appwx/index', payload);
}

export async function getAppMpSettingApi() {
  return requestClient.get<{ data: AppMpSetting | null }>('/shop/appsetting.appmp/index');
}

export async function saveAppMpSettingApi(payload: AppMpSetting) {
  return requestClient.post('/shop/appsetting.appmp/index', payload);
}

export async function getAppPaySettingApi() {
  return requestClient.get<AppPayDetailResult>('/shop/appsetting.app/pay');
}

export async function saveAppPaySettingApi(payload: AppPayFormPayload) {
  return requestClient.post<{ msg?: string }>('/shop/appsetting.app/pay', payload);
}

export async function getAppOpenSettingApi() {
  return requestClient.get<{ data: AppOpenSetting | null }>(
    '/shop/appsetting.appopen/index',
  );
}

export async function saveAppOpenSettingApi(payload: AppOpenSetting) {
  return requestClient.post<{ msg?: string }>('/shop/appsetting.appopen/index', payload);
}

export async function getAppShareSettingApi() {
  return requestClient.get<{ data: AppShareSetting | null }>(
    '/shop/appsetting.appshare/index',
  );
}

export async function saveAppShareSettingApi(payload: AppShareSetting) {
  return requestClient.post<{ msg?: string }>('/shop/appsetting.appshare/index', payload);
}

export async function getAppH5AlipaySettingApi() {
  return requestClient.get<{ data: AppH5AlipaySetting | null }>(
    '/shop/appsetting.apph5/pay',
  );
}

export async function saveAppH5AlipaySettingApi(payload: AppH5AlipaySetting) {
  return requestClient.post<{ msg?: string }>('/shop/appsetting.apph5/pay', payload);
}

export async function getAppUpdateListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<AppUpdateItem> }>(
    '/shop/appsetting.appupdate/index',
    params,
  );
}

export async function addAppUpdateApi(payload: Omit<AppUpdateItem, 'create_time' | 'update_id'>) {
  return requestClient.post<{ msg?: string }>('/shop/appsetting.appupdate/add', payload);
}

export async function editAppUpdateApi(payload: AppUpdateItem) {
  return requestClient.post<{ msg?: string }>('/shop/appsetting.appupdate/edit', payload);
}

export async function deleteAppUpdateApi(updateId: number) {
  return requestClient.post<{ msg?: string }>('/shop/appsetting.appupdate/delete', {
    update_id: updateId,
  });
}
