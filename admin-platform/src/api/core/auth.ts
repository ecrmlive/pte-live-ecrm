import { requestClient } from '#/api/request';

export namespace AuthApi {
  export interface LoginParams {
    account: string;
    password: string;
  }

  export interface TokenPair {
    access_token: string;
    refresh_token: string;
    expires_in: number;
  }

  export interface PlatformUser {
    admin_id: number;
    account: string;
    real_name: string;
    phone: string;
    roles: string;
    level: number;
  }

  export interface LoginResult {
    token: TokenPair;
    user: PlatformUser;
  }
}

export async function loginApi(data: {
  username?: string;
  account?: string;
  password: string;
  code?: string;
  codeKey?: string;
}) {
  const account = data.account || data.username || '';
  return requestClient.post<AuthApi.LoginResult>('/auth/login', {
    account,
    password: data.password,
  });
}

export async function logoutApi() {
  return requestClient.post<null>('/auth/logout', {}).catch(() => null);
}

export async function getAccessCodesApi() {
  const res = await requestClient.get<{ permissions: string[] }>(
    '/auth/permissions',
  );
  return res?.permissions || [];
}

export async function getUserInfoApi() {
  return requestClient.get<AuthApi.PlatformUser>('/auth/me');
}
