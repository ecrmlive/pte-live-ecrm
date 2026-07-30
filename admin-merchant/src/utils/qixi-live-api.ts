/** api-platform 基址规则见 docs/architecture.md。 */
export function resolveApiBaseUrl() {
  const raw =
    import.meta.env.VITE_GLOB_API_URL || import.meta.env.VITE_API_BASE_URL;
  if (raw) {
    return String(raw).replace(/\/$/, '');
  }
  return import.meta.env.MODE === 'test'
    ? 'http://127.0.0.1:18083'
    : 'http://127.0.0.1:18083';
}

/** B 端（平台/商户）统一 api-platform；直播运行时亦走同一基址 `/api/v1/shop/live/*`（终态不 HTTP 反代 api-live） */
export function resolveLiveApiBaseUrl() {
  return resolveApiBaseUrl();
}

export const QIXI_SHOP_TOKEN_KEY = 'qixiMergersMerchantAdminToken';
export const QIXI_SHOP_MENU_KEY = 'qixiMergersMerchantAdminMenus';
export const QIXI_SHOP_RENDER_MENU_KEY = 'qixiMergersMerchantAdminRenderMenus';

/** 全屏独立页（不走 BasicLayout、不重建动态菜单） */
export const STANDALONE_ROUTE_NAMES = new Set(['LiveControlCenter']);

/** DIY 装修等需要全宽、无 Vben 侧栏的 path 前缀 */
export const NO_LAYOUT_PATH_PREFIXES = [
  '/page/',
  '/devise/diy/index',
  '/setting/diy/index',
];

export function isNoBasicLayoutPath(path: string) {
  if (path === '/live/control/center') return true;
  return NO_LAYOUT_PATH_PREFIXES.some((prefix) => path.startsWith(prefix));
}

/** 后端 path → legacy 组件 path（大小写/文件名不一致） */
export const LEGACY_COMPONENT_ALIASES: Record<string, string> = {
  '/home': '/home/home',
  '/live/control/center': '/live/control/Center',
  '/live/traffic/session': '/live/traffic/SessionList',
  '/setting/Index': '/setting/index',
};

export function resolveLegacyComponentPath(path: string) {
  return LEGACY_COMPONENT_ALIASES[path] || path;
}
