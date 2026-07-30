import { APP_API_PREFIX, CALLBACK_API_PREFIX, apiBaseURL } from "@/config/env";
import { clearToken, getMerchantAppID, getToken } from "@/utils/storage";

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

export async function request<T>(
  path: string,
  options: {
    method?: RequestMethod;
    data?: Record<string, unknown> | unknown;
    auth?: boolean;
  } = {},
): Promise<T> {
  const method = options.method || "GET";
  const auth = options.auth !== false;
  const url = `${apiBaseURL()}${APP_API_PREFIX}${path.startsWith("/") ? path : `/${path}`}`;
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (auth) {
    const token = getToken();
    if (token) headers["Authori-zation"] = `Bearer ${token}`;
  }
  const merchantAppID = getMerchantAppID();
  if (merchantAppID) headers["X-AppId"] = merchantAppID;

  const res = await fetch(url, {
    method,
    headers,
    body: method === "GET" || method === "DELETE" ? undefined : JSON.stringify(options.data ?? {}),
  });

  let body: ApiResult<T> | null = null;
  try {
    body = (await res.json()) as ApiResult<T>;
  } catch {
    body = null;
  }

  if (res.status === 401 || body?.status === 401) {
    clearToken();
    throw new ApiError(401, body?.message || "请先登录");
  }

  if (res.ok) {
    if (body && typeof body.status === "number") {
      if (body.status === 0 || body.status === 200) return body.data;
      throw new ApiError(body.status, body.message || "请求失败");
    }
    return body as unknown as T;
  }

  throw new ApiError(res.status, body?.message || `HTTP ${res.status}`);
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
export async function callbackPost<T>(path: string, data?: Record<string, unknown>): Promise<T> {
  const url = `${apiBaseURL()}${CALLBACK_API_PREFIX}${path.startsWith("/") ? path : `/${path}`}`;
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data ?? {}),
  });
  let body: ApiResult<T> | null = null;
  try {
    body = (await res.json()) as ApiResult<T>;
  } catch {
    body = null;
  }
  if (res.ok) {
    if (body && typeof body.status === "number") {
      if (body.status === 0 || body.status === 200) return body.data;
      throw new ApiError(body.status, body.message || "请求失败");
    }
    return body as unknown as T;
  }
  throw new ApiError(res.status, body?.message || `HTTP ${res.status}`);
}
