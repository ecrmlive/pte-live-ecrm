import type { UserInfo } from '@vben/types';

import type { QixiLiveApiBody } from '#/utils/admin-api';
import { adminPost } from '#/utils/admin-api';
import { requestClient } from '#/api/request';

interface PlatformUserInfoResponse {
  homePath?: string;
  is_super?: number;
  roles?: string[];
  user_name?: string;
}

/** @deprecated 平台进后台请用 api-platform fetchPlatformSessionApi */
export async function getUserInfoApi(): Promise<UserInfo> {
  const res = await requestClient.post<PlatformUserInfoResponse>(
    '/admin/platform.auth/getUserInfo',
    {},
  );
  const userName = res?.user_name || '管理员';
  const roles =
    res?.roles && res.roles.length > 0 ? res.roles : ['platform_user'];
  return {
    avatar: '',
    desc: '七禧 平台超管',
    homePath: res?.homePath || '/Home',
    realName: userName,
    roles,
    userId: '0',
    username: userName,
  };
}

export function editPlatformPasswordApi(
  data: Record<string, unknown>,
  errorback?: boolean,
) {
  return adminPost<QixiLiveApiBody>('/admin/admin.user/renew', data, errorback);
}
