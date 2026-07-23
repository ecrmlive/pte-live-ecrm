import { defineStore } from 'pinia';

import {
  clearEncryptedToken,
  clearLegacyUserSession,
  getDecryptedToken,
  getLegacyUserInfo,
  setEncryptedToken,
} from '#/utils/qixi-live-token';
import { addSessionStorage } from '#/utils/base';

type BusHandler = (data: unknown) => void;

export const usePlatformUserStore = defineStore('platform-user', {
  state: () => ({
    token: getDecryptedToken() || '',
    userInfo: getLegacyUserInfo() as null | Record<string, unknown>,
    list: {} as Record<string, BusHandler[]>,
    isDefaultPassword: false,
  }),
  actions: {
    bus_on(name: string, fn: BusHandler) {
      this.list[name] = this.list[name] || [];
      this.list[name].push(fn);
    },
    bus_emit(name: string, data: unknown) {
      this.list[name]?.forEach((fn) => fn(data));
    },
    bus_off(name: string) {
      delete this.list[name];
    },
    setToken(token: string | null) {
      this.token = token || '';
      if (token) {
        setEncryptedToken(token);
      } else {
        clearEncryptedToken();
      }
    },
    afterLogin(info: { data: { token: string; update_password?: boolean; user_name: string } }) {
      this.userInfo = this.userInfo || {};
      const {
        data: { token, user_name, update_password },
      } = info;
      this.userInfo.userName = user_name;
      this.userInfo.username = user_name;
      this.isDefaultPassword = update_password || false;
      this.setToken(token);
      addSessionStorage('userInfo', this.userInfo);
    },
    afterLogout() {
      clearEncryptedToken();
      clearLegacyUserSession();
      this.token = '';
      this.userInfo = null;
      this.isDefaultPassword = false;
    },
  },
});

export function hydratePlatformUserStore() {
  const store = usePlatformUserStore();
  const token = getDecryptedToken();
  if (token) {
    store.token = token;
  }
  const userInfo = getLegacyUserInfo();
  if (userInfo) {
    store.userInfo = {
      ...(store.userInfo || {}),
      ...userInfo,
      username: userInfo.username || userInfo.userName,
    };
  }
}

export function resetPlatformUserStore() {
  usePlatformUserStore().afterLogout();
}

/** 兼容 native 页面 `useUserStore` 命名 */
export const useUserStore = usePlatformUserStore;
