import { ElMessage } from 'element-plus';

import { adminApiClient } from '#/api/request';

export interface QixiLiveApiBody<T = unknown> {
  code: number;
  data: T;
  msg: string;
}

function legacyRequest<T = unknown>(
  method: 'get' | 'post' | 'upload',
  url: string,
  data?: Record<string, unknown> | FormData,
  errorback?: boolean,
) {
  const request =
    method === 'get'
      ? adminApiClient.get<QixiLiveApiBody<T>>(url, { params: data as Record<string, unknown> })
      : method === 'upload'
        ? adminApiClient.post<QixiLiveApiBody<T>>(url, data, {
            headers: { uploadImg: 'true' },
          })
        : adminApiClient.post<QixiLiveApiBody<T>>(url, data ?? {});

  return request.catch((error) => {
    if (errorback) {
      return Promise.reject(error);
    }
    throw error;
  });
}

export function adminPost<T = unknown>(
  url: string,
  data?: Record<string, unknown> | boolean,
  errorback?: boolean,
) {
  if (typeof data === 'boolean') {
    return legacyRequest<T>('post', url, {}, data);
  }
  return legacyRequest<T>('post', url, data ?? {}, errorback ?? false);
}

export function adminGet<T = unknown>(
  url: string,
  params?: Record<string, unknown>,
  errorback?: boolean,
) {
  return legacyRequest<T>('get', url, params, errorback);
}

export function adminUpload<T = unknown>(
  url: string,
  formData: FormData,
  errorback?: boolean,
) {
  return legacyRequest<T>('upload', url, formData, errorback);
}

export function adminApiErrorMessage(error: unknown) {
  const body = error as QixiLiveApiBody | undefined;
  return body?.msg || '请求失败';
}

export function showAdminApiError(error: unknown) {
  ElMessage.error(adminApiErrorMessage(error));
}
