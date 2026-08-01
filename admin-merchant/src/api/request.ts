import type { RequestClientOptions } from '@vben/request';

import { LOGIN_PATH } from '@vben/constants';
import { preferences } from '@vben/preferences';
import {
  defaultResponseInterceptor,
  errorMessageResponseInterceptor,
  RequestClient,
} from '@vben/request';
import { useAccessStore, useUserStore } from '@vben/stores';

import { ElMessage } from 'element-plus';

import { resolveApiBaseUrl } from '#/utils/pte-live-api';
import { attachShopAppId } from '#/utils/pte-live-shop-app-id';
import {
  clearEncryptedToken,
  clearLegacyUserSession,
  getDecryptedToken,
  setEncryptedToken,
} from '#/utils/pte-live-token';

import { useAuthStore } from '#/store';
import { isApiUnauthorized, resetJwtSessionState } from '#/utils/jwt-session';
import {
  attachAPIEncryption,
  decryptAPIResponse,
} from '#/utils/api-crypto';

const apiURL = `${resolveApiBaseUrl()}/api/merchant/v1`;

function isOnLoginPage() {
  if (typeof window === 'undefined') {
    return false;
  }
  const path = window.location.pathname || '';
  const hash = window.location.hash || '';
  return path.includes(LOGIN_PATH) || hash.includes(LOGIN_PATH);
}

function formatAuthHeader(token: null | string) {
  return token ? `Bearer ${token}` : null;
}

function createRequestClient(baseURL: string, options?: RequestClientOptions) {
  const client = new RequestClient({
    ...options,
    baseURL,
    headers: {
      'Content-Type': 'application/json;charset=UTF-8',
    },
  });

  async function doReAuthenticate() {
    const accessStore = useAccessStore();
    const authStore = useAuthStore();
    accessStore.setAccessToken(null);
    clearEncryptedToken();
    clearLegacyUserSession();
    resetJwtSessionState();
    if (isOnLoginPage()) {
      return;
    }
    await authStore.logout(false);
  }

  client.addRequestInterceptor({
    fulfilled: async (config) => {
      const accessStore = useAccessStore();
      const userStore = useUserStore();
      const token = accessStore.accessToken || getDecryptedToken() || null;
      const bearer = formatAuthHeader(token);
      if (bearer) {
        config.headers['Authori-zation'] = bearer;
      }
      attachShopAppId(config, userStore.userInfo, token || undefined);
		config.headers['Accept-Language'] = preferences.app.locale;
      return attachAPIEncryption(config, baseURL);
    },
  });

  client.addResponseInterceptor({
    fulfilled: async (response) => {
      await decryptAPIResponse(response);
      const nextAuth = response.headers['authori-zation'];
      if (typeof nextAuth === 'string' && nextAuth.startsWith('Bearer ')) {
        const nextToken = nextAuth.slice(7).trim();
        const accessStore = useAccessStore();
        accessStore.setAccessToken(nextToken);
        setEncryptedToken(nextToken);
      }
      return response;
    },
  });

  client.addResponseInterceptor(
    defaultResponseInterceptor({
      codeField: 'status',
      dataField: 'data',
      successCode: 200,
    }),
  );

  client.addResponseInterceptor({
    rejected: async (error) => {
      const status = error?.response?.status;
      const code = error?.response?.data?.code ?? error?.data?.code;
      if (isApiUnauthorized(status, code)) {
        await doReAuthenticate();
        error.__authHandled = true;
      }
      throw error;
    },
  });

  client.addResponseInterceptor(
    errorMessageResponseInterceptor((msg: string, error) => {
      if (error?.__authHandled) {
        return;
      }
      const responseData = error?.response?.data ?? error?.data ?? {};
      const bodyCode = responseData?.code;
      if (bodyCode === -1) {
        return;
      }
      const errorMessage = responseData?.msg ?? responseData?.message ?? msg;
      ElMessage.error(errorMessage || msg);
    }),
  );

  return client;
}

export const requestClient = createRequestClient(apiURL, {
  responseReturn: 'data',
});

export const baseRequestClient = new RequestClient({ baseURL: apiURL });
