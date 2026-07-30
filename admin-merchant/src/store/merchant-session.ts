import { defineStore } from 'pinia';

import {
  clearEncryptedToken,
  getDecryptedToken,
  getLegacyUserInfo,
  setEncryptedToken,
} from '#/utils/qixi-live-token';

export interface MerchantSessionUserInfo {
	store_app_id?: string;
	AppID?: number | string;
  app_id?: number | string;
  logoUrl?: string;
  shopName?: string;
  shopUserId?: number | string;
  userId?: number | string;
  userName?: string;
  version?: string;
  [key: string]: unknown;
}

export const useMerchantSessionStore = defineStore('merchantSession', {
  state: () => ({
    token: getDecryptedToken() as string | null,
    userInfo: getLegacyUserInfo() as MerchantSessionUserInfo | null,
  }),
  actions: {
    hydrate() {
      const token = getDecryptedToken();
      if (token) {
        this.token = token;
      }
      const userInfo = getLegacyUserInfo();
      if (userInfo) {
        this.userInfo = { ...(this.userInfo || {}), ...userInfo };
      }
    },
    setToken(token: string | null) {
      this.token = token;
      if (token) {
        setEncryptedToken(token);
      } else {
        clearEncryptedToken();
      }
    },
    afterLogout() {
      this.token = null;
      this.userInfo = null;
      clearEncryptedToken();
    },
  },
});

export default useMerchantSessionStore;
