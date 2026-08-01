import type { UserInfo } from '@vben/types';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { requestClient } from '#/api/request';
import type { PlatformAccessMenuItem } from '#/utils/platform-menu';
import { mapMergersMenusToAccess } from '#/utils/ecrm-menu';

export interface PlatformSession {
  menus: PlatformAccessMenuItem[];
  accessCodes: string[];
  userInfo: UserInfo;
}

export async function fetchPlatformSessionApi(): Promise<PlatformSession> {
  const [user, menusRes, permissions] = await Promise.all([
    getUserInfoApi(),
    requestClient.get<{ menus: unknown[] }>('/auth/menus'),
    getAccessCodesApi(),
  ]);
  const userInfo: UserInfo = {
    userId: String(user.admin_id || user.id),
    username: user.account || user.username,
    realName: user.real_name || user.display_name || user.account,
    avatar: '',
    desc: (user.roles || []).join(', '),
    homePath: '/dashboard',
    token: '',
  };
  return {
    userInfo,
    accessCodes: permissions,
    menus: mapMergersMenusToAccess(menusRes?.menus || []),
  };
}
