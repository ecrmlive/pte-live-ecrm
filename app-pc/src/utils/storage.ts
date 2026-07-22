const TOKEN_KEY = "qixi_pc_token";
const REFRESH_KEY = "qixi_pc_refresh";
const USER_KEY = "qixi_pc_user";

export function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) || "";
}

export function getRefreshToken(): string {
  return localStorage.getItem(REFRESH_KEY) || "";
}

export function setTokenPair(access: string, refresh?: string) {
  localStorage.setItem(TOKEN_KEY, access);
  if (refresh) localStorage.setItem(REFRESH_KEY, refresh);
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(REFRESH_KEY);
}

export function getUserJSON<T>(): T | null {
  const raw = localStorage.getItem(USER_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as T;
  } catch {
    return null;
  }
}

export function setUserJSON(user: unknown) {
  localStorage.setItem(USER_KEY, JSON.stringify(user));
}

export function clearUser() {
  localStorage.removeItem(USER_KEY);
}
