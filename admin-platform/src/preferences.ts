import { defineOverridesPreferences } from '@vben/preferences';

const platformBrand = 'qixi';
const platformLogoUrl = `${import.meta.env.BASE_URL}${platformBrand}/logo.png`;

/** 平台管理端展示名称（不受偏好设置缓存覆盖） */
export const QIXI_PLATFORM_APP_NAME = '七禧多商户·平台';

/** 七禧 平台超管 — Vben 偏好（仅 zh-CN） */
export const overridesPreferences = defineOverridesPreferences({
  app: {
    accessMode: 'backend',
    defaultHomePath: '/dashboard',
    enableRefreshToken: false,
    locale: 'zh-CN',
    name: QIXI_PLATFORM_APP_NAME,
    loginExpiredMode: 'page',
    preferencesButtonPosition: 'header',
  },
  copyright: {
    enable: false,
  },
  footer: {
    enable: false,
  },
  logo: {
    enable: true,
    fit: 'contain',
    source: platformLogoUrl,
    sourceDark: platformLogoUrl,
  },
  tabbar: {
    enable: true,
    keepAlive: true,
  },
  widget: {
    fullscreen: true,
    globalSearch: false,
    languageToggle: false,
    lockScreen: false,
    notification: false,
    refresh: true,
    sidebarToggle: true,
    themeToggle: true,
    timezone: false,
  },
  theme: {
    builtinType: 'default',
    mode: 'light',
    semiDarkHeader: false,
    semiDarkSidebar: false,
  },
});
