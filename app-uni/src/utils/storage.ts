const TOKEN_KEY = "qixi_app_token";
const REFRESH_KEY = "qixi_app_refresh";
const USER_KEY = "qixi_app_user";

export function getToken(): string {
  return uni.getStorageSync(TOKEN_KEY) || "";
}

export function getRefreshToken(): string {
  return uni.getStorageSync(REFRESH_KEY) || "";
}

export function setTokenPair(access: string, refresh?: string) {
  uni.setStorageSync(TOKEN_KEY, access);
  if (refresh) uni.setStorageSync(REFRESH_KEY, refresh);
}

export function clearToken() {
  uni.removeStorageSync(TOKEN_KEY);
  uni.removeStorageSync(REFRESH_KEY);
}

export function getUserJSON<T>(): T | null {
  const raw = uni.getStorageSync(USER_KEY);
  if (!raw) return null;
  try {
    return typeof raw === "string" ? (JSON.parse(raw) as T) : (raw as T);
  } catch {
    return null;
  }
}

export function setUserJSON(user: unknown) {
  uni.setStorageSync(USER_KEY, JSON.stringify(user));
}

export function clearUser() {
  uni.removeStorageSync(USER_KEY);
}
