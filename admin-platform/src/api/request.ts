/**
 * 七禧多商户平台 API 请求适配。
 */
import type { RequestClientOptions } from '@vben/request';

import { LOGIN_PATH } from '@vben/constants';
import { preferences } from '@vben/preferences';
import {
  defaultResponseInterceptor,
  errorMessageResponseInterceptor,
  RequestClient,
} from '@vben/request';
import { useAccessStore } from '@vben/stores';

import { ElMessage } from 'element-plus';

import { resolveApiBaseUrl } from '#/utils/pte-live-api';
import { formatUserFacingApiError } from '#/utils/api-error';
import {
  clearEncryptedToken,
  clearLegacyUserSession,
  getDecryptedToken,
  setEncryptedToken,
} from '#/utils/pte-live-token';

import { useAuthStore } from '#/store';
import { isApiUnauthorized, resetJwtSessionState } from '#/utils/jwt-session';
import { attachAPIEncryption, decryptAPIResponse } from '#/utils/api-crypto';

const apiURL = `${resolveApiBaseUrl()}/api/platform/v1`;

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

function createRequestClient(
  baseURL: string,
  options?: RequestClientOptions & { skipErrorMessage?: boolean },
) {
  const { skipErrorMessage = false, ...clientOptions } = options ?? {};
  const client = new RequestClient({
    ...clientOptions,
    baseURL,
    timeout: 12_000,
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
      const token =
        accessStore.accessToken || getDecryptedToken() || null;
      const bearer = formatAuthHeader(token);
      if (bearer) {
        config.headers['Authori-zation'] = bearer;
      }
		config.headers['Accept-Language'] = preferences.app.locale;
		return attachAPIEncryption(config);
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
      }
      throw error;
    },
  });

  if (!skipErrorMessage) {
    client.addResponseInterceptor(
      errorMessageResponseInterceptor((msg: string, error) => {
        ElMessage.error(formatUserFacingApiError(error, msg));
      }),
    );
  }

  return client;
}

export const requestClient = createRequestClient(apiURL, {
  responseReturn: 'data',
});

/** 登录页等场景：不弹全局错误 Toast，由页面自行提示 */
export const publicRequestClient = createRequestClient(apiURL, {
  responseReturn: 'data',
  skipErrorMessage: true,
});

export const baseRequestClient = new RequestClient({ baseURL: apiURL });

/** 返回 `{ code, data, msg }` 全量结构，供 admin-api 使用 */
function createAdminApiClient() {
  const client = new RequestClient({
    baseURL: apiURL,
    timeout: 12_000,
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
      const token = accessStore.accessToken || getDecryptedToken() || null;
      const bearer = formatAuthHeader(token);
      if (bearer) {
        config.headers['Authori-zation'] = bearer;
      }
		config.headers['Accept-Language'] = preferences.app.locale;
		return attachAPIEncryption(config);
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

  client.addResponseInterceptor({
    fulfilled: (response) => {
      const body = response.data as { code?: number; data?: unknown; msg?: string };
      if (body?.code === 1) {
        return body as typeof response;
      }
      if (body?.code === -1 || response.status === 401) {
        void doReAuthenticate();
        return Promise.reject(body);
      }
      ElMessage.error(
        formatUserFacingApiError(body, body?.msg || '请求失败'),
      );
      return Promise.reject(body);
    },
    rejected: async (error) => {
      const status = error?.response?.status;
      if (status) {
        const fallback =
          status >= 500
            ? '接口服务异常，请稍后再试'
            : `接口请求失败 (HTTP ${status})`;
        ElMessage.error(formatUserFacingApiError(error, fallback));
      } else {
        ElMessage.error(formatUserFacingApiError(error, '接口服务无法连接'));
      }
      return Promise.reject(error);
    },
  });

  return client;
}

export const adminApiClient = createAdminApiClient();
