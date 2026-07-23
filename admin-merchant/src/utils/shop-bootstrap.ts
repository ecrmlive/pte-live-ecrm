import type { UserInfo } from '@vben/types';

import { getUserInfoApi } from '#/api/core/user';
import { requestClient } from '#/api/request';
import { convertMergersMenusToVben } from '#/utils/mergers-menu';

export class ShopBootstrapError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'ShopBootstrapError';
  }
}

export async function loadShopBootstrapData() {
  try {
    const [user, menusRes] = await Promise.all([
      getUserInfoApi(),
      requestClient.get<{ menus: unknown[] }>('/auth/menus'),
    ]);
    let permissions: string[] = [];
    try {
      const perm = await requestClient.get<{ permissions: string[] }>(
        '/auth/permissions',
      );
      permissions = perm?.permissions || [];
    } catch {
      permissions = [];
    }
    const u: any = user;
    const userInfo: UserInfo = {
      userId: String(u.admin_id || u.uid || u.id || ''),
      username: u.account || u.username || '',
      realName: u.real_name || u.realName || u.account || '',
      avatar: '',
      desc: u.mer_name || '',
      homePath: '/dashboard',
      token: '',
    };
    return {
      userInfo,
      accessCodes: permissions,
      menus: convertMergersMenusToVben(menusRes?.menus || []),
      rawMenus: menusRes?.menus || [],
    };
  } catch (error) {
    throw new ShopBootstrapError(
      error instanceof Error ? error.message : '加载失败',
    );
  }
}
