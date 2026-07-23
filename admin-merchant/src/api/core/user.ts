import type { UserInfo } from '@vben/types';

import {
  fetchShopSessionApi,
  mapShopSessionUser,
} from '#/api/core/shop-session';
import {
  applyTenantAppBranding,
  resolveShopDisplayName,
} from '#/utils/shop-display-name';

export async function getUserInfoApi(): Promise<UserInfo> {
  const session = await fetchShopSessionApi();
  const userInfo = session.userInfo;
  applyTenantAppBranding(
    (userInfo as any).shopName,
    (userInfo as any).logoUrl,
  );
  return userInfo;
}

/** 兼容旧调用：从 session 用户对象映射 */
export async function getUserInfoFromSession(user: {
  app_id?: number | string;
  logoUrl?: string;
  shop_name?: string;
  user_name?: string;
  version?: string;
}): Promise<UserInfo> {
  const shopName = resolveShopDisplayName(user?.shop_name);
  applyTenantAppBranding(user?.shop_name, user?.logoUrl);
  return mapShopSessionUser({
    ...user,
    shop_name: shopName,
    homePath: '/home',
  });
}
