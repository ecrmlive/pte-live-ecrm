import type { RouteRecordStringComponent } from '@vben/types';

import type { ShopAccessMenuItem } from '#/utils/pte-live-menu';
import { convertShopMenusToVben } from '#/utils/pte-live-menu';
import { QIXI_SHOP_MENU_KEY } from '#/utils/pte-live-api';
import { fetchShopSessionApi } from '#/api/core/shop-session';

export async function getAllMenusApi() {
  const cached = sessionStorage.getItem(QIXI_SHOP_MENU_KEY);
  if (cached) {
    try {
      const menus = JSON.parse(cached) as ShopAccessMenuItem[];
      if (menus?.length) {
        return convertShopMenusToVben(menus);
      }
    } catch {
      // ignore invalid cache
    }
  }
  const session = await fetchShopSessionApi();
  return convertShopMenusToVben(session.menus || []);
}
