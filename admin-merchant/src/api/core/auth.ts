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
    app_id: number | string;
    logoUrl?: string;
    shop_name?: string;
    token: string;
    update_password?: boolean;
    user_name: string;
    version?: string;
  }
}

export async function loginApi(data: AuthApi.LoginParams) {
  return requestClient.post<AuthApi.LoginResult>('/shop/passport/login', data);
}

export async function logoutApi() {
  const accessStore = useAccessStore();
  const token = accessStore.accessToken || getDecryptedToken();
  if (!token) {
    return null;
  }
  return requestClient.post<null>('/shop/passport/logout', {});
}

/** 平台跳转商户后台时校验已有 shop JWT（需 Header 带 token） */
export async function saasLoginApi() {
  return requestClient.post<string>('/shop/passport/saasLogin', {});
}

export async function getAccessCodesApi() {
  const session = await fetchShopSessionApi();
  return session.accessCodes;
}

