import { defineOverridesPreferences } from '@vben/preferences';

const merchantBrand = 'qixi';
const merchantLogoUrl = `${import.meta.env.BASE_URL}${merchantBrand}/logo.png`;

/** 栖息多商户商户后台 — Vben 偏好（仅 zh-CN，后端菜单） */
export const overridesPreferences = defineOverridesPreferences({
  app: {
    accessMode: 'backend',
    defaultHomePath: '/dashboard',
    enableRefreshToken: false,
    locale: 'zh-CN',
    name: '栖息多商户·商户',
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
    source: merchantLogoUrl,
    sourceDark: merchantLogoUrl,
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
});
