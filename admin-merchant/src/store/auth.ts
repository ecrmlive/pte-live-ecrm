import type { Recordable, UserInfo } from '@vben/types';

import { ref } from 'vue';
import { useRouter } from 'vue-router';

import { LOGIN_PATH } from '@vben/constants';
import { preferences } from '@vben/preferences';
import { resetAllStores, useAccessStore, useUserStore } from '@vben/stores';

import { ElMessage, ElNotification } from 'element-plus';
import { defineStore } from 'pinia';

import { requestClient } from '#/api/request';
import { $t } from '#/locales';
import {
  clearEncryptedToken,
  setEncryptedToken,
} from '#/utils/pte-live-token';
import {
  loadShopBootstrapData,
  ShopBootstrapError,
} from '#/utils/shop-bootstrap';
import { markJwtIssuedFromLogin, resetJwtSessionState } from '#/utils/jwt-session';

export const useAuthStore = defineStore('auth', () => {
  const accessStore = useAccessStore();
  const userStore = useUserStore();
  const router = useRouter();
  const loginLoading = ref(false);

  async function authLogin(params: Recordable<any>) {
    let userInfo: null | UserInfo = null;
    try {
      loginLoading.value = true;
      const loginData: any = await requestClient.post('/auth/login', {
        account: params.username || params.account,
        password: params.password,
      });
      const accessToken = loginData?.token?.access_token;
      if (!accessToken) {
        throw new Error('登录失败：未返回 token');
      }
      markJwtIssuedFromLogin();
      accessStore.setAccessToken(accessToken);
      setEncryptedToken(accessToken);

      try {
        const bootstrap = await loadShopBootstrapData();
        userInfo = bootstrap.userInfo;
        userStore.setUserInfo(userInfo);
        accessStore.setAccessCodes(bootstrap.accessCodes);
      } catch (error) {
        const message =
          error instanceof ShopBootstrapError
            ? error.message
            : '登录后加载失败，请重试';
        ElMessage.error(message);
        accessStore.setAccessToken(null);
        clearEncryptedToken();
        throw error;
      }

      accessStore.setIsAccessChecked(false);
      if (accessStore.loginExpired) {
        accessStore.setLoginExpired(false);
      } else {
        const redirectQuery = router.currentRoute.value.query.redirect as
          | string
          | undefined;
        const target = redirectQuery
          ? decodeURIComponent(redirectQuery)
          : userInfo?.homePath || preferences.app.defaultHomePath;
        await router.push(target);
      }

      ElNotification({
        message: `${$t('authentication.loginSuccessDesc')}:${userInfo?.realName}`,
        title: $t('authentication.loginSuccess'),
        type: 'success',
      });
    } finally {
      loginLoading.value = false;
    }
    return { userInfo };
  }

  async function logout(redirect = true) {
    try {
      await requestClient.post('/auth/logout', {});
    } catch {
      // ignore
    }
    resetAllStores();
    clearEncryptedToken();
    resetJwtSessionState();
    accessStore.setAccessToken(null);
    if (redirect) {
      await router.replace({
        path: LOGIN_PATH,
        query: { redirect: router.currentRoute.value.fullPath },
      });
    }
  }

  return { authLogin, logout, loginLoading };
});
