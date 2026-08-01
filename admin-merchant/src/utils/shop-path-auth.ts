import type { ShopAccessMenuItem } from './pte-live-menu';
import { QIXI_SHOP_MENU_KEY } from './pte-live-api';

function createdAuth(
  list: ShopAccessMenuItem[],
  obj: Record<string, boolean>,
) {
  if (!list?.length) {
    return;
  }
  for (const item of list) {
    const path = item.path;
    if (path && String(path).trim()) {
      obj[String(path).toLowerCase()] = true;
    }
    if (Array.isArray(item.children) && item.children.length > 0) {
      createdAuth(item.children, obj);
    }
  }
}

function readSessionMenus(): ShopAccessMenuItem[] {
  try {
    const raw = sessionStorage.getItem(QIXI_SHOP_MENU_KEY);
    if (!raw) {
      return [];
    }
    const menus = JSON.parse(raw) as ShopAccessMenuItem[];
    return Array.isArray(menus) ? menus : [];
  } catch {
    return [];
  }
}

/** 从 session 菜单树构建 path → 有权限（对齐 legacy createdAuth / v-auth） */
export function buildShopPathAuthMap(
  menus: ShopAccessMenuItem[] = readSessionMenus(),
): Record<string, boolean> {
  const authlist: Record<string, boolean> = {};
  createdAuth(menus, authlist);
  return authlist;
}

export function hasShopPathAuth(path: string): boolean {
  const normalized = String(path || '').trim().toLowerCase();
  if (!normalized) {
    return false;
  }
  const auth = buildShopPathAuthMap();
  return auth[normalized] === true;
}

/** VbenTableAction hasPermission 适配 shop path 权限 */
export function tableActionHasPermission(auth?: string | string[]): boolean {
  if (!auth) {
    return true;
  }
  const paths = Array.isArray(auth) ? auth : [auth];
  return paths.some((path) => hasShopPathAuth(path));
}

/** 刷新 sessionStorage authlist，供 legacy filters.isAuth / v-auth 使用 */
export function syncShopPathAuthCache() {
  const authlist = buildShopPathAuthMap();
  sessionStorage.setItem('authlist', JSON.stringify(authlist));
  return authlist;
}
