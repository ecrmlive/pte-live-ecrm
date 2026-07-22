/** 本仓库契约前缀：C 端走 api-app `/api/app/v1`（Nginx 反代到 :18085）。 */

const MP_DEV_BASE = "http://127.0.0.1:18085";

export function apiBaseURL(): string {
  let base = MP_DEV_BASE;
  // #ifdef H5
  base = "";
  // #endif
  return base;
}

export const APP_API_PREFIX = "/api/app/v1";
export const CALLBACK_API_PREFIX = "/api/callback/v1";
