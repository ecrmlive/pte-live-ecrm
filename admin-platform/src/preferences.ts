import { defineOverridesPreferences } from '@vben/preferences';

const platformBrand = 'qixi';
const platformLogoUrl = `${import.meta.env.BASE_URL}${platformBrand}/logo.png`;

/**
 * 统一后台展示名称。平台、商户、区域、客服、运营共用此 Vben 应用，
 * 具体可见菜单由登录角色授权，不能把“平台”写成整个系统的身份。
 */
// 侧栏品牌保持短名称，避免在 Vben 可折叠侧栏中被截断。
export const QIXI_PLATFORM_APP_NAME = '七禧多商户';
// 登录页空间充足，使用完整系统身份说明。
export const QIXI_PLATFORM_LOGIN_NAME = '七禧多商户·管理中心';

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
  // 对齐 CRMEB 多商户后台的信息密度：一级导航与完整品牌名称都不截断。
  sidebar: {
    width: 264,
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
