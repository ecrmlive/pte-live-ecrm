import axios from 'axios';
import { message } from 'ant-design-vue';

const TOKEN_KEY = 'qx_platform_access_token';
const REFRESH_KEY = 'qx_platform_refresh_token';

export const tokenStore = {
  getAccess: () => localStorage.getItem(TOKEN_KEY) || '',
  getRefresh: () => localStorage.getItem(REFRESH_KEY) || '',
  set(access: string, refresh: string) {
    localStorage.setItem(TOKEN_KEY, access);
    localStorage.setItem(REFRESH_KEY, refresh);
  },
  clear() {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(REFRESH_KEY);
  },
};

export const http = axios.create({
  baseURL: '/api/platform/v1',
  timeout: 15000,
});

http.interceptors.request.use((config) => {
  const token = tokenStore.getAccess();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

http.interceptors.response.use(
  (res) => {
    const body = res.data;
    if (body && typeof body.status === 'number' && body.status !== 200) {
      message.error(body.message || '请求失败');
      return Promise.reject(body);
    }
    return body?.data !== undefined ? { ...res, data: body.data } : res;
  },
  (err) => {
    const status = err?.response?.status;
    const msg = err?.response?.data?.message || err.message || '网络错误';
    if (status === 401) {
      tokenStore.clear();
      if (!location.pathname.startsWith('/login')) {
        location.href = `/login?redirect=${encodeURIComponent(location.pathname)}`;
      }
    } else {
      message.error(msg);
    }
    return Promise.reject(err);
  },
);
