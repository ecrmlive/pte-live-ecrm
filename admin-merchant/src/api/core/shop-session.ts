import type { UserInfo } from '@vben/types';

import { requestClient } from '#/api/request';
import type { ShopAccessMenuItem } from '#/utils/qixi-live-menu';
import { QIXI_SHOP_MENU_KEY } from '#/utils/qixi-live-api';
import {
  applyTenantAppBranding,
  resolveShopDisplayName,
} from '#/utils/shop-display-name';
import { syncShopPathAuthCache } from '#/utils/shop-path-auth';

export interface ShopSessionResponse {
  codes: string[];
  menus: ShopAccessMenuItem[];
  user: {
    app_id?: number | string;
    homePath?: string;
    is_super?: number;
    logoUrl?: string;
    roles?: string[];
    shop_name?: string;
    user_name?: string;
    version?: string;
  };
}

interface MerchantProfileResponse {
  merchant_admin_id: number;
  mer_id: number;
  mer_name: string;
  account: string;
  real_name: string;
  phone: string;
  roles: string;
}

export function mapShopSessionUser(
  user: ShopSessionResponse['user'],
): UserInfo {
  const shopName = resolveShopDisplayName(user?.shop_name);
  const userName = user?.user_name || shopName;
  const roles =
    user?.roles && user.roles.length > 0 ? user.roles : ['merchant_user'];
  return {
    avatar: user?.logoUrl || '',
    desc: shopName,
    homePath: user?.homePath || '/home',
    realName: userName,
    roles,
    userId: String(user?.app_id || ''),
    username: userName,
    token: '',
    app_id: user?.app_id,
    logoUrl: user?.logoUrl,
    shopName,
    version: user?.version,
  } as UserInfo;
}

function persistMenus(menus: ShopAccessMenuItem[]) {
  if (menus?.length) {
    sessionStorage.setItem(QIXI_SHOP_MENU_KEY, JSON.stringify(menus));
  }
}

type ShopSessionSnapshot = {
  accessCodes: string[];
  menus: ShopAccessMenuItem[];
  userInfo: UserInfo;
};

let sessionInflight: Promise<ShopSessionSnapshot> | null = null;
/** 同页会话内复用，避免 login bootstrap 与 generateAccess 各打一次 session */
let cachedSession: ShopSessionSnapshot | null = null;

export function clearShopSessionCache() {
  cachedSession = null;
  sessionInflight = null;
  try {
    sessionStorage.removeItem(QIXI_SHOP_MENU_KEY);
  } catch {
    // ignore
  }
}

async function fetchShopSessionApiInner() {
  const [profile, menuRes, permissionRes] = await Promise.all([
    requestClient.get<MerchantProfileResponse>('/auth/me'),
    requestClient.get<{ menus: ShopAccessMenuItem[] }>('/auth/menus'),
    requestClient.get<{ permissions: string[] }>('/auth/permissions'),
  ]);
  const data: ShopSessionResponse = {
    codes: permissionRes?.permissions ?? [],
    menus: menuRes?.menus ?? [],
    user: {
      app_id: profile.merchant_admin_id,
      homePath: '/dashboard',
      roles: profile.roles.split(',').filter(Boolean),
      shop_name: profile.mer_name,
      user_name: profile.real_name || profile.account,
    },
  };
  const menus = data.menus;
  persistMenus(menus);
  syncShopPathAuthCache();
  applyTenantAppBranding(data?.user?.shop_name, data?.user?.logoUrl);
  const snapshot: ShopSessionSnapshot = {
    accessCodes: data?.codes ?? [],
    menus,
    userInfo: mapShopSessionUser(
      data?.user ?? {
        homePath: '/home',
        is_super: 0,
        roles: [],
        user_name: '',
      },
    ),
  };
  cachedSession = snapshot;
  return snapshot;
}

/** api-platform：一次返回菜单、权限码、用户信息（并发 + 同页去重） */
export async function fetchShopSessionApi(options?: { force?: boolean }) {
  if (!options?.force && cachedSession) {
    return cachedSession;
  }
  if (sessionInflight) {
    return sessionInflight;
  }
  sessionInflight = fetchShopSessionApiInner().finally(() => {
    sessionInflight = null;
  });
  return sessionInflight;
}
