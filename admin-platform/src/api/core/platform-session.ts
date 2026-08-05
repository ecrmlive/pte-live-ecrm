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

function firstAccessiblePagePath(menus: PlatformAccessMenuItem[]): string {
  for (const item of menus) {
    if (item.is_route === 1 && item.path) {
      return item.path;
    }
    const childPath = firstAccessiblePagePath(item.children || []);
    if (childPath) {
      return childPath;
    }
  }
  return '';
}

export async function fetchPlatformSessionApi(): Promise<PlatformSession> {
  const [user, menusRes, permissions] = await Promise.all([
    getUserInfoApi(),
    requestClient.get<{ menus: unknown[] }>('/auth/menus'),
    getAccessCodesApi(),
  ]);
  const menus = mapMergersMenusToAccess(menusRes?.menus || []);
  const homePath = firstAccessiblePagePath(menus) || '/auth/service-error';
  const userInfo: UserInfo = {
    userId: String(user.admin_id || user.id),
    username: user.account || user.username,
    realName: user.real_name || user.display_name || user.account,
    avatar: '',
    desc: (user.roles || []).join(', '),
    homePath,
    roles: user.roles || [],
    token: '',
  };
  return {
    userInfo,
    accessCodes: permissions,
    menus,
  };
}
