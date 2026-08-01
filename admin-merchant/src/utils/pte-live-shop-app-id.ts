import type { InternalAxiosRequestConfig } from 'axios';

import { getLegacyUserInfo } from './pte-live-token';

/** 登录前匿名接口无需店铺应用标识。 */
const ANONYMOUS_SHOP_PATHS = [
  '/shop/passport/login',
  '/shop/index/base',
];

function isAnonymousShopRequest(config: InternalAxiosRequestConfig) {
  const url = String(config?.url || '').split('?')[0] || '';
  return ANONYMOUS_SHOP_PATHS.some((p) => url === p || url.endsWith(p));
}

function decodeBase64Url(value: string) {
	let normalized = value.replace(/-/g, '+').replace(/_/g, '/');
	normalized = `${normalized}${'='.repeat((4 - (normalized.length % 4)) % 4)}`;
	if (typeof atob === 'function') {
		return atob(normalized);
	}
	if (typeof Buffer !== 'undefined') {
		return Buffer.from(normalized, 'base64').toString('utf8');
	}
	return '';
}

function decodeTokenPayload(token: string) {
  if (!token) return null;
  const parts = String(token).split('.');
  if (parts.length < 2) return null;
  const raw = parts[1] || '';
  const payloadText = decodeBase64Url(raw);
  if (!payloadText) return null;
  try {
    const payload = JSON.parse(payloadText) as Record<string, unknown>;
    return payload;
  } catch {
    return null;
  }
}

function resolveShopAppIdFromToken(token?: string) {
  const payload = decodeTokenPayload(token || '');
  if (!payload) return '';

  const candidates = [payload.store_app_id];
  for (const value of candidates) {
    if (value !== undefined && value !== null && String(value).trim() !== '') {
      return String(value).trim();
    }
  }
  return '';
}

export function resolveShopAppId(
  userInfo?: Record<string, unknown> | null,
  accessToken?: string,
) {
  const fromToken = resolveShopAppIdFromToken(accessToken);
  if (fromToken) {
    return fromToken;
  }

  const candidates = [userInfo?.store_app_id];
  for (const value of candidates) {
    if (value !== undefined && value !== null && String(value).trim() !== '') {
      return String(value).trim();
    }
  }
  const stored = getLegacyUserInfo();
  if (stored) {
    const fromStorage = stored.store_app_id;
    if (
      fromStorage !== undefined &&
      fromStorage !== null &&
      String(fromStorage).trim() !== ''
    ) {
      return String(fromStorage).trim();
    }
  }
  return '';
}

export function attachShopAppId(
  config: InternalAxiosRequestConfig,
  userInfo?: Record<string, unknown> | null,
  accessToken?: string,
) {
  if (isAnonymousShopRequest(config)) {
    return config;
  }
  const appId = resolveShopAppId(userInfo, accessToken);
  if (!appId) return config;

  config.headers = config.headers || {};
  config.headers['X-AppId'] = appId;
  return config;
}
