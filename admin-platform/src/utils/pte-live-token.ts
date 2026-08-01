import CryptoJS from 'crypto-js';

import { QIXI_ADMIN_TOKEN_KEY } from './pte-live-api';

const SECRET_KEY = 'jjj_shop_single_admin_2024';

export function setEncryptedToken(token: string, name = QIXI_ADMIN_TOKEN_KEY) {
  if (!token) return;
  const encrypted = CryptoJS.AES.encrypt(JSON.stringify(token), SECRET_KEY).toString();
  sessionStorage.setItem(name, encrypted);
}

export function getDecryptedToken(name = QIXI_ADMIN_TOKEN_KEY) {
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

export function clearEncryptedToken(name = QIXI_ADMIN_TOKEN_KEY) {
  sessionStorage.removeItem(name);
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
