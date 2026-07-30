import { requestClient } from '#/api/request';
import { fetchShopSessionApi } from '#/api/core/shop-session';
import { getDecryptedToken } from '#/utils/qixi-live-token';
import { useAccessStore } from '@vben/stores';

export namespace AuthApi {
  export interface LoginParams {
    username: string;
    password: string;
    code?: string;
    codeKey?: string;
  }

  export interface LoginResult {
    token: {
      access_token: string;
      refresh_token: string;
      expires_in: number;
    };
    user: {
      merchant_admin_id: number;
      mer_id: number;
      store_id: number;
      mer_name: string;
      store_name: string;
      account: string;
      real_name: string;
      phone: string;
      roles: string;
    };
  }
}

export async function loginApi(data: AuthApi.LoginParams) {
  return requestClient.post<AuthApi.LoginResult>('/auth/login', {
    account: data.username,
    password: data.password,
  });
}

export async function logoutApi() {
  const accessStore = useAccessStore();
  const token = accessStore.accessToken || getDecryptedToken();
  if (!token) {
    return null;
  }
  return null;
}

/** 平台跳转商户后台时校验已有 merchant JWT（需 Header 带 token） */
export async function saasLoginApi() {
  return requestClient.get<{ account: string; store_id: number }>('/auth/me');
}

export async function getAccessCodesApi() {
  const session = await fetchShopSessionApi();
  return session.accessCodes;
}
