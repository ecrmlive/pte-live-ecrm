import { apiBaseURL, APP_API_PREFIX, CALLBACK_API_PREFIX } from "@/config/env";
import { clearToken, getToken } from "@/utils/storage";

export interface ApiResult<T = unknown> {
  status: number;
  message: string;
  data: T;
}

type RequestMethod = "GET" | "POST" | "PUT" | "DELETE";

export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

export function request<T>(
  path: string,
  options: {
    method?: RequestMethod;
    data?: Record<string, unknown> | unknown;
    auth?: boolean;
  } = {}
): Promise<T> {
  const method = options.method || "GET";
  const auth = options.auth !== false;
  const url = `${apiBaseURL()}${APP_API_PREFIX}${path.startsWith("/") ? path : `/${path}`}`;
  const header: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (auth) {
    const token = getToken();
    if (token) header["Authori-zation"] = `Bearer ${token}`;
  }

  return new Promise((resolve, reject) => {
    uni.request({
      url,
      method,
      data: options.data as UniApp.RequestOptions["data"],
      header,
      success: (res) => {
        const body = res.data as ApiResult<T>;
        if (res.statusCode === 401 || body?.status === 401) {
          clearToken();
          reject(new ApiError(401, body?.message || "请先登录"));
          return;
        }
        if (res.statusCode >= 200 && res.statusCode < 300) {
          if (body && typeof body.status === "number") {
            if (body.status === 0 || body.status === 200) {
              resolve(body.data);
              return;
            }
            reject(new ApiError(body.status, body.message || "请求失败"));
            return;
          }
          resolve(body as unknown as T);
          return;
        }
        reject(new ApiError(res.statusCode, body?.message || `HTTP ${res.statusCode}`));
      },
      fail: (err) => {
        reject(new ApiError(-1, err.errMsg || "网络异常"));
      },
    });
  });
}

export const http = {
  get: <T>(path: string, auth = true) => request<T>(path, { method: "GET", auth }),
  post: <T>(path: string, data?: Record<string, unknown>, auth = true) =>
    request<T>(path, { method: "POST", data, auth }),
  put: <T>(path: string, data?: Record<string, unknown>, auth = true) =>
    request<T>(path, { method: "PUT", data, auth }),
  delete: <T>(path: string, auth = true) => request<T>(path, { method: "DELETE", auth }),
};

/** 支付回调（无 JWT，前缀 /api/callback/v1）。 */
export function callbackPost<T>(path: string, data?: Record<string, unknown>): Promise<T> {
  const url = `${apiBaseURL()}${CALLBACK_API_PREFIX}${path.startsWith("/") ? path : `/${path}`}`;
  return new Promise((resolve, reject) => {
    uni.request({
      url,
      method: "POST",
      data: data as UniApp.RequestOptions["data"],
      header: { "Content-Type": "application/json" },
      success: (res) => {
        const body = res.data as ApiResult<T>;
        if (res.statusCode >= 200 && res.statusCode < 300) {
          if (body && typeof body.status === "number") {
            if (body.status === 0 || body.status === 200) {
              resolve(body.data);
              return;
            }
            reject(new ApiError(body.status, body.message || "请求失败"));
            return;
          }
          resolve(body as unknown as T);
          return;
        }
        reject(new ApiError(res.statusCode, body?.message || `HTTP ${res.statusCode}`));
      },
      fail: (err) => reject(new ApiError(-1, err.errMsg || "网络异常")),
    });
  });
}
