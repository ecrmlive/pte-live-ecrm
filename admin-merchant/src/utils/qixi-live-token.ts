import CryptoJS from 'crypto-js';

import {
  QIXI_SHOP_MENU_KEY,
  QIXI_SHOP_RENDER_MENU_KEY,
  QIXI_SHOP_TOKEN_KEY,
} from './qixi-live-api';

const SECRET_KEY = 'jjj_shop_single_admin_2024';

function readEncryptedTokenRaw(name: string) {
  return sessionStorage.getItem(name) || localStorage.getItem(name);
}

export function setEncryptedToken(token: string, name = QIXI_SHOP_TOKEN_KEY) {
  if (!token) return;
  const encrypted = CryptoJS.AES.encrypt(
    JSON.stringify(token),
    SECRET_KEY,
  ).toString();
  sessionStorage.setItem(name, encrypted);
  localStorage.setItem(name, encrypted);
}

export function getDecryptedToken(name = QIXI_SHOP_TOKEN_KEY) {
  const encrypted = readEncryptedTokenRaw(name);
  if (!encrypted) return null;
  try {
    const decrypted = CryptoJS.AES.decrypt(encrypted, SECRET_KEY);
    const token = decrypted.toString(CryptoJS.enc.Utf8);
    return JSON.parse(token) as string;
  } catch {
    return null;
  }
}

export function clearEncryptedToken(name = QIXI_SHOP_TOKEN_KEY) {
  sessionStorage.removeItem(name);
  localStorage.removeItem(name);
}

export interface LegacyShopUserInfo {
	store_app_id?: string;
	AppID?: number | string;
  app_id?: number | string;
  logoUrl?: string;
  shopName?: string;
  userName?: string;
  version?: string;
}

export function syncLegacyUserSession(info: LegacyShopUserInfo) {
  sessionStorage.setItem('userInfo', JSON.stringify(info));
}

export function getLegacyUserInfo(): LegacyShopUserInfo | null {
  const raw = sessionStorage.getItem('userInfo');
  if (!raw) return null;
  try {
    return JSON.parse(raw) as LegacyShopUserInfo;
  } catch {
    return null;
  }
}

export function clearLegacyUserSession() {
  sessionStorage.removeItem('userInfo');
  sessionStorage.removeItem(QIXI_SHOP_MENU_KEY);
  sessionStorage.removeItem(QIXI_SHOP_RENDER_MENU_KEY);
}

/** 新标签打开中控：从 opener 复制 token / userInfo / 菜单缓存 */
export function syncAuthFromOpener() {
  if (!window.opener || window.opener === window) return;

  const keys = [
    QIXI_SHOP_TOKEN_KEY,
    'userInfo',
    QIXI_SHOP_MENU_KEY,
    QIXI_SHOP_RENDER_MENU_KEY,
  ];

  try {
    const openerStorage = window.opener.sessionStorage;
    if (!openerStorage) return;
    for (const key of keys) {
      if (sessionStorage.getItem(key)) continue;
      const value = openerStorage.getItem(key);
      if (value) sessionStorage.setItem(key, value);
    }
  } catch {
    // opener 已关闭或跨域
  }
}

/** 刷新页面时：Pinia 与加密 token 对齐（Pinia 优先，避免回写陈旧 token） */
export function syncMerchantAccessToken(
  setToken: (token: string) => void,
  getToken: () => string | null | undefined,
) {
  syncAuthFromOpener();
  const stored = getDecryptedToken()?.trim() || null;
  const current = getToken()?.trim() || null;

  if (stored && !current) {
    setToken(stored);
    return stored;
  }
  if (current && !stored) {
    setEncryptedToken(current);
    return current;
  }
  if (current && stored && current !== stored) {
    setEncryptedToken(current);
    return current;
  }
  return current || stored;
}

/** 路由守卫 / bootstrap 统一取 token：Pinia 与加密存储取其一并对齐 */
export function resolveMerchantAccessToken(
  getToken: () => string | null | undefined,
  setToken: (token: string) => void,
) {
  return syncMerchantAccessToken(setToken, getToken);
}

export function hydrateAccessTokenFromLegacy(
  setToken: (token: string) => void,
) {
  syncAuthFromOpener();
  const stored = getDecryptedToken();
  if (stored) {
    setToken(stored);
  }
}
