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
    id: number;
    admin_id: number;
    username: string;
    account: string;
    display_name: string;
    real_name: string;
    roles: string[];
    data_scope_version: number;
    /** 兼容旧页面：1 表示区域角色。新权限判断以 roles 为准。 */
    is_agent: number;
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
