import axios from 'axios';
import qs from 'qs';

import { resolveApiBaseUrl } from '#/utils/qixi-live-api';
import { resolveShopAppId } from '#/utils/qixi-live-shop-app-id';
import {
  attachAPIEncryption,
  decryptAPIResponse,
} from '../../../admin-platform/src/utils/api-crypto';

export type LoginPlatform =
  | 'platform_admin'
  | 'merchant_admin'
  | 'mp_wechat'
  | 'h5'
  | 'android'
  | 'ios'
  | 'harmonyos';

export interface RefreshTokenOptions {
  loginPlatform: LoginPlatform;
  token: string;
  appId?: number | string;
}

export interface RefreshTokenResult {
  ok: boolean;
  token?: string;
  expiresAt?: number;
  clearToken?: boolean;
}

function isRefreshUnauthorized(status: number, code?: number) {
  return status === 401 || code === -1;
}

export async function refreshLiveJwtToken(
  options: RefreshTokenOptions,
): Promise<RefreshTokenResult> {
  const token = options.token?.trim();
  if (!token) {
    return { ok: false, clearToken: true };
  }

  const appId = options.appId ?? resolveShopAppId();
  const headers: Record<string, string> = {
    'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
    'Authori-zation': `Bearer ${token}`,
    'X-Login-Platform': options.loginPlatform,
  };
  if (appId) {
		headers['X-AppId'] = String(appId);
  }

  try {
    const baseURL = resolveApiBaseUrl();
    const config = await attachAPIEncryption(
      { headers, timeout: 5_000 },
      baseURL,
    );
    const res = await axios.post(
      `${baseURL}/api/v1/auth/refresh`,
		qs.stringify({ login_platform: options.loginPlatform }),
      config,
    );
    await decryptAPIResponse(res);

    const authHeader = res.headers['authori-zation'];
    const headerToken =
      typeof authHeader === 'string' && authHeader.startsWith('Bearer ')
        ? authHeader.slice(7).trim()
        : '';

    const body = res.data as {
      code?: number;
      data?: { token?: string; expires_at?: number };
    };

    if (body.code === 1) {
      const newToken = headerToken || body.data?.token || '';
      if (!newToken) {
        return { ok: false };
      }
      return {
        ok: true,
        token: newToken,
        expiresAt: body.data?.expires_at,
      };
    }

    if (isRefreshUnauthorized(res.status, body.code)) {
      return { ok: false, clearToken: true };
    }

    return { ok: false };
  } catch (error) {
    const err = error as {
      response?: { data?: { code?: number }; status?: number };
    };
    const status = err.response?.status;
    const code = err.response?.data?.code;
    if (isRefreshUnauthorized(status ?? 0, code)) {
      return { ok: false, clearToken: true };
    }
    return { ok: false };
  }
}
