import { publicRequestClient, requestClient } from '#/api/request';

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
  // 使用 publicRequestClient：退出接口缺失/失败时不应弹全局 Toast；
  // 本地 token 清理由 authStore.logout 保证。
  return publicRequestClient.post<null>('/auth/logout', {}).catch(() => null);
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
