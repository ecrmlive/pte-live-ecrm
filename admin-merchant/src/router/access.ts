import type {
  ComponentRecordType,
  GenerateMenuAndRoutesOptions,
} from '@vben/types';

import { generateAccessible } from '@vben/access';
import { preferences } from '@vben/preferences';
import { generateMenus } from '@vben/utils';

import { requestClient } from '#/api/request';
import { BasicLayout, IFrameView } from '#/layouts';
import {
  convertMergersMenusToVben,
  extractRouterRoutes,
} from '#/utils/mergers-menu';

const forbiddenComponent = () => import('#/views/_core/fallback/forbidden.vue');

async function generateAccess(options: GenerateMenuAndRoutesOptions) {
  const pageMap: ComponentRecordType = Object.fromEntries(
    Object.entries(import.meta.glob('../views/**/*.vue')),
  );
  const layoutMap: ComponentRecordType = {
    BasicLayout,
    IFrameView,
  };

  const menusRes = await requestClient.get<{ menus: unknown[] }>('/auth/menus');
  const menuTree = convertMergersMenusToVben(menusRes?.menus || []);
  const routerRoutes = extractRouterRoutes(menuTree);

  const result = await generateAccessible(preferences.app.accessMode, {
    ...options,
    fetchMenuListAsync: async () => routerRoutes,
    forbiddenComponent,
    layoutMap,
    pageMap,
  });

  return {
    ...result,
    accessibleMenus: generateMenus(menuTree, options.router),
  };
}

export { generateAccess };
