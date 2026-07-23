import type { InternalAxiosRequestConfig } from 'axios';

import qs from 'qs';

import { getLegacyUserInfo } from './qixi-live-token';

/** 登录前匿名接口：不带 app_id，由 api-platform 凭账号识别租户 */
const ANONYMOUS_SHOP_PATHS = [
  '/shop/passport/login',
  '/shop/index/base',
];

function isAnonymousShopRequest(config: InternalAxiosRequestConfig) {
  const url = String(config?.url || '').split('?')[0];
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
  const raw = parts[1];
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

  const data = payload.data as Record<string, unknown> | undefined;

  const candidates = [
    data?.AppID,
    data?.AppId,
    data?.app_id,
    data?.appId,
    payload.AppID,
    payload.app_id,
    payload.appId,
  ];
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

  const candidates = [userInfo?.AppID, userInfo?.app_id];
  for (const value of candidates) {
    if (value !== undefined && value !== null && String(value).trim() !== '') {
      return String(value).trim();
    }
  }
  const stored = getLegacyUserInfo();
  if (stored) {
    const fromStorage = stored.AppID ?? stored.app_id;
    if (
      fromStorage !== undefined &&
      fromStorage !== null &&
      String(fromStorage).trim() !== ''
    ) {
      return String(fromStorage).trim();
    }
  }
  const fromEnv = import.meta.env.VITE_APP_ID;
  if (
    fromEnv !== undefined &&
    fromEnv !== null &&
    String(fromEnv).trim() !== ''
  ) {
    return String(fromEnv).trim();
  }

  return '';
}

function normalizePostBody(data: unknown) {
  if (typeof data === 'string') {
    return qs.parse(data);
  }
  if (data && typeof data === 'object' && !(data instanceof FormData)) {
    return data as Record<string, unknown>;
  }
  return {};
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
  config.headers.AppID = appId;
  config.headers['App-Id'] = appId;

  if (config.method === 'get') {
    config.params = Object.assign({ app_id: appId }, config.params || {});
  } else if (!(config.data instanceof FormData)) {
    const wasString = typeof config.data === 'string';
    const base = normalizePostBody(config.data);
    if (!base.app_id) {
      base.app_id = appId;
      config.data = wasString ? qs.stringify(base) : base;
    }
  }
  return config;
}
