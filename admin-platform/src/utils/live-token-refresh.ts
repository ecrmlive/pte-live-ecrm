import axios from 'axios';
import qs from 'qs';

import { resolveLiveApiBaseUrl } from '#/utils/pte-live-api';
import { attachAPIEncryption, decryptAPIResponse } from '#/utils/api-crypto';

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
  appId?: number;
}

export interface RefreshTokenResult {
  ok: boolean;
  token?: string;
  expiresAt?: number;
  clearToken?: boolean;
}

/** 向 api-live 刷新 JWT，延长有效期；账号无效时返回 clearToken */
export async function refreshLiveJwtToken(
  options: RefreshTokenOptions,
): Promise<RefreshTokenResult> {
  const token = options.token?.trim();
  if (!token) {
    return { ok: false };
  }

  const headers: Record<string, string> = {
    'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
    'authori-zation': `Bearer ${token}`,
    'X-Login-Platform': options.loginPlatform,
  };
  if (options.appId && options.appId > 0) {
    headers.AppID = String(options.appId);
  }

  try {
    const baseURL = resolveLiveApiBaseUrl();
    const config = await attachAPIEncryption({
      headers,
      timeout: 5_000,
      validateStatus: () => true,
    });
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
      msg?: string;
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

    if (body.code === -1) {
      return { ok: false, clearToken: true };
    }

    return { ok: false };
  } catch {
    return { ok: false };
  }
}

export const LOGIN_PLATFORM_OPTIONS = [
  { value: 'platform_admin', label: '平台系统' },
  { value: 'merchant_admin', label: '商户系统' },
  { value: 'mp_wechat', label: '微信小程序' },
  { value: 'h5', label: 'H5' },
  { value: 'android', label: 'Android' },
  { value: 'ios', label: 'iOS' },
  { value: 'harmonyos', label: 'HarmonyOS' },
] as const;

export const USER_ROLE_OPTIONS = [
  { value: 'admin', label: '管理员' },
  { value: 'anchor', label: '主播' },
  { value: 'user', label: '用户' },
] as const;

export function formatUnixTime(ts?: number | null) {
  if (!ts) return '—';
  const d = new Date(ts * 1000);
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
}
