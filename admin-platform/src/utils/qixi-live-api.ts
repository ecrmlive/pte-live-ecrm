/** api-platform 基址规则见 docs/architecture.md。 */
export function resolveApiBaseUrl() {
  const raw = import.meta.env.VITE_GLOB_API_URL || import.meta.env.VITE_API_BASE_URL;
  if (raw) {
    return String(raw).replace(/\/$/, '');
  }
  return import.meta.env.MODE === 'test'
    ? 'http://127.0.0.1:18080'
    : 'http://127.0.0.1:18080';
}

/** B 端统一 api-platform（平台 Go 接口 + 直播运行时反代） */
export function resolveLiveApiBaseUrl() {
  return resolveApiBaseUrl();
}

/** COS CDN 根域。 */
export function resolveCosBaseUrl() {
  const raw =
    import.meta.env.VITE_COS_BASE_URL || import.meta.env.VITE_GLOB_COS_URL;
  if (raw) {
    return String(raw).replace(/\/$/, '');
  }
  return '';
}

/** 平台默认 Logo 对象键（init_file.sql / static_asset_url） */
export const QIXI_PLATFORM_LOGO_OBJECT_KEY = 'qixi-live/image/default/logo.png';

export function resolvePlatformLogoUrl(
  objectKey: string = QIXI_PLATFORM_LOGO_OBJECT_KEY,
) {
  const key = String(objectKey).replace(/^\//, '');
  return `${resolveCosBaseUrl()}/${key}`;
}

export const QIXI_PLATFORM_APP_ID = 10000;
export const QIXI_ADMIN_TOKEN_KEY = 'qixiMergersPlatformAdminToken';
