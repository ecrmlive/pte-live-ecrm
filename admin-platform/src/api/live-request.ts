import axios from 'axios';
import qs from 'qs';

import { ElMessage } from 'element-plus';

import { resolveLiveApiBaseUrl } from '#/utils/qixi-live-api';
import { applyAdminAuthorization } from '#/utils/live-api-auth';
import { attachAPIEncryption, decryptAPIResponse } from '#/utils/api-crypto';

export interface LiveApiBody<T = unknown> {
  code: number;
  data: T;
  msg: string;
}

function createLiveClient(baseURL: string) {
  const client = axios.create({
    baseURL,
    timeout: 12_000,
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
    },
  });

  client.interceptors.request.use(async (config) => {
    config.headers = applyAdminAuthorization(
      config.headers as Record<string, unknown>,
    ) as typeof config.headers;
    if (
      config.method === 'post' &&
      !(config.data instanceof FormData) &&
      typeof config.data !== 'string'
    ) {
      config.data = qs.stringify(config.data || {});
    }
    return attachAPIEncryption(config);
  });

  client.interceptors.response.use(
    async (res) => {
      await decryptAPIResponse(res);
      if (res.data.code === 1) {
        return res.data as LiveApiBody;
      }
      const msg = res.data.msg || '请求失败';
      const code = res.data.code;
      const isAuthError =
        code === -1 ||
        msg.includes('鉴权') ||
        msg.includes('token');
      const skipErrorMessage = Boolean(
        (res.config as { skipErrorMessage?: boolean }).skipErrorMessage,
      );
      // 鉴权失败由路由守卫统一提示，避免与「登录已失效」重复弹窗
      if (!isAuthError && !skipErrorMessage) {
        ElMessage.error(msg);
      }
      return Promise.reject(res.data);
    },
    () => Promise.reject(new Error('api-platform 请求失败')),
  );

  return client;
}

export const liveAdminClient = createLiveClient(resolveLiveApiBaseUrl());

export async function liveAdminPost<T = unknown>(
  url: string,
  data?: Record<string, unknown>,
  timeoutMs?: number,
  options?: { skipErrorMessage?: boolean },
) {
  return liveAdminClient.post(url, data ?? {}, {
    timeout: timeoutMs ?? 12_000,
    skipErrorMessage: options?.skipErrorMessage,
  }) as Promise<LiveApiBody<T>>;
}
