import type { Recordable, UserInfo } from '@vben/types';

import { ref } from 'vue';
import { useRouter } from 'vue-router';

import { LOGIN_PATH } from '@vben/constants';
import { preferences } from '@vben/preferences';
import { resetAllStores, useAccessStore, useUserStore } from '@vben/stores';

import { ElNotification, ElMessage } from 'element-plus';
import { defineStore } from 'pinia';

import { loginApi, logoutApi } from '#/api';
import { $t } from '#/locales';
import {
  clearEncryptedToken,
  clearLegacyUserSession,
  setEncryptedToken,
  syncLegacyUserSession,
} from '#/utils/qixi-live-token';
import {
  hydratePlatformUserStore,
  resetPlatformUserStore,
  usePlatformUserStore,
} from '#/store/platform-user';
import {
  loadPlatformBootstrapData,
  PlatformBootstrapError,
} from '#/utils/platform-bootstrap';
import { clearPlatformMenuCache } from '#/utils/platform-menu';
import { formatUserFacingApiError } from '#/utils/api-error';
import {
  markJwtIssuedFromLogin,
  resetJwtSessionState,
} from '#/utils/jwt-session';

export const useAuthStore = defineStore('auth', () => {
  const accessStore = useAccessStore();
  const userStore = useUserStore();
  const router = useRouter();

  const loginLoading = ref(false);

  async function authLogin(params: Recordable<any>) {
    let userInfo: null | UserInfo = null;
    try {
      loginLoading.value = true;
      const loginData = await loginApi({
        username: params.username || params.account,
        password: params.password,
        code: params.code,
        codeKey: params.codeKey,
      });
      const accessToken = loginData?.token?.access_token;
      if (!accessToken) {
        throw new Error('登录失败：未返回 token');
      }

      accessStore.setAccessToken(accessToken);
      setEncryptedToken(accessToken);
      usePlatformUserStore().setToken(accessToken);
      clearPlatformMenuCache();
      syncLegacyUserSession(loginData.user?.account || params.username);
      hydratePlatformUserStore();
      markJwtIssuedFromLogin();

      try {
        const { userInfo: fetchedUser, accessCodes } =
          await loadPlatformBootstrapData();
        userInfo = fetchedUser;
        userStore.setUserInfo(userInfo);
        accessStore.setAccessCodes(accessCodes);
      } catch (error) {
        const message =
          error instanceof PlatformBootstrapError
            ? error.message
            : formatUserFacingApiError(error, '登录后加载失败，请重试');
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
      await logoutApi();
    } catch {
      // ignore
    }
    resetAllStores();
    clearEncryptedToken();
    clearPlatformMenuCache();
    clearLegacyUserSession();
    resetPlatformUserStore();
    resetJwtSessionState();
    accessStore.setLoginExpired(false);
    accessStore.setIsAccessChecked(false);

    await router.replace({
      path: LOGIN_PATH,
      query: redirect
        ? {
            redirect: encodeURIComponent(router.currentRoute.value.fullPath),
          }
        : {},
    });
  }

  async function fetchUserInfo() {
    const { fetchPlatformSessionApi } = await import(
      '#/api/core/platform-session'
    );
    const session = await fetchPlatformSessionApi();
    userStore.setUserInfo(session.userInfo);
    hydratePlatformUserStore();
    return session.userInfo;
  }

  function $reset() {
    loginLoading.value = false;
  }

  return {
    $reset,
    authLogin,
    fetchUserInfo,
    loginLoading,
    logout,
  };
});
