import axios from 'axios';

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
  refreshToken: string;
}

export interface RefreshTokenResult {
  ok: boolean;
  token?: string;
  refreshToken?: string;
  clearToken?: boolean;
}

/** 统一后台只使用 refresh token 调用自己的刷新接口，不能以 access token 代替。 */
export async function refreshPlatformJwtToken(
  options: RefreshTokenOptions,
): Promise<RefreshTokenResult> {
  const refreshToken = options.refreshToken?.trim();
  if (!refreshToken) {
    return { ok: false };
  }

  const headers: Record<string, string> = {
    'Content-Type': 'application/json;charset=UTF-8',
    'Authori-zation': `Bearer ${refreshToken}`,
  };

  try {
    const baseURL = resolveLiveApiBaseUrl();
    const config = await attachAPIEncryption({
      headers,
      timeout: 5_000,
      validateStatus: () => true,
    });
    const res = await axios.post(
      `${baseURL}/api/platform/v1/auth/refresh`,
      {},
      config,
    );
    await decryptAPIResponse(res);

    const body = res.data as {
      status?: number;
      data?: { token?: { access_token?: string; refresh_token?: string } };
    };

    if (res.status === 401 || body.status === 401) {
      return { ok: false, clearToken: true };
    }
    if (res.status >= 200 && res.status < 300 && body.status === 200) {
      const nextAccessToken = body.data?.token?.access_token || '';
      const nextRefreshToken = body.data?.token?.refresh_token || '';
      if (!nextAccessToken || !nextRefreshToken) {
        return { ok: false };
      }
      return {
        ok: true,
        token: nextAccessToken,
        refreshToken: nextRefreshToken,
      };
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
