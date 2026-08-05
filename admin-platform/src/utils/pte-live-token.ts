import CryptoJS from 'crypto-js';

import { QIXI_ADMIN_TOKEN_KEY } from './pte-live-api';

const SECRET_KEY = 'jjj_shop_single_admin_2024';
const QIXI_ADMIN_REFRESH_TOKEN_KEY = `${QIXI_ADMIN_TOKEN_KEY}:refresh`;

function setEncryptedValue(token: string, name: string) {
  if (!token) return;
  const encrypted = CryptoJS.AES.encrypt(JSON.stringify(token), SECRET_KEY).toString();
  sessionStorage.setItem(name, encrypted);
}

function getDecryptedValue(name: string) {
  const encrypted = sessionStorage.getItem(name);
  if (!encrypted) return null;
  try {
    const decrypted = CryptoJS.AES.decrypt(encrypted, SECRET_KEY);
    const token = decrypted.toString(CryptoJS.enc.Utf8);
    return JSON.parse(token) as string;
  } catch {
    return null;
  }
}

export function setEncryptedToken(token: string, name = QIXI_ADMIN_TOKEN_KEY) {
  setEncryptedValue(token, name);
}

export function getDecryptedToken(name = QIXI_ADMIN_TOKEN_KEY) {
  return getDecryptedValue(name);
}

export function setEncryptedRefreshToken(token: string) {
  setEncryptedValue(token, QIXI_ADMIN_REFRESH_TOKEN_KEY);
}

export function getDecryptedRefreshToken() {
  return getDecryptedValue(QIXI_ADMIN_REFRESH_TOKEN_KEY);
}

export function clearEncryptedToken(name = QIXI_ADMIN_TOKEN_KEY) {
  sessionStorage.removeItem(name);
  if (name === QIXI_ADMIN_TOKEN_KEY) {
    sessionStorage.removeItem(QIXI_ADMIN_REFRESH_TOKEN_KEY);
  }
}

export function syncLegacyUserSession(userName: string) {
  sessionStorage.setItem(
    'userInfo',
    JSON.stringify({
      AppID: 10000,
      userName,
      username: userName,
    }),
  );
}

export function getLegacyUserInfo(): null | Record<string, unknown> {
  const raw = sessionStorage.getItem('userInfo');
  if (!raw) return null;
  try {
    return JSON.parse(raw) as Record<string, unknown>;
  } catch {
    return null;
  }
}

export function clearLegacyUserSession() {
  sessionStorage.removeItem('userInfo');
}

export function hydrateAccessTokenFromLegacy(
  setToken: (token: string) => void,
) {
  const legacy = getDecryptedToken();
  if (legacy) {
    setToken(legacy);
  }
}
