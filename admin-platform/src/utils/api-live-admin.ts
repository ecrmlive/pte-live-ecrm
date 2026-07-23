import { liveAdminPost } from '#/api/live-request';

export interface ApiLiveAdminBody<T = unknown> {
  code: number;
  data: T;
  msg: string;
}

/** 平台 RBAC / 权限管理：统一走 api-platform */
export async function apiLiveAdminPost<T = unknown>(
  url: string,
  data?: Record<string, unknown>,
  errorback?: boolean,
  timeoutMs?: number,
  options?: { skipErrorMessage?: boolean },
) {
  try {
    const res = await liveAdminPost<T>(url, data, timeoutMs, options);
    return { code: res.code, msg: res.msg, data: res.data } as ApiLiveAdminBody<T>;
  } catch (error) {
    if (errorback) {
      return Promise.reject(error);
    }
    throw error;
  }
}
